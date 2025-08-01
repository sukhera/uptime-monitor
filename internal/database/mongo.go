package database

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	Client *mongo.Client
	Name   string
}

func NewConnection(mongoURI, dbName string) (*DB, error) {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	log.Println("Connected to MongoDB successfully")

	return &DB{
		Client: client,
		Name:   dbName,
	}, nil
}

func (db *DB) Close() error {
	return db.Client.Disconnect(context.Background())
}

func (db *DB) Database() *mongo.Database {
	return db.Client.Database(db.Name)
}

func (db *DB) ServicesCollection() *mongo.Collection {
	return db.Database().Collection("services")
}

func (db *DB) StatusLogsCollection() *mongo.Collection {
	return db.Database().Collection("status_logs")
}