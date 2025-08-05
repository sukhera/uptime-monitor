package service

import (
	"github.com/sukhera/uptime-monitor/internal/shared/errors"
)

// Service-specific errors
var (
	ErrServiceNameRequired   = errors.NewValidationError("service name is required")
	ErrServiceURLRequired    = errors.NewValidationError("service URL is required")
	ErrInvalidExpectedStatus = errors.NewValidationError("expected status must be between 100 and 599")
	ErrServiceNotFound       = errors.NewNotFoundError("service not found")
	ErrServiceAlreadyExists  = errors.NewConflictError("service already exists")
	ErrServiceDisabled       = errors.NewValidationError("service is disabled")
	ErrInvalidServiceStatus  = errors.NewValidationError("invalid service status")
)
