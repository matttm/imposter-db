package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
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
	testTimeout        = 30 * time.Second
	healthCheckDelay   = 2 * time.Second
	healthCheckRetries = 15
)

// TestIntegration_ProxyWithDockerContainers tests the full proxy functionality
// using the docker-compose containers (localdb on 3306, remote db on 3307)
func TestIntegration_ProxyWithDockerContainers(t *testing.T) {
	// Skip if we're not in integration test mode
	if os.Getenv("INTEGRATION_TEST") == "" {
		t.Skip("Skipping integration test. Set INTEGRATION_TEST=1 to run")
	}

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	// Setup environment variables for the test
	setupTestEnv()

	// Wait for databases to be healthy
	t.Log("Waiting for databases to be ready...")
	localDB := waitForDatabase(t, "root", "root", "127.0.0.1", testLocalPort, "")
	require.NotNil(t, localDB, "Local database should be accessible")
	defer localDB.Close()

	remoteDB := waitForDatabase(t, "ADMIN", "ADMIN", "127.0.0.1", testRemotePort, testSchema)
	require.NotNil(t, remoteDB, "Remote database should be accessible")
	defer remoteDB.Close()

	// Verify remote database has the expected test data
	t.Log("Verifying remote database setup...")
	verifyRemoteData(t, remoteDB)

	// Setup the local database with replicated data
	t.Log("Setting up local database replication...")
	setupLocalReplication(t, localDB, remoteDB)

	// Start the proxy server
	t.Log("Starting proxy server...")
	proxyReady := make(chan bool)
	proxyErrors := make(chan error, 1)
	go startTestProxy(ctx, proxyReady, proxyErrors)

	// Wait for proxy to be ready
	select {
	case <-proxyReady:
		t.Log("Proxy server is ready")
	case err := <-proxyErrors:
		t.Fatalf("Failed to start proxy: %v", err)
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout waiting for proxy to start")
	}

	// Give proxy a moment to stabilize
	time.Sleep(500 * time.Millisecond)

	// Connect to the proxy
	t.Log("Connecting to proxy...")
	proxyDB := connectToProxy(t)
	require.NotNil(t, proxyDB, "Should be able to connect to proxy")
	defer proxyDB.Close()

	// Test 1: Verify we can see tables from remote database through proxy
	t.Run("ProxyShowsTables", func(t *testing.T) {
		tables := getTables(t, proxyDB, testSchema)
		assert.Contains(t, tables, testTable, "Proxy should show application_gates table")
		assert.Contains(t, tables, "users", "Proxy should show users table from remote")
		assert.Contains(t, tables, "applications", "Proxy should show applications table from remote")
	})

	// Test 2: Verify spoofed table reads from local database
	t.Run("SpoofedTableReadsFromLocal", func(t *testing.T) {
		// Modify data in local database
		_, err := localDB.Exec(fmt.Sprintf("USE %s", testSchema))
		require.NoError(t, err)

		// Insert a unique test row in local DB
		testGateName := fmt.Sprintf("TEST_GATE_%d", time.Now().Unix())
		_, err = localDB.Exec(
			fmt.Sprintf("INSERT INTO %s (gate_name, active_year, start_date, end_date) VALUES (?, ?, ?, ?)", testTable),
			testGateName, 2026, "2026-01-01", "2026-12-31",
		)
		require.NoError(t, err, "Should be able to insert into local database")

		// Read from proxy - should see the local data
		time.Sleep(200 * time.Millisecond) // Brief delay for consistency
		var count int
		err = proxyDB.QueryRow(
			fmt.Sprintf("SELECT COUNT(*) FROM %s.%s WHERE gate_name = ?", testSchema, testTable),
			testGateName,
		).Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 1, count, "Proxy should read the test row from local database")

		// Verify this row does NOT exist in remote database
		err = remoteDB.QueryRow(
			fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE gate_name = ?", testTable),
			testGateName,
		).Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 0, count, "Remote database should not have the test row")
	})

	// Test 3: Verify non-spoofed tables read from remote
	t.Run("NonSpoofedTableReadsFromRemote", func(t *testing.T) {
		// Read users table from proxy (should come from remote)
		var userCount int
		err := proxyDB.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s.users", testSchema)).Scan(&userCount)
		require.NoError(t, err)

		// Read users table directly from remote
		var remoteUserCount int
		err = remoteDB.QueryRow("SELECT COUNT(*) FROM users").Scan(&remoteUserCount)
		require.NoError(t, err)

		assert.Equal(t, remoteUserCount, userCount, "User count from proxy should match remote database")
		assert.Greater(t, userCount, 0, "Should have users in the database")
	})

	// Test 4: Verify local modifications are isolated
	t.Run("LocalModificationsAreIsolated", func(t *testing.T) {
		// Get initial count from remote
		var remoteCount int
		err := remoteDB.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", testTable)).Scan(&remoteCount)
		require.NoError(t, err)

		// Insert through proxy (goes to local)
		testGateName := fmt.Sprintf("ISOLATED_GATE_%d", time.Now().Unix())
		_, err = proxyDB.Exec(
			fmt.Sprintf("INSERT INTO %s.%s (gate_name, active_year, start_date, end_date) VALUES (?, ?, ?, ?)", testSchema, testTable),
			testGateName, 2026, "2026-02-01", "2026-02-28",
		)
		require.NoError(t, err, "Should be able to insert through proxy")

		// Verify remote count hasn't changed
		var newRemoteCount int
		err = remoteDB.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", testTable)).Scan(&newRemoteCount)
		require.NoError(t, err)
		assert.Equal(t, remoteCount, newRemoteCount, "Remote database should remain unchanged")

		// Verify proxy sees the new row
		var proxyHasRow int
		err = proxyDB.QueryRow(
			fmt.Sprintf("SELECT COUNT(*) FROM %s.%s WHERE gate_name = ?", testSchema, testTable),
			testGateName,
		).Scan(&proxyHasRow)
		require.NoError(t, err)
		assert.Equal(t, 1, proxyHasRow, "Proxy should see the new row from local database")
	})

	t.Log("All integration tests passed!")
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
func verifyRemoteData(t *testing.T, db *sql.DB) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM application_gates").Scan(&count)
	require.NoError(t, err, "Should be able to query application_gates")
	assert.Greater(t, count, 0, "Remote database should have application gates data")

	err = db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	require.NoError(t, err, "Should be able to query users")
	assert.Greater(t, count, 0, "Remote database should have users data")
}

// setupLocalReplication replicates the spoofed table to the local database
func setupLocalReplication(t *testing.T, localDB, remoteDB *sql.DB) {
	// Create the schema in local database
	_, err := localDB.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", testSchema))
	require.NoError(t, err)

	_, err = localDB.Exec(fmt.Sprintf("CREATE DATABASE %s", testSchema))
	require.NoError(t, err)

	_, err = localDB.Exec(fmt.Sprintf("USE %s", testSchema))
	require.NoError(t, err)

	// Get the CREATE TABLE statement from remote
	var tableName, createStmt string
	err = remoteDB.QueryRow(fmt.Sprintf("SHOW CREATE TABLE %s", testTable)).Scan(&tableName, &createStmt)
	require.NoError(t, err, "Should be able to get CREATE TABLE statement")

	// Create the table in local database
	_, err = localDB.Exec(createStmt)
	require.NoError(t, err, "Should be able to create table in local database")

	// Copy data from remote to local
	rows, err := remoteDB.Query(fmt.Sprintf("SELECT * FROM %s", testTable))
	require.NoError(t, err)
	defer rows.Close()

	cols, err := rows.Columns()
	require.NoError(t, err)

	for rows.Next() {
		values := make([]interface{}, len(cols))
		valuePtrs := make([]interface{}, len(cols))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		err = rows.Scan(valuePtrs...)
		require.NoError(t, err)

		// Build insert statement
		placeholders := ""
		for i := range cols {
			if i > 0 {
				placeholders += ", "
			}
			placeholders += "?"
		}

		insertStmt := fmt.Sprintf("INSERT INTO %s VALUES (%s)", testTable, placeholders)
		_, err = localDB.Exec(insertStmt, values...)
		require.NoError(t, err)
	}

	t.Logf("Successfully replicated %s table to local database", testTable)
}

// startTestProxy starts the proxy server in a goroutine
func startTestProxy(ctx context.Context, ready chan<- bool, errors chan<- error) {
	socket, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%s", testProxyPort))
	if err != nil {
		errors <- fmt.Errorf("failed to start proxy: %w", err)
		return
	}
	defer socket.Close()

	// Signal that proxy is ready
	ready <- true

	// Handle connections
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				originSocket, err := socket.Accept()
				if err != nil {
					// Check if we're shutting down
					select {
					case <-ctx.Done():
						return
					default:
						log.Printf("failed to accept connection: %s", err.Error())
						continue
					}
				}
				go handleConn(originSocket, testSchema, testTable)
			}
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()
}

// connectToProxy creates a connection to the proxy server
func connectToProxy(t *testing.T) *sql.DB {
	// Use the same credentials as local database
	url := fmt.Sprintf("root:root@tcp(127.0.0.1:%s)/", testProxyPort)

	var db *sql.DB
	var err error

	// Retry a few times as the proxy might need a moment to stabilize
	for i := 0; i < 5; i++ {
		db, err = sql.Open("mysql", url)
		if err == nil {
			err = db.Ping()
			if err == nil {
				return db
			}
		}
		time.Sleep(500 * time.Millisecond)
	}

	t.Fatalf("Failed to connect to proxy: %v", err)
	return nil
}

// getTables retrieves the list of tables from a database schema
func getTables(t *testing.T, db *sql.DB, schema string) []string {
	rows, err := db.Query(fmt.Sprintf("SHOW TABLES FROM %s", schema))
	require.NoError(t, err)
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var table string
		err := rows.Scan(&table)
		require.NoError(t, err)
		tables = append(tables, table)
	}

	return tables
}
