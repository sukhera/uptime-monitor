package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Interface defines the contract for database operations
type Interface interface {
	// Connection management
	Close() error
	Ping(ctx context.Context) error

	// Collection access
	ServicesCollection() *mongo.Collection
	StatusLogsCollection() *mongo.Collection

	// Database operations
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error)
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error)
}

// Repository defines the contract for domain-specific database operations
type Repository interface {
	// Service operations
	FindServices(ctx context.Context, filter interface{}) ([]interface{}, error)
	FindServiceByID(ctx context.Context, id string) (interface{}, error)
	CreateService(ctx context.Context, service interface{}) error
	UpdateService(ctx context.Context, id string, service interface{}) error
	DeleteService(ctx context.Context, id string) error

	// Status log operations
	FindStatusLogs(ctx context.Context, filter interface{}) ([]interface{}, error)
	CreateStatusLog(ctx context.Context, statusLog interface{}) error
	GetLatestStatusLogs(ctx context.Context, limit int64) ([]interface{}, error)
}
