#!/bin/bash

# Integration test runner for imposter-db
# This script starts docker-compose, runs integration tests, and cleans up

set -e

echo "🚀 Starting integration test suite for imposter-db..."

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Function to cleanup on exit
cleanup() {
    echo ""
    echo "${YELLOW}🧹 Cleaning up...${NC}"
    docker-compose down -v
    echo "${GREEN}✅ Cleanup complete${NC}"
}

# Register cleanup function
trap cleanup EXIT

# Start docker-compose services
echo "${YELLOW}📦 Starting Docker containers...${NC}"
docker compose up -d

# Wait for containers to be healthy
echo "${YELLOW}⏳ Waiting for databases to be healthy...${NC}"
max_wait=60
elapsed=0

while [ $elapsed -lt $max_wait ]; do
    local_healthy=$(docker inspect --format='{{.State.Health.Status}}' imposter-local 2>/dev/null || echo "starting")
    remote_healthy=$(docker inspect --format='{{.State.Health.Status}}' imposter-remote 2>/dev/null || echo "starting")
    
    if [ "$local_healthy" = "healthy" ] && [ "$remote_healthy" = "healthy" ]; then
        echo "${GREEN}✅ Both databases are healthy!${NC}"
        break
    fi
    
    echo "   Local DB: $local_healthy, Remote DB: $remote_healthy (${elapsed}s elapsed)"
    sleep 2
    elapsed=$((elapsed + 2))
done

if [ $elapsed -ge $max_wait ]; then
    echo "${RED}❌ Timeout waiting for databases to become healthy${NC}"
    docker compose ps
    docker compose logs
    exit 1
fi

# Give databases an extra moment to fully initialize
sleep 2

echo ""
echo "${YELLOW}🧪 Running integration tests...${NC}"
echo ""

# Run the integration tests
INTEGRATION_TEST=1 go test -v -run TestIntegration -timeout 2m

test_result=$?

echo ""
if [ $test_result -eq 0 ]; then
    echo "${GREEN}✅ All integration tests passed!${NC}"
else
    echo "${RED}❌ Integration tests failed${NC}"
    echo "${YELLOW}📋 Docker logs:${NC}"
    docker compose logs --tail=50
fi

exit $test_result
