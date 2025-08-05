package errors

import (
	"fmt"
)

// ErrorKind represents the type of error
type ErrorKind string

const (
	ErrorKindValidation   ErrorKind = "validation"
	ErrorKindNotFound     ErrorKind = "not_found"
	ErrorKindInternal     ErrorKind = "internal"
	ErrorKindConflict     ErrorKind = "conflict"
	ErrorKindUnauthorized ErrorKind = "unauthorized"
	ErrorKindForbidden    ErrorKind = "forbidden"
)

// Error represents a domain error
type Error interface {
	Kind() ErrorKind
	Error() string
	Cause() error
}

// CustomError implements the Error interface
type CustomError struct {
	kind  ErrorKind
	cause error
	msg   string
}

// New creates a new error
func New(msg string, kind ErrorKind) Error {
	return &CustomError{
		kind: kind,
		msg:  msg,
	}
}

// NewWithCause creates a new error with a cause
func NewWithCause(msg string, kind ErrorKind, cause error) Error {
	return &CustomError{
		kind:  kind,
		msg:   msg,
		cause: cause,
	}
}

// NewValidationError creates a validation error
func NewValidationError(msg string) Error {
	return New(msg, ErrorKindValidation)
}

// NewNotFoundError creates a not found error
func NewNotFoundError(msg string) Error {
	return New(msg, ErrorKindNotFound)
}

// NewInternalError creates an internal error
func NewInternalError(msg string) Error {
	return New(msg, ErrorKindInternal)
}

// NewConflictError creates a conflict error
func NewConflictError(msg string) Error {
	return New(msg, ErrorKindConflict)
}

// NewUnauthorizedError creates an unauthorized error
func NewUnauthorizedError(msg string) Error {
	return New(msg, ErrorKindUnauthorized)
}

// NewForbiddenError creates a forbidden error
func NewForbiddenError(msg string) Error {
	return New(msg, ErrorKindForbidden)
}

// Kind returns the error kind
func (e *CustomError) Kind() ErrorKind {
	return e.kind
}

// Error returns the error message
func (e *CustomError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %v", e.msg, e.cause)
	}
	return e.msg
}

// Cause returns the underlying cause
func (e *CustomError) Cause() error {
	return e.cause
}

// IsNotFound checks if the error is a not found error
func IsNotFound(err error) bool {
	if e, ok := err.(Error); ok {
		return e.Kind() == ErrorKindNotFound
	}
	return false
}

// IsValidation checks if the error is a validation error
func IsValidation(err error) bool {
	if e, ok := err.(Error); ok {
		return e.Kind() == ErrorKindValidation
	}
	return false
}

// IsInternal checks if the error is an internal error
func IsInternal(err error) bool {
	if e, ok := err.(Error); ok {
		return e.Kind() == ErrorKindInternal
	}
	return false
}
