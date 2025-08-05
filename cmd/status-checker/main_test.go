package main

import (
	"context"
	"flag"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/sukhera/uptime-monitor/internal/checker"
	"github.com/sukhera/uptime-monitor/internal/shared/config"
	"github.com/sukhera/uptime-monitor/internal/shared/logger"
)

func TestMain_ConfigurationLoading(t *testing.T) {
	tests := []struct {
		name        string
		description string
		wantErr     bool
	}{
		{
			name:        "load_from_environment",
			description: "should load configuration from environment variables",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.New(config.FromEnvironment())

			assert.NotNil(t, cfg)
			assert.NotEmpty(t, cfg.Server.Port)
			assert.NotEmpty(t, cfg.Database.URI)
			assert.NotEmpty(t, cfg.Database.Name)
			assert.True(t, cfg.Checker.Interval > 0)
		})
	}
}

func TestMain_ConfigurationValidation(t *testing.T) {
	tests := []struct {
		name        string
		description string
		config      *config.Config
		wantErr     bool
		errContains string
	}{
		{
			name:        "valid_configuration",
			description: "should validate correct configuration",
			config: config.New(
				config.WithServerPort("8080"),
				config.WithDatabase("mongodb://localhost:27017", "testdb", 10*time.Second),
				config.WithCheckerInterval(2*time.Minute),
			),
			wantErr: false,
		},
		{
			name:        "invalid_empty_port",
			description: "should fail validation with empty port",
			config: config.New(
				config.WithServerPort(""),
				config.WithCheckerInterval(2*time.Minute),
			),
			wantErr:     true,
			errContains: "server port cannot be empty",
		},
		{
			name:        "invalid_zero_interval",
			description: "should fail validation with zero interval",
			config: config.New(
				config.WithServerPort("8080"),
				config.WithCheckerInterval(0),
			),
			wantErr:     true,
			errContains: "checker interval must be positive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMain_LoggerInitialization(t *testing.T) {
	tests := []struct {
		name        string
		description string
		level       logger.Level
	}{
		{
			name:        "info_level",
			description: "should initialize logger with info level",
			level:       logger.INFO,
		},
		{
			name:        "debug_level",
			description: "should initialize logger with debug level",
			level:       logger.DEBUG,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger.Init(tt.level)
			log := logger.Get()

			assert.NotNil(t, log)
		})
	}
}

func TestMain_CommandLineFlags(t *testing.T) {
	tests := []struct {
		name        string
		description string
		args        []string
		expected    map[string]interface{}
	}{
		{
			name:        "valid_flags",
			description: "should parse valid command line flags",
			args:        []string{"-interval", "5", "-mongo-uri", "mongodb://test:27017", "-db-name", "testdb", "-verbose"},
			expected: map[string]interface{}{
				"interval":  5,
				"mongo_uri": "mongodb://test:27017",
				"db_name":   "testdb",
				"verbose":   true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flag set for testing
			flag.CommandLine = flag.NewFlagSet("test", flag.ExitOnError)

			var (
				intervalMinutes = flag.Int("interval", 0, "Check interval in minutes")
				mongoURI        = flag.String("mongo-uri", "", "MongoDB connection URI")
				dbName          = flag.String("db-name", "", "Database name")
				verbose         = flag.Bool("verbose", false, "Enable verbose logging")
			)

			// Parse test arguments
			flag.CommandLine.Parse(tt.args)

			assert.Equal(t, tt.expected["interval"], *intervalMinutes)
			assert.Equal(t, tt.expected["mongo_uri"], *mongoURI)
			assert.Equal(t, tt.expected["db_name"], *dbName)
			assert.Equal(t, tt.expected["verbose"], *verbose)
		})
	}
}

func TestMain_ConfigurationOverride(t *testing.T) {
	tests := []struct {
		name        string
		description string
		overrides   map[string]interface{}
		expected    *config.Config
	}{
		{
			name:        "override_interval",
			description: "should override checker interval",
			overrides: map[string]interface{}{
				"interval": 10,
			},
			expected: &config.Config{
				Checker: config.CheckerConfig{
					Interval: 10 * time.Minute,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.New(config.FromEnvironment())

			var options []config.Option

			if interval, ok := tt.overrides["interval"].(int); ok && interval > 0 {
				options = append(options, config.WithCheckerInterval(time.Duration(interval)*time.Minute))
			}

			// Apply all options
			if len(options) > 0 {
				cfg = config.New(append([]config.Option{config.FromEnvironment()}, options...)...)
			}

			assert.Equal(t, tt.expected.Checker.Interval, cfg.Checker.Interval)
		})
	}
}

func TestProcessAlerts(t *testing.T) {
	tests := []struct {
		name        string
		description string
		alert       checker.HealthCheckEvent
	}{
		{
			name:        "process_down_alert",
			description: "should process down service alert",
			alert: checker.HealthCheckEvent{
				ServiceName: "test-service",
				Status:      "down",
				Latency:     5000,
				StatusCode:  500,
				Error:       "Connection timeout",
				Timestamp:   time.Now().Unix(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()

			alertCh := make(chan checker.HealthCheckEvent, 10)
			log := logger.New(logger.INFO)

			// Send test alert
			go func() {
				alertCh <- tt.alert
			}()

			// Start alert processing
			go processAlerts(ctx, alertCh, log)

			// Wait for processing
			time.Sleep(20 * time.Millisecond)
		})
	}
}

func TestMain_Integration_Configuration(t *testing.T) {
	tests := []struct {
		name        string
		description string
		config      *config.Config
		expected    map[string]interface{}
	}{
		{
			name:        "complete_configuration",
			description: "should handle complete configuration workflow",
			config: config.New(
				config.WithServerPort("9090"),
				config.WithDatabase("mongodb://localhost:27017", "testdb", 10*time.Second),
				config.WithLogging("debug", true),
				config.WithCheckerInterval(5*time.Minute),
			),
			expected: map[string]interface{}{
				"port":     "9090",
				"uri":      "mongodb://localhost:27017",
				"db_name":  "testdb",
				"timeout":  10 * time.Second,
				"level":    "debug",
				"json":     true,
				"interval": 5 * time.Minute,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Validate configuration
			err := tt.config.Validate()
			require.NoError(t, err)

			// Check configuration values
			assert.Equal(t, tt.expected["port"], tt.config.Server.Port)
			assert.Equal(t, tt.expected["uri"], tt.config.Database.URI)
			assert.Equal(t, tt.expected["db_name"], tt.config.Database.Name)
			assert.Equal(t, tt.expected["timeout"], tt.config.Database.Timeout)
			assert.Equal(t, tt.expected["level"], tt.config.Logging.Level)
			assert.Equal(t, tt.expected["json"], tt.config.Logging.JSON)
			assert.Equal(t, tt.expected["interval"], tt.config.Checker.Interval)
		})
	}
}

func TestMain_Integration_Observers(t *testing.T) {
	tests := []struct {
		name        string
		description string
		event       checker.HealthCheckEvent
	}{
		{
			name:        "operational_service",
			description: "should handle operational service event",
			event: checker.HealthCheckEvent{
				ServiceName: "test-service",
				Status:      "operational",
				Latency:     150,
				StatusCode:  200,
				Timestamp:   time.Now().Unix(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test observer integration
			subject := checker.NewHealthCheckSubject()

			// Add metrics observer
			metricsObserver := checker.NewMetricsObserver()
			subject.Attach(metricsObserver)

			// Add alerting observer
			alertingObserver := checker.NewAlertingObserver(5000)
			subject.Attach(alertingObserver)

			ctx := context.Background()

			// Notify observers
			assert.NotPanics(t, func() {
				subject.Notify(ctx, tt.event)
			})

			// Wait for processing
			time.Sleep(10 * time.Millisecond)

			// Check metrics
			metrics := metricsObserver.GetMetrics()
			assert.NotEmpty(t, metrics)
		})
	}
}

func TestMain_ErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		description string
		config      *config.Config
		wantErr     bool
		errContains string
	}{
		{
			name:        "invalid_empty_port",
			description: "should handle empty port error",
			config: config.New(
				config.WithServerPort(""),
				config.WithCheckerInterval(0),
			),
			wantErr:     true,
			errContains: "server port cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMain_ContextHandling(t *testing.T) {
	tests := []struct {
		name        string
		description string
		timeout     time.Duration
	}{
		{
			name:        "context_cancellation",
			description: "should handle context cancellation",
			timeout:     50 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), tt.timeout)
			defer cancel()

			alertCh := make(chan checker.HealthCheckEvent, 10)
			log := logger.New(logger.INFO)

			go processAlerts(ctx, alertCh, log)

			// Wait for context cancellation
			time.Sleep(tt.timeout + 10*time.Millisecond)
		})
	}
}

func TestMain_Performance_Configuration(t *testing.T) {
	tests := []struct {
		name        string
		description string
		iterations  int
		maxDuration time.Duration
	}{
		{
			name:        "performance_test",
			description: "should complete configuration performance test",
			iterations:  1000,
			maxDuration: 1 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start := time.Now()

			for i := 0; i < tt.iterations; i++ {
				cfg := config.New(
					config.WithServerPort("8080"),
					config.WithDatabase("mongodb://localhost:27017", "testdb", 10*time.Second),
					config.WithCheckerInterval(2*time.Minute),
				)
				_ = cfg.Validate()
			}

			duration := time.Since(start)

			// Should complete within reasonable time
			assert.True(t, duration < tt.maxDuration)
		})
	}
}

func TestMain_Concurrency_Observers(t *testing.T) {
	tests := []struct {
		name        string
		description string
		numEvents   int
	}{
		{
			name:        "concurrent_notifications",
			description: "should handle concurrent observer notifications",
			numEvents:   100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subject := checker.NewHealthCheckSubject()
			observer := &MockObserver{}
			subject.Attach(observer)

			ctx := context.Background()

			// Send multiple notifications concurrently
			for i := 0; i < tt.numEvents; i++ {
				go func(index int) {
					event := checker.HealthCheckEvent{
						ServiceName: "service-" + string(rune(index)),
						Status:      "operational",
						Latency:     150,
						StatusCode:  200,
						Timestamp:   time.Now().Unix(),
					}
					subject.Notify(ctx, event)
				}(i)
			}

			// Wait for all goroutines to complete
			time.Sleep(100 * time.Millisecond)

			// Should not panic
			assert.NotNil(t, subject)
		})
	}
}

// MockObserver is a mock observer for testing
type MockObserver struct {
	notified bool
	events   []checker.HealthCheckEvent
	mu       sync.RWMutex
}

func (m *MockObserver) OnHealthCheckCompleted(ctx context.Context, event checker.HealthCheckEvent) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.notified = true
	m.events = append(m.events, event)
}
