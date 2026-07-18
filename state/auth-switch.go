package state

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/matttm/imposter-db/protocol"
)

// AuthSwitchState represents the authentication method switch phase
// The server requests the client to use a different authentication method
type AuthSwitchState struct {
	nonce []byte
}

func (as *AuthSwitchState) IsComplete() bool {
	// AuthSwitchState is not an exit state
	return false
}

func (as *AuthSwitchState) Handle(f *uint32, schema string, remote net.Conn, client net.Conn, username, password string, cancel context.CancelFunc) (State, error) {
	log.Println("AuthSwitchState: Handling authentication method switch...")

	// Read the AUTH_SWITCH_REQUEST packet from remote (already read in AuthenticationState)
	// We need to get it again or pass it through state
	// For now, just decode what we expect
	b, err := protocol.ReadPacket(remote)
	if err != nil {
		return nil, fmt.Errorf("AuthSwitchState: error reading AUTH_SWITCH_REQUEST: %w", err)
	}

	switchRequest := protocol.DecodeAuthSwitchRequest(protocol.CLIENT_PROTOCOL_41, b[4:])

	// Hash the password using the new auth method
	hash, err := protocol.HashPassword(
		switchRequest.PluginName,
		[]byte(switchRequest.PluginData),
		password,
	)
	if err != nil {
		return nil, fmt.Errorf("AuthSwitchState: error hashing password: %w", err)
	}

	// Create AUTH_SWITCH_RESPONSE
	resp := protocol.EncodeAuthSwitchResponse(&protocol.AuthSwitchResponse{AuthResponse: string(hash)}).Bytes()
	_, err = remote.Write(protocol.PackPayload(resp, b[3]+1))
	if err != nil {
		return nil, fmt.Errorf("AuthSwitchState: error sending AUTH_SWITCH_RESPONSE: %w", err)
	}

	log.Println("AuthSwitchState: AUTH_SWITCH_RESPONSE sent")

	// Transition to AuthSwitchResponseState to handle next response
	return &AuthSwitchResponseState{
		nonce: as.nonce,
	}, nil
}
