package state

import (
	"context"
	"net"
)

// EntryPoint is the initial state of the MySQL connection
// It transitions to HandshakeState to begin the connection protocol
type EntryPoint struct {
	started bool
}

func (e *EntryPoint) IsComplete() bool {
	// EntryPoint is not an exit state
	return false
}

func (e *EntryPoint) Handle(f *uint32, schema string, remote net.Conn, client net.Conn, username, password string, cancel context.CancelFunc) (State, error) {
	// Transition from EntryPoint to HandshakeState
	return &HandshakeState{}, nil
}
