package main

import (
	"net"
	"strings"
)

type User struct {
	Name       string
	Address    string
	Channel    chan string
	Connection net.Conn
	Server     *Server // using pointer to use same lock across all users
}

func NewUser(connection net.Conn, server *Server) *User {
	// Using address as username in v2
	user_address := connection.RemoteAddr().String()
	user := &User{
		Name:       user_address,
		Address:    user_address,
		Channel:    make(chan string),
		Connection: connection,
		Server:     server,
	}

	// start goroutine to listen incoming message
	go user.listenMessage()

	return user
}

func (user *User) online() {

	user.Server.MapLock.Lock()
	// Add new user into OnlineUserMap and store the user info
	user.Server.OnlineUserMap[user.Name] = user
	user.Server.MapLock.Unlock()

	// Broadcast message
	user.Server.Broadcast(user, "Hello World!")

}

func (user *User) offline(idle bool) {
	user.Server.MapLock.Lock()
	// Removes user from map when user going offline
	delete(user.Server.OnlineUserMap, user.Name)
	user.Server.MapLock.Unlock()
	if idle {
		user.Channel <- "User idle for too long, connection closed"
	} else {
		user.Server.Broadcast(user, "See yall next time!")
	}
}

func (user *User) processMessage(message string) {
	if message == "/online" {
		user.Server.MapLock.Lock()
		// Add new user into OnlineUserMap and store the user info
		for _, value := range user.Server.OnlineUserMap {
			onlineUserMessage := "[" + value.Address + "]" + value.Name + ": is online"
			user.Channel <- onlineUserMessage
		}
		user.Server.MapLock.Unlock()
	} else if len(message) > 7 && message[:7] == "/rename" {
		splitMessage := strings.Split(message, "|")
		if len(splitMessage) != 2 {
			user.Channel <- "Invalid format, try \"/rename|newName\" format"
			return
		}
		newName := strings.TrimSpace(splitMessage[1])
		user.Server.MapLock.Lock()
		if _, ok := user.Server.OnlineUserMap[newName]; ok {
			user.Channel <- "Username is in use, please try another name"
		} else {
			delete(user.Server.OnlineUserMap, user.Name)
			user.Server.OnlineUserMap[newName] = user
			user.Name = newName
			user.Channel <- "Username updated, you are now " + newName
		}
		user.Server.MapLock.Unlock()
	} else if len(message) > 3 && message[:3] == "/to" {
		splitMessage := strings.Split(message, "|")
		if len(splitMessage) != 3 {
			user.Channel <- "Invalid format, try \"/to|userA|messageABC\" format"
			return
		}
		recipient := strings.TrimSpace(splitMessage[1])
		directMessage := strings.TrimSpace(splitMessage[2])
		if len(directMessage) == 0 {
			user.Channel <- "Failed to send direct message: message is empty"
		}
		user.Server.MapLock.Lock()
		if recipientInfo, ok := user.Server.OnlineUserMap[recipient]; !ok {
			user.Channel <- "Failed to send direct message: recipient not exists"
		} else {
			recipientInfo.Channel <- "(DM)" + user.Name + ":" + directMessage
		}
		user.Server.MapLock.Unlock()
	} else {
		user.Server.Broadcast(user, message)
	}
}

func (user *User) listenMessage() {
	for {
		// Waiting for new message from user channel and write to user connection
		message := <-user.Channel
		user.Connection.Write([]byte(message + "\n"))
	}
}
