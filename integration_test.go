package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testProxyPort      = "13306"
	testLocalPort      = "3306"
	testRemotePort     = "3307"
	testSchema         = "TEST_DB"
	testTable          = "application_gates"
	healthCheckDelay   = 2 * time.Second
	healthCheckRetries = 15
)

// TestIntegration_TableReplication tests that main() correctly replicates
// the chosen table from the remote database to the local database
func TestIntegration_TableReplication(t *testing.T) {
	// Skip if we're not in integration test mode
	if os.Getenv("INTEGRATION_TEST") == "" {
		t.Skip("Skipping integration test. Set INTEGRATION_TEST=1 to run")
	}

	// Setup environment variables for the test
	setupTestEnv()

	// Wait for databases to be healthy
	t.Log("Waiting for databases to be ready...")
	remoteDB := waitForDatabase(t, "ADMIN", "ADMIN", "127.0.0.1", testRemotePort, testSchema)
	require.NotNil(t, remoteDB, "Remote database should be accessible")
	defer remoteDB.Close()

	// Verify remote database has the expected test data
	t.Log("Verifying remote database has test data...")
	remoteRowCount := verifyRemoteData(t, remoteDB)

	// Get the original data from remote for comparison
	t.Log("Capturing remote table data for comparison...")
	remoteRows := captureTableData(t, remoteDB, testTable)
	require.Greater(t, len(remoteRows), 0, "Remote table should have data")

	// Clear local database to ensure clean state
	t.Log("Cleaning local database...")
	localDB := waitForDatabase(t, "root", "root", "127.0.0.1", testLocalPort, "")
	require.NotNil(t, localDB, "Local database should be accessible")
	cleanLocalDatabase(t, localDB, testSchema)
	localDB.Close()

	// Set up command line flags for main()
	t.Log("Setting up flags for main()...")
	os.Args = []string{
		"imposter-db",
		"-schema", testSchema,
		"-table", testTable,
		"-fk=false", // Don't include foreign key tables for simpler test
	}

	// Reset flag parsing
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Start main() in a goroutine since it runs forever
	t.Log("Starting main() to perform table replication...")
	mainStarted := make(chan bool)
	mainErrors := make(chan error, 1)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				mainErrors <- fmt.Errorf("main() panicked: %v", r)
			}
		}()

		// Signal that we're starting
		mainStarted <- true

		// This will run until the proxy starts listening
		// The proxy listener loop will block, which is expected
		main()
	}()

	// Wait for main to start
	select {
	case <-mainStarted:
		t.Log("main() started")
	case err := <-mainErrors:
		t.Fatalf("main() failed to start: %v", err)
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout waiting for main() to start")
	}

	// Give main() time to complete the table replication
	t.Log("Waiting for table replication to complete...")
	time.Sleep(10 * time.Second)

	// Now verify that the table was copied to local database
	t.Log("Verifying table was replicated to local database...")
	localDB = waitForDatabase(t, "root", "root", "127.0.0.1", testLocalPort, testSchema)
	require.NotNil(t, localDB, "Should be able to reconnect to local database")
	defer localDB.Close()

	// Test 1: Verify the table exists in local database
	t.Run("TableExistsInLocal", func(t *testing.T) {
		var tableExists int
		err := localDB.QueryRow(
			"SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = ? AND table_name = ?",
			testSchema, testTable,
		).Scan(&tableExists)
		require.NoError(t, err)
		assert.Equal(t, 1, tableExists, "Table should exist in local database")
	})

	// Test 2: Verify row count matches
	t.Run("RowCountMatches", func(t *testing.T) {
		var localRowCount int
		err := localDB.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", testTable)).Scan(&localRowCount)
		require.NoError(t, err)
		assert.Equal(t, remoteRowCount, localRowCount,
			"Local table should have same number of rows as remote (%d rows)", remoteRowCount)
	})

	// Test 3: Verify table schema matches
	t.Run("TableSchemaMatches", func(t *testing.T) {
		localColumns := getTableColumns(t, localDB, testTable)
		remoteColumns := getTableColumns(t, remoteDB, testTable)

		assert.Equal(t, len(remoteColumns), len(localColumns), "Should have same number of columns")

		for colName, remoteType := range remoteColumns {
			localType, exists := localColumns[colName]
			assert.True(t, exists, "Column %s should exist in local table", colName)
			assert.Equal(t, remoteType, localType, "Column %s type should match", colName)
		}
	})

	// Test 4: Verify actual data was copied correctly
	t.Run("DataCopiedCorrectly", func(t *testing.T) {
		localRows := captureTableData(t, localDB, testTable)

		assert.Equal(t, len(remoteRows), len(localRows),
			"Local and remote should have same number of data rows")

		// Verify some sample rows match
		for i := 0; i < len(remoteRows) && i < 5; i++ {
			remoteRow := remoteRows[i]
			found := false
			for _, localRow := range localRows {
				if mapsEqual(remoteRow, localRow) {
					found = true
					break
				}
			}
			assert.True(t, found, "Remote row %d should exist in local database", i)
		}
	})

	// Test 5: Verify local and remote databases are independent
	t.Run("DatabasesAreIndependent", func(t *testing.T) {
		// Insert a new row in local
		testGateName := fmt.Sprintf("LOCAL_TEST_%d", time.Now().Unix())
		_, err := localDB.Exec(
			fmt.Sprintf("INSERT INTO %s (gate_name, active_year, start_date, end_date) VALUES (?, ?, ?, ?)", testTable),
			testGateName, 2026, "2026-01-11", "2026-12-31",
		)
		require.NoError(t, err, "Should be able to insert into local database")

		// Verify it exists in local
		var localCount int
		err = localDB.QueryRow(
			fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE gate_name = ?", testTable),
			testGateName,
		).Scan(&localCount)
		require.NoError(t, err)
		assert.Equal(t, 1, localCount, "New row should exist in local database")

		// Verify it does NOT exist in remote
		var remoteCount int
		err = remoteDB.QueryRow(
			fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE gate_name = ?", testTable),
			testGateName,
		).Scan(&remoteCount)
		require.NoError(t, err)
		assert.Equal(t, 0, remoteCount, "New row should NOT exist in remote database")
	})

	t.Log("✅ All integration tests passed! Table was successfully replicated from remote to local.")
}

// setupTestEnv configures environment variables for the test
func setupTestEnv() {
	os.Setenv("REMOTE_DB_PORT", testRemotePort)
	os.Setenv("REMOTE_DB_HOST", "127.0.0.1")
	os.Setenv("REMOTE_DB_USER", "ADMIN")
	os.Setenv("REMOTE_DB_PASS", "ADMIN")
	os.Setenv("REMOTE_DB_NAME", testSchema)

	os.Setenv("PROXY_PORT", testProxyPort)

	os.Setenv("LOCAL_DB_PORT", testLocalPort)
	os.Setenv("LOCAL_DB_HOST", "127.0.0.1")
	os.Setenv("LOCAL_DB_USER", "root")
	os.Setenv("LOCAL_DB_PASS", "root")
	os.Setenv("LOCAL_DB_NAME", "")
}

// waitForDatabase waits for a database to be ready and returns a connection
func waitForDatabase(t *testing.T, user, pass, host, port, dbname string) *sql.DB {
	var db *sql.DB
	var err error

	for i := 0; i < healthCheckRetries; i++ {
		url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pass, host, port, dbname)
		db, err = sql.Open("mysql", url)
		if err == nil {
			err = db.Ping()
			if err == nil {
				return db
			}
		}

		if i < healthCheckRetries-1 {
			time.Sleep(healthCheckDelay)
		}
	}

	t.Fatalf("Failed to connect to database after %d retries: %v", healthCheckRetries, err)
	return nil
}

// verifyRemoteData checks that the remote database has the expected test data
// and returns the row count for the test table
func verifyRemoteData(t *testing.T, db *sql.DB) int {
	var count int
	err := db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", testTable)).Scan(&count)
	require.NoError(t, err, "Should be able to query application_gates")
	assert.Greater(t, count, 0, "Remote database should have application gates data")

	err = db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	require.NoError(t, err, "Should be able to query users")
	assert.Greater(t, count, 0, "Remote database should have users data")

	// Return the count for the test table
	var testTableCount int
	err = db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", testTable)).Scan(&testTableCount)
	require.NoError(t, err)
	return testTableCount
}

// cleanLocalDatabase removes the test schema if it exists
func cleanLocalDatabase(t *testing.T, db *sql.DB, schema string) {
	_, err := db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", schema))
	require.NoError(t, err, "Should be able to drop test schema")
	t.Logf("Cleaned local database (dropped %s schema if it existed)", schema)
}

// captureTableData retrieves all rows from a table as a slice of maps
func captureTableData(t *testing.T, db *sql.DB, table string) []map[string]interface{} {
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s ORDER BY gate_id", table))
	require.NoError(t, err)
	defer rows.Close()

	cols, err := rows.Columns()
	require.NoError(t, err)

	var result []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(cols))
		valuePtrs := make([]interface{}, len(cols))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		err = rows.Scan(valuePtrs...)
		require.NoError(t, err)

		row := make(map[string]interface{})
		for i, col := range cols {
			row[col] = values[i]
		}
		result = append(result, row)
	}

	return result
}

// getTableColumns returns a map of column names to their data types
func getTableColumns(t *testing.T, db *sql.DB, table string) map[string]string {
	rows, err := db.Query(fmt.Sprintf("DESCRIBE %s", table))
	require.NoError(t, err)
	defer rows.Close()

	columns := make(map[string]string)
	for rows.Next() {
		var field, colType, null, key, extra string
		var defaultVal interface{}
		err := rows.Scan(&field, &colType, &null, &key, &defaultVal, &extra)
		require.NoError(t, err)
		columns[field] = colType
	}

	return columns
}

// mapsEqual compares two maps for equality
func mapsEqual(a, b map[string]interface{}) bool {
	if len(a) != len(b) {
		return false
	}

	for key, aVal := range a {
		bVal, exists := b[key]
		if !exists {
			return false
		}

		// Handle byte slices (common for date/time fields)
		aBytes, aIsBytes := aVal.([]byte)
		bBytes, bIsBytes := bVal.([]byte)
		if aIsBytes && bIsBytes {
			if string(aBytes) != string(bBytes) {
				return false
			}
			continue
		}

		// Direct comparison for other types
		if fmt.Sprintf("%v", aVal) != fmt.Sprintf("%v", bVal) {
			return false
		}
	}

	return true
}
