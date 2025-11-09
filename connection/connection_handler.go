package connection

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"net"
	"strconv"

	"com.github.redisgo/command"
)

type ConnectionHandler struct {
	conn *net.Conn
	in   chan *command.Cmd
}

func NewConnectionHandler(conn *net.Conn) *ConnectionHandler {
	return &ConnectionHandler{conn: conn, in: make(chan *command.Cmd)}
}

func (h *ConnectionHandler) Handle() {
	defer (*h.conn).Close()

	go h.read()

	for {
		select {
		case cmd := <-h.in:
			fmt.Printf("[DEBUG] Received command: %s\r\n", cmd)
			cmd.Handle(h.conn)
			//(*h.conn).Write([]byte(response))
			/*if cmd.Name == "PSYNC" && len(cmd.Args) >= 3 && cmd.Args[2] == "-1" {
				h.sendEmptyRDBfileResponse()
			}*/
		}
	}
}

func (h *ConnectionHandler) read() {
	reader := bufio.NewReader(*h.conn)

	for {
		cmd, err := command.ReadCommand(reader)
		if err != nil {
			fmt.Println("Error reading command: ", err.Error())
			break
		}
		h.in <- cmd
	}
}

func (h *ConnectionHandler) sendEmptyRDBfileResponse() {
	// Base64 encoded empty RDB file (Redis version 6.0+)
	// This is a minimal valid RDB file with no data
	sampleBase64EmptyRDBfile := "UkVESVMwMDEx+glyZWRpcy12ZXIFNy4yLjD6CnJlZGlzLWJpdHPAQPoFY3RpbWXCbQi8ZfoIdXNlZC1tZW3CsMQQAPoIYW9mLWJhc2XAAP/wbjv+wP9aog=="

	decodedRDBfile, err := base64.StdEncoding.DecodeString(sampleBase64EmptyRDBfile)
	if err != nil {
		fmt.Println("Error decoding sample base64 empty RDB file: " + err.Error())
		return
	}

	// Send RDB file in bulk string format: $<length>\r\n<rdb_bytes>\r\n
	// First send the length prefix
	lengthPrefix := []byte("$" + strconv.Itoa(len(decodedRDBfile)) + "\r\n")
	_, err = (*h.conn).Write(lengthPrefix)
	if err != nil {
		fmt.Println("Error writing RDB length prefix: " + err.Error())
		return
	}

	// Then send the RDB file bytes directly (not as string)
	_, err = (*h.conn).Write(decodedRDBfile)
	if err != nil {
		fmt.Println("Error writing RDB file bytes: " + err.Error())
		return
	}
	fmt.Println("Sent empty RDB file: " + string(decodedRDBfile))
	// Finally send the CRLF terminator
	_, err = (*h.conn).Write([]byte("\r\n"))
	if err != nil {
		fmt.Println("Error writing RDB terminator: " + err.Error())
		return
	}

	fmt.Printf("Sent empty RDB file (%d bytes) to slave\n", len(decodedRDBfile))
}
