package main

import (
	"context"
	"flag"
	"fmt"
	"slices"

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

	// NEW FLOW: Setup Federated tables
	// Create fresh local database
	ReplaceDB(localDb, s.databases[0])

	selectedTable := s.tables[0]
	log.Printf("\n=== Setting up Federated Architecture ===")
	log.Printf("Selected table for local modifications: %s", selectedTable)

	// Step 1: Create the selected table locally (non-federated)
	log.Println("\n--- Creating local table (non-federated) ---")
	selectedTableCreateCmd := QueryForTwoColumns(remoteDb, SHOW_CREATE(s.databases[0], selectedTable))[0][1]
	err := CreateLocalTableSchema(localDb, selectedTableCreateCmd)
	if err != nil {
		log.Fatalf("Failed to create local table: %v", err)
	}
	log.Printf("✓ Local table created: %s", selectedTable)

	// Step 2: Create Federated tables for all other tables
	log.Println("\n--- Creating Federated tables for other tables ---")
	allTables := QueryFor(remoteDb, SHOW_TABLE_QUERY(s.databases[0]))
	for _, table := range allTables {
		if table == selectedTable {
			continue // Skip the selected table, it's already created
		}
		log.Printf("Creating Federated table: %s", table)
		err := CreateFederatedTable(localDb, table, s.databases[0], s.databases[0])
		if err != nil {
			log.Printf("Warning: Failed to create federated table %s: %v", table, err)
			// Continue with other tables even if one fails
		}
	}

	log.Printf("\n=== Federated Architecture Setup Complete ===")
	log.Printf("• Local table (editable): %s.%s", s.databases[0], selectedTable)
	log.Printf("• Federated tables (read from remote): %d tables", len(allTables)-1)

	// Close connections
	remoteDb.Close()
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
