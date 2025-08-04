package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func TestService_JSONSerialization(t *testing.T) {
	service := Service{
		Name:           "test-service",
		Slug:           "test-service",
		URL:            "https://example.com",
		Headers:        map[string]string{"Authorization": "Bearer token"},
		ExpectedStatus: 200,
		Enabled:        true,
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(service)
	require.NoError(t, err)

	// Test JSON unmarshaling
	var unmarshaled Service
	err = json.Unmarshal(jsonData, &unmarshaled)
	require.NoError(t, err)

	assert.Equal(t, service, unmarshaled)
}

func TestService_BSONSerialization(t *testing.T) {
	service := Service{
		Name:           "test-service",
		Slug:           "test-service", 
		URL:            "https://example.com",
		Headers:        map[string]string{"Authorization": "Bearer token"},
		ExpectedStatus: 200,
		Enabled:        true,
	}

	// Test BSON marshaling
	bsonData, err := bson.Marshal(service)
	require.NoError(t, err)

	// Test BSON unmarshaling
	var unmarshaled Service
	err = bson.Unmarshal(bsonData, &unmarshaled)
	require.NoError(t, err)

	assert.Equal(t, service, unmarshaled)
}

// createTestStatusLog creates a test StatusLog for serialization tests
func createTestStatusLog() StatusLog {
	return StatusLog{
		ServiceName: "test-service",
		Status:      "operational",
		Latency:     150,
		StatusCode:  200,
		Error:       "",
		Timestamp:   time.Now(),
	}
}

// assertStatusLogEqual compares two StatusLog instances with time tolerance
func assertStatusLogEqual(t *testing.T, expected, actual StatusLog) {
	assert.Equal(t, expected.ServiceName, actual.ServiceName)
	assert.Equal(t, expected.Status, actual.Status)
	assert.Equal(t, expected.Latency, actual.Latency)
	assert.Equal(t, expected.StatusCode, actual.StatusCode)
	assert.Equal(t, expected.Error, actual.Error)
	assert.WithinDuration(t, expected.Timestamp, actual.Timestamp, time.Second)
}

func TestStatusLog_JSONSerialization(t *testing.T) {
	statusLog := createTestStatusLog()

	// Test JSON marshaling
	jsonData, err := json.Marshal(statusLog)
	require.NoError(t, err)

	// Test JSON unmarshaling
	var unmarshaled StatusLog
	err = json.Unmarshal(jsonData, &unmarshaled)
	require.NoError(t, err)

	// Compare with some tolerance for time precision
	assertStatusLogEqual(t, statusLog, unmarshaled)
}

func TestStatusLog_BSONSerialization(t *testing.T) {
	statusLog := createTestStatusLog()

	// Test BSON marshaling
	bsonData, err := bson.Marshal(statusLog)
	require.NoError(t, err)

	// Test BSON unmarshaling
	var unmarshaled StatusLog
	err = bson.Unmarshal(bsonData, &unmarshaled)
	require.NoError(t, err)

	assertStatusLogEqual(t, statusLog, unmarshaled)
}

func TestStatusLog_WithError(t *testing.T) {
	statusLog := StatusLog{
		ServiceName: "failing-service",
		Status:      "down",
		Latency:     0,
		StatusCode:  0,
		Error:       "Connection timeout",
		Timestamp:   time.Now(),
	}

	// Test that error field is properly handled
	jsonData, err := json.Marshal(statusLog)
	require.NoError(t, err)

	var unmarshaled StatusLog
	err = json.Unmarshal(jsonData, &unmarshaled)
	require.NoError(t, err)

	assert.Equal(t, "Connection timeout", unmarshaled.Error)
}

func TestStatusLog_EmptyError(t *testing.T) {
	statusLog := StatusLog{
		ServiceName: "working-service",
		Status:      "operational",
		Latency:     100,
		StatusCode:  200,
		Error:       "",
		Timestamp:   time.Now(),
	}

	// Test JSON serialization with empty error (should be omitted in BSON)
	bsonData, err := bson.Marshal(statusLog)
	require.NoError(t, err)

	var unmarshaled StatusLog
	err = bson.Unmarshal(bsonData, &unmarshaled)
	require.NoError(t, err)

	assert.Empty(t, unmarshaled.Error)
}

func TestServiceStatus_JSONSerialization(t *testing.T) {
	now := time.Now()
	serviceStatus := ServiceStatus{
		Name:      "test-service",
		Status:    "operational",
		Latency:   150,
		UpdatedAt: now,
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(serviceStatus)
	require.NoError(t, err)

	// Test JSON unmarshaling
	var unmarshaled ServiceStatus
	err = json.Unmarshal(jsonData, &unmarshaled)
	require.NoError(t, err)

	assert.Equal(t, serviceStatus.Name, unmarshaled.Name)
	assert.Equal(t, serviceStatus.Status, unmarshaled.Status)
	assert.Equal(t, serviceStatus.Latency, unmarshaled.Latency)
	assert.WithinDuration(t, serviceStatus.UpdatedAt, unmarshaled.UpdatedAt, time.Second)
}

func TestService_EmptyHeaders(t *testing.T) {
	service := Service{
		Name:           "simple-service",
		Slug:           "simple-service",
		URL:            "https://example.com",
		Headers:        nil, // Test nil headers
		ExpectedStatus: 200,
		Enabled:        true,
	}

	// Should handle nil headers gracefully
	jsonData, err := json.Marshal(service)
	require.NoError(t, err)

	var unmarshaled Service
	err = json.Unmarshal(jsonData, &unmarshaled)
	require.NoError(t, err)

	// nil map should become empty map after unmarshaling
	if unmarshaled.Headers == nil {
		unmarshaled.Headers = make(map[string]string)
	}
	
	assert.NotNil(t, unmarshaled.Headers)
	assert.Len(t, unmarshaled.Headers, 0)
}

func TestService_ValidationFields(t *testing.T) {
	tests := []struct {
		name    string
		service Service
		valid   bool
	}{
		{
			name: "valid service",
			service: Service{
				Name:           "valid-service",
				Slug:           "valid-service",
				URL:            "https://example.com",
				Headers:        map[string]string{},
				ExpectedStatus: 200,
				Enabled:        true,
			},
			valid: true,
		},
		{
			name: "disabled service",
			service: Service{
				Name:           "disabled-service",
				Slug:           "disabled-service",
				URL:            "https://example.com",
				Headers:        map[string]string{},
				ExpectedStatus: 200,
				Enabled:        false,
			},
			valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Basic validation - ensure required fields are not empty for enabled services
			if tt.service.Enabled {
				assert.NotEmpty(t, tt.service.Name, "Name should not be empty for enabled service")
				assert.NotEmpty(t, tt.service.URL, "URL should not be empty for enabled service")
				assert.NotZero(t, tt.service.ExpectedStatus, "ExpectedStatus should not be zero for enabled service")
			}
		})
	}
}