package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	// Uncomment this block to pass the first stage
)

func sendResponse(response []byte, conn net.Conn) {
	_, err := conn.Write(response)

	if err != nil {
		fmt.Println("Error writing data on connection", err.Error())
	}
}
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

	var requestBuffer []byte
	_, err = connection.Read(requestBuffer)
	if err != nil {
		log.Fatal(err)
	}

	requestBufferLines := strings.Split(string(requestBuffer), "\r\n")
	startLine := requestBufferLines[0]
	components := strings.Split(startLine, " ")
	// method := components[0]
	path := components[1]

	if path == "/" {
		fmt.Printf("Inside the root path")
		sendResponse([]byte("HTTP/1.1 200 OK\r\n\r\n"), connection)
	} else {
		fmt.Printf("Inside the 404 path")
		sendResponse([]byte("HTTP/1.1 404 Not Found\r\n\r\n"), connection)
	}

	defer connection.Close()
}
