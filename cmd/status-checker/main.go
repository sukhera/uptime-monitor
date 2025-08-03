package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/sukhera/uptime-monitor/internal/checker"
	"github.com/sukhera/uptime-monitor/internal/config"
	"github.com/sukhera/uptime-monitor/internal/database"
)

func main() {
	var (
		intervalMinutes = flag.Int("interval", 0, "Check interval in minutes (overrides env var)")
		mongoURI       = flag.String("mongo-uri", "", "MongoDB connection URI (overrides env var)")
		dbName         = flag.String("db-name", "", "Database name (overrides env var)")
		verbose        = flag.Bool("verbose", false, "Enable verbose logging")
	)
	flag.Parse()

	cfg := config.Load()

	if *intervalMinutes > 0 {
		cfg.CheckInterval = time.Duration(*intervalMinutes) * time.Minute
	}
	if *mongoURI != "" {
		cfg.MongoURI = *mongoURI
	}
	if *dbName != "" {
		cfg.DBName = *dbName
	}

	if err := cfg.Validate(); err != nil {
		log.Fatalf("[ERROR] Invalid configuration: %v", err)
	}

	if *verbose {
		log.Printf("[INFO] Configuration: interval=%v, db=%s", cfg.CheckInterval, cfg.DBName)
	}
	log.Printf("[INFO] Starting status checker with %v interval", cfg.CheckInterval)

	db, err := database.NewConnection(cfg.MongoURI, cfg.DBName)
	if err != nil {
		log.Fatalf("[ERROR] Failed to connect to MongoDB: %v", err)
	}
	defer db.Close()

	checkerService := checker.NewService(db)
	scheduler := gocron.NewScheduler(time.UTC)

	scheduler.Every(cfg.CheckInterval).Do(func() {
		ctx := context.Background()
		log.Println("[INFO] Running health checks...")
		
		if err := checkerService.RunHealthChecks(ctx); err != nil {
			log.Printf("[ERROR] Error running health checks: %v", err)
		}
	})

	log.Println("[INFO] Status checker started successfully")
	scheduler.StartBlocking()
}