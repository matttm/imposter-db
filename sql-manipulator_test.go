package main

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CreateSelectInsertionFromSchema(t *testing.T) {
	tests := []struct {
		name       string
		schemaName string
		tableName  string
		columns    [][2]string
		expected   string
	}{
		{
			name:       "Single INT column",
			schemaName: "A",
			tableName:  "PRD",
			columns: [][2]string{
				{"PRD_CD", "INT"},
			},
			expected: "SELECT CONCAT('INSERT INTO A.PRD SET ', 'PRD_CD = ', IF(ISNULL(x.PRD_CD), 0, x.PRD_CD), ';' ) AS s FROM A.PRD x;",
		},
		{
			name:       "Multiple VARCHAR and TIME columns",
			schemaName: "A",
			tableName:  "GATES",
			columns: [][2]string{
				{"NAME", "VARCHAR"},
				{"OPEN", "TIME"},
				{"CLOSE", "TIME"},
			},
			expected: "SELECT CONCAT('INSERT INTO A.GATES SET ', 'NAME = ', IF(ISNULL(x.NAME), QUOTE('N'), QUOTE(x.NAME)), ', ', 'OPEN = ', IF(ISNULL(x.OPEN), QUOTE('N'), QUOTE(x.OPEN)), ', ', 'CLOSE = ', IF(ISNULL(x.CLOSE), QUOTE('N'), QUOTE(x.CLOSE)), ';' ) AS s FROM A.GATES x;",
		},
		{
			name:       "Date and Datetime columns",
			schemaName: "B",
			tableName:  "DATES",
			columns: [][2]string{
				{"START_DATE", "date"},
				{"END_TIME", "datetime"},
			},
			expected: "SELECT CONCAT('INSERT INTO B.DATES SET ', 'START_DATE = ', IF(ISNULL(x.START_DATE), QUOTE('1990-01-01'), QUOTE(x.START_DATE)), ', ', 'END_TIME = ', IF(ISNULL(x.END_TIME), QUOTE('1990-01-01 00:00:00'), QUOTE(x.END_TIME)), ';' ) AS s FROM B.DATES x;",
		},
		{
			name:       "Mixed INT and VARCHAR columns",
			schemaName: "C",
			tableName:  "MIXED",
			columns: [][2]string{
				{"ID", "int"},
				{"DESC", "varchar"},
			},
			expected: "SELECT CONCAT('INSERT INTO C.MIXED SET ', 'ID = ', IF(ISNULL(x.ID), 0, x.ID), ', ', 'DESC = ', IF(ISNULL(x.DESC), QUOTE('N'), QUOTE(x.DESC)), ';' ) AS s FROM C.MIXED x;",
		},
		{
			name:       "Unknown type falls back to QUOTE('N')",
			schemaName: "D",
			tableName:  "UNKNOWN",
			columns: [][2]string{
				{"FOO", "BLOB"},
			},
			expected: "SELECT CONCAT('INSERT INTO D.UNKNOWN SET ', 'FOO = ', IF(ISNULL(x.FOO), QUOTE('N'), QUOTE(x.FOO)), ';' ) AS s FROM D.UNKNOWN x;",
		},
		{
			name:       "All numeric types",
			schemaName: "E",
			tableName:  "NUMS",
			columns: [][2]string{
				{"A", "int"},
				{"B", "bigint"},
				{"C", "float"},
				{"D", "double"},
			},
			expected: "SELECT CONCAT('INSERT INTO E.NUMS SET ', 'A = ', IF(ISNULL(x.A), 0, x.A), ', ', 'B = ', IF(ISNULL(x.B), 0, x.B), ', ', 'C = ', IF(ISNULL(x.C), 0, x.C), ', ', 'D = ', IF(ISNULL(x.D), 0, x.D), ';' ) AS s FROM E.NUMS x;",
		},
		{
			name:       "All string types",
			schemaName: "F",
			tableName:  "STRS",
			columns: [][2]string{
				{"A", "varchar"},
				{"B", "char"},
				{"C", "text"},
			},
			expected: "SELECT CONCAT('INSERT INTO F.STRS SET ', 'A = ', IF(ISNULL(x.A), QUOTE('N'), QUOTE(x.A)), ', ', 'B = ', IF(ISNULL(x.B), QUOTE('N'), QUOTE(x.B)), ', ', 'C = ', IF(ISNULL(x.C), QUOTE('N'), QUOTE(x.C)), ';' ) AS s FROM F.STRS x;",
		},
		{
			name:       "Date, Datetime, and Numeric",
			schemaName: "G",
			tableName:  "MIXED2",
			columns: [][2]string{
				{"D", "date"},
				{"DT", "datetime"},
				{"N", "int"},
			},
			expected: "SELECT CONCAT('INSERT INTO G.MIXED2 SET ', 'D = ', IF(ISNULL(x.D), QUOTE('1990-01-01'), QUOTE(x.D)), ', ', 'DT = ', IF(ISNULL(x.DT), QUOTE('1990-01-01 00:00:00'), QUOTE(x.DT)), ', ', 'N = ', IF(ISNULL(x.N), 0, x.N), ';' ) AS s FROM G.MIXED2 x;",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CreateSelectInsertionFromSchema(tt.schemaName, tt.tableName, tt.columns)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func extractConcatParams(input string) ([]string, error) {
	// Regular expression to match the CONCAT function and capture its parameters.
	// This regex handles potential variations in whitespace and quoting.
	re := regexp.MustCompile(`CONCAT\(([^)]*)\)`)
	match := re.FindStringSubmatch(input)

	if match == nil {
		return nil, fmt.Errorf("no CONCAT function found")
	}

	paramsString := match[1]

	// Split the parameters string by commas, handling quoted strings.
	var params []string
	inQuote := false
	currentParam := ""

	for _, r := range paramsString {
		switch r {
		case '\'':
			inQuote = !inQuote
			if !inQuote { // End of quote, trim and add
				params = append(params, currentParam)
				currentParam = ""
			}
		case ',':
			if !inQuote { // Comma outside quotes, add the parameter
				params = append(params, currentParam)
				currentParam = ""
			} else {
				currentParam += string(r) // Comma inside quotes, keep it
			}
		default:
			currentParam += string(r)
		}
	}
	if currentParam != "" { // Add the last parameter if not empty
		params = append(params, strings.TrimSpace(currentParam))
	}

	return params, nil
}
