package protocol

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type PackPayloadTest struct {
	payload  []byte
	expected []byte
}

func Test_Pack_Payload(t *testing.T) {
	_t := []PackPayloadTest{
		{
			payload:  []byte{0x01},
			expected: []byte{0x01, 0x0, 0x0, 0x1, 0x01},
		},
	}
	for _, entry := range _t {
		p := PackPayload(entry.payload, 1)
		assert.Equal(
			t,
			entry.expected,
			p,
		)
	}
}
