package protocol

import (
// "testing"
// "github.com/stretchr/testify/assert"
)

type HandshakeTestProps[T any] struct {
	encoded []byte
	decoded *T
}

var responseTable = []*HandshakeTestProps[HandshakeResponse420]{
	&HandshakeTestProps[HandshakeResponse420]{
		encoded: []byte{0xd3, 0x0, 0x0, 0x1, 0x7, 0xa6, 0x3e, 0x19, 0xff, 0xff, 0xff, 0x0, 0xff, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x4d, 0x50, 0x57, 0x33, 0x0, 0x14, 0xf, 0x47, 0xe1, 0xa0, 0xf6, 0xcc, 0x6d, 0xd1, 0x58, 0x8f, 0x79, 0xe0, 0x63, 0x15, 0xb, 0x37, 0x39, 0x97, 0x7e, 0xbe, 0x6d, 0x79, 0x73, 0x71, 0x6c, 0x5f, 0x6e, 0x61, 0x74, 0x69, 0x76, 0x65, 0x5f, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x0, 0x82, 0x10, 0x5f, 0x72, 0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x6, 0x32, 0x31, 0x2e, 0x30, 0x2e, 0x35, 0xf, 0x5f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x5, 0x38, 0x2e, 0x32, 0x2e, 0x30, 0xf, 0x5f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x6c, 0x69, 0x63, 0x65, 0x6e, 0x73, 0x65, 0x3, 0x47, 0x50, 0x4c, 0xf, 0x5f, 0x72, 0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x5f, 0x76, 0x65, 0x6e, 0x64, 0x6f, 0x72, 0x10, 0x45, 0x63, 0x6c, 0x69, 0x70, 0x73, 0x65, 0x20, 0x41, 0x64, 0x6f, 0x70, 0x74, 0x69, 0x75, 0x6d, 0xc, 0x5f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x11, 0x4d, 0x79, 0x53, 0x51, 0x4c, 0x20, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2f, 0x4a},
		decoded: &HandshakeResponse420{
			Header: &PacketHeader{
				Length:     211,
				SequenceId: 1,
			},
			ClientFlag:              425313031,
			MaxPacketSize:           16777215,
			CharacterSet:            0xff,
			Filler:                  [23]byte{},
			Username:                "MPW3",
			AuthResponseLenEnc:      nil,
			AuthResponseFixedLength: []byte{0},
			Database:                "",
			ClientPluginName:        "mysql_native_password",
			ClientAttributes:        nil,
			ZstdCompressionLevel:    data[len(data)-1],
		},
	},
}

// func Test_Handshake_Response_Decode(t *testing.T) {
// 	for _, entry := range responseTable {
// 		p, _ := DecodeHandshakeRequest(entry.encoded)
// 		assert.Equal(
// 			t,
// 			entry.decoded,
// 			p,
// 		)
// 	}
// }
// func Test_Handshake_Response_Encode(t *testing.T) {
// 	for _, entry := range responseTable {
// 		w, err := EncodeHandshakeRequest(entry.decoded)
// 		b := w.Bytes()
// 		if err != nil {
// 			panic(err)
// 		}
// 		assert.Equal(
// 			t,
// 			entry.encoded,
// 			b,
// 		)
// 	}
// }
