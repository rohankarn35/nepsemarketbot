# Makefile for building Go project

# Project name
PROJECT_NAME := nepsetelegrambot

# Go source files
SRC := main.go

# Output directories
BUILD_DIR := build

# Platforms
PLATFORMS := windows linux darwin

# Default target
.PHONY: all
all: clean build

# Clean build directory
.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)

# Build for all platforms
.PHONY: build
build: $(PLATFORMS)

# Build for specific platforms
.PHONY: windows linux darwin
windows:
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(PROJECT_NAME).exe $(SRC)
linux:
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(PROJECT_NAME)_linux $(SRC)
darwin:
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(PROJECT_NAME)_darwin $(SRC)

# Run the Go program
.PHONY: run
run:
	go run $(SRC)
