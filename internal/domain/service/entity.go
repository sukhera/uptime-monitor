package service

import (
	"time"
)

// ServiceType enum for different monitoring types
type ServiceType string

const (
	ServiceTypeHTTP    ServiceType = "http"
	ServiceTypeTCP     ServiceType = "tcp"
	ServiceTypeDNS     ServiceType = "dns"
	ServiceTypeWebhook ServiceType = "webhook"
)

// ManualStatusOverride represents manual status management
type ManualStatusOverride struct {
	Status    string     `bson:"status" json:"status"`
	Reason    string     `bson:"reason" json:"reason"`
	SetBy     string     `bson:"set_by" json:"set_by"`
	SetAt     time.Time  `bson:"set_at" json:"set_at"`
	ExpiresAt *time.Time `bson:"expires_at,omitempty" json:"expires_at,omitempty"`
}

// WebhookPayload represents incoming webhook data
type WebhookPayload struct {
	Status    string                 `json:"status"`
	Latency   *int64                 `json:"latency_ms,omitempty"`
	Message   string                 `json:"message,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	Timestamp *time.Time             `json:"timestamp,omitempty"`
}

// Service represents a monitored service (ENHANCED)
type Service struct {
	Name           string            `bson:"name" json:"name"`
	Slug           string            `bson:"slug" json:"slug"`
	URL            string            `bson:"url" json:"url"`
	Headers        map[string]string `bson:"headers" json:"headers"`
	ExpectedStatus int               `bson:"expected_status" json:"expected_status"`
	Enabled        bool              `bson:"enabled" json:"enabled"`

	// NEW: Integration Fields
	ServiceType         ServiceType            `bson:"service_type" json:"service_type"`
	WebhookURL          string                 `bson:"webhook_url" json:"webhook_url,omitempty"`
	WebhookSecret       string                 `bson:"webhook_secret" json:"-"` // Never expose in JSON
	ManualStatus        *ManualStatusOverride  `bson:"manual_status,omitempty" json:"manual_status,omitempty"`
	IntegrationMetadata map[string]interface{} `bson:"integration_metadata,omitempty" json:"integration_metadata,omitempty"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
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

	// URL is only required for non-webhook services
	if s.ServiceType != ServiceTypeWebhook && s.URL == "" {
		return ErrServiceURLRequired
	}

	if s.ExpectedStatus != 0 && (s.ExpectedStatus < 100 || s.ExpectedStatus > 599) {
		return ErrInvalidExpectedStatus
	}

	// Set default service type if not specified
	if s.ServiceType == "" {
		s.ServiceType = ServiceTypeHTTP
	}

	// Set default expected status for HTTP services
	if s.ServiceType == ServiceTypeHTTP && s.ExpectedStatus == 0 {
		s.ExpectedStatus = 200
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

// HasManualOverride returns true if the service has an active manual status override
func (s *Service) HasManualOverride() bool {
	if s.ManualStatus == nil {
		return false
	}

	// Check if the override has expired
	if s.ManualStatus.ExpiresAt != nil && time.Now().UTC().After(*s.ManualStatus.ExpiresAt) {
		return false
	}

	return true
}

// GetEffectiveStatus returns the effective status considering manual overrides
func (s *Service) GetEffectiveStatus(currentStatus string) string {
	if s.HasManualOverride() {
		return s.ManualStatus.Status
	}
	return currentStatus
}

// IsWebhookService returns true if this is a webhook-based service
func (s *Service) IsWebhookService() bool {
	return s.ServiceType == ServiceTypeWebhook
}

// RequiresURL returns true if the service type requires a URL
func (s *Service) RequiresURL() bool {
	return s.ServiceType != ServiceTypeWebhook
}

// IsValidStatus checks if a status value is valid
func IsValidStatus(status string) bool {
	validStatuses := []string{"operational", "degraded", "down", "maintenance"}
	for _, valid := range validStatuses {
		if status == valid {
			return true
		}
	}
	return false
}

// IsValidServiceType checks if a service type is valid
func IsValidServiceType(serviceType ServiceType) bool {
	validTypes := []ServiceType{ServiceTypeHTTP, ServiceTypeTCP, ServiceTypeDNS, ServiceTypeWebhook}
	for _, valid := range validTypes {
		if serviceType == valid {
			return true
		}
	}
	return false
}
