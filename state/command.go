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
	return c.commandReceived
}

func (c *CommandState) Handle(f *uint32, schema string, remote net.Conn, client net.Conn, username, password string, cancel context.CancelFunc) State {
	// TODO: Receive and parse command from client
	// Route command to appropriate handler (query, ping, etc.)
	c.commandReceived = true
	return &ResultState{}
}
