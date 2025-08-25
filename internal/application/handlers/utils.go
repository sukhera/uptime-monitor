package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/sukhera/uptime-monitor/internal/domain/service"
	"github.com/sukhera/uptime-monitor/internal/shared/errors"
	"github.com/sukhera/uptime-monitor/internal/shared/logger"
)

// BaseHandler provides common utilities for HTTP handlers
type BaseHandler struct {
	buildInfo BuildInfo
}

// NewBaseHandler creates a new base handler with build info
func NewBaseHandler(buildInfo BuildInfo) *BaseHandler {
	return &BaseHandler{buildInfo: buildInfo}
}

// LogError logs an error with structured logging
func (h *BaseHandler) LogError(message string, err error) {
	log := logger.Get()
	ctx := context.Background()
	log.Error(ctx, message, err, nil)
}

// SetJSONHeaders sets common JSON response headers
func (h *BaseHandler) SetJSONHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

// SetStatusJSONHeaders sets JSON headers with cache control for status endpoints
func (h *BaseHandler) SetStatusJSONHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("X-Service-Version", h.buildInfo.Version)
}

// SetHealthJSONHeaders sets JSON headers for health check endpoints
func (h *BaseHandler) SetHealthJSONHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
}

// WriteJSONError writes a JSON error response with logging
func (h *BaseHandler) WriteJSONError(w http.ResponseWriter, message string, err error, statusCode int) {
	h.LogError(message, err)
	h.SetJSONHeaders(w)
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"error":     message,
		"timestamp": time.Now().UTC(),
	}); err != nil {
		h.LogError("Failed to encode error response", err)
	}
}

// WriteInternalServerError writes an internal server error with logging
func (h *BaseHandler) WriteInternalServerError(w http.ResponseWriter, message string, err error) {
	h.WriteJSONError(w, message, err, http.StatusInternalServerError)
}

// WriteBadRequestError writes a bad request error with logging
func (h *BaseHandler) WriteBadRequestError(w http.ResponseWriter, message string, err error) {
	h.WriteJSONError(w, message, err, http.StatusBadRequest)
}

// WriteNotFoundError writes a not found error with logging
func (h *BaseHandler) WriteNotFoundError(w http.ResponseWriter, message string, err error) {
	h.WriteJSONError(w, message, err, http.StatusNotFound)
}

// WriteJSON encodes and writes a JSON response with error handling
func (h *BaseHandler) WriteJSON(w http.ResponseWriter, data interface{}, errorMessage string) {
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.WriteInternalServerError(w, errorMessage, err)
	}
}

// WriteJSONWithHeaders writes JSON response with custom headers
func (h *BaseHandler) WriteJSONWithHeaders(w http.ResponseWriter, data interface{}, errorMessage string, headerSetter func(http.ResponseWriter)) {
	headerSetter(w)
	h.WriteJSON(w, data, errorMessage)
}

// HandleRepositoryError handles repository errors using proper error type checking
// Returns true if error was handled, false if caller should handle it
func (h *BaseHandler) HandleRepositoryError(w http.ResponseWriter, err error, operation string) bool {
	if err == nil {
		return false
	}

	// Check for specific service errors first
	if err == service.ErrServiceNotFound {
		h.WriteNotFoundError(w, "Service not found", err)
		return true
	}

	// Check for repository implementation errors
	if err == service.ErrRepositoryNotImplemented ||
		err == service.ErrManualStatusNotImplemented ||
		err == service.ErrManualStatusClearNotImplemented {
		h.WriteInternalServerError(w, err.Error(), err)
		return true
	}

	// Check for error types using the error system
	if errors.IsNotFound(err) {
		h.WriteNotFoundError(w, "Service not found", err)
		return true
	}

	if errors.IsValidation(err) {
		h.WriteBadRequestError(w, err.Error(), err)
		return true
	}

	if errors.IsInternal(err) {
		h.WriteInternalServerError(w, "Failed to "+operation, err)
		return true
	}

	// Not handled, caller should handle
	return false
}
