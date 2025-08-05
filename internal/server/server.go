package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sukhera/uptime-monitor/internal/application/middleware"
	"github.com/sukhera/uptime-monitor/internal/shared/config"
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
			Addr:         ":" + cfg.Server.Port,
			Handler:      handler,
			ReadTimeout:  cfg.Server.ReadTimeout,
			WriteTimeout: cfg.Server.WriteTimeout,
			IdleTimeout:  cfg.Server.IdleTimeout,
		},
		config: cfg,
	}
}

// Start starts the server with graceful shutdown
func (s *Server) Start() error {
	// Create a channel to listen for errors coming from the listener
	serverErrors := make(chan error, 1)

	// Start the server in a goroutine
	go func() {
		log.Printf("Starting server on port %s", s.config.Server.Port)
		serverErrors <- s.ListenAndServe()
	}()

	// Create a channel to listen for an interrupt or terminate signal from the OS
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Blocking select waiting for either a server error or a shutdown signal
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.Printf("Received signal %v, starting graceful shutdown", sig)

		// Give outstanding requests a deadline for completion
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Gracefully shutdown the server
		if err := s.Shutdown(ctx); err != nil {
			log.Printf("Could not stop server gracefully: %v", err)
			return err
		}

		log.Println("Server stopped gracefully")
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
