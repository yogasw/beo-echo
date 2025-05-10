.PHONY: build run clean dev test

# Binary name
BINARY_NAME=mockoon-control-panel

# Default target
all: build

# Build the application
build:
	@echo "Building..."
	cd .. && go build -o backend_new/bin/$(BINARY_NAME) backend_new/main.go

# Run the application
run: build
	@echo "Running..."
	./bin/$(BINARY_NAME)

# Clean build files
clean:
	@echo "Cleaning..."
	go clean
	rm -rf bin/

# Run with hot-reload using air (requires air: go install github.com/cosmtrek/air@latest)
dev:
	@echo "Running in development mode..."
	air

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download

# Build for multiple platforms
build-all: clean
	@echo "Building for multiple platforms..."
	GOOS=linux GOARCH=amd64 go build -o bin/$(BINARY_NAME)-linux-amd64 main.go
	GOOS=darwin GOARCH=amd64 go build -o bin/$(BINARY_NAME)-darwin-amd64 main.go
	GOOS=windows GOARCH=amd64 go build -o bin/$(BINARY_NAME)-windows-amd64.exe main.go
