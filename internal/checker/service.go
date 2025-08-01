package checker

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"github.com/sukhera/uptime-monitor/internal/database"
	"github.com/sukhera/uptime-monitor/internal/models"
)

type Service struct {
	db *database.DB
}

func NewService(db *database.DB) *Service {
	return &Service{db: db}
}

func (s *Service) RunHealthChecks(ctx context.Context) error {
	cursor, err := s.db.ServicesCollection().Find(ctx, bson.M{"enabled": true})
	if err != nil {
		return fmt.Errorf("error querying services: %w", err)
	}
	defer cursor.Close(ctx)

	var services []models.Service
	if err = cursor.All(ctx, &services); err != nil {
		return fmt.Errorf("error decoding services: %w", err)
	}

	for _, service := range services {
		statusLog := s.checkService(service)
		
		if _, err := s.db.StatusLogsCollection().InsertOne(ctx, statusLog); err != nil {
			return fmt.Errorf("error inserting status log for %s: %w", service.Name, err)
		}
	}

	return nil
}

func (s *Service) checkService(service models.Service) models.StatusLog {
	req, err := http.NewRequest("GET", service.URL, nil)
	if err != nil {
		return models.StatusLog{
			ServiceName: service.Name,
			Status:      "down",
			Latency:     0,
			StatusCode:  0,
			Error:       fmt.Sprintf("Failed to create request: %v", err),
			Timestamp:   time.Now(),
		}
	}

	for k, v := range service.Headers {
		req.Header.Set(k, v)
	}

	client := http.Client{Timeout: 10 * time.Second}
	start := time.Now()
	resp, err := client.Do(req)
	latency := time.Since(start).Milliseconds()

	statusLog := models.StatusLog{
		ServiceName: service.Name,
		Latency:     latency,
		Timestamp:   time.Now(),
	}

	if err != nil {
		statusLog.Status = "down"
		statusLog.Error = fmt.Sprintf("Request failed: %v", err)
		return statusLog
	}
	defer resp.Body.Close()

	statusLog.StatusCode = resp.StatusCode

	if resp.StatusCode == service.ExpectedStatus {
		statusLog.Status = "operational"
	} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		statusLog.Status = "degraded"
	} else {
		statusLog.Status = "down"
		statusLog.Error = fmt.Sprintf("Unexpected status code: %d", resp.StatusCode)
	}

	return statusLog
}