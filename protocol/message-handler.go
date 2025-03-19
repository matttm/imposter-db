package protocol

import "io"

type MessageHandler struct {
	Capabilities uint32 // provided bby the server
	ClientFlags  uint32
	ServerStatus uint32
}

func HandleMessage(r io.Reader) any {
	return nil
}
