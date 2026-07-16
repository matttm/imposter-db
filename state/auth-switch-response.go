package state

import (
	"context"
	"net"
)

// AuthSwitchResponseState represents the response to an authentication method switch
// The client responds with credentials using the new authentication method
type AuthSwitchResponseState struct {
	responseReceived bool
}

func (asr *AuthSwitchResponseState) IsComplete() bool {
	return asr.responseReceived
}

func (asr *AuthSwitchResponseState) Handle(f *uint32, schema string, remote net.Conn, client net.Conn, username, password string, cancel context.CancelFunc) State {
	// TODO: Receive and validate authentication switch response from client
	// May need additional exchanges (AuthMoreDataState) or go directly to ReadyState
	asr.responseReceived = true
	return &ReadyState{}
}
