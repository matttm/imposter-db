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
	return r.resultSent
}

func (r *ResultState) Handle(f *uint32, schema string, remote net.Conn, client net.Conn, username, password string, cancel context.CancelFunc) State {
	// TODO: Send query results to client
	r.resultSent = true
	// Return to CommandState to wait for the next command
	return &CommandState{}
}
