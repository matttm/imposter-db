package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestCreateSelectInsertionFromSchema(t *testing.T) {
	tests := []struct {
		name     string
		db       string
		table    string
		columns  [][2]string
		expected string
	}{
		{
			name:    "Basic",
			db:      "mydb",
			table:   "users",
			columns: [][2]string{{"id", "int"}, {"name", "varchar"}, {"created_at", "datetime"}},
			expected: "SELECT CONCAT('INSERT INTO mydb.users SET ', " +
				"id = IF(ISNULL(x.id), 0, CAST(x.id AS CHAR)), ', ', " +
				"name = IF(ISNULL(x.name), QUOTE('N'), QUOTE(x.name)), ', ', " +
				"created_at = IF(ISNULL(x.created_at), QUOTE('1990-01-01 00:00:00'), QUOTE(x.created_at)), ';' ) AS s FROM mydb.users x;",
		},
		{
			name:     "EmptyColumns",
			db:       "db",
			table:    "tbl",
			columns:  [][2]string{},
			expected: "SELECT CONCAT('INSERT INTO db.tbl SET ', , ';' ) AS s FROM db.tbl x;",
		},
		{
			name:    "NullHandling",
			db:      "s",
			table:   "t",
			columns: [][2]string{{"score", "int"}, {"nickname", "varchar"}, {"birthdate", "date"}},
			expected: "SELECT CONCAT('INSERT INTO s.t SET ', " +
				"score = IF(ISNULL(x.score), 0, CAST(x.score AS CHAR)), ', ', " +
				"nickname = IF(ISNULL(x.nickname), QUOTE('N'), QUOTE(x.nickname)), ', ', " +
				"birthdate = IF(ISNULL(x.birthdate), QUOTE('1990-01-01'), QUOTE(x.birthdate)), ';' ) AS s FROM s.t x;",
		},
		{
			name:    "SingleColumn",
			db:      "a",
			table:   "b",
			columns: [][2]string{{"active", "tinyint"}},
			expected: "SELECT CONCAT('INSERT INTO a.b SET ', " +
				"active = IF(ISNULL(x.active), 0, CAST(x.active AS CHAR)), ';' ) AS s FROM a.b x;",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CreateSelectInsertionFromSchema(tt.db, tt.table, tt.columns)
			if got != tt.expected {
				t.Errorf("expected:\n%q\ngot:\n%q", tt.expected, got)
			}
		})
	}
}
func TestCreateTableStatement(t *testing.T) {
	tests := []struct {
		name      string
		dbName    string
		tableName string
		columns   map[string]string
		want      string
		wantErr   bool
	}{
		{
			name:      "Basic",
			dbName:    "mydb",
			tableName: "users",
			columns: map[string]string{
				"id":   "INT",
				"name": "VARCHAR(255)",
			},
			want:    "CREATE TABLE `mydb`.`users` (`id` INT, `name` VARCHAR(255));",
			wantErr: false,
		},
		{
			name:      "EmptyColumns",
			dbName:    "db",
			tableName: "tbl",
			columns:   map[string]string{},
			want:      "",
			wantErr:   true,
		},
		{
			name:      "SpecialCharsInNames",
			dbName:    "schema",
			tableName: "table-name",
			columns: map[string]string{
				"col 1": "TEXT",
				"c@l2":  "BLOB",
			},
			want:    "CREATE TABLE `schema`.`table-name` (`col 1` TEXT, `c@l2` BLOB);",
			wantErr: false,
		},
		{
			name:      "SingleColumn",
			dbName:    "a",
			tableName: "b",
			columns: map[string]string{
				"only": "BOOLEAN",
			},
			want:    "CREATE TABLE `a`.`b` (`only` BOOLEAN);",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateTableStatement(tt.dbName, tt.tableName, tt.columns)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateTableStatement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// Since map iteration order is random, only check for error or that all expected parts are present
			if !tt.wantErr {
				for k, v := range tt.columns {
					colDef := fmt.Sprintf("`%s` %s", k, v)
					if !strings.Contains(got, colDef) {
						t.Errorf("expected column definition %q in result %q", colDef, got)
					}
				}
				if !strings.HasPrefix(got, fmt.Sprintf("CREATE TABLE `%s`.`%s` (", tt.dbName, tt.tableName)) {
					t.Errorf("expected prefix not found in result: %q", got)
				}
				if !strings.HasSuffix(got, ");") {
					t.Errorf("expected suffix not found in result: %q", got)
				}
			} else if got != "" {
				t.Errorf("expected empty string on error, got %q", got)
			}
		})
	}
}
