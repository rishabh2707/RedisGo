package command

import "com.github.redisgo/util"

func (cmd *Cmd) handleREPLCONFCommand() string {
	if len(cmd.Args) < 2 {
		return "-ERR wrong number of arguments for 'replconf' command\r\n"
	}
	key := cmd.Args[1]
	switch key {
	case "listening-port":
		return util.ReturnOkResponse()
	case "capa":
		return util.ReturnOkResponse()
	default:
		return util.ReturnErrorResponse("unknown key: " + key)
	}
}
