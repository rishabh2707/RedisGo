package connection

import (
	"bufio"
	"fmt"
	"net"

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
			switch cmd.Name {
			case "PING":
				(*h.conn).Write([]byte("+PONG\r\n"))
			default:
				(*h.conn).Write([]byte("-ERR unknown command '" + cmd.Name + "'\r\n"))
			}
		default:
			fmt.Println("[DEBUG] Quit signal received")
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
