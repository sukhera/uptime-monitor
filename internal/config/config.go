package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	MongoURI      string
	DBName        string
	Port          string
	CheckInterval time.Duration
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

	checkInterval := 2 * time.Minute
	if intervalStr := os.Getenv("CHECK_INTERVAL_MINUTES"); intervalStr != "" {
		if minutes, err := strconv.Atoi(intervalStr); err == nil && minutes > 0 {
			checkInterval = time.Duration(minutes) * time.Minute
		}
	}

	return &Config{
		MongoURI:      mongoURI,
		DBName:        dbName,
		Port:          port,
		CheckInterval: checkInterval,
	}
}

func (c *Config) Validate() error {
	if c.CheckInterval <= 0 {
		return fmt.Errorf("check interval must be greater than 0")
	}
	if c.MongoURI == "" {
		return fmt.Errorf("mongo URI cannot be empty")
	}
	if c.DBName == "" {
		return fmt.Errorf("database name cannot be empty")
	}
	if c.Port == "" {
		return fmt.Errorf("port cannot be empty")
	}
	return nil
}
