package handlers

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/sukhera/uptime-monitor/internal/domain/service"
	mongodb "github.com/sukhera/uptime-monitor/internal/infrastructure/database/mongo"
	"github.com/sukhera/uptime-monitor/internal/shared/logger"
	"github.com/sukhera/uptime-monitor/internal/shared/utils"
)

// WebhookHandler handles webhook operations
type WebhookHandler struct {
	*BaseHandler
	db     *mongodb.Database
	logger logger.Logger
}

// NewWebhookHandler creates a new webhook handler
func NewWebhookHandler(db *mongodb.Database, buildInfo BuildInfo) *WebhookHandler {
	return &WebhookHandler{
		BaseHandler: NewBaseHandler(buildInfo),
		db:          db,
		logger:      logger.Get(),
	}
}

// HandleWebhook processes incoming webhook status updates
func (h *WebhookHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if r.Method != http.MethodPost {
		h.WriteJSONError(w, "Method not allowed", fmt.Errorf("method %s not allowed", r.Method), http.StatusMethodNotAllowed)
		return
	}

	// Extract service slug from URL path
	serviceSlug := h.extractServiceSlugFromPath(r.URL.Path)
	if serviceSlug == "" {
		h.WriteBadRequestError(w, "Invalid service slug in URL", fmt.Errorf("could not extract service slug from path"))
		return
	}
	
	// Sanitize and validate service slug
	serviceSlug = utils.SanitizeUserInput(serviceSlug)
	if !utils.ValidateSlug(serviceSlug) {
		h.WriteBadRequestError(w, "Invalid service slug format", fmt.Errorf("service slug must contain only alphanumeric characters, hyphens, and underscores"))
		return
	}

	// Read and parse webhook payload first (before database operations)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.WriteBadRequestError(w, "Invalid request body", err)
		return
	}
	defer func() { _ = r.Body.Close() }()

	// Parse webhook payload
	var payload service.WebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		h.WriteBadRequestError(w, "Invalid JSON payload", err)
		return
	}

	// Sanitize payload fields
	payload.Status = utils.SanitizeUserInput(payload.Status)
	payload.Message = utils.SanitizeUserInput(payload.Message)
	
	// Sanitize metadata
	if payload.Metadata != nil {
		payload.Metadata = utils.SanitizeMap(payload.Metadata)
	}

	// Validate payload
	if err := h.validateWebhookPayload(&payload); err != nil {
		h.WriteBadRequestError(w, "Invalid payload", err)
		return
	}

	// Get service from database to validate webhook
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	svc, err := h.getServiceBySlug(ctx, serviceSlug)
	if err != nil {
		if h.HandleRepositoryError(w, err, "retrieve service") {
			return
		}
		// Fallback for unhandled errors
		h.WriteInternalServerError(w, "Failed to retrieve service", err)
		return
	}

	if !svc.IsWebhookService() {
		h.WriteBadRequestError(w, "Service not configured for webhooks", service.ErrWebhookNotConfigured)
		return
	}

	// Validate webhook signature if secret is configured
	if svc.WebhookSecret != "" {
		signature := r.Header.Get("X-Webhook-Signature")
		if !h.validateWebhookSignature(body, signature, svc.WebhookSecret) {
			h.WriteJSONError(w, "Invalid signature", service.ErrWebhookInvalidSignature, http.StatusUnauthorized)
			return
		}
	}

	// Create status log entry
	timestamp := time.Now().UTC()
	if payload.Timestamp != nil {
		timestamp = *payload.Timestamp
	}

	latency := int64(0)
	if payload.Latency != nil {
		latency = *payload.Latency
	}

	statusLog := service.StatusLog{
		ServiceName: svc.Name,
		Status:      payload.Status,
		Latency:     latency,
		StatusCode:  200, // Default for webhook updates
		Error:       payload.Message,
		Timestamp:   timestamp,
	}

	// Save to database
	if err := h.saveStatusLog(ctx, &statusLog); err != nil {
		h.WriteInternalServerError(w, "Failed to process webhook", err)
		return
	}

	// Log successful webhook processing - logger handles sanitization internally
	h.logger.Info(ctx, "Webhook processed successfully", logger.Fields{
		"service": svc.Name,
		"slug":    serviceSlug,
		"status":  payload.Status,
		"latency": latency,
	})

	// Return success response
	h.SetJSONHeaders(w)
	w.WriteHeader(http.StatusOK)
	h.WriteJSON(w, map[string]interface{}{
		"message":   "Status updated successfully",
		"service":   svc.Name,
		"status":    payload.Status,
		"timestamp": timestamp,
	}, "Failed to encode webhook response")
}

// extractServiceSlugFromPath extracts the service slug from the URL path
// Expected format: /api/v1/webhook/{service-slug}
func (h *WebhookHandler) extractServiceSlugFromPath(urlPath string) string {
	// Handle empty path
	if urlPath == "" {
		return ""
	}

	// Parse the URL path
	parsedURL, err := url.Parse(urlPath)
	if err != nil {
		return ""
	}

	// Check if the path ends with a trailing slash
	if strings.HasSuffix(parsedURL.Path, "/") {
		return ""
	}

	// Split path and get the last segment
	dir, file := path.Split(parsedURL.Path)
	if file == "" {
		// If no file, get the last directory
		cleanDir := path.Clean(dir)
		if cleanDir == "." || cleanDir == "/" {
			return ""
		}
		_, file = path.Split(cleanDir)
	}

	return file
}

// validateWebhookSignature validates HMAC-SHA256 webhook signature
func (h *WebhookHandler) validateWebhookSignature(body []byte, signature, secret string) bool {
	if signature == "" {
		return false
	}

	// Remove "sha256=" prefix if present
	if len(signature) > 7 && signature[:7] == "sha256=" {
		signature = signature[7:]
	}

	// Calculate expected signature
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expectedMAC := mac.Sum(nil)
	expectedSignature := hex.EncodeToString(expectedMAC)

	// Use constant-time comparison to prevent timing attacks
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

// validateWebhookPayload validates the webhook payload
func (h *WebhookHandler) validateWebhookPayload(payload *service.WebhookPayload) error {
	if payload.Status == "" {
		return fmt.Errorf("status is required")
	}

	// Use security utility for status validation
	if !utils.ValidateStatusValue(payload.Status) {
		return fmt.Errorf("invalid status value: %s", payload.Status)
	}

	// Validate latency if provided
	if payload.Latency != nil {
		if *payload.Latency < 0 || *payload.Latency > 300000 { // Max 5 minutes
			return fmt.Errorf("latency must be between 0 and 300000ms")
		}
	}

	// Validate message length if provided
	if len(payload.Message) > 1000 {
		return fmt.Errorf("message too long (max 1000 characters)")
	}

	return nil
}

// getServiceBySlug retrieves a service by its slug
func (h *WebhookHandler) getServiceBySlug(ctx context.Context, slug string) (*service.Service, error) {
	if h.db == nil {
		return nil, service.ErrRepositoryNotImplemented
	}

	// Create service repository and use it
	serviceRepo := mongodb.NewServiceRepository(h.db)
	return serviceRepo.GetBySlug(ctx, slug)
}

// saveStatusLog saves a status log to the database
func (h *WebhookHandler) saveStatusLog(ctx context.Context, log *service.StatusLog) error {
	if h.db == nil {
		return service.ErrRepositoryNotImplemented
	}

	// Create status log repository and use it
	statusLogRepo := mongodb.NewStatusLogRepository(h.db)
	return statusLogRepo.Create(ctx, log)
}

// GetWebhookExamples returns webhook usage examples for documentation
func (h *WebhookHandler) GetWebhookExamples(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.WriteJSONError(w, "Method not allowed", fmt.Errorf("method %s not allowed", r.Method), http.StatusMethodNotAllowed)
		return
	}

	examples := map[string]interface{}{
		"webhook_url_format": "/api/v1/webhook/{service-slug}",
		"method":             "POST",
		"headers": map[string]string{
			"Content-Type":        "application/json",
			"X-Webhook-Signature": "sha256=your_hmac_signature (optional)",
		},
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
		"curl_example": `curl -X POST /api/v1/webhook/my-service \
  -H "Content-Type: application/json" \
  -H "X-Webhook-Signature: sha256=YOUR_SIGNATURE" \
  -d '{"status": "operational", "latency_ms": 150}'`,
		"signature_calculation": map[string]string{
			"algorithm":   "HMAC-SHA256",
			"secret":      "Your webhook secret (configured per service)",
			"payload":     "Raw JSON request body",
			"header":      "X-Webhook-Signature: sha256={hex_signature}",
			"description": "Calculate HMAC-SHA256 of the raw request body using your webhook secret",
		},
	}

	h.SetJSONHeaders(w)
	w.WriteHeader(http.StatusOK)
	h.WriteJSON(w, examples, "Failed to encode webhook examples")
}
