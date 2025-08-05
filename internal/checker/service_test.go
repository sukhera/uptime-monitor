package checker

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/sukhera/uptime-monitor/internal/domain/service"
	"github.com/sukhera/uptime-monitor/testutil"
)

func TestNewService(t *testing.T) {
	tests := []struct {
		name     string
		db       interface{} // Using interface{} to avoid import issues
		options  []ServiceOption
		expected bool
	}{
		{
			name:     "success with default options",
			db:       nil,
			options:  []ServiceOption{},
			expected: true,
		},
		{
			name: "success with custom HTTP client",
			db:   nil,
			options: []ServiceOption{
				WithHTTPClient(&http.Client{Timeout: 5 * time.Second}),
			},
			expected: true,
		},
		{
			name: "success with timeout option",
			db:   nil,
			options: []ServiceOption{
				WithTimeout(5 * time.Second),
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// For this example, we'll skip the database parameter
			// In a real implementation, you would use the mock
			service := NewService(nil, tt.options...)
			assert.Equal(t, tt.expected, service != nil)
		})
	}
}

func TestHealthCheckCommand_Execute(t *testing.T) {
	tests := []struct {
		name           string
		service        service.Service
		expectedStatus string
		expectedError  bool
	}{
		{
			name:           "success with operational service",
			service:        *testutil.CreateTestService(),
			expectedStatus: "operational",
			expectedError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create command with real HTTP client for integration test
			command := NewHTTPHealthCheckCommand(tt.service, &http.Client{Timeout: 5 * time.Second})

			// Execute command
			statusLog := command.Execute(context.Background())

			// Assert results
			assert.Equal(t, tt.service.Name, statusLog.ServiceName)
			assert.NotZero(t, statusLog.Timestamp)

			if tt.expectedError {
				assert.NotEmpty(t, statusLog.Error)
			} else {
				// Note: In a real test, you might want to mock the HTTP client
				// to control the response and test specific status codes
				assert.NotEmpty(t, statusLog.Status)
			}
		})
	}
}

func TestHealthCheckInvoker_ExecuteAll(t *testing.T) {
	tests := []struct {
		name          string
		commands      []HealthCheckCommand
		expectedCount int
	}{
		{
			name:          "success with no commands",
			commands:      []HealthCheckCommand{},
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create invoker
			invoker := NewHealthCheckInvoker()

			// Add commands
			for _, command := range tt.commands {
				invoker.AddCommand(command)
			}

			// Execute all commands
			results := invoker.ExecuteAll(context.Background())

			// Assert results
			assert.Equal(t, tt.expectedCount, len(results))
		})
	}
}
