#!/bin/bash

# Integration test runner for imposter-db
# This script starts docker-compose, runs integration tests, and cleans up

set -e

# Compose file can be overridden (e.g., docker-compose.mysql8.yml)
# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

COMPOSE_FILES=()
if [ -n "${COMPOSE_FILE:-}" ]; then
    COMPOSE_FILES=("$COMPOSE_FILE")
else
    COMPOSE_FILES=(docker-compose.mysql8.yml docker-compose.mysql91.yml)
fi

echo "🚀 Starting integration test suite for imposter-db..."
echo "${YELLOW}📄 Compose files: ${COMPOSE_FILES[*]}${NC}"

overall_status=0

run_for_compose() {
    local compose_file="$1"
    local -a compose_cmd=(docker compose -f "$compose_file")

    # Function to cleanup on exit
    cleanup() {
        echo ""
        echo "${YELLOW}🧹 Cleaning up...${NC}"
        "${compose_cmd[@]}" down -v
        echo "${GREEN}✅ Cleanup complete${NC}"
    }

    # Register cleanup function for this run
    trap cleanup EXIT

    echo ""
    echo "${YELLOW}▶️  Running integration tests with: ${compose_file}${NC}"

    # Start docker-compose services
    echo "${YELLOW}📦 Starting Docker containers...${NC}"
    "${compose_cmd[@]}" up -d

    # Wait for containers to be healthy
    echo "${YELLOW}⏳ Waiting for databases to be healthy...${NC}"
    max_wait=60
    elapsed=0

    while [ $elapsed -lt $max_wait ]; do
        local_id=$("${compose_cmd[@]}" ps -q localdb 2>/dev/null || true)
        remote_id=$("${compose_cmd[@]}" ps -q db 2>/dev/null || true)
        local_healthy=$(docker inspect --format='{{.State.Health.Status}}' "$local_id" 2>/dev/null || echo "starting")
        remote_healthy=$(docker inspect --format='{{.State.Health.Status}}' "$remote_id" 2>/dev/null || echo "starting")
        
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
        "${compose_cmd[@]}" ps
        "${compose_cmd[@]}" logs
        exit 1
    fi

    # Give databases an extra moment to fully initialize
    sleep 2

    echo ""
    echo "${YELLOW}🧪 Running integration tests...${NC}"
    echo ""

    # Set environment variables for the test
    export REMOTE_DB_PORT=3307
    export REMOTE_DB_HOST=127.0.0.1
    export REMOTE_DB_USER=ADMIN
    export REMOTE_DB_PASS=ADMIN
    export REMOTE_DB_NAME=TEST_DB
    export PROXY_PORT=13306
    export LOCAL_DB_PORT=3306
    export LOCAL_DB_HOST=127.0.0.1
    export LOCAL_DB_USER=root
    export LOCAL_DB_PASS=root
    export LOCAL_DB_NAME=""

    # Run the integration tests
    set +e
    INTEGRATION_TEST=1 go test -v -run TestIntegration -timeout 2m
    test_result=$?
    set -e

    echo ""
    if [ $test_result -eq 0 ]; then
        echo "${GREEN}✅ Integration tests passed for ${compose_file}!${NC}"
    else
        echo "${RED}❌ Integration tests failed for ${compose_file}${NC}"
        echo "${YELLOW}📋 Docker logs:${NC}"
        "${compose_cmd[@]}" logs --tail=50
    fi

    cleanup
    trap - EXIT

    return $test_result
}

for compose_file in "${COMPOSE_FILES[@]}"; do
    if ! run_for_compose "$compose_file"; then
        overall_status=1
    fi
done

if [ $overall_status -eq 0 ]; then
    echo ""
    echo "${GREEN}✅ All integration tests passed for all compose files!${NC}"
else
    echo ""
    echo "${RED}❌ One or more integration test runs failed${NC}"
fi

exit $overall_status
