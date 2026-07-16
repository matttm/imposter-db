package state

import (
	"context"
	"net"
)

// HandshakeResponseState represents the phase where the client sends handshake response
// The client responds to the initial server handshake with its capabilities and initial auth data
type HandshakeResponseState struct {
	responseReceived bool
}

func (hr *HandshakeResponseState) IsComplete() bool {
	return hr.responseReceived
}

func (hr *HandshakeResponseState) Handle(f *uint32, schema string, remote net.Conn, client net.Conn, username, password string, cancel context.CancelFunc) State {
	// TODO: Receive and parse HandshakeResponse from client
	// Validate initial authentication data if provided
	hr.responseReceived = true
	// May continue to AuthMoreDataState for multi-step auth, or ReadyState if auth is successful
	return &AuthenticationState{}
}
