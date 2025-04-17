package protocol

import (
	"bytes"
)

// doc https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_connection_phase_packets_protocol_auth_switch_request.html
type AuthSwitchRequest struct {
	status     byte
	pluginName string
	pluginData string
}

func DecodeAuthSwitchRequest(capabilities uint32, b []byte) *AuthSwitchRequest {
	r := bytes.NewReader(b)
	p := &AuthSwitchRequest{}
	p.status = ReadByte(r)
	p.pluginName = ReadNullTerminatedString(r)
	p.pluginData = ReadStringEOF(r)
	return p
}
