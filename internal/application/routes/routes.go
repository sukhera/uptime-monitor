package routes

import (
	"net/http"

	"github.com/sukhera/uptime-monitor/internal/application/handlers"
)

// SetupRoutes configures and returns the HTTP router with all routes
func SetupRoutes(
	statusHandler *handlers.StatusHandler,
	webhookHandler *handlers.WebhookHandler,
	integrationHandler *handlers.IntegrationHandler,
) *http.ServeMux {
	router := http.NewServeMux()

	// Existing API routes (v1)
	router.HandleFunc("/api/v1/status", statusHandler.GetStatus)
	router.HandleFunc("/api/v1/health", statusHandler.HealthCheck)
	router.HandleFunc("/api/v1/incidents", statusHandler.GetIncidents)
	router.HandleFunc("/api/v1/maintenance", statusHandler.GetMaintenance)
	router.HandleFunc("/api/v1/test", statusHandler.GetTest)
	router.HandleFunc("/api/v1/debug", statusHandler.GetDebug)

	// NEW: Webhook routes
	router.HandleFunc("/api/v1/webhook/", handleWebhookRoutes(webhookHandler))
	router.HandleFunc("/api/v1/webhook/examples", webhookHandler.GetWebhookExamples)

	// NEW: Integration routes
	router.HandleFunc("/api/v1/integration/services/", handleIntegrationServiceRoutes(integrationHandler))
	router.HandleFunc("/api/v1/integration/services/bulk-import", integrationHandler.BulkImport)
	router.HandleFunc("/api/v1/integration/docs", integrationHandler.GetAPIDocumentation)

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

// handleWebhookRoutes handles webhook-related routes with path parameter extraction
func handleWebhookRoutes(handler *handlers.WebhookHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// The path will be like /api/v1/webhook/{service-slug}
		// We'll extract the service slug and delegate to the handler
		handler.HandleWebhook(w, r)
	}
}

// handleIntegrationServiceRoutes handles integration service routes with path parameter extraction
func handleIntegrationServiceRoutes(handler *handlers.IntegrationHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Handle different integration service endpoints:
		// - GET /api/v1/integration/services/{id} -> GetIntegrationDetails
		// - POST /api/v1/integration/services/{id}/manual-status -> SetManualStatus
		// - DELETE /api/v1/integration/services/{id}/manual-status -> ClearManualStatus

		path := r.URL.Path

		// Check if this is a manual-status route
		if len(path) > len("/api/v1/integration/services/") {
			if r.Method == http.MethodPost && endsWithManualStatus(path) {
				handler.SetManualStatus(w, r)
				return
			} else if r.Method == http.MethodDelete && endsWithManualStatus(path) {
				handler.ClearManualStatus(w, r)
				return
			}
		}

		// Default to integration details
		if r.Method == http.MethodGet {
			handler.GetIntegrationDetails(w, r)
			return
		}

		// Method not allowed
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// endsWithManualStatus checks if the path ends with /manual-status
func endsWithManualStatus(path string) bool {
	suffix := "/manual-status"
	return len(path) >= len(suffix) && path[len(path)-len(suffix):] == suffix
}

// GetRoutes returns a map of all registered routes for documentation
func GetRoutes() map[string]string {
	return map[string]string{
		"GET /api/v1/status":                                     "Get current system status",
		"GET /api/v1/health":                                     "Health check endpoint",
		"GET /api/v1/incidents":                                  "Get incidents list",
		"GET /api/v1/maintenance":                                "Get maintenance schedule",
		"GET /api/v1/test":                                       "Test endpoint",
		"GET /api/v1/debug":                                      "Debug endpoint",
		"POST /api/v1/webhook/{service-slug}":                    "Update service status via webhook",
		"GET /api/v1/webhook/examples":                           "Get webhook usage examples",
		"GET /api/v1/integration/services/{id}":                  "Get integration details for service",
		"POST /api/v1/integration/services/{id}/manual-status":   "Set manual status override",
		"DELETE /api/v1/integration/services/{id}/manual-status": "Clear manual status override",
		"POST /api/v1/integration/services/bulk-import":          "Bulk import services",
		"GET /api/v1/integration/docs":                           "Get API documentation",
	}
}
