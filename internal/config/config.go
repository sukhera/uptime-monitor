package config

import "os"

type Config struct {
	MongoURI string
	DBName   string
	Port     string
}

func Load() *Config {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "statuspage"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		MongoURI: mongoURI,
		DBName:   dbName,
		Port:     port,
	}
}