package protocol

import (
	"bytes"
	"encoding/binary"
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
	// this string is null terminated, sp look dfor null
	serverVersionEndIdx := bytes.IndexByte(data, 0x00)
	payload.ServerVersion = string(data[:serverVersionEndIdx])

	data = data[serverVersionEndIdx+1:] // we add 1 so we move past null byte (nullbyte = data[serverVEIdx+1]

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
	var capabilities uint32 = payload.GetCapabilities()
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

	var authPluginDataPart2Len uint8
	if int(payload.AuthPluginDataLen) > 0 {
		authPluginDataPart2Len = Max(13, payload.AuthPluginDataLen-8) //subtract the first part
	}

	payload.AuthPluginDataPart2 = make([]byte, authPluginDataPart2Len)
	if _, err := buffer.Read(payload.AuthPluginDataPart2); err != nil {
		return payload, err
	}

	if capabilities&CLIENT_PLUGIN_AUTH != 0 {
		// TODO:
		name := []byte{}
		for true {
			b, err := buffer.ReadByte()
			if err != nil || b == 0x0 {
				break
			}
			name = append(name, b)
		}
		payload.AuthPluginName = string(name)
	}

	return payload, nil
}
func Encode(p *HandshakePacket) (*bytes.Buffer, error) {
	var b []byte
	w := bytes.NewBuffer(b)
	var h uint32
	h |= p.Header.Length | uint32(p.Header.SequenceId)<<24
	if err := binary.Write(w, binary.LittleEndian, &h); err != nil {
		return w, err
	}
	if err := binary.Write(w, binary.LittleEndian, p.ProtocolVersion); err != nil {
		return w, err
	}
	if err := binary.Write(w, binary.LittleEndian, []byte(p.ServerVersion)); err != nil {
		return w, err
	}
	null := []byte{0x0} // for terminating abo e string
	_ = binary.Write(w, binary.LittleEndian, null)
	if err := binary.Write(w, binary.LittleEndian, p.ThreadID); err != nil {
		return w, err
	}
	if err := binary.Write(w, binary.LittleEndian, p.AuthPluginDataPart1); err != nil {
		return w, err
	}
	if err := binary.Write(w, binary.LittleEndian, p.Filler); err != nil {
		return w, err
	}
	if err := binary.Write(w, binary.LittleEndian, p.CapabilityFlags1); err != nil {
		return w, err
	}
	if err := binary.Write(w, binary.LittleEndian, p.CharacterSet); err != nil {
		return w, err
	}
	if err := binary.Write(w, binary.LittleEndian, p.StatusFlags); err != nil {
		return w, err
	}
	if err := binary.Write(w, binary.LittleEndian, p.CapabilityFlags2); err != nil {
		return w, err
	}
	var capabilities uint32 = p.GetCapabilities()
	if capabilities&CLIENT_PLUGIN_AUTH != 0 {
		if err := binary.Write(w, binary.LittleEndian, p.AuthPluginDataLen); err != nil {
			return w, err
		}
	} else {
		if err := binary.Write(w, binary.LittleEndian, null); err != nil {
			return w, err
		}
	}
	if err := binary.Write(w, binary.LittleEndian, p.Reserved); err != nil {
		return w, err
	}
	if err := binary.Write(w, binary.LittleEndian, p.AuthPluginDataPart2); err != nil {
		return w, err
	}
	if capabilities&CLIENT_PLUGIN_AUTH != 0 {
		if err := binary.Write(w, binary.LittleEndian, []byte(p.AuthPluginName)); err != nil {
			return w, err
		}
		_ = binary.Write(w, binary.LittleEndian, null)
	}
	return w, nil
}
func (p *HandshakePacket) GetCapabilities() uint32 {
	var capabilities uint32
	capabilities |= uint32(p.CapabilityFlags1)
	capabilities |= uint32(p.CapabilityFlags2) << 16
	return capabilities
}
