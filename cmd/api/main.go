package main

import (
	"log"
	"net/http"

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

	log.Printf("API server starting on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, handler))
}