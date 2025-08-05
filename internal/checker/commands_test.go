package checker

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/sukhera/uptime-monitor/internal/domain/service"
)

// MockHTTPClient is a simple mock for testing
type MockHTTPClient struct {
	responses map[string]*http.Response
	errors    map[string]error
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if err, exists := m.errors[req.URL.String()]; exists {
		return nil, err
	}
	if resp, exists := m.responses[req.URL.String()]; exists {
		return resp, nil
	}
	return &http.Response{StatusCode: 200}, nil
}

func TestHTTPHealthCheckCommand_New(t *testing.T) {
	service := service.Service{
		Name:           "test-service",
		URL:            "http://test.com",
		ExpectedStatus: 200,
	}
	client := &MockHTTPClient{}

	command := NewHTTPHealthCheckCommand(service, client)

	assert.Equal(t, service, command.service)
	assert.Equal(t, client, command.client)
}

func TestHTTPHealthCheckCommand_GetServiceName(t *testing.T) {
	service := service.Service{
		Name: "test-service",
		URL:  "http://test.com",
	}
	client := &MockHTTPClient{}

	command := NewHTTPHealthCheckCommand(service, client)
	name := command.GetServiceName()

	assert.Equal(t, "test-service", name)
}

func TestHTTPHealthCheckCommand_Execute_Success(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	service := service.Service{
		Name:           "test-service",
		URL:            server.URL,
		ExpectedStatus: 200,
		Headers: map[string]string{
			"User-Agent": "test-agent",
		},
	}
	client := &http.Client{Timeout: 5 * time.Second}

	command := NewHTTPHealthCheckCommand(service, client)
	ctx := context.Background()

	result := command.Execute(ctx)

	assert.Equal(t, "test-service", result.ServiceName)
	assert.Equal(t, "operational", result.Status)
	assert.Equal(t, 200, result.StatusCode)
	assert.True(t, result.Latency >= 0)
	assert.Empty(t, result.Error)
}

func TestHTTPHealthCheckCommand_Execute_ExpectedStatusMismatch(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(201) // Different from expected 200
		w.Write([]byte("Created"))
	}))
	defer server.Close()

	service := service.Service{
		Name:           "test-service",
		URL:            server.URL,
		ExpectedStatus: 200,
	}
	client := &http.Client{Timeout: 5 * time.Second}

	command := NewHTTPHealthCheckCommand(service, client)
	ctx := context.Background()

	result := command.Execute(ctx)

	assert.Equal(t, "test-service", result.ServiceName)
	assert.Equal(t, "degraded", result.Status) // Should be degraded for 2xx but not expected status
	assert.Equal(t, 201, result.StatusCode)
	assert.True(t, result.Latency >= 0) // Latency can be 0 for very fast responses
	assert.Empty(t, result.Error)
}

func TestHTTPHealthCheckCommand_Execute_ErrorStatus(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	service := service.Service{
		Name:           "test-service",
		URL:            server.URL,
		ExpectedStatus: 200,
	}
	client := &http.Client{Timeout: 5 * time.Second}

	command := NewHTTPHealthCheckCommand(service, client)
	ctx := context.Background()

	result := command.Execute(ctx)

	assert.Equal(t, "test-service", result.ServiceName)
	assert.Equal(t, "down", result.Status)
	assert.Equal(t, 500, result.StatusCode)
	assert.True(t, result.Latency >= 0) // Latency can be 0 for very fast responses
	assert.Contains(t, result.Error, "Unexpected status code: 500")
}

func TestHTTPHealthCheckCommand_Execute_NetworkError(t *testing.T) {
	service := service.Service{
		Name:           "test-service",
		URL:            "http://localhost:99999", // Invalid port
		ExpectedStatus: 200,
	}
	client := &http.Client{Timeout: 1 * time.Second}

	command := NewHTTPHealthCheckCommand(service, client)
	ctx := context.Background()

	result := command.Execute(ctx)

	assert.Equal(t, "test-service", result.ServiceName)
	assert.Equal(t, "down", result.Status)
	assert.Equal(t, 0, result.StatusCode)
	assert.True(t, result.Latency >= 0) // Latency can be 0 for very fast responses
	assert.Contains(t, result.Error, "Request failed after 3 attempts")
}

func TestHTTPHealthCheckCommand_Execute_WithHeaders(t *testing.T) {
	var receivedHeaders http.Header
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedHeaders = r.Header
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	service := service.Service{
		Name:           "test-service",
		URL:            server.URL,
		ExpectedStatus: 200,
		Headers: map[string]string{
			"Authorization": "Bearer token123",
			"User-Agent":    "test-agent",
		},
	}
	client := &http.Client{Timeout: 5 * time.Second}

	command := NewHTTPHealthCheckCommand(service, client)
	ctx := context.Background()

	result := command.Execute(ctx)

	assert.Equal(t, "operational", result.Status)
	assert.Equal(t, "Bearer token123", receivedHeaders.Get("Authorization"))
	assert.Equal(t, "test-agent", receivedHeaders.Get("User-Agent"))
}

func TestHealthCheckInvoker_New(t *testing.T) {
	invoker := NewHealthCheckInvoker()

	assert.NotNil(t, invoker)
	assert.NotNil(t, invoker.commands)
	assert.Empty(t, invoker.commands)
}

func TestHealthCheckInvoker_AddCommand(t *testing.T) {
	invoker := NewHealthCheckInvoker()
	service := service.Service{Name: "test-service", URL: "http://test.com"}
	client := &MockHTTPClient{}
	command := NewHTTPHealthCheckCommand(service, client)

	invoker.AddCommand(command)

	assert.Len(t, invoker.commands, 1)
	assert.Equal(t, command, invoker.commands[0])
}

func TestHealthCheckInvoker_ExecuteAll_Empty(t *testing.T) {
	invoker := NewHealthCheckInvoker()
	ctx := context.Background()

	results := invoker.ExecuteAll(ctx)

	assert.Empty(t, results)
}

func TestHealthCheckInvoker_ExecuteAll_SingleCommand(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	invoker := NewHealthCheckInvoker()
	service := service.Service{
		Name:           "test-service",
		URL:            server.URL,
		ExpectedStatus: 200,
	}
	client := &http.Client{Timeout: 5 * time.Second}
	command := NewHTTPHealthCheckCommand(service, client)

	invoker.AddCommand(command)
	ctx := context.Background()

	results := invoker.ExecuteAll(ctx)

	assert.Len(t, results, 1)
	assert.Equal(t, "test-service", results[0].ServiceName)
	assert.Equal(t, "operational", results[0].Status)
}

func TestHealthCheckInvoker_ExecuteAll_MultipleCommands(t *testing.T) {
	// Create multiple test servers
	server1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	}))
	defer server1.Close()

	server2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		w.Write([]byte("Error"))
	}))
	defer server2.Close()

	invoker := NewHealthCheckInvoker()

	// Add first command
	service1 := service.Service{
		Name:           "service-1",
		URL:            server1.URL,
		ExpectedStatus: 200,
	}
	client1 := &http.Client{Timeout: 5 * time.Second}
	command1 := NewHTTPHealthCheckCommand(service1, client1)
	invoker.AddCommand(command1)

	// Add second command
	service2 := service.Service{
		Name:           "service-2",
		URL:            server2.URL,
		ExpectedStatus: 200,
	}
	client2 := &http.Client{Timeout: 5 * time.Second}
	command2 := NewHTTPHealthCheckCommand(service2, client2)
	invoker.AddCommand(command2)

	ctx := context.Background()
	results := invoker.ExecuteAll(ctx)

	assert.Len(t, results, 2)

	// Find results by service name
	var result1, result2 service.StatusLog
	for _, result := range results {
		if result.ServiceName == "service-1" {
			result1 = result
		} else if result.ServiceName == "service-2" {
			result2 = result
		}
	}

	assert.Equal(t, "operational", result1.Status)
	assert.Equal(t, "down", result2.Status)
}

func TestHealthCheckInvoker_ExecuteSequential(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	invoker := NewHealthCheckInvoker()
	service := service.Service{
		Name:           "test-service",
		URL:            server.URL,
		ExpectedStatus: 200,
	}
	client := &http.Client{Timeout: 5 * time.Second}
	command := NewHTTPHealthCheckCommand(service, client)

	invoker.AddCommand(command)
	ctx := context.Background()

	results := invoker.ExecuteSequential(ctx)

	assert.Len(t, results, 1)
	assert.Equal(t, "test-service", results[0].ServiceName)
	assert.Equal(t, "operational", results[0].Status)
}

func TestHTTPHealthCheckCommand_Execute_ContextCancellation(t *testing.T) {
	service := service.Service{
		Name:           "test-service",
		URL:            "http://localhost:99999", // Invalid port
		ExpectedStatus: 200,
	}
	client := &http.Client{Timeout: 1 * time.Second}

	command := NewHTTPHealthCheckCommand(service, client)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	result := command.Execute(ctx)

	assert.Equal(t, "test-service", result.ServiceName)
	assert.Equal(t, "down", result.Status)
}

func TestHealthCheckInvoker_ConcurrentExecution(t *testing.T) {
	// Create test server that takes some time to respond
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(50 * time.Millisecond) // Simulate network delay
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	invoker := NewHealthCheckInvoker()

	// Add multiple commands
	for i := 0; i < 5; i++ {
		service := service.Service{
			Name:           "service-" + string(rune(i)),
			URL:            server.URL,
			ExpectedStatus: 200,
		}
		client := &http.Client{Timeout: 5 * time.Second}
		command := NewHTTPHealthCheckCommand(service, client)
		invoker.AddCommand(command)
	}

	ctx := context.Background()
	start := time.Now()
	results := invoker.ExecuteAll(ctx)
	duration := time.Since(start)

	assert.Len(t, results, 5)
	// Concurrent execution should be faster than sequential
	assert.True(t, duration < 300*time.Millisecond) // Should be much faster than 5 * 50ms
}

func TestHTTPHealthCheckCommand_Execute_RetryLogic(t *testing.T) {
	// Test retry logic with network errors (not HTTP status errors)
	// The implementation only retries on network errors, not HTTP status errors
	service := service.Service{
		Name:           "test-service",
		URL:            "http://localhost:99999", // Invalid port to cause network error
		ExpectedStatus: 200,
	}
	client := &http.Client{Timeout: 1 * time.Second}

	command := NewHTTPHealthCheckCommand(service, client)
	ctx := context.Background()

	result := command.Execute(ctx)

	assert.Equal(t, "down", result.Status)
	assert.Contains(t, result.Error, "Request failed after 3 attempts")
}

func TestHealthCheckInvoker_Integration(t *testing.T) {
	// Test complete workflow with multiple services
	server1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	}))
	defer server1.Close()

	server2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(201)
		w.Write([]byte("Created"))
	}))
	defer server2.Close()

	invoker := NewHealthCheckInvoker()

	// Add operational service
	service1 := service.Service{
		Name:           "operational-service",
		URL:            server1.URL,
		ExpectedStatus: 200,
	}
	client1 := &http.Client{Timeout: 5 * time.Second}
	command1 := NewHTTPHealthCheckCommand(service1, client1)
	invoker.AddCommand(command1)

	// Add degraded service
	service2 := service.Service{
		Name:           "degraded-service",
		URL:            server2.URL,
		ExpectedStatus: 200,
	}
	client2 := &http.Client{Timeout: 5 * time.Second}
	command2 := NewHTTPHealthCheckCommand(service2, client2)
	invoker.AddCommand(command2)

	ctx := context.Background()
	results := invoker.ExecuteAll(ctx)

	assert.Len(t, results, 2)

	// Verify results
	operationalCount := 0
	degradedCount := 0
	for _, result := range results {
		if result.Status == "operational" {
			operationalCount++
		} else if result.Status == "degraded" {
			degradedCount++
		}
	}

	assert.Equal(t, 1, operationalCount)
	assert.Equal(t, 1, degradedCount)
}
