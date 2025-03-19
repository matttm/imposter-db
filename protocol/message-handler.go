package protocol

import "io"

type MessageHandler struct {
	capabilities uint32 // provided bby the server
	clientFlags  uint32
	serverStatus uint32
}

func HandleMessage(r io.Reader) any {
	return nil
}
