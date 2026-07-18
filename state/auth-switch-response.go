package state

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/matttm/imposter-db/protocol"
)

// AuthSwitchResponseState represents the response to an authentication method switch
// The client responds with credentials using the new authentication method
type AuthSwitchResponseState struct {
	nonce []byte
}

func (asr *AuthSwitchResponseState) IsComplete() bool {
	// AuthSwitchResponseState is not an exit state
	return false
}

func (asr *AuthSwitchResponseState) Handle(f *uint32, schema string, remote net.Conn, client net.Conn, username, password string, cancel context.CancelFunc) (State, error) {
	log.Println("AuthSwitchResponseState: Reading response after auth switch...")

	// Read response from server
	b, err := protocol.ReadPacket(remote)
	if err != nil {
		return nil, fmt.Errorf("AuthSwitchResponseState: error reading from server: %w", err)
	}

	// Check for OK packet
	if protocol.IsOkPacket(b) {
		log.Println("AuthSwitchResponseState: OK packet received")
		if client != nil {
			_, err := client.Write(b)
			if err != nil {
				return nil, fmt.Errorf("AuthSwitchResponseState: error sending OK packet to client: %w", err)
			}
		}
		return &ReadyState{}, nil
	}

	// Check for AuthMoreData (for additional multi-step authentication)
	if len(b) > 4 && b[4] == protocol.AUTH_MORE_DATA {
		log.Println("AuthSwitchResponseState: AUTH_MORE_DATA received")
		return &AuthMoreDataState{
			nonce: asr.nonce,
		}, nil
	}

	return nil, fmt.Errorf("AuthSwitchResponseState: unexpected packet type %x", b[4])
}
