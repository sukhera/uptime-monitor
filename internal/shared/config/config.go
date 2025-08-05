package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Logging  LoggingConfig
	Checker  CheckerConfig
}

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// DatabaseConfig holds database-specific configuration
type DatabaseConfig struct {
	URI     string
	Name    string
	Timeout time.Duration
}

// LoggingConfig holds logging-specific configuration
type LoggingConfig struct {
	Level string
	JSON  bool
}

// CheckerConfig holds checker-specific configuration
type CheckerConfig struct {
	Interval time.Duration
}

// Option is a function that configures a Config
type Option func(*Config)

// WithServerPort sets the server port
func WithServerPort(port string) Option {
	return func(c *Config) {
		c.Server.Port = port
	}
}

// WithServerTimeouts sets server timeouts
func WithServerTimeouts(read, write, idle time.Duration) Option {
	return func(c *Config) {
		c.Server.ReadTimeout = read
		c.Server.WriteTimeout = write
		c.Server.IdleTimeout = idle
	}
}

// WithDatabase sets database configuration
func WithDatabase(uri, name string, timeout time.Duration) Option {
	return func(c *Config) {
		c.Database.URI = uri
		c.Database.Name = name
		c.Database.Timeout = timeout
	}
}

// WithLogging sets logging configuration
func WithLogging(level string, json bool) Option {
	return func(c *Config) {
		c.Logging.Level = level
		c.Logging.JSON = json
	}
}

// WithCheckerInterval sets the checker interval
func WithCheckerInterval(interval time.Duration) Option {
	return func(c *Config) {
		c.Checker.Interval = interval
	}
}

// FromEnvironment loads configuration from environment variables
func FromEnvironment() Option {
	return func(c *Config) {
		c.Server.Port = getEnv("PORT", "8080")
		c.Server.ReadTimeout = getDurationEnv("READ_TIMEOUT", 15*time.Second)
		c.Server.WriteTimeout = getDurationEnv("WRITE_TIMEOUT", 15*time.Second)
		c.Server.IdleTimeout = getDurationEnv("IDLE_TIMEOUT", 60*time.Second)

		c.Database.URI = getEnv("MONGO_URI", "mongodb://localhost:27017")
		c.Database.Name = getEnv("DB_NAME", "statuspage")
		c.Database.Timeout = getDurationEnv("DB_TIMEOUT", 10*time.Second)

		c.Logging.Level = getEnv("LOG_LEVEL", "info")
		c.Logging.JSON = getBoolEnv("LOG_JSON", false)

		c.Checker.Interval = getDurationEnv("CHECK_INTERVAL", 2*time.Minute)
	}
}

// New creates a new Config with the given options
func New(options ...Option) *Config {
	config := &Config{
		Server: ServerConfig{
			Port:         "8080",
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		Database: DatabaseConfig{
			URI:     "mongodb://localhost:27017",
			Name:    "statuspage",
			Timeout: 10 * time.Second,
		},
		Logging: LoggingConfig{
			Level: "info",
			JSON:  false,
		},
		Checker: CheckerConfig{
			Interval: 2 * time.Minute,
		},
	}

	for _, option := range options {
		option(config)
	}

	return config
}

// Load loads configuration from environment variables (backward compatibility)
func Load() *Config {
	return New(FromEnvironment())
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Server.Port == "" {
		return fmt.Errorf("server port cannot be empty")
	}

	if c.Database.URI == "" {
		return fmt.Errorf("database URI cannot be empty")
	}

	if c.Database.Name == "" {
		return fmt.Errorf("database name cannot be empty")
	}

	if c.Checker.Interval <= 0 {
		return fmt.Errorf("checker interval must be positive")
	}

	return nil
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
