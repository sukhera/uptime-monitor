package checker

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
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
		if _, err := w.Write([]byte("OK")); err != nil {
			t.Errorf("Failed to write response: %v", err)
		}
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
		if _, err := w.Write([]byte("Created")); err != nil {
			t.Errorf("Failed to write response: %v", err)
		}
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
		if _, err := w.Write([]byte("Internal Server Error")); err != nil {
			t.Errorf("Failed to write response: %v", err)
		}
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
	assert.Contains(t, result.Error, "request failed after 3 attempts")
}

func TestHTTPHealthCheckCommand_Execute_WithHeaders(t *testing.T) {
	var receivedHeaders http.Header
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedHeaders = r.Header
		w.WriteHeader(200)
		if _, err := w.Write([]byte("OK")); err != nil {
			t.Errorf("Failed to write response: %v", err)
		}
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

func setupTestServer(t *testing.T) (*httptest.Server, *HealthCheckInvoker, HealthCheckCommand) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		if _, err := w.Write([]byte("OK")); err != nil {
			t.Errorf("Failed to write response: %v", err)
		}
	}))

	invoker := NewHealthCheckInvoker()
	service := service.Service{
		Name:           "test-service",
		URL:            server.URL,
		ExpectedStatus: 200,
	}
	client := &http.Client{Timeout: 5 * time.Second}
	command := NewHTTPHealthCheckCommand(service, client)

	invoker.AddCommand(command)
	return server, invoker, command
}

func TestHealthCheckInvoker_ExecuteAll_SingleCommand(t *testing.T) {
	server, invoker, _ := setupTestServer(t)
	defer server.Close()

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
		if _, err := w.Write([]byte("OK")); err != nil {
			t.Errorf("Failed to write response: %v", err)
		}
	}))
	defer server1.Close()

	server2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		if _, err := w.Write([]byte("Error")); err != nil {
			t.Errorf("Failed to write response: %v", err)
		}
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
		switch result.ServiceName {
		case "service-1":
			result1 = result
		case "service-2":
			result2 = result
		}
	}

	assert.Equal(t, "operational", result1.Status)
	assert.Equal(t, "down", result2.Status)
}

func TestHealthCheckInvoker_ExecuteSequential(t *testing.T) {
	server, invoker, _ := setupTestServer(t)
	defer server.Close()

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
	// Track concurrent requests to verify they happen simultaneously
	var mu sync.Mutex
	requestTimes := make([]time.Time, 0)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Record when this request started
		mu.Lock()
		requestTimes = append(requestTimes, time.Now())
		mu.Unlock()

		// Simulate some processing time
		time.Sleep(50 * time.Millisecond)
		w.WriteHeader(200)
		if _, err := w.Write([]byte("OK")); err != nil {
			t.Errorf("Failed to write response: %v", err)
		}
	}))
	defer server.Close()

	// Create invoker with no jitter to avoid random delays that break timing tests
	config := DefaultWorkerPoolConfig()
	config.JitterMaxDuration = 0 // Remove jitter for predictable timing
	config.WorkerCount = 5       // Ensure we have enough workers
	invoker := NewHealthCheckInvokerWithConfig(config)

	// Add multiple commands
	numServices := 3
	for i := 0; i < numServices; i++ {
		service := service.Service{
			Name:           fmt.Sprintf("service-%d", i),
			URL:            server.URL,
			ExpectedStatus: 200,
		}
		client := &http.Client{Timeout: 5 * time.Second}
		command := NewHTTPHealthCheckCommand(service, client)
		invoker.AddCommand(command)
	}

	ctx := context.Background()
	results := invoker.ExecuteAll(ctx)

	// Test that we got all results
	assert.Len(t, results, numServices)

	// Test that all results are successful
	for _, result := range results {
		assert.Equal(t, "operational", result.Status)
	}

	// Test concurrency by checking that requests started close together
	mu.Lock()
	times := make([]time.Time, len(requestTimes))
	copy(times, requestTimes)
	mu.Unlock()

	assert.Len(t, times, numServices, "Expected %d concurrent requests", numServices)

	if len(times) >= 2 {
		// Find the time span between first and last request start
		var earliest, latest time.Time
		earliest = times[0]
		latest = times[0]

		for _, t := range times {
			if t.Before(earliest) {
				earliest = t
			}
			if t.After(latest) {
				latest = t
			}
		}

		// If requests are truly concurrent, they should start within a small time window
		// Allow more generous time for race detection overhead
		maxStartSpread := 500 * time.Millisecond // Increased for race detection
		actualSpread := latest.Sub(earliest)

		assert.True(t, actualSpread < maxStartSpread,
			"Requests should start concurrently. Time spread: %v (max allowed: %v)",
			actualSpread, maxStartSpread)
	}
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
	assert.Contains(t, result.Error, "request failed after 3 attempts")
}

func TestHealthCheckInvoker_Integration(t *testing.T) {
	// Test complete workflow with multiple services
	server1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		if _, err := w.Write([]byte("OK")); err != nil {
			t.Errorf("Failed to write response: %v", err)
		}
	}))
	defer server1.Close()

	server2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(201)
		if _, err := w.Write([]byte("Created")); err != nil {
			t.Errorf("Failed to write response: %v", err)
		}
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
		switch result.Status {
		case "operational":
			operationalCount++
		case "degraded":
			degradedCount++
		}
	}

	assert.Equal(t, 1, operationalCount)
	assert.Equal(t, 1, degradedCount)
}
