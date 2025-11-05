package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"

	"com.github.redisgo/config"
	"com.github.redisgo/connection"
	"com.github.redisgo/database"
	"com.github.redisgo/util"
)

func main() {
	fmt.Println("Welcome to RedisGo")

	port := flag.String("port", "6379", "The port to listen on")
	replicaof := flag.String("replicaof", "", "The address of the master to replicate from")
	flag.Parse()
	config.InitServerConfig()
	config.ServerConfig.ReplicaOf = strings.Join(strings.Split(*replicaof, " "), ":")
	config.ServerConfig.Role = "master"
	config.ServerConfig.Port = *port
	if *replicaof != "" {
		config.ServerConfig.Role = "slave"
	}
	if config.ServerConfig.Role == "master" {
		config.ServerConfig.Master_replid = util.GenerateID(40)
	}

	fmt.Println("Starting a tcp server on port " + *port)

	ln, err := net.Listen("tcp", "0.0.0.0:"+*port)
	if err != nil {
		fmt.Println("Failed to bind to port " + *port)
		os.Exit(1)
	}

	fmt.Println("Server started on port " + *port)

	database.InitDataBase()

	if config.ServerConfig.Role == "slave" {
		masterConnection, err := net.Dial("tcp", config.ServerConfig.ReplicaOf)
		if err != nil {
			fmt.Println("Failed to connect to master: ", err.Error())
			os.Exit(1)
		}
		err = util.Handshake(&masterConnection)
		if err != nil {
			fmt.Println("Failed to handshake with master: ", err.Error())
			os.Exit(1)
		}
	}

	//test 1
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
