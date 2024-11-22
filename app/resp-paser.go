package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	ARRAY         = '*'
	SIMPLE_STRING = '+'
	BULK_STRING   = '$'
	GET           = "GET"
	ECHO          = "ECHO"
	PING          = "PING"
	SET           = "SET"
)

func parse_req(req []byte) []string {
	switch req[0] {
	case ARRAY:
		return readArray(req)
	// case BULK_STRING:
	// 	return readBulkString(req)
	// case SIMPLE_STRING:
	// 	return readSimpleString(req)
	default:
		return []string{}
	}
}

func readArray(req []byte) []string {
	arr := strings.Split(string(req), "\r\n")
	cmd := &[]string{}
	for index, t := range arr {
		switch strings.ToLower(t) {
		case "echo", "get":
			*cmd = append(*cmd, t, arr[index+2])
		case "set":
			*cmd = append(*cmd, t, arr[index+2], arr[index+4])
		case "px":
			*cmd = append(*cmd, t, arr[index+2])
        case "ping":
            *cmd = append(*cmd, PING)
		default:
			continue

		}
	}
	return *cmd
}

/* func readBulkString(req []byte) []string
func readSimpleString(req []byte) []string */

func handleRequest(req []byte, db *map[string]*data) ([]byte, error) {
	dm := *db
	cmd := parse_req(req)
	switch cmd[0] {
	case ECHO:
		res := fmt.Sprint("+", cmd[1], "\r\n")
		return []byte(res), nil

	case PING:
		return []byte("+PONG\r\n"), nil

	case SET:
		d := &data{
			value:   cmd[2],
			created: time.Now(),
			expire:  0,
		}

		if len(cmd) > 3 {
			ex, err := strconv.ParseInt(cmd[4], 10, 64)
			if err != nil {
				return nil, err
			}
			d.expire = ex
		}

		dm[cmd[1]] = d
		return []byte("+OK\r\n"), nil

	case GET:
		item, prs := dm[cmd[1]]
		if !prs {
			return []byte("$-1\r\n"), nil
		}
		req := fmt.Sprintf("$%d\r\n%s\r\n", len(item.value), item.value)
		return []byte(req), nil

	default:
		return []byte{}, nil
	}
}

func clearExpired(db *map[string]*data) {
	df := *db
	for i := range df {
		item := df[i]
		if item.expire != 0 {
			now := time.Now()
			p := now.Sub(item.created).Milliseconds()
			fmt.Println(p)
			if p > item.expire {
				fmt.Println("Hello")
				delete(df, i)
			}
		}
	}
}
