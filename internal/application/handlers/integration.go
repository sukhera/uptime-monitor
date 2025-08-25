package handlers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/sukhera/uptime-monitor/internal/domain/service"
	mongodb "github.com/sukhera/uptime-monitor/internal/infrastructure/database/mongo"
	"github.com/sukhera/uptime-monitor/internal/shared/logger"
)

// IntegrationHandler handles integration operations like manual status overrides and bulk imports
type IntegrationHandler struct {
	*BaseHandler
	db      *mongodb.Database
	logger  logger.Logger
	baseURL string // For generating webhook URLs
}

// NewIntegrationHandler creates a new integration handler
func NewIntegrationHandler(db *mongodb.Database, baseURL string, buildInfo BuildInfo) *IntegrationHandler {
	return &IntegrationHandler{
		BaseHandler: NewBaseHandler(buildInfo),
		db:          db,
		logger:      logger.Get(),
		baseURL:     baseURL,
	}
}

// SetManualStatus sets manual status override for a service
func (h *IntegrationHandler) SetManualStatus(w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if r.Method != http.MethodPost {
		h.WriteJSONError(w, "Method not allowed", fmt.Errorf("method %s not allowed", r.Method), http.StatusMethodNotAllowed)
		return
	}

	// Extract service ID from URL path
	serviceID := h.extractServiceIDFromPath(r.URL.Path)
	if serviceID == "" {
		h.WriteBadRequestError(w, "Invalid service ID in URL", fmt.Errorf("could not extract service ID from path: %s", r.URL.Path))
		return
	}

	// Parse request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.WriteBadRequestError(w, "Invalid request body", err)
		return
	}
	defer func() { _ = r.Body.Close() }()

	var request struct {
		Status    string     `json:"status"`
		Reason    string     `json:"reason"`
		ExpiresAt *time.Time `json:"expires_at,omitempty"`
	}

	if err := json.Unmarshal(body, &request); err != nil {
		h.WriteBadRequestError(w, "Invalid JSON payload", err)
		return
	}

	// Validate request
	if request.Status == "" {
		h.WriteBadRequestError(w, "Status is required", service.ErrManualStatusRequired)
		return
	}

	if request.Reason == "" {
		h.WriteBadRequestError(w, "Reason is required", service.ErrManualReasonRequired)
		return
	}

	if !service.IsValidStatus(request.Status) {
		h.WriteBadRequestError(w, "Invalid status value", service.ErrInvalidServiceStatus)
		return
	}

	// Create manual status override
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Get current user (implement auth as needed)
	userID := h.getCurrentUser(r)

	manualOverride := &service.ManualStatusOverride{
		Status:    request.Status,
		Reason:    request.Reason,
		SetBy:     userID,
		SetAt:     time.Now().UTC(),
		ExpiresAt: request.ExpiresAt,
	}

	// Set manual status (placeholder - will implement with repository update)
	err = h.setManualStatus(ctx, serviceID, manualOverride)
	if err != nil {
		if h.HandleRepositoryError(w, err, "update status") {
			return
		}
		// Fallback for unhandled errors
		h.WriteInternalServerError(w, "Failed to update status", err)
		return
	}

	// Sanitize user-controlled input for logging to prevent log injection
	safeServiceID := strings.ReplaceAll(strings.ReplaceAll(serviceID, "\n", ""), "\r", "")
	safeStatus := strings.ReplaceAll(strings.ReplaceAll(request.Status, "\n", ""), "\r", "")
	safeReason := strings.ReplaceAll(strings.ReplaceAll(request.Reason, "\n", ""), "\r", "")
	safeUserID := strings.ReplaceAll(strings.ReplaceAll(userID, "\n", ""), "\r", "")
	
	h.logger.Info(ctx, "Manual status set", logger.Fields{
		"service_id": safeServiceID,
		"status":     safeStatus,
		"reason":     safeReason,
		"user":       safeUserID,
	})

	// Return success response
	h.SetJSONHeaders(w)
	w.WriteHeader(http.StatusOK)
	h.WriteJSON(w, map[string]interface{}{
		"message":    "Manual status set successfully",
		"service_id": serviceID,
		"status":     request.Status,
		"reason":     request.Reason,
		"set_by":     userID,
		"set_at":     manualOverride.SetAt,
		"expires_at": request.ExpiresAt,
	}, "Failed to encode manual status response")
}

// ClearManualStatus removes manual status override for a service
func (h *IntegrationHandler) ClearManualStatus(w http.ResponseWriter, r *http.Request) {
	// Only accept DELETE requests
	if r.Method != http.MethodDelete {
		h.WriteJSONError(w, "Method not allowed", fmt.Errorf("method %s not allowed", r.Method), http.StatusMethodNotAllowed)
		return
	}

	// Extract service ID from URL path
	serviceID := h.extractServiceIDFromPath(r.URL.Path)
	if serviceID == "" {
		h.WriteBadRequestError(w, "Invalid service ID in URL", fmt.Errorf("could not extract service ID from path: %s", r.URL.Path))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Clear manual status (placeholder - will implement with repository update)
	err := h.clearManualStatus(ctx, serviceID)
	if err != nil {
		if h.HandleRepositoryError(w, err, "clear manual status") {
			return
		}
		// Fallback for unhandled errors
		h.WriteInternalServerError(w, "Failed to clear manual status", err)
		return
	}

	userID := h.getCurrentUser(r)
	
	// Sanitize user-controlled input for logging to prevent log injection
	safeServiceID := strings.ReplaceAll(strings.ReplaceAll(serviceID, "\n", ""), "\r", "")
	safeUserID := strings.ReplaceAll(strings.ReplaceAll(userID, "\n", ""), "\r", "")
	
	h.logger.Info(ctx, "Manual status cleared", logger.Fields{
		"service_id": safeServiceID,
		"user":       safeUserID,
	})

	// Return success response
	h.SetJSONHeaders(w)
	w.WriteHeader(http.StatusOK)
	h.WriteJSON(w, map[string]interface{}{
		"message":    "Manual status cleared successfully",
		"service_id": serviceID,
		"cleared_by": userID,
		"cleared_at": time.Now().UTC(),
	}, "Failed to encode clear manual status response")
}

// GetIntegrationDetails returns webhook URL and integration information
func (h *IntegrationHandler) GetIntegrationDetails(w http.ResponseWriter, r *http.Request) {
	// Only accept GET requests
	if r.Method != http.MethodGet {
		h.WriteJSONError(w, "Method not allowed", fmt.Errorf("method %s not allowed", r.Method), http.StatusMethodNotAllowed)
		return
	}

	// Extract service ID from URL path
	serviceID := h.extractServiceIDFromPath(r.URL.Path)
	if serviceID == "" {
		h.WriteBadRequestError(w, "Invalid service ID in URL", fmt.Errorf("could not extract service ID from path: %s", r.URL.Path))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Get service (placeholder - will implement with repository update)
	svc, err := h.getServiceByID(ctx, serviceID)
	if err != nil {
		if h.HandleRepositoryError(w, err, "retrieve service") {
			return
		}
		// Fallback for unhandled errors
		h.WriteInternalServerError(w, "Failed to retrieve service", err)
		return
	}

	// Generate webhook URL if service is webhook type and URL not exists
	if svc.ServiceType == service.ServiceTypeWebhook && svc.WebhookURL == "" {
		webhookURL := fmt.Sprintf("%s/api/v1/webhook/%s", h.baseURL, svc.Slug)
		svc.WebhookURL = webhookURL

		// Generate webhook secret if not exists
		if svc.WebhookSecret == "" {
			svc.WebhookSecret = h.generateWebhookSecret()
		}

		// Update service with webhook details (placeholder)
		if err := h.updateService(ctx, svc); err != nil {
			h.logger.Error(ctx, "Failed to update webhook details", err, logger.Fields{"service_id": serviceID})
		}
	}

	// Build response
	response := map[string]interface{}{
		"service_id":           svc.Slug,
		"service_name":         svc.Name,
		"service_type":         svc.ServiceType,
		"webhook_url":          svc.WebhookURL,
		"has_secret":           svc.WebhookSecret != "",
		"manual_status":        svc.ManualStatus,
		"integration_metadata": svc.IntegrationMetadata,
	}

	// Include webhook examples for webhook services
	if svc.ServiceType == service.ServiceTypeWebhook {
		response["webhook_examples"] = h.getWebhookExamples(svc.WebhookURL)
	}

	h.SetJSONHeaders(w)
	w.WriteHeader(http.StatusOK)
	h.WriteJSON(w, response, "Failed to encode integration details response")
}

// BulkImport handles bulk service import from JSON/YAML
func (h *IntegrationHandler) BulkImport(w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if r.Method != http.MethodPost {
		h.WriteJSONError(w, "Method not allowed", fmt.Errorf("method %s not allowed", r.Method), http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.WriteBadRequestError(w, "Invalid request body", err)
		return
	}
	defer func() { _ = r.Body.Close() }()

	var request struct {
		Services []service.Service `json:"services"`
	}

	if err := json.Unmarshal(body, &request); err != nil {
		h.WriteBadRequestError(w, "Invalid JSON payload", err)
		return
	}

	if len(request.Services) == 0 {
		h.WriteBadRequestError(w, "No services provided", fmt.Errorf("services array is empty"))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	var imported, failed int
	var errors []string

	// Process each service
	for i, svc := range request.Services {
		// Set timestamps
		now := time.Now().UTC()
		svc.CreatedAt = now
		svc.UpdatedAt = now

		// Generate slug if not provided
		if svc.Slug == "" {
			svc.Slug = h.generateSlug(svc.Name)
		}

		// Validate service
		if err := svc.Validate(); err != nil {
			errors = append(errors, fmt.Sprintf("Service %d (%s): %v", i+1, svc.Name, err))
			failed++
			continue
		}

		// Generate webhook details for webhook services
		if svc.ServiceType == service.ServiceTypeWebhook {
			svc.WebhookURL = fmt.Sprintf("%s/api/v1/webhook/%s", h.baseURL, svc.Slug)
			svc.WebhookSecret = h.generateWebhookSecret()
		}

		// Create service (placeholder - will implement with repository update)
		if err := h.createService(ctx, &svc); err != nil {
			errors = append(errors, fmt.Sprintf("Service %d (%s): %v", i+1, svc.Name, err))
			failed++
			continue
		}

		imported++
	}

	// Build response
	response := map[string]interface{}{
		"imported": imported,
		"failed":   failed,
		"total":    len(request.Services),
	}

	if len(errors) > 0 {
		response["errors"] = errors
	}

	status := http.StatusOK
	if failed > 0 && imported == 0 {
		status = http.StatusBadRequest
	} else if failed > 0 {
		status = http.StatusPartialContent
	}

	h.logger.Info(ctx, "Bulk import completed", logger.Fields{
		"imported": imported,
		"failed":   failed,
		"total":    len(request.Services),
	})

	h.SetJSONHeaders(w)
	w.WriteHeader(status)
	h.WriteJSON(w, response, "Failed to encode bulk import response")
}

// GetAPIDocumentation returns built-in API documentation
func (h *IntegrationHandler) GetAPIDocumentation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.WriteJSONError(w, "Method not allowed", fmt.Errorf("method %s not allowed", r.Method), http.StatusMethodNotAllowed)
		return
	}

	docs := map[string]interface{}{
		"title":       "Status Page Integration API",
		"version":     h.buildInfo.Version,
		"description": "API endpoints for webhook integration, manual status management, and bulk operations",
		"base_url":    h.baseURL,
		"endpoints": map[string]interface{}{
			"webhook": map[string]interface{}{
				"url":         "/api/v1/webhook/{service-slug}",
				"method":      "POST",
				"description": "Update service status via webhook",
				"headers":     []string{"Content-Type: application/json", "X-Webhook-Signature: sha256=signature (optional)"},
				"example":     h.getWebhookExamples(""),
			},
			"manual_status": map[string]interface{}{
				"set": map[string]interface{}{
					"url":         "/api/v1/integration/services/{id}/manual-status",
					"method":      "POST",
					"description": "Set manual status override for a service",
					"payload": map[string]interface{}{
						"status":     "operational|degraded|down|maintenance",
						"reason":     "Reason for manual override",
						"expires_at": "2024-01-01T12:00:00Z (optional)",
					},
				},
				"clear": map[string]interface{}{
					"url":         "/api/v1/integration/services/{id}/manual-status",
					"method":      "DELETE",
					"description": "Clear manual status override for a service",
				},
			},
			"integration_details": map[string]interface{}{
				"url":         "/api/v1/integration/services/{id}",
				"method":      "GET",
				"description": "Get integration details for a service including webhook URL",
			},
			"bulk_import": map[string]interface{}{
				"url":         "/api/v1/integration/services/bulk-import",
				"method":      "POST",
				"description": "Import multiple services at once",
				"payload": map[string]interface{}{
					"services": []map[string]interface{}{
						{
							"name":                 "Service Name",
							"slug":                 "service-slug (optional)",
							"url":                  "https://example.com (required for http/tcp/dns)",
							"service_type":         "http|tcp|dns|webhook",
							"expected_status":      200,
							"enabled":              true,
							"headers":              map[string]string{"Authorization": "Bearer token"},
							"integration_metadata": map[string]interface{}{"key": "value"},
						},
					},
				},
			},
		},
		"authentication": map[string]interface{}{
			"webhook_signatures": map[string]string{
				"algorithm":   "HMAC-SHA256",
				"header":      "X-Webhook-Signature",
				"format":      "sha256={hex_signature}",
				"description": "Optional HMAC-SHA256 signature of request body using webhook secret",
			},
		},
		"status_values": []string{"operational", "degraded", "down", "maintenance"},
		"service_types": []string{"http", "tcp", "dns", "webhook"},
	}

	h.SetJSONHeaders(w)
	w.WriteHeader(http.StatusOK)
	h.WriteJSON(w, docs, "Failed to encode API documentation")
}

// Helper methods

// extractServiceIDFromPath extracts the service ID from the URL path
func (h *IntegrationHandler) extractServiceIDFromPath(urlPath string) string {
	// Parse the URL path and extract the service ID
	parsedURL, err := url.Parse(urlPath)
	if err != nil {
		return ""
	}

	// Expected format: /api/v1/integration/services/{id}/...
	// Split path and verify it matches the expected pattern
	parts := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")

	// Must have at least 5 parts: api, v1, integration, services, {id}
	if len(parts) < 5 {
		return ""
	}

	// Check if it matches the expected integration API pattern
	if parts[0] == "api" && parts[1] == "v1" && parts[2] == "integration" && parts[3] == "services" {
		return parts[4]
	}

	return ""
}

// getCurrentUser returns the current user (placeholder for auth implementation)
func (h *IntegrationHandler) getCurrentUser(r *http.Request) string {
	// TODO: Implement based on your authentication system
	// For now, return a placeholder
	return "system"
}

// generateWebhookSecret generates a secure webhook secret
func (h *IntegrationHandler) generateWebhookSecret() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to timestamp-based secret
		return hex.EncodeToString([]byte(fmt.Sprintf("webhook-%d", time.Now().UnixNano())))
	}
	return hex.EncodeToString(bytes)
}

// generateSlug generates a URL-friendly slug from a name
func (h *IntegrationHandler) generateSlug(name string) string {
	// Simple slug generation - convert to lowercase and replace spaces with hyphens
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")

	// Remove non-alphanumeric characters except hyphens
	result := strings.Builder{}
	for _, r := range slug {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}

	return result.String()
}

// getWebhookExamples returns webhook usage examples
func (h *IntegrationHandler) getWebhookExamples(webhookURL string) map[string]interface{} {
	if webhookURL == "" {
		webhookURL = "{webhook_url}"
	}

	return map[string]interface{}{
		"curl": fmt.Sprintf(`curl -X POST %s \
  -H "Content-Type: application/json" \
  -H "X-Webhook-Signature: sha256=YOUR_SIGNATURE" \
  -d '{"status": "operational", "latency_ms": 150, "message": "All systems operational"}'`, webhookURL),
		"payload_example": map[string]interface{}{
			"status":     "operational",
			"latency_ms": 150,
			"message":    "All systems operational",
			"metadata": map[string]interface{}{
				"region":  "us-west-1",
				"version": "1.2.3",
			},
			"timestamp": "2024-01-01T12:00:00Z",
		},
		"status_values": []string{"operational", "degraded", "down", "maintenance"},
	}
}

// Implementation methods using repository pattern

func (h *IntegrationHandler) setManualStatus(ctx context.Context, serviceID string, override *service.ManualStatusOverride) error {
	if h.db == nil {
		return service.ErrManualStatusNotImplemented
	}

	serviceRepo := mongodb.NewServiceRepository(h.db)
	return serviceRepo.SetManualStatus(ctx, serviceID, override)
}

func (h *IntegrationHandler) clearManualStatus(ctx context.Context, serviceID string) error {
	if h.db == nil {
		return service.ErrManualStatusClearNotImplemented
	}

	serviceRepo := mongodb.NewServiceRepository(h.db)
	return serviceRepo.ClearManualStatus(ctx, serviceID)
}

func (h *IntegrationHandler) getServiceByID(ctx context.Context, serviceID string) (*service.Service, error) {
	if h.db == nil {
		return nil, service.ErrServiceNotFound
	}

	serviceRepo := mongodb.NewServiceRepository(h.db)

	// Try by ID first, then by slug
	svc, err := serviceRepo.GetByID(ctx, serviceID)
	if err != nil {
		// If ID lookup fails, try by slug
		return serviceRepo.GetBySlug(ctx, serviceID)
	}

	return svc, nil
}

func (h *IntegrationHandler) updateService(ctx context.Context, svc *service.Service) error {
	if h.db == nil {
		return service.ErrRepositoryNotImplemented
	}

	serviceRepo := mongodb.NewServiceRepository(h.db)
	return serviceRepo.Update(ctx, svc)
}

func (h *IntegrationHandler) createService(ctx context.Context, svc *service.Service) error {
	if h.db == nil {
		return service.ErrRepositoryNotImplemented
	}

	serviceRepo := mongodb.NewServiceRepository(h.db)
	return serviceRepo.Create(ctx, svc)
}
