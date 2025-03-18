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

func StripPacketHeader(r io.Reader) *PacketHeader {
	header := &PacketHeader{}
	x := ReadNBytes(r, 3)
	header.Length = binary.LittleEndian.Uint32(x)
	header.SequenceId = ReadByte(r)
	return header
}
