.PHONY: build install clean test

# Build the binary
build:
	go build -o clse ./cmd/clse

# Build and install to /usr/local/bin
install: build
	sudo mv clse /usr/local/bin/clse

# Clean build artifacts
clean:
	rm -f clse
	go clean

# Run tests
test:
	go test ./...

# Build for multiple platforms
build-all:
	GOOS=darwin GOARCH=amd64 go build -o clse-darwin-amd64 ./cmd/clse
	GOOS=darwin GOARCH=arm64 go build -o clse-darwin-arm64 ./cmd/clse
	GOOS=linux GOARCH=amd64 go build -o clse-linux-amd64 ./cmd/clse

# Run locally (development)
run:
	go run ./cmd/clse

# Install dependencies
deps:
	go mod download
	go mod tidy

# Show help
help:
	@echo "Available targets:"
	@echo "  build      - Build the clse binary"
	@echo "  install    - Build and install to /usr/local/bin"
	@echo "  clean      - Remove build artifacts"
	@echo "  test       - Run tests"
	@echo "  build-all  - Build for multiple platforms"
	@echo "  run        - Run the CLI locally"
	@echo "  deps       - Install and tidy dependencies"
