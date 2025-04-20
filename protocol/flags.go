package protocol

// defined in: https://dev.mysql.com/doc/dev/mysql-server/latest/group__group__cs__capabilities__flags.html
const (
	CLIENT_LONG_PASSWORD                  uint32 = 0x01000000
	CLIENT_FOUND_ROWS                     uint32 = 0x02000000
	CLIENT_LONG_FLAG                      uint32 = 0x04000000
	CLIENT_CONNECT_WITH_DB                uint32 = 0x08000000
	CLIENT_NO_SCHEMA                      uint32 = 0x10000000
	CLIENT_COMPRESS                       uint32 = 0x20000000
	CLIENT_ODBC                           uint32 = 0x40000000
	CLIENT_LOCAL_FILES                    uint32 = 0x80000000
	CLIENT_IGNORE_SPACE                   uint32 = 0x00010000
	CLIENT_PROTOCOL_41                    uint32 = 0x00020000
	CLIENT_INTERACTIVE                    uint32 = 0x00040000
	CLIENT_SSL                            uint32 = 0x00080000
	CLIENT_IGNORE_SIGPIPE                 uint32 = 0x00100000
	CLIENT_TRANSACTIONS                   uint32 = 0x00200000
	CLIENT_RESERVED                       uint32 = 0x00400000
	CLIENT_RESERVED2                      uint32 = 0x00800000
	CLIENT_MULTI_STATEMENTS               uint32 = 0x00010000 // Note: Bit 16
	CLIENT_MULTI_RESULTS                  uint32 = 0x00020000 // Note: Bit 17
	CLIENT_PS_MULTI_RESULTS               uint32 = 0x00040000 // Note: Bit 18
	CLIENT_PLUGIN_AUTH                    uint32 = 0x00080000 // Note: Bit 19
	CLIENT_CONNECT_ATTRS                  uint32 = 0x00100000 // Note: Bit 20
	CLIENT_PLUGIN_AUTH_LENENC_CLIENT_DATA uint32 = 0x00200000 // Note: Bit 21
	CLIENT_CAN_HANDLE_EXPIRED_PASSWORDS   uint32 = 0x00400000 // Note: Bit 22
	CLIENT_SESSION_TRACK                  uint32 = 0x00800000 // Note: Bit 23
	CLIENT_DEPRECATE_EOF                  uint32 = 0x01000000 // Note: Bit 24
	CLIENT_OPTIONAL_RESULTSET_METADATA    uint32 = 0x02000000 // Note: Bit 25
	CLIENT_ZSTD_COMPRESSION_ALGORITHM     uint32 = 0x04000000 // Note: Bit 26
	CLIENT_QUERY_ATTRIBUTES               uint32 = 0x08000000 // Note: Bit 27
	MULTI_FACTOR_AUTHENTICATION           uint32 = 0x10000000 // Note: Bit 28
	CLIENT_CAPABILITY_EXTENSION           uint32 = 0x20000000 // Note: Bit 29
	CLIENT_SSL_VERIFY_SERVER_CERT         uint32 = 0x40000000 // Note: Bit 30
	CLIENT_REMEMBER_OPTIONS               uint32 = 0x80000000 // Note: Bit 31
	CLIENT_SECURE_CONNECTION              uint32 = 0x00008000 // Note: Bit 15 (already in correct position)
)

// Define the SERVER_STATUS constants
// TODO: CONVERT TO BIG ENDIAN
const (
	SERVER_STATUS_IN_TRANS             = 1 << iota // 0x0001
	SERVER_STATUS_AUTOCOMMIT                       // 0x0002
	SERVER_MORE_RESULTS_EXISTS                     // 0x0004
	SERVER_QUERY_NO_GOOD_INDEX_USED                // 0x0008
	SERVER_QUERY_NO_INDEX_USED                     // 0x0010
	SERVER_STATUS_CURSOR_EXISTS                    // 0x0020
	SERVER_STATUS_LAST_ROW_SENT                    // 0x0040
	SERVER_STATUS_DB_DROPPED                       // 0x0080
	SERVER_STATUS_NO_BACKSLASH_ESCAPES             // 0x0100
	SERVER_STATUS_METADATA_CHANGED                 // 0x0200
	SERVER_QUERY_WAS_SLOW                          // 0x0400
	SERVER_PS_OUT_PARAMS                           // 0x0800
	SERVER_STATUS_IN_TRANS_READONLY                // 0x1000
	SERVER_SESSION_STATE_CHANGED                   // 0x2000
)
