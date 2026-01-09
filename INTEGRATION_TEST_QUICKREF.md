# Integration Test Quick Reference

## 🚀 Quick Start

Run the complete integration test suite:
```bash
make integration-test
```

Or directly:
```bash
./run-integration-tests.sh
```

## 📋 What Gets Tested

✅ **Proxy exposes all tables** from remote database  
✅ **Spoofed table queries** are served from local database  
✅ **Non-spoofed tables** are forwarded to remote database  
✅ **Local modifications** remain isolated from remote  

## 🔧 Test Infrastructure

| Component | Port | Purpose |
|-----------|------|---------|
| Local DB (MySQL) | 3306 | Stores spoofed `application_gates` table |
| Remote DB (MySQL) | 3307 | Full test dataset with all tables |
| Proxy Server | 13306 | Routes queries between local & remote |

## 📝 Common Commands

```bash
# Run all tests (unit + integration)
make test

# Run only integration tests
make integration-test

# Run integration tests and keep containers up for debugging
make integration-test-keep

# Start Docker containers manually
make docker-up

# Stop and clean up
make docker-down

# View container logs
make docker-logs

# Check container status
make docker-status
```

## 🐛 Debugging

If tests fail, check:

1. **Container health**:
   ```bash
   docker ps
   docker-compose logs
   ```

2. **Port availability**:
   ```bash
   lsof -i :3306
   lsof -i :3307
   lsof -i :13306
   ```

3. **Run with verbose output**:
   ```bash
   INTEGRATION_TEST=1 go test -v -run TestIntegration
   ```

## 🔍 Test Structure

```
integration_test.go
├── TestIntegration_ProxyWithDockerContainers
    ├── Setup (databases + proxy)
    ├── ProxyShowsTables
    ├── SpoofedTableReadsFromLocal
    ├── NonSpoofedTableReadsFromRemote
    └── LocalModificationsAreIsolated
```

## ⚙️ Environment Variables

Tests automatically set these:
- `INTEGRATION_TEST=1` (enables integration tests)
- Database connection configs for local, remote, and proxy

## 📚 More Details

See [INTEGRATION_TESTING.md](INTEGRATION_TESTING.md) for:
- Detailed test descriptions
- CI/CD integration examples
- Troubleshooting guide
- How to extend tests
