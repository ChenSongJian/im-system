package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

func (server *Server) handleConnection(connection net.Conn) {
	fmt.Println("connection created succesfully")
}

func (server *Server) start() {
	// Creates a TCP listener on the specified IP address and port
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.Ip, server.Port))
	if err != nil {
		fmt.Println("net.Listen error:", err)
		return
	}

	defer listener.Close()

	for {
		// Accepting incoming connections
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept error:", err)
			continue
		}
		// For each accepted connection, it spawns a new goroutine (using go) to handle the connection
		go server.handleConnection(connection)
	}

}
