package checker

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/sukhera/uptime-monitor/internal/models"
	"github.com/sukhera/uptime-monitor/mocks"
)

func TestNewService(t *testing.T) {
	mockDB := mocks.NewMockDatabaseInterface(t)

	service := NewService(mockDB)

	assert.NotNil(t, service)
	assert.NotNil(t, service.client)
	assert.Equal(t, 10*time.Second, service.client.Timeout)
	assert.Equal(t, mockDB, service.db)
}

func TestService_checkService(t *testing.T) {
	tests := []struct {
		name           string
		service        models.Service
		serverResponse func(w http.ResponseWriter, r *http.Request)
		expectedStatus string
		expectedError  string
		serverDelay    time.Duration
	}{
		{
			name: "successful check - operational",
			service: models.Service{
				Name:           "test-service",
				URL:            "", // Will be set to test server URL
				ExpectedStatus: 200,
				Headers:        map[string]string{"Authorization": "Bearer token"},
			},
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				// Verify headers were set
				assert.Equal(t, "Bearer token", r.Header.Get("Authorization"))
				w.WriteHeader(200)
				if _, err := w.Write([]byte("OK")); err != nil {
					t.Errorf("Failed to write response: %v", err)
				}
			},
			expectedStatus: "operational",
		},
		{
			name: "successful check - degraded",
			service: models.Service{
				Name:           "test-service",
				URL:            "",
				ExpectedStatus: 200,
				Headers:        map[string]string{},
			},
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(201) // Different from expected but still 2xx
			},
			expectedStatus: "degraded",
		},
		{
			name: "server error - down",
			service: models.Service{
				Name:           "test-service",
				URL:            "",
				ExpectedStatus: 200,
				Headers:        map[string]string{},
			},
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(500)
			},
			expectedStatus: "down",
			expectedError:  "Unexpected status code: 500",
		},
		{
			name: "timeout - down",
			service: models.Service{
				Name:           "test-service",
				URL:            "",
				ExpectedStatus: 200,
				Headers:        map[string]string{},
			},
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				// Simulate slow response
				time.Sleep(2 * time.Second) // Longer than client timeout (1s)
				w.WriteHeader(200)
			},
			serverDelay:    2 * time.Second,
			expectedStatus: "down",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			server := httptest.NewServer(http.HandlerFunc(tt.serverResponse))
			defer server.Close()

			// Update service URL to point to test server
			tt.service.URL = server.URL

			// Create service with mock DB
			mockDB := mocks.NewMockDatabaseInterface(t)
			service := NewService(mockDB)

			// Override client timeout for timeout test
			if tt.serverDelay > 0 {
				service.client.Timeout = 1 * time.Second
			}

			// Run the check
			result := service.checkService(tt.service)

			// Assertions
			assert.Equal(t, tt.service.Name, result.ServiceName)
			assert.Equal(t, tt.expectedStatus, result.Status)

			if tt.expectedError != "" {
				assert.Contains(t, result.Error, tt.expectedError)
			} else if tt.serverDelay == 0 {
				assert.Empty(t, result.Error)
			}

			assert.NotZero(t, result.Timestamp)
			assert.GreaterOrEqual(t, result.Latency, int64(0))

			// For successful requests, verify status code is set
			if tt.expectedStatus != "down" || tt.serverDelay == 0 {
				assert.NotZero(t, result.StatusCode)
			}
		})
	}
}

func TestService_checkService_InvalidURL(t *testing.T) {
	service := models.Service{
		Name:           "test-service",
		URL:            "://invalid-url",
		ExpectedStatus: 200,
		Headers:        map[string]string{},
	}

	mockDB := mocks.NewMockDatabaseInterface(t)
	checker := NewService(mockDB)

	result := checker.checkService(service)

	assert.Equal(t, "test-service", result.ServiceName)
	assert.Equal(t, "down", result.Status)
	assert.Contains(t, result.Error, "Failed to create request")
	assert.Equal(t, int64(0), result.Latency)
	assert.Equal(t, 0, result.StatusCode)
}

func TestService_checkService_WithRetries(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount < 3 {
			// Simulate network error by closing connection
			hj, ok := w.(http.Hijacker)
			if ok {
				conn, _, _ := hj.Hijack()
				conn.Close()
				return
			}
		}
		// Third attempt succeeds
		w.WriteHeader(200)
	}))
	defer server.Close()

	service := models.Service{
		Name:           "test-service",
		URL:            server.URL,
		ExpectedStatus: 200,
	}

	mockDB := mocks.NewMockDatabaseInterface(t)
	checker := NewService(mockDB)

	result := checker.checkService(service)

	assert.Equal(t, "operational", result.Status)
	assert.Equal(t, 3, callCount) // Should have retried
}

func TestService_RunHealthChecks_DatabaseInterface(t *testing.T) {
	// This test demonstrates how the service interacts with the database interface
	// For full testing of RunHealthChecks, you would need to mock mongo.Collection operations
	// which requires additional interface wrapping beyond the scope of this example

	t.Run("service properly uses database interface", func(t *testing.T) {
		mockDB := mocks.NewMockDatabaseInterface(t)

		// Set up expectations for the database interface methods
		mockServicesCollection := &mongo.Collection{}
		mockStatusLogsCollection := &mongo.Collection{}

		mockDB.EXPECT().ServicesCollection().Return(mockServicesCollection).Maybe()
		mockDB.EXPECT().StatusLogsCollection().Return(mockStatusLogsCollection).Maybe()

		service := NewService(mockDB)

		// Verify the service was created with the mock database
		assert.NotNil(t, service)
		assert.Equal(t, mockDB, service.db)

		// Note: Full testing of RunHealthChecks would require mocking mongo.Collection
		// operations, which is beyond the current interface abstraction level
	})
}

func TestService_checkURL(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer server.Close()

	service := models.Service{
		Name:           "test-service",
		URL:            server.URL,
		ExpectedStatus: 200,
	}

	mockDB := mocks.NewMockDatabaseInterface(t)
	checker := NewService(mockDB)

	// Test the concurrent check function
	statusLogs := make(chan models.StatusLog, 1)

	var wgReal sync.WaitGroup
	wgReal.Add(1)

	go checker.checkURL(service, &wgReal, statusLogs)

	// Wait for goroutine to complete
	wgReal.Wait()
	close(statusLogs)

	// Get the result
	select {
	case result := <-statusLogs:
		assert.Equal(t, "test-service", result.ServiceName)
		assert.Equal(t, "operational", result.Status)
	default:
		t.Fatal("No status log received")
	}
}

func TestService_Integration(t *testing.T) {
	// Test that NewService creates a working service
	mockDB := mocks.NewMockDatabaseInterface(t)
	service := NewService(mockDB)

	require.NotNil(t, service)
	require.NotNil(t, service.client)
	require.NotNil(t, service.db)

	// Test a simple service check
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write([]byte("OK"))
	}))
	defer server.Close()

	testService := models.Service{
		Name:           "integration-test",
		URL:            server.URL,
		ExpectedStatus: 200,
		Headers:        map[string]string{"User-Agent": "test-agent"},
	}

	result := service.checkService(testService)

	assert.Equal(t, "integration-test", result.ServiceName)
	assert.Equal(t, "operational", result.Status)
	assert.Equal(t, 200, result.StatusCode)
	assert.Empty(t, result.Error)
	assert.GreaterOrEqual(t, result.Latency, int64(0))
	assert.NotZero(t, result.Timestamp)
}
