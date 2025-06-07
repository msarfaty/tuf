# Makefile for tuf CLI Go project

APP_NAME := tuf
PKG := ./...
BIN_DIR := bin
BIN := $(BIN_DIR)/$(APP_NAME)

GO_FILES := $(shell find . -type f -name '*.go' -not -path "./vendor/*")

.DEFAULT_GOAL := help

.PHONY: help
help:  ## Show this help.
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## Build the CLI binary
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN) .

.PHONY: run
run: build ## Build and run
	@$(BIN)

.PHONY: clean
clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf $(BIN_DIR) coverage.out

.PHONY: fmt
fmt: ## Format code
	@echo "Formatting..."
	go fmt $(PKG)

.PHONY: lint
lint: ## Lint code (requires golangci-lint)
	@echo "Linting..."
	golangci-lint run

.PHONY: test
test: ## Run tests
	@echo "Running tests..."
	go test -v $(PKG)

.PHONY: cover
cover: ## Run tests with coverage
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out $(PKG)
	go tool cover -func=coverage.out

.PHONY: deps
deps: ## Install dependencies (e.g., golangci-lint)
	@echo "Installing dependencies..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

