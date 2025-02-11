package main

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/dolthub/go-mysql-server/memory"
)

var ()

func handleConn(c net.Conn, provider *memory.DbProvider) {
	p := InitializeProxy(c, provider)

	log.Printf("new connection: %s\n", c.RemoteAddr())
	// defer p.CloseProxy()
	for {
		if err := p.server.HandleCommand(); err != nil {
			if strings.Contains(err.Error(), "connection closed") {
				continue
			}
			panic(err)
		}
	}
}
func main() {
	log.Printf("Checking for available databases...")
	o := InitOverseerConnection()
	// start proxying
	socket, err := net.Listen("tcp", "127.0.0.1:3307")
	if err != nil {
		log.Fatalf("failed to start proxy: %s", err.Error())
	}
	fmt.Printf("Listening on localhost:%d\n", 3307)
	// inputTables := []string{"ACO_MS_DB.APLCTN_RVW_PRD"}
	provider := InitEmptyDatabase()
	for {
		originSocket, err := socket.Accept()
		if err != nil {
			log.Fatalf("failed to accept connection: %s", err.Error())
		}
		go handleConn(originSocket, provider)
	}

}

const COM_QUERY = byte(0x03)
