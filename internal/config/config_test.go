package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig_Load(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		expected *Config
		resetEnv bool
	}{
		{
			name:     "default values when no env vars set",
			envVars:  map[string]string{},
			resetEnv: true,
			expected: &Config{
				MongoURI:      "mongodb://localhost:27017",
				DBName:        "statuspage",
				Port:          "8080",
				CheckInterval: 2 * time.Minute,
			},
		},
		{
			name: "custom values from env vars",
			envVars: map[string]string{
				"MONGO_URI":              "mongodb://custom:27017",
				"DB_NAME":                "custom_db",
				"PORT":                   "9090",
				"CHECK_INTERVAL_MINUTES": "5",
			},
			expected: &Config{
				MongoURI:      "mongodb://custom:27017",
				DBName:        "custom_db",
				Port:          "9090",
				CheckInterval: 5 * time.Minute,
			},
		},
		{
			name: "invalid interval defaults to 2 minutes",
			envVars: map[string]string{
				"CHECK_INTERVAL_MINUTES": "invalid",
			},
			resetEnv: true,
			expected: &Config{
				MongoURI:      "mongodb://localhost:27017",
				DBName:        "statuspage",
				Port:          "8080",
				CheckInterval: 2 * time.Minute,
			},
		},
		{
			name: "zero interval defaults to 2 minutes",
			envVars: map[string]string{
				"CHECK_INTERVAL_MINUTES": "0",
			},
			resetEnv: true,
			expected: &Config{
				MongoURI:      "mongodb://localhost:27017",
				DBName:        "statuspage",
				Port:          "8080",
				CheckInterval: 2 * time.Minute,
			},
		},
		{
			name: "negative interval defaults to 2 minutes",
			envVars: map[string]string{
				"CHECK_INTERVAL_MINUTES": "-5",
			},
			resetEnv: true,
			expected: &Config{
				MongoURI:      "mongodb://localhost:27017",
				DBName:        "statuspage",
				Port:          "8080",
				CheckInterval: 2 * time.Minute,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Store original env vars
			originalVars := make(map[string]string)
			envKeys := []string{"MONGO_URI", "DB_NAME", "PORT", "CHECK_INTERVAL_MINUTES"}

			for _, key := range envKeys {
				originalVars[key] = os.Getenv(key)
			}

			// Clean env if needed
			if tt.resetEnv {
				for _, key := range envKeys {
					os.Unsetenv(key)
				}
			}

			// Set test env vars
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			// Run test
			config := Load()

			// Assertions
			assert.Equal(t, tt.expected.MongoURI, config.MongoURI)
			assert.Equal(t, tt.expected.DBName, config.DBName)
			assert.Equal(t, tt.expected.Port, config.Port)
			assert.Equal(t, tt.expected.CheckInterval, config.CheckInterval)

			// Restore original env vars
			for key, value := range originalVars {
				if value == "" {
					os.Unsetenv(key)
				} else {
					os.Setenv(key, value)
				}
			}
		})
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name      string
		config    *Config
		expectErr bool
		errMsg    string
	}{
		{
			name: "valid config",
			config: &Config{
				MongoURI:      "mongodb://localhost:27017",
				DBName:        "testdb",
				Port:          "8080",
				CheckInterval: 2 * time.Minute,
			},
			expectErr: false,
		},
		{
			name: "zero check interval",
			config: &Config{
				MongoURI:      "mongodb://localhost:27017",
				DBName:        "testdb",
				Port:          "8080",
				CheckInterval: 0,
			},
			expectErr: true,
			errMsg:    "check interval must be greater than 0",
		},
		{
			name: "negative check interval",
			config: &Config{
				MongoURI:      "mongodb://localhost:27017",
				DBName:        "testdb",
				Port:          "8080",
				CheckInterval: -1 * time.Second,
			},
			expectErr: true,
			errMsg:    "check interval must be greater than 0",
		},
		{
			name: "empty mongo URI",
			config: &Config{
				MongoURI:      "",
				DBName:        "testdb",
				Port:          "8080",
				CheckInterval: 2 * time.Minute,
			},
			expectErr: true,
			errMsg:    "mongo URI cannot be empty",
		},
		{
			name: "empty database name",
			config: &Config{
				MongoURI:      "mongodb://localhost:27017",
				DBName:        "",
				Port:          "8080",
				CheckInterval: 2 * time.Minute,
			},
			expectErr: true,
			errMsg:    "database name cannot be empty",
		},
		{
			name: "empty port",
			config: &Config{
				MongoURI:      "mongodb://localhost:27017",
				DBName:        "testdb",
				Port:          "",
				CheckInterval: 2 * time.Minute,
			},
			expectErr: true,
			errMsg:    "port cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()

			if tt.expectErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestLoad_Integration(t *testing.T) {
	// Test that Load() + Validate() works together
	config := Load()
	err := config.Validate()

	// Should not error with default values
	assert.NoError(t, err)

	// Should have sensible defaults
	assert.NotEmpty(t, config.MongoURI)
	assert.NotEmpty(t, config.DBName)
	assert.NotEmpty(t, config.Port)
	assert.Positive(t, config.CheckInterval)
}
