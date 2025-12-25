# order-service - Makefile
#
# order-api - RESTful API with DDD + CQRS Pattern
# Copyright (c) 2024-2026 order-service. All rights reserved.

# Build configuration
BINARY_NAME := order-service
VERSION ?= 1.0.0
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "unknown")
PRODUCT_NAME := order-service
BUILD_TIME := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
GO_VERSION := $(shell go version | cut -d ' ' -f 3)

# Directories
BUILD_DIR := ./build
CONFIG_DIR := ./configs
DIST_DIR := ./dist
MIGRATION_DIR := ./migrations

# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod

# Build flags
LDFLAGS := -s -w \
	-X 'main.Version=$(VERSION)' \
	-X 'main.GitCommit=$(GIT_COMMIT)' \
	-X 'main.BuildTime=$(BUILD_TIME)'

# Platforms for cross-compilation
PLATFORMS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64

# Colors for output
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
	@echo "  make build              - Build the application"
	@echo "  make run                - Run the application"
	@echo "  make dev                - Run with hot reload (requires air)"
	@echo ""
	@echo "$(YELLOW)Testing Commands:$(NC)"
	@echo "  make test               - Run all tests"
	@echo "  make test-unit          - Run unit tests"
	@echo "  make test-integration   - Run integration tests"
	@echo "  make test-coverage      - Run tests with coverage"
	@echo ""
	@echo "$(YELLOW)Database Commands:$(NC)"
	@echo "  make migrate-up         - Run database migrations"
	@echo "  make migrate-down       - Rollback database migrations"
	@echo "  make migrate-create NAME=name - Create new migration"
	@echo ""
	@echo "$(YELLOW)Code Quality:$(NC)"
	@echo "  make lint               - Run linter"
	@echo "  make lint-fix           - Run linter with auto-fix"
	@echo "  make fmt                - Format code"
	@echo "  make vet                - Run go vet"
	@echo "  make check              - Run all checks (fmt, vet, lint, test)"
	@echo ""
	@echo "$(YELLOW)Docker Commands:$(NC)"
	@echo "  make docker-build       - Build Docker image"
	@echo "  make docker-run         - Run Docker container"
	@echo "  make docker-compose-up  - Start all services"
	@echo ""
	@echo "$(YELLOW)Documentation:$(NC)"
	@echo "  make docs               - Generate API documentation"
	@echo "  make swagger            - Open Swagger UI"
	@echo ""
	@echo "$(YELLOW)Dependencies:$(NC)"
	@echo "  make deps               - Download dependencies"
	@echo "  make deps-update        - Update dependencies"
	@echo "  make tidy               - Tidy go modules"
	@echo ""
	@echo "$(YELLOW)Other Commands:$(NC)"
	@echo "  make clean              - Clean build artifacts"
	@echo "  make install            - Install binary to /usr/local/bin"
	@echo "  make uninstall          - Uninstall binary"
	@echo "  make docker-build       - Build Docker image"
	@echo "  make docker-push        - Push Docker image"
	@echo "  make version            - Show version information"
	@echo ""
	@echo "$(YELLOW)Configuration:$(NC)"
	@echo "  VERSION=$(VERSION)"
	@echo "  GIT_COMMIT=$(GIT_COMMIT)"
	@echo "  GIT_BRANCH=$(GIT_BRANCH)"
	@echo "  BUILD_TIME=$(BUILD_TIME)"
	@echo "  GO_VERSION=$(GO_VERSION)"

## Build commands
build:
	@echo "$(GREEN)Building $(BINARY_NAME)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/api
	@echo "$(GREEN)Build complete: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"

## Development commands
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

## Test commands
test: test-unit test-integration
	@echo "$(GREEN)All tests completed$(NC)"

test-unit:
	@echo "$(GREEN)Running unit tests...$(NC)"
	@$(GOTEST) -v -timeout 5m -coverprofile=coverage-unit.out ./tests/unit/...

test-integration:
	@echo "$(GREEN)Running integration tests...$(NC)"
	@$(GOTEST) -v -timeout 5m -coverprofile=coverage-integration.out ./tests/integration/...

test-e2e:
	@echo "$(GREEN)Running E2E tests...$(NC)"
	@$(GOTEST) -v -timeout 10m ./tests/e2e/...

test-all: test-unit test-integration test-e2e
	@echo "$(GREEN)All tests completed$(NC)"

test-coverage:
	@echo "$(GREEN)Generating coverage reports...$(NC)"
	@$(GOCMD) tool cover -html=coverage-unit.out -o coverage-unit.html 2>/dev/null || true
	@$(GOCMD) tool cover -html=coverage-integration.out -o coverage-integration.html 2>/dev/null || true
	@echo "$(GREEN)Coverage reports generated$(NC)"

test-short:
	@echo "$(GREEN)Running short tests (skip E2E)...$(NC)"
	@./scripts/test.sh short

bench:
	@echo "$(GREEN)Running benchmarks...$(NC)"
	@$(GOTEST) -bench=. -benchmem ./...

## Dependencies
deps:
	@echo "$(GREEN)Downloading dependencies...$(NC)"
	@$(GOMOD) download
	@$(GOMOD) tidy
	@echo "$(GREEN)Dependencies downloaded$(NC)"

deps-update:
	@echo "$(GREEN)Updating dependencies...$(NC)"
	@$(GOGET) -u ./...
	@$(GOMOD) tidy
	@echo "$(GREEN)Dependencies updated$(NC)"

tidy:
	@echo "$(GREEN)Tidying go modules...$(NC)"
	@$(GOMOD) tidy
	@echo "$(GREEN)Go modules tidied$(NC)"

## Code quality
lint:
	@echo "$(GREEN)Running linter...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "$(YELLOW)golangci-lint not installed, skipping...$(NC)"; \
	fi

lint-fix:
	@echo "$(GREEN)Running linter with auto-fix...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --fix ./...; \
	else \
		echo "$(YELLOW)golangci-lint not installed, skipping...$(NC)"; \
	fi

fmt:
	@echo "$(GREEN)Formatting code...$(NC)"
	@$(GOCMD) fmt ./...
	@echo "$(GREEN)Code formatted$(NC)"

vet:
	@echo "$(GREEN)Running go vet...$(NC)"
	@$(GOCMD) vet ./...
	@echo "$(GREEN)Vet complete$(NC)"

check: fmt vet lint test
	@echo "$(GREEN)All checks passed$(NC)"

## Database
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

## Docker
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

## Documentaion
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

## Version info
version:
	@echo "$(GREEN)$(PRODUCT_NAME)$(NC)"
	@echo "  Version:      $(VERSION)"
	@echo "  Git Commit:   $(GIT_COMMIT)"
	@echo "  Git Branch:   $(GIT_BRANCH)"
	@echo "  Build Time:   $(BUILD_TIME)"
	@echo "  Go Version:   $(GO_VERSION)"

# =============================================================================
# CI-Specific Targets
# =============================================================================
# These targets are optimized for CI/CD pipelines with proper exit codes,
# coverage output, and race detection.

## CI: Check formatting (fails if code needs formatting)
fmt-check:
	@echo "$(GREEN)Checking code formatting...$(NC)"
	@if [ -n "$$(gofmt -l .)" ]; then \
		echo "$(RED)The following files need formatting:$(NC)"; \
		gofmt -l .; \
		exit 1; \
	fi
	@echo "$(GREEN)Code formatting OK$(NC)"

## CI: Run staticcheck
staticcheck:
	@echo "$(GREEN)Running staticcheck...$(NC)"
	@if command -v staticcheck >/dev/null 2>&1; then \
		staticcheck ./...; \
	else \
		echo "$(YELLOW)Installing staticcheck...$(NC)"; \
		go install honnef.co/go/tools/cmd/staticcheck@latest; \
		staticcheck ./...; \
	fi

## CI: Verify dependencies
verify:
	@echo "$(GREEN)Verifying dependencies...$(NC)"
	@$(GOMOD) verify
	@echo "$(GREEN)Dependencies verified$(NC)"

## CI: Download and verify dependencies
deps-verify: deps verify
	@echo "$(GREEN)Dependencies downloaded and verified$(NC)"

## CI: Run unit tests with race detection and coverage
test-unit-ci:
	@echo "$(GREEN)Running unit tests (CI mode)...$(NC)"
	@$(GOTEST) -v -race -timeout 5m -coverprofile=coverage-unit.out -covermode=atomic ./tests/unit/...

## CI: Run integration tests with race detection and coverage
test-integration-ci:
	@echo "$(GREEN)Running integration tests (CI mode)...$(NC)"
	@$(GOTEST) -v -race -timeout 10m -coverprofile=coverage-integration.out -covermode=atomic ./tests/integration/...

## CI: Run E2E tests
test-e2e-ci:
	@echo "$(GREEN)Running E2E tests (CI mode)...$(NC)"
	@$(GOTEST) -v -timeout 15m ./tests/e2e/...

## CI: Run security scan with gosec
security:
	@echo "$(GREEN)Running security scan...$(NC)"
	@if command -v gosec >/dev/null 2>&1; then \
		gosec -no-fail -fmt sarif -out gosec-results.sarif ./...; \
	else \
		echo "$(YELLOW)gosec not installed, skipping...$(NC)"; \
	fi

## CI: Run govulncheck
govulncheck:
	@echo "$(GREEN)Running govulncheck...$(NC)"
	@if command -v govulncheck >/dev/null 2>&1; then \
		govulncheck ./... || true; \
	else \
		echo "$(YELLOW)Installing govulncheck...$(NC)"; \
		go install golang.org/x/vuln/cmd/govulncheck@latest; \
		govulncheck ./... || true; \
	fi

## CI: Merge coverage files
coverage-merge:
	@echo "$(GREEN)Merging coverage files...$(NC)"
	@if command -v gocovmerge >/dev/null 2>&1; then \
		if [ -f coverage-integration.out ]; then \
			gocovmerge coverage-unit.out coverage-integration.out > coverage-merged.out; \
		else \
			cp coverage-unit.out coverage-merged.out; \
		fi; \
	else \
		echo "$(YELLOW)Installing gocovmerge...$(NC)"; \
		go install github.com/wadey/gocovmerge@latest; \
		if [ -f coverage-integration.out ]; then \
			gocovmerge coverage-unit.out coverage-integration.out > coverage-merged.out; \
		else \
			cp coverage-unit.out coverage-merged.out; \
		fi; \
	fi
	@echo "$(GREEN)Coverage merged to coverage-merged.out$(NC)"

## CI: Generate coverage report
coverage-report: coverage-merge
	@echo "$(GREEN)Generating coverage report...$(NC)"
	@$(GOCMD) tool cover -func=coverage-merged.out | tee coverage-summary.txt
	@$(GOCMD) tool cover -html=coverage-merged.out -o coverage.html
	@echo "$(GREEN)Coverage report generated$(NC)"

## CI: Complete lint pipeline
ci-lint: deps-verify fmt-check vet staticcheck lint
	@echo "$(GREEN)CI lint pipeline completed$(NC)"

## CI: Complete test pipeline
ci-test: test-unit-ci test-integration-ci
	@echo "$(GREEN)CI test pipeline completed$(NC)"

## CI: Complete build verification for a specific platform
ci-build:
	@echo "$(GREEN)Building for CI ($(GOOS)/$(GOARCH))...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@OUTPUT="$(BUILD_DIR)/$(BINARY_NAME)-$(GOOS)-$(GOARCH)"; \
	if [ "$(GOOS)" = "windows" ]; then OUTPUT="$${OUTPUT}.exe"; fi; \
	CGO_ENABLED=0 $(GOBUILD) -ldflags "$(LDFLAGS)" -o $${OUTPUT} ./cmd/api; \
	echo "$(GREEN)Built: $${OUTPUT}$(NC)"


.DEFAULT_GOAL := help
