package main

import (
	"fmt"
	"log"
	"strings"
)

func CreateSelectInsertionFromSchema(schemaName, tableName string, columns [][2]string) string {
	var sb strings.Builder
	write(&sb, `SELECT CONCAT('INSERT INTO %s.%s SET ', `, schemaName, tableName)
	for i, v := range columns {
		name := v[0]
		typ := strings.ToLower(v[1])
		var expr string
		switch typ {
		case "int", "integer", "bigint", "tinyint", "smallint", "mediumint", "decimal", "numeric", "float", "double":
			expr = fmt.Sprintf("'%s = ', IF(ISNULL(x.%s), 0, x.%s)", name, name, name)
		case "date":
			expr = fmt.Sprintf("'%s = ', IF(ISNULL(x.%s), QUOTE('1990-01-01'), QUOTE(x.%s))", name, name, name)
		case "datetime":
			expr = fmt.Sprintf("'%s = ', IF(ISNULL(x.%s), QUOTE('1990-01-01 00:00:00'), QUOTE(x.%s))", name, name, name)
		case "varchar", "char", "text", "time":
			expr = fmt.Sprintf("'%s = ', IF(ISNULL(x.%s), QUOTE('N'), QUOTE(x.%s))", name, name, name)
		default:
			expr = fmt.Sprintf("'%s = ', IF(ISNULL(x.%s), QUOTE('N'), QUOTE(x.%s))", name, name, name)
		}
		write(&sb, expr)
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
	log.Println(t)
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
