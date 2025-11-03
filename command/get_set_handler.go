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
