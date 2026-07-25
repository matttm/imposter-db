package protocol

import (
	"bytes"
	"encoding/binary"
)

type AuthSwitchResponse struct {
	AuthResponse string
}

func EncodeAuthSwitchResponse(res *AuthSwitchResponse) *bytes.Buffer {
	b := []byte{}
	buffer := bytes.NewBuffer(b)
	if err := binary.Write(buffer, binary.LittleEndian, []byte(res.AuthResponse)); err != nil {
		panic(err)
	}
	return buffer
}
