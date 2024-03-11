package main

import (
	"flag"
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	Connection net.Conn
}

func NewClient(ip string, port int) *Client {
	client := &Client{
		ServerIp:   ip,
		ServerPort: port,
	}

	connection, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return nil
	}
	client.Name = connection.LocalAddr().String()
	client.Connection = connection
	return client
}

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "Server IP, default 127.0.0.1")
	flag.IntVar(&serverPort, "port", 8888, "Server Port, default 8888")
}

func main() {
	flag.Parse()
	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println("Failed to start new client")
		return
	}

	fmt.Println("Client started! client name:", client.Name)

	select {}
}
