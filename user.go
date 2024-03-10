package main

import (
	"net"
)

type User struct {
	Name       string
	Address    string
	Channel    chan string
	Connection net.Conn
}

func NewUser(connection net.Conn) *User {
	// Using address as username in v2
	user_address := connection.RemoteAddr().String()
	user := &User{
		Name:       user_address,
		Address:    user_address,
		Channel:    make(chan string),
		Connection: connection,
	}

	// start goroutine to listen incoming message
	go user.listenMessage()

	return user
}

func (user *User) listenMessage() {
	for {
		// Waiting for new message from user channel and write to user connection
		message := <-user.Channel
		user.Connection.Write([]byte(message + "\n"))
	}
}
