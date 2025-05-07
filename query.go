package main

import "fmt"

var (
	SHOW_DB_QUERY        = "show databases;"
	FETCH_FOREIGN_TABLES = func(table, column string) string {
		return fmt.Sprintf(`
			SELECT TABLE_NAME, TABLE_SCHEMA
			FROM information_schema.KEY_COLUMN_USAGE
			WHERE REFERENCED_TABLE_NAME = '%s'
			AND REFERENCED_COLUMN_NAME = '%s';
			`, table, column)
	}
	SHOW_TABLE_QUERY = func(db string) string {
		return fmt.Sprintf(`
			SELECT TABLE_NAME
			FROM INFORMATION_SCHEMA.TABLES 
			WHERE TABLE_TYPE = 'BASE TABLE' AND TABLE_SCHEMA = '%s'
			AND TABLE_ROWS < 300 AND TABLE_ROWS > 35;
			`, db)
	}
	SHOW_CREATE = func(dbName, tableName string) string {
		return fmt.Sprintf("SHOW CREATE TABLE %s.%s;", dbName, tableName)
	}
	SELECT_COLUMNS = func(tableName string) string {
		return fmt.Sprintf("SELECT DISTINCT COLUMN_NAME, DATA_TYPE FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='%s';", tableName)
	}
	DROP_DB = func(dbName string) string {
		return fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName)
	}
	CREATE_DB = func(dbName string) string {
		return fmt.Sprintf("CREATE DATABASE %s", dbName)
	}
	USE_DB = func(dbName string) string {
		return fmt.Sprintf("USE %s", dbName)
	}
)
