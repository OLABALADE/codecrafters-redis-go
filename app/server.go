package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	defer l.Close()

	conn, err := l.Accept()
	for {
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		handleClient(conn)
	}

}
func handleClient(conn net.Conn) {
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)

	if err != nil {
        fmt.Println(err)
		return
	}

	fmt.Println(string(buf))
	conn.Write([]byte("+PONG\r\n"))
}
