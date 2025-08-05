package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/sukhera/uptime-monitor/internal/checker"
	"github.com/sukhera/uptime-monitor/internal/infrastructure/database/mongo"
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
	if err := viper.BindPFlag("checker.interval", checkerCmd.Flags().Lookup("interval")); err != nil {
		log.Fatalf("Failed to bind checker.interval flag: %v", err)
	}
	if err := viper.BindPFlag("database.url", checkerCmd.Flags().Lookup("db-url")); err != nil {
		log.Fatalf("Failed to bind database.url flag: %v", err)
	}
	if err := viper.BindPFlag("database.name", checkerCmd.Flags().Lookup("db-name")); err != nil {
		log.Fatalf("Failed to bind database.name flag: %v", err)
	}
}

func runChecker(cmd *cobra.Command, args []string) {
	// Parse interval
	checkInterval, err := time.ParseDuration(checkInterval)
	if err != nil {
		log.Fatalf("Invalid interval format: %v", err)
	}

	// Initialize database
	db, err := mongo.NewConnection(checkerDBURL, checkerDBName)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
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

		log.Println("Shutting down checker...")
		cancel()
	}()

	// Start health checking loop
	log.Printf("Starting health checker with %s interval", checkInterval)
	log.Printf("Connected to database: %s/%s", checkerDBURL, checkerDBName)

	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	// Run initial check
	if err := service.RunHealthChecks(ctx); err != nil {
		log.Printf("Initial health check failed: %v", err)
	}

	// Continue checking at intervals
	for {
		select {
		case <-ticker.C:
			if err := service.RunHealthChecks(ctx); err != nil {
				log.Printf("Health check failed: %v", err)
			}
		case <-ctx.Done():
			log.Println("Checker stopped")
			return
		}
	}
}
