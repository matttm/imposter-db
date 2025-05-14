package protocol

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
)

// CompleteHandshakeV10 performs the MySQL protocol handshake (version 10) between a client and a remote server.
// It acts as a proxy, relaying handshake packets between the client and server, handling authentication negotiation,
// and supporting both standard and full authentication flows (including RSA public key exchange if required).
//
// Parameters:
//   - f: Pointer to a uint32 to store the client capability flags negotiated during handshake.
//   - schema: The default database/schema to use for the connection.
//   - remote: The net.Conn representing the connection to the remote MySQL server.
//   - client: The net.Conn representing the connection to the client (may be nil for headless mode).
//   - username: The username to authenticate with.
//   - password: The password to authenticate with.
//   - cancel: A context.CancelFunc to allow cancellation of the handshake process.
//
// The function panics on unrecoverable errors and logs key handshake steps for debugging.
// It supports both SSL and non-SSL handshakes, but assumes SSL is not negotiated.
// The function handles AuthSwitchRequest and full authentication (including public key retrieval and password encryption).
func CompleteHandshakeV10(f *uint32, schema string, remote net.Conn, client net.Conn, username, password string, cancel context.CancelFunc) {
	// function writes given []byte  to client if not null
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
	// function that reads from clientt if not null
	//
	// when null, the function gi rn as arg, is invoked to
	//  get []byte resembling a response
	clientRead := func(f func() []byte) []byte {
		if client == nil {
			return f()
		}
		b, err := ReadPacket(client)
		// read ok
		if err != nil {
			panic(err)
		}
		log.Printf("%d bytes read from client", len(b))
		return b
	}
	var b []byte
	// read handshake request
	log.Println("Entering connection phase (without SSL)...")
	b, _ = ReadPacket(remote)
	req, _ := DecodeHandshakeRequest(b[4:])
	nonce := append(req.AuthPluginDataPart1[:], req.AuthPluginDataPart2...)
	log.Println("HandshakeRequest read from server")
	// got the salt aNd responded with my scramble
	clientWrite(b)
	// TODO: REFACTOR CLOSURE
	lazy := func() []byte { return NewHandshakeResponse(f, schema, req, username, password) }
	b = clientRead(lazy)
	if client != nil {
		_response, _ := DecodeHandshakeResponse(b[4:])
		*f = _response.ClientFlag
	}
	//
	log.Println("Sending HandshakeResponse")
	_, err := remote.Write(b)
	if err != nil {
		panic(err)
	}
	b, _ = ReadPacket(remote)
	log.Printf("%d bytes read from the server", len(b))
	if isOkPacket(b) {
		clientWrite(b)
		log.Println("Ok packet sent to client")
		return
	}
	// TODO: spend more time on this case
	// checking for auth switch
	if b[4] == AUTH_SWITCH_REQUEST {
		log.Printf("AuthSwitchRequest received")
		switchRequest := DecodeAuthSwitchRequest(CLIENT_CAPABILITIES, b[4:])
		hash, err := hashPassword(
			switchRequest.pluginName,
			[]byte(switchRequest.pluginData),
			password,
		)
		if err != nil {
			panic(err)
		}
		b = []byte{}
		b = append(b, EncodeAuthSwitchResponse(&AuthSwitchResponse{data: string(hash)}).Bytes()...)
		clientWrite(PackPayload(b, 3))
	}
	// if not ok packet, then Prootocol::AuthMoreData
	if b[4] != AUTH_MORE_DATA {
		log.Panicf("AuthMoreData was expected but got %x", b)
	}
	// getting auth switch request -- should have header 0x01 followed by 0x04 indicating perform full auth (not cached)
	if b[5] == FAST_AUTH_SUCCESS {
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
	if b[5] != PERFORM_FULL_AUTH {
		log.Panicf("Expecting perform_full_authentication")
	}
	// clientWrite(b)
	log.Println("Requesting server's public key")
	// since im not doing an ssl -- ask for rsa public key
	lazy = func() []byte { return PackPayload([]byte{0x02}, 3) }
	reqKeyPacket := clientRead(lazy)
	_, err = remote.Write(reqKeyPacket)
	if err != nil {
		panic(err)
	}
	pem, _ := ReadPacket(remote)
	pem = pem[4:] // removing header
	pem = pem[1:] // removing header for AuthMoreData 0x01
	lazy = func() []byte { return encryptPassword(pem, []byte(password), nonce) }
	e := clientRead(lazy)
	b = PackPayload(e, 5)
	_, err = remote.Write(b)
	if err != nil {
		panic(err)
	}
	_, _ = ReadPacket(remote)
	return
}
func NewHandshakeResponse(f *uint32, schema string, req *HandshakeV10Payload, username, password string) []byte {
	log.Println("=============== START 'respondToHandshakeReq'")
	nonce := append(req.AuthPluginDataPart1[:], req.AuthPluginDataPart2...)
	hashed, err := hashPassword(
		req.AuthPluginName,
		nonce,
		password,
	)
	if err != nil {
		panic(err)
	}
	res := HandshakeResponse41{
		ClientFlag:           *f,
		MaxPacketSize:        16777215,
		CharacterSet:         0xff,
		Filler:               [23]byte{},
		Username:             username,
		AuthResponseLen:      uint8(len(hashed)),
		AuthResponse:         string(hashed),
		Database:             schema,
		ClientPluginName:     req.AuthPluginName,
		ClientAttributes:     nil,
		ZstdCompressionLevel: 0,
	}
	b, _ := EncodeHandshakeResponse(&res)
	log.Println("=============== END 'respondToHandshakeReq'")
	return PackPayload(b.Bytes(), 0x01)
}

// HandleMessage processes a single MySQL protocol message from the client connection.
// It reads the next packet from the client, determines the command type, and routes the message
// to the appropriate backend (remote or local database) based on the command and query content.
// For COM_QUERY commands, it inspects the query to decide whether to route to the local or remote
// database, depending on whether the query references a spoofed table name. The function handles
// forwarding packets between the client and the selected backend, including special handling for
// EOF packets depending on client capabilities. For COM_QUIT, it triggers the provided cancel function.
// Unused and unknown commands are logged or ignored. Any errors encountered during packet reading
// or writing will cause the function to panic.
func HandleMessage(clientFlags uint32, client, remote, localDb net.Conn, spoofedTableName string, cancel context.CancelFunc) {
	// i assume next message is a command
	packet, err := ReadPacket(client)
	if len(packet) <= 4 {
		return
	}
	if err != nil {
		panic(err)
	}
	log.Printf("Received command code %x", packet[4])
	cmd := Command(packet[4])
	switch cmd {
	case COM_QUIT:
		cancel()
	case COM_SLEEP, COM_INIT_DB, COM_FIELD_LIST, COM_CREATE_DB, COM_DROP_DB, COM_STATISTICS, COM_CONNECT, COM_DEBUG, COM_PING, COM_TIME, COM_DELAYED_INSERT, COM_CHANGE_USER, COM_BINLOG_DUMP, COM_TABLE_DUMP, COM_CONNECT_OUT, COM_REGISTER_SLAVE, COM_STMT_PREPARE, COM_STMT_EXECUTE, COM_STMT_SEND_LONG_DATA, COM_STMT_CLOSE, COM_STMT_RESET, COM_SET_OPTION, COM_STMT_FETCH, COM_DAEMON, COM_BINLOG_DUMP_GTID, COM_RESET_CONNECTION, COM_CLONE, COM_SUBSCRIBE_GROUP_REPLICATION_STREAM, COM_END:
		_, err = remote.Write(packet)
		if err != nil {
			panic(err)
		}
		packet = ReadPackets(remote, cancel)
		_, err = client.Write(packet)
		if err != nil {
			panic(err)
		}
	case COM_QUERY:
		var queried net.Conn
		if DecodeQuery(clientFlags, packet[4:]).Contains(spoofedTableName) {
			fmt.Println("Routing to local")
			queried = localDb
		} else {
			fmt.Println("Routing to remote")
			queried = remote
		}
		_, err = queried.Write(packet)
		if err != nil {
			panic(err)
		}
		// getting columns
		packet = ReadPackets(queried, cancel)
		_, err = client.Write(packet)
		if err != nil {
			panic(err)
		}
		if clientFlags&CLIENT_DEPRECATE_EOF != 0 {
		} else {
			// this case is when a client is not using eof deperecated
			// so in this instance, when doing a query, the server sends 2 EOFs--an intermediate one following
			// the fieldset and one more following the rows
			//
			// getting rows
			packet = ReadPackets(queried, cancel)
			_, err = client.Write(packet)
			if err != nil {
				panic(err)
			}
		}
		log.Println("Finished query handling")
	case COM_UNUSED_1, COM_UNUSED_2, COM_UNUSED_4, COM_UNUSED_5:
		fmt.Println("Unused Command")
	default:
		fmt.Println("Unknown Command")
	}
	return
}

// ReadPackets continuously reads packets from the provided net.Conn connection until a timeout,
// an EOF, or an OK packet is received. It appends each packet's payload to a byte slice and returns
// the accumulated packets. If a timeout occurs, it returns the packets read so far. If EOF is encountered,
// it logs the closure, calls the provided cancel function, and returns the packets. Any other error
// will cause a panic. The function logs the sequence ID of each received packet and logs when an OK
// packet is received.
func ReadPackets(c net.Conn, cancel context.CancelFunc) []byte {
	packets := []byte{}
	for {
		// we continue reading packets until a timeout, meaning
		// we have no more packets to be read
		payload, err := ReadPacket(c)
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
		log.Printf("receiving seq id %d", payload[3])
		packets = append(packets, payload...)
		if isOkPacket(payload) {
			log.Println("Received an OK packet")
			return packets
		}
	}
}

// ReadPacket reads a packet from the given net.Conn according to the MySQL protocol.
// It first reads the 4-byte packet header to determine the payload size and sequence ID,
// then reads the payload of the specified size. If the packet indicates an error (ERR_PACKET),
// it decodes the error packet and panics with the error message. Returns the full packet
// (header + payload) or an error if reading fails.
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
		e := DecodeErrPacket(0, packet[4:])
		panic(e.ErrorMessage)
	}
	return packet, nil
}
