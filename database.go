package main

import (
	"context"
	"log"

	"github.com/dolthub/go-mysql-server/memory"
	"github.com/dolthub/go-mysql-server/sql"
)

var (
	dbName    = "mydb"
	tableName = "mytable"
	address   = "localhost"
	port      = 3306
)

var ctx *sql.Context = nil

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
