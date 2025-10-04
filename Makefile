.PHONY: help setup dev dev-backend dev-frontend migrate clean build docker-build docker-up docker-down

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

setup: ## Initial setup - install dependencies and configure environment
	@bash scripts/setup.sh

dev: ## Start all services in development mode
	@bash scripts/dev.sh

dev-backend: ## Start only the backend server
	@cd backend && go run cmd/server/main.go

dev-frontend: ## Start only the frontend server
	@cd frontend && npm run dev

migrate: ## Run database migrations
	@cd backend && go run cmd/migrate/main.go

clean: ## Clean up generated files and stop services
	@docker-compose down
	@cd backend && rm -f server
	@cd frontend && rm -rf .next node_modules

build: ## Build backend and frontend for production
	@echo "Building backend..."
	@cd backend && go build -o server cmd/server/main.go
	@echo "Building frontend..."
	@cd frontend && npm run build

docker-build: ## Build Docker images
	@docker-compose build

docker-up: ## Start services with Docker Compose
	@docker-compose up -d

docker-down: ## Stop Docker Compose services
	@docker-compose down

test-backend: ## Run backend tests
	@cd backend && go test ./...

test-frontend: ## Run frontend tests
	@cd frontend && npm test

lint-backend: ## Run backend linter
	@cd backend && go fmt ./...

lint-frontend: ## Run frontend linter
	@cd frontend && npm run lint

test-data: ## Generate test transaction data
	@cd backend && go run cmd/test-data/main.go
