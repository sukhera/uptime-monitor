package container

import (
	"context"
	"fmt"
	"sync"

	"github.com/sukhera/uptime-monitor/internal/application/handlers"
	"github.com/sukhera/uptime-monitor/internal/application/middleware"
	"github.com/sukhera/uptime-monitor/internal/application/routes"
	"github.com/sukhera/uptime-monitor/internal/checker"
	"github.com/sukhera/uptime-monitor/internal/domain/service"
	"github.com/sukhera/uptime-monitor/internal/infrastructure/database"
	mongodb "github.com/sukhera/uptime-monitor/internal/infrastructure/database/mongo"
	"github.com/sukhera/uptime-monitor/internal/server"
	"github.com/sukhera/uptime-monitor/internal/shared/config"
	"github.com/sukhera/uptime-monitor/internal/shared/logger"
)

// ContainerOption is a function that configures the container
type ContainerOption func(*Container) error

// Container manages application dependencies with proper interfaces
type Container struct {
	mu       sync.RWMutex
	services map[string]interface{}
	config   *config.Config
	logger   logger.Logger
}

// New creates a new dependency injection container
func New(cfg *config.Config, opts ...ContainerOption) (*Container, error) {
	container := &Container{
		services: make(map[string]interface{}),
		config:   cfg,
		logger:   logger.Get(),
	}

	// Apply all options
	for _, opt := range opts {
		if err := opt(container); err != nil {
			return nil, fmt.Errorf("failed to apply container option: %w", err)
		}
	}

	return container, nil
}

// WithDatabase adds a database service to the container
func WithDatabase(db database.Interface) ContainerOption {
	return func(c *Container) error {
		c.Register("database", db)
		return nil
	}
}

// WithServiceRepository adds a service repository to the container
func WithServiceRepository(repo service.Repository) ContainerOption {
	return func(c *Container) error {
		c.Register("service_repository", repo)
		return nil
	}
}

// WithStatusHandler adds a status handler to the container
func WithStatusHandler(handler *handlers.StatusHandler) ContainerOption {
	return func(c *Container) error {
		c.Register("status_handler", handler)
		return nil
	}
}

// WithCheckerService adds a checker service to the container
func WithCheckerService(svc checker.ServiceInterface) ContainerOption {
	return func(c *Container) error {
		c.Register("checker", svc)
		return nil
	}
}

// WithHTTPServer adds an HTTP server to the container
func WithHTTPServer(srv server.Interface) ContainerOption {
	return func(c *Container) error {
		c.Register("http_server", srv)
		return nil
	}
}

// Register registers a service with the container
func (c *Container) Register(name string, service interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.services[name] = service
}

// Get retrieves a service from the container
func (c *Container) Get(name string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	service, exists := c.services[name]
	return service, exists
}

// GetDatabase returns the database service
func (c *Container) GetDatabase() (database.Interface, error) {
	if db, exists := c.Get("database"); exists {
		return db.(database.Interface), nil
	}

	// Create new database connection using functional options
	db, err := mongodb.NewConnection(c.config.Database.URI, c.config.Database.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to create database connection: %w", err)
	}

	c.Register("database", db)
	return db, nil
}

// GetServiceRepository returns the service repository
func (c *Container) GetServiceRepository() (service.Repository, error) {
	if repo, exists := c.Get("service_repository"); exists {
		return repo.(service.Repository), nil
	}

	// Get database dependency
	db, err := c.GetDatabase()
	if err != nil {
		return nil, fmt.Errorf("failed to get database: %w", err)
	}

	// Create service repository using functional options
	mongoDB, ok := db.(*mongodb.Database)
	if !ok {
		return nil, fmt.Errorf("database is not MongoDB implementation")
	}

	repo := mongodb.NewServiceRepository(mongoDB)
	c.Register("service_repository", repo)
	return repo, nil
}

// GetStatusHandler returns the status handler
func (c *Container) GetStatusHandler() (*handlers.StatusHandler, error) {
	if handler, exists := c.Get("status_handler"); exists {
		return handler.(*handlers.StatusHandler), nil
	}

	// Get database dependency
	db, err := c.GetDatabase()
	if err != nil {
		return nil, fmt.Errorf("failed to get database: %w", err)
	}

	// Create status handler using functional options
	mongoDB, ok := db.(*mongodb.Database)
	if !ok {
		return nil, fmt.Errorf("database is not MongoDB implementation")
	}

	handler := handlers.NewStatusHandler(mongoDB)
	c.Register("status_handler", handler)
	return handler, nil
}

// GetCheckerService returns the checker service
func (c *Container) GetCheckerService() (checker.ServiceInterface, error) {
	if service, exists := c.Get("checker"); exists {
		return service.(checker.ServiceInterface), nil
	}

	// Get database dependency
	db, err := c.GetDatabase()
	if err != nil {
		return nil, fmt.Errorf("failed to get database: %w", err)
	}

	// Create checker service with functional options
	checkerService := checker.NewService(db,
		checker.WithTimeout(c.config.Database.Timeout),
	)

	c.Register("checker", checkerService)
	return checkerService, nil
}

// GetHTTPServer returns the HTTP server
func (c *Container) GetHTTPServer() (server.Interface, error) {
	if srv, exists := c.Get("http_server"); exists {
		return srv.(server.Interface), nil
	}

	// Get status handler
	statusHandler, err := c.GetStatusHandler()
	if err != nil {
		return nil, fmt.Errorf("failed to get status handler: %w", err)
	}

	// Setup routes
	router := routes.SetupRoutes(statusHandler)

	// Apply middleware
	corsMiddleware := middleware.NewCORS()
	handler := corsMiddleware.Handler(router)

	// Create server using functional options
	srv := server.New(handler, c.config)
	c.Register("http_server", srv)
	return srv, nil
}

// GetConfig returns the configuration
func (c *Container) GetConfig() *config.Config {
	return c.config
}

// GetLogger returns the logger
func (c *Container) GetLogger() logger.Logger {
	return c.logger
}

// Shutdown gracefully shuts down all services
func (c *Container) Shutdown(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var lastErr error
	for name, service := range c.services {
		if closer, ok := service.(interface{ Close() error }); ok {
			if err := closer.Close(); err != nil {
				lastErr = err
				c.logger.Error(ctx, "Failed to close service", err, logger.Fields{"service": name})
			}
		}
		delete(c.services, name)
	}

	return lastErr
}

// MustGetDatabase returns the database service or panics
func (c *Container) MustGetDatabase() database.Interface {
	db, err := c.GetDatabase()
	if err != nil {
		panic(fmt.Sprintf("failed to get database: %v", err))
	}
	return db
}

// MustGetServiceRepository returns the service repository or panics
func (c *Container) MustGetServiceRepository() service.Repository {
	repo, err := c.GetServiceRepository()
	if err != nil {
		panic(fmt.Sprintf("failed to get service repository: %v", err))
	}
	return repo
}

// MustGetStatusHandler returns the status handler or panics
func (c *Container) MustGetStatusHandler() *handlers.StatusHandler {
	handler, err := c.GetStatusHandler()
	if err != nil {
		panic(fmt.Sprintf("failed to get status handler: %v", err))
	}
	return handler
}

// MustGetCheckerService returns the checker service or panics
func (c *Container) MustGetCheckerService() checker.ServiceInterface {
	svc, err := c.GetCheckerService()
	if err != nil {
		panic(fmt.Sprintf("failed to get checker service: %v", err))
	}
	return svc
}

// MustGetHTTPServer returns the HTTP server or panics
func (c *Container) MustGetHTTPServer() server.Interface {
	srv, err := c.GetHTTPServer()
	if err != nil {
		panic(fmt.Sprintf("failed to get HTTP server: %v", err))
	}
	return srv
}
