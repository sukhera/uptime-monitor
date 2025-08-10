package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/sukhera/uptime-monitor/internal/application/routes"
	"github.com/sukhera/uptime-monitor/internal/container"
	"github.com/sukhera/uptime-monitor/internal/shared/config"
	"github.com/sukhera/uptime-monitor/internal/shared/logger"
)

func main() {
	// Create context that cancels on interrupt signals
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Load configuration with proper precedence (flags > env > config file)
	cfg := config.LoadFromViper()

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	// Initialize logger
	log := logger.Get()
	log.Info(ctx, "Starting API server", logger.Fields{
		"port":     cfg.Server.Port,
		"database": cfg.Database.URI,
	})

	// Create dependency injection container using libnexus-style patterns
	container, err := container.New(cfg)
	if err != nil {
		log.Fatal(ctx, "Failed to create container", err, logger.Fields{})
	}

	// Get HTTP server from container
	httpServer, err := container.GetHTTPServer()
	if err != nil {
		log.Fatal(ctx, "Failed to get HTTP server", err, logger.Fields{})
	}

	// Log registered routes
	logRoutes(log)

	// Start server with context
	log.Info(ctx, "Starting HTTP server", logger.Fields{
		"address": httpServer.GetAddr(),
	})

	if err := httpServer.StartWithContext(ctx); err != nil {
		log.Fatal(ctx, "Server failed to start", err, logger.Fields{})
	}
}

// logRoutes logs all registered routes for debugging
func logRoutes(log logger.Logger) {
	ctx := context.Background()
	routes := routes.GetRoutes()
	log.Info(ctx, "Registered routes", logger.Fields{
		"count": len(routes),
	})
	for route, description := range routes {
		log.Info(ctx, "Route registered", logger.Fields{
			"route":       route,
			"description": description,
		})
	}
}
