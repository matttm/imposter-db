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
	FETCH_GRAPH_EDGES = func(name string) string {
		return fmt.Sprintf(`
			SELECT
			PK_KCU.TABLE_NAME AS ParentTableName,
			FK_KCU.TABLE_NAME AS ChildTableName
			FROM
			INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS AS RC
			INNER JOIN
			INFORMATION_SCHEMA.KEY_COLUMN_USAGE AS FK_KCU
			ON RC.CONSTRAINT_SCHEMA = FK_KCU.CONSTRAINT_SCHEMA
			AND RC.CONSTRAINT_NAME = FK_KCU.CONSTRAINT_NAME
			INNER JOIN
			INFORMATION_SCHEMA.KEY_COLUMN_USAGE AS PK_KCU
			ON RC.UNIQUE_CONSTRAINT_SCHEMA = PK_KCU.CONSTRAINT_SCHEMA
			AND RC.UNIQUE_CONSTRAINT_NAME = PK_KCU.CONSTRAINT_NAME
			INNER JOIN
			INFORMATION_SCHEMA.TABLES AS PT_Tables -- Join to get size information for the parent table
			ON PK_KCU.TABLE_SCHEMA = PT_Tables.TABLE_SCHEMA
			AND PK_KCU.TABLE_NAME = PT_Tables.TABLE_NAME
			WHERE PK_KCU.TABLE_NAME = "%s";
			`, name)
	}
	FETCH_TABLES_SIZES = func(schema string, gtSize int) string {
		return fmt.Sprintf(`
			SELECT
			TABLE_SCHEMA,
			TABLE_NAME,
			DATA_LENGTH,
			INDEX_LENGTH,
			(DATA_LENGTH + INDEX_LENGTH) AS TotalBytes,
			(DATA_LENGTH + INDEX_LENGTH) / (1024.0 * 1024.0 * 1024.0) AS TotalGB
			FROM
			information_schema.TABLES
			WHERE
			TABLE_SCHEMA = '%s';
			`, schema)
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
