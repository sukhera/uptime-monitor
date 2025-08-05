package mongo

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Interface defines the interface for database operations
type Interface interface {
	ServicesCollection() *mongo.Collection
	StatusLogsCollection() *mongo.Collection
	Close() error
	Ping(ctx context.Context) error
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error)
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error)
}

type Database struct {
	Client *mongo.Client
	Name   string
}

func NewConnection(mongoURI, dbName string) (*Database, error) {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	log.Println("Connected to MongoDB successfully")

	return &Database{
		Client: client,
		Name:   dbName,
	}, nil
}

func (db *Database) Close() error {
	return db.Client.Disconnect(context.Background())
}

func (db *Database) Ping(ctx context.Context) error {
	return db.Client.Ping(ctx, nil)
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
