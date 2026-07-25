package state

import (
	"context"
	"log"
	"net"
)

// RunHandshakeStateMachine orchestrates the MySQL handshake authentication process
// using the state pattern. It transitions through various states (Handshake, SSL,
// HandshakeResponse, Authentication, etc.) until reaching ReadyState (the exit point).
//
// Parameters:
//   - f: Pointer to a uint32 to store the client capability flags negotiated during handshake.
//   - schema: The default database/schema to use for the connection.
//   - remote: The net.Conn representing the connection to the remote MySQL server.
//   - client: The net.Conn representing the connection to the client (may be nil for headless mode).
//   - username: The username to authenticate with.
//   - password: The password to authenticate with.
//   - cancel: A context.CancelFunc to allow cancellation of the handshake process.
//
// The function will panic if any state returns an error. All errors from individual states
// are converted to panics here, keeping the state implementations clean.
func RunHandshakeStateMachine(f *uint32, schema string, remote net.Conn, client net.Conn, username, password string, cancel context.CancelFunc) {
	log.Println("Starting handshake state machine")
	
	// Begin with the EntryPoint state
	currentState := State(&EntryPoint{})
	
	// Transition through states until we reach an exit state (IsComplete returns true)
	for !currentState.IsComplete() {
		log.Printf("In state: %T", currentState)
		nextState, err := currentState.Handle(f, schema, remote, client, username, password, cancel)
		if err != nil {
			log.Panicf("Error in state %T: %v", currentState, err)
		}
		currentState = nextState
	}
	
	log.Printf("Handshake complete, reached exit state: %T", currentState)
}
