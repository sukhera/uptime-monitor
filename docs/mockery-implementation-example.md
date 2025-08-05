# Mockery Implementation Example

This document demonstrates how to implement mockery with table-driven tests following Go best practices, based on the patterns from the CLM-SVC example.

## Overview

The implementation follows these key principles:
1. **Table-driven tests** for comprehensive test coverage
2. **Mockery with expecter pattern** for better test readability
3. **Proper mock setup and verification** following Go testing best practices
4. **Integration with CI/CD** for automated mock generation

## Configuration

### .mockery.yaml

```yaml
with-expecter: true
packages:
  github.com/sukhera/uptime-monitor/internal/infrastructure/database:
    config:
      inpackage: true
      mockname: "MockDatabaseInterface"
      filename: "MockDatabaseInterface.go"
    interfaces:
      Interface:
  
  github.com/sukhera/uptime-monitor/internal/infrastructure/database/mongo:
    config:
      inpackage: true
      mockname: "MockInterface"
      filename: "MockInterface.go"
    interfaces:
      Interface:
  
  github.com/sukhera/uptime-monitor/internal/domain/service:
    config:
      inpackage: true
      mockname: "MockRepository"
      filename: "MockRepository.go"
    interfaces:
      Repository:
  
  github.com/sukhera/uptime-monitor/internal/checker:
    config:
      inpackage: true
    interfaces:
      HTTPClient:
      ServiceInterface:
      HealthCheckCommand:
      HealthCheckObserver:
      ServiceOption:
```

### Makefile Commands

```makefile
generate-mocks: ## Generate mocks using Mockery
	@echo "$(BLUE)Generating mocks...$(NC)"
	@go install github.com/vektra/mockery/v2@latest
	@mockery
	@echo "$(GREEN)✓ Mocks generated!$(NC)"

test-with-mocks: ## Run tests with generated mocks
	@echo "$(BLUE)Running tests with mocks...$(NC)"
	@make generate-mocks
	@go test -v ./internal/...
	@echo "$(GREEN)✓ Tests with mocks completed!$(NC)"
```

## Test Structure

### 1. Test Utilities (testutil/helper.go)

```go
package testutil

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sukhera/uptime-monitor/internal/domain/service"
)

// TestContext creates a test context with timeout
func TestContext(t *testing.T) context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(cancel)
	return ctx
}

// CreateTestService creates a test service with default values
func CreateTestService() *service.Service {
	return &service.Service{
		Name:           "Test Service",
		Slug:           "test-service",
		URL:            "https://example.com",
		Headers:        map[string]string{"User-Agent": "TestAgent"},
		ExpectedStatus: 200,
		Enabled:        true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}

// CreateTestStatusLog creates a test status log with default values
func CreateTestStatusLog() *service.StatusLog {
	return &service.StatusLog{
		ServiceName: "Test Service",
		Status:      "operational",
		Latency:     100,
		StatusCode:  200,
		Timestamp:   time.Now(),
	}
}

// CreateTestHTTPRequest creates a test HTTP request
func CreateTestHTTPRequest(method, path string, body io.Reader) *http.Request {
	req := httptest.NewRequest(method, path, body)
	req.Header.Set("Content-Type", "application/json")
	return req
}

// CreateTestHTTPResponse creates a test HTTP response recorder
func CreateTestHTTPResponse() *httptest.ResponseRecorder {
	return httptest.NewRecorder()
}

// AssertJSONResponse asserts that a response body matches expected JSON
func AssertJSONResponse(t *testing.T, responseBody []byte, expected any) {
	var actual map[string]interface{}
	err := json.Unmarshal(responseBody, &actual)
	require.NoError(t, err)

	expectedBytes, err := json.Marshal(expected)
	require.NoError(t, err)

	var expectedMap map[string]interface{}
	err = json.Unmarshal(expectedBytes, &expectedMap)
	require.NoError(t, err)

	assert.Equal(t, expectedMap, actual)
}
```

### 2. Example Test Implementation

#### Handler Tests (internal/application/handlers/status_test.go)

```go
package handlers

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/sukhera/uptime-monitor/internal/infrastructure/database/mongo"
	"github.com/sukhera/uptime-monitor/testutil"
)

func TestStatusHandler_HealthCheck(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		setupMock      func(*mongo.MockInterface)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:   "success",
			method: "GET",
			path:   "/api/health",
			setupMock: func(mockDB *mongo.MockInterface) {
				mockDB.EXPECT().Ping(mock.Anything).Return(nil).Maybe()
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"status": "healthy",
			},
		},
		{
			name:   "database error",
			method: "GET",
			path:   "/api/health",
			setupMock: func(mockDB *mongo.MockInterface) {
				mockDB.EXPECT().Ping(mock.Anything).Return(assert.AnError).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "internal server error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock database
			mockDB := mongo.NewMockInterface(t)
			
			// Setup mock expectations
			if tt.setupMock != nil {
				tt.setupMock(mockDB)
			}
			
			// Create handler with mock
			handler := NewStatusHandler(mockDB)
			req := testutil.CreateTestHTTPRequest(tt.method, tt.path, nil)
			w := testutil.CreateTestHTTPResponse()

			// Execute handler
			handler.HealthCheck(w, req)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)
			
			if tt.expectedBody != nil {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedBody, response)
			}
		})
	}
}
```

#### Service Tests (internal/checker/service_test.go)

```go
package checker

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/sukhera/uptime-monitor/internal/infrastructure/database/mongo"
	"github.com/sukhera/uptime-monitor/testutil"
)

func TestNewService(t *testing.T) {
	tests := []struct {
		name     string
		db       mongo.Interface
		options  []ServiceOption
		expected bool
	}{
		{
			name:     "success with default options",
			db:       &mongo.MockInterface{},
			options:  []ServiceOption{},
			expected: true,
		},
		{
			name: "success with custom HTTP client",
			db:   &mongo.MockInterface{},
			options: []ServiceOption{
				WithHTTPClient(&http.Client{Timeout: 5 * time.Second}),
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewService(tt.db, tt.options...)
			assert.Equal(t, tt.expected, service != nil)
		})
	}
}

func TestService_RunHealthChecks(t *testing.T) {
	tests := []struct {
		name          string
		setupMock     func(*mongo.MockInterface, *MockHTTPClient)
		expectError   bool
		expectedError string
	}{
		{
			name: "success with operational services",
			setupMock: func(mockDB *mongo.MockInterface, mockClient *MockHTTPClient) {
				// Setup database expectations
				mockDB.EXPECT().ServicesCollection().Return(&mongo.Collection{}).Once()
				
				// Setup HTTP client expectations
				mockClient.EXPECT().Do(mock.Anything).Return(&http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(bytes.NewBufferString("OK")),
				}, nil).Once()
			},
			expectError: false,
		},
		{
			name: "database error",
			setupMock: func(mockDB *mongo.MockInterface, mockClient *MockHTTPClient) {
				mockDB.EXPECT().ServicesCollection().Return(nil, assert.AnError).Once()
			},
			expectError:   true,
			expectedError: "error querying services",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mocks
			mockDB := mongo.NewMockInterface(t)
			mockClient := NewMockHTTPClient(t)
			
			// Setup mock expectations
			if tt.setupMock != nil {
				tt.setupMock(mockDB, mockClient)
			}
			
			// Create service with mocks
			service := NewService(mockDB, WithHTTPClient(mockClient))
			
			// Run health checks
			err := service.RunHealthChecks(context.Background())
			
			// Assert results
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
```

## Mock Generation Process

### 1. Generate Mocks

```bash
# Generate all mocks defined in .mockery.yaml
make generate-mocks

# Or run directly
mockery
```

### 2. Verify Generated Mocks

After generation, you should see files like:
- `internal/infrastructure/database/mongo/MockInterface.go`
- `internal/domain/service/MockRepository.go`
- `internal/checker/MockHTTPClient.go`

### 3. Use Mocks in Tests

```go
// Create mock
mockDB := mongo.NewMockInterface(t)

// Setup expectations using expecter pattern
mockDB.EXPECT().ServicesCollection().Return(&mongo.Collection{}).Once()
mockDB.EXPECT().Ping(mock.Anything).Return(nil).Maybe()

// Use mock in your test
handler := NewStatusHandler(mockDB)

// Expectations are automatically verified
```

## Best Practices

### 1. Table-Driven Tests

Always use table-driven tests for comprehensive coverage:

```go
func TestFunction(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		setupMock     func(*MockInterface)
		expected      string
		expectError   bool
	}{
		{
			name:  "success case",
			input: "test",
			setupMock: func(mock *MockInterface) {
				mock.EXPECT().Method(mock.Anything).Return("result", nil).Once()
			},
			expected:    "result",
			expectError: false,
		},
		{
			name:  "error case",
			input: "error",
			setupMock: func(mock *MockInterface) {
				mock.EXPECT().Method(mock.Anything).Return("", assert.AnError).Once()
			},
			expected:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := NewMockInterface(t)
			if tt.setupMock != nil {
				tt.setupMock(mock)
			}
			
			result, err := FunctionUnderTest(mock, tt.input)
			
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
```

### 2. Mock Expectations

Use the expecter pattern for better readability:

```go
// Good - using expecter pattern
mockDB.EXPECT().ServicesCollection().Return(&mongo.Collection{}).Once()
mockDB.EXPECT().Ping(mock.Anything).Return(nil).Maybe()

// Avoid - using old pattern
mockDB.On("ServicesCollection").Return(&mongo.Collection{})
mockDB.On("Ping", mock.Anything).Return(nil)
```

### 3. Mock Verification

Let mockery handle verification automatically:

```go
func TestWithMock(t *testing.T) {
	mockDB := mongo.NewMockInterface(t)
	
	// Setup expectations
	mockDB.EXPECT().ServicesCollection().Return(&mongo.Collection{}).Once()
	
	// Use mock
	handler := NewStatusHandler(mockDB)
	
	// Expectations are automatically verified when mock goes out of scope
}
```

### 4. Test Organization

Organize tests by functionality:

```go
// Group related tests
func TestService_Create(t *testing.T) { /* ... */ }
func TestService_Get(t *testing.T) { /* ... */ }
func TestService_Update(t *testing.T) { /* ... */ }
func TestService_Delete(t *testing.T) { /* ... */ }
```

## Integration with CI/CD

### GitHub Actions

```yaml
# .github/workflows/ci.yml
- name: Generate Mocks
  run: |
    go install github.com/vektra/mockery/v2@latest
    mockery

- name: Run Tests
  run: |
    go test -v -race -coverprofile=coverage.out ./...
```

### Pre-commit Hook

```bash
#!/bin/bash
# scripts/hooks/pre-commit

echo "Generating mocks..."
mockery

echo "Running tests..."
go test -v ./...
```

## Troubleshooting

### Common Issues

1. **Mock not found**
   ```bash
   # Check if interface exists
   go list -f '{{.Dir}}' ./internal/checker
   
   # Regenerate mocks
   mockery --all
   ```

2. **Import errors**
   ```bash
   # Clean and regenerate
   go clean -modcache
   go mod tidy
   mockery
   ```

3. **Configuration errors**
   ```bash
   # Validate configuration
   mockery --config .mockery.yaml --dry-run
   ```

### Debugging

```bash
# Verbose output
mockery --verbose

# Check specific interface
mockery --dir internal/checker --name HTTPClient --output mocks
```

## Example Project Structure

```
project/
├── .mockery.yaml                 # Mockery configuration
├── Makefile                      # Build commands including mock generation
├── testutil/
│   └── helper.go                 # Test utilities
├── internal/
│   ├── application/
│   │   └── handlers/
│   │       ├── status.go
│   │       └── status_test.go    # Handler tests with mocks
│   ├── checker/
│   │   ├── service.go
│   │   └── service_test.go       # Service tests with mocks
│   └── infrastructure/
│       └── database/
│           └── mongo/
│               ├── mongo.go
│               └── MockInterface.go  # Generated mock
└── docs/
    └── mockery-implementation-example.md
```

This implementation provides a solid foundation for using mockery with table-driven tests following Go best practices, similar to the patterns used in the CLM-SVC example. 