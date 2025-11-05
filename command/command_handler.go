package command

import (
	"strconv"
	"strings"

	"com.github.redisgo/database"
	"com.github.redisgo/util"
)

func (cmnd *Cmd) Handle() string {
	switch strings.ToUpper(cmnd.Name) {
	case "PING":
		return cmnd.handlePingCommand()
	case "ECHO":
		return cmnd.handleEchoCommand()
	case "SET":
		return cmnd.handleSetCommand()
	case "GET":
		return cmnd.handleGetCommand()
	case "RPUSH":
		return cmnd.handleRPushCommand()
	case "LPUSH":
		return cmnd.handleLPushCommand()
	case "LRANGE":
		return cmnd.handleLRangeCommand()
	case "LLEN":
		return cmnd.handleLLenCommand()
	case "LPOP":
		return cmnd.handleLPopCommand()
	case "BLPOP":
		return cmnd.handleBLPopCommand()
	case "TYPE":
		return cmnd.handleTypeCommand()
	case "XADD":
		return cmnd.handleXAddCommand()
	case "INCR":
		return cmnd.handleIncrCommand()
	case "MULTI":
		return cmnd.handleMultiExecCommand()
	case "EXEC":
		return cmnd.handleExecCommand()
	case "DISCARD":
		return cmnd.handleDiscardCommand()
	case "INFO":
		return cmnd.handleInfoCommand()
	case "REPLCONF":
		return cmnd.handleREPLCONFCommand()
	case "PSYNC":
		return cmnd.handlePsyncCommand()
	default:
		return "-ERR unknown command '" + cmnd.Name + "'\r\n"
	}
}

func (cmnd *Cmd) handlePingCommand() string {
	if cmnd.checkMultiExists() {
		return util.ReturnQueuedResponse()
	}
	return "+PONG\r\n"
}

func (cmnd *Cmd) handleEchoCommand() string {
	if cmnd.checkMultiExists() {
		return util.ReturnQueuedResponse()
	}
	return util.ReturnBulkStringResponse(cmnd.Args[1])
}

func (cmnd *Cmd) handleTypeCommand() string {
	if len(cmnd.Args) < 2 {
		return "-ERR wrong number of arguments for 'type' command\r\n"
	}
	if cmnd.checkMultiExists() {
		return util.ReturnQueuedResponse()
	}
	key := cmnd.Args[1]

	value, ok := database.Get(key)
	if !ok {
		return util.ReturnBulkStringResponse("none")
	}

	switch value.(type) {
	case *object:
		return util.ReturnBulkStringResponse("string")
	case *objectList:
		return util.ReturnBulkStringResponse("list")
	case *streamList:
		return util.ReturnBulkStringResponse("stream")
	default:
		return util.ReturnBulkStringResponse("none")
	}
}

func (cmnd *Cmd) checkMultiExists() bool {
	multiLock := database.GetKeyLock("MULTI-" + strconv.FormatUint(util.GetGoroutineID(), 10))
	multiLock.Lock()
	response, ok := database.Get("MULTI-" + strconv.FormatUint(util.GetGoroutineID(), 10))
	if ok {
		queue := response.(*Queue)
		queue.Commands = append(queue.Commands, cmnd)
		multiLock.Unlock()
		return true
	}
	multiLock.Unlock()
	return false
}

func (cmnd *Cmd) handleDiscardCommand() string {
	if cmnd.checkMultiExists() {
		multiLock := database.GetKeyLock("MULTI-" + strconv.FormatUint(util.GetGoroutineID(), 10))
		multiLock.Lock()
		database.Delete("MULTI")
		multiLock.Unlock()
		return util.ReturnOkResponse()
	}
	return util.ReturnErrorResponse("DISCARD without MULTI")
}
