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

func DecodeHandshakeResponse(capabilities uint32, b []byte) (*HandshakeResponse41, error) {
	p := &HandshakeResponse41{}
	r := bytes.NewReader(b)
	// to mske backwsrds compatable, flags are stored in 2 16-bit parts, so
	// I'll resd them seperately and shift into a uint32
	var partA, partB uint16
	err := binary.Read(r, binary.LittleEndian, &partA)
	if err != nil {
		panic(err)
	}
	err = binary.Read(r, binary.LittleEndian, &partB)
	if err != nil {
		panic(err)
	}
	p.ClientFlag |= uint32(partB)<<16 | uint32(partA)
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
func EncodeHandshakeResponse(capabilities uint32, p *HandshakeResponse41) (*bytes.Buffer, error) {
	var b []byte
	w := bytes.NewBuffer(b)
	// to mske backwsrds compatable, flags are stored in 2 16-bit parts, so
	// I'll resd them seperately and shift into a uint32
	var partA, partB uint16
	partA |= uint16(p.ClientFlag)
	partB |= uint16(p.ClientFlag >> 16)
	err := binary.Write(w, binary.LittleEndian, &partA)
	if err != nil {
		panic(err)
	}
	err = binary.Write(w, binary.LittleEndian, &partB)
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

// TODDO: refactor method to be an enum
func encryptPassword(method string, salt []byte, password string) ([]byte, error) {
	if isNonASCIIorEmpty(method) {
		return []byte{}, fmt.Errorf("Authentication method is undecipherable")
	}
	if authMeth, ok := authMap[method]; ok {
		// https://dev.mysql.com/doc/dev/mysql-server/8.0.40/page_protocol_connection_phase_authentication_methods_native_password_authentication.html
		stage1 := authMeth.Fn([]byte(password))
		dub := authMeth.Fn(stage1[:])
		stage2 := authMeth.Fn(append(salt, dub[:]...))

		scrambled := make([]byte, authMeth.Sz)
		for i := 0; i < authMeth.Sz; i++ {
			scrambled[i] = stage1[i] ^ stage2[i]
		}
		return scrambled, nil
	}
	return []byte{}, fmt.Errorf("Unknown authentication method: %s", method)
}
func xorBytes(a, b []byte) []byte {
	if len(a) != len(b) {
		return nil // Return nil if slices have different lengths
	}
	result := make([]byte, len(a))
	for i := range a {
		result[i] = a[i] ^ b[i]
	}
	return result
}
