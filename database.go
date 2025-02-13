package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/dolthub/go-mysql-server/driver"
	_ "github.com/go-mysql-org/go-mysql/driver"
)

var (
	dbName    = "mydb"
	tableName = "mytable"
	address   = "localhost"
	port      = "3306"
	host      = os.Getenv("DB_HOST")
	user      = os.Getenv("DB_USER")
	pass      = os.Getenv("DB_PASS")
)

func InitializeDatabase(user, pass, host, port, dbname, driverName string) *sql.DB {
	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pass, host, port, dbname)
	fmt.Printf("Connecting to %s...\n", url)
	db, err := sql.Open(
		driverName,
		url,
	)
	if err != nil {
		fmt.Println("Error while connecting to database")
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		fmt.Println("Error while pinging database")
		panic(err)
	}
	log.Println("Database was successfully connected to")
	return db
}

func InitOverseerConnection() *sql.DB {
	sql.Register("sqle", driver.New(factory{}, nil))
	// create connection to ask user what should be imposed
	return InitializeDatabase(user, pass, host, port, "", "sqle")
}

func InitEmptyDatabase() *sql.DB {
	db, err := sql.Open("sqle", "")
	if err != nil {
		fmt.Println("Error while connecting to database")
		panic(err)
	}
	log.Println("Database provider init")
	return db

}

func QueryForPropety(c *sql.DB, query string) []string {
	props := []string{}
	return props
}
