package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/dolthub/go-mysql-server/memory"
	"github.com/go-mysql-org/go-mysql/client"
	"github.com/go-mysql-org/go-mysql/mysql"
)

var (
	port = 3306
	host = os.Getenv("DB_HOST")
	user = os.Getenv("DB_USER")
	pass = os.Getenv("DB_PASS")
)

func handleConn(c net.Conn, provider *memory.DbProvider) {
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
	// create connection to ask user what should be imposed
	conductor, err := client.Connect(fmt.Sprintf("%s:%d", host, port), user, pass, "")
	if err != nil {
		panic(err)
	}
	log.Printf("Checking for available databases...")
	r, err := conductor.Execute(SHOW_DB_QUERY)
	if err != nil {
		panic(err)
	}
	//
	// Close result for reuse memory (it's not necessary but very useful)
	defer r.Close()

	// Handle resultset
	// v, _ := r.GetInt(0, 0)
	// v, _ = r.GetIntByName(0, "id")

	// Direct access to fields
	for _, row := range r.Values {
		for _, val := range row {
			// _ := val.Value() // interface{}
			// or
			if val.Type == mysql.FieldValueTypeString {
				log.Print(string(val.Value().([]uint8)))
			}
		}
	}

	// lets see what schemas are available

	// start proxying
	socket, err := net.Listen("tcp", "127.0.0.1:3307")
	if err != nil {
		log.Fatalf("failed to start proxy: %s", err.Error())
	}
	fmt.Printf("Listening on localhost:%d\n", 3307)
	// inputTables := []string{"ACO_MS_DB.APLCTN_RVW_PRD"}
	provider := InitEmptyDatabase()
	for {
		originSocket, err := socket.Accept()
		if err != nil {
			log.Fatalf("failed to accept connection: %s", err.Error())
		}
		go handleConn(originSocket, provider)
	}

}

const COM_QUERY = byte(0x03)
