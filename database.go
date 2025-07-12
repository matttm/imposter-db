package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var (
	port   = "3306"
	host   = os.Getenv("DB_HOST")
	user   = os.Getenv("DB_USER")
	pass   = os.Getenv("DB_PASS")
	dbName = os.Getenv("DB_NAME")
)

// InitializeDatabase establishes a connection to a MySQL database using the provided
// user credentials, host, port, and database name. It returns a pointer to the sql.DB
// instance if the connection is successful. The function will panic if it fails to
// connect or ping the database.
//
// Parameters:
//   - user:   The username for database authentication.
//   - pass:   The password for database authentication.
//   - host:   The hostname or IP address of the MySQL server.
//   - port:   The port number on which the MySQL server is listening.
//   - dbname: The name of the database to connect to.
//
// Returns:
//   - *sql.DB: A pointer to the established database connection.
func InitializeDatabase(user, pass, host, port, dbname string) *sql.DB {
	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pass, host, port, dbname)
	fmt.Printf("Connecting to %s...\n", url)
	db, err := sql.Open(
		"mysql",
		url,
	)
	if err != nil {
		fmt.Println("Error while connecting to database")
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		fmt.Println("Error while pinging database")
		panic(err)
	}
	log.Println("Database was successfully connected to")
	return db
}

// InitRemoteConnection initializes and returns a remote database connection.
// It prompts the user for connection parameters and establishes a connection
// using the provided credentials and database information.
//
// Returns:
//
//	*sql.DB: A pointer to the established database connection.
func InitRemoteConnection() *sql.DB {
	// create connection to ask user what should be imposed
	return InitializeDatabase(user, pass, host, port, dbName)
}

// InitLocalDatabase initializes and returns a connection to a local MySQL database
// using default credentials and connection parameters.
// It returns a pointer to the sql.DB instance.
// The database name is left empty, so the connection is established without selecting a specific database.
func InitLocalDatabase() *sql.DB {
	return InitializeDatabase("root", "mypassword", "127.0.0.1", "3306", "")
}

// QueryFor executes the provided SQL query on the given database connection,
// retrieves the first column of each resulting row as a string, and returns
// a slice containing these string values. If an error occurs during query
// execution, the function logs the error and panics.
//
// Parameters:
//
//	db    - The database connection to use for the query.
//	query - The SQL query string to execute.
//
// Returns:
//
//	A slice of strings containing the values from the first column of each row.
func QueryFor(db *sql.DB, query string) []string {
	props := []string{}
	log.Printf(query)
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Error querying database")
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var s string
		rows.Scan(&s)
		props = append(props, s)
	}
	return props
}

// QueryForTwoColumns executes the provided SQL query using the given database connection,
// expecting each row in the result set to contain exactly two columns of string values.
// It returns a slice of [2]string arrays, where each array represents a row with its two column values.
// If an error occurs during query execution, the function logs the query, prints an error message,
// and panics with the encountered error.
func QueryForTwoColumns(db *sql.DB, query string) [][2]string {
	props := [][2]string{}
	log.Printf(query)
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Error while connecting to database")
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		s := [2]string{}
		rows.Scan(&s[0], &s[1])
		props = append(props, s)
	}
	return props
}
func SelectOneDynamic(db *sql.DB, query string, params ...any) []any {
	log.Printf(query)
	row, err := db.Query(query, params...)
	if err != nil {
		fmt.Println("Error while querying database")
		panic(err)
	}
	defer row.Close()
	c, _ := row.Columns()
	n := len(c)
	s := make([]any, n)
	var _s float64
	row.Scan(&_s)
	fmt.Println(_s)
	return s
}

// ReplaceDB drops the specified database if it exists, recreates it, and switches the connection to use the new database.
// It sets the SQL mode to an empty string before performing these operations.
// If any step fails, the function panics with the encountered error.
func ReplaceDB(db *sql.DB, dbName string) {
	_, err := db.Exec("SET sql_mode=''")
	if err != nil {

		panic(err)
	}
	_, err = db.Exec(DROP_DB(dbName))
	if err != nil {
		fmt.Println("Error while dropping imposter database")
		panic(err)
	}
	_, err = db.Exec(CREATE_DB(dbName))
	if err != nil {
		fmt.Println("Error while creating imposter database")
		panic(err)
	}
	_, err = db.Exec(USE_DB(dbName))
	if err != nil {
		fmt.Println("Error while using database")
		panic(err)
	}
}

// Populate creates a table in the provided database using the given createQuery,
// then executes a series of insert statements to populate the table with data.
// If table creation fails, the function panics. Insert errors are logged and ignored.
// Parameters:
//   - db:       The database connection.
//   - dbName:   The name of the database (used for logging).
//   - createQuery: The SQL statement to create the table.
//   - inserts:  A slice of SQL insert statements to populate the table.
func Populate(db *sql.DB, dbName, createQuery string, inserts []string) {
	_, err := db.Exec(createQuery)
	log.Printf("INFO::CREATE_QUERY %s\n", createQuery)
	if err != nil {
		fmt.Printf("Error while creating spoofed table: %s\n", dbName)
		panic(err)
	}
	for _, ins := range inserts {
		// log.Println(ins)
		_, err = db.Exec(ins)
		if err != nil {
			fmt.Println("Error while inserting spoofed data")
			// there wrete some inserts that errored because they had bad data in db,
			// so it threw when afdded
			//
			// just ignore it
		}
	}
}
