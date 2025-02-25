package main

import "fmt"

var (
	SHOW_DB_QUERY    = "show databases;"
	SHOW_TABLE_QUERY = func(db string) string {
		return fmt.Sprintf(`
			SELECT TABLE_NAME
			FROM INFORMATION_SCHEMA.TABLES 
			WHERE TABLE_TYPE = 'BASE TABLE' AND TABLE_SCHEMA = '%s'
			AND TABLE_ROWS < 200 AND TABLE_ROWS > 35;
			`, db)
	}
	SHOW_CREATE = func(dbName, tableName string) string {
		return fmt.Sprintf("SHOW CREATE TABLE %s.%s;", dbName, tableName)
	}
	SELECT_COLUMNS = func(tableName string) string {
		return fmt.Sprintf("SELECT COLUMN_NAME, DATA_TYPE FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='%s';", tableName)
	}
)
