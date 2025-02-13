package main

import (
	"context"

	"github.com/dolthub/go-mysql-server/driver"
	"github.com/dolthub/go-mysql-server/memory"
	"github.com/dolthub/go-mysql-server/sql"
)

type factory struct{}

var inmemCtx *sql.Context = nil

func (factory) Resolve(name string, options *driver.Options) (string, sql.DatabaseProvider, error) {
	provider := memory.NewDBProvider(
		createTestDatabase(),
	)
	return name, provider, nil
}

func createTestDatabase() *memory.Database {
	const (
		dbName    = "mydb"
		tableName = "mytable"
	)

	db := memory.NewDatabase(dbName)
	pro := memory.NewDBProvider(db)
	inmemCtx = sql.NewContext(context.Background(), sql.WithSession(memory.NewSession(sql.NewBaseSession(), pro)))
	return db
}
