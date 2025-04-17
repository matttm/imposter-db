package protocol

// doc https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_com_query_response_text_resultset_column_definition.html#sect_protocol_com_query_response_text_resultset_column_definition_41
//
// ColumnDefinition41 represents the structure for defining a column, version 41 of the protocol.
type ColumnDefinition41 struct {
	Catalog                   string  // The catalog used. Currently always "def" - represented as a normal string
	Schema                    string  // Schema name
	Table                     string  // Virtual table name
	OrgTable                  string  // Physical table name
	Name                      string  // Virtual column name
	OrgName                   string  // Physical column name
	LengthOfFixedLengthFields uint8   // Length of fixed length fields [0x0c] -  using uint8
	CharacterSet              uint16  // The column character set as defined in Character Set
	ColumnLength              uint32  // Maximum length of the field
	Type                      uint8   // Type of the column as defined in enum_field_types - using uint8
	Flags                     uint16  // Flags as defined in Column Definition Flags
	Decimals                  uint8   // Max shown decimal digits: 0x00 for integers and static strings, 0x1f for dynamic strings, double, float, 0x00 to 0x51 for decimals
	Reserved                  [2]byte // Reserved for future use.
	DefaultValue              *string // Default value (NULL if 0xFB) - using a pointer to a string
}
type Column *string

// doc https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_com_query_response_text_resultset_row.html
type ResultsetRow struct {
	columns []Column
}
