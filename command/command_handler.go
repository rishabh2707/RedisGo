package command

import "strings"

func (cmnd *Cmd) Handle() string {
	switch strings.ToUpper(cmnd.Name) {
	case "PING":
		return cmnd.handlePingCommand()
	default:
		return "-ERR unknown command '" + cmnd.Name + "'\r\n"
	}
}

func (cmnd *Cmd) handlePingCommand() string {
	return "+PONG\r\n"
}
