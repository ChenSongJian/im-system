package main

import (
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	Connection net.Conn
}

func NewClient(serverIp string, serverPort int) *Client {
	ipPort := fmt.Sprintf("%s:%d", serverIp, serverPort)
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		Name:       ipPort,
	}

	connection, err := net.Dial("tcp", ipPort)
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return nil
	}
	client.Connection = connection
	return client
}

func main() {
	client := NewClient("127.0.0.1", 8888)
	if client == nil {
		fmt.Println("Failed to start new client")
		return
	}

	fmt.Println("Client started!")

	select {}
}
