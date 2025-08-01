package main

import (
	"context"
	"log"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/sukhera/uptime-monitor/internal/checker"
	"github.com/sukhera/uptime-monitor/internal/config"
	"github.com/sukhera/uptime-monitor/internal/database"
)

func main() {
	cfg := config.Load()

	db, err := database.NewConnection(cfg.MongoURI, cfg.DBName)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer db.Close()

	checkerService := checker.NewService(db)
	scheduler := gocron.NewScheduler(time.UTC)

	scheduler.Every(2).Minutes().Do(func() {
		ctx := context.Background()
		log.Println("Running health checks...")
		
		if err := checkerService.RunHealthChecks(ctx); err != nil {
			log.Printf("Error running health checks: %v", err)
		} else {
			log.Println("Health checks completed")
		}
	})

	log.Println("Starting status checker...")
	scheduler.StartBlocking()
}