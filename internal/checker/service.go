package checker

import (
	"context"
	"fmt"
	"log"
	"net/http"
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
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			log.Printf("[ERROR] Failed to close cursor: %v", err)
		}
	}()

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
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			log.Printf("[ERROR] Failed to close cursor: %v", err)
		}
	}()

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
