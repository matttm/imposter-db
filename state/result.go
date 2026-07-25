package state

import (
	"context"
	"net"
)

// ResultState represents the result sending phase
// The server sends query results back to the client
type ResultState struct {
	resultSent bool
}

func (r *ResultState) IsComplete() bool {
	// ResultState is not an exit state
	return false
}

func (r *ResultState) Handle(f *uint32, schema string, remote net.Conn, client net.Conn, username, password string, cancel context.CancelFunc) (State, error) {
	// TODO: Send query results to client
	// Results can be from CommandPhase packets (OK_Packet, ERR_Packet, ResultSet, etc.)
	// Return to CommandState to wait for the next command
	return &CommandState{}, nil
}
