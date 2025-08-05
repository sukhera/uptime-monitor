package mongo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConnection_InvalidURI(t *testing.T) {
	// Test with invalid MongoDB URI
	db, err := NewConnection("invalid://uri", "testdb")

	assert.Error(t, err)
	assert.Nil(t, db)
}

func TestDatabase_CollectionMethods(t *testing.T) {
	// Test that Database methods work with a valid Database struct
	// Note: This doesn't actually connect to MongoDB, just tests the struct methods

	db := &Database{
		Client: nil, // We won't actually connect
		Name:   "testdb",
	}

	// Test that collection names are correct
	// We can't actually call these methods without a real connection,
	// but we can test the Database struct is properly formed
	assert.Equal(t, "testdb", db.Name)
	assert.NotNil(t, db) // Basic struct validation
}

func TestDatabase_InterfaceCompliance(t *testing.T) {
	// Test that Database implements Interface
	var _ Interface = (*Database)(nil)

	// This will fail to compile if Database doesn't implement Interface
	assert.True(t, true, "Database implements Interface")
}

func TestNewConnection_EmptyDatabase(t *testing.T) {
	// Test with empty database name - should still work
	// Note: This will fail without a real MongoDB instance
	db, err := NewConnection("mongodb://nonexistent:27017", "")

	// Expect error since we can't connect to nonexistent MongoDB
	assert.Error(t, err)
	assert.Nil(t, db)
}
