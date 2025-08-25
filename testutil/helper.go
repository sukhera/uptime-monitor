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

// Marshall marshals a value to JSON and returns a reader
func Marshall(t *testing.T, v any) io.Reader {
	b, err := json.Marshal(v)
	require.NoError(t, err)
	return bytes.NewReader(b)
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

// CreateMockHTTPResponse creates a mock HTTP response for testing
func CreateMockHTTPResponse(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(bytes.NewBufferString(body)),
		Header:     make(http.Header),
	}
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
		ServiceType:    service.ServiceTypeHTTP,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}

// CreateTestWebhookService creates a test webhook service
func CreateTestWebhookService() *service.Service {
	return &service.Service{
		Name:          "Test Webhook Service",
		Slug:          "test-webhook-service",
		Enabled:       true,
		ServiceType:   service.ServiceTypeWebhook,
		WebhookURL:    "http://localhost:8080/api/v1/webhook/test-webhook-service",
		WebhookSecret: "test-webhook-secret",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

// CreateTestManualStatusOverride creates a test manual status override
func CreateTestManualStatusOverride() *service.ManualStatusOverride {
	now := time.Now().UTC()
	future := now.Add(1 * time.Hour)
	return &service.ManualStatusOverride{
		Status:    "maintenance",
		Reason:    "Scheduled maintenance window",
		SetBy:     "test-user@example.com",
		SetAt:     now,
		ExpiresAt: &future,
	}
}

// CreateTestWebhookPayload creates a test webhook payload
func CreateTestWebhookPayload() *service.WebhookPayload {
	now := time.Now().UTC()
	latency := int64(150)
	return &service.WebhookPayload{
		Status:    "operational",
		Latency:   &latency,
		Message:   "All systems operational",
		Timestamp: &now,
		Metadata: map[string]interface{}{
			"region":  "us-west-1",
			"version": "1.0.0",
		},
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

// CreateTestServiceStatus creates a test service status with default values
func CreateTestServiceStatus() *service.ServiceStatus {
	return &service.ServiceStatus{
		Name:      "Test Service",
		Status:    "operational",
		Latency:   100,
		UpdatedAt: time.Now(),
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

// AssertHTTPResponse asserts HTTP response status and optionally body
func AssertHTTPResponse(t *testing.T, rec *httptest.ResponseRecorder, expectedStatus int, expectedBody ...string) {
	assert.Equal(t, expectedStatus, rec.Code)

	if len(expectedBody) > 0 {
		assert.Contains(t, rec.Body.String(), expectedBody[0])
	}
}

// CreateTestCursor creates a mock cursor for testing
func CreateTestCursor(t *testing.T, documents []interface{}) *mongo.Cursor {
	// This is a simplified mock - in real tests you might want to use a more sophisticated approach
	return &mongo.Cursor{}
}

// CreateTestFindOptions creates test find options
func CreateTestFindOptions() *options.FindOptions {
	return options.Find()
}

// CreateTestInsertOneResult creates a test insert one result
func CreateTestInsertOneResult() *mongo.InsertOneResult {
	return &mongo.InsertOneResult{
		InsertedID: "test-id",
	}
}

// CreateTestUpdateResult creates a test update result
func CreateTestUpdateResult() *mongo.UpdateResult {
	return &mongo.UpdateResult{
		MatchedCount:  1,
		ModifiedCount: 1,
		UpsertedCount: 0,
	}
}

// CreateTestDeleteResult creates a test delete result
func CreateTestDeleteResult() *mongo.DeleteResult {
	return &mongo.DeleteResult{
		DeletedCount: 1,
	}
}
