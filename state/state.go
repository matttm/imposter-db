package state
import (
	"context"
	"net"
)
type State interface {
	IsComplete() bool
	Handle(f *uint32, schema string, remote net.Conn, client net.Conn, username, password string, cancel context.CancelFunc) State
}