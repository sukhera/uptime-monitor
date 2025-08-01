package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sukhera/uptime-monitor/internal/database"
	"github.com/sukhera/uptime-monitor/internal/models"
)

type StatusHandler struct {
	db *database.DB
}

func NewStatusHandler(db *database.DB) *StatusHandler {
	return &StatusHandler{db: db}
}

func (h *StatusHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	opts := options.Find().SetSort(bson.D{{"timestamp", -1}}).SetLimit(100)
	cursor, err := h.db.StatusLogsCollection().Find(ctx, bson.M{}, opts)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Printf("Error querying status: %v", err)
		return
	}
	defer cursor.Close(ctx)

	var logs []bson.M
	if err = cursor.All(ctx, &logs); err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Printf("Error decoding status: %v", err)
		return
	}

	serviceMap := make(map[string]models.ServiceStatus)
	for _, log := range logs {
		serviceName := log["service_name"].(string)
		if _, exists := serviceMap[serviceName]; !exists {
			// Convert primitive.DateTime to time.Time
			var updatedAt time.Time
			if timestamp, ok := log["timestamp"].(time.Time); ok {
				updatedAt = timestamp
			} else if primitiveDateTime, ok := log["timestamp"].(primitive.DateTime); ok {
				updatedAt = primitiveDateTime.Time()
			}
			
			serviceMap[serviceName] = models.ServiceStatus{
				Name:      serviceName,
				Status:    log["status"].(string),
				Latency:   log["latency_ms"].(int64),
				UpdatedAt: updatedAt,
			}
		}
	}

	var statuses []models.ServiceStatus
	for _, status := range serviceMap {
		statuses = append(statuses, status)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(statuses)
}

func (h *StatusHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}