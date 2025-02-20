package main

import (
	"fmt"
	"strings"
)

func CreateConcatFromSchema(schemaName, tableName string, columns []string) string {
	var sb strings.Builder
	write(sb, `
		CONCAT(
		'INSERT INTO ',
		'%s.%s' x
		SET
	`, dbName, tableName)
	for i, v := range columns {
		write(
			sb,
			`x.%s = %s`,
			v,
			v,
		)
		if i < len(columns)-1 {
			write(sb, ", ")
		}
	}
	write(sb, `';'
		)
`)
	return sb.String()
}
func write(sb strings.Builder, template string, a ...any) (int, error) {
	s := fmt.Sprintf(template, a...)
	return sb.WriteString(s)
}
