package command

import (
	"net"
	"strconv"
	"strings"
	"time"

	"com.github.redisgo/config"
	"com.github.redisgo/database"
	"com.github.redisgo/util"
)

type object struct {
	Value           string
	ExpireInMillies int64
	Time            time.Time
}

func (cmd *Cmd) handleSetCommand(conn *net.Conn) {
	if len(cmd.Args) < 3 {
		writeResponse(conn, "-ERR wrong number of arguments for 'set' command\r\n")
		return
	}
	if cmd.checkMultiExists() {
		writeResponse(conn, util.ReturnQueuedResponse())
		return
	}

	lock := database.GetKeyLock(cmd.Args[1])
	defer lock.Unlock()

	key := cmd.Args[1]
	value := cmd.Args[2]

	Object := &object{Value: value, ExpireInMillies: 0, Time: time.Time{}}

	if len(cmd.Args) == 5 {
		timeUnit := strings.ToUpper(cmd.Args[3])
		timeValue, err := strconv.Atoi(cmd.Args[4])
		if err != nil {
			writeResponse(conn, "-ERR invalid expire time\r\n")
			return
		}

		switch timeUnit {
		case "PX":
			Object.ExpireInMillies = int64(timeValue)
		case "EX":
			Object.ExpireInMillies = int64(timeValue) * 1000
		default:
			writeResponse(conn, "-ERR unknown time unit\r\n")
			return
		}

		Object.Time = time.Now().Add(time.Duration(Object.ExpireInMillies) * time.Millisecond)
	}
	lock.Lock()
	database.Store(key, Object)
	if config.ServerConfig.Role == "master" {
		writeResponse(conn, util.ReturnOkResponse())
		go cmd.propagateToSlaves()
	}
}

func (cmd *Cmd) handleGetCommand(conn *net.Conn) {
	if len(cmd.Args) < 2 {
		writeResponse(conn, "-ERR wrong number of arguments for 'get' command\r\n")
		return
	}

	if cmd.checkMultiExists() {
		writeResponse(conn, util.ReturnQueuedResponse())
		return
	}

	lock := database.GetKeyLock(cmd.Args[1])

	key := cmd.Args[1]
	lock.Lock()
	value, ok := database.Get(key)
	if !ok {
		lock.Unlock()
		writeResponse(conn, util.ReturnNullResponse())
		return
	}
	now := time.Now()
	Object := value.(*object)
	if Object.ExpireInMillies > 0 && now.After(Object.Time) {
		database.Delete(key)
		lock.Unlock()
		writeResponse(conn, util.ReturnNullResponse())
		return
	}
	lock.Unlock()
	writeResponse(conn, util.ReturnBulkStringResponse(Object.Value))
}

func (cmd *Cmd) handleIncrCommand(conn *net.Conn) {
	if len(cmd.Args) < 2 {
		writeResponse(conn, "-ERR wrong number of arguments for 'incr' command\r\n")
		return
	}

	if cmd.checkMultiExists() {
		writeResponse(conn, util.ReturnQueuedResponse())
		return
	}

	key := cmd.Args[1]
	lock := database.GetKeyLock(key)

	lock.Lock()
	response, ok := database.Get(key)
	if !ok {
		database.Store(key, &object{Value: "1", ExpireInMillies: 0, Time: time.Time{}})
		lock.Unlock()
		if config.ServerConfig.Role == "master" {
			writeResponse(conn, util.ReturnIntegerResponse(1))
			go cmd.propagateToSlaves()
		}
		return
	}
	value := response.(*object).Value
	now := time.Now()
	if response.(*object).ExpireInMillies > 0 && response.(*object).Time.Before(now) {
		database.Store(key, &object{Value: "1", ExpireInMillies: 0, Time: now.Add(1 * time.Second)})
		lock.Unlock()
		if config.ServerConfig.Role == "master" {
			writeResponse(conn, util.ReturnIntegerResponse(1))
			go cmd.propagateToSlaves()
			return
		}
	}

	valueInt, err := strconv.Atoi(value)
	if err != nil {
		lock.Unlock()
		writeResponse(conn, "-ERR value is not an integer\r\n")
		return
	}
	valueInt++
	database.Store(key, &object{Value: strconv.Itoa(valueInt), ExpireInMillies: 0, Time: time.Time{}})
	lock.Unlock()
	if config.ServerConfig.Role == "master" {
		writeResponse(conn, util.ReturnIntegerResponse(valueInt))
		go cmd.propagateToSlaves()
	}
}
