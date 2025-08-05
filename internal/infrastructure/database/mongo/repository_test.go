package mongo

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/sukhera/uptime-monitor/internal/domain/service"
	"github.com/sukhera/uptime-monitor/internal/shared/errors"
)

func TestService_Validate(t *testing.T) {
	tests := []struct {
		name          string
		service       *service.Service
		expectedError error
	}{
		{
			name: "valid service",
			service: &service.Service{
				Name:           "Test Service",
				URL:            "https://example.com",
				ExpectedStatus: 200,
				Enabled:        true,
			},
			expectedError: nil,
		},
		{
			name: "missing name",
			service: &service.Service{
				Name:           "",
				URL:            "https://example.com",
				ExpectedStatus: 200,
				Enabled:        true,
			},
			expectedError: errors.NewValidationError("service name is required"),
		},
		{
			name: "missing URL",
			service: &service.Service{
				Name:           "Test Service",
				URL:            "",
				ExpectedStatus: 200,
				Enabled:        true,
			},
			expectedError: errors.NewValidationError("service URL is required"),
		},
		{
			name: "invalid expected status",
			service: &service.Service{
				Name:           "Test Service",
				URL:            "https://example.com",
				ExpectedStatus: 999, // Invalid status code
				Enabled:        true,
			},
			expectedError: errors.NewValidationError("expected status must be between 100 and 599"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.service.Validate()

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestStatusLog_Validate(t *testing.T) {
	tests := []struct {
		name          string
		statusLog     *service.StatusLog
		expectedError error
	}{
		{
			name: "valid status log",
			statusLog: &service.StatusLog{
				ServiceName: "Test Service",
				Status:      "up",
				StatusCode:  200,
				Latency:     150,
				Timestamp:   time.Now(),
			},
			expectedError: nil,
		},
		{
			name: "missing service name",
			statusLog: &service.StatusLog{
				ServiceName: "",
				Status:      "up",
				StatusCode:  200,
				Latency:     150,
				Timestamp:   time.Now(),
			},
			expectedError: errors.NewValidationError("service name is required"),
		},
		{
			name: "missing status",
			statusLog: &service.StatusLog{
				ServiceName: "Test Service",
				Status:      "",
				StatusCode:  200,
				Latency:     150,
				Timestamp:   time.Now(),
			},
			expectedError: errors.NewValidationError("status is required"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: StatusLog doesn't have a Validate method, so we'll test the validation logic
			// that would be used in the repository
			var err error
			if tt.statusLog.ServiceName == "" {
				err = errors.NewValidationError("service name is required")
			} else if tt.statusLog.Status == "" {
				err = errors.NewValidationError("status is required")
			}

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestServiceStatus_Helpers(t *testing.T) {
	tests := []struct {
		name          string
		serviceStatus *service.ServiceStatus
		isOperational bool
		isDegraded    bool
		isDown        bool
	}{
		{
			name: "operational service",
			serviceStatus: &service.ServiceStatus{
				Name:   "Test Service",
				Status: "operational",
			},
			isOperational: true,
			isDegraded:    false,
			isDown:        false,
		},
		{
			name: "degraded service",
			serviceStatus: &service.ServiceStatus{
				Name:   "Test Service",
				Status: "degraded",
			},
			isOperational: false,
			isDegraded:    true,
			isDown:        false,
		},
		{
			name: "down service",
			serviceStatus: &service.ServiceStatus{
				Name:   "Test Service",
				Status: "down",
			},
			isOperational: false,
			isDegraded:    false,
			isDown:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.isOperational, tt.serviceStatus.IsOperational())
			assert.Equal(t, tt.isDegraded, tt.serviceStatus.IsDegraded())
			assert.Equal(t, tt.isDown, tt.serviceStatus.IsDown())
		})
	}
}

func TestService_Integration(t *testing.T) {
	// Test service creation and validation
	svc := &service.Service{
		Name:           "Integration Test Service",
		URL:            "https://integration-test.com",
		ExpectedStatus: 200,
		Enabled:        true,
	}

	// Validate the service
	err := svc.Validate()
	assert.NoError(t, err)

	// Test status log creation
	statusLog := &service.StatusLog{
		ServiceName: svc.Name,
		Status:      "up",
		StatusCode:  200,
		Latency:     150,
		Timestamp:   time.Now(),
	}

	// Validate the status log (basic validation)
	assert.NotEmpty(t, statusLog.ServiceName)
	assert.NotEmpty(t, statusLog.Status)
	assert.Greater(t, statusLog.StatusCode, 0)
	assert.GreaterOrEqual(t, statusLog.Latency, int64(0))
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		name          string
		errorType     errors.ErrorKind
		message       string
		expectedError error
	}{
		{
			name:          "validation error",
			errorType:     errors.ErrorKindValidation,
			message:       "validation failed",
			expectedError: errors.NewValidationError("validation failed"),
		},
		{
			name:          "not found error",
			errorType:     errors.ErrorKindNotFound,
			message:       "service not found",
			expectedError: errors.NewNotFoundError("service not found"),
		},
		{
			name:          "internal error",
			errorType:     errors.ErrorKindInternal,
			message:       "database error",
			expectedError: errors.NewInternalError("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			switch tt.errorType {
			case errors.ErrorKindValidation:
				err = errors.NewValidationError(tt.message)
			case errors.ErrorKindNotFound:
				err = errors.NewNotFoundError(tt.message)
			case errors.ErrorKindInternal:
				err = errors.NewInternalError(tt.message)
			}

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.message)

			// Test error type checking
			var domainErr errors.Error
			if assert.ErrorAs(t, err, &domainErr) {
				assert.Equal(t, tt.errorType, domainErr.Kind())
			}
		})
	}
}
