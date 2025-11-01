.PHONY: build test clean run

# Build the server binary
build:
	go build -o bin/server cmd/server/main.go

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Run the server
run:
	go run cmd/server/main.go

# Install dependencies
deps:
	go mod download

# Tidy dependencies
tidy:
	go mod tidy

# Format code
fmt:
	go fmt ./...


# Build and run
all: build run