package command

import (
	"net"
	"strconv"

	"com.github.redisgo/database"
	"com.github.redisgo/util"
)

func (cmd *Cmd) handleLRangeCommand(conn *net.Conn) {
	if len(cmd.Args) < 4 {
		writeResponse(conn, "-ERR wrong number of arguments for 'lrange' command\r\n")
		return
	}
	if cmd.checkMultiExists() {
		writeResponse(conn, util.ReturnQueuedResponse())
		return
	}

	key := cmd.Args[1]
	start, err1 := strconv.Atoi(cmd.Args[2])
	end, err2 := strconv.Atoi(cmd.Args[3])
	if err1 != nil || err2 != nil {
		writeResponse(conn, "-ERR invalid start or end\r\n")
		return
	}

	lock := database.GetKeyLock(key)
	//defer lock.Unlock()

	lock.Lock()
	ObjectList, ok := database.Get(key)
	if !ok {
		lock.Unlock()
		writeResponse(conn, util.ReturnEmptyArrayResponse())
		return
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

	if start > end || start >= objectListLength {
		lock.Unlock()
		writeResponse(conn, util.ReturnEmptyArrayResponse())
		return
	}

	if end >= objectListLength {
		end = objectListLength - 1
	}

	lock.Unlock()
	writeResponse(conn, util.ReturnArrayResponse(ObjectListValue[start:end+1]))
}
