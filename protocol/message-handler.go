package protocol

import (
	"fmt"
	"io"
	"net"
)

type MessageHandler struct {
	Capabilities uint32 // provided bby the server
	ClientFlags  uint32
	ServerStatus uint32
}

// Method for handling messages when handshake has been done
func (h *MessageHandler) HandleMessage(r io.Reader) any {
	return nil
}

// Function that upgrades a tcp connection into a mysql protocol by completing a simple handshake
func NewMessageHandler(client, remote net.Conn) (*MessageHandler, error) {
	// NOTE: i am using io.Reader/Writers right now in handshaking but will probably use raw byte during command phase
	// TODO: optimize and decide on io pkg or raw bytes
	mh := &MessageHandler{}
	_ = BisectPayload(remote)
	req, err := DecodeHandshakeRequest(remote)
	if err != nil {
		panic(err)
	}
	mh.Capabilities = req.GetCapabilities()
	b, err := EncodeHandshakeRequest(req)
	if err != nil {
		panic(err)
	}
	// forwarding the req from server to client
	client.Write(b.Bytes())
	h := BisectPayload(client)
	response := make([]byte, h.Length)
	_, err = client.Read(response)
	res, err := DecodeHandshakeResponse(mh.Capabilities, response)
	b, err = EncodeHandshakeResponse(mh.Capabilities, res)
	// forwarding to remote
	remote.Write(b.Bytes())
	// now i expect ok
	_ = BisectPayload(remote)
	ok := DecodeOkPacket(mh.Capabilities, remote)
	if ok.Header != 0x0 {
		fmt.Printf("OK packet was not received in response to a simple handshake as expected")
	}
	// See "Important settings" section.
	return mh, nil
}
