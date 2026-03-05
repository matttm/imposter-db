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

// InitializeProxy sets up a new Proxy instance by establishing connections between a client and a remote MySQL server,
// as well as a local MySQL server. It performs the initial handshake for both connections, configures client flags,
// and prepares the Proxy for further use. The function requires the client connection, remote host, schema, table name,
// a cancellation function, and user credentials. It returns a pointer to the initialized Proxy.
//
// Parameters:
//   - client: net.Conn representing the incoming client connection.
//   - host: string specifying the remote MySQL server address.
//   - schema: string specifying the database schema to use.
//   - tableName: string specifying the table name to use.
//   - cancel: context.CancelFunc to allow cancellation of ongoing operations.
//   - user: string for the username (used for local connection).
//   - pass: string for the password (used for local connection).
//
// Returns:
//   - *Proxy: a pointer to the initialized Proxy instance.
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

// HandleCommand processes a command received by the Proxy instance.
// It delegates the handling of the message to the HandleMessage function,
// passing along the relevant client flags, client connection, remote connection,
// local database reference, absolute table name, and a cancellation function.
func (p *Proxy) HandleCommand() {
	HandleMessage(p.clientFlags, p.client, p.remote, p.localDb, p.absoluteTableName, p.cancel)
}

// CloseProxy gracefully closes all connections associated with the Proxy instance,
// including the remote, client, and local database connections.
func (p *Proxy) CloseProxy() {
	p.remote.Close()
	p.client.Close()
	p.localDb.Close()
}
