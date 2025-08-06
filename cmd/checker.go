package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/sukhera/uptime-monitor/internal/checker"
	"github.com/sukhera/uptime-monitor/internal/infrastructure/database/mongo"
	"github.com/sukhera/uptime-monitor/internal/shared/config"
	"github.com/sukhera/uptime-monitor/internal/shared/logger"
)

var checkerCmd = &cobra.Command{
	Use:   "checker",
	Short: "Start the health monitoring service",
	Long: `Start the health monitoring service that:

- Monitors configured services for uptime
- Performs HTTP health checks
- Stores results in database
- Supports configurable intervals
- Handles retries and timeouts

Example:
  status-page checker --interval 30s --db-url mongodb://localhost:27017`,
	Run: runChecker,
}

var (
	checkInterval string
	checkerDBURL  string
	checkerDBName string
)

func init() {
	rootCmd.AddCommand(checkerCmd)

	// Checker-specific flags
	checkerCmd.Flags().StringVarP(&checkInterval, "interval", "i", "30s", "Health check interval")
	checkerCmd.Flags().StringVar(&checkerDBURL, "db-url", "mongodb://localhost:27017", "MongoDB connection URL")
	checkerCmd.Flags().StringVar(&checkerDBName, "db-name", "status_page", "MongoDB database name")

	// Bind flags to viper
	ctx := context.Background()
	log := logger.Get()
	
	if err := viper.BindPFlag("checker.interval", checkerCmd.Flags().Lookup("interval")); err != nil {
		log.Fatal(ctx, "Failed to bind checker.interval flag", err, nil)
	}
	if err := viper.BindPFlag("database.url", checkerCmd.Flags().Lookup("db-url")); err != nil {
		log.Fatal(ctx, "Failed to bind database.url flag", err, nil)
	}
	if err := viper.BindPFlag("database.name", checkerCmd.Flags().Lookup("db-name")); err != nil {
		log.Fatal(ctx, "Failed to bind database.name flag", err, nil)
	}
}

func runChecker(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	log := logger.Get()
	
	// Load configuration with proper precedence (flags > env > config file)
	cfg := config.LoadFromViper()
	
	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatal(ctx, "Invalid configuration", err, logger.Fields{})
	}

	// Initialize database
	db, err := mongo.NewConnection(cfg.Database.URI, cfg.Database.Name)
	if err != nil {
		log.Fatal(ctx, "Failed to connect to database", err, logger.Fields{"db_url": cfg.Database.URI, "db_name": cfg.Database.Name})
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Error(ctx, "Error closing database connection", err, nil)
		}
	}()

	// Initialize checker service
	service := checker.NewService(db)

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Info(ctx, "Shutting down checker", nil)
		cancel()
	}()

	// Start health checking loop
	log.Info(ctx, "Starting health checker", logger.Fields{
		"interval": cfg.Checker.Interval.String(),
		"db_url": cfg.Database.URI,
		"db_name": cfg.Database.Name,
	})

	ticker := time.NewTicker(cfg.Checker.Interval)
	defer ticker.Stop()

	// Run initial check
	if err := service.RunHealthChecks(ctx); err != nil {
		log.Error(ctx, "Initial health check failed", err, nil)
	}

	// Continue checking at intervals
	for {
		select {
		case <-ticker.C:
			if err := service.RunHealthChecks(ctx); err != nil {
				log.Error(ctx, "Health check failed", err, nil)
			}
		case <-ctx.Done():
			log.Info(ctx, "Checker stopped", nil)
			return
		}
	}
}
