# Functional Options Pattern Implementation

## Why Functional Options Pattern?

The functional options pattern was implemented to address several issues with the previous configuration approach:

### Problems with Previous Approach

1. **Inconsistent Configuration Structure**: The original `config.go` defined a `Config` struct with nested configs, but the status-checker was trying to access fields like `CheckInterval`, `MongoURI`, and `DBName` that didn't exist in the current structure.

2. **Mixed Configuration Sources**: The status-checker used both environment variables and command-line flags with manual override logic that was error-prone and hard to maintain.

3. **Tight Coupling**: Configuration was tightly coupled to environment variables, making testing and flexibility difficult.

4. **No Validation**: The original `Load()` function didn't validate configuration, leading to runtime errors.

5. **Poor Testability**: It was difficult to test different configuration scenarios without setting environment variables.

## Benefits of Functional Options Pattern

### 1. **Flexibility and Composability**

```go
// Before: Hard to customize
cfg := config.Load()

// After: Easy to customize with options
cfg := config.New(
    config.WithServerPort("9090"),
    config.WithDatabase("mongodb://custom:27017", "custom_db", 15*time.Second),
    config.WithLogging("debug", true),
    config.WithCheckerInterval(5*time.Minute),
)
```

### 2. **Better Testing**

```go
// Before: Had to set environment variables
os.Setenv("MONGO_URI", "mongodb://test:27017")
os.Setenv("DB_NAME", "test_db")
cfg := config.Load()

// After: Clean, isolated tests
cfg := config.New(
    config.WithDatabase("mongodb://test:27017", "test_db", 5*time.Second),
)
```

### 3. **Clear Configuration Intent**

```go
// Before: Magic values scattered throughout
checkerService := checker.NewService(db)

// After: Explicit configuration
checkerService := checker.NewService(db,
    checker.WithTimeout(30*time.Second),
    checker.WithHTTPClient(customClient),
)
```

### 4. **Backward Compatibility**

The old `Load()` function still works:

```go
// Still works for existing code
cfg := config.Load()
```

### 5. **Validation and Error Handling**

```go
cfg := config.New(
    config.WithServerPort(""), // This will be caught by validation
)
if err := cfg.Validate(); err != nil {
    log.Fatalf("Invalid configuration: %v", err)
}
```

## Implementation Details

### Configuration Options

```go
type Option func(*Config)

func WithServerPort(port string) Option {
    return func(c *Config) {
        c.Server.Port = port
    }
}

func WithDatabase(uri, name string, timeout time.Duration) Option {
    return func(c *Config) {
        c.Database.URI = uri
        c.Database.Name = name
        c.Database.Timeout = timeout
    }
}
```

### Service Options

```go
type ServiceOption func(*Service)

func WithHTTPClient(client HTTPClient) ServiceOption {
    return func(s *Service) {
        s.client = client
    }
}

func WithTimeout(timeout time.Duration) ServiceOption {
    return func(s *Service) {
        if s.client == nil {
            s.client = &http.Client{Timeout: timeout}
        } else {
            if httpClient, ok := s.client.(*http.Client); ok {
                httpClient.Timeout = timeout
            }
        }
    }
}
```

## Usage Examples

### Basic Usage

```go
// Default configuration
cfg := config.New()

// Environment-based configuration
cfg := config.New(config.FromEnvironment())

// Custom configuration
cfg := config.New(
    config.WithServerPort("9090"),
    config.WithDatabase("mongodb://localhost:27017", "myapp", 10*time.Second),
)
```

### Command-line Overrides

```go
// Start with environment config
cfg := config.New(config.FromEnvironment())

// Apply command-line overrides
var options []config.Option
if *intervalMinutes > 0 {
    options = append(options, config.WithCheckerInterval(time.Duration(*intervalMinutes)*time.Minute))
}
if *mongoURI != "" {
    options = append(options, config.WithDatabase(*mongoURI, cfg.Database.Name, cfg.Database.Timeout))
}

// Apply all options
if len(options) > 0 {
    cfg = config.New(append([]config.Option{config.FromEnvironment()}, options...)...)
}
```

### Testing

```go
func TestServiceWithCustomClient(t *testing.T) {
    mockClient := &MockHTTPClient{}
    service := checker.NewService(db, checker.WithHTTPClient(mockClient))
    
    // Test with mock client
    // ...
}
```

## Migration Guide

### For Existing Code

1. **Replace direct field access**:
   ```go
   // Before
   cfg.MongoURI
   cfg.DBName
   cfg.CheckInterval
   
   // After
   cfg.Database.URI
   cfg.Database.Name
   cfg.Checker.Interval
   ```

2. **Update service creation**:
   ```go
   // Before
   service := checker.NewService(db)
   
   // After (still works)
   service := checker.NewService(db)
   
   // Or with options
   service := checker.NewService(db, checker.WithTimeout(30*time.Second))
   ```

3. **Add validation**:
   ```go
   cfg := config.Load()
   if err := cfg.Validate(); err != nil {
       log.Fatalf("Invalid configuration: %v", err)
   }
   ```

## Best Practices

1. **Always validate configuration** before using it
2. **Use descriptive option names** that clearly indicate their purpose
3. **Provide sensible defaults** in the `New()` function
4. **Keep options focused** - each option should configure one aspect
5. **Document options** with clear examples
6. **Maintain backward compatibility** when possible

## Future Enhancements

1. **Configuration Profiles**: Pre-defined configurations for different environments
2. **Configuration Validation**: More comprehensive validation rules
3. **Configuration Hot-reload**: Dynamic configuration updates
4. **Configuration Encryption**: Secure storage of sensitive configuration
5. **Configuration Metrics**: Monitoring configuration usage and changes 