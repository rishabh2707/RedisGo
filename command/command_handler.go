package command

import (
	"strings"

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
	default:
		return "-ERR unknown command '" + cmnd.Name + "'\r\n"
	}
}

func (cmnd *Cmd) handlePingCommand() string {
	return "+PONG\r\n"
}

func (cmnd *Cmd) handleEchoCommand() string {
	return util.ParseNormalResponse(cmnd.Args[1])
}
