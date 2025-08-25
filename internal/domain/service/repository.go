package service

import (
	"context"
	"time"
)

// Repository defines the interface for service data access
type Repository interface {
	// Basic CRUD operations
	Create(ctx context.Context, service *Service) error
	GetByID(ctx context.Context, id string) (*Service, error)
	GetBySlug(ctx context.Context, slug string) (*Service, error)
	GetAll(ctx context.Context) ([]*Service, error)
	GetEnabled(ctx context.Context) ([]*Service, error)
	Update(ctx context.Context, service *Service) error
	Delete(ctx context.Context, id string) error

	// Integration-specific operations
	FindByType(ctx context.Context, serviceType ServiceType) ([]*Service, error)
	BulkCreate(ctx context.Context, services []*Service) error
	SetManualStatus(ctx context.Context, serviceID string, override *ManualStatusOverride) error
	ClearManualStatus(ctx context.Context, serviceID string) error

	// Status log operations
	SaveStatusLog(ctx context.Context, log *StatusLog) error
	GetLatestStatus(ctx context.Context) ([]*ServiceStatus, error)
	GetStatusHistory(ctx context.Context, serviceName string, limit int) ([]*StatusLog, error)
}

// StatusLogRepository defines operations for status logs (separate interface for cleaner separation)
type StatusLogRepository interface {
	Create(ctx context.Context, log *StatusLog) error
	FindByServiceName(ctx context.Context, serviceName string, limit int) ([]*StatusLog, error)
	FindRecent(ctx context.Context, limit int) ([]*StatusLog, error)
	DeleteOlderThan(ctx context.Context, cutoff time.Time) error
	GetLatestByService(ctx context.Context) (map[string]*StatusLog, error)
}
