.PHONY: help test integration-test unit-test clean docker-up docker-down docker-logs

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

test: unit-test ## Run all tests (unit + integration)
	@echo "Running unit tests first..."
	@$(MAKE) unit-test
	@echo ""
	@echo "Running integration tests..."
	@$(MAKE) integration-test

unit-test: ## Run unit tests only
	go test -v -run 'Test[^I]' -timeout 30s

integration-test: ## Run integration tests with Docker
	./run-integration-tests.sh

integration-test-keep: docker-up ## Run integration tests but keep containers running
	@echo "⏳ Waiting for databases to be healthy..."
	@sleep 10
	INTEGRATION_TEST=1 go test -v -run TestIntegration -timeout 2m

docker-up: ## Start Docker containers
	docker compose up -d

docker-down: ## Stop and remove Docker containers
	docker compose down -v

docker-logs: ## Show Docker container logs
	docker compose logs -f

docker-status: ## Show status of Docker containers
	docker compose ps

clean: docker-down ## Clean up everything (containers, volumes, binaries)
	rm -f imposter-db
	go clean -testcache

build: ## Build the binary
	go build -o imposter-db

run: build docker-up ## Build and run the application
	@echo "Waiting for databases to start..."
	@sleep 10
	./imposter-db -schema TEST_DB -table application_gates

deps: ## Download dependencies
	go mod download
	go mod tidy

.DEFAULT_GOAL := help
