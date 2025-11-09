package util

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"com.github.redisgo/config"
)

func Handshake(connection *net.Conn) error {
	(*connection).Write([]byte("*1\r\n$4\r\nping\r\n"))
	response, err := bufio.NewReader(*connection).ReadString('\n')
	if err != nil {
		fmt.Println("Failed to read response from server: " + (*connection).RemoteAddr().String() + ", " + err.Error())
		os.Exit(1)
	}
	if response != "+PONG\r\n" {
		fmt.Println("Server: " + (*connection).RemoteAddr().String() + " is not responding to ping")
		os.Exit(1)
	}
	fmt.Println("Got the response from server: " + (*connection).RemoteAddr().String() + ", " + response)

	(*connection).Write([]byte("*3\r\n$8\r\nREPLCONF\r\n$14\r\nlistening-port\r\n$4\r\n" + config.ServerConfig.Port + "\r\n"))
	response, err = bufio.NewReader(*connection).ReadString('\n')
	if err != nil {
		fmt.Println("Failed to read response from server: " + (*connection).RemoteAddr().String() + ", " + err.Error())
		os.Exit(1)
	}
	if response != "+OK\r\n" {
		fmt.Println("Server: " + (*connection).RemoteAddr().String() + " is not responding to REPLCONF")
		os.Exit(1)
	}
	fmt.Println("Got the response from server: " + (*connection).RemoteAddr().String() + ", " + response)

	(*connection).Write([]byte("*3\r\n$8\r\nREPLCONF\r\n$4\r\ncapa\r\n$6\r\npsync2\r\n"))
	response, err = bufio.NewReader(*connection).ReadString('\n')
	if err != nil {
		fmt.Println("Failed to read response from server: " + (*connection).RemoteAddr().String() + ", " + err.Error())
		os.Exit(1)
	}
	if response != "+OK\r\n" {
		fmt.Println("Server: " + (*connection).RemoteAddr().String() + " is not responding to REPLCONF")
		os.Exit(1)
	}
	fmt.Println("Got the response from server: " + (*connection).RemoteAddr().String() + ", " + response)
	(*connection).Write([]byte("*3\r\n$5\r\nPSYNC\r\n$1\r\n?\r\n$2\r\n-1\r\n"))
	response, err = bufio.NewReader(*connection).ReadString('\n')
	if err != nil {
		fmt.Println("Failed to read response from server: " + (*connection).RemoteAddr().String() + ", " + err.Error())
		os.Exit(1)
	}
	fmt.Println("Got the response from server: " + (*connection).RemoteAddr().String() + ", " + response)
	response, err = bufio.NewReader(*connection).ReadString('\n')
	if err != nil {
		fmt.Println("Failed to read response from server: " + (*connection).RemoteAddr().String() + ", " + err.Error())
		os.Exit(1)
	}
	fmt.Println("Got the response from server: " + (*connection).RemoteAddr().String() + ", " + response)
	response, err = bufio.NewReader(*connection).ReadString('\n')
	if err != nil {
		fmt.Println("Failed to read response from server: " + (*connection).RemoteAddr().String() + ", " + err.Error())
		os.Exit(1)
	}
	fmt.Println("Got the response from server: " + (*connection).RemoteAddr().String() + ", " + response)
	response, err = bufio.NewReader(*connection).ReadString('\n')
	if err != nil {
		fmt.Println("Failed to read response from server: " + (*connection).RemoteAddr().String() + ", " + err.Error())
		os.Exit(1)
	}
	fmt.Println("Got the response from server: " + (*connection).RemoteAddr().String() + ", " + response)
	response, err = bufio.NewReader(*connection).ReadString('\n')
	if err != nil {
		fmt.Println("Failed to read response from server: " + (*connection).RemoteAddr().String() + ", " + err.Error())
		os.Exit(1)
	}
	fmt.Println("Got the response from server: " + (*connection).RemoteAddr().String() + ", " + response)
	return nil
}
