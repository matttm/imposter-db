package protocol

import (
	"context"

	"fmt"
	"log"
	"net"
)

type Proxy struct {
	client    net.Conn
	remote    net.Conn
	localDb   net.Conn
	tableName string
	cancel    context.CancelFunc
}

func InitializeProxy(client net.Conn, host string, tableName string, cancel context.CancelFunc) *Proxy {
	p := &Proxy{}
	p.cancel = cancel

	// im going to build up the tcp connectin to mysql protocol
	// remote, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, 3306))
	// if err != nil {
	// 	panic(err)
	// }
	// TODO: refactor so i can provide user credentials
	local, err := net.Dial("tcp", fmt.Sprintf("%s:%d", "localhost", 3306))
	if err != nil {
		panic(err)
	}
	log.Println("Creating raw tcp connection for local")
	// create struct that implements interface Client, in ./sql.go
	// var _remote Client
	// _remote.respondToHandshakeReq = func(b []byte) []byte {
	// 	_, err := client.Write(b)
	// 	// read handshake rexponde
	// 	b = ReadPackets(client, cancel)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	return b
	// }
	// _remote.handleOkResponse = func(ok []byte) {
	// 	_, err := client.Write(ok)
	// 	// read ok
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }
	var _local_cb Client
	_local_cb.respondToHandshakeReq = func(req []byte) []byte {
		log.Println("=============== START 'respondToHandshakeReq'")
		// tear off header
		seq := req[3]
		req = req[4:]
		_req, _ := DecodeHandshakeRequest(req)
		log.Println("Decoding HandshakeRequest via docker connection")
		p, err := hashPassword(
			_req.AuthPluginName,
			append(_req.AuthPluginDataPart1[:], _req.AuthPluginDataPart2...),
			"mysql_password",
		)
		if err != nil {
			SaveToFile(req, "failed-codings", "authentication-decoding-failure")
			panic(err)
		}
		res := HandshakeResponse41{
			ClientFlag: 423536135, // TODO: SHOULD THID BE REAS AS 2 UINT16?
			// ClientFlag:           _req.GetCapabilities(),
			MaxPacketSize:        16777215,
			CharacterSet:         0xff,
			Filler:               [23]byte{},
			Username:             "root",
			AuthResponseLen:      0,
			AuthResponse:         string(p),
			Database:             "",
			ClientPluginName:     _req.AuthPluginName,
			ClientAttributes:     nil,
			ZstdCompressionLevel: 0,
		}
		b, _ := EncodeHandshakeResponse(0, &res)
		log.Println("Encoding HandshakeResponse via docker connection")
		log.Println("=============== END 'respondToHandshakeReq'")
		return PackPayload(b.Bytes(), seq+byte(1))
	}
	_local_cb.handleOkResponse = func(ok []byte) {
		// nothing to be done here
		// NOTE: i think i need to remove thise once rremote is going to run too, as to not send dup OK packets
		_, err := client.Write(ok)
		// read ok
		if err != nil {
			panic(err)
		}
	}
	// CompleteSimpleHandshakeV10(remote, _remote, cancel)
	CompleteHandshakeV10(local, client, _local_cb, cancel)
	log.Println("Handshake protocol with remote was successful")

	// p.remote = remote
	p.client = client // TODO: wrap this `c` as to not have raw data
	p.localDb = local
	p.tableName = tableName
	return p
}
func (p *Proxy) HandleCommand() {
	HandleMessage(p.client, p.remote, p.localDb, p.cancel)
}

// func (p *Proxy) QueryRemote(query string, args ...interface{}) (*sql.Result, error) {
// 	if p.remote == nil {
// 		log.Panicf("Error: remote is nil")
// 	}
// 	return p.remote.Execute(query, args...)
// }

func (p *Proxy) CloseProxy() {
	// p.remote.Close()
	p.client.Close()
}
