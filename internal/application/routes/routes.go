package routes

import (
	"net/http"

	"github.com/sukhera/uptime-monitor/internal/application/handlers"
)

// SetupRoutes configures and returns the HTTP router with all routes
func SetupRoutes(statusHandler *handlers.StatusHandler) *http.ServeMux {
	router := http.NewServeMux()

	// Add versioned routes (v1)
	router.HandleFunc("/api/v1/status", statusHandler.GetStatus)
	router.HandleFunc("/api/v1/health", statusHandler.HealthCheck)
	router.HandleFunc("/api/v1/incidents", statusHandler.GetIncidents)
	router.HandleFunc("/api/v1/maintenance", statusHandler.GetMaintenance)
	router.HandleFunc("/api/v1/test", statusHandler.GetTest)
	router.HandleFunc("/api/v1/debug", statusHandler.GetDebug)

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
	router.HandleFunc("/api/test", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/api/v1/test", http.StatusMovedPermanently)
	})
	router.HandleFunc("/debug", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/api/v1/debug", http.StatusMovedPermanently)
	})

	return router
}

// GetRoutes returns a map of all registered routes for documentation
func GetRoutes() map[string]string {
	return map[string]string{
		"GET /api/v1/status":      "Get current system status",
		"GET /api/v1/health":      "Health check endpoint",
		"GET /api/v1/incidents":   "Get incidents list",
		"GET /api/v1/maintenance": "Get maintenance schedule",
		"GET /api/v1/test":        "Test endpoint",
		"GET /api/v1/debug":       "Debug endpoint",
	}
}
