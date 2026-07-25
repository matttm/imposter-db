package state

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/matttm/imposter-db/protocol"
)

// HandshakeState represents the initial handshake phase of the MySQL connection
// The server sends a handshake packet to the client
type HandshakeState struct{}

func (h *HandshakeState) IsComplete() bool {
	// HandshakeState is not an exit state
	return false
}

func (h *HandshakeState) Handle(f *uint32, schema string, remote net.Conn, client net.Conn, username, password string, cancel context.CancelFunc) (State, error) {
	// Read handshake request from remote (server)
	log.Println("HandshakeState: Reading handshake request from server...")
	b, err := protocol.ReadPacket(remote)
	if err != nil {
		return nil, fmt.Errorf("HandshakeState: error reading handshake request: %w", err)
	}

	// Send handshake to client (if not in headless mode)
	if client != nil {
		_, err := client.Write(b)
		if err != nil {
			return nil, fmt.Errorf("HandshakeState: error sending handshake to client: %w", err)
		}
		log.Printf("Sent %d bytes handshake to client", len(b))
	}

	// Decode handshake to extract auth plugin data (nonce)
	req, err := protocol.DecodeHandshakeRequest(b[4:])
	if err != nil {
		return nil, fmt.Errorf("HandshakeState: error decoding handshake: %w", err)
	}
	log.Println("HandshakeRequest read from server")

	// Store the handshake request for later use in authentication
	return &SSLRequestState{
		handshakeRequest: req,
	}, nil
}
