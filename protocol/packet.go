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
func PackPayload(b []byte) []byte {
	h := make([]byte, 4)
	l := uint32(len(b))
	binary.LittleEndian.PutUint32(h, l)
	return append(h, b...)

}
