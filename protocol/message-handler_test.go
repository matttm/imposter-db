package protocol

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var commqneTestTable = []*HandshakeTestProps[HandshakeResponse41]{}

func Test_Command_Decode(t *testing.T) {
	for _, entry := range responseTable {
		p, _ := DecodeHandshakeResponse(capForResp, entry.encoded)
		assert.Equal(
			t,
			entry.decoded,
			p,
		)
	}
}
func Test_Command_Encode(t *testing.T) {
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
