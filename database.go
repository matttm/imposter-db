package main

import (
	"log"
	"net"
	"os"

	"github.com/go-mysql-org/go-mysql/server"
)

var conn *server.Conn

//	type ServerConfig struct {
//	    Host     string `env:"DB_HOST"`
//	    Port     int    `env:"DB_PORT"`
//	    User     string `env:"DB_USER"`
//	    Password string `env:"DB_PASSWORD"`
//	    Database string `env:"DB_DATABASE"`
//	}
func InitializeProxy(c net.Conn) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")

	_conn, err := server.NewConn(c, user, pass, server.EmptyHandler{})
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.

	log.Fatalln("Database was successfully connected to")

	conn = _conn
}

func CloseDB() {
	conn.Close()
}

func GetDatabase() *server.Conn {
	if conn == nil {
		log.Fatalf("Error: database not initialized")
	}
	return conn
}
