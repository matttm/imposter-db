#!/bin/bash

# Quick test runner for integration tests with proper environment variables

export INTEGRATION_TEST=1
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

echo "🧪 Running integration test with environment variables..."
echo "   Remote: $REMOTE_DB_USER@$REMOTE_DB_HOST:$REMOTE_DB_PORT/$REMOTE_DB_NAME"
echo "   Local: $LOCAL_DB_USER@$LOCAL_DB_HOST:$LOCAL_DB_PORT"
echo ""

go test -v -run TestIntegration -timeout 2m
