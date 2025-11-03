package command

import (
	"strconv"
	"strings"
)

func (cmnd *Cmd) Handle() string {
	switch strings.ToUpper(cmnd.Name) {
	case "PING":
		return cmnd.handlePingCommand()
	case "ECHO":
		return cmnd.handleEchoCommand()
	default:
		return "-ERR unknown command '" + cmnd.Name + "'\r\n"
	}
}

func (cmnd *Cmd) handlePingCommand() string {
	return "+PONG\r\n"
}

func (cmnd *Cmd) handleEchoCommand() string {
	return "$" + strconv.Itoa(len(cmnd.Args[1])) + "\r\n" + cmnd.Args[1] + "\r\n"
}
