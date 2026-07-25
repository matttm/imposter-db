package state

import (
	"context"
	"net"
)

// State interface defines the behavior for each state in the MySQL connection protocol
type State interface {
	// IsComplete returns true if this state is an exit point in the FSA
	IsComplete() bool
	// Handle processes the state logic and returns the next state and any error encountered
	Handle(f *uint32, schema string, remote net.Conn, client net.Conn, username, password string, cancel context.CancelFunc) (State, error)
}
