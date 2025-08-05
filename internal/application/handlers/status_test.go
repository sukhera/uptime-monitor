package handlers

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sukhera/uptime-monitor/testutil"
)

func TestNewStatusHandler(t *testing.T) {
	tests := []struct {
		name     string
		db       interface{}
		expected bool
	}{
		{
			name:     "success with nil database",
			db:       nil,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewStatusHandler(nil)
			assert.Equal(t, tt.expected, handler != nil)
		})
	}
}

func TestStatusHandler_HealthCheck(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		expectedFields []string
	}{
		{
			name:           "success",
			method:         "GET",
			path:           "/api/health",
			expectedStatus: http.StatusOK,
			expectedFields: []string{"status", "timestamp", "version"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create handler with nil database (for unit testing)
			handler := NewStatusHandler(nil)
			req := testutil.CreateTestHTTPRequest(tt.method, tt.path, nil)
			w := testutil.CreateTestHTTPResponse()

			// Execute handler
			handler.HealthCheck(w, req)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Parse response
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			// Check that expected fields are present
			for _, field := range tt.expectedFields {
				assert.Contains(t, response, field)
			}

			// Check specific values
			assert.Equal(t, "healthy", response["status"])
		})
	}
}

func TestStatusHandler_GetStatus(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{
			name:           "success with nil database",
			method:         "GET",
			path:           "/api/status",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create handler with nil database
			handler := NewStatusHandler(nil)
			req := testutil.CreateTestHTTPRequest(tt.method, tt.path, nil)
			w := testutil.CreateTestHTTPResponse()

			// Execute handler
			handler.GetStatus(w, req)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Assert headers
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
			assert.Equal(t, "no-cache, no-store, must-revalidate", w.Header().Get("Cache-Control"))
			assert.Equal(t, "no-cache", w.Header().Get("Pragma"))
			assert.Equal(t, "0", w.Header().Get("Expires"))

			// Parse response as array (since GetStatus returns an array)
			var response []interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			// Should return an empty array when database is nil
			assert.Equal(t, 0, len(response))
		})
	}
}
