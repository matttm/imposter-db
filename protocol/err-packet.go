package protocol

import (
	"bytes"
	"encoding/binary"
)

// Documentation: https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_basic_ok_packet.html

// OKPacket represents the structure of an OK packet
//
// ERR: header 0xff
type ErrPacket struct {
	Header         byte   // 1 byte (0x00 or 0xFE)
	ErrorCode      int16  //  error-code
	SqlStateMarker string // [1] # marker of the SQL state
	SqlState       string // [5] SQL state
	ErrorMessage   string // <EOF> human readable error message
}

func DecodeErrPacket(capabilities uint32, b []byte) *ErrPacket {
	p := &ErrPacket{}
	r := bytes.NewReader(b)
	p.Header = ReadByte(r)

	_ = binary.Read(r, binary.LittleEndian, &p.ErrorCode)
	if capabilities&CLIENT_PROTOCOL_41 != 0 {
		p.SqlStateMarker = ReadFixedLengthString(r, 1)
		p.SqlState = ReadFixedLengthString(r, 5)
	}
	p.ErrorMessage = ReadStringEOF(r)
	// Handle the affected rows and last insert ID logic (to be determined by packet content)
	return p
}
