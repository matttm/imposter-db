package protocol

import (
	"encoding/binary"
	"io"
)

// Documentation: https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_basic_ok_packet.html

// OKPacket represents the structure of an OK packet
// These rules distinguish whether the packet represents OK or EOF:
//
// OK: header = 0 and length of packet > 7
// EOF: header = 0xfe and length of packet < 9
type OKPacket struct {
	Header           byte   // 1 byte (0x00 or 0xFE)
	AffectedRows     uint64 // Variable length encoded unsigned integer (LENENC)
	LastInsertID     uint64 // Variable length encoded unsigned integer (LENENC)
	StatusFlags      uint16 // 2 bytes if CLIENT_PROTOCOL_41 or CLIENT_TRANSACTIONS
	Warnings         uint16 // 2 bytes if CLIENT_PROTOCOL_41
	Info             string // String (LENENC) for human-readable status information
	SessionStateInfo string // String (LENENC) for session state info, only if SERVER_SESSION_STATE_CHANGED
}

func DecodeOkPacket(capabilities uint32, r io.Reader) *OKPacket {
	p := &OKPacket{}
	h := make([]byte, 1)
	_, _ = r.Read(h)
	p.Header = h[0]
	p.AffectedRows, _ = ReadVarLengthInt(r)
	p.LastInsertID, _ = ReadVarLengthInt(r)

	if capabilities&CLIENT_PROTOCOL_41 != 0 {
		_ = binary.Read(r, binary.LittleEndian, &p.StatusFlags)
		_ = binary.Read(r, binary.LittleEndian, &p.Warnings)
	} else if capabilities&CLIENT_SESSION_TRACK != 0 {
		_ = binary.Read(r, binary.LittleEndian, &p.StatusFlags)
	}
	if capabilities&CLIENT_SESSION_TRACK != 0 {
		if p.StatusFlags&SERVER_SESSION_STATE_CHANGED != 0 || p.StatusFlags != 0 {
			p.Info = ReadLengthEncodedString(r)
		}
		if p.StatusFlags&SERVER_SESSION_STATE_CHANGED != 0 {
			p.SessionStateInfo = ReadLengthEncodedString(r)
		}
	} else {
		// rezt of packet is for this field
		rest, _ := io.ReadAll(r)
		p.Info = string(rest)
	}
	// Handle the affected rows and last insert ID logic (to be determined by packet content)
	return p
}
