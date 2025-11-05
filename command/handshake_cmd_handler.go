package command

import (
	"com.github.redisgo/config"
	"com.github.redisgo/util"
)

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

func (cmd *Cmd) handlePsyncCommand() string {
	if len(cmd.Args) < 3 {
		return "-ERR wrong number of arguments for 'psync' command\r\n"
	}
	master_repl_offset := cmd.Args[2]
	if master_repl_offset == "-1" {
		master_repl_offset = "0"
	}
	return "+FULLRESYNC " + config.ServerConfig.Master_replid + " " + master_repl_offset + "\r\n"
}
