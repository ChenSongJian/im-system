package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	Ip                      string
	Port                    int
	OnlineUserMap           map[string]*User // key: username, value: user info
	MapLock                 sync.RWMutex
	BroadcastMessageChannel chan string
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:                      ip,
		Port:                    port,
		OnlineUserMap:           make(map[string]*User),
		BroadcastMessageChannel: make(chan string),
	}
	return server
}

func (server *Server) Broadcast(user *User, message string) {
	// Format message and pass to server BroadcastMessageChannel
	broadcastMessage := "[" + user.Address + "]" + user.Name + ": " + message
	server.BroadcastMessageChannel <- broadcastMessage
}

func (server *Server) ListenBroadcast() {
	for {
		// get broadcast message from server BroadcastMessageChannel
		message := <-server.BroadcastMessageChannel

		server.MapLock.Lock()
		// pass the message to all the users
		for _, user := range server.OnlineUserMap {
			user.Channel <- message
		}
		server.MapLock.Unlock()
	}
}

func (server *Server) handleConnection(connection net.Conn) {
	// Create new user from the connection
	user := NewUser(connection)

	server.MapLock.Lock()
	// Add new user into OnlineUserMap and store the user info
	server.OnlineUserMap[user.Name] = user
	server.MapLock.Unlock()

	// Broadcast message
	server.Broadcast(user, "Hello World!")

	select {}
}

func (server *Server) start() {
	// Creates a TCP listener on the specified IP address and port
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.Ip, server.Port))
	if err != nil {
		fmt.Println("net.Listen error:", err)
		return
	}

	defer listener.Close()

	// Create goroutine to listen broadcast message
	go server.ListenBroadcast()

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
