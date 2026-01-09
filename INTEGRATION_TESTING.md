# Integration Testing Guide

## Overview

The integration test suite validates the full proxy functionality using the Docker containers defined in `docker-compose.yml`. The tests verify that:

1. The proxy correctly routes queries for spoofed tables to the local database
2. Queries for non-spoofed tables are forwarded to the remote database
3. Local modifications remain isolated from the remote database
4. The proxy handles multiple simultaneous client connections

## Test Architecture

The integration test uses:
- **Local Database** (`imposter-local`): MySQL on port 3306 - stores the spoofed table
- **Remote Database** (`imposter-remote`): MySQL on port 3307 - contains the full test dataset
- **Proxy Server**: Listens on port 13306 and routes queries between local and remote

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

### TestIntegration_ProxyWithDockerContainers

The main integration test includes four sub-tests:

#### 1. ProxyShowsTables
Verifies that the proxy correctly exposes all tables from the remote database, including the spoofed table.

#### 2. SpoofedTableReadsFromLocal
Tests that queries against the spoofed table (`application_gates`) are served from the local database:
- Inserts a unique row into the local database
- Queries through the proxy to verify the row is visible
- Confirms the row does NOT exist in the remote database

#### 3. NonSpoofedTableReadsFromRemote
Validates that queries for non-spoofed tables (e.g., `users`) are correctly routed to the remote database.

#### 4. LocalModificationsAreIsolated
Confirms that inserts/updates through the proxy only affect the local database:
- Inserts data through the proxy
- Verifies the remote database remains unchanged
- Confirms the proxy sees the new data from local

## Test Data

The remote database is initialized with test data from `init.sql`, which creates:
- `user_types`: Reference data for user types
- `users`: User accounts
- `application_gates`: Timeline/enrollment gates (this is the spoofed table)
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

If ports 3306, 3307, or 13306 are already in use:
1. Stop conflicting services
2. Or modify the ports in `docker-compose.yml` and update the test constants

### Test Timeouts

If tests timeout, it may indicate:
- Databases are slow to initialize (increase `healthCheckRetries`)
- Network connectivity issues
- Proxy not starting correctly

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

1. Add a new `t.Run()` block in `TestIntegration_ProxyWithDockerContainers`
2. Use the existing database connections (`localDB`, `remoteDB`, `proxyDB`)
3. Follow the pattern of verifying isolation between local and remote

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
- **Test execution**: 5-10 seconds
- **Cleanup**: 5 seconds

Total time: ~30-45 seconds per run.
