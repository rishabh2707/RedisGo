package command

import (
	"fmt"
	"strconv"

	"com.github.redisgo/database"
	"com.github.redisgo/util"
)

type Queue struct {
	Commands []*Cmd
}

func (cmd *Cmd) handleMultiExecCommand() string {
	if len(cmd.Args) < 1 {
		return "-ERR wrong number of arguments for 'multi' command\r\n"
	}

	lock := database.GetKeyLock("MULTI-" + strconv.FormatUint(util.GetGoroutineID(), 10))
	lock.Lock()
	database.Store("MULTI-"+strconv.FormatUint(util.GetGoroutineID(), 10), &Queue{Commands: []*Cmd{}})
	lock.Unlock()
	return util.ReturnOkResponse()
}

func (cmd *Cmd) handleExecCommand() string {
	if len(cmd.Args) < 1 {
		return "-ERR wrong number of arguments for 'exec' command\r\n"
	}

	lock := database.GetKeyLock("MULTI-" + strconv.FormatUint(util.GetGoroutineID(), 10))

	lock.Lock()
	response, ok := database.Get("MULTI-" + strconv.FormatUint(util.GetGoroutineID(), 10))
	if !ok {
		lock.Unlock()
		return fmt.Sprintf("-ERR EXEC without MULTI\r\n")
	}
	queue := response.(*Queue)
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
	}
	return result
}
