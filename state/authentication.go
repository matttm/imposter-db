package state

import (
	"context"
	"net"
)

// AuthenticationState represents the authentication phase
// The client sends credentials to the server for validation
type AuthenticationState struct {
	authenticated bool
}

func (a *AuthenticationState) IsComplete() bool {
	return a.authenticated
}

func (a *AuthenticationState) Handle(f *uint32, schema string, remote net.Conn, client net.Conn, username, password string, cancel context.CancelFunc) State {
	// TODO: Validate authentication data from HandshakeResponse
	// If successful, go to ReadyState; if need more data, go to AuthMoreDataState
	a.authenticated = true
	// For single-step auth, transition to ReadyState
	// For multi-step auth, transition to AuthMoreDataState
	return &ReadyState{}
}
