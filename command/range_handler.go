package command

import (
	"strconv"

	"com.github.redisgo/database"
	"com.github.redisgo/util"
)

func (cmd *Cmd) handleLRangeCommand() string {
	if len(cmd.Args) < 4 {
		return "-ERR wrong number of arguments for 'lrange' command\r\n"
	}

	key := cmd.Args[1]
	start, err1 := strconv.Atoi(cmd.Args[2])
	end, err2 := strconv.Atoi(cmd.Args[3])
	if err1 != nil || err2 != nil {
		return "-ERR invalid start or end\r\n"
	}

	lock := database.GetKeyLock(key)
	defer lock.Unlock()

	lock.Lock()
	ObjectList, ok := database.Get(key)
	if !ok {
		return util.ReturnEmptyArrayResponse()
	}
	ObjectListValue := ObjectList.(*objectList).Value
	objectListLength := len(ObjectListValue)

	if start < 0 {
		start = start + objectListLength
	}
	if end < 0 {
		end = objectListLength + end
	}

	if start < 0 {
		start = 0
	}
	if end < 0 {
		end = 0
	}

	return util.ReturnArrayResponse(ObjectListValue[start : end+1])
}
