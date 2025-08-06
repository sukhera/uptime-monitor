package handlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sukhera/uptime-monitor/internal/domain/service"
	mongodb "github.com/sukhera/uptime-monitor/internal/infrastructure/database/mongo"
)

// BuildInfo holds build-time information
type BuildInfo struct {
	Version   string
	Commit    string
	BuildDate string
}

type StatusHandler struct {
	*BaseHandler
	db *mongodb.Database
}

func NewStatusHandler(db *mongodb.Database, buildInfo BuildInfo) *StatusHandler {
	return &StatusHandler{
		BaseHandler: NewBaseHandler(buildInfo),
		db:          db,
	}
}

// GetStatus retrieves the current status of all monitored services
func (h *StatusHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// If no database is available, return empty status array
	if h.db == nil {
		h.SetStatusJSONHeaders(w)
		h.WriteJSON(w, []service.ServiceStatus{}, "failed to encode empty response")
		return
	}

	// Add timeout to prevent long-running queries
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	opts := options.Find().SetSort(bson.D{primitive.E{Key: "timestamp", Value: -1}}).SetLimit(100)
	cursor, err := h.db.StatusLogsCollection().Find(ctx, bson.M{}, opts)
	if err != nil {
		h.WriteInternalServerError(w, "failed to query status logs", err)
		return
	}
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			h.LogError("failed to close cursor", err)
		}
	}()

	var logs []bson.M
	if err = cursor.All(ctx, &logs); err != nil {
		h.WriteInternalServerError(w, "failed to decode status logs", err)
		return
	}

	serviceMap := make(map[string]service.ServiceStatus)
	for _, log := range logs {
		serviceName, ok := log["service_name"].(string)
		if !ok {
			h.LogError("invalid service name in log", errors.New("service_name is not a string"))
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
				h.LogError("invalid timestamp format", errors.New("timestamp is not a valid time format"))
				continue
			}

			status, ok := log["status"].(string)
			if !ok {
				h.LogError("invalid status in log", errors.New("status is not a string"))
				continue
			}

			latency, ok := log["latency_ms"].(int64)
			if !ok {
				h.LogError("invalid latency in log", errors.New("latency_ms is not an int64"))
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

	h.SetStatusJSONHeaders(w)
	h.WriteJSON(w, statuses, "failed to encode response")
}

// HealthCheck provides a simple health check endpoint
func (h *StatusHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Test database connectivity if database is available
	if h.db != nil {
		if err := h.db.Client.Ping(ctx, nil); err != nil {
			h.SetJSONHeaders(w)
			w.WriteHeader(http.StatusServiceUnavailable)
			h.WriteJSON(w, map[string]interface{}{
				"status":    "unhealthy",
				"error":     "database connection failed",
				"timestamp": time.Now().UTC(),
			}, "failed to encode health check response")
			return
		}
	}

	h.SetHealthJSONHeaders(w)
	h.WriteJSON(w, map[string]interface{}{
		"status":     "healthy",
		"timestamp":  time.Now().UTC(),
		"version":    h.buildInfo.Version,
		"commit":     h.buildInfo.Commit,
		"build_date": h.buildInfo.BuildDate,
	}, "failed to encode health check response")
}

// GetIncidents returns a list of incidents
func (h *StatusHandler) GetIncidents(w http.ResponseWriter, r *http.Request) {
	h.SetJSONHeaders(w)
	
	// For now, return empty incidents array directly
	response := []interface{}{}
	
	h.WriteJSON(w, response, "failed to encode incidents response")
}

// GetMaintenance returns maintenance schedule
func (h *StatusHandler) GetMaintenance(w http.ResponseWriter, r *http.Request) {
	h.SetJSONHeaders(w)
	
	// For now, return empty maintenance array directly
	response := []interface{}{}
	
	h.WriteJSON(w, response, "failed to encode maintenance response")
}

// GetTest returns a test response
func (h *StatusHandler) GetTest(w http.ResponseWriter, r *http.Request) {
	h.SetJSONHeaders(w)
	h.WriteJSON(w, map[string]string{"message": "test route works"}, "failed to encode test response")
}

// GetDebug returns debug information
func (h *StatusHandler) GetDebug(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	if _, err := w.Write([]byte("Debug route works")); err != nil {
		h.WriteInternalServerError(w, "failed to write debug response", err)
		return
	}
}

