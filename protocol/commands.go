package protocol

type Command byte

const (
	COM_SLEEP                              Command = 0x00
	COM_QUIT                               Command = 0x01
	COM_INIT_DB                            Command = 0x02
	COM_QUERY                              Command = 0x03
	COM_FIELD_LIST                         Command = 0x04
	COM_CREATE_DB                          Command = 0x05
	COM_DROP_DB                            Command = 0x06
	COM_UNUSED_2                           Command = 0x07
	COM_UNUSED_1                           Command = 0x08
	COM_STATISTICS                         Command = 0x09
	COM_UNUSED_4                           Command = 0x0A
	COM_CONNECT                            Command = 0x0B
	COM_UNUSED_5                           Command = 0x0C
	COM_DEBUG                              Command = 0x0D
	COM_PING                               Command = 0x0E
	COM_TIME                               Command = 0x0F
	COM_DELAYED_INSERT                     Command = 0x10
	COM_CHANGE_USER                        Command = 0x11
	COM_BINLOG_DUMP                        Command = 0x12
	COM_TABLE_DUMP                         Command = 0x13
	COM_CONNECT_OUT                        Command = 0x14
	COM_REGISTER_SLAVE                     Command = 0x15
	COM_STMT_PREPARE                       Command = 0x16
	COM_STMT_EXECUTE                       Command = 0x17
	COM_STMT_SEND_LONG_DATA                Command = 0x18
	COM_STMT_CLOSE                         Command = 0x19
	COM_STMT_RESET                         Command = 0x1A
	COM_SET_OPTION                         Command = 0x1B
	COM_STMT_FETCH                         Command = 0x1C
	COM_DAEMON                             Command = 0x1D
	COM_BINLOG_DUMP_GTID                   Command = 0x1E
	COM_RESET_CONNECTION                   Command = 0x1F
	COM_CLONE                              Command = 0x20
	COM_SUBSCRIBE_GROUP_REPLICATION_STREAM Command = 0x21
	COM_END                                Command = 0xFE
)

var commandInfo = map[Command]string{
	COM_SLEEP:            "Currently refused by the server. See dispatch_command. Also used internally to mark the start of a session.",
	COM_QUIT:             "See COM_QUIT.",
	COM_INIT_DB:          "See COM_INIT_DB.",
	COM_QUERY:            "See COM_QUERY.",
	COM_FIELD_LIST:       "Deprecated. See COM_FIELD_LIST.",
	COM_CREATE_DB:        "Currently refused by the server. See dispatch_command.",
	COM_DROP_DB:          "Currently refused by the server. See dispatch_command.",
	COM_UNUSED_2:         "Removed, used to be COM_REFRESH.",
	COM_UNUSED_1:         "Removed, used to be COM_SHUTDOWN.",
	COM_STATISTICS:       "See COM_STATISTICS.",
	COM_UNUSED_4:         "Removed, used to be COM_PROCESS_INFO.",
	COM_CONNECT:          "Currently refused by the server.",
	COM_UNUSED_5:         "Removed, used to be COM_PROCESS_KILL.",
	COM_DEBUG:            "See COM_DEBUG.",
	COM_PING:             "See COM_PING.",
	COM_TIME:             "Currently refused by the server.",
	COM_DELAYED_INSERT:   "Functionality removed.",
	COM_CHANGE_USER:      "See COM_CHANGE_USER.",
	COM_BINLOG_DUMP:      "See COM_BINLOG_DUMP.",
	COM_DAEMON:           "Currently refused by the server. Used internally to mark the session as a daemon.",
	COM_RESET_CONNECTION: "See COM_RESET_CONNECTION.",
	COM_END:              "Not a real command. Refused.",
}
