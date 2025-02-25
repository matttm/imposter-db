package main

import (
	"fmt"
	"strings"
)

func CreateSelectInsertionFromSchema(schemaName, tableName string, columns [][2]string) string {
	var sb strings.Builder
	write(&sb, `SELECT CONCAT('INSERT INTO ', '%s.%s ', 'SET ', `, schemaName, tableName)
	for i, v := range columns {
		name := v[0]
		_type := v[1]
		nullVal := getNullValue(_type)
		isNull := "x.%s"
		setVal := `x.%s`
		if nullVal == `""` {
			isNull = "CAST(x.%s AS CHAR)"
			setVal = "CONCAT('\"', x.%s, '\"')"
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
	vals := []string{"INT", "INTEGER", "BIGINT", "TINYINT", "SMALLINT", "MEDIUMINT", "DECIMAL", "NUMERIC", "FLOAT", "DOUBLE"}
	if contains(t, vals) {
		return 0
	}
	return `""`
}

func contains(v string, a []string) bool {
	for _, i := range a {
		if i == v {
			return true
		}
	}
	return false
}
