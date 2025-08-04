package main

import (
	"log"
	"net/http"
	"time"

	"github.com/sukhera/uptime-monitor/internal/api/routes"
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

	router := routes.Setup(db)
	handler := routes.WithMiddleware(router)

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("API server starting on port %s", cfg.Port)
	log.Fatal(server.ListenAndServe())
}