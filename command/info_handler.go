package command

import (
	"net"

	"com.github.redisgo/util"
)

func (cmd *Cmd) handleInfoCommand(conn *net.Conn) {
	if len(cmd.Args) < 2 {
		writeResponse(conn, "-ERR wrong number of arguments for 'info' command\r\n")
		return
	}

	key := cmd.Args[1]

	switch key {
	case "replication":
		writeResponse(conn, util.ReturnReplicationInfoResponse())
	default:
		writeResponse(conn, util.ReturnErrorResponse("unknown key: "+key))
	}
}
