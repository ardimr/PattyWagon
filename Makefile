# Simple Makefile for a Go project

# Build the application
all: build test

build:
	@echo "Building..."
	
	
	@go build -o main cmd/api/main.go

# Run the application
run:
	@go run cmd/api/main.go
# Create DB container
docker-run:
	@if docker compose up --build 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up --build; \
	fi

# Shutdown DB container
docker-down:
	@if docker compose down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down; \
	fi

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v
# Integrations Tests for the application
itest:
	@echo "Running integration tests..."
	@go test ./internal/database -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

lint:
	@golangci-lint run

# Live Reload
watch:
	@if command -v air > /dev/null; then \
            air; \
            echo "Watching...";\
        else \
            read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                go install github.com/air-verse/air@latest; \
                air; \
                echo "Watching...";\
            else \
                echo "You chose not to install air. Exiting..."; \
                exit 1; \
            fi; \
        fi

db-migrate-create:
	@goose create $(file) sql

db-migrate-up:
	@goose up

db-migrate-down:
	@goose down

db-generate-sql:
	@sqlc generate

deploy-local:
	@if docker compose --profile deploy up --build 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose --profile deploy up --build; \
	fi

test-file:
	@hey -n 1 -c 1 -q 1 -m POST \
	-H "Content-Type: multipart/form-data; boundary=db7c097185310d799ace6ee193ebad21" -T "multipart/form-data" \
	-D './internal/server/testdata/image-50KB.jpg' 'http://localhost:8080/v1/file'

load-test-file:
	@hey -z 5s -c 100 -m POST -D './internal/server/testdata/image-50KB.jpg' 'http://localhost:8080/v1/file'

.PHONY: all build run test clean watch lint docker-run docker-down itest db-migrate-create db-migrate-up db-migrate-down db-generate-sql
