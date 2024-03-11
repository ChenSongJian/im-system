package main

import (
	"bufio"
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

func (client *Client) Broadcast() {
	var message string
	fmt.Println("You chose 1. Broadcast message, please enter your message, enter \"exit\" to exit Broadcast mode")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		message = scanner.Text()
	}
	fmt.Println(message)
	for message != "exit" {
		if len(message) != 0 {
			if _, err := client.Connection.Write([]byte(message + "\n")); err != nil {
				fmt.Println("Broadcast message failed:", err)
			}
		} else {
			fmt.Println("please enter your message, enter \"exit\" to exit Broadcast mode")
		}
		message = ""
		if scanner.Scan() {
			message = scanner.Text()
		}
	}
}

func (client *Client) ListUsers() {
	if _, err := client.Connection.Write([]byte("/online\n")); err != nil {
		fmt.Println("ListUsers failed:", err)
	}

}

func (client *Client) SlideIntoDM() {
	var recipient string
	var message string
	fmt.Println("You chose 3. Rename, please enter recipient username, enter \"exit\" to exit Direct Message mode")
	client.ListUsers()
	fmt.Scanln(&recipient)

	scanner := bufio.NewScanner(os.Stdin)
	if recipient != "exit" {
		fmt.Println("please enter your message, enter \"exit\" to exit Direct Message mode")
		message = ""
		if scanner.Scan() {
			message = scanner.Text()
		}
		for message != "exit" {
			if len(message) != 0 {
				if _, err := client.Connection.Write([]byte("/to|" + recipient + "|" + message + "\n")); err != nil {
					fmt.Println("Direct Message failed:", err)
				} else {
					fmt.Println("please enter your message, enter \"exit\" to exit Direct Message mode")
				}
				message = ""
				if scanner.Scan() {
					message = scanner.Text()
				}
			}
		}
	}
}

func (client *Client) Rename() {
	var name string
	fmt.Println("You chose 3. Rename, please enter the new username")
	fmt.Scanln(&name)
	if _, err := client.Connection.Write([]byte("/rename|" + name + "\n")); err != nil {
		fmt.Println("Rename failed:", err)
	}
}

func (client *Client) Run() {
	for client.Choice != 0 {
		for !client.menu() {
		}

		switch client.Choice {
		case 1:
			client.Broadcast()
		case 2:
			client.SlideIntoDM()
		case 3:
			client.Rename()
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
