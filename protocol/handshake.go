package protocol

import (
	"bytes"
	"encoding/binary"
)

// MySql Protocol::HandshakeV10
//
// Definition can be fond at: https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_connection_phase_packets_protocol_handshake_v10.html
type HandshakeV10 struct {
	ProtocolVersion     uint8
	ServerVersion       string
	ThreadID            uint32
	AuthPluginDataPart1 [8]byte
	Filler              uint8
	CapabilityFlags1    uint16
	CharacterSet        uint8
	StatusFlags         uint16
	CapabilityFlags2    uint16
	AuthPluginDataLen   uint8
	Reserved            [10]byte
	AuthPluginDataPart2 []byte
	AuthPluginName      string
}

func Decode(data []byte) (*HandshakeV10, error) {
	payload := &HandshakeV10{}

	// TODO: FIND DOC FOR THIS HEADER
	_ = data[:4]

	payload.ProtocolVersion = data[4]
	data = data[5:] // this reassignment is done so the read bytes arnt included in search

	// this string is null terminated, sp look dfor null
	term := []byte{0x00}
	serverVersionEndIdx := bytes.Index(data, term)
	payload.ServerVersion = string(data[:serverVersionEndIdx+1]) // we add 1 so null byte is included

	data = data[:serverVersionEndIdx+2] // we add 2 so we move past null byte

	buffer := bytes.NewReader(data)
	if err := binary.Read(buffer, binary.LittleEndian, &payload.ThreadID); err != nil {
		return payload, err
	}

	if err := binary.Read(buffer, binary.LittleEndian, &payload.AuthPluginDataPart1); err != nil {
		return payload, err
	}

	if err := binary.Read(buffer, binary.LittleEndian, &payload.Filler); err != nil {
		return payload, err
	}

	if err := binary.Read(buffer, binary.LittleEndian, &payload.CapabilityFlags1); err != nil {
		return payload, err
	}

	if err := binary.Read(buffer, binary.LittleEndian, &payload.CharacterSet); err != nil {
		return payload, err
	}

	if err := binary.Read(buffer, binary.LittleEndian, &payload.StatusFlags); err != nil {
		return payload, err
	}

	if err := binary.Read(buffer, binary.LittleEndian, &payload.CapabilityFlags2); err != nil {
		return payload, err
	}

	if capabilities&CLIENT_PLUGIN_AUTH != 0 {
		if err := binary.Read(buffer, binary.LittleEndian, &payload.AuthPluginDataLen); err != nil {
			return payload, err
		}
	} else {
		payload.AuthPluginDataLen = 0
		var tempFiller int8
		if err := binary.Read(buffer, binary.LittleEndian, &tempFiller); err != nil {
			return payload, err
		}
	}

	if err := binary.Read(buffer, binary.LittleEndian, &payload.Reserved); err != nil {
		return payload, err
	}

	var authPluginDataPart2Len int
	if int(payload.AuthPluginDataLen) > 0 {
		authPluginDataPart2Len = int(payload.AuthPluginDataLen) - 8 //subtract the first part
		if authPluginDataPart2Len < 13 {
			authPluginDataPart2Len = 13
		}
	}

	payload.AuthPluginDataPart2 = make([]byte, authPluginDataPart2Len)
	if _, err := buffer.Read(payload.AuthPluginDataPart2); err != nil {
		return payload, err
	}

	if capabilities&CLIENT_PLUGIN_AUTH != 0 {
		authPluginNameBytes, err := buffer.ReadBytes(0)
		if err != nil {
			return payload, err
		}
		payload.AuthPluginName = string(authPluginNameBytes[:len(authPluginNameBytes)-1])
	}

	return payload, nil
}
