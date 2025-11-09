package command

import (
	"encoding/base64"
	"fmt"
	"net"
	"strconv"

	"com.github.redisgo/config"
	"com.github.redisgo/util"
)

func (cmd *Cmd) handleREPLCONFCommand(conn *net.Conn) {
	if len(cmd.Args) < 2 {
		writeResponse(conn, "-ERR wrong number of arguments for 'replconf' command\r\n")
		return
	}
	key := cmd.Args[1]
	switch key {
	case "listening-port":
		writeResponse(conn, util.ReturnOkResponse())
	case "capa":
		writeResponse(conn, util.ReturnOkResponse())
	default:
		writeResponse(conn, util.ReturnErrorResponse("unknown key: "+key))
	}
}

func (cmd *Cmd) handlePsyncCommand(conn *net.Conn) {
	if len(cmd.Args) < 3 {
		writeResponse(conn, "-ERR wrong number of arguments for 'psync' command\r\n")
		return
	}
	master_repl_offset := cmd.Args[2]
	if master_repl_offset == "-1" {
		master_repl_offset = "0"
	}
	writeResponse(conn, "+FULLRESYNC "+config.ServerConfig.Master_replid+" "+master_repl_offset+"\r\n")
	// Send empty RDB file after FULLRESYNC response
	if cmd.Args[2] == "-1" {
		sendEmptyRDBfile(conn)
	}
}

func sendEmptyRDBfile(conn *net.Conn) {
	// Base64 encoded empty RDB file (Redis version 6.0+)
	// This is a minimal valid RDB file with no data
	const base64EmptyRDBFile = "UkVESVMwMDEx+glyZWRpcy12ZXIFNy4yLjD6CnJlZGlzLWJpdHPAQPoFY3RpbWXCbQi8ZfoIdXNlZC1tZW3CsMQQAPoIYW9mLWJhc2XAAP/wbjv+wP9aog=="

	decodedRDBfile, err := base64.StdEncoding.DecodeString(base64EmptyRDBFile)
	if err != nil {
		fmt.Println("Error decoding empty RDB file: " + err.Error())
		return
	}

	// Send RDB file in bulk string format: $<length>\r\n<rdb_bytes>\r\n
	lengthPrefix := []byte("$" + strconv.Itoa(len(decodedRDBfile)) + "\r\n")
	if _, err := (*conn).Write(lengthPrefix); err != nil {
		fmt.Println("Error writing RDB length prefix: " + err.Error())
		return
	}

	// Send the RDB file bytes directly (not as string)
	if _, err := (*conn).Write(decodedRDBfile); err != nil {
		fmt.Println("Error writing RDB file bytes: " + err.Error())
		return
	}

	// Finally send the CRLF terminator
	if _, err := (*conn).Write([]byte("\r\n")); err != nil {
		fmt.Println("Error writing RDB terminator: " + err.Error())
		return
	}

	fmt.Printf("Sent empty RDB file (%d bytes) to slave\n", len(decodedRDBfile))
}
