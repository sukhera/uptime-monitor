package checker

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/sukhera/uptime-monitor/internal/domain/service"
	mongodb "github.com/sukhera/uptime-monitor/internal/infrastructure/database/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	statusDown        = "down"
	statusOperational = "operational"
	statusDegraded    = "degraded"
)

// HTTPClient interface for mocking HTTP requests
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// ServiceInterface defines the interface for the health checker service
type ServiceInterface interface {
	RunHealthChecks(ctx context.Context) error
}

type Service struct {
	db     mongodb.Interface
	client HTTPClient
}

// ServiceOption is a function that configures a Service
type ServiceOption func(*Service)

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(client HTTPClient) ServiceOption {
	return func(s *Service) {
		s.client = client
	}
}

// WithTimeout sets the HTTP client timeout
func WithTimeout(timeout time.Duration) ServiceOption {
	return func(s *Service) {
		if s.client == nil {
			s.client = &http.Client{Timeout: timeout}
		} else {
			if httpClient, ok := s.client.(*http.Client); ok {
				httpClient.Timeout = timeout
			}
		}
	}
}

// NewService creates a new Service with the given options
func NewService(db mongodb.Interface, options ...ServiceOption) *Service {
	service := &Service{
		db: db,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	for _, option := range options {
		option(service)
	}

	return service
}

// NewServiceWithClient creates a service with a custom HTTP client (useful for testing)
// Deprecated: Use NewService with WithHTTPClient option instead
func NewServiceWithClient(db mongodb.Interface, client HTTPClient) *Service {
	return NewService(db, WithHTTPClient(client))
}

// RunHealthChecks runs health checks using the command pattern
func (s *Service) RunHealthChecks(ctx context.Context) error {
	cursor, err := s.db.ServicesCollection().Find(ctx, bson.M{"enabled": true})
	if err != nil {
		return fmt.Errorf("error querying services: %w", err)
	}
	defer cursor.Close(ctx)

	var services []service.Service
	if err = cursor.All(ctx, &services); err != nil {
		return fmt.Errorf("error decoding services: %w", err)
	}

	// Create command invoker
	invoker := NewHealthCheckInvoker()

	// Create commands for each service
	for _, service := range services {
		command := NewHTTPHealthCheckCommand(service, s.client)
		invoker.AddCommand(command)
	}

	// Execute all health check commands concurrently
	statusLogs := invoker.ExecuteAll(ctx)

	// Store results in database
	for _, statusLog := range statusLogs {
		if _, err := s.db.StatusLogsCollection().InsertOne(ctx, statusLog); err != nil {
			// Log error but continue with other results
			log.Printf("[ERROR] Failed to insert status log for %s: %v", statusLog.ServiceName, err)
		}
	}

	return nil
}

// RunHealthChecksWithObservers runs health checks and notifies observers
func (s *Service) RunHealthChecksWithObservers(ctx context.Context, subject *HealthCheckSubject) error {
	cursor, err := s.db.ServicesCollection().Find(ctx, bson.M{"enabled": true})
	if err != nil {
		return fmt.Errorf("error querying services: %w", err)
	}
	defer cursor.Close(ctx)

	var services []service.Service
	if err = cursor.All(ctx, &services); err != nil {
		return fmt.Errorf("error decoding services: %w", err)
	}

	// Create command invoker
	invoker := NewHealthCheckInvoker()

	// Create commands for each service
	for _, service := range services {
		command := NewHTTPHealthCheckCommand(service, s.client)
		invoker.AddCommand(command)
	}

	// Execute all health check commands concurrently
	statusLogs := invoker.ExecuteAll(ctx)

	// Store results and notify observers
	for _, statusLog := range statusLogs {
		// Store in database
		if _, err := s.db.StatusLogsCollection().InsertOne(ctx, statusLog); err != nil {
			log.Printf("[ERROR] Failed to insert status log for %s: %v", statusLog.ServiceName, err)
		}

		// Notify observers
		event := HealthCheckEvent{
			ServiceName: statusLog.ServiceName,
			Status:      statusLog.Status,
			Latency:     statusLog.Latency,
			StatusCode:  statusLog.StatusCode,
			Error:       statusLog.Error,
			Timestamp:   statusLog.Timestamp.Unix(),
		}
		subject.Notify(ctx, event)
	}

	return nil
}

func (s *Service) checkURL(svc service.Service, wg *sync.WaitGroup, statusLogs chan<- service.StatusLog) {
	defer wg.Done()

	statusLog := s.checkService(svc)
	log.Printf("[INFO] %s: %s (status: %d, latency: %dms)",
		svc.Name, statusLog.Status, statusLog.StatusCode, statusLog.Latency)

	statusLogs <- statusLog
}

func (s *Service) checkService(svc service.Service) service.StatusLog {
	const maxRetries = 3
	const retryDelay = 500 * time.Millisecond

	req, err := http.NewRequest("GET", svc.URL, nil)
	if err != nil {
		log.Printf("[ERROR] Failed to create request for %s: %v", svc.Name, err)
		return service.StatusLog{
			ServiceName: svc.Name,
			Status:      "down",
			Latency:     0,
			StatusCode:  0,
			Error:       fmt.Sprintf("Failed to create request: %v", err),
			Timestamp:   time.Now(),
		}
	}

	for k, v := range svc.Headers {
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
				attempt, svc.Name, retryDelay, err)
			time.Sleep(retryDelay)
		}
	}

	statusLog := service.StatusLog{
		ServiceName: svc.Name,
		Latency:     latency,
		Timestamp:   time.Now(),
	}

	if lastErr != nil && resp == nil {
		statusLog.Status = statusDown
		statusLog.Error = fmt.Sprintf("Request failed after %d attempts: %v", maxRetries, lastErr)
		log.Printf("[ERROR] Request failed for %s after %d attempts: %v", svc.Name, maxRetries, lastErr)
		return statusLog
	}
	defer resp.Body.Close()

	statusLog.StatusCode = resp.StatusCode

	if resp.StatusCode == svc.ExpectedStatus {
		statusLog.Status = statusOperational
	} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		statusLog.Status = statusDegraded
	} else {
		statusLog.Status = statusDown
		statusLog.Error = fmt.Sprintf("Unexpected status code: %d", resp.StatusCode)
	}

	return statusLog
}
