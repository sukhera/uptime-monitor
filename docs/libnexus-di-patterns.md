# Libnexus-Style Dependency Injection Patterns

This document explains how to use the libnexus-style dependency injection patterns implemented in our project.

## Overview

Our dependency injection container follows the patterns used in the libnexus codebase, providing:

- **Functional Options Pattern** for flexible configuration
- **Interface-First Design** for better testability
- **Safe and Must Methods** for different use cases
- **Graceful Shutdown** with proper resource cleanup
- **Comprehensive Error Handling** with custom error types

## Basic Usage

### 1. Creating a Container

```go
import (
    "context"
    "log"
    
    "github.com/sukhera/uptime-monitor/internal/container"
    "github.com/sukhera/uptime-monitor/internal/shared/config"
)

func main() {
    ctx := context.Background()
    
    // Create configuration using functional options
    cfg := config.New(
        config.FromEnvironment(),
        config.WithDatabase("mongodb://localhost:27017", "uptime_monitor", 30),
    )
    
    // Create container (libnexus style)
    container, err := container.New(cfg)
    if err != nil {
        log.Fatal("Failed to create container:", err)
    }
    
    // Use the container...
}
```

### 2. Getting Services

#### Safe Methods (with error handling)
```go
// Get database with error handling
database, err := container.GetDatabase()
if err != nil {
    log.Fatal("Failed to get database:", err)
}

// Get service repository with error handling
repo, err := container.GetServiceRepository()
if err != nil {
    log.Fatal("Failed to get repository:", err)
}

// Get HTTP server with error handling
httpServer, err := container.GetHTTPServer()
if err != nil {
    log.Fatal("Failed to get HTTP server:", err)
}
```

#### Must Methods (panics on error)
```go
// Must get methods - useful for initialization
database := container.MustGetDatabase()
repo := container.MustGetServiceRepository()
httpServer := container.MustGetHTTPServer()
```

#### Direct Access
```go
// Get any service by name
if service, exists := container.Get("database"); exists {
    if database, ok := service.(database.Interface); ok {
        // Use database
    }
}
```

### 3. Graceful Shutdown

```go
// Graceful shutdown with context
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

if err := container.Shutdown(ctx); err != nil {
    log.Printf("Error during shutdown: %v", err)
}
```

## Advanced Patterns

### 1. Container with Custom Services (Testing)

```go
// Create container with custom services for testing
container, err := container.New(cfg,
    container.WithDatabase(mockDatabase),
    container.WithServiceRepository(mockRepository),
    container.WithStatusHandler(mockHandler),
)
if err != nil {
    log.Fatal("Failed to create container:", err)
}
```

### 2. Functional Options Pattern

```go
// Create configuration with multiple options
cfg := config.New(
    config.FromEnvironment(),
    config.WithDatabase("mongodb://localhost:27017", "uptime_monitor", 30),
    config.WithServerPort("8080"),
    config.WithLogging("info", true),
)

// Create container with options
container, err := container.New(cfg)
if err != nil {
    log.Fatal("Failed to create container:", err)
}
```

### 3. Error Handling Patterns

```go
// Pattern 1: Safe get with error handling
database, err := container.GetDatabase()
if err != nil {
    log.Printf("Database error: %v", err)
    return
}

// Pattern 2: Must get (panics on error) - useful for initialization
defer func() {
    if r := recover(); r != nil {
        log.Printf("Must get panic: %v", r)
    }
}()

// This would panic if there's an error
// repo := container.MustGetServiceRepository()

// Pattern 3: Graceful shutdown with error handling
if err := container.Shutdown(ctx); err != nil {
    log.Printf("Shutdown error: %v", err)
}
```

### 4. Advanced Usage

```go
// Get configuration
config := container.GetConfig()

// Get logger
logger := container.GetLogger()

// Get all services
services := []string{"database", "service_repository", "status_handler", "http_server"}
for _, serviceName := range services {
    if service, exists := container.Get(serviceName); exists {
        log.Printf("%s: %T", serviceName, service)
    }
}
```

## Real-World Examples

### 1. API Server Main

```go
package main

import (
    "context"
    "log"
    
    "github.com/sukhera/uptime-monitor/internal/container"
    "github.com/sukhera/uptime-monitor/internal/shared/config"
    "github.com/sukhera/uptime-monitor/internal/shared/logger"
)

func main() {
    ctx := context.Background()
    
    // Load configuration using functional options pattern
    cfg := config.New(
        config.FromEnvironment(),
    )
    
    // Validate configuration
    if err := cfg.Validate(); err != nil {
        log.Fatalf("Invalid configuration: %v", err)
    }
    
    // Initialize logger
    log := logger.Get()
    log.Info(ctx, "Starting API server", logger.Fields{
        "port":     cfg.Server.Port,
        "database": cfg.Database.URI,
    })
    
    // Create dependency injection container using libnexus-style patterns
    container, err := container.New(cfg)
    if err != nil {
        log.Fatal(ctx, "Failed to create container", err, logger.Fields{})
    }
    
    // Get HTTP server from container
    httpServer, err := container.GetHTTPServer()
    if err != nil {
        log.Fatal(ctx, "Failed to get HTTP server", err, logger.Fields{})
    }
    
    // Start server
    log.Info(ctx, "Starting HTTP server", logger.Fields{
        "address": httpServer.GetAddr(),
    })
    
    if err := httpServer.Start(); err != nil {
        log.Fatal(ctx, "Server failed to start", err, logger.Fields{})
    }
}
```

### 2. Status Checker Main

```go
package main

import (
    "context"
    "time"
    
    "github.com/go-co-op/gocron"
    "github.com/sukhera/uptime-monitor/internal/checker"
    "github.com/sukhera/uptime-monitor/internal/container"
    "github.com/sukhera/uptime-monitor/internal/shared/config"
    "github.com/sukhera/uptime-monitor/internal/shared/logger"
)

func main() {
    ctx := context.Background()
    
    // Load configuration
    cfg := config.New(
        config.FromEnvironment(),
    )
    
    // Validate configuration
    if err := cfg.Validate(); err != nil {
        log.Fatal(ctx, "Invalid configuration", err, logger.Fields{})
    }
    
    // Initialize logger
    log := logger.Get()
    
    // Create dependency injection container
    container, err := container.New(cfg)
    if err != nil {
        log.Fatal(ctx, "Failed to create container", err, logger.Fields{})
    }
    
    // Get checker service from container
    checkerService, err := container.GetCheckerService()
    if err != nil {
        log.Fatal(ctx, "Failed to get checker service", err, logger.Fields{})
    }
    
    // Setup observers for health check events
    subject := checker.NewHealthCheckSubject()
    
    // Add logging observer
    loggingObserver := checker.NewLoggingObserver(log)
    subject.Attach(loggingObserver)
    
    // Start scheduler
    scheduler := gocron.NewScheduler(time.UTC)
    
    _, err = scheduler.Every(cfg.Checker.Interval).Do(func() {
        runHealthChecks(ctx, checkerService, subject, log)
    })
    if err != nil {
        log.Fatal(ctx, "Failed to schedule health checks", err, logger.Fields{})
    }
    
    log.Info(ctx, "Status checker started successfully", logger.Fields{})
    scheduler.StartBlocking()
}
```

## Container Options

### Available Options

```go
// WithDatabase adds a database service to the container
container.WithDatabase(db database.Interface)

// WithServiceRepository adds a service repository to the container
container.WithServiceRepository(repo service.Repository)

// WithStatusHandler adds a status handler to the container
container.WithStatusHandler(handler *handlers.StatusHandler)

// WithCheckerService adds a checker service to the container
container.WithCheckerService(svc checker.ServiceInterface)

// WithHTTPServer adds an HTTP server to the container
container.WithHTTPServer(srv server.Interface)
```

### Using Options

```go
// Create container with custom services
container, err := container.New(cfg,
    container.WithDatabase(customDatabase),
    container.WithServiceRepository(customRepository),
    container.WithStatusHandler(customHandler),
)
```

## Benefits

### 1. Testability
- Easy to mock interfaces using mockery
- Clear separation between domain and infrastructure
- Functional options allow easy testing configurations

### 2. Maintainability
- Clear separation of concerns
- Interface-first design
- Consistent patterns throughout the codebase

### 3. Flexibility
- Easy to swap implementations
- Functional options for configuration
- Custom service injection for testing

### 4. Type Safety
- Strong typing throughout
- Interface contracts ensure consistency
- Compile-time error checking

### 5. Error Handling
- Comprehensive error handling with custom types
- Safe and must methods for different use cases
- Graceful shutdown with proper resource cleanup

## Best Practices

### 1. Always Use Interfaces
```go
// Good: Use interfaces
func NewService(db database.Interface) *Service {
    return &Service{db: db}
}

// Avoid: Use concrete types
func NewService(db *mongodb.Database) *Service {
    return &Service{db: db}
}
```

### 2. Use Functional Options
```go
// Good: Functional options pattern
container, err := container.New(cfg,
    container.WithDatabase(db),
    container.WithServiceRepository(repo),
)

// Avoid: Manual service creation
container := &Container{}
container.Register("database", db)
container.Register("repository", repo)
```

### 3. Handle Errors Properly
```go
// Good: Proper error handling
database, err := container.GetDatabase()
if err != nil {
    return fmt.Errorf("failed to get database: %w", err)
}

// Avoid: Ignoring errors
database := container.MustGetDatabase() // Panics on error
```

### 4. Use Graceful Shutdown
```go
// Good: Graceful shutdown with timeout
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

if err := container.Shutdown(ctx); err != nil {
    log.Printf("Error during shutdown: %v", err)
}
```

## Migration Guide

### From Old Container Usage

**Before:**
```go
// Old way
cfg := config.Load()
db, err := mongodb.NewConnection(cfg.Database.URI, cfg.Database.Name)
if err != nil {
    log.Fatalf("Failed to connect to database: %v", err)
}
defer db.Close()

statusHandler := handlers.NewStatusHandler(db)
router := routes.SetupRoutes(statusHandler)
srv := server.New(router, cfg)
```

**After:**
```go
// New libnexus-style way
cfg := config.New(config.FromEnvironment())
if err := cfg.Validate(); err != nil {
    log.Fatalf("Invalid configuration: %v", err)
}

container, err := container.New(cfg)
if err != nil {
    log.Fatal("Failed to create container:", err)
}

httpServer, err := container.GetHTTPServer()
if err != nil {
    log.Fatal("Failed to get HTTP server:", err)
}

if err := httpServer.Start(); err != nil {
    log.Fatal("Server failed to start:", err)
}
```

## Conclusion

The libnexus-style dependency injection patterns provide:

1. **Clean Architecture** with proper separation of concerns
2. **Interface-First Design** for better testability
3. **Functional Options Pattern** for flexible configuration
4. **Comprehensive Error Handling** with custom error types
5. **Graceful Shutdown** with proper resource cleanup
6. **Type Safety** with strong typing throughout
7. **Maintainability** with clear patterns and consistent design
8. **Testability** with easy mocking and clear interfaces

These patterns make the codebase more maintainable, testable, and scalable while following Go best practices and the proven patterns from the libnexus codebase. 