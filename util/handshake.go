package util

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func Handshake(connection *net.Conn) error {
	(*connection).Write([]byte("*1\r\n$4\r\nping\r\n"))
	response, err := bufio.NewReader(*connection).ReadString('\n')
	if err != nil {
		fmt.Println("Failed to read response from master: ", err.Error())
		os.Exit(1)
	}
	if response != "+PONG\r\n" {
		fmt.Println("Master is not responding to ping")
		os.Exit(1)
	}
	fmt.Println("Connected to master: ", response)
	return nil
}
