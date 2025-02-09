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
	server *server.Conn // proxy server-side -- from client to server
	client *client.Conn // proxy server0side 00 from server to real db
}

func InitializeProxy(c net.Conn) *Proxy {
	host := os.Getenv("DB_HOST")
	// port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")

	p := &Proxy{}
	_client, err := client.Connect(fmt.Sprintf("%s:%d", host, 3306), user, pass, "")
	if err != nil {
		panic(err)
	}
	_conn, err := server.NewConn(c, "root", "", NewRemoteHandler(_client))
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.

	log.Println("Database was successfully connected to")

	p.server = _conn
	p.client = _client
	return p
}

func (p *Proxy) CloseProxy() {
	p.server.Close()
	p.client.Close()
}
