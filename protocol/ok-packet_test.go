package protocol

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _cap uint32 = 0 | CLIENT_PROTOCOL_41

var ok_table = []*HandshakeTestProps[OKPacket]{
	&HandshakeTestProps[OKPacket]{
		encoded: []byte{0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0},
		decoded: &OKPacket{
			Header:       0x0,
			AffectedRows: 0,
			LastInsertID: 0,
			StatusFlags:  0x0002,
			Warnings:     0,
		},
	},
}

func Test_Ok_Packet_Decode(t *testing.T) {
	for _, entry := range ok_table {
		p := DecodeOkPacket(_cap, bytes.NewReader(entry.encoded))
		assert.Equal(
			t,
			entry.decoded,
			p,
		)
	}
}
