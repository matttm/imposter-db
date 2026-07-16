package state

import (
	"context"
	"net"
)

// SSLRequestState represents the optional SSL request phase
// The client can request to establish an SSL connection
type SSLRequestState struct {
	sslRequested bool
}

func (s *SSLRequestState) IsComplete() bool {
	return s.sslRequested
}

func (s *SSLRequestState) Handle(f *uint32, schema string, remote net.Conn, client net.Conn, username, password string, cancel context.CancelFunc) State {
	// TODO: Check if client sent SSLRequest or HandshakeResponse
	// If SSLRequest, establish SSL connection then wait for HandshakeResponse
	// If no SSL, proceed directly to HandshakeResponse
	s.sslRequested = true
	return &HandshakeResponseState{}
}
