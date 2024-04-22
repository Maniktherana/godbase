package main

import (
	"fmt"
	"github.com/maniktherana/godbase/pkg/Database"
	"github.com/maniktherana/godbase/pkg/aof"
	"github.com/maniktherana/godbase/pkg/handler"
	"github.com/maniktherana/godbase/pkg/resp"
	"github.com/maniktherana/godbase/pkg/writer"
	"net"
	"strings"
)

func handleConnection(conn net.Conn, kv *Database.Kv, aof *aof.Aof) {
	defer conn.Close()
	kv.Clients[conn.RemoteAddr().String()] = conn

	for {
		r := resp.NewResp(conn)
		value, err := r.Read()
		if err != nil {
			fmt.Println(err)
			return
		}

		if value.Typ != "array" {
			fmt.Println("Invalid request, expected array")
			continue
		}

		if len(value.Array) == 0 {
			fmt.Println("Invalid request, expected array length > 0")
			continue
		}

		command := strings.ToUpper(value.Array[0].Bulk)
		args := value.Array[1:]

		writer := writer.NewWriter(conn)

		handler, ok := handler.Handlers[command]
		if !ok {
			fmt.Println("Invalid command: ", command)
			err := writer.Write(resp.Value{Typ: "string", Str: ""})
			if err != nil {
				fmt.Println("Error writing response:", err)
				break
			}
			continue
		}

		if command == "SET" || command == "HSET" {
			err := aof.Write(value)
			if err != nil {
				fmt.Println("Error writing response:", err)
				break
			}
		}

		result := handler(args, kv)
		writer.Write(result)
	}
}

func main() {
	// Create a new server
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}

	kv := Database.NewKv()
	fmt.Println("Listening on port :6379")

	aof, err := aof.NewAof("database.aof")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer aof.Close()

	aof.Read(func(value resp.Value) {
		command := strings.ToUpper(value.Array[0].Bulk)
		args := value.Array[1:]

		handler, ok := handler.Handlers[command]
		if !ok {
			fmt.Println("Invalid command: ", command)
			return
		}

		handler(args, kv)
	})

	defer l.Close()

	for {
		// Listen for connections
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		go handleConnection(conn, kv, aof)
	}
}
