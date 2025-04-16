package protocol

import (
	"bytes"
	"encoding/binary"
)

type AuthSwitchResponse struct {
	data string
}

func EncodeAuthSwitchResponse(res *AuthSwitchResponse) []byte {
	b := []byte{}
	buffer := bytes.NewBuffer(b)
	if err := binary.Write(buffer, binary.LittleEndian, []byte(res.data)); err != nil {
		panic(err)
	}
	return buffer.Bytes()
}
