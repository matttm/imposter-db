package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Definition can be found at https://dev.mysql.com/doc/dev/mysql-server/8.4.3/page_protocol_basic_packets.html
type PacketHeader struct {
	Length     uint32
	SequenceId uint8
}

// MySql Protocol::HandshakePacket
//
// Definition can be fond at: https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_connection_phase_packets_protocol_handshake_v10.html
type HandshakePacket struct {
	Header              *PacketHeader
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

func Decode(data []byte) (*HandshakePacket, error) {
	payload := &HandshakePacket{}
	header := &PacketHeader{}
	var payloadLength uint32 = 0
	payloadLength |= uint32(data[0]) | (uint32(data[1]) << 8) | (uint32(data[2]) << 16)
	header.Length = payloadLength
	header.SequenceId = data[3]
	payload.Header = header

	payload.ProtocolVersion = data[4]

	data = data[5:] // this reassignment is done so the read bytes arnt included in search
	fmt.Printf("right before serV %02x\n", data[0])
	// this string is null terminated, sp look dfor null
	serverVersionEndIdx := bytes.IndexByte(data, 0x00)
	payload.ServerVersion = string(data[:serverVersionEndIdx])

	data = data[serverVersionEndIdx+1:] // we add 1 so we move past null byte (nullbyte = data[serverVEIdx+1]
	fmt.Printf("right after serV %02x\n", data[0])

	buffer := bytes.NewReader(data)
	if err := binary.Read(buffer, binary.LittleEndian, &payload.ThreadID); err != nil {
		return payload, err
	}
	fmt.Printf("%02x", payload.ThreadID)

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
	var capabilities uint32
	capabilities |= uint32(payload.CapabilityFlags1) << 16
	capabilities |= uint32(payload.CapabilityFlags2)
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
		// TODO: figurre this prop out
		// authPluginNameBytes, err := buffer.ReadBytes(0)
		// if err != nil {
		// 	return payload, err
		// }
		// payload.AuthPluginName = string(authPluginNameBytes[:len(authPluginNameBytes)-1])
	}

	return payload, nil
}
