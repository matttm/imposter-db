package state

import (
	"context"
	"net"
)

// CommandState represents the command execution phase
// The client sends commands (queries, etc.) to the server
type CommandState struct {
	commandReceived bool
}

func (c *CommandState) IsComplete() bool {
	// CommandState is not an exit state within handshake phase
	return false
}

func (c *CommandState) Handle(f *uint32, schema string, remote net.Conn, client net.Conn, username, password string, cancel context.CancelFunc) (State, error) {
	// TODO: Receive and parse command from client
	// Route command to appropriate handler (query, ping, etc.)
	return &ResultState{}, nil
}
