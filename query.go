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
	FETCH_GRAPH_EDGES = func(sname, tname string) string {
		return fmt.Sprintf(`
			SELECT DISTINCT
			kcu.TABLE_NAME AS referencing_table_name,
			kcu.REFERENCED_TABLE_NAME AS referenced_table_name
			FROM
			INFORMATION_SCHEMA.KEY_COLUMN_USAGE AS kcu
			WHERE
			kcu.TABLE_SCHEMA = '%s'
			AND ((
			kcu.TABLE_NAME = '%s'
			AND kcu.REFERENCED_TABLE_NAME IS NOT NULL
			) OR (
			kcu.REFERENCED_TABLE_NAME = '%s'
			AND kcu.TABLE_NAME IS NOT NULL
			));
			`, sname, tname, tname)
	}
	FETCH_PARENT_GRAPH_EDGES = func(sname, tname string) string {
		return fmt.Sprintf(`
			SELECT DISTINCT
			kcu.REFERENCED_TABLE_NAME AS referenced_table_name,
			kcu.TABLE_NAME AS referencing_table_name
			FROM
			INFORMATION_SCHEMA.KEY_COLUMN_USAGE AS kcu
			WHERE
			kcu.TABLE_SCHEMA = '%s'
			AND kcu.TABLE_NAME = '%s'
			AND kcu.REFERENCED_TABLE_NAME IS NOT NULL
			`, sname, tname)
	}
	FETCH_TABLES_SIZES = func(schema, inArg string) string {
		return fmt.Sprintf(`
			SELECT
			SUM((DATA_LENGTH + INDEX_LENGTH) / (1024.0 * 1024.0 * 1024.0)) AS TotalGB
			FROM
			information_schema.TABLES t
			WHERE
			t.TABLE_SCHEMA = '%s'
			AND t.TABLE_NAME IN %s;
			`, schema, inArg)
	}
	SHOW_TABLE_QUERY = func(db string) string {
		return fmt.Sprintf(`
			SELECT TABLE_NAME
			FROM INFORMATION_SCHEMA.TABLES 
			WHERE TABLE_TYPE = 'BASE TABLE' AND TABLE_SCHEMA = '%s';`, db)
		// AND TABLE_ROWS < 300 AND TABLE_ROWS > 35;
		// `, db)
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
