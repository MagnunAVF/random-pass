.PHONY: help build run test test-verbose test-coverage lint fmt vet clean docker-up docker-down install-tools docker-test-up docker-test-down test-integration docker-build docker-prod-up docker-prod-down

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
	@echo "  make docker-up        - Start docker services (dev profile)"
	@echo "  make docker-down      - Stop docker services (dev profile)"
	@echo "  make docker-test-up   - Start docker services (test profile)"
	@echo "  make docker-test-down - Stop docker services (test profile)"
	@echo "  make test-integration - Start test profile and run tests"
	@echo "  make docker-build     - Build the production Docker image"
	@echo "  make docker-prod-up   - Build and start all services (prod profile)"
	@echo "  make docker-prod-down - Stop all services (prod profile)"
	@echo "  make db-shell       - Open psql shell (requires docker-up)"
	@echo "  make install-tools  - Install development tools"
	@echo "  make all            - Run fmt, vet, lint, test, and build"
	@echo ""
	@echo "Env vars (with defaults):"
	@echo "  POSTGRES_DSN  postgres://postgres:postgres@localhost:5432/randompass?sslmode=disable"
	@echo "  REDIS_ADDR    localhost:6379"
	@echo "  PORT          3000"
	@echo "  JWT_SECRET    change-me-in-production"

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
	docker compose --profile dev up -d
	@echo "Docker services started"

docker-test-up:
	@echo "Starting Docker test services (Redis)..."
	docker compose --profile test up -d
	@echo "Docker test services started"

docker-down:
	@echo "Stopping Docker services..."
	docker compose --profile dev down
	@echo "Docker services stopped"

docker-test-down:
	@echo "Stopping Docker test services..."
	docker compose --profile test down
	@echo "Docker test services stopped"

test-integration: docker-test-up
	@echo "Running integration tests..."
	$(GOTEST) ./...

docker-build:
	@echo "Building production Docker image..."
	docker build -t $(BINARY_NAME):latest .
	@echo "Image built: $(BINARY_NAME):latest"

docker-prod-up:
	@echo "Starting production services..."
	docker compose --profile prod up -d --build
	@echo "Production services started on http://localhost:3000"

docker-prod-down:
	@echo "Stopping production services..."
	docker compose --profile prod down
	@echo "Production services stopped"

db-shell:
	docker exec -it $(BINARY_NAME)-postgres psql -U postgres -d randompass

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
