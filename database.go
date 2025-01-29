package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/go-mysql-org/go-mysql/client"
	"github.com/go-mysql-org/go-mysql/server"
)

type Proxy struct {
	in  *server.Conn // proxy server-side -- from client to server
	out *client.Conn // proxy server0side 00 from server to real db
}

func InitializeProxyIn(c net.Conn) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")

	p := &Proxy{}
	_conn, err := server.NewConn(c, user, pass, server.EmptyHandler{})
	if err != nil {
		panic(err)
	}
	_client, err := client.Connect(fmt.Sprintf("%s:%s", host, port), user, pass, "ACO_MS_DB")
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.

	log.Fatalln("Database was successfully connected to")

	p.in = _conn
	p.out = _client
}

func CloseDB() {
	proxyIn.Close()
}

func GetDatabase() *server.Conn {
	if proxyIn == nil {
		log.Fatalf("Error: database not initialized")
	}
	return proxyIn
}
