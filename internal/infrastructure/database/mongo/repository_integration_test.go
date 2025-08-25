package mongo

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/sukhera/uptime-monitor/internal/domain/service"
	"github.com/sukhera/uptime-monitor/testutil"
)

// TestServiceRepository_Integration tests the new integration methods
// These tests focus on the new functionality added for webhook and integration features

func TestServiceRepository_FindByType(t *testing.T) {
	// Test with nil database - this tests method signature and validation logic
	repo := NewServiceRepository(nil)
	ctx := testutil.TestContext(t)

	tests := []struct {
		name        string
		serviceType service.ServiceType
		expectError bool
	}{
		{
			name:        "find HTTP services",
			serviceType: service.ServiceTypeHTTP,
			expectError: true, // Will error with nil DB
		},
		{
			name:        "find webhook services",
			serviceType: service.ServiceTypeWebhook,
			expectError: true, // Will error with nil DB
		},
		{
			name:        "find TCP services",
			serviceType: service.ServiceTypeTCP,
			expectError: true, // Will error with nil DB
		},
		{
			name:        "find DNS services",
			serviceType: service.ServiceTypeDNS,
			expectError: true, // Will error with nil DB
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			services, err := repo.FindByType(ctx, tt.serviceType)

			if tt.expectError {
				assert.Error(t, err) // Nil DB will cause panic/error
				assert.Nil(t, services)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, services)
			}
		})
	}
}

func TestServiceRepository_BulkCreate(t *testing.T) {
	repo := NewServiceRepository(nil)
	ctx := testutil.TestContext(t)

	now := time.Now().UTC()

	validServices := []*service.Service{
		{
			Name:        "Test Service 1",
			Slug:        "test-service-1",
			URL:         "https://example1.com",
			ServiceType: service.ServiceTypeHTTP,
			Enabled:     true,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			Name:        "Test Service 2",
			Slug:        "test-service-2",
			URL:         "https://example2.com",
			ServiceType: service.ServiceTypeHTTP,
			Enabled:     true,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}

	invalidServices := []*service.Service{
		{
			// Missing name
			Slug:        "invalid-service",
			URL:         "https://example.com",
			ServiceType: service.ServiceTypeHTTP,
			Enabled:     true,
		},
	}

	tests := []struct {
		name        string
		services    []*service.Service
		expectError bool
		errorMsg    string
	}{
		{
			name:        "empty services array",
			services:    []*service.Service{},
			expectError: true,
			errorMsg:    "no services provided for bulk create",
		},
		{
			name:        "nil services array",
			services:    nil,
			expectError: true,
			errorMsg:    "no services provided for bulk create",
		},
		{
			name:        "invalid services - validation error",
			services:    invalidServices,
			expectError: true,
			errorMsg:    "invalid service at index",
		},
		{
			name:        "valid services",
			services:    validServices,
			expectError: true, // Mock will return error for actual DB operation
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.BulkCreate(ctx, tt.services)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestServiceRepository_SetManualStatus(t *testing.T) {
	repo := NewServiceRepository(nil)
	ctx := testutil.TestContext(t)

	now := time.Now().UTC()
	future := now.Add(1 * time.Hour)

	validOverride := &service.ManualStatusOverride{
		Status:    "maintenance",
		Reason:    "Scheduled maintenance",
		SetBy:     "admin@example.com",
		SetAt:     now,
		ExpiresAt: &future,
	}

	tests := []struct {
		name        string
		serviceID   string
		override    *service.ManualStatusOverride
		expectError bool
		errorMsg    string
	}{
		{
			name:        "nil override",
			serviceID:   "test-service",
			override:    nil,
			expectError: true,
			errorMsg:    "manual status override cannot be nil",
		},
		{
			name:        "valid override with ObjectID",
			serviceID:   "507f1f77bcf86cd799439011", // Valid ObjectID
			override:    validOverride,
			expectError: true, // Mock will return error
		},
		{
			name:        "valid override with slug",
			serviceID:   "test-service",
			override:    validOverride,
			expectError: true, // Mock will return error
		},
		{
			name:        "empty service ID",
			serviceID:   "",
			override:    validOverride,
			expectError: true, // Mock will return error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.SetManualStatus(ctx, tt.serviceID, tt.override)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestServiceRepository_ClearManualStatus(t *testing.T) {
	repo := NewServiceRepository(nil)
	ctx := testutil.TestContext(t)

	tests := []struct {
		name        string
		serviceID   string
		expectError bool
	}{
		{
			name:        "clear with ObjectID",
			serviceID:   "507f1f77bcf86cd799439011", // Valid ObjectID
			expectError: true,                       // Mock will return error
		},
		{
			name:        "clear with slug",
			serviceID:   "test-service",
			expectError: true, // Mock will return error
		},
		{
			name:        "empty service ID",
			serviceID:   "",
			expectError: true, // Mock will return error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.ClearManualStatus(ctx, tt.serviceID)

			if tt.expectError {
				assert.Error(t, err) // Mock will return error
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestStatusLogRepository tests the new StatusLogRepository
func TestNewStatusLogRepository(t *testing.T) {
	repo := NewStatusLogRepository(nil)

	assert.NotNil(t, repo)
	assert.Nil(t, repo.db)
}

func TestStatusLogRepository_Create(t *testing.T) {
	repo := NewStatusLogRepository(nil)
	ctx := testutil.TestContext(t)

	now := time.Now().UTC()

	validLog := &service.StatusLog{
		ServiceName: "Test Service",
		Status:      "operational",
		Latency:     100,
		StatusCode:  200,
		Timestamp:   now,
	}

	tests := []struct {
		name        string
		log         *service.StatusLog
		expectError bool
		errorMsg    string
	}{
		{
			name:        "nil log",
			log:         nil,
			expectError: true,
			errorMsg:    "status log cannot be nil",
		},
		{
			name: "missing service name",
			log: &service.StatusLog{
				Status:    "operational",
				Latency:   100,
				Timestamp: now,
			},
			expectError: true,
			errorMsg:    "service name is required",
		},
		{
			name: "missing status",
			log: &service.StatusLog{
				ServiceName: "Test Service",
				Latency:     100,
				Timestamp:   now,
			},
			expectError: true,
			errorMsg:    "status is required",
		},
		{
			name: "missing timestamp - should be auto-set",
			log: &service.StatusLog{
				ServiceName: "Test Service",
				Status:      "operational",
				Latency:     100,
			},
			expectError: true, // Mock will return error for DB operation
		},
		{
			name:        "valid log",
			log:         validLog,
			expectError: true, // Mock will return error for DB operation
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(ctx, tt.log)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestStatusLogRepository_FindByServiceName(t *testing.T) {
	repo := NewStatusLogRepository(nil)
	ctx := testutil.TestContext(t)

	tests := []struct {
		name        string
		serviceName string
		limit       int
		expectError bool
		errorMsg    string
	}{
		{
			name:        "empty service name",
			serviceName: "",
			limit:       10,
			expectError: true,
			errorMsg:    "service name is required",
		},
		{
			name:        "valid service name",
			serviceName: "Test Service",
			limit:       10,
			expectError: true, // Mock will return error
		},
		{
			name:        "zero limit",
			serviceName: "Test Service",
			limit:       0,
			expectError: true, // Mock will return error
		},
		{
			name:        "large limit",
			serviceName: "Test Service",
			limit:       1000,
			expectError: true, // Mock will return error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logs, err := repo.FindByServiceName(ctx, tt.serviceName, tt.limit)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, logs)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, logs)
			}
		})
	}
}

func TestStatusLogRepository_FindRecent(t *testing.T) {
	repo := NewStatusLogRepository(nil)
	ctx := testutil.TestContext(t)

	tests := []struct {
		name        string
		limit       int
		expectError bool
	}{
		{
			name:        "find recent logs",
			limit:       10,
			expectError: true, // Mock will return error
		},
		{
			name:        "zero limit",
			limit:       0,
			expectError: true, // Mock will return error
		},
		{
			name:        "large limit",
			limit:       1000,
			expectError: true, // Mock will return error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logs, err := repo.FindRecent(ctx, tt.limit)

			if tt.expectError {
				assert.Error(t, err) // Mock will return error
				assert.Nil(t, logs)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, logs)
			}
		})
	}
}

func TestStatusLogRepository_DeleteOlderThan(t *testing.T) {
	repo := NewStatusLogRepository(nil)
	ctx := testutil.TestContext(t)

	now := time.Now().UTC()
	past := now.Add(-24 * time.Hour)
	future := now.Add(1 * time.Hour)

	tests := []struct {
		name        string
		cutoff      time.Time
		expectError bool
	}{
		{
			name:        "delete old logs",
			cutoff:      past,
			expectError: true, // Mock will return error
		},
		{
			name:        "delete with future cutoff",
			cutoff:      future,
			expectError: true, // Mock will return error
		},
		{
			name:        "delete with zero time",
			cutoff:      time.Time{},
			expectError: true, // Mock will return error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.DeleteOlderThan(ctx, tt.cutoff)

			if tt.expectError {
				assert.Error(t, err) // Mock will return error
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestStatusLogRepository_GetLatestByService(t *testing.T) {
	repo := NewStatusLogRepository(nil)
	ctx := testutil.TestContext(t)

	// Test the method (will fail with mock, but tests the signature)
	result, err := repo.GetLatestByService(ctx)

	// Mock will return error
	assert.Error(t, err)
	assert.Nil(t, result)
}

// Test helper functions for ObjectID handling
func TestObjectIDHandling(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		isValidOID bool
	}{
		{
			name:       "valid ObjectID",
			input:      "507f1f77bcf86cd799439011",
			isValidOID: true,
		},
		{
			name:       "invalid ObjectID - too short",
			input:      "507f1f77bcf86cd79943",
			isValidOID: false,
		},
		{
			name:       "invalid ObjectID - non-hex chars",
			input:      "507f1f77bcf86cd799439xyz",
			isValidOID: false,
		},
		{
			name:       "service slug",
			input:      "my-service-name",
			isValidOID: false,
		},
		{
			name:       "empty string",
			input:      "",
			isValidOID: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := primitive.ObjectIDFromHex(tt.input)

			if tt.isValidOID {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

// Benchmark tests for performance-critical operations
func BenchmarkServiceRepository_BulkCreate_Validation(b *testing.B) {
	repo := NewServiceRepository(nil)
	ctx := context.Background()

	now := time.Now().UTC()
	services := make([]*service.Service, 100)
	for i := 0; i < 100; i++ {
		services[i] = &service.Service{
			Name:        "Benchmark Service",
			Slug:        "benchmark-service",
			URL:         "https://example.com",
			ServiceType: service.ServiceTypeHTTP,
			Enabled:     true,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// This will fail at DB operation but will test validation performance
		_ = repo.BulkCreate(ctx, services)
	}
}

func BenchmarkStatusLogRepository_Create_Validation(b *testing.B) {
	repo := NewStatusLogRepository(nil)
	ctx := context.Background()

	log := &service.StatusLog{
		ServiceName: "Benchmark Service",
		Status:      "operational",
		Latency:     100,
		StatusCode:  200,
		Timestamp:   time.Now().UTC(),
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// This will fail at DB operation but will test validation performance
		_ = repo.Create(ctx, log)
	}
}
