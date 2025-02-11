package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/dolthub/go-mysql-server/memory"
	"github.com/dolthub/go-mysql-server/sql"
	"github.com/go-mysql-org/go-mysql/client"
	"github.com/go-mysql-org/go-mysql/mysql"
)

var (
	dbName    = "mydb"
	tableName = "mytable"
	address   = "localhost"
	port      = 3306
	host      = os.Getenv("DB_HOST")
	user      = os.Getenv("DB_USER")
	pass      = os.Getenv("DB_PASS")
)

var ctx *sql.Context = nil

func InitOverseerConnection() *client.Conn {
	// create connection to ask user what should be imposed
	conductor, err := client.Connect(fmt.Sprintf("%s:%d", host, port), user, pass, "")
	if err != nil {
		panic(err)
	}
	return conductor
}

func InitEmptyDatabase() *memory.DbProvider {
	pro := createTestDatabase()
	// engine := sqle.NewDefault(pro)
	// session := memory.NewSession(sql.NewBaseSession(), pro)
	log.Println("Database provider init")
	return pro

}

func createTestDatabase() *memory.DbProvider {
	db := memory.NewDatabase(dbName)
	db.BaseDatabase.EnablePrimaryKeyIndexes()

	pro := memory.NewDBProvider(db)
	session := memory.NewSession(sql.NewBaseSession(), pro)
	ctx = sql.NewContext(context.Background(), sql.WithSession(session))
	return pro
}

func QueryForPropety(c *client.Conn, query string) []string {
	r, err := c.Execute(SHOW_DB_QUERY)
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

}
