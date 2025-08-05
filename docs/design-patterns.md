# Design Patterns Implementation Guide

## Overview

This document outlines the design patterns implemented in the status-checker and broader codebase to improve maintainability, testability, and extensibility.

## 🎯 **Implemented Patterns**

### 1. **Functional Options Pattern** ✅ **IMPLEMENTED**

**Purpose**: Flexible configuration management with backward compatibility.

**Implementation**:
```go
// Configuration with options
cfg := config.New(
    config.WithServerPort("9090"),
    config.WithDatabase("mongodb://custom:27017", "custom_db", 15*time.Second),
    config.WithLogging("debug", true),
    config.WithCheckerInterval(5*time.Minute),
)
```

**Benefits**:
- ✅ Flexible configuration composition
- ✅ Better testability (no environment variables needed)
- ✅ Clear configuration intent
- ✅ Backward compatibility maintained

### 2. **Dependency Injection Container** ✅ **IMPLEMENTED**

**Purpose**: Manage service dependencies and improve testability.

**Implementation**:
```go
// Initialize container
container := container.New(cfg)

// Get services from container
checkerService, err := container.GetCheckerService()
if err != nil {
    log.Fatal(ctx, "Failed to get checker service", err, logger.Fields{})
}
```

**Benefits**:
- ✅ Loose coupling between components
- ✅ Easy service replacement for testing
- ✅ Centralized dependency management
- ✅ Automatic resource cleanup

### 3. **Structured Logging with Context** ✅ **IMPLEMENTED**

**Purpose**: Improve observability and debugging capabilities.

**Implementation**:
```go
// Initialize structured logging
logger.Init(logger.INFO)
log := logger.Get()

// Log with context and fields
log.Info(ctx, "Health check completed", logger.Fields{
    "service_name": "api",
    "status":       "operational",
    "latency_ms":   150,
})
```

**Benefits**:
- ✅ Consistent log format across the application
- ✅ Context-aware logging (request IDs, user IDs)
- ✅ Structured data for log aggregation
- ✅ Configurable log levels

### 4. **Command Pattern** ✅ **IMPLEMENTED**

**Purpose**: Encapsulate health check operations and make them more modular.

**Implementation**:
```go
// Create health check commands
invoker := NewHealthCheckInvoker()
for _, service := range services {
    command := NewHTTPHealthCheckCommand(service, client)
    invoker.AddCommand(command)
}

// Execute commands concurrently
statusLogs := invoker.ExecuteAll(ctx)
```

**Benefits**:
- ✅ Encapsulated health check logic
- ✅ Easy to add new health check types
- ✅ Concurrent execution support
- ✅ Better testability

### 5. **Observer Pattern** ✅ **IMPLEMENTED**

**Purpose**: Decouple health check events from their handlers.

**Implementation**:
```go
// Setup observers
subject := NewHealthCheckSubject()
subject.Attach(NewLoggingObserver(logger))
subject.Attach(NewMetricsObserver())
subject.Attach(NewAlertingObserver(5000))

// Notify observers of events
subject.Notify(ctx, event)
```

**Benefits**:
- ✅ Loose coupling between events and handlers
- ✅ Easy to add new event handlers
- ✅ Asynchronous event processing
- ✅ Extensible alerting system

## 🔧 **Pattern Integration**

### Status Checker Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Configuration │    │   DI Container  │    │  Command Pattern│
│   (Functional   │───▶│   (Dependency   │───▶│  (Health Check  │
│    Options)     │    │   Injection)    │    │   Commands)     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │                        │
                                ▼                        ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  Observer Pattern│    │ Structured Logging│    │  Database Layer │
│  (Event Handling)│    │  (Context + Fields)│    │  (MongoDB)     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### Flow Diagram

```
1. Configuration Loaded (Functional Options)
   ↓
2. DI Container Initialized
   ↓
3. Services Retrieved from Container
   ↓
4. Observers Attached (Logging, Metrics, Alerting)
   ↓
5. Health Check Commands Created
   ↓
6. Commands Executed Concurrently
   ↓
7. Results Stored + Observers Notified
   ↓
8. Alerts Processed Asynchronously
```

## 📊 **Benefits Achieved**

### 1. **Maintainability**
- Clear separation of concerns
- Modular components
- Easy to understand and modify

### 2. **Testability**
- Dependency injection enables easy mocking
- Command pattern allows isolated testing
- Observer pattern decouples event handling

### 3. **Extensibility**
- Easy to add new health check types
- Simple to add new event observers
- Flexible configuration options

### 4. **Observability**
- Structured logging with context
- Metrics collection
- Alerting system

### 5. **Performance**
- Concurrent health check execution
- Asynchronous event processing
- Efficient resource management

## 🧪 **Testing Examples**

### Functional Options Testing
```go
func TestConfigWithOptions(t *testing.T) {
    cfg := config.New(
        config.WithServerPort("9090"),
        config.WithDatabase("mongodb://test:27017", "test_db", 5*time.Second),
    )
    
    assert.Equal(t, "9090", cfg.Server.Port)
    assert.Equal(t, "mongodb://test:27017", cfg.Database.URI)
}
```

### Command Pattern Testing
```go
func TestHealthCheckCommand(t *testing.T) {
    mockClient := &MockHTTPClient{}
    service := models.Service{Name: "test", URL: "http://test.com"}
    
    command := NewHTTPHealthCheckCommand(service, mockClient)
    result := command.Execute(context.Background())
    
    assert.Equal(t, "test", result.ServiceName)
}
```

### Observer Pattern Testing
```go
func TestObserverNotification(t *testing.T) {
    subject := NewHealthCheckSubject()
    mockObserver := &MockObserver{}
    subject.Attach(mockObserver)
    
    event := HealthCheckEvent{ServiceName: "test", Status: "down"}
    subject.Notify(context.Background(), event)
    
    assert.True(t, mockObserver.Notified)
}
```

## 🚀 **Future Enhancements**

### 1. **Strategy Pattern**
- Different health check strategies (HTTP, TCP, DNS, etc.)
- Pluggable health check algorithms

### 2. **Factory Pattern**
- Health check command factories
- Observer factories based on configuration

### 3. **Chain of Responsibility**
- Health check result processing pipeline
- Alert escalation chain

### 4. **Template Method Pattern**
- Base health check template with customizable steps
- Standardized health check workflow

### 5. **State Pattern**
- Service state management (operational, degraded, down)
- State transition logic

## 📈 **Performance Impact**

### Before Patterns
- Tight coupling between components
- Hard to test individual parts
- Manual configuration management
- Basic logging without structure

### After Patterns
- Loose coupling with DI container
- Easy unit testing with mocks
- Flexible configuration with options
- Structured logging with context
- Concurrent health check execution
- Asynchronous event processing

## 🔍 **Monitoring and Debugging**

### Structured Logs
```json
{
  "level": "INFO",
  "message": "Health check completed",
  "timestamp": "2025-08-04T12:30:00Z",
  "service_name": "api",
  "status": "operational",
  "latency_ms": 150,
  "status_code": 200,
  "request_id": "req-123"
}
```

### Metrics Collection
```json
{
  "global": {
    "total_checks": 150
  },
  "service_api": {
    "status": "operational",
    "latency_ms": 150,
    "status_code": 200,
    "timestamp": 1628087400
  }
}
```

### Alert Processing
```go
// Alert event structure
type AlertEvent struct {
    ServiceName string
    Status      string
    Latency     int64
    Severity    string
    Timestamp   time.Time
}
```

## 📚 **Best Practices**

### 1. **Configuration Management**
- Use functional options for flexible configuration
- Validate configuration early
- Provide sensible defaults

### 2. **Dependency Management**
- Use DI container for service management
- Keep dependencies minimal
- Implement proper cleanup

### 3. **Logging**
- Use structured logging with context
- Include relevant fields in log messages
- Configure appropriate log levels

### 4. **Testing**
- Mock dependencies for unit tests
- Test each pattern in isolation
- Use table-driven tests for configuration

### 5. **Error Handling**
- Provide meaningful error messages
- Log errors with context
- Implement proper error recovery

This implementation demonstrates how design patterns can significantly improve code quality, maintainability, and extensibility while maintaining backward compatibility and performance. 