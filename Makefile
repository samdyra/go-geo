# Makefile for auth-api project

# Go related variables
BINARY_NAME=go-geo
MAIN_PACKAGE=./cmd/server

# Docker related variables
POSTGRES_CONTAINER=go-geo

# Load environment variables from .env file
include .env
export

# PHONY targets
.PHONY: build run test clean deps start-db stop-db migrate help

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) $(MAIN_PACKAGE)

# Run the application
run: build
	@echo "Running $(BINARY_NAME)..."
	@./$(BINARY_NAME)

# Run tests
test:
	@echo "Running tests..."
	@go test ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@go clean

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download

# Start PostgreSQL database
start-db:
	@echo "Starting PostgreSQL..."
	@docker run --name go-geo-db \
		-e POSTGRES_USER=$(DB_USER) \
		-e POSTGRES_PASSWORD=$(DB_PASSWORD) \
		-e POSTGRES_DB=$(DB_NAME) \
		-p $(DB_PORT):5432 \
		-d postgis/postgis:13-3.3

# Stop PostgreSQL database
stop-db:
	@echo "Stopping PostgreSQL..."
	@docker stop $(POSTGRES_CONTAINER)
	@docker rm $(POSTGRES_CONTAINER)

# Run database migrations (requires golang-migrate)
migrate-up:
	@echo "Running database migrations..."
	@migrate -path $(PWD)/migrations/up -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" up

migrate-down:
	@echo "Running database migrations..."
	@migrate -path $(PWD)/migrations/down -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" down

# Help command
help:
	@echo "Available commands:"
	@echo "  make build      - Build the application"
	@echo "  make run        - Run the application"
	@echo "  make test       - Run tests"
	@echo "  make clean      - Clean build artifacts"
	@echo "  make deps       - Download dependencies"
	@echo "  make start-db   - Start PostgreSQL database"
	@echo "  make stop-db    - Stop PostgreSQL database"
	@echo "  make migrate    - Run database migrations"

docs:
	@echo "Generating API documentation..."
	@swag init -g ./cmd/server/main.go

# Default target
all: clean build

docs:
	swag init -g cmd/server/main.go


# Add this to run multiple commands in sequence
setup: deps install-swag install-lint migrate docs