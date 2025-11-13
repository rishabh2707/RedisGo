package command

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"com.github.redisgo/config"
	"com.github.redisgo/database"
	"com.github.redisgo/replication"
	"com.github.redisgo/util"
)

func (cmnd *Cmd) Handle(conn *net.Conn) {
	switch strings.ToUpper(cmnd.Name) {
	case "PING":
		cmnd.handlePingCommand(conn)
	case "ECHO":
		cmnd.handleEchoCommand(conn)
	case "SET":
		cmnd.handleSetCommand(conn)
	case "GET":
		cmnd.handleGetCommand(conn)
	case "RPUSH":
		cmnd.handleRPushCommand(conn)
	case "LPUSH":
		cmnd.handleLPushCommand(conn)
	case "LRANGE":
		cmnd.handleLRangeCommand(conn)
	case "LLEN":
		cmnd.handleLLenCommand(conn)
	case "LPOP":
		cmnd.handleLPopCommand(conn)
	case "BLPOP":
		cmnd.handleBLPopCommand(conn)
	case "TYPE":
		cmnd.handleTypeCommand(conn)
	case "XADD":
		cmnd.handleXAddCommand(conn)
	case "INCR":
		cmnd.handleIncrCommand(conn)
	case "MULTI":
		cmnd.handleMultiExecCommand(conn)
	case "EXEC":
		cmnd.handleExecCommand(conn)
	case "DISCARD":
		cmnd.handleDiscardCommand(conn)
	case "INFO":
		cmnd.handleInfoCommand(conn)
	case "REPLCONF":
		cmnd.handleREPLCONFCommand(conn)
	case "PSYNC":
		cmnd.handlePsyncCommand(conn)
	case "CONFIG":
		cmnd.handleGetConfigCommand(conn)
	default:
		writeResponse(conn, "-ERR unknown command '"+cmnd.Name+"'\r\n")
	}
}

func (cmnd *Cmd) handlePingCommand(conn *net.Conn) {
	if cmnd.checkMultiExists() {
		writeResponse(conn, util.ReturnQueuedResponse())
		return
	}
	writeResponse(conn, "+PONG\r\n")
	if config.ServerConfig.Role == "slave" {
		replication.SetReplicaOffset(replication.GetReplicaOffset() + int64(len(cmnd.ToRespFormat())))
	}
}

func (cmnd *Cmd) handleEchoCommand(conn *net.Conn) {
	if cmnd.checkMultiExists() {
		writeResponse(conn, util.ReturnQueuedResponse())
		return
	}
	if len(cmnd.Args) < 2 {
		writeResponse(conn, "-ERR wrong number of arguments for 'echo' command\r\n")
		return
	}
	writeResponse(conn, util.ReturnBulkStringResponse(cmnd.Args[1]))
}

func (cmnd *Cmd) handleTypeCommand(conn *net.Conn) {
	if len(cmnd.Args) < 2 {
		writeResponse(conn, "-ERR wrong number of arguments for 'type' command\r\n")
		return
	}
	if cmnd.checkMultiExists() {
		writeResponse(conn, util.ReturnQueuedResponse())
		return
	}
	key := cmnd.Args[1]

	value, ok := database.Get(key)
	if !ok {
		writeResponse(conn, util.ReturnBulkStringResponse("none"))
		return
	}

	switch value.(type) {
	case *object:
		writeResponse(conn, util.ReturnBulkStringResponse("string"))
	case *objectList:
		writeResponse(conn, util.ReturnBulkStringResponse("list"))
	case *streamList:
		writeResponse(conn, util.ReturnBulkStringResponse("stream"))
	default:
		writeResponse(conn, util.ReturnBulkStringResponse("none"))
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

func (cmnd *Cmd) handleDiscardCommand(conn *net.Conn) {
	if cmnd.checkMultiExists() {
		multiLock := database.GetKeyLock("MULTI-" + strconv.FormatUint(util.GetGoroutineID(), 10))
		multiLock.Lock()
		database.Delete("MULTI-" + strconv.FormatUint(util.GetGoroutineID(), 10))
		multiLock.Unlock()
		writeResponse(conn, util.ReturnOkResponse())
		return
	}
	writeResponse(conn, util.ReturnErrorResponse("DISCARD without MULTI"))
}

func writeResponse(conn *net.Conn, response string) {
	(*conn).Write([]byte(response))
}

func (cmnd *Cmd) propagateToSlaves() {
	if config.ServerConfig.Role == "master" {
		for _, slaveConn := range replication.GetSlaveConnections() {
			request := cmnd.ToRespFormat()
			fmt.Println("Propagating to slave: " + slaveConn.RemoteAddr().String() + ", " + request)
			slaveConn.Write([]byte(request))
		}
	}
}
