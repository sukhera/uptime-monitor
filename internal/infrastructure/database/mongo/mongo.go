package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/sukhera/uptime-monitor/internal/shared/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Interface defines the interface for database operations
type Interface interface {
	ServicesCollection() *mongo.Collection
	StatusLogsCollection() *mongo.Collection
	IncidentsCollection() *mongo.Collection
	MaintenancesCollection() *mongo.Collection
	Close() error
	Ping(ctx context.Context) error
	HealthCheck(ctx context.Context) error
	EnsureIndexes(ctx context.Context) error
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error)
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error)
}

type Database struct {
	Client  *mongo.Client
	Name    string
	timeout time.Duration
}

func NewConnection(mongoURI, dbName string) (*Database, error) {
	return NewConnectionWithTimeout(mongoURI, dbName, 10*time.Second)
}

func NewConnectionWithTimeout(mongoURI, dbName string, timeout time.Duration) (*Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping with timeout to ensure connection is ready
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	db := &Database{
		Client:  client,
		Name:    dbName,
		timeout: timeout,
	}

	// Health check on startup
	if err := db.HealthCheck(ctx); err != nil {
		return nil, fmt.Errorf("database health check failed: %w", err)
	}

	// Ensure indexes are created
	if err := db.EnsureIndexes(ctx); err != nil {
		return nil, fmt.Errorf("failed to create database indexes: %w", err)
	}

	log := logger.Get()
	log.Info(ctx, "Connected to MongoDB successfully", logger.Fields{
		"db_name":   dbName,
		"mongo_uri": mongoURI,
		"timeout":   timeout.String(),
	})

	return db, nil
}

func (db *Database) Close() error {
	return db.Client.Disconnect(context.Background())
}

func (db *Database) Ping(ctx context.Context) error {
	return db.Client.Ping(ctx, nil)
}

// HealthCheck performs a comprehensive health check of the database
func (db *Database) HealthCheck(ctx context.Context) error {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, db.timeout)
	defer cancel()

	// Test basic connectivity
	if err := db.Client.Ping(ctxWithTimeout, nil); err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}

	// Test database access
	database := db.Database()
	if database == nil {
		return fmt.Errorf("failed to get database instance")
	}

	// Test collections exist (create if not)
	collections := []string{"services", "status_logs", "incidents", "maintenances"}
	for _, collName := range collections {
		collection := database.Collection(collName)
		if collection == nil {
			return fmt.Errorf("failed to get collection: %s", collName)
		}

		// Test basic operation on each collection
		if _, err := collection.CountDocuments(ctxWithTimeout, bson.M{}); err != nil {
			return fmt.Errorf("failed to count documents in %s: %w", collName, err)
		}
	}

	return nil
}

// EnsureIndexes creates necessary indexes for optimal performance
func (db *Database) EnsureIndexes(ctx context.Context) error {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, db.timeout)
	defer cancel()

	log := logger.Get()

	// Services collection indexes
	servicesIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "slug", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("services_slug_unique"),
		},
		{
			Keys:    bson.D{{Key: "name", Value: 1}},
			Options: options.Index().SetName("services_name"),
		},
		{
			Keys:    bson.D{{Key: "enabled", Value: 1}},
			Options: options.Index().SetName("services_enabled"),
		},
	}

	if _, err := db.ServicesCollection().Indexes().CreateMany(ctxWithTimeout, servicesIndexes); err != nil {
		return fmt.Errorf("failed to create services indexes: %w", err)
	}

	// Status logs collection indexes with TTL for retention
	statusLogsIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "service_id", Value: 1}, {Key: "created_at", Value: -1}},
			Options: options.Index().SetName("status_logs_service_created"),
		},
		{
			Keys:    bson.D{{Key: "service_name", Value: 1}},
			Options: options.Index().SetName("status_logs_service_name"),
		},
		{
			Keys:    bson.D{{Key: "created_at", Value: 1}},
			Options: options.Index().SetExpireAfterSeconds(2592000).SetName("status_logs_ttl"), // 30 days TTL
		},
		{
			Keys:    bson.D{{Key: "status", Value: 1}},
			Options: options.Index().SetName("status_logs_status"),
		},
	}

	if _, err := db.StatusLogsCollection().Indexes().CreateMany(ctxWithTimeout, statusLogsIndexes); err != nil {
		return fmt.Errorf("failed to create status_logs indexes: %w", err)
	}

	// Incidents collection indexes
	incidentsIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "status", Value: 1}, {Key: "created_at", Value: -1}},
			Options: options.Index().SetName("incidents_status_created"),
		},
		{
			Keys:    bson.D{{Key: "created_at", Value: -1}},
			Options: options.Index().SetName("incidents_created_desc"),
		},
		{
			Keys:    bson.D{{Key: "affected_services", Value: 1}},
			Options: options.Index().SetName("incidents_affected_services"),
		},
	}

	if _, err := db.IncidentsCollection().Indexes().CreateMany(ctxWithTimeout, incidentsIndexes); err != nil {
		return fmt.Errorf("failed to create incidents indexes: %w", err)
	}

	// Maintenances collection indexes
	maintenancesIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "scheduled_start", Value: 1}},
			Options: options.Index().SetName("maintenances_scheduled_start"),
		},
		{
			Keys:    bson.D{{Key: "status", Value: 1}},
			Options: options.Index().SetName("maintenances_status"),
		},
		{
			Keys:    bson.D{{Key: "affected_services", Value: 1}},
			Options: options.Index().SetName("maintenances_affected_services"),
		},
	}

	if _, err := db.MaintenancesCollection().Indexes().CreateMany(ctxWithTimeout, maintenancesIndexes); err != nil {
		return fmt.Errorf("failed to create maintenances indexes: %w", err)
	}

	log.Info(ctx, "Database indexes created successfully", logger.Fields{
		"collections": []string{"services", "status_logs", "incidents", "maintenances"},
	})

	return nil
}

func (db *Database) Database() *mongo.Database {
	return db.Client.Database(db.Name)
}

func (db *Database) ServicesCollection() *mongo.Collection {
	return db.Database().Collection("services")
}

func (db *Database) StatusLogsCollection() *mongo.Collection {
	return db.Database().Collection("status_logs")
}

func (db *Database) IncidentsCollection() *mongo.Collection {
	return db.Database().Collection("incidents")
}

func (db *Database) MaintenancesCollection() *mongo.Collection {
	return db.Database().Collection("maintenances")
}

// Implement the database interface methods
func (db *Database) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	return db.ServicesCollection().Find(ctx, filter, opts...)
}

func (db *Database) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	return db.ServicesCollection().FindOne(ctx, filter, opts...)
}

func (db *Database) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return db.ServicesCollection().InsertOne(ctx, document, opts...)
}

func (db *Database) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return db.ServicesCollection().UpdateOne(ctx, filter, update, opts...)
}

func (db *Database) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return db.ServicesCollection().DeleteOne(ctx, filter, opts...)
}

func (db *Database) CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	return db.ServicesCollection().CountDocuments(ctx, filter, opts...)
}
