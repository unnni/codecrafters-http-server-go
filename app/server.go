package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"path"
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

func getUrlPath(buffer []byte, byteSize int) (string, string) {
	httpRequest := strings.Split(string(buffer[:byteSize]), "\r\n")
	for _, val := range httpRequest {
		fmt.Println(val)
	}
	httpStatus := httpRequest[0]
	httpPath := strings.Split(httpStatus, " ")
	fmt.Println("---------<><<>")
	fmt.Println(httpPath)
	return httpPath[0], httpPath[1]
}

func readFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error", err)
		return "", errors.New("Not Found")
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	fileSize := fileInfo.Size()
	buffer := make([]byte, fileSize)
	_, err = file.Read(buffer)
	if err != nil {
		fmt.Println("Error", err)
		return "", err
	}
	return string(buffer), nil
}

func readRequestFile(buffer []byte, n int) []byte {
	buffer = buffer[:n]
	lines := bytes.Split(buffer, []byte("\r\n"))
	fmt.Println("readRequestFile--=---------->>")
	fmt.Println(string(lines[len(lines)-1]))
	return lines[len(lines)-1]
}
func handleConnection(conn net.Conn, dir string) {
	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Print("Failed to read contents of HTTP Request", err.Error())
		os.Exit(1)
	}
	httpMethod, httpPath := getUrlPath(buffer, n)
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

	} else if strings.HasPrefix(httpPath, "/files/") {

		if httpMethod == "GET" {
			filePath := strings.TrimPrefix(httpPath, "/files/")
			fileContent, err := readFile(path.Join(dir, filePath))
			if err != nil && err.Error() != "Not Found" {
				return
			} else if err != nil && err.Error() == "Not Found" {
				sendResponse([]byte("HTTP/1.1 404 Not Found\r\n\r\n"), conn)
			}

			responseBuffer := ([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", len(fileContent), fileContent)))
			sendResponse(responseBuffer, conn)
		} else if httpMethod == "POST" {
			fileName := strings.TrimPrefix(httpPath, "/files/")
			file, err := os.Create(path.Join(dir, fileName))
			fileContent := readRequestFile(buffer, n)
			if err != nil {
				return
			}
			defer file.Close()
			_, err = file.Write(fileContent)

			if err != nil {
				fmt.Println("Error writing to file: ", err.Error())
				os.Exit(1)
			}
			responseBuffer := []byte("HTTP/1.1 201 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: 0\r\n\r\n")
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

	var directory string
	flag.StringVar(&directory, "directory", "", "path to file directory")
	flag.Parse()

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
		go handleConnection(connection, directory)
	}
}
