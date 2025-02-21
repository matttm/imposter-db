package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type CreateConcatTest struct {
	schemaName string
	tableName  string
	columns    []string
	expected   string
}

func Test_CreateConcat(t *testing.T) {
	table := []CreateConcatTest{
		{
			schemaName: "A",
			tableName:  "PRD",
			columns:    []string{"PRD_CD"},
			expected:   "SELECT CONCAT('INSERT INTO ', 'A.PRD ', 'SET ', 'PRD_CD = ', x.PRD_CD, ';' ) AS s FROM A.PRD x;",
		},
		{
			schemaName: "A",
			tableName:  "GATES",
			columns:    []string{"NAME", "OPEN", "CLOSE"},
			expected:   "SELECT CONCAT('INSERT INTO ', 'A.GATES ', 'SET ', 'NAME = ', x.NAME, ', ', 'OPEN = ', x.OPEN, ', ', 'CLOSE = ', x.CLOSE, ';' ) AS s FROM A.GATES x;",
		},
	}
	for _, v := range table {
		assert.Equal(
			t,
			v.expected,
			CreateSelectInsertionFromSchema(v.schemaName, v.tableName, v.columns),
		)
	}
}
