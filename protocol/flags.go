package protocol

// defined in: https://dev.mysql.com/doc/dev/mysql-server/latest/group__group__cs__capabilities__flags.html
const (
	CLIENT_LONG_PASSWORD                  = 0b00000000000000000000000000000001
	CLIENT_FOUND_ROWS                     = 0b00000000000000000000000000000010
	CLIENT_LONG_FLAG                      = 0b00000000000000000000000000000100
	CLIENT_CONNECT_WITH_DB                = 0b00000000000000000000000000001000
	CLIENT_NO_SCHEMA                      = 0b00000000000000000000000000010000
	CLIENT_COMPRESS                       = 0b00000000000000000000000000100000
	CLIENT_ODBC                           = 0b00000000000000000000000001000000
	CLIENT_LOCAL_FILES                    = 0b00000000000000000000000010000000
	CLIENT_IGNORE_SPACE                   = 0b00000000000000000000000100000000
	CLIENT_PROTOCOL_41                    = 0b00000000000000000000001000000000
	CLIENT_INTERACTIVE                    = 0b00000000000000000000010000000000
	CLIENT_SSL                            = 0b00000000000000000000100000000000
	CLIENT_IGNORE_SIGPIPE                 = 0b00000000000000000001000000000000
	CLIENT_TRANSACTIONS                   = 0b00000000000000000010000000000000
	CLIENT_RESERVED                       = 0b00000000000000000100000000000000
	CLIENT_RESERVED2                      = 0b00000000000000001000000000000000
	CLIENT_MULTI_STATEMENTS               = 0b00000000000000010000000000000000
	CLIENT_MULTI_RESULTS                  = 0b00000000000000100000000000000000
	CLIENT_PS_MULTI_RESULTS               = 0b00000000000001000000000000000000
	CLIENT_PLUGIN_AUTH                    = 0b00000000000010000000000000000000
	CLIENT_CONNECT_ATTRS                  = 0b00000000000100000000000000000000
	CLIENT_PLUGIN_AUTH_LENENC_CLIENT_DATA = 0b00000000001000000000000000000000
	CLIENT_CAN_HANDLE_EXPIRED_PASSWORDS   = 0b00000000010000000000000000000000
	CLIENT_SESSION_TRACK                  = 0b00000000100000000000000000000000
	CLIENT_DEPRECATE_EOF                  = 0b00000001000000000000000000000000
	CLIENT_OPTIONAL_RESULTSET_METADATA    = 0b00000010000000000000000000000000
	CLIENT_ZSTD_COMPRESSION_ALGORITHM     = 0b00000100000000000000000000000000
	CLIENT_QUERY_ATTRIBUTES               = 0b00001000000000000000000000000000
	MULTI_FACTOR_AUTHENTICATION           = 0b00010000000000000000000000000000
	CLIENT_CAPABILITY_EXTENSION           = 0b00100000000000000000000000000000
	CLIENT_SSL_VERIFY_SERVER_CERT         = 0b01000000000000000000000000000000
	CLIENT_REMEMBER_OPTIONS               = 0b10000000000000000000000000000000
)
