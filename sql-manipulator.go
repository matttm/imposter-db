package main

import (
	"fmt"
	"strings"
)

// CreateTableStatement generates a SQL CREATE TABLE statement.
//
// Parameters:
//
//	dbName:    The name of the database (schema).
//	tableName: The name of the table to create.
//	columns:   A map where the key is the column name and the value is the SQL data type.
//
// Returns:
//
//	A string containing the CREATE TABLE statement, or an error if
//	there are no columns.
func CreateTableStatement(dbName, tableName string, columns map[string]string) (string, error) {
	if len(columns) == 0 {
		return "", fmt.Errorf("no columns provided for table %s.%s", dbName, tableName)
	}

	// Build the column definitions.
	columnDefinitions := make([]string, 0, len(columns))
	for columnName, columnType := range columns {
		//  Properly escape the column name.  In real-world code, this
		//  is crucial to prevent SQL injection.  This simple backtick
		//  escaping is suitable for MySQL, but other databases may
		//  require different escaping.  Consider using a library for
		//  database-specific escaping.
		escapedColumnName := "`" + columnName + "`"
		columnDefinitions = append(columnDefinitions, fmt.Sprintf("%s %s", escapedColumnName, columnType))
	}
	// Join the column definitions with commas.
	columnList := strings.Join(columnDefinitions, ", ")

	// Construct the full CREATE TABLE statement.
	//  Again, use backticks for the table name.
	escapedTableName := "`" + tableName + "`"
	statement := fmt.Sprintf("CREATE TABLE `%s`.%s (%s);", dbName, escapedTableName, columnList)
	return statement, nil
}

func CreateSelectInsertionFromSchema(schemaName, tableName string, columns [][2]string) string {
	var sb strings.Builder
	write(&sb, `SELECT CONCAT('INSERT INTO %s.%s SET ', `, schemaName, tableName)
	for i, v := range columns {
		name := v[0]
		_type := v[1]
		nullVal := getNullValue(_type)
		isNull := "x.%s"
		setVal := `x.%s`
		if nullVal != 0 {
			isNull = "CAST(x.%s AS CHAR)"
			setVal = "QUOTE(x.%s)"
		}
		isNull = fmt.Sprintf(isNull, name)
		setVal = fmt.Sprintf(setVal, name)
		write(
			&sb,
			`'%s = ', IF(ISNULL(%s), %v, %s)`,
			name,
			isNull,
			nullVal,
			setVal,
		)
		if i < len(columns)-1 {
			write(&sb, ", ', ', ")
		}
	}
	write(&sb, `, ';' ) AS s FROM %s.%s x;`, schemaName, tableName)
	return sb.String()
}
func write(sb *strings.Builder, template string, a ...any) (int, error) {
	s := fmt.Sprintf(template, a...)
	return sb.WriteString(s)
}
func getNullValue(t string) any {
	intTypes := []string{"int", "integer", "bigint", "tinyint", "smallint", "mediumint", "decimal", "numeric", "float", "double"}
	if contains(t, intTypes) {
		return 0
	}
	if t == "date" {
		return "QUOTE('1990-01-01')"
	}
	if t == "datetime" {
		return "QUOTE('1990-01-01 00:00:00')"
	}
	return `QUOTE('N')`
}

func contains(v string, a []string) bool {
	for _, i := range a {
		if i == v {
			return true
		}
	}
	return false
}
