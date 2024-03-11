package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Connection net.Conn
	Choice     int
}

func (cilent *Client) menu() bool {
	choice := -1
	fmt.Println("1. Broadcast message")
	fmt.Println("2. Direct message")
	fmt.Println("3. Rename")
	fmt.Println("0. Exit")
	fmt.Scanln(&choice)
	if choice >= 0 && choice <= 3 {
		cilent.Choice = choice
		return true
	} else {
		fmt.Println("Invalid choice!")
		return false
	}
}

func (client *Client) Rename() bool {
	var name string
	fmt.Println("You chose 3. Rename, please enter the new username")
	fmt.Scanln(&name)
	if _, err := client.Connection.Write([]byte("/rename|" + name + "\n")); err != nil {
		fmt.Println("Rename failed:", err)
		return false
	}
	return true
}

func (client *Client) Run() {
	for client.Choice != 0 {
		for !client.menu() {
		}

		switch client.Choice {
		case 1:
			fallthrough // implement later
		case 2:
			fmt.Println("Your choice is", client.Choice)
			break
		case 3:
			client.Rename()
			break
		}
	}
}

func NewClient(ip string, port int) *Client {
	client := &Client{
		ServerIp:   ip,
		ServerPort: port,
		Choice:     -1,
	}

	connection, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return nil
	}
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
	fmt.Println("Client started!")

	go func() {
		io.Copy(os.Stdout, client.Connection)
	}()

	client.Run()
}
