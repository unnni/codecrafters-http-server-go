package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	// Uncomment this block to pass the first stage
)

func sendResponse(response []byte, conn net.Conn) {
	writeBuffer := []byte(response)
	_, err := conn.Write(writeBuffer)

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
	connection.Read(requestBuffer)
	requestBufferString := string(requestBuffer)
	requestBufferLines := strings.Split(requestBufferString, "\r\n")

	if len(requestBufferLines) > 0 {
		firstLine := requestBufferLines[0]
		firstLineEl := strings.Split(firstLine, " ")
		if len(firstLineEl) > 1 {
			path := firstLineEl[1]
			if path == "/" {
				sendResponse([]byte("HTTP/1.1 200 OK\r\n\r\n"), connection)
			} else {
				sendResponse([]byte("HTTP/1.1 404 Not Found\r\n\r\n"), connection)
			}

		}
	}
	defer connection.Close()
}
