package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"strings"
)

var ()

type selection struct {
	database []string
	table    []string
}

func handleConn(c net.Conn, provider *sql.DB) {
	p := InitializeProxy(c, provider)

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
	databases := QueryFor(o, SHOW_DB_QUERY)
	s.database = PromptSelection("Choose database", databases)
	log.Printf("You chose %s", s.database[0])

	table := QueryFor(o, SHOW_TABLE_QUERY(s.database[0]))
	s.table = PromptSelection("Choose table", table)
	var provider *sql.DB = InitEmptyDatabase()
	log.Printf("You chose %s", s.table[0])

	createCommand := QueryForSecondColumn(o, SHOW_CREATE(s.database[0], s.table[0]))
	columns := QueryFor(o, SELECT_COLUMNS(s.table[0]))

	log.Println(createCommand)
	log.Println(columns)

	insertTemplate := CreateSelectInsertionFromSchema(s.database[0], s.table[0], columns)
	log.Println(insertTemplate)

	inserts := QueryFor(o, insertTemplate)
	log.Panicln(inserts[0])

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
		go handleConn(originSocket, provider)
	}

}

const COM_QUERY = byte(0x03)
