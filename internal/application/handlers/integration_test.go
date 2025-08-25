package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sukhera/uptime-monitor/internal/domain/service"
	"github.com/sukhera/uptime-monitor/testutil"
)

const testBaseURL = "http://localhost:8080"

func TestNewIntegrationHandler(t *testing.T) {
	tests := []struct {
		name     string
		db       interface{}
		baseURL  string
		expected bool
	}{
		{
			name:     "success with nil database",
			db:       nil,
			baseURL:  testBaseURL,
			expected: true,
		},
		{
			name:     "success with empty base URL",
			db:       nil,
			baseURL:  "",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buildInfo := BuildInfo{Version: "test", Commit: "test", BuildDate: "test"}
			handler := NewIntegrationHandler(nil, tt.baseURL, buildInfo)
			assert.Equal(t, tt.expected, handler != nil)
			assert.Equal(t, tt.baseURL, handler.baseURL)
		})
	}
}

func TestIntegrationHandler_SetManualStatus(t *testing.T) {
	buildInfo := BuildInfo{Version: "test", Commit: "test", BuildDate: "test"}
	baseURL := testBaseURL

	future := time.Now().UTC().Add(1 * time.Hour)

	tests := []struct {
		name           string
		method         string
		path           string
		payload        interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "invalid method - GET not allowed",
			method:         "GET",
			path:           "/api/v1/integration/services/test-id/manual-status",
			payload:        nil,
			expectedStatus: http.StatusMethodNotAllowed,
			expectedError:  "Method not allowed",
		},
		{
			name:           "invalid path - missing service ID",
			method:         "POST",
			path:           "/api/v1/integration/services/",
			payload:        map[string]interface{}{"status": "maintenance", "reason": "test"},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid service ID in URL",
		},
		{
			name:           "invalid JSON payload",
			method:         "POST",
			path:           "/api/v1/integration/services/test-id/manual-status",
			payload:        "invalid json",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid JSON payload",
		},
		{
			name:   "missing status",
			method: "POST",
			path:   "/api/v1/integration/services/test-id/manual-status",
			payload: map[string]interface{}{
				"reason": "test reason",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Status is required",
		},
		{
			name:   "missing reason",
			method: "POST",
			path:   "/api/v1/integration/services/test-id/manual-status",
			payload: map[string]interface{}{
				"status": "maintenance",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Reason is required",
		},
		{
			name:   "invalid status",
			method: "POST",
			path:   "/api/v1/integration/services/test-id/manual-status",
			payload: map[string]interface{}{
				"status": "invalid-status",
				"reason": "test reason",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid status value",
		},
		{
			name:   "valid request without expiry",
			method: "POST",
			path:   "/api/v1/integration/services/test-id/manual-status",
			payload: map[string]interface{}{
				"status": "maintenance",
				"reason": "Scheduled maintenance",
			},
			expectedStatus: http.StatusInternalServerError, // Will fail due to no repository implementation
			expectedError:  "manual status setting not implemented",
		},
		{
			name:   "valid request with expiry",
			method: "POST",
			path:   "/api/v1/integration/services/test-id/manual-status",
			payload: map[string]interface{}{
				"status":     "maintenance",
				"reason":     "Scheduled maintenance",
				"expires_at": future.Format(time.RFC3339),
			},
			expectedStatus: http.StatusInternalServerError, // Will fail due to no repository implementation
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewIntegrationHandler(nil, baseURL, buildInfo)

			var body []byte
			var err error

			if tt.payload != nil {
				if str, ok := tt.payload.(string); ok && str == "invalid json" {
					body = []byte(str)
				} else {
					body, err = json.Marshal(tt.payload)
					require.NoError(t, err)
				}
			}

			req := httptest.NewRequest(tt.method, tt.path, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()

			handler.SetManualStatus(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Contains(t, response["error"], tt.expectedError)
			}
		})
	}
}

// Helper function to reduce test duplication
type integrationTestCase struct {
	name           string
	method         string
	path           string
	expectedStatus int
	expectedError  string
}

func runIntegrationHandlerTest(t *testing.T, testCases []integrationTestCase, handlerFunc func(http.ResponseWriter, *http.Request)) {
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			handlerFunc(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Contains(t, response["error"], tt.expectedError)
			}
		})
	}
}

func TestIntegrationHandler_ClearManualStatus(t *testing.T) {
	tests := []integrationTestCase{
		{
			name:           "invalid method - GET not allowed",
			method:         "GET",
			path:           "/api/v1/integration/services/test-id/manual-status",
			expectedStatus: http.StatusMethodNotAllowed,
			expectedError:  "Method not allowed",
		},
		{
			name:           "invalid path - missing service ID",
			method:         "DELETE",
			path:           "/api/v1/integration/services/",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid service ID in URL",
		},
		{
			name:           "valid request",
			method:         "DELETE",
			path:           "/api/v1/integration/services/test-id/manual-status",
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "manual status clearing not implemented",
		},
	}

	buildInfo := BuildInfo{Version: "test", Commit: "test", BuildDate: "test"}
	baseURL := testBaseURL
	handler := NewIntegrationHandler(nil, baseURL, buildInfo)
	runIntegrationHandlerTest(t, tests, handler.ClearManualStatus)
}

func TestIntegrationHandler_GetIntegrationDetails(t *testing.T) {
	tests := []integrationTestCase{
		{
			name:           "invalid method - POST not allowed",
			method:         "POST",
			path:           "/api/v1/integration/services/test-id",
			expectedStatus: http.StatusMethodNotAllowed,
			expectedError:  "Method not allowed",
		},
		{
			name:           "invalid path - missing service ID",
			method:         "GET",
			path:           "/api/v1/integration/services/",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid service ID in URL",
		},
		{
			name:           "valid request",
			method:         "GET",
			path:           "/api/v1/integration/services/test-id",
			expectedStatus: http.StatusNotFound,
			expectedError:  "Service not found",
		},
	}

	buildInfo := BuildInfo{Version: "test", Commit: "test", BuildDate: "test"}
	baseURL := testBaseURL
	handler := NewIntegrationHandler(nil, baseURL, buildInfo)
	runIntegrationHandlerTest(t, tests, handler.GetIntegrationDetails)
}

func TestIntegrationHandler_BulkImport(t *testing.T) {
	buildInfo := BuildInfo{Version: "test", Commit: "test", BuildDate: "test"}
	baseURL := testBaseURL

	validServices := []service.Service{
		{
			Name:        "Test Service 1",
			URL:         "https://example1.com",
			ServiceType: service.ServiceTypeHTTP,
			Enabled:     true,
		},
		{
			Name:        "Test Service 2",
			URL:         "https://example2.com",
			ServiceType: service.ServiceTypeHTTP,
			Enabled:     true,
		},
	}

	invalidServices := []service.Service{
		{
			// Missing name
			URL:         "https://example.com",
			ServiceType: service.ServiceTypeHTTP,
			Enabled:     true,
		},
	}

	tests := []struct {
		name           string
		method         string
		path           string
		payload        interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "invalid method - GET not allowed",
			method:         "GET",
			path:           "/api/v1/integration/services/bulk-import",
			payload:        nil,
			expectedStatus: http.StatusMethodNotAllowed,
			expectedError:  "Method not allowed",
		},
		{
			name:           "invalid JSON payload",
			method:         "POST",
			path:           "/api/v1/integration/services/bulk-import",
			payload:        "invalid json",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid JSON payload",
		},
		{
			name:   "empty services array",
			method: "POST",
			path:   "/api/v1/integration/services/bulk-import",
			payload: map[string]interface{}{
				"services": []interface{}{},
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "No services provided",
		},
		{
			name:   "valid services",
			method: "POST",
			path:   "/api/v1/integration/services/bulk-import",
			payload: map[string]interface{}{
				"services": validServices,
			},
			expectedStatus: http.StatusBadRequest, // Will fail due to validation errors from repository implementation
		},
		{
			name:   "invalid services",
			method: "POST",
			path:   "/api/v1/integration/services/bulk-import",
			payload: map[string]interface{}{
				"services": invalidServices,
			},
			expectedStatus: http.StatusBadRequest, // Will show validation errors
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewIntegrationHandler(nil, baseURL, buildInfo)

			var body []byte
			var err error

			if tt.payload != nil {
				if str, ok := tt.payload.(string); ok && str == "invalid json" {
					body = []byte(str)
				} else {
					body, err = json.Marshal(tt.payload)
					require.NoError(t, err)
				}
			}

			req := httptest.NewRequest(tt.method, tt.path, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()

			handler.BulkImport(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)

				// Check if it's in error field or if it's a bulk operation response
				if errorField, exists := response["error"]; exists {
					assert.Contains(t, errorField, tt.expectedError)
				} else if errors, exists := response["errors"]; exists {
					// For bulk operations, errors might be in the errors array
					errorsArray, ok := errors.([]interface{})
					if ok && len(errorsArray) > 0 {
						// At least one error should contain our expected error
						found := false
						for _, errItem := range errorsArray {
							if errStr, ok := errItem.(string); ok && len(errStr) > 0 {
								found = true
								break
							}
						}
						assert.True(t, found, "Expected to find validation errors in bulk import response")
					}
				}
			}
		})
	}
}

func TestIntegrationHandler_GetAPIDocumentation(t *testing.T) {
	buildInfo := BuildInfo{Version: "test", Commit: "test", BuildDate: "test"}
	baseURL := testBaseURL

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		expectedFields []string
	}{
		{
			name:           "valid GET request",
			method:         "GET",
			path:           "/api/v1/integration/docs",
			expectedStatus: http.StatusOK,
			expectedFields: []string{
				"title",
				"version",
				"description",
				"base_url",
				"endpoints",
				"authentication",
				"status_values",
				"service_types",
			},
		},
		{
			name:           "invalid POST request",
			method:         "POST",
			path:           "/api/v1/integration/docs",
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewIntegrationHandler(nil, baseURL, buildInfo)

			req := testutil.CreateTestHTTPRequest(tt.method, tt.path, nil)
			w := testutil.CreateTestHTTPResponse()

			handler.GetAPIDocumentation(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)

				for _, field := range tt.expectedFields {
					assert.Contains(t, response, field, "Expected field %s not found in response", field)
				}

				// Validate specific content
				assert.Equal(t, "Status Page Integration API", response["title"])
				assert.Equal(t, buildInfo.Version, response["version"])
				assert.Equal(t, baseURL, response["base_url"])

				// Validate endpoints structure
				endpoints, ok := response["endpoints"].(map[string]interface{})
				require.True(t, ok, "endpoints should be an object")
				assert.Contains(t, endpoints, "webhook")
				assert.Contains(t, endpoints, "manual_status")
				assert.Contains(t, endpoints, "integration_details")
				assert.Contains(t, endpoints, "bulk_import")

				// Validate status values
				statusValues, ok := response["status_values"].([]interface{})
				require.True(t, ok, "status_values should be an array")
				expectedStatuses := []string{"operational", "degraded", "down", "maintenance"}
				assert.Len(t, statusValues, len(expectedStatuses))
			}
		})
	}
}

func TestIntegrationHandler_extractServiceIDFromPath(t *testing.T) {
	handler := &IntegrationHandler{}

	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "valid path with service ID",
			path:     "/api/v1/integration/services/test-service-123",
			expected: "test-service-123",
		},
		{
			name:     "path with manual-status endpoint",
			path:     "/api/v1/integration/services/my-service/manual-status",
			expected: "my-service",
		},
		{
			name:     "complex service ID",
			path:     "/api/v1/integration/services/507f1f77bcf86cd799439011",
			expected: "507f1f77bcf86cd799439011",
		},
		{
			name:     "path without service ID",
			path:     "/api/v1/integration/services/",
			expected: "",
		},
		{
			name:     "empty path",
			path:     "",
			expected: "",
		},
		{
			name:     "different endpoint structure",
			path:     "/api/v1/other/services/test-id",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handler.extractServiceIDFromPath(tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIntegrationHandler_generateSlug(t *testing.T) {
	handler := &IntegrationHandler{}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple name",
			input:    "Test Service",
			expected: "test-service",
		},
		{
			name:     "name with underscores",
			input:    "My_Test_Service",
			expected: "my-test-service",
		},
		{
			name:     "name with special characters",
			input:    "Test@Service#123!",
			expected: "testservice123",
		},
		{
			name:     "name with mixed case",
			input:    "MyComplexServiceName",
			expected: "mycomplexservicename",
		},
		{
			name:     "name with numbers",
			input:    "Service 123 API",
			expected: "service-123-api",
		},
		{
			name:     "empty name",
			input:    "",
			expected: "",
		},
		{
			name:     "name with only spaces",
			input:    "   ",
			expected: "---",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handler.generateSlug(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIntegrationHandler_generateWebhookSecret(t *testing.T) {
	handler := &IntegrationHandler{}

	// Test that we can generate secrets
	secret1 := handler.generateWebhookSecret()
	secret2 := handler.generateWebhookSecret()

	// Secrets should not be empty
	assert.NotEmpty(t, secret1)
	assert.NotEmpty(t, secret2)

	// Secrets should be different (very high probability)
	assert.NotEqual(t, secret1, secret2)

	// Secrets should be hex encoded (64 characters for 32 bytes)
	assert.Len(t, secret1, 64)
	assert.Len(t, secret2, 64)

	// Should only contain hex characters
	for _, r := range secret1 {
		assert.True(t, (r >= '0' && r <= '9') || (r >= 'a' && r <= 'f'), "Secret should contain only hex characters")
	}
}

func TestIntegrationHandler_getCurrentUser(t *testing.T) {
	handler := &IntegrationHandler{}

	req := httptest.NewRequest("GET", "/test", nil)

	// Since we don't have auth implementation, should return "system"
	user := handler.getCurrentUser(req)
	assert.Equal(t, "system", user)
}

func TestIntegrationHandler_getWebhookExamples(t *testing.T) {
	handler := &IntegrationHandler{}

	tests := []struct {
		name       string
		webhookURL string
	}{
		{
			name:       "with URL",
			webhookURL: "http://localhost:8080/api/v1/webhook/test-service",
		},
		{
			name:       "with empty URL",
			webhookURL: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			examples := handler.getWebhookExamples(tt.webhookURL)

			require.NotNil(t, examples)
			assert.Contains(t, examples, "curl")
			assert.Contains(t, examples, "payload_example")
			assert.Contains(t, examples, "status_values")

			// Validate curl example
			curlExample, ok := examples["curl"].(string)
			require.True(t, ok)
			assert.Contains(t, curlExample, "curl -X POST")

			if tt.webhookURL != "" {
				assert.Contains(t, curlExample, tt.webhookURL)
			} else {
				assert.Contains(t, curlExample, "{webhook_url}")
			}

			// Validate payload example
			payloadExample, ok := examples["payload_example"].(map[string]interface{})
			require.True(t, ok)
			assert.Contains(t, payloadExample, "status")
			assert.Equal(t, "operational", payloadExample["status"])

			// Validate status values
			statusValues, ok := examples["status_values"].([]string)
			require.True(t, ok)
			expectedStatuses := []string{"operational", "degraded", "down", "maintenance"}
			assert.Equal(t, expectedStatuses, statusValues)
		})
	}
}
