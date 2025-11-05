package command

import "com.github.redisgo/util"

func (cmd *Cmd) handleInfoCommand() string {
	if len(cmd.Args) < 2 {
		return "-ERR wrong number of arguments for 'info' command\r\n"
	}

	key := cmd.Args[1]

	switch key {
	case "replication":
		return util.ReturnReplicationInfoResponse()
	default:
		return util.ReturnErrorResponse("unknown key: " + key)
	}

}
