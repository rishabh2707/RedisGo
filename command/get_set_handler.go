package command

import (
	"strconv"
	"strings"
	"time"

	"com.github.redisgo/database"
	"com.github.redisgo/util"
)

type object struct {
	Value           string
	ExpireInMillies int64
	Time            time.Time
}

func (cmd *Cmd) handleSetCommand() string {
	if len(cmd.Args) < 3 {
		return "-ERR wrong number of arguments for 'set' command\r\n"
	}
	if cmd.checkMultiExists() {
		return util.ReturnQueuedResponse()
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
			return "-ERR invalid expire time\r\n"
		}

		switch timeUnit {
		case "PX":
			Object.ExpireInMillies = int64(timeValue)
		case "EX":
			Object.ExpireInMillies = int64(timeValue) * 1000
		default:
			return "-ERR unknown time unit\r\n"
		}

		Object.Time = time.Now().Add(time.Duration(Object.ExpireInMillies) * time.Millisecond)
	}
	lock.Lock()
	database.Store(key, Object)
	return util.ReturnOkResponse()
}

func (cmd *Cmd) handleGetCommand() string {
	if len(cmd.Args) < 2 {
		return "-ERR wrong number of arguments for 'get' command\r\n"
	}

	if cmd.checkMultiExists() {
		return util.ReturnQueuedResponse()
	}

	lock := database.GetKeyLock(cmd.Args[1])
	defer lock.Unlock()

	key := cmd.Args[1]
	lock.Lock()
	value, ok := database.Get(key)
	if !ok {
		return util.ReturnNullResponse()
	}
	now := time.Now()
	Object := value.(*object)
	if Object.ExpireInMillies > 0 && now.After(Object.Time) {
		database.Delete(key)
		return util.ReturnNullResponse()
	}
	return util.ParseNormalResponse(Object.Value)
}

func (cmd *Cmd) handleIncrCommand() string {
	if len(cmd.Args) < 2 {
		return "-ERR wrong number of arguments for 'incr' command\r\n"
	}

	if cmd.checkMultiExists() {
		return util.ReturnQueuedResponse()
	}

	key := cmd.Args[1]
	lock := database.GetKeyLock(key)
	defer lock.Unlock()

	lock.Lock()
	response, ok := database.Get(key)
	if !ok {
		database.Store(key, &object{Value: "1", ExpireInMillies: 0, Time: time.Time{}})
		return util.ReturnIntegerResponse(1)
	}
	value := response.(*object).Value
	now := time.Now()
	if response.(*object).ExpireInMillies > 0 && response.(*object).Time.Before(now) {
		database.Store(key, &object{Value: "1", ExpireInMillies: 0, Time: now.Add(1 * time.Second)})
		return util.ReturnIntegerResponse(1)
	}

	valueInt, err := strconv.Atoi(value)
	if err != nil {
		return "-ERR value is not an integer\r\n"
	}
	valueInt++
	database.Store(key, &object{Value: strconv.Itoa(valueInt), ExpireInMillies: 0, Time: time.Time{}})
	return util.ReturnIntegerResponse(valueInt)
}
