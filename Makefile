# Makefile for SHELA - SHELA Helps Evaluate Linux Access

APP_NAME := shela
SRC := main.go
BIN_DIR := bin
BUILD_TIME := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
VERSION := 1.0.0

.PHONY: all build run clean fmt vet lint deps

all: build

build: deps fmt vet
	@echo ">> Building $(APP_NAME)..."
	@mkdir -p $(BIN_DIR)
	@go build -ldflags="-s -w -X main.buildVersion=$(VERSION) -X main.buildTime=$(BUILD_TIME)" -o $(BIN_DIR)/$(APP_NAME) $(SRC)
	@chmod 755 $(BIN_DIR)/$(APP_NAME)

run: build
	@./$(BIN_DIR)/$(APP_NAME)

fmt:
	@echo ">> Running go fmt..."
	@go fmt ./...

vet:
	@echo ">> Running go vet..."
	@go vet ./...

lint:
	@echo ">> Running golint (if installed)..."
	@golint ./... || echo "golint not found. Install with 'go install golang.org/x/lint/golint@latest'"

deps:
	@echo ">> Ensuring dependencies..."
	@go mod tidy

clean:
	@echo ">> Cleaning up..."
	@rm -rf $(BIN_DIR)
