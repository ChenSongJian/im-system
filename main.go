package main

func main() {
	// create server hosting at localhost port 8888 and start server
	server := NewServer("127.0.0.1", 8888)
	server.start()
}
