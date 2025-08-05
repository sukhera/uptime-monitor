package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig_New(t *testing.T) {
	tests := []struct {
		name     string
		options  []Option
		expected *Config
	}{
		{
			name:    "default configuration",
			options: []Option{},
			expected: &Config{
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
			},
		},
		{
			name: "custom server configuration",
			options: []Option{
				WithServerPort("9090"),
				WithServerTimeouts(30*time.Second, 30*time.Second, 120*time.Second),
			},
			expected: &Config{
				Server: ServerConfig{
					Port:         "9090",
					ReadTimeout:  30 * time.Second,
					WriteTimeout: 30 * time.Second,
					IdleTimeout:  120 * time.Second,
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
			},
		},
		{
			name: "custom database configuration",
			options: []Option{
				WithDatabase("mongodb://custom:27017", "custom_db", 15*time.Second),
			},
			expected: &Config{
				Server: ServerConfig{
					Port:         "8080",
					ReadTimeout:  15 * time.Second,
					WriteTimeout: 15 * time.Second,
					IdleTimeout:  60 * time.Second,
				},
				Database: DatabaseConfig{
					URI:     "mongodb://custom:27017",
					Name:    "custom_db",
					Timeout: 15 * time.Second,
				},
				Logging: LoggingConfig{
					Level: "info",
					JSON:  false,
				},
				Checker: CheckerConfig{
					Interval: 2 * time.Minute,
				},
			},
		},
		{
			name: "custom logging configuration",
			options: []Option{
				WithLogging("debug", true),
			},
			expected: &Config{
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
					Level: "debug",
					JSON:  true,
				},
				Checker: CheckerConfig{
					Interval: 2 * time.Minute,
				},
			},
		},
		{
			name: "custom checker configuration",
			options: []Option{
				WithCheckerInterval(5 * time.Minute),
			},
			expected: &Config{
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
					Interval: 5 * time.Minute,
				},
			},
		},
		{
			name: "multiple options",
			options: []Option{
				WithServerPort("9090"),
				WithDatabase("mongodb://custom:27017", "custom_db", 15*time.Second),
				WithLogging("debug", true),
				WithCheckerInterval(5 * time.Minute),
			},
			expected: &Config{
				Server: ServerConfig{
					Port:         "9090",
					ReadTimeout:  15 * time.Second,
					WriteTimeout: 15 * time.Second,
					IdleTimeout:  60 * time.Second,
				},
				Database: DatabaseConfig{
					URI:     "mongodb://custom:27017",
					Name:    "custom_db",
					Timeout: 15 * time.Second,
				},
				Logging: LoggingConfig{
					Level: "debug",
					JSON:  true,
				},
				Checker: CheckerConfig{
					Interval: 5 * time.Minute,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := New(tt.options...)
			assert.Equal(t, tt.expected, config)
		})
	}
}

func TestConfig_FromEnvironment(t *testing.T) {
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
			},
		},
		{
			name: "custom values from env vars",
			envVars: map[string]string{
				"MONGO_URI":      "mongodb://custom:27017",
				"DB_NAME":        "custom_db",
				"PORT":           "9090",
				"CHECK_INTERVAL": "5m",
				"LOG_LEVEL":      "debug",
				"LOG_JSON":       "true",
				"READ_TIMEOUT":   "30s",
				"WRITE_TIMEOUT":  "30s",
				"IDLE_TIMEOUT":   "120s",
				"DB_TIMEOUT":     "15s",
			},
			expected: &Config{
				Server: ServerConfig{
					Port:         "9090",
					ReadTimeout:  30 * time.Second,
					WriteTimeout: 30 * time.Second,
					IdleTimeout:  120 * time.Second,
				},
				Database: DatabaseConfig{
					URI:     "mongodb://custom:27017",
					Name:    "custom_db",
					Timeout: 15 * time.Second,
				},
				Logging: LoggingConfig{
					Level: "debug",
					JSON:  true,
				},
				Checker: CheckerConfig{
					Interval: 5 * time.Minute,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.resetEnv {
				os.Clearenv()
			}

			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}

			config := New(FromEnvironment())
			assert.Equal(t, tt.expected, config)
		})
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid configuration",
			config: &Config{
				Server: ServerConfig{Port: "8080"},
				Database: DatabaseConfig{
					URI:  "mongodb://localhost:27017",
					Name: "statuspage",
				},
				Checker: CheckerConfig{Interval: 2 * time.Minute},
			},
			wantErr: false,
		},
		{
			name: "empty server port",
			config: &Config{
				Server: ServerConfig{Port: ""},
				Database: DatabaseConfig{
					URI:  "mongodb://localhost:27017",
					Name: "statuspage",
				},
				Checker: CheckerConfig{Interval: 2 * time.Minute},
			},
			wantErr: true,
		},
		{
			name: "empty database URI",
			config: &Config{
				Server: ServerConfig{Port: "8080"},
				Database: DatabaseConfig{
					URI:  "",
					Name: "statuspage",
				},
				Checker: CheckerConfig{Interval: 2 * time.Minute},
			},
			wantErr: true,
		},
		{
			name: "empty database name",
			config: &Config{
				Server: ServerConfig{Port: "8080"},
				Database: DatabaseConfig{
					URI:  "mongodb://localhost:27017",
					Name: "",
				},
				Checker: CheckerConfig{Interval: 2 * time.Minute},
			},
			wantErr: true,
		},
		{
			name: "invalid checker interval",
			config: &Config{
				Server: ServerConfig{Port: "8080"},
				Database: DatabaseConfig{
					URI:  "mongodb://localhost:27017",
					Name: "statuspage",
				},
				Checker: CheckerConfig{Interval: 0},
			},
			wantErr: true,
		},
		{
			name: "negative checker interval",
			config: &Config{
				Server: ServerConfig{Port: "8080"},
				Database: DatabaseConfig{
					URI:  "mongodb://localhost:27017",
					Name: "statuspage",
				},
				Checker: CheckerConfig{Interval: -1 * time.Minute},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestConfig_Load_BackwardCompatibility(t *testing.T) {
	// Test that Load() still works for backward compatibility
	config := Load()
	require.NotNil(t, config)
	assert.NotEmpty(t, config.Server.Port)
	assert.NotEmpty(t, config.Database.URI)
	assert.NotEmpty(t, config.Database.Name)
	assert.True(t, config.Checker.Interval > 0)
}
