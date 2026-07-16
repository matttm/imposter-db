package state

import (
	"context"
	"net"
)

// AuthMoreDataState represents the multi-step authentication phase
// The server sends authentication data and waits for client responses
// This state is used when authentication requires multiple exchanges
type AuthMoreDataState struct {
	authExchangeComplete bool
}

func (amd *AuthMoreDataState) IsComplete() bool {
	return amd.authExchangeComplete
}

func (amd *AuthMoreDataState) Handle(f *uint32, schema string, remote net.Conn, client net.Conn, username, password string, cancel context.CancelFunc) State {
	// TODO: Send AuthMoreData packet to client and receive response
	// Continue exchanging data until authentication succeeds or fails
	amd.authExchangeComplete = true
	// After successful multi-step auth, transition to ReadyState
	return &ReadyState{}
}
