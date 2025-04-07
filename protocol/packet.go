package protocol

import (
	"encoding/binary"
	"io"
)

// Definition can be found at https://dev.mysql.com/doc/dev/mysql-server/8.4.3/page_protocol_basic_packets.html
type PacketHeader struct {
	Length     uint32
	SequenceId uint8
}

const (
	// Server Responses/Status (Examples - Check your server documentation)
	OK_PACKET        byte = 0x00 // Successful operation
	ERR_PACKET       byte = 0xFF // Error occurred
	EOF_PACKET       byte = 0xFE // End of result set
	AUTH_SWITCH      byte = 0xFE // Authentication method switch (same value as AUTH_SWITCH_REQUEST)
	MORE_RESULTS_SET byte = 0xFB // More result sets following
)

type Packet[T any] struct {
	h       *PacketHeader
	payload *T
}

func BisectPayload(r io.Reader) *PacketHeader {
	header := &PacketHeader{}
	x := ReadNBytes(r, 3)
	x = append(x, 0x0)
	header.Length = binary.LittleEndian.Uint32(x)
	header.SequenceId = ReadByte(r)
	return header
}
func PackPayload(b []byte, seq byte) []byte {
	h := make([]byte, 4)
	l := uint32(len(b))
	binary.LittleEndian.PutUint32(h, l)
	h[3] = seq
	return append(h, b...)

}
