package state

import (
	"context"
	"net"
)

// State interface defines the behavior for each state in the MySQL connection protocol
type State interface {
	IsComplete() bool
	Handle(f *uint32, schema string, remote net.Conn, client net.Conn, username, password string, cancel context.CancelFunc) State
}
