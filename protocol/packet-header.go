package protocol

// Definition can be found at https://dev.mysql.com/doc/dev/mysql-server/8.4.3/page_protocol_basic_packets.html
type PacketHeader struct {
	Length     uint32
	SequenceId uint8
}

func StripPacketHeader(b []byte) (*PacketHeader, int) {
	header := &PacketHeader{}
	var payloadLength uint32 = 0
	payloadLength |= uint32(b[0]) | (uint32(b[1]) << 8) | (uint32(b[2]) << 16)
	header.Length = payloadLength
	header.SequenceId = b[3]
	return header, 4
}
