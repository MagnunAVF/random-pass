.PHONY: help build run test test-verbose test-coverage lint fmt vet clean docker-up docker-down install-tools

BINARY_NAME=random-pass
BIN_DIR=bin
COV_DIR=cov
MAIN_PATH=./main.go
GO=go
GOTEST=$(GO) test
GOVET=$(GO) vet
GOFMT=gofmt
GOLINT=golangci-lint

help:
	@echo "Available targets:"
	@echo "  make build          - Build the application binary"
	@echo "  make run            - Run the application"
	@echo "  make test           - Run all tests"
	@echo "  make test-verbose   - Run tests with verbose output"
	@echo "  make test-coverage  - Run tests with coverage report"
	@echo "  make lint           - Run linter (golangci-lint)"
	@echo "  make fmt            - Format code with gofmt"
	@echo "  make vet            - Run go vet"
	@echo "  make clean          - Remove build artifacts"
	@echo "  make docker-up      - Start docker services (Redis)"
	@echo "  make docker-down    - Stop docker services"
	@echo "  make install-tools  - Install development tools"
	@echo "  make all            - Run fmt, vet, lint, test, and build"

build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BIN_DIR)
	$(GO) build -o $(BIN_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BIN_DIR)/$(BINARY_NAME)"

run: build
	@echo "Starting $(BINARY_NAME)..."
	./$(BIN_DIR)/$(BINARY_NAME)

test:
	@echo "Running tests..."
	$(GOTEST) ./...

test-verbose:
	@echo "Running tests (verbose)..."
	$(GOTEST) -v ./...

test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -cover ./...
	@echo "Generating coverage report..."
	@mkdir -p $(COV_DIR)
	$(GOTEST) -coverprofile=$(COV_DIR)/coverage.out ./...
	$(GO) tool cover -html=$(COV_DIR)/coverage.out -o $(COV_DIR)/coverage.html
	@echo "Coverage report generated: $(COV_DIR)/coverage.html"

lint:
	@echo "Running linter..."
	@if command -v $(GOLINT) > /dev/null; then \
		$(GOLINT) run ./...; \
	else \
		echo "golangci-lint not installed. Run 'make install-tools' to install it."; \
		echo "Falling back to go vet..."; \
		$(MAKE) vet; \
	fi

fmt:
	@echo "Formatting code..."
	$(GOFMT) -w .
	@echo "Code formatted"

vet:
	@echo "Running go vet..."
	$(GOVET) ./...

clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BIN_DIR)
	rm -rf $(COV_DIR)
	@echo "Clean complete"

docker-up:
	@echo "Starting Docker services..."
	docker-compose up -d
	@echo "Docker services started"

docker-down:
	@echo "Stopping Docker services..."
	docker-compose down
	@echo "Docker services stopped"

install-tools:
	@echo "Installing development tools..."
	@echo "Installing golangci-lint..."
	@if ! command -v golangci-lint > /dev/null; then \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		echo "golangci-lint installed"; \
	else \
		echo "golangci-lint already installed"; \
	fi

all: fmt vet lint test build
	@echo "All checks passed!"
