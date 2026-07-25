package state

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/matttm/imposter-db/protocol"
)

// AuthenticationState represents the authentication phase
// The client sends credentials to the server for validation
type AuthenticationState struct {
	handshakeRequest *protocol.HandshakeV10Payload
	nonce            []byte
}

func (a *AuthenticationState) IsComplete() bool {
	// AuthenticationState is not an exit state
	return false
}

func (a *AuthenticationState) Handle(f *uint32, schema string, remote net.Conn, client net.Conn, username, password string, cancel context.CancelFunc) (State, error) {
	log.Println("AuthenticationState: Reading authentication response from server...")

	// Read response from server
	b, err := protocol.ReadPacket(remote)
	if err != nil {
		return nil, fmt.Errorf("AuthenticationState: error reading from server: %w", err)
	}
	log.Printf("Received %d bytes from server", len(b))

	// Check for OK packet (authentication successful)
	if protocol.IsOkPacket(b) {
		log.Println("AuthenticationState: OK packet received - authentication successful")
		// Write to client if not headless
		if client != nil {
			_, err := client.Write(b)
			if err != nil {
				return nil, fmt.Errorf("AuthenticationState: error sending OK packet to client: %w", err)
			}
		}
		// Authentication complete, move to ready state
		return &ReadyState{}, nil
	}

	// Check for AUTH_SWITCH_REQUEST (0xFE)
	if len(b) > 4 && b[4] == protocol.AUTH_SWITCH_REQUEST {
		log.Println("AuthenticationState: AUTH_SWITCH_REQUEST received")
		return &AuthSwitchState{
			nonce: a.nonce,
		}, nil
	}

	// Check for AUTH_MORE_DATA (0x01)
	if len(b) > 4 && b[4] == protocol.AUTH_MORE_DATA {
		log.Println("AuthenticationState: AUTH_MORE_DATA received")
		return &AuthMoreDataState{
			nonce: a.nonce,
		}, nil
	}

	return nil, fmt.Errorf("AuthenticationState: unexpected packet type %x", b[4])
}
