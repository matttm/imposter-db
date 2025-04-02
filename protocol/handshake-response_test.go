package protocol

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var capForResp uint32 = 0 | CLIENT_PROTOCOL_41 | CLIENT_PLUGIN_AUTH_LENENC_CLIENT_DATA | CLIENT_PLUGIN_AUTH

var responseTable = []*HandshakeTestProps[HandshakeResponse41]{
	&HandshakeTestProps[HandshakeResponse41]{
		encoded: []byte{0x7, 0xa6, 0x3e, 0x19, 0xff, 0xff, 0xff, 0x0, 0xff, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x4d, 0x50, 0x57, 0x33, 0x0, 0x14, 0xf, 0x47, 0xe1, 0xa0, 0xf6, 0xcc, 0x6d, 0xd1, 0x58, 0x8f, 0x79, 0xe0, 0x63, 0x15, 0xb, 0x37, 0x39, 0x97, 0x7e, 0xbe, 0x6d, 0x79, 0x73, 0x71, 0x6c, 0x5f, 0x6e, 0x61, 0x74, 0x69, 0x76, 0x65, 0x5f, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x0}, // the following bytes are the connection attr to be implemented later  , 0x82, 0x10, 0x5f, 0x72, 0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x6, 0x32, 0x31, 0x2e, 0x30, 0x2e, 0x35, 0xf, 0x5f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x5, 0x38, 0x2e, 0x32, 0x2e, 0x30, 0xf, 0x5f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x6c, 0x69, 0x63, 0x65, 0x6e, 0x73, 0x65, 0x3, 0x47, 0x50, 0x4c, 0xf, 0x5f, 0x72, 0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x5f, 0x76, 0x65, 0x6e, 0x64, 0x6f, 0x72, 0x10, 0x45, 0x63, 0x6c, 0x69, 0x70, 0x73, 0x65, 0x20, 0x41, 0x64, 0x6f, 0x70, 0x74, 0x69, 0x75, 0x6d, 0xc, 0x5f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x11, 0x4d, 0x79, 0x53, 0x51, 0x4c, 0x20, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2f, 0x4a},
		decoded: &HandshakeResponse41{
			ClientFlag:           423536135, // TODO: SHOULD THID BE REAS AS 2 UINT16?
			MaxPacketSize:        16777215,
			CharacterSet:         0xff,
			Filler:               [23]byte{},
			Username:             "MPW3",
			AuthResponseLen:      0,
			AuthResponse:         string([]byte{0xf, 0x47, 0xe1, 0xa0, 0xf6, 0xcc, 0x6d, 0xd1, 0x58, 0x8f, 0x79, 0xe0, 0x63, 0x15, 0xb, 0x37, 0x39, 0x97, 0x7e, 0xbe}),
			Database:             "",
			ClientPluginName:     "mysql_native_password",
			ClientAttributes:     nil,
			ZstdCompressionLevel: 0,
		},
	},
}

func Test_Handshake_Response_Decode(t *testing.T) {
	for _, entry := range responseTable {
		p, _ := DecodeHandshakeResponse(capForResp, entry.encoded)
		assert.Equal(
			t,
			entry.decoded,
			p,
		)
	}
}
func Test_Handshake_Response_Encode(t *testing.T) {
	for _, entry := range responseTable {
		w, err := EncodeHandshakeResponse(capForResp, entry.decoded)
		b := w.Bytes()
		if err != nil {
			panic(err)
		}
		assert.Equal(
			t,
			entry.encoded,
			b,
		)
	}
}

func Test_Handshake_Response_Password_Encode(t *testing.T) {
	type PasswordTest struct {
		Encoded []byte
		Decoded string
		Salt    []byte
	}
	passwordTable := []PasswordTest{
		{
			Encoded: []byte{0xf2, 0xe3, 0xdc, 0x61, 0x10, 0x5, 0xcd, 0x84, 0x2d, 0xbd, 0x14, 0x3e, 0x68, 0x4a, 0xaf, 0x54, 0x9f, 0xe2, 0x8d, 0x37},
			Decoded: "Softrams10#",
			Salt:    []byte{0x5b, 0x2d, 0x59, 0x41, 0x41, 0x25, 0x20, 0x27 /* */, 0x20, 0x75, 0x4b, 0x25, 0x4a, 0x57, 0x1f, 0x7a, 0x9, 0x5d, 0x1, 0x6b},
		},
	}
	for _, entry := range passwordTable {
		e, _ := encryptPassword("mysql_native_password", entry.Salt, string(entry.Decoded))
		assert.Equal(
			t,
			entry.Encoded,
			e,
		)
	}
}
