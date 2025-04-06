package protocol

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type Client struct {
	respondToHandshakeReq func([]byte) []byte
	handleOkResponse      func([]byte)
}

// Function CompleteSimpleHandshakeV10
//
// Receives packets from `remote` conn and calls the respective client funcs for a simple mysql handshake
func CompleteSimpleHandshakeV10(remote net.Conn, client Client, cancel context.CancelFunc) {
	var b []byte
	// read handshake request
	log.Println("Entering connection phase...")
	b = ReadPackets(remote, cancel)
	log.Println("HandshakeRequest read from server")
	b = client.respondToHandshakeReq(b)
	log.Println("Executed client callback 'respondToHandshakeReq'")
	_, err := remote.Write(b)
	if err != nil {
		panic(err)
	}
	log.Println("Bytes from callback were sent to the server")
	// read ok packet
	b = ReadPackets(remote, cancel)
	log.Println("Packet read from server")
	client.handleOkResponse(b)
}

// Method for handling messages when handshake has been done
func HandleMessage(client, remote, localDb net.Conn, cancel context.CancelFunc) {
	// i assume next message is a command
	packet := ReadPackets(client, cancel)
	if len(packet) <= 4 {
		return
	}
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
		if err != nil {
			panic(err)
		}
		packet = ReadPackets(remote, cancel)
		_, err = client.Write(packet)
		if err != nil {
			panic(err)
		}
	default:
		fmt.Println("Unknown Command")
	}
	return
}

// Reads a connection until there are no bytes to be read ATM
func ReadPackets(c net.Conn, cancel context.CancelFunc) []byte {
	packets := []byte{}
	// hand here, until we have a payload to read
	payload, err := ReadPacket(c)
	if err != nil {
		panic(err)
	}
	packets = append(packets, payload...)
	// since we have one pavket, we just want to check if there
	// are more without hanging, so set a timeout
	c.SetReadDeadline(time.Now().Add(50 * time.Millisecond))
	defer c.SetReadDeadline(time.Time{}) // Reset deadline after function exits
	for {
		// we continue reading packets until a timeout, meaning
		// we have no more packets to be read
		payload, err = ReadPacket(c)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				// Timeout occurred, return whatever we've read so far
				return packets
			}
			if err == io.EOF {
				log.Println("Closing client connection")
				cancel()
				return packets
			}
			panic(err)
		}
		packets = append(packets, payload...)
	}
}

// ReadPacket reads one mysql packet, by examing the length encodings
func ReadPacket(c net.Conn) ([]byte, error) {
	packet := []byte{}
	h := make([]byte, 4)
	_, err := io.ReadFull(c, h)
	if err != nil {
		return nil, err
	}
	seqId := h[3]
	sz := binary.LittleEndian.Uint32(append(h[:3], 0x0))
	h[3] = seqId
	payload := make([]byte, sz)
	_, err = io.ReadFull(c, payload)
	if err != nil {
		panic(err)
	}
	packet = append(h, payload...)
	return packet, nil
}

// TODO: implement async ver
// func ReadAllMySQLPacketsNonBlocking(conn net.Conn) ([][]byte, error) {
// 	var packets [][]byte
// 	buf := make([]byte, 4096)
// 	offset := 0
//
// 	dataChan := make(chan []byte, 1)
// 	errChan := make(chan error, 1)
//
// 	go func() {
// 		n, err := conn.Read(buf[offset:])
// 		if err != nil {
// 			errChan <- err
// 			return
// 		}
// 		dataChan <- buf[:n]
// 	}()
//
// 	select {
// 	case data := <-dataChan:
// 		copy(buf[offset:], data)
// 		offset += len(data)
// 	case err := <-errChan:
// 		if err == io.EOF {
// 			return packets, nil
// 		}
// 		return nil, err
// 	case <-time.After(50 * time.Millisecond): // Timeout to prevent hanging
// 		return packets, nil
// 	}
//
// 	// Process packets as before...
// 	return packets, nil
// }
