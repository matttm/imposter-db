package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	socket, err := net.Listen("tcp", ":3307")
	if err != nil {
		log.Fatalf("failed to start proxy: %s", err.Error())
	}
	fmt.Printf("Listening on localhost%d", 3307)
	originSocket, err := socket.Accept()
	p := InitializeProxy(originSocket)

	log.Printf("new connection: %s", originSocket.RemoteAddr())
	if err != nil {
		log.Fatalf("failed to accept connection: %s", err.Error())
	}
	for {
		p.in.HandleCommand()
	}

}

const COM_QUERY = byte(0x03)
