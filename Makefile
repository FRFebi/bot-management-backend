.PHONY: help build run test clean docker-build docker-run docker-dev deps

# Default target
help: ## Show this help message
	@echo 'Usage: make <target>'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

deps: ## Install dependencies
	go mod download
	go mod tidy

build: ## Build the application
	CGO_ENABLED=0 go build -o bin/bot-management-backend ./cmd/server

run: ## Run the application
	go run ./cmd/server

test: ## Run tests
	go test -v ./...

clean: ## Clean build artifacts
	rm -rf bin/
	go clean

docker-build: ## Build Docker image
	docker build -t bot-management-backend .

docker-run: ## Run Docker container
	docker run --rm -p 8080:8080 bot-management-backend

docker-dev: ## Run development environment with Docker Compose
	docker-compose -f docker-compose.dev.yml up --build

docker-dev-down: ## Stop development environment
	docker-compose -f docker-compose.dev.yml down

lint: ## Run linter
	golangci-lint run

fmt: ## Format code
	go fmt ./...

mod-tidy: ## Tidy go modules
	go mod tidy

.DEFAULT_GOAL := help