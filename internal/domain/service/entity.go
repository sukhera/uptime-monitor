package service

import (
	"time"
)

// Service represents a monitored service
type Service struct {
	Name           string            `bson:"name" json:"name"`
	Slug           string            `bson:"slug" json:"slug"`
	URL            string            `bson:"url" json:"url"`
	Headers        map[string]string `bson:"headers" json:"headers"`
	ExpectedStatus int               `bson:"expected_status" json:"expected_status"`
	Enabled        bool              `bson:"enabled" json:"enabled"`
	CreatedAt      time.Time         `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time         `bson:"updated_at" json:"updated_at"`
}

// ServiceStatus represents the current status of a service
type ServiceStatus struct {
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	Latency   int64     `json:"latency_ms"`
	UpdatedAt time.Time `json:"updated_at"`
	Error     string    `json:"error,omitempty"`
}

// StatusLog represents a health check result
type StatusLog struct {
	ServiceName string    `bson:"service_name" json:"service_name"`
	Status      string    `bson:"status" json:"status"`
	Latency     int64     `bson:"latency_ms" json:"latency_ms"`
	StatusCode  int       `bson:"status_code" json:"status_code"`
	Error       string    `bson:"error,omitempty" json:"error,omitempty"`
	Timestamp   time.Time `bson:"timestamp" json:"timestamp"`
}

// Validate validates the service entity
func (s *Service) Validate() error {
	if s.Name == "" {
		return ErrServiceNameRequired
	}
	if s.URL == "" {
		return ErrServiceURLRequired
	}
	if s.ExpectedStatus < 100 || s.ExpectedStatus > 599 {
		return ErrInvalidExpectedStatus
	}
	return nil
}

// IsOperational returns true if the service is operational
func (ss *ServiceStatus) IsOperational() bool {
	return ss.Status == "operational"
}

// IsDegraded returns true if the service is degraded
func (ss *ServiceStatus) IsDegraded() bool {
	return ss.Status == "degraded"
}

// IsDown returns true if the service is down
func (ss *ServiceStatus) IsDown() bool {
	return ss.Status == "down"
}
