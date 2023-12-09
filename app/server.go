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

	connection, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	defer connection.Close()
	_, err = connection.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	writeBuffer := []byte("HTTP/1.1 200 OK\r\n\r\n")
	_, err = connection.Write([]byte(writeBuffer))
	if err != nil {
		fmt.Println("Error sending response", err.Error())
		os.Exit(1)
	}
}
