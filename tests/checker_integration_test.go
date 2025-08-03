package tests

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sukhera/uptime-monitor/internal/checker"
	"github.com/sukhera/uptime-monitor/internal/config"
	"github.com/sukhera/uptime-monitor/internal/database"
	"github.com/sukhera/uptime-monitor/internal/models"
)

func TestCheckerService_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Setup test database
	cfg := &config.Config{
		MongoURI: "mongodb://localhost:27017",
		DBName:   "statuspage_test_checker",
	}

	db, err := database.NewConnection(cfg.MongoURI, cfg.DBName)
	if err != nil {
		t.Skipf("MongoDB not available for integration test: %v", err)
		return
	}
	require.NoError(t, err)
	defer db.Close()

	// Clean up test data
	ctx := context.Background()
	_, err = db.ServicesCollection().DeleteMany(ctx, map[string]interface{}{})
	require.NoError(t, err)
	_, err = db.StatusLogsCollection().DeleteMany(ctx, map[string]interface{}{})
	require.NoError(t, err)

	// Create test server
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/healthy":
			w.WriteHeader(200)
			w.Write([]byte("OK"))
		case "/slow":
			time.Sleep(100 * time.Millisecond)
			w.WriteHeader(200)
			w.Write([]byte("Slow but OK"))
		case "/error":
			w.WriteHeader(500)
			w.Write([]byte("Internal Server Error"))
		default:
			w.WriteHeader(404)
			w.Write([]byte("Not Found"))
		}
	}))
	defer testServer.Close()

	// Insert test services
	testServices := []models.Service{
		{
			Name:           "healthy-service",
			Slug:           "healthy-service",
			URL:            testServer.URL + "/healthy",
			Headers:        map[string]string{"User-Agent": "test-agent"},
			ExpectedStatus: 200,
			Enabled:        true,
		},
		{
			Name:           "slow-service",
			Slug:           "slow-service",
			URL:            testServer.URL + "/slow",
			Headers:        map[string]string{},
			ExpectedStatus: 200,
			Enabled:        true,
		},
		{
			Name:           "error-service",
			Slug:           "error-service",
			URL:            testServer.URL + "/error",
			Headers:        map[string]string{},
			ExpectedStatus: 200,
			Enabled:        true,
		},
		{
			Name:           "disabled-service",
			Slug:           "disabled-service",
			URL:            testServer.URL + "/healthy",
			Headers:        map[string]string{},
			ExpectedStatus: 200,
			Enabled:        false, // This should not be checked
		},
	}

	for _, service := range testServices {
		_, err := db.ServicesCollection().InsertOne(ctx, service)
		require.NoError(t, err)
	}

	// Create checker service and run health checks
	checkerService := checker.NewService(db)
	err = checkerService.RunHealthChecks(ctx)
	require.NoError(t, err)

	// Verify results
	cursor, err := db.StatusLogsCollection().Find(ctx, map[string]interface{}{})
	require.NoError(t, err)
	defer cursor.Close(ctx)

	var statusLogs []models.StatusLog
	err = cursor.All(ctx, &statusLogs)
	require.NoError(t, err)

	// Should have 3 status logs (disabled service should not be checked)
	assert.Len(t, statusLogs, 3)

	// Group logs by service name for easier testing
	logsByService := make(map[string]models.StatusLog)
	for _, log := range statusLogs {
		logsByService[log.ServiceName] = log
	}

	// Verify healthy service
	healthyLog, exists := logsByService["healthy-service"]
	require.True(t, exists)
	assert.Equal(t, "operational", healthyLog.Status)
	assert.Equal(t, 200, healthyLog.StatusCode)
	assert.Empty(t, healthyLog.Error)
	assert.GreaterOrEqual(t, healthyLog.Latency, int64(0))

	// Verify slow service
	slowLog, exists := logsByService["slow-service"]
	require.True(t, exists)
	assert.Equal(t, "operational", slowLog.Status)
	assert.Equal(t, 200, slowLog.StatusCode)
	assert.Empty(t, slowLog.Error)
	assert.GreaterOrEqual(t, slowLog.Latency, int64(100)) // Should take at least 100ms

	// Verify error service
	errorLog, exists := logsByService["error-service"]
	require.True(t, exists)
	assert.Equal(t, "down", errorLog.Status)
	assert.Equal(t, 500, errorLog.StatusCode)
	assert.Contains(t, errorLog.Error, "Unexpected status code: 500")
	assert.GreaterOrEqual(t, errorLog.Latency, int64(0))

	// Verify disabled service was not checked
	_, exists = logsByService["disabled-service"]
	assert.False(t, exists)
}

func TestCheckerService_ConcurrentChecks(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Setup test database
	cfg := &config.Config{
		MongoURI: "mongodb://localhost:27017",
		DBName:   "statuspage_test_concurrent",
	}

	db, err := database.NewConnection(cfg.MongoURI, cfg.DBName)
	if err != nil {
		t.Skipf("MongoDB not available for integration test: %v", err)
		return
	}
	require.NoError(t, err)
	defer db.Close()

	// Clean up test data
	ctx := context.Background()
	_, err = db.ServicesCollection().DeleteMany(ctx, map[string]interface{}{})
	require.NoError(t, err)
	_, err = db.StatusLogsCollection().DeleteMany(ctx, map[string]interface{}{})
	require.NoError(t, err)

	// Create test server with artificial delay
	requestCount := 0
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		time.Sleep(50 * time.Millisecond) // Small delay to test concurrency
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	}))
	defer testServer.Close()

	// Insert multiple test services to verify concurrent execution
	numServices := 5
	for i := 0; i < numServices; i++ {
		service := models.Service{
			Name:           fmt.Sprintf("concurrent-service-%d", i),
			Slug:           fmt.Sprintf("concurrent-service-%d", i),
			URL:            testServer.URL,
			Headers:        map[string]string{},
			ExpectedStatus: 200,
			Enabled:        true,
		}
		_, err := db.ServicesCollection().InsertOne(ctx, service)
		require.NoError(t, err)
	}

	// Measure execution time
	start := time.Now()
	checkerService := checker.NewService(db)
	err = checkerService.RunHealthChecks(ctx)
	require.NoError(t, err)
	duration := time.Since(start)

	// Verify that concurrent execution is faster than sequential
	// Sequential would take at least numServices * 50ms = 250ms
	// Concurrent should be significantly faster (closer to 50ms + overhead)
	maxExpectedDuration := time.Duration(numServices) * 25 * time.Millisecond // Allow some overhead
	assert.Less(t, duration, maxExpectedDuration, 
		"Concurrent execution should be faster than sequential")

	// Verify all services were checked
	cursor, err := db.StatusLogsCollection().Find(ctx, map[string]interface{}{})
	require.NoError(t, err)
	defer cursor.Close(ctx)

	var statusLogs []models.StatusLog
	err = cursor.All(ctx, &statusLogs)
	require.NoError(t, err)

	assert.Len(t, statusLogs, numServices)
	assert.Equal(t, numServices, requestCount)
}

func TestCheckerService_ConfigIntegration(t *testing.T) {
	// Test that the checker service works with the config system
	cfg := config.Load()
	err := cfg.Validate()
	require.NoError(t, err)

	// This test verifies that the config validation works with realistic values
	assert.Positive(t, cfg.CheckInterval)
	assert.NotEmpty(t, cfg.MongoURI)
	assert.NotEmpty(t, cfg.DBName)
	assert.NotEmpty(t, cfg.Port)
}