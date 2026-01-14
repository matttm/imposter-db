# Integration Testing Guide

## Overview

The integration test suite validates that the `main()` function correctly replicates a chosen table from the remote database to the local database using the Docker containers defined in `docker-compose.yml`. The tests verify that:

1. The chosen table is successfully created in the local database
2. Table schema (columns and types) matches the remote database
3. All data is correctly copied from remote to local
4. Row counts match between databases
5. Local and remote databases remain independent after replication

## Test Architecture

The integration test uses:
- **Local Database** (`imposter-local`): MySQL on port 3306 - target for table replication
- **Remote Database** (`imposter-remote`): MySQL on port 3307 - source containing the test dataset
- **Main Function**: Called with flags `-schema TEST_DB -table application_gates -fk=false`

## Running the Tests

### Quick Start

The easiest way to run the integration tests is using the provided shell script:

```bash
./run-integration-tests.sh
```

This script will:
1. Start both Docker containers
2. Wait for databases to become healthy
3. Run the integration test suite
4. Clean up containers and volumes

### Manual Execution

If you prefer to run tests manually:

1. Start the Docker containers:
```bash
docker-compose up -d
```

2. Wait for containers to be healthy (check with `docker ps`):
```bash
docker ps
```

3. Run the integration tests:
```bash
INTEGRATION_TEST=1 go test -v -run TestIntegration -timeout 2m
```

4. Clean up when done:
```bash
docker-compose down -v
```

## Test Cases

### TestIntegration_TableReplication

The main integration test calls `main()` with appropriate flags and includes five sub-tests:

#### 1. TableExistsInLocal
Verifies that the specified table (`application_gates`) was successfully created in the local database after running `main()`.

#### 2. RowCountMatches
Confirms that the local table has the exact same number of rows as the remote table, ensuring complete data replication.

#### 3. TableSchemaMatches
Validates that the table structure is identical:
- Same number of columns
- Column names match
- Column data types match

#### 4. DataCopiedCorrectly
Verifies that actual data was copied correctly by:
- Capturing all rows from both databases
- Comparing the first 5 rows to ensure they match
- Checking that row counts are identical

#### 5. DatabasesAreIndependent
Confirms that local and remote databases remain independent:
- Inserts a test row into the local database
- Verifies the row exists in local but NOT in remote
- Ensures modifications to local don't affect remote

## Test Data

The remote database is initialized with test data from `init.sql`, which creates:
- `user_types`: Reference data for user types
- `users`: User accounts
- `application_gates`: Timeline/enrollment gates (this is the replicated table)
- `applications`: Links users to gates

Each table contains 20 rows of sample data.

## Environment Variables

The tests configure the following environment variables:

```bash
# Remote Database (port 3307)
REMOTE_DB_HOST=127.0.0.1
REMOTE_DB_PORT=3307
REMOTE_DB_USER=ADMIN
REMOTE_DB_PASS=ADMIN
REMOTE_DB_NAME=TEST_DB

# Local Database (port 3306)
LOCAL_DB_HOST=127.0.0.1
LOCAL_DB_PORT=3306
LOCAL_DB_USER=root
LOCAL_DB_PASS=root
LOCAL_DB_NAME=

# Proxy
PROXY_PORT=13306
```

## Troubleshooting

### Databases Not Starting

If containers fail to start or become healthy:
```bash
docker-compose logs
docker-compose ps
```

### Port Conflicts

If ports 3306 or 3307 are already in use:
1. Stop conflicting services
2. Or modify the ports in `docker-compose.yml` and update the test constants

### Test Timeouts

If tests timeout, it may indicate:
- Databases are slow to initialize (increase `healthCheckRetries`)
- Network connectivity issues
- `main()` taking longer than 10 seconds to replicate data

Check logs with:
```bash
docker-compose logs localdb
docker-compose logs db
```

### Cleanup Issues

If cleanup fails, manually remove containers and volumes:
```bash
docker-compose down -v
docker volume prune -f
```

## CI/CD Integration

To run these tests in a CI pipeline:

```yaml
# Example GitHub Actions
- name: Run Integration Tests
  run: |
    chmod +x run-integration-tests.sh
    ./run-integration-tests.sh
```

Or with explicit steps:
```yaml
- name: Start Docker Compose
  run: docker-compose up -d

- name: Wait for Health
  run: |
    timeout 60 bash -c 'until docker inspect --format="{{.State.Health.Status}}" imposter-local | grep -q healthy; do sleep 2; done'
    timeout 60 bash -c 'until docker inspect --format="{{.State.Health.Status}}" imposter-remote | grep -q healthy; do sleep 2; done'

- name: Run Tests
  run: INTEGRATION_TEST=1 go test -v -run TestIntegration -timeout 2m

- name: Cleanup
  if: always()
  run: docker-compose down -v
```

## Extending the Tests

To add new test cases:

1. Add a new `t.Run()` block in `TestIntegration_TableReplication`
2. Use the existing database connections (`localDB`, `remoteDB`)
3. Follow the pattern of verifying table replication and database independence

Example:
```go
t.Run("MyNewTestCase", func(t *testing.T) {
    // Your test logic here
    // Use require.NoError() for fatal errors
    // Use assert.Equal() for non-fatal assertions
})
```

## Performance Considerations

The integration test suite typically takes:
- **Container startup**: 15-30 seconds
- **Table replication** (via main()): 10-15 seconds
- **Test verification**: 2-5 seconds
- **Cleanup**: 5 seconds

Total time: ~35-55 seconds per run.
