package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

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

func InitRemoteConnection() *sql.DB {
	// create connection to ask user what should be imposed
	return InitializeDatabase(user, pass, host, port, dbName)
}

func InitLocalDatabase() *sql.DB {
	return InitializeDatabase("root", "mypassword", "127.0.0.1", "3306", "")
}
func QueryFor(db *sql.DB, query string) []string {
	props := []string{}
	log.Printf(query)
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Error querying database")
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var s string
		rows.Scan(&s)
		props = append(props, s)
	}
	return props
}
func QueryForTwoColumns(db *sql.DB, query string) [][2]string {
	props := [][2]string{}
	log.Printf(query)
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Error while connecting to database")
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		s := [2]string{}
		rows.Scan(&s[0], &s[1])
		props = append(props, s)
	}
	return props
}
func ReplaceDB(db *sql.DB, dbName string) {
	_, err := db.Exec("SET sql_mode=''")
	if err != nil {

		panic(err)
	}
	_, err = db.Exec(DROP_DB(dbName))
	if err != nil {
		fmt.Println("Error while dropping imposter database")
		panic(err)
	}
	_, err = db.Exec(CREATE_DB(dbName))
	if err != nil {
		fmt.Println("Error while creating imposter database")
		panic(err)
	}
	_, err = db.Exec(USE_DB(dbName))
	if err != nil {
		fmt.Println("Error while using database")
		panic(err)
	}
}
func Populate(db *sql.DB, dbName, createQuery string, inserts []string) {
	_, err := db.Exec(createQuery)
	log.Printf("INFO::CREATE_QUERY %s\n", createQuery)
	if err != nil {
		fmt.Printf("Error while creating spoofed table: %s\n", dbName)
		panic(err)
	}
	for _, ins := range inserts {
		// log.Println(ins)
		_, err = db.Exec(ins)
		if err != nil {
			fmt.Println("Error while inserting spoofed data")
			// there wrete some inserts that errored because they had bad data in db,
			// so it threw when afdded
			//
			// just ignore it
		}
	}
}
