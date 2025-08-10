package server

import (
	"context"
	"net/http"
)

// Interface defines the contract for HTTP server operations
type Interface interface {
	// Start starts the server
	Start() error

	// StartWithContext starts the server with context for graceful shutdown
	StartWithContext(ctx context.Context) error

	// Stop gracefully stops the server
	Stop(ctx context.Context) error

	// GetAddr returns the server address
	GetAddr() string

	// GetHandler returns the HTTP handler
	GetHandler() http.Handler
}

// Middleware defines the contract for HTTP middleware
type Middleware interface {
	Handler(next http.Handler) http.Handler
}

// Router defines the contract for HTTP routing
type Router interface {
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
	Handle(pattern string, handler http.Handler)
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}
