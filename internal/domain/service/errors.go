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

	// Integration-specific errors
	ErrInvalidServiceType      = errors.NewValidationError("invalid service type")
	ErrWebhookInvalidSignature = errors.NewValidationError("invalid webhook signature")
	ErrWebhookInvalidPayload   = errors.NewValidationError("invalid webhook payload")
	ErrWebhookServiceNotFound  = errors.NewNotFoundError("webhook service not found")
	ErrWebhookNotConfigured    = errors.NewValidationError("service not configured for webhooks")
	ErrManualStatusRequired    = errors.NewValidationError("manual status is required")
	ErrManualReasonRequired    = errors.NewValidationError("manual status reason is required")
	ErrManualStatusExpired     = errors.NewValidationError("manual status override has expired")
	ErrBulkImportFailed        = errors.NewValidationError("bulk import operation failed")
	ErrInvalidWebhookSecret    = errors.NewValidationError("invalid webhook secret")

	// Repository operation errors
	ErrRepositoryNotImplemented        = errors.NewInternalError("service repository method not implemented")
	ErrManualStatusNotImplemented      = errors.NewInternalError("manual status setting not implemented")
	ErrManualStatusClearNotImplemented = errors.NewInternalError("manual status clearing not implemented")
)
