package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Documentation can be found at https://dev.mysql.com/doc/dev/mysql-server/8.4.3/page_protocol_connection_phase_packets_protocol_handshake_response.html
type HandshakeResponse420 struct {
	Header                 *PacketHeader
	ClientFlag             uint32
	MaxPacketSize          uint32
	CharacterSet           uint8
	Filler                 [23]byte
	Username               string
	AuthResponseLen        uint8  // Length encoded auth response
	AuthResponse           string // Fixed length auth response
	Database               string
	ClientPluginName       string
	ClientAttributesLength uint64
	ClientAttributes       map[string]string
	ZstdCompressionLevel   uint8
}

func DecodeHandshakeResponse(b []byte) (*HandshakeResponse420, error) {

	p := &HandshakeResponse420{}
	h, headerSz := StripPacketHeader(b)
	p.Header = h
	b = b[headerSz:]
	r := bytes.NewReader(b)
	_ = binary.Read(r, binary.LittleEndian, &p.ClientFlag)
	_ = binary.Read(r, binary.LittleEndian, &p.MaxPacketSize)
	_ = binary.Read(r, binary.LittleEndian, &p.CharacterSet)
	_ = binary.Read(r, binary.LittleEndian, &p.Filler)
	p.Username = ReadNullTerminatedString(r)
	// TODO: DOUBLE-CHECK
	if p.ClientFlag&CLIENT_PLUGIN_AUTH_LENENC_CLIENT_DATA != 0 {
		_ = binary.Read(r, binary.LittleEndian, &p.AuthResponseLen)
		p.AuthResponse = ReadFixedLengthString(r, uint(p.AuthResponseLen))
	} else {
		_ = binary.Read(r, binary.LittleEndian, &p.AuthResponseLen)
		p.AuthResponse = ReadFixedLengthString(r, uint(p.AuthResponseLen))
	}
	if p.ClientFlag&CLIENT_CONNECT_WITH_DB != 0 {
		p.Database = ReadNullTerminatedString(r)
	}
	if p.ClientFlag&CLIENT_PLUGIN_AUTH != 0 {
		p.ClientPluginName = ReadNullTerminatedString(r)
	}
	if p.ClientFlag&CLIENT_CONNECT_ATTRS != 0 {
		fmt.Println("Connection attributtes not implementeded")
	}
	if p.ClientFlag&CLIENT_ZSTD_COMPRESSION_ALGORITHM != 0 {
		_ = binary.Read(r, binary.LittleEndian, &p.ZstdCompressionLevel)
	}

	return p, nil
}
