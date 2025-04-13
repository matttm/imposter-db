package protocol

import (
	"context"

	"fmt"
	"log"
	"net"
)

var CLIENT_CAPABILITIES uint32 = CLIENT_LONG_PASSWORD |
	CLIENT_LONG_FLAG |
	CLIENT_PROTOCOL_41 |
	CLIENT_PLUGIN_AUTH |
	CLIENT_SECURE_CONNECTION | //  CLIENT_PLUGIN_AUTH_LENENC_CLIENT_DATA |
	CLIENT_TRANSACTIONS |
	CLIENT_MULTI_RESULTS |
	CLIENT_MULTI_STATEMENTS |
	CLIENT_DEPRECATE_EOF

type Proxy struct {
	client    net.Conn
	remote    net.Conn
	localDb   net.Conn
	tableName string
	cancel    context.CancelFunc
}

func InitializeProxy(client net.Conn, host string, tableName string, cancel context.CancelFunc, user, pass string) *Proxy {
	p := &Proxy{}
	p.cancel = cancel

	// im going to build up the tcp connectin to mysql protocol
	log.Printf("Connection intializing with %s:%s@%s", user, pass, host)
	remote, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, 3306))
	if err != nil {
		panic(err)
	}
	local, err := net.Dial("tcp", fmt.Sprintf("%s:%d", "127.0.0.1", 3306))
	if err != nil {
		panic(err)
	}
	log.Println("Creating raw tcp connection for local")
	// create struct that implements interface Client, in ./sql.go
	CompleteHandshakeV10(remote, client, user, pass, cancel)
	CompleteHandshakeV10(local, nil, "root", "mypassword", cancel)
	log.Println("Handshake protocol with remote was successful")

	p.remote = remote
	p.client = client
	p.localDb = local
	p.tableName = tableName
	return p
}
func (p *Proxy) HandleCommand() {
	HandleMessage(p.client, p.remote, p.localDb, p.cancel)
}

func (p *Proxy) CloseProxy() {
	p.remote.Close()
	p.client.Close()
	p.localDb.Close()
}
