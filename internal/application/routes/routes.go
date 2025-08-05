package routes

import (
	"net/http"

	"github.com/sukhera/uptime-monitor/internal/application/handlers"
)

// SetupRoutes configures and returns the HTTP router with all routes
func SetupRoutes(statusHandler *handlers.StatusHandler) *http.ServeMux {
	router := http.NewServeMux()

	// API routes
	router.HandleFunc("/api/status", statusHandler.GetStatus)
	router.HandleFunc("/api/health", statusHandler.HealthCheck)
	router.HandleFunc("/api/incidents", statusHandler.GetIncidents)
	router.HandleFunc("/api/maintenance", statusHandler.GetMaintenance)

	// Debug routes (only in development)
	router.HandleFunc("/api/test", statusHandler.GetTest)
	router.HandleFunc("/debug", statusHandler.GetDebug)

	return router
}

// GetRoutes returns a map of all registered routes for documentation
func GetRoutes() map[string]string {
	return map[string]string{
		"GET /api/status":      "Get current system status",
		"GET /api/health":      "Health check endpoint",
		"GET /api/incidents":   "Get incidents list",
		"GET /api/maintenance": "Get maintenance schedule",
		"GET /api/test":        "Test endpoint",
		"GET /debug":           "Debug endpoint",
	}
}
