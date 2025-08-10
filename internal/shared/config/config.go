package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/spf13/viper"
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
	Port              string
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
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
func WithServerTimeouts(read, readHeader, write, idle time.Duration) Option {
	return func(c *Config) {
		c.Server.ReadTimeout = read
		c.Server.ReadHeaderTimeout = readHeader
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
		c.Server.ReadHeaderTimeout = getDurationEnv("READ_HEADER_TIMEOUT", 5*time.Second)
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
			Port:              "8080",
			ReadTimeout:       15 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
			WriteTimeout:      15 * time.Second,
			IdleTimeout:       60 * time.Second,
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

// LoadFromViper loads configuration with proper precedence (flags > env > config file)
func LoadFromViper() *Config {
	// Set viper defaults first
	setViperDefaults()

	config := &Config{
		Server: ServerConfig{
			Port:              viper.GetString("server.port"),
			ReadTimeout:       viper.GetDuration("server.read_timeout"),
			ReadHeaderTimeout: viper.GetDuration("server.read_header_timeout"),
			WriteTimeout:      viper.GetDuration("server.write_timeout"),
			IdleTimeout:       viper.GetDuration("server.idle_timeout"),
		},
		Database: DatabaseConfig{
			URI:     viper.GetString("database.url"),
			Name:    viper.GetString("database.name"),
			Timeout: viper.GetDuration("database.timeout"),
		},
		Logging: LoggingConfig{
			Level: viper.GetString("logging.level"),
			JSON:  viper.GetBool("logging.json"),
		},
		Checker: CheckerConfig{
			Interval: viper.GetDuration("checker.interval"),
		},
	}

	return config
}

// setViperDefaults sets default values in viper
func setViperDefaults() {
	// Server defaults
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.read_timeout", "15s")
	viper.SetDefault("server.read_header_timeout", "5s")
	viper.SetDefault("server.write_timeout", "15s")
	viper.SetDefault("server.idle_timeout", "60s")

	// Database defaults
	viper.SetDefault("database.url", "mongodb://localhost:27017")
	viper.SetDefault("database.name", "statuspage")
	viper.SetDefault("database.timeout", "10s")

	// Logging defaults
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.json", false)

	// Checker defaults
	viper.SetDefault("checker.interval", "2m")

	// API defaults (for consistency with current flags)
	viper.SetDefault("api.port", "8080")

	// Web defaults
	viper.SetDefault("web.port", "3000")
	viper.SetDefault("web.api_url", "http://localhost:8080")
	viper.SetDefault("web.static_dir", "./web/react-status-page/dist")
}


// Validate validates the configuration
func (c *Config) Validate() error {
	// Server validation
	if c.Server.Port == "" {
		return fmt.Errorf("server port cannot be empty")
	}

	if c.Server.ReadTimeout <= 0 {
		return fmt.Errorf("server read timeout must be positive")
	}

	if c.Server.ReadHeaderTimeout <= 0 {
		return fmt.Errorf("server read header timeout must be positive")
	}

	if c.Server.WriteTimeout <= 0 {
		return fmt.Errorf("server write timeout must be positive")
	}

	if c.Server.IdleTimeout <= 0 {
		return fmt.Errorf("server idle timeout must be positive")
	}

	// Database validation
	if c.Database.URI == "" {
		return fmt.Errorf("database URI cannot be empty (required)")
	}

	if c.Database.Name == "" {
		return fmt.Errorf("database name cannot be empty (required)")
	}

	if c.Database.Timeout <= 0 {
		return fmt.Errorf("database timeout must be positive")
	}

	// Checker validation
	if c.Checker.Interval <= 0 {
		return fmt.Errorf("checker interval must be positive")
	}

	// Logging validation
	if c.Logging.Level == "" {
		return fmt.Errorf("logging level cannot be empty")
	}

	validLevels := []string{"debug", "info", "warn", "error", "fatal", "panic"}
	isValidLevel := false
	for _, level := range validLevels {
		if c.Logging.Level == level {
			isValidLevel = true
			break
		}
	}
	if !isValidLevel {
		return fmt.Errorf("invalid logging level: %s (must be one of: debug, info, warn, error, fatal, panic)", c.Logging.Level)
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
