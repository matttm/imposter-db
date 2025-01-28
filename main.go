package main

import (
	"fmt"
	"github.com/go-mysql-org/go-mysql/server"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	InitializeProxy()
	proxy, err := net.Listen("tcp", ":3307")
	if err != nil {
		log.Fatalf("failed to start proxy: %s", err.Error())
	}
	fmt.Printf("Listening on localhost%d", 3307)
	for {
		conn, err := proxy.Accept()

		log.Printf("new connection: %s", conn.RemoteAddr())
		if err != nil {
			log.Fatalf("failed to accept connection: %s", err.Error())
		}

	}
}

const COM_QUERY = byte(0x03)
