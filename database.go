package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-mysql-org/go-mysql/client"
)

var conn *client.Conn

func InitializeDatabase() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	// TODO: ADD VALIDATION
	_conn, err := client.Connect(fmt.Sprintf("%s:%s", host, port), user, pass, "ACO_MS_DB")
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.

	log.Fatalln("Database was successfully connected to")

	conn = _conn
}

func CloseDB() error {
	return conn.Close()
}

func GetDatabase() *client.Conn {
	if conn == nil {
		log.Fatalf("Error: database not initialized")
	}
	return conn
}
