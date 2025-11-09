package command

import (
	"net"

	"com.github.redisgo/config"
	"com.github.redisgo/util"
)

func (cmd *Cmd) handleGetConfigCommand(conn *net.Conn) {
	if len(cmd.Args) < 2 {
		writeResponse(conn, "-ERR wrong number of arguments for 'config' command\r\n")
		return
	}
	if cmd.checkMultiExists() {
		writeResponse(conn, util.ReturnQueuedResponse())
		return
	}

	if len(cmd.Args) > 2 {
		key := cmd.Args[2]
		switch key {
		case "dir":
			writeResponse(conn, util.ReturnBulkStringResponse(config.RDBConfig.Dir))
		case "dbfilename":
			writeResponse(conn, util.ReturnBulkStringResponse(config.RDBConfig.DbFileName))
		default:
			writeResponse(conn, util.ReturnErrorResponse("unknown key: "+key))
		}
		return
	}

	writeResponse(conn, util.ReturnArrayResponse([]string{"dir", config.RDBConfig.Dir, "dbfilename", config.RDBConfig.DbFileName}))
}
