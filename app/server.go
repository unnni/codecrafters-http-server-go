package main

import (
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
	// Uncomment this block to pass the first stage
)

func sendResponse(response []byte, conn net.Conn) {
	_, err := conn.Write(response)

	if err != nil {
		fmt.Println("Error writing data on connection", err.Error())
	}
}

func getUrlPath(buffer []byte, byteSize int) string {
	httpRequest := strings.Split(string(buffer[:byteSize]), "\r\n")
	for _, val := range httpRequest {
		fmt.Println(val)
	}
	httpStatus := httpRequest[0]
	httpPath := strings.Split(httpStatus, " ")
	return httpPath[1]
}

func handleConnection(conn net.Conn) {
	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Print("Failed to read contents of HTTP Request", err.Error())
		os.Exit(1)
	}
	httpPath := getUrlPath(buffer, n)
	if httpPath == "/" {
		sendResponse([]byte("HTTP/1.1 200 OK\r\n\r\n"), conn)
	} else if strings.HasPrefix(httpPath, "/echo/") {
		responseBody := strings.TrimPrefix(httpPath, "/echo/")
		responseBuffer := []byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s\r\n\r\n", len(responseBody), responseBody))
		sendResponse(responseBuffer, conn)
	} else if strings.HasPrefix(httpPath, "/user-agent") {
		pattern := `User-Agent: (.+?)(?:\r\n|$)`
		regex := regexp.MustCompile(pattern)
		req := string(buffer[:])
		matches := regex.FindStringSubmatch(req)
		if len(matches) > 1 {
			content := matches[1]
			responseBuffer := ([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(content), content)))
			sendResponse(responseBuffer, conn)
		}
	} else {
		sendResponse([]byte("HTTP/1.1 404 Not Found\r\n\r\n"), conn)
	}

	defer conn.Close()
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

	defer l.Close()

	for {
		connection, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleConnection(connection)
	}
}
