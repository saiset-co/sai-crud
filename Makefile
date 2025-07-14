# SAI Service Microservice Makefile

# Build configuration
BINARY_NAME=sai-service
VERSION?=1.0.0

# Go configuration
GO_VERSION=1.21
GOOS?=linux
GOARCH?=amd64
CGO_ENABLED?=0

# Docker configuration
DOCKER_IMAGE=sai-service
DOCKER_TAG?=latest

# Environment configuration
ENV_FILE?=.env

# Default target
.DEFAULT_GOAL := help

# Colors for output
GREEN=\033[0;32m
YELLOW=\033[1;33m
RED=\033[0;31m
NC=\033[0m # No Color

## Help
.PHONY: help
help: ## Show this help message
	@echo "$(GREEN)SAI service Microservice$(NC)"
	@echo "$(YELLOW)Available commands:$(NC)"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(GREEN)%-20s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST)

## Environment Setup
.PHONY: setup-env
setup-env: ## Create .env file from .env.example
	@echo "$(YELLOW)Setting up environment file...$(NC)"
	@if [ -f ".env" ]; then \
		echo "$(GREEN).env file already exists!$(NC)"; \
	elif [ -f ".env.example" ]; then \
		cp .env.example .env; \
		echo "$(GREEN).env file created from .env.example$(NC)"; \
		echo "$(YELLOW)Please edit .env file with your configuration$(NC)"; \
	else \
		echo "$(RED)Error: .env.example file not found!$(NC)"; \
		echo "$(YELLOW)Please create .env.example file first or create .env manually$(NC)"; \
		exit 1; \
	fi

.PHONY: check-env
check-env: ## Check if .env file exists and create if needed
	@if [ ! -f "$(ENV_FILE)" ]; then \
		echo "$(YELLOW).env file not found. Creating...$(NC)"; \
		$(MAKE) setup-env; \
	fi

## Development
.PHONY: deps
deps: ## Download Go dependencies
	@echo "$(YELLOW)Downloading Go dependencies...$(NC)"
	@go mod download
	@go mod tidy
	@echo "$(GREEN)Dependencies downloaded!$(NC)"

.PHONY: config
config: check-env ## Generate config.yaml from template using environment variables
	@echo "$(YELLOW)Generating configuration from template...$(NC)"
	@if [ ! -f "config.yaml.template" ]; then \
		echo "$(RED)Error: config.yaml.template not found!$(NC)"; \
		exit 1; \
	fi
	@echo "$(YELLOW)Loading environment variables from $(ENV_FILE)...$(NC)"
	@set -a; . ./$(ENV_FILE); set +a; envsubst < ./config.yaml.template > ./config.yaml
	@echo "$(GREEN)Configuration generated at ./config.yaml$(NC)"

## Build
.PHONY: build
build: config ## Build the application binary
	@echo "$(YELLOW)Building $(BINARY_NAME)...$(NC)"
	@CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) go build \
		-ldflags="-w -s -X main.version=$(VERSION) -extldflags '-static'" \
		-a -installsuffix cgo \
		-o $(BINARY_NAME) \
		./cmd/main.go
	@echo "$(GREEN)Build complete: $(BINARY_NAME)$(NC)"

## Run
.PHONY: run
run: config ## Run the application locally
	@echo "$(YELLOW)Starting SAI Storage locally...$(NC)"
	@if [ ! -f "./config.yaml" ]; then \
		echo "$(RED)Configuration not found. Generating...$(NC)"; \
		$(MAKE) config; \
	fi
	@go run ./cmd/main.go

## Docker
.PHONY: docker-build
docker-build: ## Build Docker image
	@echo "$(YELLOW)Building Docker image $(DOCKER_IMAGE):$(DOCKER_TAG)...$(NC)"
	@docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .
	@echo "$(GREEN)Docker image built: $(DOCKER_IMAGE):$(DOCKER_TAG)$(NC)"

.PHONY: docker-run
docker-run: check-env ## Run Docker container
	@echo "$(YELLOW)Running Docker container...$(NC)"
	@docker run --rm --env-file .env -p 8080:8080 $(DOCKER_IMAGE):$(DOCKER_TAG)

## Docker Compose
.PHONY: up
up: check-env ## Start all services with docker-compose
	@echo "$(YELLOW)Starting all services...$(NC)"
	@docker-compose up -d
	@echo "$(GREEN)Services started!$(NC)"

.PHONY: down
down: ## Stop all services
	@echo "$(YELLOW)Stopping all services...$(NC)"
	@docker-compose down
	@echo "$(GREEN)Services stopped!$(NC)"

.PHONY: logs
logs: ## Show logs from all services
	@docker-compose logs -f

.PHONY: logs-app
logs-app: ## Show logs from application only
	@docker-compose logs -f sai-service

.PHONY: restart
restart: down up ## Restart all services

.PHONY: rebuild
rebuild: check-env ## Rebuild and restart all services
	@echo "$(YELLOW)Rebuilding and restarting services...$(NC)"
	@docker-compose down
	@docker-compose build --no-cache
	@docker-compose up -d
	@echo "$(GREEN)Services rebuilt and restarted!$(NC)"

## Code Quality
.PHONY: lint
lint: ## Run linter
	@echo "$(YELLOW)Running linter...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "$(RED)golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest$(NC)"; \
	fi

.PHONY: fmt
fmt: ## Format Go code
	@echo "$(YELLOW)Formatting Go code...$(NC)"
	@go fmt ./...
	@echo "$(GREEN)Code formatted!$(NC)"

.PHONY: vet
vet: ## Run go vet
	@echo "$(YELLOW)Running go vet...$(NC)"
	@go vet ./...

## Cleanup
.PHONY: clean
clean: ## Clean build artifacts and generated files
	@echo "$(YELLOW)Cleaning build artifacts...$(NC)"
	@rm -f $(BINARY_NAME)
	@rm -f coverage.out coverage.html
	@rm -f config.yaml
	@echo "$(GREEN)Cleanup complete!$(NC)"

.PHONY: clean-docker
clean-docker: ## Clean Docker images and volumes
	@echo "$(YELLOW)Cleaning Docker resources...$(NC)"
	@docker-compose down -v --remove-orphans
	@docker system prune -f
	@echo "$(GREEN)Docker cleanup complete!$(NC)"

.PHONY: clean-all
clean-all: clean clean-docker ## Clean everything