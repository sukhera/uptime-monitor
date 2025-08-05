package container

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sukhera/uptime-monitor/internal/checker"
	"github.com/sukhera/uptime-monitor/internal/shared/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MockDatabase is a simple mock for testing
type MockDatabase struct{}

func (m *MockDatabase) Close() error                            { return nil }
func (m *MockDatabase) Ping(ctx context.Context) error          { return nil }
func (m *MockDatabase) ServicesCollection() *mongo.Collection   { return nil }
func (m *MockDatabase) StatusLogsCollection() *mongo.Collection { return nil }
func (m *MockDatabase) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	return nil, nil
}
func (m *MockDatabase) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	return nil
}
func (m *MockDatabase) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return nil, nil
}
func (m *MockDatabase) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return nil, nil
}
func (m *MockDatabase) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return nil, nil
}
func (m *MockDatabase) CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	return 0, nil
}

// MockServiceInterface is a simple mock for testing
type MockServiceInterface struct{}

func (m *MockServiceInterface) RunHealthChecks(ctx context.Context) error           { return nil }
func (m *MockServiceInterface) AddObserver(observer checker.HealthCheckObserver)    {}
func (m *MockServiceInterface) RemoveObserver(observer checker.HealthCheckObserver) {}

func TestContainer_New(t *testing.T) {
	cfg := config.New()
	container, err := New(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, container)
	assert.NotNil(t, container.config)
	assert.NotNil(t, container.logger)
}

func TestContainer_WithDatabase(t *testing.T) {
	cfg := config.New()
	mockDB := &MockDatabase{}
	container, err := New(cfg, WithDatabase(mockDB))
	assert.NoError(t, err)
	assert.NotNil(t, container)

	// Test getting database
	db, err := container.GetDatabase()
	assert.NoError(t, err)
	assert.Equal(t, mockDB, db)
}

func TestContainer_WithCheckerService(t *testing.T) {
	cfg := config.New()
	mockService := &MockServiceInterface{}
	container, err := New(cfg, WithCheckerService(mockService))
	assert.NoError(t, err)
	assert.NotNil(t, container)

	// Test getting checker service
	service, err := container.GetCheckerService()
	assert.NoError(t, err)
	assert.Equal(t, mockService, service)
}

func TestContainer_Shutdown(t *testing.T) {
	cfg := config.New()
	mockDB := &MockDatabase{}
	container, err := New(cfg, WithDatabase(mockDB))
	assert.NoError(t, err)

	// Should not panic
	assert.NotPanics(t, func() {
		err := container.Shutdown(context.Background())
		assert.NoError(t, err)
	})
}

func TestContainer_Get(t *testing.T) {
	cfg := config.New()
	container, err := New(cfg)
	assert.NoError(t, err)

	// Test getting config
	config := container.GetConfig()
	assert.NotNil(t, config)

	// Test getting logger
	logger := container.GetLogger()
	assert.NotNil(t, logger)

	// Test getting database (should return error initially)
	// Skip this test in CI as it tries to connect to MongoDB
	if testing.Short() {
		t.Skip("Skipping database connection test in short mode (CI)")
	}

	db, err := container.GetDatabase()
	assert.Error(t, err)
	assert.Nil(t, db)

	// Test getting checker service (should return error initially)
	service, err := container.GetCheckerService()
	assert.Error(t, err)
	assert.Nil(t, service)
}

func TestContainer_Integration(t *testing.T) {
	// Create container with all dependencies
	cfg := config.New()
	mockDB := &MockDatabase{}
	mockService := &MockServiceInterface{}

	container, err := New(cfg,
		WithDatabase(mockDB),
		WithCheckerService(mockService),
	)
	assert.NoError(t, err)

	// Verify all dependencies are set
	assert.Equal(t, cfg, container.GetConfig())
	assert.NotNil(t, container.GetLogger())

	db, err := container.GetDatabase()
	assert.NoError(t, err)
	assert.Equal(t, mockDB, db)

	service, err := container.GetCheckerService()
	assert.NoError(t, err)
	assert.Equal(t, mockService, service)

	// Test shutdown
	err = container.Shutdown(context.Background())
	assert.NoError(t, err)
}
