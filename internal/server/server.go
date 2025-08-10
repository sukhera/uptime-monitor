package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/sukhera/uptime-monitor/internal/application/middleware"
	"github.com/sukhera/uptime-monitor/internal/shared/config"
	"github.com/sukhera/uptime-monitor/internal/shared/logger"
)

// Server represents the HTTP server
type Server struct {
	*http.Server
	config *config.Config
}

// New creates a new server instance
func New(handler http.Handler, cfg *config.Config) *Server {
	return &Server{
		Server: &http.Server{
			Addr:              ":" + cfg.Server.Port,
			Handler:           handler,
			ReadTimeout:       cfg.Server.ReadTimeout,
			ReadHeaderTimeout: cfg.Server.ReadHeaderTimeout,
			WriteTimeout:      cfg.Server.WriteTimeout,
			IdleTimeout:       cfg.Server.IdleTimeout,
		},
		config: cfg,
	}
}

// Start starts the server with graceful shutdown
func (s *Server) Start() error {
	return s.StartWithContext(context.Background())
}

// StartWithContext starts the server with context for graceful shutdown
func (s *Server) StartWithContext(ctx context.Context) error {
	log := logger.Get()
	
	// Create a channel to listen for errors coming from the listener
	serverErrors := make(chan error, 1)

	// Start the server in a goroutine
	go func() {
		log.Info(ctx, "Starting server", logger.Fields{"port": s.config.Server.Port})
		serverErrors <- s.ListenAndServe()
	}()

	// Blocking select waiting for either a server error or context cancellation
	select {
	case err := <-serverErrors:
		if err == http.ErrServerClosed {
			return nil
		}
		return fmt.Errorf("server error: %w", err)

	case <-ctx.Done():
		log.Info(ctx, "Context cancelled, starting graceful shutdown", logger.Fields{"reason": ctx.Err().Error()})

		// Give outstanding requests a deadline for completion
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Gracefully shutdown the server
		if err := s.Shutdown(shutdownCtx); err != nil {
			log.Error(ctx, "Could not stop server gracefully", err, nil)
			return err
		}

		log.Info(ctx, "Server stopped gracefully", nil)
		return nil
	}
}

// Stop gracefully stops the server
func (s *Server) Stop(ctx context.Context) error {
	return s.Shutdown(ctx)
}

// GetAddr returns the server address
func (s *Server) GetAddr() string {
	return s.Addr
}

// GetHandler returns the HTTP handler
func (s *Server) GetHandler() http.Handler {
	return s.Handler
}

// ApplyMiddleware applies middleware to the server handler
func (s *Server) ApplyMiddleware(middlewares ...middleware.Middleware) {
	s.Handler = middleware.Chain(s.Handler, middlewares...)
}
