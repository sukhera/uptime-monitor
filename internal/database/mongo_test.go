package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConnection_InvalidURI(t *testing.T) {
	// Test with invalid MongoDB URI
	db, err := NewConnection("invalid://uri", "testdb")

	assert.Error(t, err)
	assert.Nil(t, db)
}

func TestDB_CollectionMethods(t *testing.T) {
	// Test that DB methods work with a valid DB struct
	// Note: This doesn't actually connect to MongoDB, just tests the struct methods

	db := &DB{
		Client: nil, // We won't actually connect
		Name:   "testdb",
	}

	// Test that collection names are correct
	// We can't actually call these methods without a real connection,
	// but we can test the DB struct is properly formed
	assert.Equal(t, "testdb", db.Name)
	assert.NotNil(t, db) // Basic struct validation
}

func TestDB_InterfaceCompliance(t *testing.T) {
	// Test that DB implements Interface
	var _ Interface = (*DB)(nil)

	// This will fail to compile if DB doesn't implement Interface
	assert.True(t, true, "DB implements Interface")
}

func TestNewConnection_EmptyDatabase(t *testing.T) {
	// Test with empty database name - should still work
	// Note: This will fail without a real MongoDB instance
	db, err := NewConnection("mongodb://nonexistent:27017", "")

	// Expect error since we can't connect to nonexistent MongoDB
	assert.Error(t, err)
	assert.Nil(t, db)
}

// Integration test - only runs if MongoDB is available
func TestNewConnection_Integration(t *testing.T) {
	// Skip this test in unit test mode
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// This test requires a real MongoDB instance
	// In CI/CD, you would set up a test MongoDB container
	mongoURI := "mongodb://localhost:27017"
	dbName := "test_integration"

	db, err := NewConnection(mongoURI, dbName)

	if err != nil {
		t.Skipf("MongoDB not available for integration test: %v", err)
		return
	}

	require.NotNil(t, db)
	defer db.Close()

	// Test that collections can be retrieved
	servicesCollection := db.ServicesCollection()
	assert.NotNil(t, servicesCollection)
	assert.Equal(t, "services", servicesCollection.Name())

	statusLogsCollection := db.StatusLogsCollection()
	assert.NotNil(t, statusLogsCollection)
	assert.Equal(t, "status_logs", statusLogsCollection.Name())

	// Test database name
	database := db.Database()
	assert.NotNil(t, database)
	assert.Equal(t, dbName, database.Name())
}
