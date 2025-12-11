.PHONY: help run build test clean install-tools swagger migrate-up migrate-down migrate-create migrate-force migrate-version

# Load environment variables
include .env
export

# Application info
APP_NAME := go-api-starter
VERSION := 1.0.0
BUILD_DIR := bin

# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod
BINARY_NAME := server
MAIN_PATH := ./cmd/server

# Database connection string
DB_URL := postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

# Colors for output
COLOR_RESET := \033[0m
COLOR_BOLD := \033[1m
COLOR_GREEN := \033[32m
COLOR_YELLOW := \033[33m
COLOR_BLUE := \033[34m

help: ## Show this help message
	@echo '$(COLOR_BOLD)$(APP_NAME) - Available Commands:$(COLOR_RESET)'
	@echo ''
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(COLOR_BLUE)%-20s$(COLOR_RESET) %s\n", $$1, $$2}'
	@echo ''

install-tools: ## Install required development tools
	@echo "$(COLOR_GREEN)Installing development tools...$(COLOR_RESET)"
	go install github.com/swaggo/swag/cmd/swag@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@echo "$(COLOR_GREEN)✅ Tools installed successfully$(COLOR_RESET)"

deps: ## Download Go module dependencies
	@echo "$(COLOR_GREEN)Downloading dependencies...$(COLOR_RESET)"
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "$(COLOR_GREEN)✅ Dependencies downloaded$(COLOR_RESET)"

build: ## Build the application
	@echo "$(COLOR_GREEN)Building $(APP_NAME)...$(COLOR_RESET)"
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "$(COLOR_GREEN)✅ Build completed: $(BUILD_DIR)/$(BINARY_NAME)$(COLOR_RESET)"

run: swagger ## Run the application
	@echo "$(COLOR_GREEN)Starting $(APP_NAME)...$(COLOR_RESET)"
	$(GOCMD) run $(MAIN_PATH)/main.go

dev: ## Run with hot reload (requires air)
	@echo "$(COLOR_GREEN)Starting development server with hot reload...$(COLOR_RESET)"
	air

test: ## Run tests
	@echo "$(COLOR_GREEN)Running tests...$(COLOR_RESET)"
	$(GOTEST) -v -cover ./...

test-coverage: ## Run tests with coverage
	@echo "$(COLOR_GREEN)Running tests with coverage...$(COLOR_RESET)"
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "$(COLOR_GREEN)✅ Coverage report: coverage.html$(COLOR_RESET)"

clean: ## Clean build files
	@echo "$(COLOR_YELLOW)Cleaning...$(COLOR_RESET)"
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
	@echo "$(COLOR_GREEN)✅ Cleaned$(COLOR_RESET)"

swagger: ## Generate Swagger documentation
	@echo "$(COLOR_GREEN)Generating Swagger documentation...$(COLOR_RESET)"
	swag init -g cmd/server/main.go -o docs
	@echo "$(COLOR_GREEN)✅ Swagger docs generated$(COLOR_RESET)"

swagger-open: ## Open Swagger UI in browser
	@echo "$(COLOR_GREEN)Opening Swagger UI...$(COLOR_RESET)"
	@open http://localhost:$(APP_PORT)/swagger/index.html || xdg-open http://localhost:$(APP_PORT)/swagger/index.html

# Database Migration Commands
migrate-up: ## Run all pending migrations
	@echo "$(COLOR_GREEN)Running migrations...$(COLOR_RESET)"
	migrate -path migrations -database "$(DB_URL)" up
	@echo "$(COLOR_GREEN)✅ Migrations completed$(COLOR_RESET)"

migrate-down: ## Rollback last migration
	@echo "$(COLOR_YELLOW)Rolling back last migration...$(COLOR_RESET)"
	migrate -path migrations -database "$(DB_URL)" down 1
	@echo "$(COLOR_GREEN)✅ Rollback completed$(COLOR_RESET)"

migrate-drop: ## Drop all migrations (WARNING: deletes all data)
	@echo "$(COLOR_YELLOW)⚠️  WARNING: This will drop all tables and data!$(COLOR_RESET)"
	@read -p "Are you sure? [y/N]: " confirm && [ "$$confirm" = "y" ]
	migrate -path migrations -database "$(DB_URL)" drop
	@echo "$(COLOR_GREEN)✅ All migrations dropped$(COLOR_RESET)"

migrate-create: ## Create a new migration (usage: make migrate-create name=create_table_name)
	@if [ -z "$(name)" ]; then \
		echo "$(COLOR_YELLOW)Error: Migration name required. Usage: make migrate-create name=your_migration_name$(COLOR_RESET)"; \
		exit 1; \
	fi
	@echo "$(COLOR_GREEN)Creating migration: $(name)$(COLOR_RESET)"
	migrate create -ext sql -dir migrations -seq $(name)
	@echo "$(COLOR_GREEN)✅ Migration files created$(COLOR_RESET)"

migrate-force: ## Force migration version (usage: make migrate-force version=1)
	@if [ -z "$(version)" ]; then \
		echo "$(COLOR_YELLOW)Error: Version required. Usage: make migrate-force version=1$(COLOR_RESET)"; \
		exit 1; \
	fi
	@echo "$(COLOR_YELLOW)Forcing migration to version $(version)...$(COLOR_RESET)"
	migrate -path migrations -database "$(DB_URL)" force $(version)
	@echo "$(COLOR_GREEN)✅ Migration version forced$(COLOR_RESET)"

migrate-version: ## Show current migration version
	@echo "$(COLOR_GREEN)Current migration version:$(COLOR_RESET)"
	migrate -path migrations -database "$(DB_URL)" version

# Database Commands
db-create: ## Create database
	@echo "$(COLOR_GREEN)Creating database $(DB_NAME)...$(COLOR_RESET)"
	createdb -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) $(DB_NAME)
	@echo "$(COLOR_GREEN)✅ Database created$(COLOR_RESET)"

db-drop: ## Drop database (WARNING: deletes all data)
	@echo "$(COLOR_YELLOW)⚠️  WARNING: This will delete the entire database!$(COLOR_RESET)"
	@read -p "Are you sure? [y/N]: " confirm && [ "$$confirm" = "y" ]
	dropdb -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) $(DB_NAME)
	@echo "$(COLOR_GREEN)✅ Database dropped$(COLOR_RESET)"

db-reset: db-drop db-create migrate-up ## Reset database (drop, create, migrate)
	@echo "$(COLOR_GREEN)✅ Database reset completed$(COLOR_RESET)"

# Docker Commands
docker-build: ## Build Docker image
	@echo "$(COLOR_GREEN)Building Docker image...$(COLOR_RESET)"
	docker build -t $(APP_NAME):$(VERSION) .
	@echo "$(COLOR_GREEN)✅ Docker image built$(COLOR_RESET)"

docker-run: ## Run Docker container
	@echo "$(COLOR_GREEN)Running Docker container...$(COLOR_RESET)"
	docker run -p $(APP_PORT):$(APP_PORT) --env-file .env $(APP_NAME):$(VERSION)

docker-compose-up: ## Start all services with docker-compose
	@echo "$(COLOR_GREEN)Starting services with docker-compose...$(COLOR_RESET)"
	docker-compose up -d
	@echo "$(COLOR_GREEN)✅ Services started$(COLOR_RESET)"

docker-compose-down: ## Stop all services
	@echo "$(COLOR_YELLOW)Stopping services...$(COLOR_RESET)"
	docker-compose down
	@echo "$(COLOR_GREEN)✅ Services stopped$(COLOR_RESET)"

docker-compose-logs: ## View docker-compose logs
	docker-compose logs -f

# Setup Commands
setup: install-tools deps db-create migrate-up swagger ## Complete project setup
	@echo "$(COLOR_GREEN)✅ Project setup completed!$(COLOR_RESET)"
	@echo "$(COLOR_BLUE)Run 'make run' to start the server$(COLOR_RESET)"

init: ## Initialize new project from scratch
	@echo "$(COLOR_GREEN)Initializing project...$(COLOR_RESET)"
	cp .env.example .env
	@echo "$(COLOR_YELLOW)Please configure .env file before continuing$(COLOR_RESET)"
	@echo "$(COLOR_GREEN)Then run: make setup$(COLOR_RESET)"

# Linting and Formatting
lint: ## Run golangci-lint
	@echo "$(COLOR_GREEN)Running linter...$(COLOR_RESET)"
	golangci-lint run ./...

fmt: ## Format code
	@echo "$(COLOR_GREEN)Formatting code...$(COLOR_RESET)"
	$(GOCMD) fmt ./...
	gofmt -s -w .

# Project Information
info: ## Show project information
	@echo "$(COLOR_BOLD)Project Information:$(COLOR_RESET)"
	@echo "  Name:     $(APP_NAME)"
	@echo "  Version:  $(VERSION)"
	@echo "  Go:       $(shell go version)"
	@echo "  Port:     $(APP_PORT)"
	@echo "  Database: $(DB_NAME)"
	@echo ""

# Default target
.DEFAULT_GOAL := help