package protocol

import (
	_ "github.com/go-sql-driver/mysql"

	"database/sql"
	"fmt"
	"log"
	"net"
)

type Proxy struct {
	client    net.Conn
	remote    net.Conn
	localDb   *sql.DB
	tableName string
	handler   *MessageHandler
}

func InitializeProxy(c net.Conn, host string, db *sql.DB, tableName string) *Proxy {
	p := &Proxy{}
	// TODO: implement handshake protocol here?
	// i dont think i  can use the below as it would hide the handshake to me
	// _conn, err := server.NewConn(c, "root", "", NewRemoteHandler(_client, tableName, db))

	// im going to build up the tcp connectin to mysql protocol
	remote, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, 3306))
	if err != nil {
		panic(err)
	}
	mh, _ := NewMessageHandler(c, remote)
	log.Println("Handshake protocol with remote was successful")

	p.remote = remote
	p.client = c // TODO: wrap this `c` as to not have raw data
	p.tableName = tableName
	p.localDb = db
	p.handler = mh
	return p
}
func (p *Proxy) HandleCommand() {
	p.handler.HandleMessage(p.client, p.remote, p.localDb)
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
