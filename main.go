package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"com.github.redisgo/connection"
	"com.github.redisgo/database"
)

func main() {
	fmt.Println("Welcome to RedisGo")

	port := flag.String("port", "6379", "The port to listen on")

	flag.Parse()

	fmt.Println("Starting a tcp server on port " + *port)

	ln, err := net.Listen("tcp", "0.0.0.0:"+*port)
	if err != nil {
		fmt.Println("Failed to bind to port " + *port)
		os.Exit(1)
	}

	fmt.Println("Server started on port " + *port)

	database.InitDataBase()
	//
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
