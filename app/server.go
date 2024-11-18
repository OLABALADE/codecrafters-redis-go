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
    var res string
	body := strings.Split(string(req), "\r\n")
	cmd := strings.ToLower(body[2])

	switch cmd {
	case "echo":
		mes := body[4]
		res = fmt.Sprint("+", mes, "\r\n")
		return res

    case "ping":
        res = "+PONG\r\n"
        return res
	}

	return ""
}
