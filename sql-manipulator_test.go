package main

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type CreateConcatTest struct {
	schemaName string
	tableName  string
	columns    []string
	expected   string
	concatEx   string
}

func Test_CreateConcat(t *testing.T) {
	table := []CreateConcatTest{
		{
			schemaName: "A",
			tableName:  "PRD",
			columns:    []string{"PRD_CD"},
			expected:   "SELECT CONCAT('INSERT INTO ', 'A.PRD ', 'SET ', 'PRD_CD = ', x.PRD_CD, ';' ) AS s FROM A.PRD x;",
			concatEx:   "INSERT INTO A.PRD SET PRD_CD = x.PRD_CD ;",
		},
		{
			schemaName: "A",
			tableName:  "GATES",
			columns:    []string{"NAME", "OPEN", "CLOSE"},
			expected:   "SELECT CONCAT('INSERT INTO ', 'A.GATES ', 'SET ', 'NAME = ', x.NAME, ', ', 'OPEN = ', x.OPEN, ', ', 'CLOSE = ', x.CLOSE, ';' ) AS s FROM A.GATES x;",
			concatEx:   "INSERT INTO A.GATES SET NAME = x.NAME , OPEN = x.OPEN , CLOSE = x.CLOSE ;",
		},
	}
	for _, v := range table {
		output := CreateSelectInsertionFromSchema(v.schemaName, v.tableName, v.columns)
		assert.Equal(
			t,
			v.expected,
			output,
		)
		extracted, err := extractConcatParams(output)
		if err != nil {
			panic(err)
		}
		space := regexp.MustCompile(`\s+`)
		assert.Equal(
			t,
			space.ReplaceAllString(strings.Join(extracted, ""), " "),
			v.concatEx,
		)
	}
}

//	func extractParenthesesContent(input string) (string, error) {
//		openParenIndex := strings.Index(input, "(")
//		closeParenIndex := strings.Index(input, ")")
//
//		if openParenIndex == -1 || closeParenIndex == -1 || openParenIndex >= closeParenIndex {
//			return "", fmt.Errorf("invalid input format: missing or misplaced parentheses")
//		}
//
//		content := input[openParenIndex+1 : closeParenIndex]
//		return content, nil
//	}
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
