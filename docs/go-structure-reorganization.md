# Go Structure Reorganization Plan

Based on analysis of libnexus repository patterns, here's the reorganization plan:

## Current Issues vs. libnexus Patterns

### Current Issues:
1. Mixed concerns in `cmd/` directory
2. Inconsistent naming conventions
3. Missing clear separation of layers
4. Scattered test organization
5. No clear domain separation

### libnexus Best Practices:
1. **Clean Package Organization**: Each package has single responsibility
2. **Interface-First Design**: Clear interfaces before implementations
3. **Functional Options Pattern**: Used for configuration
4. **Domain-Driven Structure**: Clear separation of concerns
5. **Comprehensive Testing**: Tests alongside source files
6. **Error Handling**: Custom error types
7. **Dependency Injection**: Clear interfaces and implementations

## Proposed New Structure

```
status_page_starter/
├── cmd/                          # Application entry points (like libnexus)
│   ├── api/                      # API server
│   │   └── main.go
│   ├── status-checker/           # Status checker service
│   │   └── main.go
│   └── web/                      # Web server (if needed)
│       └── main.go
├── internal/                     # Private application code
│   ├── domain/                   # Domain models and business logic
│   │   ├── service/              # Service domain
│   │   │   ├── entity.go         # Service entity
│   │   │   ├── repository.go     # Repository interface
│   │   │   ├── service.go        # Service interface
│   │   │   └── types.go          # Domain types
│   │   ├── incident/             # Incident domain
│   │   │   ├── entity.go
│   │   │   ├── repository.go
│   │   │   ├── service.go
│   │   │   └── types.go
│   │   └── healthcheck/          # Health check domain
│   │       ├── entity.go
│   │       ├── repository.go
│   │       ├── service.go
│   │       └── types.go
│   ├── infrastructure/           # External concerns (like libnexus)
│   │   ├── database/             # Database implementations
│   │   │   ├── mongo/
│   │   │   │   ├── connection.go
│   │   │   │   ├── repository.go
│   │   │   │   └── migrations.go
│   │   │   └── interfaces.go
│   │   ├── messaging/            # Message queue implementations
│   │   ├── cache/                # Cache implementations
│   │   └── external/             # External service clients
│   ├── application/              # Application services
│   │   ├── handlers/             # HTTP handlers
│   │   │   ├── status.go
│   │   │   ├── incidents.go
│   │   │   └── health.go
│   │   ├── middleware/           # HTTP middleware
│   │   │   ├── cors.go
│   │   │   ├── auth.go
│   │   │   ├── logging.go
│   │   │   └── chain.go
│   │   └── routes/               # Route definitions
│   │       └── routes.go
│   ├── server/                   # Server implementations
│   │   ├── http.go
│   │   └── grpc.go              # Future use
│   └── shared/                   # Shared utilities (like libnexus)
│       ├── config/               # Configuration
│       │   ├── config.go
│       │   └── options.go
│       ├── logger/               # Logging
│       │   └── logger.go
│       ├── errors/               # Error handling
│       │   └── errors.go
│       └── utils/                # Utility functions
│           └── utils.go
├── pkg/                          # Public packages (if any)
├── api/                          # API definitions
│   ├── openapi/
│   └── proto/                    # Protocol buffers (future)
├── deployments/                  # Deployment configurations
├── scripts/                      # Build and deployment scripts
├── docs/                         # Documentation
├── test/                         # Integration tests
│   ├── api/
│   ├── database/
│   └── e2e/
├── mocks/                        # Generated mocks
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## Migration Steps

### Phase 1: Create New Structure
1. Create new directory structure
2. Move existing code to new locations
3. Update imports and package declarations

### Phase 2: Refactor Code
1. Implement interface-first design
2. Add functional options pattern
3. Improve error handling
4. Add comprehensive testing

### Phase 3: Update Dependencies
1. Update build scripts
2. Update documentation
3. Update deployment configurations

## Key Improvements from libnexus Patterns

### 1. Interface-First Design
```go
// Before: Direct implementation
type Service struct {
    db *mongo.Database
}

// After: Interface-first (like libnexus)
type ServiceRepository interface {
    GetService(ctx context.Context, id string) (*Service, error)
    CreateService(ctx context.Context, service *Service) error
    UpdateService(ctx context.Context, service *Service) error
    DeleteService(ctx context.Context, id string) error
}
```

### 2. Functional Options Pattern
```go
// Like libnexus config pattern
type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    Logging  LoggingConfig
}

type Option func(*Config)

func WithServerPort(port string) Option {
    return func(c *Config) {
        c.Server.Port = port
    }
}

func New(opts ...Option) *Config {
    cfg := &Config{}
    for _, opt := range opts {
        opt(cfg)
    }
    return cfg
}
```

### 3. Error Handling
```go
// Like libnexus errors package
type ErrorKind string

const (
    ErrorKindValidation ErrorKind = "validation"
    ErrorKindNotFound   ErrorKind = "not_found"
    ErrorKindInternal   ErrorKind = "internal"
)

type Error interface {
    Kind() ErrorKind
    Error() string
    Cause() error
}
```

### 4. Dependency Injection
```go
// Like libnexus container pattern
type Container struct {
    config     *config.Config
    logger     *logger.Logger
    database   database.Database
    services   map[string]interface{}
}

func NewContainer(cfg *config.Config) *Container {
    return &Container{
        config:   cfg,
        services: make(map[string]interface{}),
    }
}
```

## Benefits

1. **Maintainability**: Clear separation of concerns
2. **Testability**: Easy to mock and test individual components
3. **Scalability**: Easy to add new domains and features
4. **Readability**: Clear structure makes code easier to understand
5. **Standards Compliance**: Follows Go community best practices
6. **Consistency**: Consistent patterns throughout the codebase 