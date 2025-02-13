package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/dolthub/go-mysql-server/driver"
	_ "github.com/go-sql-driver/mysql"
)

var (
	port   = "3306"
	host   = os.Getenv("DB_HOST")
	user   = os.Getenv("DB_USER")
	pass   = os.Getenv("DB_PASS")
	dbName = os.Getenv("DB_NAME")
)

func InitializeDatabase(user, pass, host, port, dbname string) *sql.DB {
	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pass, host, port, dbname)
	fmt.Printf("Connecting to %s...\n", url)
	db, err := sql.Open(
		"mysql",
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
	// create connection to ask user what should be imposed
	return InitializeDatabase(user, pass, host, port, dbName)
}

func InitEmptyDatabase() *sql.DB {
	sql.Register("sqle", driver.New(factory{}, nil))
	db, err := sql.Open("sqle", "")
	if err != nil {
		fmt.Println("Error while connecting to database")
		panic(err)
	}
	log.Println("Database provider init")
	return db

}
func QueryFor(db *sql.DB, query string) []string {
	props := []string{}
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Error while connecting to database")
		panic(err)
	}
	for rows.Next() {
		var s string
		rows.Scan(&s)
		props = append(props, s)
	}
	return props
}
