package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleClient(conn)
	}

}
func handleClient(conn net.Conn) {
	defer conn.Close()
	for {
		buf := make([]byte, 1024)
		_, err := conn.Read(buf)

		if err != nil {
			fmt.Println(err)
			return
		}

		par_req := parseRequest(buf)
		conn.Write([]byte(par_req))
	}
}

func parseRequest(req []byte) string {
	data := map[string]string{}
	body := strings.Split(string(req), "\r\n")
	cmd := strings.ToLower(body[2])

	switch cmd {
	case "echo":
		mes := body[4]
		return fmt.Sprint("+", mes, "\r\n")

	case "ping":
		return "+PONG\r\n"

	case "set":
		data[body[4]] = body[5]
		return "+OK\r\n"

	case "get":
		item, prs := data[body[4]]
		if !prs {
			return "$-1\r\n"
		}
        return fmt.Sprintf("$%d\r\n%s\r\n", len(item), item)
	}

	return ""
}
