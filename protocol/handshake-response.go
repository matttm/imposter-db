package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Documentation can be found at https://dev.mysql.com/doc/dev/mysql-server/8.4.3/page_protocol_connection_phase_packets_protocol_handshake_response.html
type HandshakeResponse41 struct {
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

func DecodeHandshakeResponse(b []byte) (*HandshakeResponse41, error) {
	p := &HandshakeResponse41{}
	r := bytes.NewReader(b)
	err := binary.Read(r, binary.LittleEndian, &p.ClientFlag)
	if err != nil {
		panic(err)
	}
	capabilities := p.ClientFlag // TODO: refactor
	err = binary.Read(r, binary.LittleEndian, &p.MaxPacketSize)
	if err != nil {
		panic(err)
	}
	err = binary.Read(r, binary.LittleEndian, &p.CharacterSet)
	if err != nil {
		panic(err)
	}
	err = binary.Read(r, binary.LittleEndian, &p.Filler)
	if err != nil {
		panic(err)
	}
	p.Username = ReadNullTerminatedString(r)
	// TODO: DOUBLE-CHECK
	if capabilities&CLIENT_PLUGIN_AUTH_LENENC_CLIENT_DATA != 0 {
		p.AuthResponse = ReadLengthEncodedString(r)
	} else {
		p.AuthResponseLen = ReadByte(r)
		p.AuthResponse = ReadFixedLengthString(r, uint64(p.AuthResponseLen))
	}
	if capabilities&CLIENT_CONNECT_WITH_DB != 0 {
		p.Database = ReadNullTerminatedString(r)
	}
	if capabilities&CLIENT_PLUGIN_AUTH != 0 {
		p.ClientPluginName = ReadNullTerminatedString(r)
	}
	if capabilities&CLIENT_CONNECT_ATTRS != 0 {
		fmt.Println("Connection attributtes not implementeded")
	}
	if capabilities&CLIENT_ZSTD_COMPRESSION_ALGORITHM != 0 {
		fmt.Println("Zstd compression not implementeded")
		// err = binary.Read(r, binary.LittleEndian, &p.ZstdCompressionLevel)
	}
	return p, nil
}
func EncodeHandshakeResponse(p *HandshakeResponse41) (*bytes.Buffer, error) {
	var b []byte
	w := bytes.NewBuffer(b)
	capabilities := p.ClientFlag
	err := binary.Write(w, binary.LittleEndian, &p.ClientFlag)
	if err != nil {
		panic(err)
	}
	err = binary.Write(w, binary.LittleEndian, &p.MaxPacketSize)
	if err != nil {
		panic(err)
	}
	err = binary.Write(w, binary.LittleEndian, &p.CharacterSet)
	if err != nil {
		panic(err)
	}
	err = binary.Write(w, binary.LittleEndian, &p.Filler)
	if err != nil {
		panic(err)
	}
	WriteNullTerminatedString(w, p.Username)
	// TODO: DOUBLE-CHECK
	WriteLengthEncodedString(w, p.AuthResponse)
	if capabilities&CLIENT_CONNECT_WITH_DB != 0 {
		WriteNullTerminatedString(w, p.Database)
	}
	if capabilities&CLIENT_PLUGIN_AUTH != 0 {
		WriteNullTerminatedString(w, p.ClientPluginName)
	}
	if capabilities&CLIENT_CONNECT_ATTRS != 0 {
		fmt.Println("Connection attributtes not implementeded")
	}
	if capabilities&CLIENT_ZSTD_COMPRESSION_ALGORITHM != 0 {
		fmt.Println("Zstd compression not implementeded")
		// err = binary.Read(r, binary.LittleEndian, &p.ZstdCompressionLevel)
	}
	return w, nil
}
