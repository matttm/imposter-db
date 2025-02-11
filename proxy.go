package main

import (
	"fmt"
	"log"
	"net"

	"github.com/dolthub/go-mysql-server/memory"
	"github.com/go-mysql-org/go-mysql/client"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/server"
)

type Proxy struct {
	server *server.Conn // proxy server-side -- from client to server
	client *client.Conn // proxy server0side 00 from server to real db
	p      *memory.DbProvider
}

func InitializeProxy(c net.Conn, pro *memory.DbProvider) *Proxy {
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
	p.p = pro
	return p
}

func (p *Proxy) QueryRemote(query string, args ...interface{}) (*mysql.Result, error) {
	if p.client == nil {
		log.Panicf("Error: client is nil")
	}
	return p.client.Execute(query, args...)
}

func (p *Proxy) CloseProxy() {
	p.server.Close()
	p.client.Close()
}
