package protocol

import (
	"bytes"
)

// doc https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_connection_phase_packets_protocol_auth_switch_request.html
type AuthSwitchRequest struct {
	Status     byte
	PluginName string
	PluginData string
}

func DecodeAuthSwitchRequest(capabilities uint32, b []byte) *AuthSwitchRequest {
	r := bytes.NewReader(b)
	p := &AuthSwitchRequest{}
	p.Status = ReadByte(r)
	p.PluginName = ReadNullTerminatedString(r)
	p.PluginData = ReadStringEOF(r)
	return p
}
