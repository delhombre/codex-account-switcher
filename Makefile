.PHONY: build install test lint clean

BINARY_NAME=cxa
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS=-ldflags "-X main.version=$(VERSION)"

# Build the binary
build:
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) ./cmd/cxa

# Install to GOPATH/bin
install:
	go install $(LDFLAGS) ./cmd/cxa

# Run all tests
test:
	go test -v -race -cover ./...

# Run linter
lint:
	golangci-lint run

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Run the TUI in development
dev:
	go run ./cmd/cxa

# Build for all platforms
release:
	goreleaser release --snapshot --clean
