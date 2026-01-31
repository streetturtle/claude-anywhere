.PHONY: build install clean test lint

# Build the binary
build:
	go build -o claude-anywhere ./cmd/claude-anywhere

# Build and install to /usr/local/bin
install: build
	sudo mv claude-anywhere /usr/local/bin/claude-anywhere

# Clean build artifacts
clean:
	rm -f claude-anywhere
	go clean

# Run tests
test:
	go test ./...

# Run linter
lint:
	golangci-lint run ./...

# Build for multiple platforms
build-all:
	GOOS=darwin GOARCH=amd64 go build -o claude-anywhere-darwin-amd64 ./cmd/claude-anywhere
	GOOS=darwin GOARCH=arm64 go build -o claude-anywhere-darwin-arm64 ./cmd/claude-anywhere

# Run locally (development)
run:
	go run ./cmd/claude-anywhere

# Install dependencies
deps:
	go mod download
	go mod tidy

# Show help
help:
	@echo "Available targets:"
	@echo "  build      - Build the claude-anywhere binary"
	@echo "  install    - Build and install to /usr/local/bin"
	@echo "  clean      - Remove build artifacts"
	@echo "  test       - Run tests"
	@echo "  lint       - Run golangci-lint"
	@echo "  build-all  - Build for multiple platforms"
	@echo "  run        - Run the CLI locally"
	@echo "  deps       - Install and tidy dependencies"
