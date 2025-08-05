package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sukhera/uptime-monitor/internal/domain/service"
	mongodb "github.com/sukhera/uptime-monitor/internal/infrastructure/database/mongo"
)

type StatusHandler struct {
	db *mongodb.Database
}

func NewStatusHandler(db *mongodb.Database) *StatusHandler {
	return &StatusHandler{db: db}
}

// GetStatus retrieves the current status of all monitored services
func (h *StatusHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// If no database is available, return empty status array
	if h.db == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

		// Return empty array
		if err := json.NewEncoder(w).Encode([]service.ServiceStatus{}); err != nil {
			h.logError("failed to encode empty response", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		return
	}

	// Add timeout to prevent long-running queries
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	opts := options.Find().SetSort(bson.D{primitive.E{Key: "timestamp", Value: -1}}).SetLimit(100)
	cursor, err := h.db.StatusLogsCollection().Find(ctx, bson.M{}, opts)
	if err != nil {
		h.logError("failed to query status logs", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			h.logError("failed to close cursor", err)
		}
	}()

	var logs []bson.M
	if err = cursor.All(ctx, &logs); err != nil {
		h.logError("failed to decode status logs", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	serviceMap := make(map[string]service.ServiceStatus)
	for _, log := range logs {
		serviceName, ok := log["service_name"].(string)
		if !ok {
			h.logError("invalid service name in log", fmt.Errorf("service_name is not a string"))
			continue
		}

		if _, exists := serviceMap[serviceName]; !exists {
			// Convert primitive.DateTime to time.Time
			var updatedAt time.Time
			if timestamp, ok := log["timestamp"].(time.Time); ok {
				updatedAt = timestamp
			} else if primitiveDateTime, ok := log["timestamp"].(primitive.DateTime); ok {
				updatedAt = primitiveDateTime.Time()
			} else {
				h.logError("invalid timestamp format", fmt.Errorf("timestamp is not a valid time format"))
				continue
			}

			status, ok := log["status"].(string)
			if !ok {
				h.logError("invalid status in log", fmt.Errorf("status is not a string"))
				continue
			}

			latency, ok := log["latency_ms"].(int64)
			if !ok {
				h.logError("invalid latency in log", fmt.Errorf("latency_ms is not an int64"))
				continue
			}

			serviceMap[serviceName] = service.ServiceStatus{
				Name:      serviceName,
				Status:    status,
				Latency:   latency,
				UpdatedAt: updatedAt,
			}
		}
	}

	var statuses []service.ServiceStatus
	for _, status := range serviceMap {
		statuses = append(statuses, status)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	// Return array directly instead of object
	if err := json.NewEncoder(w).Encode(statuses); err != nil {
		h.logError("failed to encode response", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// HealthCheck provides a simple health check endpoint
func (h *StatusHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Test database connectivity if database is available
	if h.db != nil {
		if err := h.db.Client.Ping(ctx, nil); err != nil {
			h.logError("database health check failed", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			if err := json.NewEncoder(w).Encode(map[string]interface{}{
				"status":    "unhealthy",
				"error":     "database connection failed",
				"timestamp": time.Now().UTC(),
			}); err != nil {
				h.logError("failed to encode health check response", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().UTC(),
		"version":   "1.0.0", // TODO: Get from build info
	}); err != nil {
		h.logError("failed to encode health check response", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// GetIncidents returns a list of incidents
func (h *StatusHandler) GetIncidents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// For now, return empty incidents array directly
	response := []interface{}{}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logError("failed to encode incidents response", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// GetMaintenance returns maintenance schedule
func (h *StatusHandler) GetMaintenance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// For now, return empty maintenance array directly
	response := []interface{}{}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logError("failed to encode maintenance response", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// GetTest returns a test response
func (h *StatusHandler) GetTest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "test route works"}); err != nil {
		h.logError("failed to encode test response", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// GetDebug returns debug information
func (h *StatusHandler) GetDebug(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	if _, err := w.Write([]byte("Debug route works")); err != nil {
		h.logError("failed to write debug response", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// logError provides structured error logging
func (h *StatusHandler) logError(message string, err error) {
	// TODO: Replace with proper structured logging (e.g., logrus, zap)
	fmt.Printf("[ERROR] %s: %v\n", message, err)
}
