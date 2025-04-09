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

// Function CompleteHandshakeV10
//
// Receives packets from `remote` conn and calls the respective client funcs for a simple mysql handshake
func CompleteHandshakeV10(remote net.Conn, client net.Conn, clientFn Client, cancel context.CancelFunc) {
	clientWrite := func(b []byte) {
		if client == nil {
			return
		}
		n, err := client.Write(b)
		// read ok
		if err != nil {
			panic(err)
		}
		log.Printf("%d bytes sent to client", n)
	}
	var b []byte
	// read handshake request
	log.Println("Entering connection phase (without SSL)...")
	b, _ = ReadPacket(remote)
	log.Println("HandshakeRequest read from server")
	// got the salt aNd responded with my scramble
	clientWrite(b) // NOTE: thinking i have to keep client in-the-loop
	b = clientFn.respondToHandshakeReq(b)
	log.Println("Executed client callback 'respondToHandshakeReq'")
	_, err := remote.Write(b)
	if err != nil {
		panic(err)
	}
	log.Println("Bytes from callback were sent to the server")
	b, _ = ReadPacket(remote)
	clientWrite(b) // NOTE: thinking i have to keep client in-the-loop
	log.Println("Packet read from server")
	if isOkPacket(b) {
		return
	}
	// if not ok packet, then Prootocol::AuthMoreData
	// getting auth switch request -- should have header 0x01 followed by 0x04 indicating perform full auth (not cached)
	if b[4] != AUTH_MORE_DATA {
		log.Panicf("AuthMoreData was expected but got %x", b)
	}
	if b[5] == 0x03 {
		// this is FAST_AUTH_SUCCESS
		log.Println("FAST_AUTH_SUCCESS received")
		b, _ = ReadPacket(remote)
		if isOkPacket(b) {
			log.Println("OK packet received")
			clientWrite(b)
			return
		} else {
			log.Panic("Received FAST_AUTH_SUCCESS followed by non-OK packet")
		}
	}
	// since im not doing an ssl -- ask for rsa public key
	reqKeyPacket := PackPayload([]byte{0x02}, 3)
	_, err = remote.Write(reqKeyPacket)
	if err != nil {
		panic(err)
	}
	pem, _ := ReadPacket(remote)
	pem = pem[4:] // removing header
	pem = pem[1:] // removing header for AuthMoreData 0x01
	e := encryptPassword(pem, []byte("mypassword"))
	b = PackPayload(e, 5)
	_, err = remote.Write(b)
	if err != nil {
		panic(err)
	}
	_, _ = ReadPacket(remote)
}

// Method for handling messages when handshake has been done
func HandleMessage(client, remote, localDb net.Conn, cancel context.CancelFunc) {
	// i assume next message is a command
	packet := ReadPackets(client, cancel)
	if len(packet) <= 4 {
		return
	}
	log.Printf("Received command code %x", packet[4])
	log.Printf("command pack is %x", packet)
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
