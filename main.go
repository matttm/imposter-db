package main

import (
	"context"
	"fmt"
	"strings"

	"log"
	"net"

	"github.com/matttm/imposter-db/protocol"
)

var ()

type selection struct {
	databases []string
	tables    []string
}

func handleConn(c net.Conn, schema, tableName string) {
	ctx, cancel := context.WithCancel(context.Background()) // Create a cancelable context
	p := protocol.InitializeProxy(c, host, schema, tableName, cancel, user, pass)

	log.Printf("new connection: %s\n", c.RemoteAddr())
	defer p.CloseProxy()
	for {
		select {
		case <-ctx.Done():
			return // Exit loop when context is done
		default:
			p.HandleCommand()
		}
	}
}
func main() {
	s := selection{}

	remoteDb := InitRemoteConnection()
	defer remoteDb.Close()
	log.Println("Remote database init")
	// localDb := InitLocalDatabase()
	// defer localDb.Close()
	// log.Println("Local database init")

	log.Printf("Checking for available databases...")
	databases := QueryFor(remoteDb, SHOW_DB_QUERY)
	s.databases = PromptSelection("Choose database", databases)
	if len(s.databases) < 1 {
		log.Panic("Error: no selection made")
	}
	log.Printf("You chose %s", s.databases[0])

	table := QueryFor(remoteDb, SHOW_TABLE_QUERY(s.databases[0]))
	s.tables = PromptSelection("Choose table", table)
	if len(s.tables) < 1 {
		log.Panic("Error: no selection made")
	}
	log.Printf("You chose %s", s.tables[0])

	// ReplaceDB(localDb, s.databases[0])

	// create all referencing tables in localDb
	foreignTables := QueryForTwoColumns(remoteDb, FETCH_GRAPH_EDGES(s.databases[0], s.tables[0])) // columns[0][0] should be primary key

	// getting heirarchical ordering
	inverseTopologicalOrdering, _ := topologicalSort(foreignTables)
	// TODO: move this code to manip service
	for i, v := range inverseTopologicalOrdering {
		inverseTopologicalOrdering[i] = fmt.Sprintf("'%s'", v)
	}
	topoString := strings.Join(inverseTopologicalOrdering, ",")
	inParam := fmt.Sprintf("(%s)", topoString)
	estimated, _ := SelectOneDynamic(remoteDb, FETCH_TABLES_SIZES(s.databases[0], inParam))[0].(float64)
	MAX := 0.05
	if estimated > MAX {
		log.Panicf("Error: total tables size %f GB exceeds %f GB", estimated, MAX)
	} else {
		fmt.Printf("Estimated replication size: %f", estimated)
	}

	return

	// appenc foreign tables to table slice
	for _, tableName := range inverseTopologicalOrdering {
		s.tables = append(s.tables, tableName)
	}
	// for _, table := range s.tables {
	// 	log.Printf("Replicating %s", table)
	// 	// get data to create template
	// 	createCommand := QueryForTwoColumns(remoteDb, SHOW_CREATE(s.databases[0], table))[0][1]
	// 	columns := QueryForTwoColumns(remoteDb, SELECT_COLUMNS(table))
	//
	// 	// log.Println(createCommand)
	// 	// log.Println(columns)
	// 	// form the select query that results in inserts
	// 	insertTemplate := CreateSelectInsertionFromSchema(s.databases[0], table, columns)
	// 	// get an insert for each row
	// 	inserts := QueryFor(remoteDb, insertTemplate)
	// 	Populate(localDb, s.databases[0], createCommand, inserts)
	// }
	// // close db as were going to open it again in raw tcp form
	// localDb.Close()
	//
	// // start proxying
	// socket, err := net.Listen("tcp", "127.0.0.1:3307")
	// if err != nil {
	// 	log.Fatalf("failed to start proxy: %s", err.Error())
	// }
	// fmt.Printf("Listening on localhost:%d\n", 3307)
	// for {
	// 	originSocket, err := socket.Accept()
	// 	if err != nil {
	// 		log.Fatalf("failed to accept connection: %s", err.Error())
	// 	}
	// 	go handleConn(originSocket, s.databases[0], s.tables[0])
	// }
	//
}

const COM_QUERY = byte(0x03)
