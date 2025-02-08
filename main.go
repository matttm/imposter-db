package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

func main() {
	socket, err := net.Listen("tcp", "127.0.0.1:3307")
	if err != nil {
		log.Fatalf("failed to start proxy: %s", err.Error())
	}
	fmt.Printf("Listening on localhost:%d\n", 3307)
	for {
		originSocket, err := socket.Accept()
		go func(c net.Conn) {
			p := InitializeProxy(originSocket)

			log.Printf("new connection: %s\n", originSocket.RemoteAddr())
			if err != nil {
				log.Fatalf("failed to accept connection: %s", err.Error())
			}
			defer p.CloseDB()
			for {
				if err := p.server.HandleCommand(); err != nil {
					if strings.Contains(err.Error(), "connection closed") {
						continue
					}
					panic(err)
				}
			}
		}(originSocket)
	}

}

const COM_QUERY = byte(0x03)
