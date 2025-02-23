package main

import (
	"fmt"
	"strings"
)

func CreateSelectInsertionFromSchema(schemaName, tableName string, columns []string) string {
	var sb strings.Builder
	write(&sb, `SELECT CONCAT('INSERT INTO ', '%s.%s ', 'SET ', `, schemaName, tableName)
	for i, v := range columns {
		write(
			&sb,
			`'%s = ', COALESCE(x.%s, "NULL")`,
			v,
			v,
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
