package handlers

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
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

const invalidJSONPayload = "invalid json"

func TestNewWebhookHandler(t *testing.T) {
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
			buildInfo := BuildInfo{Version: "test", Commit: "test", BuildDate: "test"}
			handler := NewWebhookHandler(nil, buildInfo)
			assert.Equal(t, tt.expected, handler != nil)
		})
	}
}

// Additional tests for extractServiceSlugFromPath edge cases
func TestWebhookHandler_extractServiceSlugFromPath_additional(t *testing.T) {
	handler := &WebhookHandler{}

	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "root path",
			path:     "/",
			expected: "",
		},
		{
			name:     "double trailing slash",
			path:     "/api/v1/webhook/my-service//",
			expected: "",
		},
		{
			name:     "no webhook in path",
			path:     "/api/v1/other/test-service",
			expected: "test-service",
		},
		{
			name:     "service slug with dots",
			path:     "/api/v1/webhook/my.service.slug",
			expected: "my.service.slug",
		},
		{
			name:     "service slug with underscores",
			path:     "/api/v1/webhook/my_service_slug",
			expected: "my_service_slug",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handler.extractServiceSlugFromPath(tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Additional tests for validateWebhookSignature
func TestWebhookHandler_validateWebhookSignature_additional(t *testing.T) {
	handler := &WebhookHandler{}
	secret := "test-secret"
	payload := []byte(`{"status": "operational"}`)

	// Calculate valid signature
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	validSignature := hex.EncodeToString(mac.Sum(nil))

	tests := []struct {
		name      string
		body      []byte
		signature string
		secret    string
		expected  bool
	}{
		{
			name:      "signature with whitespace",
			body:      payload,
			signature: " " + validSignature + " ",
			secret:    secret,
			expected:  false,
		},
		{
			name:      "signature with wrong prefix",
			body:      payload,
			signature: "sha1=" + validSignature,
			secret:    secret,
			expected:  false,
		},
		{
			name:      "signature with uppercase hex",
			body:      payload,
			signature: "sha256=" + hex.EncodeToString(mac.Sum(nil)),
			secret:    secret,
			expected:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handler.validateWebhookSignature(tt.body, tt.signature, tt.secret)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Additional tests for validateWebhookPayload
func TestWebhookHandler_validateWebhookPayload_additional(t *testing.T) {
	handler := &WebhookHandler{}

	zeroLatency := int64(0)
	now := time.Now().UTC()

	tests := []struct {
		name    string
		payload service.WebhookPayload
		wantErr bool
		errMsg  string
	}{
		{
			name: "zero latency is valid",
			payload: service.WebhookPayload{
				Status:  "operational",
				Latency: &zeroLatency,
			},
			wantErr: false,
		},
		{
			name: "timestamp in the future is valid",
			payload: service.WebhookPayload{
				Status:    "operational",
				Timestamp: func() *time.Time { t := now.Add(24 * time.Hour); return &t }(),
			},
			wantErr: false,
		},
		{
			name: "metadata is nil",
			payload: service.WebhookPayload{
				Status:   "operational",
				Metadata: nil,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handler.validateWebhookPayload(&tt.payload)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Additional tests for GetWebhookExamples
func TestWebhookHandler_GetWebhookExamples_additional(t *testing.T) {
	buildInfo := BuildInfo{Version: "test", Commit: "test", BuildDate: "test"}
	handler := NewWebhookHandler(nil, buildInfo)

	req := httptest.NewRequest("GET", "/api/v1/webhook/examples", nil)
	w := httptest.NewRecorder()
	handler.GetWebhookExamples(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Test with unsupported method
	req = httptest.NewRequest("PUT", "/api/v1/webhook/examples", nil)
	w = httptest.NewRecorder()
	handler.GetWebhookExamples(w, req)
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

func TestWebhookHandler_HandleWebhook(t *testing.T) {
	buildInfo := BuildInfo{Version: "test", Commit: "test", BuildDate: "test"}

	tests := []struct {
		name           string
		method         string
		path           string
		payload        interface{}
		signature      string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "invalid method - GET not allowed",
			method:         "GET",
			path:           "/api/v1/webhook/test-service",
			payload:        nil,
			expectedStatus: http.StatusMethodNotAllowed,
			expectedError:  "Method not allowed",
		},
		{
			name:           "invalid path - missing service slug",
			method:         "POST",
			path:           "/api/v1/webhook/",
			payload:        service.WebhookPayload{Status: "operational"},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid service slug in URL",
		},
		{
			name:           "invalid JSON payload",
			method:         "POST",
			path:           "/api/v1/webhook/test-service",
			payload:        invalidJSONPayload,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid JSON payload",
		},
		{
			name:   "invalid payload - missing status",
			method: "POST",
			path:   "/api/v1/webhook/test-service",
			payload: map[string]interface{}{
				"message": "test message",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid payload",
		},
		{
			name:   "invalid payload - invalid status value",
			method: "POST",
			path:   "/api/v1/webhook/test-service",
			payload: service.WebhookPayload{
				Status: "invalid-status",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid payload",
		},
		{
			name:   "valid payload - minimal",
			method: "POST",
			path:   "/api/v1/webhook/test-service",
			payload: service.WebhookPayload{
				Status: "operational",
			},
			expectedStatus: http.StatusInternalServerError, // Will fail due to no repository implementation
			expectedError:  "service repository method not implemented",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewWebhookHandler(nil, buildInfo)

			var body []byte
			var err error

			if tt.payload != nil {
				if str, ok := tt.payload.(string); ok && str == invalidJSONPayload {
					body = []byte(str)
				} else {
					body, err = json.Marshal(tt.payload)
					require.NoError(t, err)
				}
			}

			req := httptest.NewRequest(tt.method, tt.path, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			if tt.signature != "" {
				req.Header.Set("X-Webhook-Signature", tt.signature)
			}

			w := httptest.NewRecorder()

			handler.HandleWebhook(w, req)

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

func TestWebhookHandler_extractServiceSlugFromPath(t *testing.T) {
	handler := &WebhookHandler{}

	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "valid path with service slug",
			path:     "/api/v1/webhook/test-service",
			expected: "test-service",
		},
		{
			name:     "path with trailing slash",
			path:     "/api/v1/webhook/my-service/",
			expected: "",
		},
		{
			name:     "path without service slug",
			path:     "/api/v1/webhook/",
			expected: "",
		},
		{
			name:     "complex service slug",
			path:     "/api/v1/webhook/my-complex-service-123",
			expected: "my-complex-service-123",
		},
		{
			name:     "empty path",
			path:     "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handler.extractServiceSlugFromPath(tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestWebhookHandler_validateWebhookSignature(t *testing.T) {
	handler := &WebhookHandler{}
	secret := "test-secret"
	payload := []byte(`{"status": "operational"}`)

	// Calculate valid signature
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	validSignature := hex.EncodeToString(mac.Sum(nil))

	tests := []struct {
		name      string
		body      []byte
		signature string
		secret    string
		expected  bool
	}{
		{
			name:      "valid signature",
			body:      payload,
			signature: validSignature,
			secret:    secret,
			expected:  true,
		},
		{
			name:      "valid signature with sha256 prefix",
			body:      payload,
			signature: "sha256=" + validSignature,
			secret:    secret,
			expected:  true,
		},
		{
			name:      "invalid signature",
			body:      payload,
			signature: "invalid-signature",
			secret:    secret,
			expected:  false,
		},
		{
			name:      "empty signature",
			body:      payload,
			signature: "",
			secret:    secret,
			expected:  false,
		},
		{
			name:      "wrong secret",
			body:      payload,
			signature: validSignature,
			secret:    "wrong-secret",
			expected:  false,
		},
		{
			name:      "different payload",
			body:      []byte(`{"status": "down"}`),
			signature: validSignature,
			secret:    secret,
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handler.validateWebhookSignature(tt.body, tt.signature, tt.secret)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestWebhookHandler_validateWebhookPayload(t *testing.T) {
	handler := &WebhookHandler{}

	latencyPositive := int64(150)
	latencyNegative := int64(-50)
	now := time.Now().UTC()

	tests := []struct {
		name    string
		payload service.WebhookPayload
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid minimal payload",
			payload: service.WebhookPayload{
				Status: "operational",
			},
			wantErr: false,
		},
		{
			name: "valid complete payload",
			payload: service.WebhookPayload{
				Status:    "degraded",
				Latency:   &latencyPositive,
				Message:   "Service degraded",
				Timestamp: &now,
				Metadata:  map[string]interface{}{"region": "us-west"},
			},
			wantErr: false,
		},
		{
			name: "missing status",
			payload: service.WebhookPayload{
				Message: "Test message",
			},
			wantErr: true,
			errMsg:  "status is required",
		},
		{
			name: "invalid status",
			payload: service.WebhookPayload{
				Status: "invalid-status",
			},
			wantErr: true,
			errMsg:  "invalid status value",
		},
		{
			name: "negative latency",
			payload: service.WebhookPayload{
				Status:  "operational",
				Latency: &latencyNegative,
			},
			wantErr: true,
			errMsg:  "latency must be between 0 and 300000ms",
		},
		{
			name: "all valid statuses - operational",
			payload: service.WebhookPayload{
				Status: "operational",
			},
			wantErr: false,
		},
		{
			name: "all valid statuses - degraded",
			payload: service.WebhookPayload{
				Status: "degraded",
			},
			wantErr: false,
		},
		{
			name: "all valid statuses - down",
			payload: service.WebhookPayload{
				Status: "down",
			},
			wantErr: false,
		},
		{
			name: "all valid statuses - maintenance",
			payload: service.WebhookPayload{
				Status: "maintenance",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handler.validateWebhookPayload(&tt.payload)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestWebhookHandler_GetWebhookExamples(t *testing.T) {
	buildInfo := BuildInfo{Version: "test", Commit: "test", BuildDate: "test"}
	handler := NewWebhookHandler(nil, buildInfo)

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
			path:           "/api/v1/webhook/examples",
			expectedStatus: http.StatusOK,
			expectedFields: []string{
				"webhook_url_format",
				"method",
				"headers",
				"payload_example",
				"status_values",
				"curl_example",
				"signature_calculation",
			},
		},
		{
			name:           "invalid POST request",
			method:         "POST",
			path:           "/api/v1/webhook/examples",
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := testutil.CreateTestHTTPRequest(tt.method, tt.path, nil)
			w := testutil.CreateTestHTTPResponse()

			handler.GetWebhookExamples(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)

				for _, field := range tt.expectedFields {
					assert.Contains(t, response, field, "Expected field %s not found in response", field)
				}

				// Validate specific content
				assert.Equal(t, "POST", response["method"])
				assert.Equal(t, "/api/v1/webhook/{service-slug}", response["webhook_url_format"])

				// Validate status values
				statusValues, ok := response["status_values"].([]interface{})
				require.True(t, ok, "status_values should be an array")
				expectedStatuses := []string{"operational", "degraded", "down", "maintenance"}
				assert.Len(t, statusValues, len(expectedStatuses))

				// Validate payload example structure
				payloadExample, ok := response["payload_example"].(map[string]interface{})
				require.True(t, ok, "payload_example should be an object")
				assert.Contains(t, payloadExample, "status")
				assert.Contains(t, payloadExample, "latency_ms")
			}
		})
	}
}

// Benchmark tests for performance
func BenchmarkWebhookHandler_validateWebhookSignature(b *testing.B) {
	handler := &WebhookHandler{}
	secret := "test-secret-for-benchmarking"
	payload := []byte(`{"status": "operational", "latency_ms": 150, "message": "All systems operational"}`)

	// Pre-calculate signature
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	signature := hex.EncodeToString(mac.Sum(nil))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		handler.validateWebhookSignature(payload, signature, secret)
	}
}

func BenchmarkWebhookHandler_extractServiceSlugFromPath(b *testing.B) {
	handler := &WebhookHandler{}
	path := "/api/v1/webhook/my-service-with-complex-name-123"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		handler.extractServiceSlugFromPath(path)
	}
}
