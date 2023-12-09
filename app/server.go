package main

import (
	"fmt"
	"net"
	"os"
	// Uncomment this block to pass the first stage
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	connection, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	// writeBuffer := []byte("HTTP/1.1 200 OK\r\n\r\n")
	// _, err = connection.Write(writeBuffer)
	if err != nil {
		fmt.Println("Error sending response", err.Error())
		os.Exit(1)
	}

	var requestBuffer []byte
	connection.Read(requestBuffer)
	fmt.Printf(string(requestBuffer))
	defer connection.Close()
}
