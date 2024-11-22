package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

type data struct {
	value   string
	created time.Time
	expire  int64
}

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
	db := &map[string]*data{}
	for {
		buf := make([]byte, 1024)
		_, err := conn.Read(buf)

		if err != nil {
			fmt.Println(err)
			return
		}
		clearExpired(db)
		res, err := handleRequest(buf, db)
		conn.Write(res)
	}
}
