package main

import (
	"context"
	"flag"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/sukhera/uptime-monitor/internal/checker"
	"github.com/sukhera/uptime-monitor/internal/container"
	"github.com/sukhera/uptime-monitor/internal/shared/config"
	"github.com/sukhera/uptime-monitor/internal/shared/logger"
)

func main() {
	var (
		intervalMinutes = flag.Int("interval", 0, "Check interval in minutes (overrides env var)")
		mongoURI        = flag.String("mongo-uri", "", "MongoDB connection URI (overrides env var)")
		dbName          = flag.String("db-name", "", "Database name (overrides env var)")
		verbose         = flag.Bool("verbose", false, "Enable verbose logging")
	)
	flag.Parse()

	// Initialize structured logging
	logLevel := logger.INFO
	if *verbose {
		logLevel = logger.DEBUG
	}
	logger.Init(logLevel)
	log := logger.Get()

	ctx := context.Background()

	// Load configuration with functional options
	cfg := config.New(config.FromEnvironment())

	// Apply command-line overrides using functional options
	var options []config.Option

	if *intervalMinutes > 0 {
		options = append(options, config.WithCheckerInterval(time.Duration(*intervalMinutes)*time.Minute))
	}
	if *mongoURI != "" {
		options = append(options, config.WithDatabase(*mongoURI, cfg.Database.Name, cfg.Database.Timeout))
	}
	if *dbName != "" {
		options = append(options, config.WithDatabase(cfg.Database.URI, *dbName, cfg.Database.Timeout))
	}

	// Apply all options
	if len(options) > 0 {
		cfg = config.New(append([]config.Option{config.FromEnvironment()}, options...)...)
	}

	if err := cfg.Validate(); err != nil {
		log.Fatal(ctx, "Invalid configuration", err, logger.Fields{})
	}

	if *verbose {
		log.Info(ctx, "Configuration loaded", logger.Fields{
			"interval": cfg.Checker.Interval.String(),
			"db_name":  cfg.Database.Name,
		})
	}

	// Initialize dependency injection container
	container, err := container.New(cfg)
	if err != nil {
		log.Fatal(ctx, "Failed to create container", err, logger.Fields{})
	}

	// Get services from container
	checkerService, err := container.GetCheckerService()
	if err != nil {
		log.Fatal(ctx, "Failed to get checker service", err, logger.Fields{})
	}

	// Setup observers for health check events
	subject := checker.NewHealthCheckSubject()

	// Add logging observer
	loggingObserver := checker.NewLoggingObserver(log)
	subject.Attach(loggingObserver)

	// Add metrics observer
	metricsObserver := checker.NewMetricsObserver()
	subject.Attach(metricsObserver)

	// Add alerting observer
	alertingObserver := checker.NewAlertingObserver(5000) // 5 second threshold
	subject.Attach(alertingObserver)

	// Start alert processing goroutine
	go processAlerts(ctx, alertingObserver.GetAlertChannel(), log)

	log.Info(ctx, "Starting status checker", logger.Fields{
		"interval": cfg.Checker.Interval.String(),
	})

	// Setup scheduler with enhanced health check function
	scheduler := gocron.NewScheduler(time.UTC)

	_, err = scheduler.Every(cfg.Checker.Interval).Do(func() {
		runHealthChecks(ctx, checkerService, subject, log)
	})
	if err != nil {
		log.Fatal(ctx, "Failed to schedule health checks", err, logger.Fields{})
	}

	log.Info(ctx, "Status checker started successfully", logger.Fields{})
	scheduler.StartBlocking()
}

// runHealthChecks runs health checks with enhanced logging and metrics
func runHealthChecks(ctx context.Context, service checker.ServiceInterface, subject *checker.HealthCheckSubject, log logger.Logger) {
	log.Info(ctx, "Running health checks", logger.Fields{})

	if err := service.RunHealthChecks(ctx); err != nil {
		log.Error(ctx, "Error running health checks", err, logger.Fields{})
	}
}

// processAlerts processes alerts from the alerting observer
func processAlerts(ctx context.Context, alertCh <-chan checker.HealthCheckEvent, log logger.Logger) {
	for {
		select {
		case alert := <-alertCh:
			log.Warn(ctx, "Service alert triggered", logger.Fields{
				"service_name": alert.ServiceName,
				"status":       alert.Status,
				"latency_ms":   alert.Latency,
				"status_code":  alert.StatusCode,
				"error":        alert.Error,
			})

			// Here you could integrate with external alerting systems
			// like PagerDuty, Slack, email, etc.

		case <-ctx.Done():
			return
		}
	}
}
