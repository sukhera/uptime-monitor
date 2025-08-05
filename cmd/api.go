package cmd

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/sukhera/uptime-monitor/internal/application/handlers"
	// "github.com/sukhera/uptime-monitor/internal/application/middleware"
	mongodb "github.com/sukhera/uptime-monitor/internal/infrastructure/database/mongo"
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
	if err := viper.BindPFlag("api.port", apiCmd.Flags().Lookup("port")); err != nil {
		log.Fatalf("Failed to bind api.port flag: %v", err)
	}
	if err := viper.BindPFlag("database.url", apiCmd.Flags().Lookup("db-url")); err != nil {
		log.Fatalf("Failed to bind database.url flag: %v", err)
	}
	if err := viper.BindPFlag("database.name", apiCmd.Flags().Lookup("db-name")); err != nil {
		log.Fatalf("Failed to bind database.name flag: %v", err)
	}
}

func runAPI(cmd *cobra.Command, args []string) {
	// Initialize database
	db, err := mongodb.NewConnection(dbURL, dbName)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	// Initialize handlers
	statusHandler := handlers.NewStatusHandler(db)

	// Setup routes using gorilla/mux
	router := http.NewServeMux()

	// Add routes
	router.HandleFunc("/api/status", statusHandler.GetStatus)
	router.HandleFunc("/api/health", statusHandler.HealthCheck)
	router.HandleFunc("/api/incidents", statusHandler.GetIncidents)
	router.HandleFunc("/api/maintenance", statusHandler.GetMaintenance)
	router.HandleFunc("/api/test", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]string{"message": "test route works"}); err != nil {
			log.Printf("Error encoding test response: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	})
	router.HandleFunc("/debug", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		if _, err := w.Write([]byte("Debug route works")); err != nil {
			log.Printf("Error writing debug response: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	})

	// Debug: Print registered routes
	log.Printf("Registered routes:")
	log.Printf("  GET /api/status")
	log.Printf("  GET /api/health")
	log.Printf("  GET /api/incidents")
	log.Printf("  GET /api/maintenance")

	// Add CORS middleware
	// corsMiddleware := middleware.NewCORS()
	// handler := corsMiddleware.Handler(router)
	handler := router // Temporarily use router directly

	// Create server
	server := &http.Server{
		Addr:              ":" + port,
		Handler:           handler,
		ReadHeaderTimeout: 30 * time.Second, // Prevent Slowloris attacks
	}

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Println("Shutting down server...")
		if err := server.Close(); err != nil {
			log.Printf("Error shutting down server: %v", err)
		}
	}()

	// Start server
	log.Printf("Starting API server on port %s", port)
	log.Printf("Health check: http://localhost:%s/api/health", port)
	log.Printf("Status endpoint: http://localhost:%s/api/status", port)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}
