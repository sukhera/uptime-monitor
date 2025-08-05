# Mockery Setup Guide

## Overview

This guide explains how to use [Mockery](https://github.com/vektra/mockery) for generating mocks in the Status Page project. Mockery is a mock code autogenerator for Go that works with the `stretchr/testify/mock` package.

## Installation

Mockery is already installed in this project. You can verify the installation:

```bash
mockery --version
# Output: v2.53.4
```

## Configuration

The project uses a `.mockery.yaml` configuration file located in the root directory:

```yaml
with-expecter: true
packages:
  github.com/sukhera/uptime-monitor/internal/database:
    config:
      dir: "mocks"
      mockname: "MockDatabaseInterface"
      filename: "MockDatabaseInterface.go"
    interfaces:
      Interface:
```

### Configuration Options

- **`with-expecter: true`**: Enables the expecter pattern for better test readability
- **`dir: "mocks"`**: Output directory for generated mocks
- **`mockname`**: Custom name for the generated mock
- **`filename`**: Custom filename for the generated mock

## Generating Mocks

### Using Configuration File

```bash
# Generate all mocks defined in .mockery.yaml
mockery

# Generate specific interface
mockery --dir internal/database --name Interface --output mocks
```

### Command Line Options

```bash
# Generate mock for specific interface
mockery --dir <package-dir> --name <interface-name> --output <output-dir>

# Example
mockery --dir internal/database --name Interface --output mocks
```

## Using Generated Mocks

### Basic Usage

```go
package handlers_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/sukhera/uptime-monitor/mocks"
)

func TestWithMock(t *testing.T) {
    // Create mock
    mockDB := mocks.NewMockDatabaseInterface(t)
    
    // Setup expectations
    mockDB.On("ServicesCollection").Return(mockCollection)
    
    // Use mock in your test
    handler := NewStatusHandler(mockDB)
    
    // Verify expectations
    mockDB.AssertExpectations(t)
}
```

### Expecter Pattern (Recommended)

```go
func TestWithExpecter(t *testing.T) {
    // Create mock with expecter
    mockDB := mocks.NewMockDatabaseInterface(t)
    
    // Setup expectations using expecter pattern
    mockDB.EXPECT().ServicesCollection().Return(mockCollection)
    
    // Use mock in your test
    handler := NewStatusHandler(mockDB)
    
    // Expectations are automatically verified
}
```

## Available Mocks

### MockDatabaseInterface

Generated from `internal/database/mongo.go`:

```go
type Interface interface {
    ServicesCollection() *mongo.Collection
    StatusLogsCollection() *mongo.Collection
    Close() error
}
```

**Usage Example:**

```go
func TestStatusHandler_GetStatus(t *testing.T) {
    mockDB := mocks.NewMockDatabaseInterface(t)
    
    // Setup mock expectations
    mockDB.EXPECT().StatusLogsCollection().Return(mockCollection)
    
    handler := NewStatusHandler(mockDB)
    // ... test implementation
}
```

## Best Practices

### 1. Use Integration Tests When Possible

For database operations, prefer integration tests with real database connections:

```go
func TestStatusHandler_Integration(t *testing.T) {
    db, err := database.NewConnection("mongodb://localhost:27017", "test_db")
    if err != nil {
        t.Skip("MongoDB not available for testing")
    }
    defer db.Close()
    
    handler := NewStatusHandler(db)
    // ... test with real database
}
```

### 2. Mock External Dependencies

Use mocks for external services and dependencies:

```go
func TestService_WithMockHTTPClient(t *testing.T) {
    mockClient := mocks.NewMockHTTPClient(t)
    
    // Setup HTTP client expectations
    mockClient.EXPECT().Do(mock.Anything).Return(mockResponse, nil)
    
    service := NewService(mockClient)
    // ... test implementation
}
```

### 3. Test Error Scenarios

Use mocks to test error conditions:

```go
func TestDatabaseError(t *testing.T) {
    mockDB := mocks.NewMockDatabaseInterface(t)
    
    // Setup mock to return error
    mockDB.EXPECT().StatusLogsCollection().Return(nil, assert.AnError)
    
    handler := NewStatusHandler(mockDB)
    // ... test error handling
}
```

### 4. Verify Expectations

Always verify that all mock expectations were met:

```go
func TestWithVerification(t *testing.T) {
    mockDB := mocks.NewMockDatabaseInterface(t)
    
    // Setup expectations
    mockDB.EXPECT().ServicesCollection().Return(mockCollection)
    
    // Run test
    handler := NewStatusHandler(mockDB)
    // ... test implementation
    
    // Verify all expectations were met
    mockDB.AssertExpectations(t)
}
```

## Adding New Mocks

### 1. Create Interface

First, define the interface you want to mock:

```go
// internal/service/service.go
type ServiceInterface interface {
    DoSomething() error
    GetData() ([]byte, error)
}
```

### 2. Update Configuration

Add the new interface to `.mockery.yaml`:

```yaml
packages:
  github.com/sukhera/uptime-monitor/internal/service:
    config:
      dir: "mocks"
      mockname: "MockServiceInterface"
      filename: "MockServiceInterface.go"
    interfaces:
      ServiceInterface:
```

### 3. Generate Mock

```bash
mockery
```

### 4. Use in Tests

```go
func TestWithNewMock(t *testing.T) {
    mockService := mocks.NewMockServiceInterface(t)
    
    mockService.EXPECT().DoSomething().Return(nil)
    mockService.EXPECT().GetData().Return([]byte("test"), nil)
    
    // Use mock in your test
}
```

## Troubleshooting

### Common Issues

1. **Interface Not Found**
   ```
   WRN no such interface interface=ServiceInterface
   ```
   - Check that the interface exists in the specified package
   - Verify the interface name spelling

2. **Import Errors**
   ```go
   undefined: mocks.MockCollection
   ```
   - Mockery only generates mocks for interfaces, not concrete types
   - Use real types or create wrapper interfaces

3. **Configuration Errors**
   ```
   ERR use of unsupported options detected
   ```
   - Check `.mockery.yaml` syntax
   - Remove unsupported options

### Debugging

```bash
# Run with verbose output
mockery --verbose

# Check configuration
mockery --config .mockery.yaml --dry-run
```

## Integration with CI/CD

### GitHub Actions

Add mock generation to your CI pipeline:

```yaml
# .github/workflows/ci.yml
- name: Generate Mocks
  run: |
    go install github.com/vektra/mockery/v2@latest
    mockery
```

### Pre-commit Hook

Add mock generation to your pre-commit hook:

```bash
# scripts/hooks/pre-commit
#!/bin/bash
echo "Generating mocks..."
mockery
```

## Advanced Configuration

### Custom Mock Names

```yaml
packages:
  github.com/sukhera/uptime-monitor/internal/database:
    config:
      mockname: "MockDB"  # Custom mock name
      filename: "db_mock.go"  # Custom filename
```

### Multiple Interfaces

```yaml
packages:
  github.com/sukhera/uptime-monitor/internal/database:
    config:
      dir: "mocks"
    interfaces:
      Interface:
      AnotherInterface:
```

### Package-level Configuration

```yaml
packages:
  github.com/sukhera/uptime-monitor/internal/database:
    config:
      dir: "mocks"
      mockname: "MockDatabaseInterface"
      filename: "MockDatabaseInterface.go"
    interfaces:
      Interface:
  
  github.com/sukhera/uptime-monitor/internal/service:
    config:
      dir: "mocks"
      mockname: "MockServiceInterface"
      filename: "MockServiceInterface.go"
    interfaces:
      ServiceInterface:
```

## Resources

- [Mockery Documentation](https://vektra.github.io/mockery/)
- [Testify Mock Package](https://github.com/stretchr/testify#mock-package)
- [Go Testing Best Practices](https://golang.org/doc/tutorial/add-a-test)

## Example Test Files

See the following files for complete examples:
- `internal/api/handlers/status_test.go` - Integration tests
- `internal/checker/service_test.go` - Unit tests with mocks 