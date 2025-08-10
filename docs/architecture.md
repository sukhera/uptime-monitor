# Architecture Overview

## System Components

The Status Page application consists of several key components:

### 1. API Service (`cmd/api/`)
- RESTful API server built with Go
- Provides status data endpoints
- Handles CORS for web frontend
- Connects to MongoDB for data persistence

### 2. Status Checker (`cmd/status-checker/`)
- Background service for health monitoring
- Runs scheduled checks every 2 minutes
- Stores results in MongoDB
- Configurable service definitions

### 3. Web Frontend (`web/`)
- Modern responsive dashboard
- Auto-refreshing status display
- Modular JavaScript architecture
- Served via Nginx

### 4. Database (MongoDB)
- Stores service configurations
- Maintains status history logs
- Supports horizontal scaling

## Data Flow

```
[Services] -> [Status Checker] -> [MongoDB] -> [API] -> [Web Frontend] -> [Users]
```

1. Status Checker polls configured services
2. Results stored in MongoDB
3. API serves latest status data
4. Web frontend displays real-time status
5. Auto-refresh keeps data current

## Directory Structure

```
status_page_starter/
├── cmd/                     # Command-line applications
│   ├── api/                # API server entry point
│   │   └── main.go         # API server main
│   ├── status-checker/     # Health checker service
│   │   ├── main.go         # Status checker main
│   │   └── main_test.go    # Status checker tests
│   ├── api.go              # API command
│   ├── checker.go          # Checker command
│   ├── root.go             # Root command
│   └── web.go              # Web server command
├── configs/                 # Configuration files
│   ├── dev/                # Development configurations
│   ├── docker/             # Docker configurations
│   ├── env/                # Environment configurations
│   └── nginx/              # Nginx configurations
├── data/                    # Data files and seeds
│   ├── seed.js             # Database seed data
│   └── status.json         # Status data
├── deployments/             # Deployment configurations
│   └── kubernetes/         # Kubernetes manifests
│       ├── mongodb.yaml    # MongoDB deployment
│       └── namespace.yaml  # Namespace configuration
├── docs/                    # Documentation
│   ├── api.md              # API documentation
│   ├── architecture.md     # Architecture overview
│   ├── best-practices.md   # Best practices guide
│   ├── configuration.md    # Configuration guide
│   ├── design-patterns.md  # Design patterns
│   ├── functional-options-pattern.md # Functional options
│   ├── go-reorganization-progress.md # Go reorganization
│   ├── go-structure-proposal.md # Structure proposal
│   ├── go-structure-reorganization.md # Reorganization
│   ├── libnexus-di-patterns.md # DI patterns
│   ├── mockery-golangci-implementation.md # Mockery setup
│   ├── mockery-implementation-example.md # Mockery examples
│   └── mockery-setup.md    # Mockery configuration
├── examples/                # Code examples
│   └── functional-options-demo.go # Functional options demo
├── feature-request/         # Feature request documentation
│   ├── golang-improvements.md # Go improvements
│   ├── integration-tests-ci-cd.md # Integration tests
│   ├── makefile-infrastructure.md # Makefile infrastructure
│   ├── migration.md        # Migration guide
│   ├── performance-and-scalability-improvements.md # Performance
│   ├── readme-cleanup.md   # README cleanup
│   ├── readme-improvements.md # README improvements
│   └── software-architect-recommendations.md # Architect recommendations
├── internal/                # Internal application code
│   ├── application/         # Application layer
│   │   ├── handlers/       # HTTP handlers
│   │   │   ├── status.go   # Status handler
│   │   │   └── status_test.go # Status handler tests
│   │   ├── middleware/     # HTTP middleware
│   │   │   ├── chain.go    # Middleware chain
│   │   │   ├── cors.go     # CORS middleware
│   │   │   └── security.go # Security middleware
│   │   └── routes/         # Route definitions
│   │       └── routes.go   # Route setup
│   ├── checker/            # Health checking logic
│   │   ├── commands.go     # Health check commands
│   │   ├── commands_test.go # Command tests
│   │   ├── observer.go     # Observer pattern
│   │   ├── observer_test.go # Observer tests
│   │   ├── service.go      # Checker service
│   │   └── service_test.go # Service tests
│   ├── container/          # Dependency injection
│   │   ├── container.go    # Container implementation
│   │   └── container_test.go # Container tests
│   ├── domain/             # Domain models
│   │   ├── healthcheck/    # Health check domain
│   │   ├── incident/       # Incident domain
│   │   └── service/        # Service domain
│   │       ├── entity.go   # Service entity
│   │       ├── errors.go   # Service errors
│   │       └── repository.go # Service repository
│   ├── infrastructure/     # Infrastructure layer
│   │   ├── cache/          # Caching implementations
│   │   ├── database/       # Database implementations
│   │   │   ├── interfaces.go # Database interfaces
│   │   │   └── mongo/      # MongoDB implementation
│   │   │       ├── mongo.go # MongoDB connection
│   │   │       ├── mongo_test.go # MongoDB tests
│   │   │       ├── repository.go # MongoDB repository
│   │   │       └── repository_test.go # Repository tests
│   │   ├── external/       # External service integrations
│   │   └── messaging/      # Messaging implementations
│   ├── server/             # Server implementations
│   │   ├── interfaces.go   # Server interfaces
│   │   └── server.go       # Server implementation
│   └── shared/             # Shared utilities
│       ├── config/         # Configuration management
│       │   ├── config.go   # Configuration implementation
│       │   └── config_test.go # Configuration tests
│       ├── errors/         # Error handling
│       │   └── errors.go   # Error definitions
│       ├── logger/         # Logging utilities
│       │   ├── logger.go   # Logger implementation
│       │   └── logger_test.go # Logger tests
│       └── utils/          # General utilities
├── mocks/                   # Generated mock files
├── pkg/                     # Public packages
├── scripts/                 # Utility scripts
│   ├── backup/             # Backup scripts
│   ├── db/                 # Database scripts
│   │   ├── backup/         # Database backup
│   │   └── migrations/     # Database migrations
│   ├── hooks/              # Git hooks
│   ├── infra/              # Infrastructure scripts
│   ├── lint/               # Linting scripts
│   ├── maintenance/        # Maintenance scripts
│   ├── monitor/            # Monitoring scripts
│   ├── test/               # Testing scripts
│   └── utils/              # Utility scripts
├── test/                    # Test files
│   ├── api/                # API tests
│   ├── database/           # Database tests
│   └── e2e/                # End-to-end tests
├── testutil/                # Test utilities
│   └── helper.go           # Test helper functions
├── web/                     # Frontend application
│   ├── react-status-page/  # React application (Vite + Tailwind)
│   │   ├── dist/           # Built assets (ignored in git)
│   │   ├── public/         # Public assets
│   │   │   └── vite.svg    # Vite logo
│   │   ├── src/            # React source code
│   │   │   ├── assets/     # Static assets
│   │   │   │   └── react.svg # React logo
│   │   │   ├── components/ # React components
│   │   │   │   ├── Dashboard/ # Dashboard components
│   │   │   │   │   ├── ServiceCard.jsx # Service card
│   │   │   │   │   ├── StatusDashboard.jsx # Status dashboard
│   │   │   │   │   └── StatusIndicator.jsx # Status indicator
│   │   │   │   ├── Incidents/ # Incident components
│   │   │   │   │   ├── IncidentCard.jsx # Incident card
│   │   │   │   │   ├── IncidentManager.jsx # Incident manager
│   │   │   │   │   └── MaintenanceSchedule.jsx # Maintenance
│   │   │   │   ├── Layout/ # Layout components
│   │   │   │   │   ├── Footer.jsx # Footer
│   │   │   │   │   ├── Header.jsx # Header
│   │   │   │   │   └── ThemeToggle.jsx # Theme toggle
│   │   │   │   └── common/ # Common components
│   │   │   │       ├── ErrorBoundary.jsx # Error boundary
│   │   │   │       └── LoadingSpinner.jsx # Loading spinner
│   │   │   ├── contexts/   # React contexts
│   │   │   │   └── ThemeContext.jsx # Theme context
│   │   │   ├── hooks/      # Custom hooks
│   │   │   │   ├── useApi.js # API hook
│   │   │   │   └── usePolling.js # Polling hook
│   │   │   ├── services/   # API services
│   │   │   ├── utils/      # Utility functions
│   │   │   ├── App.jsx     # Main app component
│   │   │   ├── index.css   # Global styles
│   │   │   └── main.jsx    # App entry point
│   │   ├── Dockerfile      # Docker configuration
│   │   ├── eslint.config.js # ESLint configuration
│   │   ├── index.html      # HTML entry point
│   │   ├── nginx.conf      # Nginx configuration
│   │   ├── package.json    # NPM dependencies
│   │   ├── postcss.config.js # PostCSS configuration
│   │   ├── tailwind.config.js # Tailwind CSS configuration
│   │   └── vite.config.js  # Vite configuration
│   ├── Dockerfile.dev      # Development Docker config
│   ├── jest.config.js      # Jest testing configuration
│   ├── package-lock.json   # NPM lock file
│   └── public/             # Shared public assets
├── .golangci.yml           # Go linter configuration
├── .mockery.yaml           # Mockery configuration
├── CHANGELOG.md            # Change log
├── docker-compose.dev.yml  # Development Docker Compose
├── docker-compose.prod.yml # Production Docker Compose
├── docker-compose.yml      # Docker Compose configuration
├── go.mod                  # Go module definition
├── go.sum                  # Go module checksums
├── main.go                 # Application entry point
└── README.md               # Project documentation
```