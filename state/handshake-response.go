package state

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/matttm/imposter-db/protocol"
)

// HandshakeResponseState represents the phase where the client sends handshake response
// The client responds to the initial server handshake with its capabilities and initial auth data
type HandshakeResponseState struct {
	handshakeRequest *protocol.HandshakeV10Payload
}

func (hr *HandshakeResponseState) IsComplete() bool {
	// HandshakeResponseState is not an exit state
	return false
}

func (hr *HandshakeResponseState) Handle(f *uint32, schema string, remote net.Conn, client net.Conn, username, password string, cancel context.CancelFunc) (State, error) {
	log.Println("HandshakeResponseState: Processing handshake response...")

	// Generate or read handshake response from client
	var responsePacket []byte

	if client == nil {
		// Headless mode: generate response programmatically
		responsePacket = protocol.NewHandshakeResponse(f, schema, hr.handshakeRequest, username, password)
	} else {
		// Read from client
		var err error
		responsePacket, err = protocol.ReadPacket(client)
		if err != nil {
			return nil, fmt.Errorf("HandshakeResponseState: error reading handshake response from client: %w", err)
		}
		log.Printf("Received %d bytes handshake response from client", len(responsePacket))

		// Extract client flags from response
		if len(responsePacket) > 4 {
			response, err := protocol.DecodeHandshakeResponse(responsePacket[4:])
			if err != nil {
				return nil, fmt.Errorf("HandshakeResponseState: error decoding handshake response: %w", err)
			}
			*f = response.ClientFlag
		}
	}

	// Send handshake response to remote server
	log.Println("HandshakeResponseState: Sending handshake response to server...")
	_, err := remote.Write(responsePacket)
	if err != nil {
		return nil, fmt.Errorf("HandshakeResponseState: error sending handshake response to remote: %w", err)
	}

	// Extract and store nonce for password encryption
	nonce := append(hr.handshakeRequest.AuthPluginDataPart1[:], hr.handshakeRequest.AuthPluginDataPart2...)

	// Transition to authentication state
	return &AuthenticationState{
		handshakeRequest: hr.handshakeRequest,
		nonce:            nonce,
	}, nil
}
