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

// Function CompleteHandshakeV10
//
// Receives packets from `remote` conn and calls the respective client funcs for a simple mysql handshake
func CompleteHandshakeV10(remote net.Conn, client net.Conn, username, password string, cancel context.CancelFunc) {
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
	clientRead := func(_defaultRead []byte) []byte {
		if client == nil {
			return _defaultRead
		}
		b, err := ReadPacket(client)
		// read ok
		if err != nil {
			panic(err)
		}
		return b
	}
	var b []byte
	// read handshake request
	log.Println("Entering connection phase (without SSL)...")
	b, _ = ReadPacket(remote)
	log.Println("HandshakeRequest read from server")
	// got the salt aNd responded with my scramble
	clientWrite(b) // NOTE: thinking i have to keep client in-the-loop
	b = clientRead(makeHandshakeResponseFromRequest(b, username, password))
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
	e := encryptPassword(pem, []byte(password))
	b = PackPayload(e, 5)
	_, err = remote.Write(b)
	if err != nil {
		panic(err)
	}
	_, _ = ReadPacket(remote)
}
func makeHandshakeResponseFromRequest(req []byte, username, password string) []byte {
	log.Println("=============== START 'respondToHandshakeReq'")
	// tear off header
	seq := req[3]
	req = req[4:]
	_req, _ := DecodeHandshakeRequest(req)
	log.Println("Decoding HandshakeRequest via docker connection")
	p, err := hashPassword(
		_req.AuthPluginName,
		append(_req.AuthPluginDataPart1[:], _req.AuthPluginDataPart2...),
		password,
	)
	if err != nil {
		SaveToFile(req, "failed-codings", "authentication-decoding-failure")
		panic(err)
	}
	res := HandshakeResponse41{
		ClientFlag: CLIENT_CAPABILITIES,
		// ClientFlag:           _req.GetCapabilities(),
		MaxPacketSize:        16777215,
		CharacterSet:         0xff,
		Filler:               [23]byte{},
		Username:             username,
		AuthResponseLen:      uint8(len(p)),
		AuthResponse:         string(p),
		Database:             "",
		ClientPluginName:     _req.AuthPluginName,
		ClientAttributes:     nil,
		ZstdCompressionLevel: 0,
	}
	b, _ := EncodeHandshakeResponse(CLIENT_CAPABILITIES, &res)
	log.Println("Encoding HandshakeResponse via docker connection")
	log.Println("=============== END 'respondToHandshakeReq'")
	return PackPayload(b.Bytes(), seq+byte(1))
}

// Method for handling messages when handshake has been done
func HandleMessage(client, remote, localDb net.Conn, cancel context.CancelFunc) {
	// i assume next message is a command
	packet := ReadPackets(client, cancel)
	if len(packet) <= 4 {
		return
	}
	log.Printf("Received command code %x", packet[4])
	cmd := Command(packet[4])
	switch cmd {
	case COM_SLEEP, COM_QUIT, COM_INIT_DB, COM_FIELD_LIST, COM_CREATE_DB, COM_DROP_DB, COM_STATISTICS, COM_CONNECT, COM_DEBUG, COM_PING, COM_TIME, COM_DELAYED_INSERT, COM_CHANGE_USER, COM_BINLOG_DUMP, COM_TABLE_DUMP, COM_CONNECT_OUT, COM_REGISTER_SLAVE, COM_STMT_PREPARE, COM_STMT_EXECUTE, COM_STMT_SEND_LONG_DATA, COM_STMT_CLOSE, COM_STMT_RESET, COM_SET_OPTION, COM_STMT_FETCH, COM_DAEMON, COM_BINLOG_DUMP_GTID, COM_RESET_CONNECTION, COM_CLONE, COM_SUBSCRIBE_GROUP_REPLICATION_STREAM, COM_END, COM_QUERY:
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
	case COM_UNUSED_1, COM_UNUSED_2, COM_UNUSED_4, COM_UNUSED_5:
		fmt.Println("Unused Command")
		log.Panicf("packet: %x", packet)
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
	//
	// just check for error and panic here?
	if packet[4] == ERR_PACKET {
		panic("Error occured")
	}
	return packet, nil
}
