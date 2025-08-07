# Status Page Starter - Comprehensive Makefile
# =============================================

# Variables
PROJECT_NAME := status-page-starter
DOCKER_COMPOSE := docker-compose
DOCKER_COMPOSE_DEV := docker-compose -f docker-compose.yml
DOCKER_COMPOSE_PROD := docker-compose -f docker-compose.yml -f docker-compose.prod.yml
DOCKER_COMPOSE_TEST := docker-compose -f docker-compose.test.yml

# Build variables for version injection
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE := $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
LDFLAGS := -X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.BuildDate=$(BUILD_DATE) -s -w

# Colors for output
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[1;33m
BLUE := \033[0;34m
NC := \033[0m # No Color

# Default target
.DEFAULT_GOAL := help

# Helper function to find web directory
define find_web_dir
$(shell if [ -f web/package.json ]; then echo "web"; elif [ -f web/react-status-page/package.json ]; then echo "web/react-status-page"; else echo ""; fi)
endef

# Phony targets
.PHONY: help setup dev prod test clean build deploy docs lint format security backup restore monitor logs

##@ Help
help: ## Display this help message
	@echo "$(BLUE)Status Page Starter - Available Commands$(NC)"
	@echo "========================================"
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Setup & Installation
setup: ## Initial project setup
	@echo "$(YELLOW)Setting up project...$(NC)"
	@make setup-env
	@make setup-go
	@make setup-node
	@make setup-docker
	@make setup-git-hooks
	@echo "$(GREEN)✓ Project setup completed!$(NC)"

setup-env: ## Setup environment configuration
	@echo "$(BLUE)Setting up environment configuration...$(NC)"
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "$(YELLOW)Created .env file from .env.example$(NC)"; \
		echo "$(YELLOW)Please edit .env with your configuration before proceeding$(NC)"; \
	else \
		echo "$(GREEN)✓ .env file already exists$(NC)"; \
	fi

setup-go: ## Setup Go environment
	@echo "$(BLUE)Setting up Go environment...$(NC)"
	@go mod download
	@go mod tidy
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install golang.org/x/tools/cmd/goimports@latest

setup-node: ## Setup Node.js environment
	@echo "$(BLUE)Setting up Node.js environment...$(NC)"
	@if [ -f web/react-status-page/package.json ]; then \
		cd web/react-status-page && npm install && npm audit fix || true; \
	else \
		echo "$(YELLOW)No package.json found - skipping Node.js setup$(NC)"; \
	fi

setup-docker: ## Setup Docker environment
	@echo "$(BLUE)Setting up Docker environment...$(NC)"
	@docker --version || (echo "$(RED)Docker not installed$(NC)" && exit 1)
	@docker-compose --version || (echo "$(RED)Docker Compose not installed$(NC)" && exit 1)
	@docker network create status-network || true

setup-git-hooks: ## Setup Git hooks
	@echo "$(BLUE)Setting up Git hooks...$(NC)"
	@cp scripts/hooks/pre-commit .git/hooks/pre-commit || true
	@chmod +x .git/hooks/pre-commit || true

##@ Development
dev: ## Start development environment
	@echo "$(YELLOW)Starting development environment...$(NC)"
	@$(DOCKER_COMPOSE_DEV) up -d
	@make wait-for-services
	@make seed-db
	@echo "$(GREEN)✓ Development environment ready!$(NC)"
	@echo "$(BLUE)Access points:$(NC)"
	@echo "  - Status Page: http://localhost"
	@echo "  - API Health: http://localhost/api/health"
	@echo "  - MongoDB: mongodb://localhost:27017"

dev-logs: ## Follow development logs
	@$(DOCKER_COMPOSE_DEV) logs -f

dev-rebuild: ## Rebuild and restart development environment
	@echo "$(YELLOW)Rebuilding development environment...$(NC)"
	@$(DOCKER_COMPOSE_DEV) down
	@$(DOCKER_COMPOSE_DEV) build --no-cache
	@$(DOCKER_COMPOSE_DEV) up -d
	@make wait-for-services

dev-frontend: ## Start only frontend development
	@echo "$(YELLOW)Starting frontend development server...$(NC)"
	@if [ -f web/react-status-page/package.json ]; then \
		cd web/react-status-page && npm run dev; \
	else \
		echo "$(RED)No package.json found - cannot start frontend$(NC)"; \
	fi

dev-api: ## Run API server locally
	@echo "$(YELLOW)Starting API server locally...$(NC)"
	@export MONGO_URI=mongodb://localhost:27017/status_page && go run ./cmd/api

dev-checker: ## Run status checker locally
	@echo "$(YELLOW)Starting status checker locally...$(NC)"
	@export MONGO_URI=mongodb://localhost:27017/status_page && go run ./cmd/status-checker

##@ Testing & Quality
test: ## Run all tests
	@echo "$(YELLOW)Running all tests...$(NC)"
	@make test-go
	@make test-frontend
	@echo "$(GREEN)✓ All tests completed!$(NC)"

test-go: ## Run Go tests
	@echo "$(BLUE)Running Go tests...$(NC)"
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out
	@go tool cover -html=coverage.out -o coverage.html

test-frontend: ## Run frontend tests
	@echo "$(BLUE)Running frontend tests...$(NC)"
	@if [ -f web/react-status-page/package.json ]; then \
		cd web/react-status-page && npm test; \
	else \
		echo "$(YELLOW)No package.json found - skipping frontend tests$(NC)"; \
	fi

test-integration: ## Run integration tests
	@echo "$(BLUE)Running integration tests...$(NC)"
	@make db-start
	@go test -v -tags=integration ./...
	@make db-stop

test-e2e: ## Run end-to-end tests
	@echo "$(BLUE)Running end-to-end tests...$(NC)"
	@make dev
	@go test -v -tags=e2e ./tests/e2e/...
	@make clean

generate-mocks: ## Generate mocks using Mockery
	@echo "$(BLUE)Generating mocks...$(NC)"
	@go install github.com/vektra/mockery/v2@latest
	@mockery
	@echo "$(GREEN)✓ Mocks generated!$(NC)"

test-with-mocks: ## Run tests with generated mocks
	@echo "$(BLUE)Running tests with mocks...$(NC)"
	@make generate-mocks
	@go test -v ./internal/application/handlers/...
	@go test -v ./internal/checker/...
	@echo "$(GREEN)✓ Tests with mocks completed!$(NC)"


##@ Code Quality
lint: ## Run all linters
	@echo "$(YELLOW)Running linters...$(NC)"
	@make lint-go
	@make lint-frontend
	@echo "$(GREEN)✓ Linting completed!$(NC)"

lint-go: ## Run Go linter
	@echo "$(BLUE)Running Go linter...$(NC)"
	@golangci-lint run ./...

lint-frontend: ## Run frontend linter
	@echo "$(BLUE)Running frontend linter...$(NC)"
	@if [ -f web/react-status-page/package.json ]; then \
		cd web/react-status-page && npm run lint; \
	else \
		echo "$(YELLOW)No package.json found - skipping frontend linting$(NC)"; \
	fi

format: ## Format all code
	@echo "$(YELLOW)Formatting code...$(NC)"
	@make format-go
	@make format-frontend
	@echo "$(GREEN)✓ Code formatting completed!$(NC)"

format-go: ## Format Go code
	@echo "$(BLUE)Formatting Go code...$(NC)"
	@go fmt ./...
	@goimports -w .

format-frontend: ## Format frontend code
	@echo "$(BLUE)Formatting frontend code...$(NC)"
	@if [ -f web/react-status-page/package.json ]; then \
		echo "$(YELLOW)No format script configured - skipping frontend formatting$(NC)"; \
	else \
		echo "$(YELLOW)No package.json found - skipping frontend formatting$(NC)"; \
	fi

##@ Security
security: ## Run security scans
	@echo "$(YELLOW)Running security scans...$(NC)"
	@make security-go
	@make security-frontend
	@make security-docker
	@echo "$(GREEN)✓ Security scans completed!$(NC)"

security-go: ## Run Go security scan
	@echo "$(BLUE)Running Go security scan...$(NC)"
	@go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	@gosec ./...

security-frontend: ## Run frontend security scan
	@echo "$(BLUE)Running frontend security scan...$(NC)"
	@cd web && npm audit
	@cd web && npm audit fix || true

security-docker: ## Run Docker security scan
	@echo "$(BLUE)Running Docker security scan...$(NC)"
	@docker run --rm -v /var/run/docker.sock:/var/run/docker.sock aquasec/trivy image $(PROJECT_NAME):latest || true

##@ Build & Deploy
build: ## Build all services
	@echo "$(YELLOW)Building all services...$(NC)"
	@make build-go
	@make build-frontend
	@make build-docker
	@echo "$(GREEN)✓ Build completed!$(NC)"

build-go: ## Build Go binaries with version injection
	@echo "$(BLUE)Building Go binaries with version $(VERSION)...$(NC)"
	@mkdir -p bin
	@echo "Building status-page binary..."
	@go build -ldflags "$(LDFLAGS)" -o bin/status-page .
	@echo "Building individual service binaries..."
	@go build -ldflags "$(LDFLAGS)" -o bin/status-page-api ./cmd/api
	@go build -ldflags "$(LDFLAGS)" -o bin/status-page-checker ./cmd/status-checker
	@echo "$(GREEN)✓ Go binaries built successfully$(NC)"

build-frontend: ## Build frontend
	@echo "$(BLUE)Building frontend...$(NC)"
	@if [ -f web/react-status-page/package.json ]; then \
		cd web/react-status-page && npm run build; \
	else \
		echo "$(YELLOW)No package.json found - skipping frontend build$(NC)"; \
	fi

build-docker: ## Build Docker images
	@echo "$(BLUE)Building Docker images...$(NC)"
	@$(DOCKER_COMPOSE_DEV) build

deploy-dev: ## Deploy to development
	@echo "$(YELLOW)Deploying to development...$(NC)"
	@make build
	@$(DOCKER_COMPOSE_DEV) up -d
	@make wait-for-services
	@echo "$(GREEN)✓ Development deployment completed!$(NC)"

deploy-prod: ## Deploy to production
	@echo "$(YELLOW)Deploying to production...$(NC)"
	@make build
	@$(DOCKER_COMPOSE_PROD) up -d
	@make wait-for-services
	@make health-check
	@echo "$(GREEN)✓ Production deployment completed!$(NC)"

deploy-staging: ## Deploy to staging
	@echo "$(YELLOW)Deploying to staging...$(NC)"
	@$(DOCKER_COMPOSE) -f docker-compose.staging.yml up -d
	@make wait-for-services
	@echo "$(GREEN)✓ Staging deployment completed!$(NC)"

##@ Database
db-start: ## Start database only
	@echo "$(YELLOW)Starting database...$(NC)"
	@$(DOCKER_COMPOSE_DEV) up -d mongo
	@echo "$(GREEN)✓ Database started!$(NC)"

db-stop: ## Stop database
	@$(DOCKER_COMPOSE_DEV) stop mongo

db-shell: ## Connect to database shell
	@$(DOCKER_COMPOSE_DEV) exec mongo mongosh

seed-db: ## Seed database with sample data
	@echo "$(BLUE)Seeding database...$(NC)"
	@./scripts/seed-db.sh
	@echo "$(GREEN)✓ Database seeded!$(NC)"

backup-db: ## Backup database
	@echo "$(YELLOW)Creating database backup...$(NC)"
	@mkdir -p backups
	@docker-compose exec mongo mongodump --out /tmp/backup
	@docker cp $$(docker-compose ps -q mongo):/tmp/backup ./backups/backup-$$(date +%Y%m%d-%H%M%S)
	@echo "$(GREEN)✓ Database backup completed!$(NC)"

restore-db: ## Restore database (usage: make restore-db BACKUP=backup-20240101-120000)
	@echo "$(YELLOW)Restoring database from $(BACKUP)...$(NC)"
	@docker cp ./backups/$(BACKUP) $$(docker-compose ps -q mongo):/tmp/restore
	@docker-compose exec mongo mongorestore /tmp/restore
	@echo "$(GREEN)✓ Database restored!$(NC)"

##@ Monitoring & Maintenance
status: ## Show service status
	@echo "$(BLUE)Service Status:$(NC)"
	@$(DOCKER_COMPOSE_DEV) ps

logs: ## Show all logs
	@$(DOCKER_COMPOSE_DEV) logs --tail=100 -f

logs-api: ## Show API logs
	@$(DOCKER_COMPOSE_DEV) logs --tail=100 -f api

logs-checker: ## Show status checker logs
	@$(DOCKER_COMPOSE_DEV) logs --tail=100 -f status-checker

logs-web: ## Show web server logs
	@$(DOCKER_COMPOSE_DEV) logs --tail=100 -f web

health-check: ## Run health checks
	@echo "$(BLUE)Running health checks...$(NC)"
	@curl -f http://localhost/api/health || (echo "$(RED)API health check failed$(NC)" && exit 1)
	@curl -f http://localhost || (echo "$(RED)Web health check failed$(NC)" && exit 1)
	@echo "$(GREEN)✓ All health checks passed!$(NC)"

monitor: ## Show system monitoring
	@echo "$(BLUE)System Monitoring:$(NC)"
	@docker stats --no-stream

##@ Cleanup
clean: ## Clean up development environment
	@echo "$(YELLOW)Cleaning up...$(NC)"
	@$(DOCKER_COMPOSE_DEV) down -v
	@docker system prune -f
	@rm -f coverage.out coverage.html
	@cd web && rm -rf node_modules/.cache dist
	@echo "$(GREEN)✓ Cleanup completed!$(NC)"

clean-all: ## Deep clean (removes all containers, images, volumes)
	@echo "$(RED)WARNING: This will remove all Docker containers, images, and volumes!$(NC)"
	@read -p "Are you sure? (y/N): " confirm && [ "$$confirm" = "y" ] || exit 1
	@$(DOCKER_COMPOSE_DEV) down -v --remove-orphans
	@docker system prune -a -f --volumes
	@echo "$(GREEN)✓ Deep cleanup completed!$(NC)"

##@ Documentation
docs: ## Generate documentation
	@echo "$(YELLOW)Generating documentation...$(NC)"
	@make docs-api
	@make docs-frontend
	@echo "$(GREEN)✓ Documentation generated!$(NC)"

docs-api: ## Generate API documentation
	@echo "$(BLUE)Generating API documentation...$(NC)"
	@swag init -g ./cmd/api/main.go -o ./docs/swagger

docs-frontend: ## Generate frontend documentation
	@echo "$(BLUE)Generating frontend documentation...$(NC)"
	@cd web && npm run docs || echo "Frontend docs not configured"

docs-serve: ## Serve documentation locally
	@echo "$(BLUE)Serving documentation at http://localhost:8080$(NC)"
	@cd docs && python -m http.server 8080

##@ Utilities
wait-for-services: ## Wait for services to be ready
	@echo "$(BLUE)Waiting for services to be ready...$(NC)"
	@./scripts/wait-for-services.sh

version: ## Show version information
	@echo "$(BLUE)Version Information:$(NC)"
	@echo "Go: $$(go version)"
	@echo "Node: $$(node --version)"
	@echo "Docker: $$(docker --version)"
	@echo "Docker Compose: $$(docker-compose --version)"

env-check: ## Check environment requirements
	@echo "$(BLUE)Checking environment...$(NC)"
	@./scripts/env-check.sh

##@ CI/CD
ci-setup: ## Setup CI environment
	@make setup
	@make lint
	@make security

ci-test: ## CI test pipeline
	@make test
	@make security

ci-build: ## CI build pipeline
	@make build
	@make test-integration

ci-deploy: ## CI deploy pipeline
	@make ci-build
	@make deploy-staging

# mocks: mocks-clean mocks-generate
# 
# mocks-generate:
# 	go generate ./... && mockery --all --inpackage --with-expecter=true
# 
# mocks-clean:
# 	find . -name "mock_*.go" -type f -delete