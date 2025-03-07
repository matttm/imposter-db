package main

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/go-mysql-org/go-mysql/client"
)

var ()

type selection struct {
	database []string
	table    []string
}

func handleConn(c net.Conn, tableName string, db *client.Conn) {
	p := InitializeProxy(c, tableName, db)

	log.Printf("new connection: %s\n", c.RemoteAddr())
	// defer p.CloseProxy()
	for {
		if err := p.server.HandleCommand(); err != nil {
			if strings.Contains(err.Error(), "connection closed") {
				continue
			}
			panic(err)
		}
	}
}
func main() {
	s := selection{}
	log.Printf("Checking for available databases...")

	o := InitOverseerConnection()
	defer o.Close()
	databases := QueryFor(o, SHOW_DB_QUERY)
	s.database = PromptSelection("Choose database", databases)
	log.Printf("You chose %s", s.database[0])

	table := QueryFor(o, SHOW_TABLE_QUERY(s.database[0]))
	s.table = PromptSelection("Choose table", table)
	log.Printf("You chose %s", s.table[0])

	createCommand := QueryForTwoColumns(o, SHOW_CREATE(s.database[0], s.table[0]))[0][1]
	columns := QueryForTwoColumns(o, SELECT_COLUMNS(s.table[0]))

	log.Println(createCommand)
	log.Println(columns)

	insertTemplate := CreateSelectInsertionFromSchema(s.database[0], s.table[0], columns)
	log.Println(insertTemplate)

	inserts := QueryFor(o, insertTemplate)
	for _, v := range inserts {
		log.Println(v)
	}
	var localDb *client.Conn = InitLocalDatabase()
	defer localDb.Close()
	log.Println("Database provider init")
	Populate(localDb, s.database[0], createCommand, inserts)

	// start proxying
	socket, err := net.Listen("tcp", "127.0.0.1:3307")
	if err != nil {
		log.Fatalf("failed to start proxy: %s", err.Error())
	}
	fmt.Printf("Listening on localhost:%d\n", 3307)
	// inputTables := []string{"ACO_MS_DB.APLCTN_RVW_PRD"}
	for {
		originSocket, err := socket.Accept()
		if err != nil {
			log.Fatalf("failed to accept connection: %s", err.Error())
		}
		go handleConn(originSocket, s.table[0], localDb)
	}

}

const COM_QUERY = byte(0x03)
