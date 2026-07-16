package state

import (
	"context"
	"net"
)

// AuthSwitchState represents the authentication method switch phase
// The server requests the client to use a different authentication method
type AuthSwitchState struct {
	switchRequested bool
}

func (as *AuthSwitchState) IsComplete() bool {
	return as.switchRequested
}

func (as *AuthSwitchState) Handle(f *uint32, schema string, remote net.Conn, client net.Conn, username, password string, cancel context.CancelFunc) State {
	// TODO: Send AuthSwitchRequest packet to client with new authentication method
	as.switchRequested = true
	return &AuthSwitchResponseState{}
}
