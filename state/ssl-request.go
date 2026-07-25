package state

import (
	"context"
	"log"
	"net"

	"github.com/matttm/imposter-db/protocol"
)

// SSLRequestState represents the optional SSL request phase
// The client can request to establish an SSL connection
type SSLRequestState struct {
	handshakeRequest *protocol.HandshakeV10Payload
}

func (s *SSLRequestState) IsComplete() bool {
	// SSLRequestState is not an exit state
	return false
}

func (s *SSLRequestState) Handle(f *uint32, schema string, remote net.Conn, client net.Conn, username, password string, cancel context.CancelFunc) (State, error) {
	// TODO: Check if client sent SSLRequest or HandshakeResponse
	// For now, assume no SSL and proceed directly to HandshakeResponse
	log.Println("SSLRequestState: Skipping SSL (not implemented)")

	return &HandshakeResponseState{
		handshakeRequest: s.handshakeRequest,
	}, nil
}
