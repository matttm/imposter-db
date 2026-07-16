package state

import (
	"context"
	"net"
)

// HandshakeState represents the initial handshake phase of the MySQL connection
// The server sends a handshake packet to the client
type HandshakeState struct {
	initialized bool
}

func (h *HandshakeState) IsComplete() bool {
	return h.initialized
}

func (h *HandshakeState) Handle(f *uint32, schema string, remote net.Conn, client net.Conn, username, password string, cancel context.CancelFunc) State {
	// TODO: Send handshake packet to client from remote connection
	h.initialized = true
	// Client may send SSLRequest or go directly to HandshakeResponse
	return &SSLRequestState{}
}
