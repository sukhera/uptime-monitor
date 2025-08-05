package service

import (
	"context"
)

// Repository defines the interface for service data access
type Repository interface {
	// Create creates a new service
	Create(ctx context.Context, service *Service) error

	// GetByID retrieves a service by ID
	GetByID(ctx context.Context, id string) (*Service, error)

	// GetBySlug retrieves a service by slug
	GetBySlug(ctx context.Context, slug string) (*Service, error)

	// GetAll retrieves all services
	GetAll(ctx context.Context) ([]*Service, error)

	// GetEnabled retrieves all enabled services
	GetEnabled(ctx context.Context) ([]*Service, error)

	// Update updates a service
	Update(ctx context.Context, service *Service) error

	// Delete deletes a service
	Delete(ctx context.Context, id string) error

	// SaveStatusLog saves a status log entry
	SaveStatusLog(ctx context.Context, log *StatusLog) error

	// GetLatestStatus retrieves the latest status for all services
	GetLatestStatus(ctx context.Context) ([]*ServiceStatus, error)

	// GetStatusHistory retrieves status history for a service
	GetStatusHistory(ctx context.Context, serviceName string, limit int) ([]*StatusLog, error)
}
