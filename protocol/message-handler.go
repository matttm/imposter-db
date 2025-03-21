package protocol

import (
	"database/sql"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"
)

type MessageHandler struct {
	Capabilities uint32 // provided bby the server
	ClientFlags  uint32
	ServerStatus uint32
}

// Method for handling messages when handshake has been done
func (h *MessageHandler) HandleMessage(client, remote net.Conn, localDb *sql.DB) any {
	// i assume next message is a command
	packet := ReadPacket(client)
	cmd := Command(packet[4])
	switch cmd {
	case COM_SLEEP:
	case COM_QUIT:
	case COM_INIT_DB:
	case COM_FIELD_LIST:
	case COM_CREATE_DB:
	case COM_DROP_DB:
	case COM_UNUSED_2:
	case COM_UNUSED_1:
	case COM_STATISTICS:
	case COM_UNUSED_4:
	case COM_CONNECT:
	case COM_UNUSED_5:
	case COM_DEBUG:
	case COM_PING:
	case COM_TIME:
	case COM_DELAYED_INSERT:
	case COM_CHANGE_USER:
	case COM_BINLOG_DUMP:
	case COM_TABLE_DUMP:
	case COM_CONNECT_OUT:
	case COM_REGISTER_SLAVE:
	case COM_STMT_PREPARE:
	case COM_STMT_EXECUTE:
	case COM_STMT_SEND_LONG_DATA:
	case COM_STMT_CLOSE:
	case COM_STMT_RESET:
	case COM_SET_OPTION:
	case COM_STMT_FETCH:
	case COM_DAEMON:
	case COM_BINLOG_DUMP_GTID:
	case COM_RESET_CONNECTION:
	case COM_CLONE:
	case COM_SUBSCRIBE_GROUP_REPLICATION_STREAM:
	case COM_END:
	case COM_QUERY:
		fmt.Println("Routing to usual remote")
		_, err := remote.Write(packet)
		packet = ReadPacket(remote)
		if err != nil {
			panic(err)
		}
		_, err = client.Write(packet)
		if err != nil {
			panic(err)
		}
	default:
		fmt.Println("Unknown Command")
	}
	return nil
}

// Function that upgrades a tcp connection into a mysql protocol by completing a simple handshake
func NewMessageHandler(client, remote net.Conn) (*MessageHandler, error) {
	// NOTE: replace generic forwarding
	// TODO: optimize and decide on io pkg or raw bytes
	SetTCPNoDelay(client)
	SetTCPNoDelay(remote)
	mh := &MessageHandler{}
	var b []byte
	b = ReadPacket(remote)
	_, err := client.Write(b)
	b = ReadPacket(client)
	if err != nil {
		panic(err)
	}
	_, err = remote.Write(b)
	if err != nil {
		panic(err)
	}
	b = ReadPacket(remote)
	_, err = client.Write(b)
	if err != nil {
		panic(err)
	}
	return mh, nil
}
func ReadPacket(c io.Reader) []byte {
	h := make([]byte, 4)
	_, err := io.ReadFull(c, h)
	if err != nil {
		panic(err)
	}
	seqId := h[3]
	sz := binary.LittleEndian.Uint32(append(h[:3], 0x0))
	h[3] = seqId
	payload := make([]byte, sz)
	_, err = io.ReadFull(c, payload)
	if err != nil {
		panic(err)
	}
	return append(h, payload...)
}
func SetTCPNoDelay(conn net.Conn) {
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		tcpConn.SetNoDelay(true)
	}
}

// TODO: implement async ver
func ReadAllMySQLPacketsNonBlocking(conn net.Conn) ([][]byte, error) {
	var packets [][]byte
	buf := make([]byte, 4096)
	offset := 0

	dataChan := make(chan []byte, 1)
	errChan := make(chan error, 1)

	go func() {
		n, err := conn.Read(buf[offset:])
		if err != nil {
			errChan <- err
			return
		}
		dataChan <- buf[:n]
	}()

	select {
	case data := <-dataChan:
		copy(buf[offset:], data)
		offset += len(data)
	case err := <-errChan:
		if err == io.EOF {
			return packets, nil
		}
		return nil, err
	case <-time.After(500 * time.Millisecond): // Timeout to prevent hanging
		return packets, nil
	}

	// Process packets as before...
	return packets, nil
}
