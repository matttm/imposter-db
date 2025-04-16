package main

import (
	"context"
	// "database/sql"
	"fmt"
	"log"
	"net"

	"github.com/matttm/imposter-db/protocol"
)

var ()

type selection struct {
	database []string
	table    []string
}

func handleConn(c net.Conn, tableName string) {
	ctx, cancel := context.WithCancel(context.Background()) // Create a cancelable context
	p := protocol.InitializeProxy(c, host, tableName, cancel, user, pass)

	log.Printf("new connection: %s\n", c.RemoteAddr())
	// defer c.Close()
	defer p.CloseProxy()
	for {
		select {
		case <-ctx.Done():
			return // Exit loop when context is done
		// TODO: add monitoring here
		default:
			p.HandleCommand()
		}
	}
}
func main() {
	// s := selection{}
	// log.Printf("Checking for available databases...")
	//
	// o := InitOverseerConnection()
	// defer o.Close()
	// databases := QueryFor(o, SHOW_DB_QUERY)
	// s.database = PromptSelection("Choose database", databases)
	// log.Printf("You chose %s", s.database[0])
	//
	// table := QueryFor(o, SHOW_TABLE_QUERY(s.database[0]))
	// s.table = PromptSelection("Choose table", table)
	// log.Printf("You chose %s", s.table[0])
	//
	// createCommand := QueryForTwoColumns(o, SHOW_CREATE(s.database[0], s.table[0]))[0][1]
	// columns := QueryForTwoColumns(o, SELECT_COLUMNS(s.table[0]))
	//
	// log.Println(createCommand)
	// log.Println(columns)
	//
	// insertTemplate := CreateSelectInsertionFromSchema(s.database[0], s.table[0], columns)
	//
	// inserts := QueryFor(o, insertTemplate)
	// var localDb *sql.DB = InitLocalDatabase()
	// log.Println("Database provider init")
	// Populate(localDb, s.database[0], createCommand, inserts)
	// // close db as were going to open it again in raw tcp form
	// localDb.Close()

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
		go handleConn(originSocket, "") // fmt.Sprintf("%s.%s", s.database[0], s.table[0]))
	}

}

const COM_QUERY = byte(0x03)
