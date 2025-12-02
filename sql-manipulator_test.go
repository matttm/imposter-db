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
	columns    [][2]string
	expected   string
	concatEx   string
}

func Test_CreateConcat(t *testing.T) {
	table := []CreateConcatTest{
		{
			schemaName: "A",
			tableName:  "PRD",
			columns:    [][2]string{{"PRD_CD", "varchar"}},
			expected:   "SELECT CONCAT('INSERT INTO A.PRD SET ', 'PRD_CD = ', IF(ISNULL(CAST(x.PRD_CD AS CHAR)), QUOTE('N'), QUOTE(x.PRD_CD)), ';' ) AS s FROM A.PRD x;",
			concatEx:   "INSERT INTO A.PRD SET PRD_CD = IF(ISNULL(CAST(x.PRD_CD AS CHAR)), QUOTE(N), QUOTE(x.PRD_CD)) ;",
		},
		{
			schemaName: "A",
			tableName:  "GATES",
			columns:    [][2]string{{"NAME", "varchar"}, {"OPEN", "varchar"}, {"CLOSE", "varchar"}},
			expected:   "SELECT CONCAT('INSERT INTO A.GATES SET ', 'NAME = ', IF(ISNULL(CAST(x.NAME AS CHAR)), QUOTE('N'), QUOTE(x.NAME)), ', ', 'OPEN = ', IF(ISNULL(CAST(x.OPEN AS CHAR)), QUOTE('N'), QUOTE(x.OPEN)), ', ', 'CLOSE = ', IF(ISNULL(CAST(x.CLOSE AS CHAR)), QUOTE('N'), QUOTE(x.CLOSE)), ';' ) AS s FROM A.GATES x;",
			concatEx:   "INSERT INTO A.GATES SET NAME = IF(ISNULL(CAST(x.NAME AS CHAR)), QUOTE(N), QUOTE(x.NAME)) , OPEN = IF(ISNULL(CAST(x.OPEN AS CHAR)), QUOTE(N), QUOTE(x.OPEN)) , CLOSE = IF(ISNULL(CAST(x.CLOSE AS CHAR)), QUOTE(N), QUOTE(x.CLOSE)) ;",
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
	re := regexp.MustCompile(`CONCAT\((.*)\)`)
	match := re.FindStringSubmatch(input)

	if match == nil {
		return nil, fmt.Errorf("no CONCAT function found")
	}

	paramsString := match[1]

	// Split the parameters string by commas, handling quoted strings and nested parentheses.
	var params []string
	inQuote := false
	parenDepth := 0
	currentParam := ""

	for _, r := range paramsString {
		switch r {
		case '\'':
			inQuote = !inQuote
			if !inQuote { // End of quote, add the quoted content
				params = append(params, currentParam)
				currentParam = ""
			}
		case '(':
			if !inQuote {
				parenDepth++
				currentParam += string(r)
			} else {
				currentParam += string(r)
			}
		case ')':
			if !inQuote {
				parenDepth--
				currentParam += string(r)
			} else {
				currentParam += string(r)
			}
		case ',':
			if !inQuote && parenDepth == 0 { // Comma outside quotes and top-level parentheses
				if strings.TrimSpace(currentParam) != "" {
					params = append(params, strings.TrimSpace(currentParam))
				}
				currentParam = ""
			} else {
				currentParam += string(r) // Comma inside quotes or nested parens, keep it
			}
		default:
			currentParam += string(r)
		}
	}
	if currentParam != "" { // Add the last parameter if not empty
		trimmed := strings.TrimSpace(currentParam)
		if trimmed != "" {
			params = append(params, trimmed)
		}
	}

	return params, nil
}
