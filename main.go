package main

import (
	"context"
	"flag"
	"fmt"
	"slices"
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
	p := protocol.InitializeProxy(c, localHost, schema, tableName, cancel, localUser, localPass)

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
	schemaFlag := flag.String("schema", "", "a string of the schema name")
	tableFlag := flag.String("table", "", "a string of the table name")
	fkFlag := flag.Bool("fk", false, "a boolean indicating whether foreign tables of the chosen table should be created")
	flag.Parse()

	remoteDb := InitRemoteConnection()
	defer remoteDb.Close()
	log.Println("Remote database init")
	localDb := InitLocalDatabase()
	defer localDb.Close()
	log.Println("Local database init")

	log.Printf("Checking for available databases...")
	databases := QueryFor(remoteDb, SHOW_DB_QUERY)
	if *schemaFlag == "" {
		s.databases = PromptSelection("Choose database", databases)
		if len(s.databases) < 1 {
			log.Panic("Error: no selection made")
		}
		if len(s.databases) > 1 {
			log.Panic("Error: one selection is currently supported")
		}
	} else {
		if !slices.Contains(databases, *schemaFlag) {
			panic("Fatal: provided schema is not visible on connection")
		}
		s.databases = []string{*schemaFlag}
	}
	log.Printf("You chose %s", s.databases[0])

	tables := QueryFor(remoteDb, SHOW_TABLE_QUERY(s.databases[0]))
	if *tableFlag == "" {
		s.tables = PromptSelection("Choose table", tables)
		if len(s.tables) < 1 {
			log.Panic("Error: no selection made")
		}
		if len(s.tables) > 1 {
			log.Panic("Error: one selection is currently supported")
		}
	} else {
		if !slices.Contains(tables, *tableFlag) {
			panic("Fatal: provided table is not visible on connection")
		}
		s.tables = []string{*tableFlag}
	}
	log.Printf("You chose %s", s.tables[0])

	ReplaceDB(localDb, s.databases[0])

	var foreignTables [][2]string
	if *fkFlag == false {
		// create all referencing tables in localDb
		// foreignTables = QueryForTwoColumns(remoteDb, FETCH_GRAPH_EDGES(s.databases[0], s.tables[0]))
		// just thid table
		foreignTables = [][2]string{{"", s.tables[0]}}
	} else {
		// copy all child tables
		foreignTables = QueryForTwoColumns(remoteDb, FETCH_PARENT_GRAPH_EDGES(s.databases[0], s.tables[0]))
	}
	// size check
	log.Printf("Starting topological sort: %v\n", foreignTables)
	// getting heirarchical ordering
	inverseTopologicalOrdering, _ := topologicalSort(foreignTables)
	// TODO: move this code to manip service
	var stringified []string
	for _, v := range inverseTopologicalOrdering {
		stringified = append(stringified, fmt.Sprintf("'%s'", v))
	}
	topoString := strings.Join(stringified, ",")
	inParam := fmt.Sprintf("(%s)", topoString)
	estimated := SelectOneDynamic(remoteDb, FETCH_TABLES_SIZES(s.databases[0], inParam))
	MAX := 0.05
	if *estimated > MAX {
		log.Panicf("Error: total tables size %f GB exceeds %f GB", *estimated, MAX)
		// log.Printf("Falling back to ignoring foreign keys")
	} else {
		log.Printf("Estimated replication size: %f", *estimated)
		s.tables = []string{}
		for _, tableName := range inverseTopologicalOrdering {
			s.tables = append(s.tables, tableName)
		}
	}

	// appenc foreign tables to table slice
	for _, table := range s.tables {

		// if table is empty, skip
		if len(table) == 0 {
			continue
		}
		log.Printf("Replicating %s", table)
		// get data to create template
		createCommand := QueryForTwoColumns(remoteDb, SHOW_CREATE(s.databases[0], table))[0][1]
		columns := QueryForTwoColumns(remoteDb, SELECT_COLUMNS(table))

		log.Println(createCommand)
		log.Println(columns)
		// form the select query that results in inserts
		insertTemplate := CreateSelectInsertionFromSchema(s.databases[0], table, columns)
		// get an insert for each row
		inserts := QueryFor(remoteDb, insertTemplate)
		Populate(localDb, s.databases[0], createCommand, inserts)
	}
	// close db as were going to open it again in raw tcp form
	localDb.Close()

	// start proxying
	// TODO: put in env vars
	socket, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%s", proxyPort))
	if err != nil {
		log.Fatalf("failed to start proxy: %s", err.Error())
	}
	fmt.Printf("Listening on localhost:%s\n", proxyPort)
	for {
		originSocket, err := socket.Accept()
		if err != nil {
			log.Fatalf("failed to accept connection: %s", err.Error())
		}
		go handleConn(originSocket, s.databases[0], s.tables[0])
	}

}

const COM_QUERY = byte(0x03)
