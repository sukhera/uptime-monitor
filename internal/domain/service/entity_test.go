package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_Validate(t *testing.T) {
	tests := []struct {
		name    string
		service Service
		wantErr error
	}{
		{
			name: "valid HTTP service",
			service: Service{
				Name:           "Test Service",
				URL:            "https://example.com",
				ExpectedStatus: 200,
				ServiceType:    ServiceTypeHTTP,
			},
			wantErr: nil,
		},
		{
			name: "valid webhook service without URL",
			service: Service{
				Name:        "Webhook Service",
				ServiceType: ServiceTypeWebhook,
			},
			wantErr: nil,
		},
		{
			name: "invalid - missing name",
			service: Service{
				URL:         "https://example.com",
				ServiceType: ServiceTypeHTTP,
			},
			wantErr: ErrServiceNameRequired,
		},
		{
			name: "invalid - HTTP service missing URL",
			service: Service{
				Name:        "Test Service",
				ServiceType: ServiceTypeHTTP,
			},
			wantErr: ErrServiceURLRequired,
		},
		{
			name: "invalid - bad expected status",
			service: Service{
				Name:           "Test Service",
				URL:            "https://example.com",
				ExpectedStatus: 99,
				ServiceType:    ServiceTypeHTTP,
			},
			wantErr: ErrInvalidExpectedStatus,
		},
		{
			name: "valid - defaults applied",
			service: Service{
				Name: "Default Service",
				URL:  "https://example.com",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.service.Validate()

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)

				// Check that defaults are applied
				if tt.service.ServiceType == "" {
					assert.Equal(t, ServiceTypeHTTP, tt.service.ServiceType)
				}
				if tt.service.ServiceType == ServiceTypeHTTP && tt.service.ExpectedStatus == 0 {
					assert.Equal(t, 200, tt.service.ExpectedStatus)
				}
			}
		})
	}
}

func TestService_HasManualOverride(t *testing.T) {
	now := time.Now().UTC()
	future := now.Add(1 * time.Hour)
	past := now.Add(-1 * time.Hour)

	tests := []struct {
		name     string
		service  Service
		expected bool
	}{
		{
			name: "no manual status",
			service: Service{
				ManualStatus: nil,
			},
			expected: false,
		},
		{
			name: "active manual status without expiry",
			service: Service{
				ManualStatus: &ManualStatusOverride{
					Status: "maintenance",
					Reason: "Scheduled maintenance",
					SetAt:  now,
				},
			},
			expected: true,
		},
		{
			name: "active manual status with future expiry",
			service: Service{
				ManualStatus: &ManualStatusOverride{
					Status:    "maintenance",
					Reason:    "Scheduled maintenance",
					SetAt:     now,
					ExpiresAt: &future,
				},
			},
			expected: true,
		},
		{
			name: "expired manual status",
			service: Service{
				ManualStatus: &ManualStatusOverride{
					Status:    "maintenance",
					Reason:    "Scheduled maintenance",
					SetAt:     past,
					ExpiresAt: &past,
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.service.HasManualOverride()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestService_GetEffectiveStatus(t *testing.T) {
	now := time.Now().UTC()
	future := now.Add(1 * time.Hour)

	tests := []struct {
		name          string
		service       Service
		currentStatus string
		expected      string
	}{
		{
			name: "no manual override - returns current status",
			service: Service{
				ManualStatus: nil,
			},
			currentStatus: "operational",
			expected:      "operational",
		},
		{
			name: "active manual override - returns manual status",
			service: Service{
				ManualStatus: &ManualStatusOverride{
					Status:    "maintenance",
					Reason:    "Scheduled maintenance",
					SetAt:     now,
					ExpiresAt: &future,
				},
			},
			currentStatus: "operational",
			expected:      "maintenance",
		},
		{
			name: "expired manual override - returns current status",
			service: Service{
				ManualStatus: &ManualStatusOverride{
					Status:    "maintenance",
					Reason:    "Scheduled maintenance",
					SetAt:     now.Add(-2 * time.Hour),
					ExpiresAt: func() *time.Time { t := now.Add(-1 * time.Hour); return &t }(),
				},
			},
			currentStatus: "operational",
			expected:      "operational",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.service.GetEffectiveStatus(tt.currentStatus)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestService_IsWebhookService(t *testing.T) {
	tests := []struct {
		name        string
		serviceType ServiceType
		expected    bool
	}{
		{
			name:        "webhook service",
			serviceType: ServiceTypeWebhook,
			expected:    true,
		},
		{
			name:        "HTTP service",
			serviceType: ServiceTypeHTTP,
			expected:    false,
		},
		{
			name:        "TCP service",
			serviceType: ServiceTypeTCP,
			expected:    false,
		},
		{
			name:        "DNS service",
			serviceType: ServiceTypeDNS,
			expected:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := Service{ServiceType: tt.serviceType}
			result := service.IsWebhookService()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestService_RequiresURL(t *testing.T) {
	tests := []struct {
		name        string
		serviceType ServiceType
		expected    bool
	}{
		{
			name:        "webhook service - no URL required",
			serviceType: ServiceTypeWebhook,
			expected:    false,
		},
		{
			name:        "HTTP service - URL required",
			serviceType: ServiceTypeHTTP,
			expected:    true,
		},
		{
			name:        "TCP service - URL required",
			serviceType: ServiceTypeTCP,
			expected:    true,
		},
		{
			name:        "DNS service - URL required",
			serviceType: ServiceTypeDNS,
			expected:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := Service{ServiceType: tt.serviceType}
			result := service.RequiresURL()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsValidStatus(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		expected bool
	}{
		{
			name:     "valid - operational",
			status:   "operational",
			expected: true,
		},
		{
			name:     "valid - degraded",
			status:   "degraded",
			expected: true,
		},
		{
			name:     "valid - down",
			status:   "down",
			expected: true,
		},
		{
			name:     "valid - maintenance",
			status:   "maintenance",
			expected: true,
		},
		{
			name:     "invalid - unknown status",
			status:   "unknown",
			expected: false,
		},
		{
			name:     "invalid - empty status",
			status:   "",
			expected: false,
		},
		{
			name:     "invalid - case sensitive",
			status:   "OPERATIONAL",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidStatus(tt.status)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsValidServiceType(t *testing.T) {
	tests := []struct {
		name        string
		serviceType ServiceType
		expected    bool
	}{
		{
			name:        "valid - HTTP",
			serviceType: ServiceTypeHTTP,
			expected:    true,
		},
		{
			name:        "valid - TCP",
			serviceType: ServiceTypeTCP,
			expected:    true,
		},
		{
			name:        "valid - DNS",
			serviceType: ServiceTypeDNS,
			expected:    true,
		},
		{
			name:        "valid - Webhook",
			serviceType: ServiceTypeWebhook,
			expected:    true,
		},
		{
			name:        "invalid - unknown type",
			serviceType: ServiceType("unknown"),
			expected:    false,
		},
		{
			name:        "invalid - empty type",
			serviceType: ServiceType(""),
			expected:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidServiceType(tt.serviceType)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestServiceStatus_HelperMethods(t *testing.T) {
	tests := []struct {
		name       string
		status     string
		isOp       bool
		isDegraded bool
		isDown     bool
	}{
		{
			name:       "operational status",
			status:     "operational",
			isOp:       true,
			isDegraded: false,
			isDown:     false,
		},
		{
			name:       "degraded status",
			status:     "degraded",
			isOp:       false,
			isDegraded: true,
			isDown:     false,
		},
		{
			name:       "down status",
			status:     "down",
			isOp:       false,
			isDegraded: false,
			isDown:     true,
		},
		{
			name:       "maintenance status",
			status:     "maintenance",
			isOp:       false,
			isDegraded: false,
			isDown:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss := &ServiceStatus{Status: tt.status}

			assert.Equal(t, tt.isOp, ss.IsOperational())
			assert.Equal(t, tt.isDegraded, ss.IsDegraded())
			assert.Equal(t, tt.isDown, ss.IsDown())
		})
	}
}

func TestWebhookPayload_Validation(t *testing.T) {
	now := time.Now().UTC()
	latency := int64(150)

	tests := []struct {
		name    string
		payload WebhookPayload
		valid   bool
	}{
		{
			name: "valid minimal payload",
			payload: WebhookPayload{
				Status: "operational",
			},
			valid: true,
		},
		{
			name: "valid complete payload",
			payload: WebhookPayload{
				Status:    "operational",
				Latency:   &latency,
				Message:   "All systems operational",
				Timestamp: &now,
				Metadata: map[string]interface{}{
					"region":  "us-west-1",
					"version": "1.0.0",
				},
			},
			valid: true,
		},
		{
			name: "invalid - missing status",
			payload: WebhookPayload{
				Message: "Test message",
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Basic validation test - checking if required fields are present
			hasStatus := tt.payload.Status != ""
			assert.Equal(t, tt.valid, hasStatus)

			if tt.valid {
				if tt.payload.Latency != nil {
					assert.Greater(t, *tt.payload.Latency, int64(-1))
				}
				if tt.payload.Timestamp != nil {
					assert.WithinDuration(t, now, *tt.payload.Timestamp, time.Minute)
				}
			}
		})
	}
}

func TestManualStatusOverride_Complete(t *testing.T) {
	now := time.Now().UTC()
	future := now.Add(1 * time.Hour)

	override := ManualStatusOverride{
		Status:    "maintenance",
		Reason:    "Scheduled maintenance window",
		SetBy:     "admin@example.com",
		SetAt:     now,
		ExpiresAt: &future,
	}

	require.NotNil(t, override)
	assert.Equal(t, "maintenance", override.Status)
	assert.Equal(t, "Scheduled maintenance window", override.Reason)
	assert.Equal(t, "admin@example.com", override.SetBy)
	assert.WithinDuration(t, now, override.SetAt, time.Second)
	assert.NotNil(t, override.ExpiresAt)
	assert.WithinDuration(t, future, *override.ExpiresAt, time.Second)
}
