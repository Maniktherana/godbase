package main

import (
	"fmt"
	"net"

	"github.com/maniktherana/godbase/pkg/resp"
	"github.com/maniktherana/godbase/pkg/writer"
)

func main() {
	fmt.Println("Listening on port :6379")

	// Create a new server
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Listen for connections
	conn, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()

	for {
		resp := resp.NewResp(conn)
		value, err := resp.Read()
		if err != nil {
			fmt.Println(err)
			return
		}

		_ = value

		writer := writer.NewWriter(conn)
		v := resp.Value{Typ: "string", Str: "OK"}
		writer.Write(v)
	}
}
