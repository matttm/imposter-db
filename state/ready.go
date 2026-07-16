package state

import (
	"context"
	"net"
)

// ReadyState represents the state after successful authentication
// The connection is established and ready to receive commands
type ReadyState struct {
	initialized bool
}

func (r *ReadyState) IsComplete() bool {
	return r.initialized
}

func (r *ReadyState) Handle(f *uint32, schema string, remote net.Conn, client net.Conn, username, password string, cancel context.CancelFunc) State {
	// TODO: Send OK_Packet to confirm successful authentication
	r.initialized = true
	// Transition to CommandState to handle incoming commands
	return &CommandState{}
}
