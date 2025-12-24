# order-service - Makefile
#
# order-api - RESTful API with DDD + CQRS Pattern
# Copyright (c) 2024-2026 order-service. All rights reserved.

# Build configuration
BINARY_NAME := order-service
VERSION ?= 1.0.0
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
GO_VERSION := $(shell go version | cut -d ' ' -f 3)

# Directories
BUILD_DIR := ./build
CONFIG_DIR := ./configs
MIGRATION_DIR := ./migrations

# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GOTEST := $(GOCMD) test
GOMOD := $(GOCMD) mod

# Build flags
LDFLAGS := -s -w \
	-X 'main.Version=$(VERSION)' \
	-X 'main.GitCommit=$(GIT_COMMIT)' \
	-X 'main.BuildTime=$(BUILD_TIME)'

# Colors
GREEN := \033[0;32m
YELLOW := \033[0;33m
RED := \033[0;31m
NC := \033[0m

.PHONY: all build run test clean help deps lint fmt migrate-up migrate-down docker-build

all: build

help:
	@echo "$(GREEN)order-service - Build System$(NC)"
	@echo ""
	@echo "$(YELLOW)Build Commands:$(NC)"
	@echo "  make build        - Build the application"
	@echo "  make run          - Run the application"
	@echo "  make dev          - Run with hot reload (requires air)"
	@echo ""
	@echo "$(YELLOW)Testing Commands:$(NC)"
	@echo "  make test         - Run all tests"
	@echo "  make test-unit    - Run unit tests"
	@echo "  make test-integration - Run integration tests"
	@echo "  make test-coverage - Run tests with coverage"
	@echo ""
	@echo "$(YELLOW)Database Commands:$(NC)"
	@echo "  make migrate-up   - Run database migrations"
	@echo "  make migrate-down - Rollback database migrations"
	@echo "  make migrate-create NAME=name - Create new migration"
	@echo ""
	@echo "$(YELLOW)Code Quality:$(NC)"
	@echo "  make lint         - Run linter"
	@echo "  make fmt          - Format code"
	@echo "  make vet          - Run go vet"
	@echo ""
	@echo "$(YELLOW)Docker Commands:$(NC)"
	@echo "  make docker-build - Build Docker image"
	@echo "  make docker-run   - Run Docker container"
	@echo "  make docker-compose-up - Start all services"
	@echo ""
	@echo "$(YELLOW)Documentation:$(NC)"
	@echo "  make docs         - Generate API documentation"
	@echo "  make swagger      - Open Swagger UI"

build:
	@echo "$(GREEN)Building $(BINARY_NAME)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/api
	@echo "$(GREEN)Build complete: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"

run: build
	@echo "$(GREEN)Starting $(BINARY_NAME)...$(NC)"
	$(BUILD_DIR)/$(BINARY_NAME)

dev:
	@echo "$(GREEN)Starting in development mode...$(NC)"
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "$(YELLOW)air not installed. Install with: go install github.com/cosmtrek/air@latest$(NC)"; \
		$(GOCMD) run ./cmd/api; \
	fi

test:
	@echo "$(GREEN)Running tests...$(NC)"
	$(GOTEST) -v -race ./...

test-unit:
	@echo "$(GREEN)Running unit tests...$(NC)"
	$(GOTEST) -v -short ./tests/unit/...

test-integration:
	@echo "$(GREEN)Running integration tests...$(NC)"
	$(GOTEST) -v ./tests/integration/...

test-coverage:
	@echo "$(GREEN)Running tests with coverage...$(NC)"
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report: coverage.html$(NC)"

deps:
	@echo "$(GREEN)Downloading dependencies...$(NC)"
	$(GOMOD) download
	$(GOMOD) tidy

lint:
	@echo "$(GREEN)Running linter...$(NC)"
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "$(YELLOW)golangci-lint not installed$(NC)"; \
	fi

fmt:
	@echo "$(GREEN)Formatting code...$(NC)"
	$(GOCMD) fmt ./...

vet:
	@echo "$(GREEN)Running go vet...$(NC)"
	$(GOCMD) vet ./...

migrate-up:
	@echo "$(GREEN)Running migrations...$(NC)"
	@if command -v migrate > /dev/null; then \
		migrate -path $(MIGRATION_DIR) -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" up; \
	else \
		echo "$(YELLOW)migrate not installed. Install with: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest$(NC)"; \
	fi

migrate-down:
	@echo "$(GREEN)Rolling back migrations...$(NC)"
	migrate -path $(MIGRATION_DIR) -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" down 1

migrate-create:
	@echo "$(GREEN)Creating migration: $(NAME)...$(NC)"
	migrate create -ext sql -dir $(MIGRATION_DIR) -seq $(NAME)

clean:
	@echo "$(GREEN)Cleaning...$(NC)"
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
	$(GOCMD) clean

docker-build:
	@echo "$(GREEN)Building Docker image...$(NC)"
	docker build -t order-service:$(VERSION) .

docker-run:
	@echo "$(GREEN)Running Docker container...$(NC)"
	docker run -p 8080:8080 --env-file .env order-service:$(VERSION)

docker-compose-up:
	@echo "$(GREEN)Starting services...$(NC)"
	docker compose --profile all up -d

docker-compose-down:
	@echo "$(GREEN)Stopping services...$(NC)"
	docker compose --profile all down

docs:
	@echo "$(GREEN)Documentation available at:$(NC)"
	@echo "  - OpenAPI: docs/api/openapi.yaml"
	@echo "  - Swagger: docs/api/swagger.json"
	@echo "  - ERD: docs/diagrams/ERD.md"
	@echo "  - DFD: docs/diagrams/DFD.md"
	@echo "  - Postman: docs/postman/collection.json"

swagger:
	@echo "$(GREEN)Opening Swagger UI...$(NC)"
	@if command -v swagger > /dev/null; then \
		swagger serve docs/api/swagger.json; \
	else \
		echo "$(YELLOW)swagger-cli not installed$(NC)"; \
		echo "Open docs/api/swagger.json in https://editor.swagger.io"; \
	fi

.DEFAULT_GOAL := help
