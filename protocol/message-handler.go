package protocol

import (
	"encoding/binary"
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
	SetTCPNoDelay(client)
	SetTCPNoDelay(remote)
	mh := &MessageHandler{}
	var b []byte
	b = ReadPacket(remote)
	fmt.Printf("Packet: %02x", b)
	n, err := client.Write(b)
	fmt.Printf("Forwarded bytes: %d\n", n)
	b = ReadPacket(client)
	fmt.Printf("Packet: %02x", b)
	if err != nil {
		panic(err)
	}
	n, err = remote.Write(b)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Forwarded bytes: %d\n", n)
	b = ReadPacket(remote)
	fmt.Printf("Packet: %02x", b)
	n, err = client.Write(b)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Forwarded bytes: %d\n", n)
	return mh, nil
}
func ReadPacket(c io.Reader) []byte {
	h := make([]byte, 4)
	_, err := io.ReadFull(c, h)
	if err != nil {
		panic(err)
	}
	seqId := h[3]
	sz := binary.LittleEndian.Uint32(append(h[:3], 0x0))
	h[3] = seqId
	fmt.Printf("Expected Payload Length: %d\n", sz)
	payload := make([]byte, sz)
	n, err := io.ReadFull(c, payload)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Actual Payload Length: %d\n", n)
	return append(h, payload...)
}
func SetTCPNoDelay(conn net.Conn) {
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		tcpConn.SetNoDelay(true)
	}
}
