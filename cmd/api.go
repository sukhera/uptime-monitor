package cmd

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/sukhera/uptime-monitor/internal/application/handlers"
	"github.com/sukhera/uptime-monitor/internal/application/middleware"
	mongodb "github.com/sukhera/uptime-monitor/internal/infrastructure/database/mongo"
	"github.com/sukhera/uptime-monitor/internal/shared/config"
	"github.com/sukhera/uptime-monitor/internal/shared/logger"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Start the REST API server",
	Long: `Start the REST API server that provides:

- Status endpoint (/api/status) - Get current service status
- Health endpoint (/api/health) - API health check
- CORS support for web dashboard
- Structured JSON responses
- Database integration for persistence

Example:
  status-page api --port 8080 --db-url mongodb://localhost:27017`,
	Run: runAPI,
}

var (
	port   string
	dbURL  string
	dbName string
)

func init() {
	rootCmd.AddCommand(apiCmd)

	// API-specific flags
	apiCmd.Flags().StringVarP(&port, "port", "p", "8080", "Port to run the API server on")
	apiCmd.Flags().StringVar(&dbURL, "db-url", "mongodb://localhost:27017", "MongoDB connection URL")
	apiCmd.Flags().StringVar(&dbName, "db-name", "status_page", "MongoDB database name")

	// Bind flags to viper
	bindFlagToViper(apiCmd, "api.port", "port")
	bindFlagToViper(apiCmd, "database.url", "db-url")
	bindFlagToViper(apiCmd, "database.name", "db-name")
}

func runAPI(cmd *cobra.Command, args []string) {
	// Initialize logger
	log := logger.Get()
	ctx := context.Background()

	// Load configuration with proper precedence (flags > env > config file)
	cfg := config.LoadFromViper()

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatal(ctx, "Invalid configuration", err, logger.Fields{})
	}

	// Initialize database
	db, err := mongodb.NewConnection(cfg.Database.URI, cfg.Database.Name)
	if err != nil {
		log.Fatal(ctx, "Failed to connect to database", err, logger.Fields{"db_url": cfg.Database.URI, "db_name": cfg.Database.Name})
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Error(ctx, "Error closing database connection", err, nil)
		}
	}()

	// Get build info
	version, commit, buildDate := GetBuildInfo()
	buildInfo := handlers.BuildInfo{
		Version:   version,
		Commit:    commit,
		BuildDate: buildDate,
	}
	if buildInfo.Version == "" {
		buildInfo.Version = "dev"
	}

	// Initialize handlers
	statusHandler := handlers.NewStatusHandler(db, buildInfo)

	// Setup routes using gorilla/mux
	router := http.NewServeMux()

	// Add versioned routes (v1)
	router.HandleFunc("/api/v1/status", statusHandler.GetStatus)
	router.HandleFunc("/api/v1/health", statusHandler.HealthCheck)
	router.HandleFunc("/api/v1/incidents", statusHandler.GetIncidents)
	router.HandleFunc("/api/v1/maintenance", statusHandler.GetMaintenance)

	// Backward compatibility - redirect old routes to v1
	router.HandleFunc("/api/status", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/api/v1/status", http.StatusMovedPermanently)
	})
	router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/api/v1/health", http.StatusMovedPermanently)
	})
	router.HandleFunc("/api/incidents", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/api/v1/incidents", http.StatusMovedPermanently)
	})
	router.HandleFunc("/api/maintenance", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/api/v1/maintenance", http.StatusMovedPermanently)
	})

	// Test endpoint (versioned)
	router.HandleFunc("/api/v1/test", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]string{"message": "test route works"}); err != nil {
			log.Error(ctx, "Error encoding test response", err, nil)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	})
	// Debug endpoint (versioned)
	router.HandleFunc("/api/v1/debug", statusHandler.GetDebug)

	// Backward compatibility - redirect old debug route to v1
	router.HandleFunc("/debug", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/api/v1/debug", http.StatusMovedPermanently)
	})

	// Debug: Print registered routes
	log.Info(ctx, "Registered routes", logger.Fields{
		"v1_routes": []string{
			"GET /api/v1/status",
			"GET /api/v1/health",
			"GET /api/v1/incidents",
			"GET /api/v1/maintenance",
			"GET /api/v1/test",
			"GET /api/v1/debug",
		},
		"legacy_redirects": []string{
			"GET /api/status → /api/v1/status",
			"GET /api/health → /api/v1/health",
			"GET /api/incidents → /api/v1/incidents",
			"GET /api/maintenance → /api/v1/maintenance",
			"GET /debug → /api/v1/debug",
		},
	})

	// Add middleware chain
	var handler http.Handler = router

	// Add API versioning middleware
	handler = middleware.APIVersion("v1")(handler)

	// CORS middleware is available but not enabled by default
	// corsMiddleware := middleware.NewCORS()
	// handler = corsMiddleware.Handler(handler)

	// Create server
	apiPort := viper.GetString("api.port")
	server := &http.Server{
		Addr:              ":" + apiPort,
		Handler:           handler,
		ReadHeaderTimeout: 30 * time.Second, // Prevent Slowloris attacks
	}

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Info(ctx, "Shutting down server", nil)

		// Create context with timeout for graceful shutdown
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Error(ctx, "Error during graceful shutdown, forcing close", err, nil)
			if err := server.Close(); err != nil {
				log.Error(ctx, "Error forcing server close", err, nil)
			}
		} else {
			log.Info(ctx, "Server shutdown completed gracefully", nil)
		}
	}()

	// Start server
	log.Info(ctx, "Starting API server", logger.Fields{
		"port":              apiPort,
		"health_check_url":  "http://localhost:" + apiPort + "/api/v1/health",
		"status_url":        "http://localhost:" + apiPort + "/api/v1/status",
		"legacy_health_url": "http://localhost:" + apiPort + "/api/health (redirects to v1)",
		"legacy_status_url": "http://localhost:" + apiPort + "/api/status (redirects to v1)",
	})

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(ctx, "Failed to start server", err, logger.Fields{"port": apiPort})
	}
}
