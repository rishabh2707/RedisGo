package command

import (
	"net"
	"strconv"

	"com.github.redisgo/database"
	"com.github.redisgo/util"
)

type Queue struct {
	Commands []*Cmd
}

func (cmd *Cmd) handleMultiExecCommand(conn *net.Conn) {
	if len(cmd.Args) < 1 {
		writeResponse(conn, "-ERR wrong number of arguments for 'multi' command\r\n")
		return
	}

	lock := database.GetKeyLock("MULTI-" + strconv.FormatUint(util.GetGoroutineID(), 10))
	lock.Lock()
	database.Store("MULTI-"+strconv.FormatUint(util.GetGoroutineID(), 10), &Queue{Commands: []*Cmd{}})
	lock.Unlock()
	writeResponse(conn, util.ReturnOkResponse())
}

func (cmd *Cmd) handleExecCommand(conn *net.Conn) {
	if len(cmd.Args) < 1 {
		writeResponse(conn, "-ERR wrong number of arguments for 'exec' command\r\n")
		return
	}

	lock := database.GetKeyLock("MULTI-" + strconv.FormatUint(util.GetGoroutineID(), 10))

	lock.Lock()
	_, ok := database.Get("MULTI-" + strconv.FormatUint(util.GetGoroutineID(), 10))
	if !ok {
		lock.Unlock()
		writeResponse(conn, "-ERR EXEC without MULTI\r\n")
		return
	}
	/*queue := response.(*Queue)
	responses := make([]string, 0)
	database.Delete("MULTI-" + strconv.FormatUint(util.GetGoroutineID(), 10))
	lock.Unlock()
	for _, queuedCmd := range queue.Commands {
		response := queuedCmd.Handle()
		responses = append(responses, response)
	}
	result := "*" + strconv.Itoa(len(responses)) + "\r\n"
	for _, response := range responses {
		result += response
	}*/
	writeResponse(conn, util.ReturnEmptyArrayResponse())
}
