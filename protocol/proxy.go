package protocol

import (
	"context"
	"strings"

	"fmt"
	"log"
	"net"
)

var CLIENT_CAPABILITIES uint32 = CLIENT_LONG_PASSWORD |
	CLIENT_LONG_FLAG |
	CLIENT_PROTOCOL_41 |
	CLIENT_PLUGIN_AUTH |
	CLIENT_SECURE_CONNECTION |
	CLIENT_QUERY_ATTRIBUTES |
	CLIENT_TRANSACTIONS |
	CLIENT_MULTI_RESULTS |
	CLIENT_MULTI_STATEMENTS |
	CLIENT_DEPRECATE_EOF

type Proxy struct {
	client            net.Conn
	clientFlags       uint32
	remote            net.Conn
	localDb           net.Conn
	absoluteTableName string
	cancel            context.CancelFunc
}

func InitializeProxy(client net.Conn, host string, schema, tableName string, cancel context.CancelFunc, user, pass string) *Proxy {
	p := &Proxy{}
	p.cancel = cancel

	var remote net.Conn
	var local net.Conn
	connect := func(f *uint32, _schema, host, _user, _pass string, _client net.Conn) net.Conn {
		// im going to build up the tcp connectin to mysql protocol
		log.Printf("Connection intializing with %s:%s@%s", _user, _pass, host)
		r, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, 3306))
		if err != nil {
			panic(err)
		}
		CompleteHandshakeV10(f, _schema, r, _client, _user, _pass, cancel)
		return r
	}
	remote = connect(&p.clientFlags, schema, host, "", "", client) // no need to provide user/pass when client is non-nil
	p.clientFlags = p.clientFlags &^ CLIENT_CONNECT_ATTRS
	local = connect(&p.clientFlags, schema, "127.0.0.1", "root", "mypassword", nil)
	log.Println("Handshake protocol with remote was successful")

	log.Printf("--------------flags--------------")
	log.Printf("Flag DEPRECATE_EOF set: %t", p.clientFlags&CLIENT_DEPRECATE_EOF != 0)
	log.Printf("Flag PROTOCOL 41   set: %t", p.clientFlags&CLIENT_PROTOCOL_41 != 0)
	log.Printf("Flag SESSION TRACK set: %t", p.clientFlags&CLIENT_SESSION_TRACK != 0)
	log.Printf("Flag PLUGIN AUTH   set: %t", p.clientFlags&CLIENT_PLUGIN_AUTH != 0)
	log.Printf("Flag SECURE CONN   set: %t", p.clientFlags&CLIENT_SECURE_CONNECTION != 0)
	log.Printf("--------------flags--------------")
	p.remote = remote
	p.client = client
	p.localDb = local
	schema = strings.ToLower(schema)
	tableName = strings.ToLower(tableName)
	p.absoluteTableName = fmt.Sprintf("%s.%s", schema, tableName)
	return p
}
func (p *Proxy) HandleCommand() {
	HandleMessage(p.clientFlags, p.client, p.remote, p.localDb, p.absoluteTableName, p.cancel)
}

func (p *Proxy) CloseProxy() {
	p.remote.Close()
	p.client.Close()
	p.localDb.Close()
}
