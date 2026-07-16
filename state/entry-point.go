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
	return e.started
}

func (e *EntryPoint) Handle(f *uint32, schema string, remote net.Conn, client net.Conn, username, password string, cancel context.CancelFunc) State {
	// Transition from EntryPoint to HandshakeState
	e.started = true
	return &HandshakeState{}
}
