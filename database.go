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

func InitializeProxy(c net.Conn) *Proxy {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")

	p := &Proxy{}
	_conn, err := server.NewConn(c, "root", "", server.EmptyHandler{})
	if err != nil {
		panic(err)
	}
	_client, err := client.Connect(fmt.Sprintf("%s:%s", host, port), user, pass, "ACO_MS_DB")
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.

	log.Println("Database was successfully connected to")

	p.in = _conn
	p.out = _client
	return p
}

func (p *Proxy) CloseDB() {
	p.in.Close()
	p.out.Close()
}
