package protocol

import (
	"context"

	"fmt"
	"log"
	"net"
)

type Proxy struct {
	client    net.Conn
	remote    net.Conn
	localDb   net.Conn
	tableName string
	cancel    context.CancelFunc
}

func InitializeProxy(client net.Conn, host string, tableName string, cancel context.CancelFunc) *Proxy {
	p := &Proxy{}
	p.cancel = cancel

	// im going to build up the tcp connectin to mysql protocol
	remote, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, 3306))
	if err != nil {
		panic(err)
	}
	// TODO: refactor so i can provide user credentials
	local, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, 3306))
	if err != nil {
		panic(err)
	}
	CompleteSimpleHandshakeV10(client, remote, cancel)
	CompleteSimpleHandshakeV10(client, remote, cancel)
	log.Println("Handshake protocol with remote was successful")

	p.remote = remote
	p.client = client // TODO: wrap this `c` as to not have raw data
	p.localDb = local
	p.tableName = tableName
	return p
}
func (p *Proxy) HandleCommand() {
	HandleMessage(p.client, p.remote, p.localDb, p.cancel)
}

// func (p *Proxy) QueryRemote(query string, args ...interface{}) (*sql.Result, error) {
// 	if p.remote == nil {
// 		log.Panicf("Error: remote is nil")
// 	}
// 	return p.remote.Execute(query, args...)
// }

func (p *Proxy) CloseProxy() {
	p.remote.Close()
	p.client.Close()
}
