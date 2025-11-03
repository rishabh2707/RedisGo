package main

import (
	"fmt"
	"net"
	"os"

	"com.github.redisgo/connection"
)

func main() {
	fmt.Println("Welcome to RedisGo")
	fmt.Println("Starting a tcp server on port 6379")

	ln, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	fmt.Println("Server started on port 6379")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		handler := connection.NewConnectionHandler(&conn)
		go handler.Handle()
	}
}
