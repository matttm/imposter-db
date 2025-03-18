package protocol

// Documentation: https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_basic_ok_packet.html

// OKPacket represents the structure of an OK packet
// These rules distinguish whether the packet represents OK or EOF:
//
// OK: header = 0 and length of packet > 7
// EOF: header = 0xfe and length of packet < 9
type OKPacket struct {
	Header           byte   // 1 byte (0x00 or 0xFE)
	AffectedRows     uint64 // Variable length encoded unsigned integer (LENENC)
	LastInsertID     uint64 // Variable length encoded unsigned integer (LENENC)
	StatusFlags      uint16 // 2 bytes if CLIENT_PROTOCOL_41 or CLIENT_TRANSACTIONS
	Warnings         uint16 // 2 bytes if CLIENT_PROTOCOL_41
	Info             string // String (LENENC) for human-readable status information
	SessionStateInfo string // String (LENENC) for session state info, only if SERVER_SESSION_STATE_CHANGED
	Capabilities     uint32 // Capabilities flags to check conditional fields
}

func Decode(b []byte) *OKPacket {
	packet := &OKPacket{
		Capabilities: capabilities,
		StatusFlags:  statusFlags,
		StatusFlags2: statusFlags2,
	}

	if capabilities&CLIENT_PROTOCOL_41 != 0 {
		// Handle the additional fields for CLIENT_PROTOCOL_41
		packet.Warnings = 0 // placeholder, adjust as necessary
	}

	if capabilities&CLIENT_SESSION_TRACK != 0 {
		if sessionStateChanged || info != "" {
			packet.Info = info
		}

		if sessionStateChanged {
			packet.SessionStateInfo = sessionStateInfo
		}
	} else {
		packet.Info = info
	}

	// Handle the affected rows and last insert ID logic (to be determined by packet content)

	return packet
}
