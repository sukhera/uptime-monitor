# Development Guide

## Pre-PR Checklist

Before creating a pull request, always run the following commands:

```bash
# Run tests
go test ./...

# Run linting and formatting
gofmt -w .
golangci-lint run

# Run vet for static analysis
go vet ./...
```

## Test Commands

- `go test ./...` - Run all tests
- `go test -v ./...` - Run tests with verbose output
- `go test -race ./...` - Run tests with race detection

## Build Commands

- `go build ./...` - Build all packages
- `go build ./cmd/api` - Build the API server
- `make build` - Build all binaries (if Makefile exists)

## Development Commands

- `go run cmd/api/main.go` - Run the API server directly
- `go mod tidy` - Clean up dependencies

## Code Quality

- Use `gofmt` to format code
- Run `golangci-lint` before committing
- Ensure all tests pass
- Use secure random number generation (`crypto/rand` instead of `math/rand`)