package checker

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"github.com/sukhera/uptime-monitor/internal/database"
	"github.com/sukhera/uptime-monitor/internal/models"
)

type Service struct {
	db     database.DatabaseInterface
	client *http.Client
}

func NewService(db database.DatabaseInterface) *Service {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	return &Service{
		db:     db,
		client: client,
	}
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

	log.Printf("[INFO] Checking %d services", len(services))

	var wg sync.WaitGroup
	statusLogs := make(chan models.StatusLog, len(services))

	for _, service := range services {
		wg.Add(1)
		go s.checkURL(service, &wg, statusLogs)
	}

	go func() {
		wg.Wait()
		close(statusLogs)
	}()

	for statusLog := range statusLogs {
		if _, err := s.db.StatusLogsCollection().InsertOne(ctx, statusLog); err != nil {
			log.Printf("[ERROR] Failed to insert status log for %s: %v", statusLog.ServiceName, err)
		}
	}

	log.Printf("[INFO] Health checks completed for %d services", len(services))
	return nil
}

func (s *Service) checkURL(service models.Service, wg *sync.WaitGroup, statusLogs chan<- models.StatusLog) {
	defer wg.Done()

	statusLog := s.checkService(service)
	log.Printf("[INFO] %s: %s (status: %d, latency: %dms)", 
		service.Name, statusLog.Status, statusLog.StatusCode, statusLog.Latency)
	
	statusLogs <- statusLog
}

func (s *Service) checkService(service models.Service) models.StatusLog {
	const maxRetries = 3
	const retryDelay = 500 * time.Millisecond

	req, err := http.NewRequest("GET", service.URL, nil)
	if err != nil {
		log.Printf("[ERROR] Failed to create request for %s: %v", service.Name, err)
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

	var resp *http.Response
	var latency int64
	var lastErr error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		start := time.Now()
		resp, err = s.client.Do(req)
		latency = time.Since(start).Milliseconds()

		if err == nil {
			break
		}

		lastErr = err
		if attempt < maxRetries {
			log.Printf("[WARN] Attempt %d failed for %s, retrying in %v: %v", 
				attempt, service.Name, retryDelay, err)
			time.Sleep(retryDelay)
		}
	}

	statusLog := models.StatusLog{
		ServiceName: service.Name,
		Latency:     latency,
		Timestamp:   time.Now(),
	}

	if lastErr != nil && resp == nil {
		statusLog.Status = "down"
		statusLog.Error = fmt.Sprintf("Request failed after %d attempts: %v", maxRetries, lastErr)
		log.Printf("[ERROR] Request failed for %s after %d attempts: %v", service.Name, maxRetries, lastErr)
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