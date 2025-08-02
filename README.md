# Status Page Starter

A comprehensive uptime monitoring solution with real-time service monitoring, built with Go, MongoDB, and modern web technologies. Features a modern, responsive dashboard with dark mode, real-time updates, incident tracking, and maintenance scheduling. Now includes **complete infrastructure automation** with comprehensive Makefile commands for development, testing, deployment, and maintenance.

## ✨ Features

- 🔄 **Real-time Service Monitoring** with automated health checks
- 📊 **Modern Dashboard** with dark mode and responsive design  
- 🚨 **Incident Management** with severity tracking and resolution workflow
- 🔧 **Maintenance Scheduling** with automated notifications
- ⚡ **High Performance** Go backend with MongoDB storage
- 🐳 **Containerized** deployment with Docker and Docker Compose
- 🛠️ **Complete Infrastructure Automation** with 50+ Makefile commands
- 🧪 **Comprehensive Testing** with automated CI/CD pipeline
- 📈 **Monitoring & Alerting** with performance metrics and health checks

## 🚀 Quick Start

### Prerequisites
- Docker and Docker Compose
- Go 1.21+
- Node.js 18+
- Git

### One-Command Setup
```bash
git clone https://github.com/sukhera/uptime-monitor.git
cd uptime-monitor
make setup    # Complete project setup
make dev      # Start development environment
```

### Access Points
- **Status Dashboard**: http://localhost
- **API Endpoint**: http://localhost/api/status  
- **API Health**: http://localhost/api/health
- **MongoDB**: mongodb://localhost:27017


### 🔧 Tech Stack

#### Backend
- **Go 1.21+**: High-performance backend services
- **MongoDB**: Document-based data storage
- **Docker**: Containerized deployment
- **Nginx**: Reverse proxy and static serving

#### Frontend
- **React**: Modern JavaScript library for building user interfaces
- **Tailwind CSS**: Utility-first CSS framework for rapid UI development
- **Vite**: Fast build tool and development server
- **TypeScript**: Type-safe JavaScript development

#### Infrastructure  
- **Docker Compose**: Multi-service orchestration
- **Makefile Automation**: 50+ commands for complete workflow automation
- **GitHub Actions**: CI/CD pipeline ready
- **Monitoring & Alerting**: Performance metrics and health checks
- **Database Management**: Migrations, optimization, and automated backups

## 🛠️ Infrastructure Automation

This project includes comprehensive infrastructure automation through a feature-rich Makefile with 50+ commands organized by category.

### Available Commands

```bash
make help           # Show all available commands with descriptions
```

#### 🏗️ Setup & Installation
```bash
make setup          # Complete project setup (Go, Node.js, Docker, Git hooks)
make setup-go       # Setup Go environment and tools
make setup-node     # Setup Node.js dependencies
make setup-docker   # Setup Docker environment
make setup-git-hooks # Install Git pre-commit hooks
make env-check      # Check environment requirements
```

#### 🚀 Development
```bash
make dev            # Start complete development environment
make dev-logs       # Follow development logs
make dev-rebuild    # Rebuild and restart development environment
make dev-frontend   # Start frontend development server only
make dev-api        # Run API server locally
make dev-checker    # Run status checker locally
```

#### 🧪 Testing & Quality
```bash
make test           # Run all tests (Go, frontend, integration)
make test-go        # Run Go tests with coverage
make test-frontend  # Run frontend tests
make test-integration # Run integration tests
make test-e2e       # Run end-to-end tests
make lint           # Run all linters
make lint-go        # Run Go linter (golangci-lint)
make lint-frontend  # Run frontend linter
make format         # Format all code
make format-go      # Format Go code
make format-frontend # Format frontend code
```

#### 🔒 Security
```bash
make security       # Run all security scans
make security-go    # Run Go security scan (gosec)
make security-frontend # Run frontend security audit
make security-docker # Run Docker security scan (trivy)
```

#### 🏗️ Build & Deploy
```bash
make build          # Build all services
make build-frontend # Build frontend for production
make build-docker   # Build Docker images
make deploy-dev     # Deploy to development
make deploy-staging # Deploy to staging
make deploy-prod    # Deploy to production (with safety checks)
```

#### 🗄️ Database Management
```bash
make db-start       # Start database only
make db-stop        # Stop database
make db-shell       # Connect to database shell
make seed-db        # Seed database with sample data
make backup-db      # Create database backup
make restore-db BACKUP=<name> # Restore database from backup
```

#### 📊 Monitoring & Maintenance
```bash
make status         # Show service status
make logs           # Show all logs
make logs-api       # Show API logs only
make logs-checker   # Show status checker logs only
make logs-web       # Show web server logs only
make health-check   # Run health checks
make monitor        # Show system monitoring dashboard
```

#### 🧹 Cleanup
```bash
make clean          # Clean development environment
make clean-all      # Deep clean (removes all containers, images, volumes)
```

#### 📚 Documentation
```bash
make docs           # Generate all documentation
make docs-api       # Generate API documentation (Swagger)
make docs-frontend  # Generate frontend documentation
make docs-serve     # Serve documentation locally
```

#### 🔧 Utilities
```bash
make wait-for-services # Wait for all services to be ready
make version        # Show version information for all tools
```

#### 🚀 CI/CD
```bash
make ci-setup       # Setup CI environment
make ci-test        # CI test pipeline
make ci-build       # CI build pipeline
make ci-deploy      # CI deploy pipeline
```

## 🔧 Configuration

### Adding Services
Services are stored in MongoDB. You can add them via the seed script or directly:

```javascript
// Example service configuration
{
  name: "My API",
  slug: "my-api",
  url: "https://api.example.com/health",
  headers: {
    "Authorization": "Bearer token"
  },
  expected_status: 200,
  enabled: true
}
```

### Environment Configuration

Copy the example environment file and configure as needed:
```bash
cp .env.example .env
# Edit .env with your specific configuration
```

Key configuration areas:
- **Database**: MongoDB connection and authentication
- **Security**: JWT secrets and SSL/TLS certificates
- **Monitoring**: Alert thresholds and webhook URLs
- **Storage**: Backup retention and cloud storage settings

## 📁 Project Structure

```
status_page_starter/
├── cmd/                        # Application entry points
│   ├── api/                   # API server main
│   └── status-checker/        # Status checker main
├── internal/                  # Private application code
│   ├── api/                  # API handlers, middleware, routes
│   ├── checker/              # Health checking logic
│   ├── database/             # Database connections
│   ├── models/               # Data models
│   └── config/               # Configuration management
├── configs/                  # Configuration files
│   ├── docker/              # Docker configurations
│   │   ├── Dockerfile.api.dev        # Development API Dockerfile
│   │   ├── Dockerfile.api.prod       # Production API Dockerfile
│   │   └── Dockerfile.status-checker # Status checker Dockerfile
│   ├── dev/                 # Development configurations
│   │   └── air.toml         # Go hot reloading configuration
│   ├── nginx/               # Nginx configurations
│   └── env/                 # Environment templates
├── web/                     # React Frontend Application
│   ├── src/                 # TypeScript source files
│   │   ├── components/      # React components
│   │   │   ├── StatusDashboard.tsx   # Main dashboard component
│   │   │   ├── IncidentManager.tsx   # Incident tracking component
│   │   │   └── ui/          # Shared UI components
│   │   ├── services/        # API service layer
│   │   │   └── api.ts       # HTTP client with TypeScript
│   │   ├── types/           # TypeScript type definitions
│   │   ├── hooks/           # Custom React hooks
│   │   ├── utils/           # Utility functions
│   │   └── main.tsx         # Application entry point
│   ├── public/              # Static assets
│   ├── dist/                # Production build output
│   ├── package.json         # Dependencies and scripts
│   ├── vite.config.ts       # Vite configuration
│   ├── tailwind.config.js   # Tailwind CSS configuration
│   ├── jest.config.js       # Jest testing configuration
│   ├── Dockerfile.dev       # Development frontend Dockerfile
│   └── tsconfig.json        # TypeScript configuration
├── scripts/                 # Infrastructure automation scripts
│   ├── infra/              # Infrastructure scripts
│   │   └── deploy.sh       # Multi-environment deployment
│   ├── hooks/              # Git hooks
│   │   └── pre-commit      # Pre-commit quality checks
│   ├── utils/              # Utility scripts
│   │   ├── reset-dev.sh    # Reset development environment
│   │   └── quick-start.sh  # Quick development setup
│   ├── test/               # Testing automation
│   │   └── run-all-tests.sh # Comprehensive test runner
│   ├── lint/               # Linting automation
│   │   └── run-all-linters.sh # Multi-language linting
│   ├── db/                 # Database management
│   │   ├── migrate.sh      # Database migration system
│   │   ├── optimize.sh     # Database optimization
│   │   ├── cleanup.sh      # Data cleanup and maintenance
│   │   └── migrations/     # Migration files
│   ├── monitor/            # Monitoring and alerting
│   │   ├── system-monitor.sh    # System monitoring dashboard
│   │   ├── performance-monitor.sh # Performance monitoring
│   │   └── log-aggregator.sh    # Log aggregation and analysis
│   ├── maintenance/        # Automated maintenance
│   │   └── auto-maintenance.sh  # Scheduled maintenance tasks
│   ├── wait-for-services.sh # Service startup orchestration
│   ├── env-check.sh        # Environment validation
│   └── seed-db.sh          # Database seeding
├── data/                    # Data and seed files
├── docs/                    # API and architecture documentation
├── tests/                   # Test files
├── deployments/             # Deployment configurations (K8s, Helm)
├── backups/                 # Database backups
├── logs/                    # Application logs
├── test-results/            # Test output and coverage reports
├── lint-results/            # Linting results and reports
├── reports/                 # Monitoring and maintenance reports
├── Makefile                 # Comprehensive automation commands (50+)
├── .golangci.yml           # Go linting configuration
├── docker-compose.yml      # Main service orchestration
├── docker-compose.dev.yml  # Development environment overrides
├── docker-compose.prod.yml # Production environment configuration
├── go.mod                  # Go module definition
└── README.md              # This comprehensive guide
```


## 🛠️ Development

### Automated Development Workflow

The project includes comprehensive automation for all development tasks:

#### Complete Setup (One Command)
```bash
make setup          # Install all dependencies, setup Git hooks, validate environment
make dev            # Start complete development environment with hot reloading
```

#### Individual Development Commands
```bash
# Environment validation
make env-check      # Check all prerequisites (Go, Node.js, Docker, etc.)

# Development servers
make dev-frontend   # Start frontend with hot reloading (Vite)
make dev-api        # Start API server with hot reloading (Air)
make dev-checker    # Start status checker locally

# Database operations
make db-start       # Start MongoDB only
make seed-db        # Populate with sample data
make db-shell       # Connect to MongoDB shell

# Development utilities
make dev-logs       # Follow all service logs
make dev-rebuild    # Clean rebuild of development environment
```

#### Hot Reloading & Live Development
- **Go Backend**: Automatic recompilation and restart with Air
- **React Frontend**: Instant updates with Vite HMR
- **Volume Mounts**: Live code changes without rebuilding containers
- **Synchronized Services**: Automatic service orchestration

#### Quality Assurance
```bash
make lint           # Run all linters (Go, TypeScript, CSS, Shell, YAML)
make format         # Auto-format all code
make test           # Run comprehensive test suite
make security       # Security scanning and vulnerability checks
```

### Database Schema

#### Services Collection
```javascript
{
  name: "Service Name",
  slug: "service-slug",
  url: "https://service.com/health",
  headers: {}, // Optional custom headers
  expected_status: 200,
  enabled: true
}
```

#### Status Logs Collection
```javascript
{
  service_name: "Service Name",
  status: "operational|degraded|down",
  latency_ms: 150,
  status_code: 200,
  error: "Error message if any",
  timestamp: ISODate("2024-01-01T00:00:00Z")
}
```

## 📊 Monitoring & Maintenance

### Comprehensive Monitoring System

#### System Health Dashboard
```bash
make monitor        # Real-time system monitoring dashboard
make status         # Service status overview
make health-check   # Application health verification
```

#### Log Management
```bash
make logs           # Follow all service logs
make logs-api       # API server logs only
make logs-checker   # Status checker logs only
make logs-web       # Web server logs only
```

#### Performance Monitoring
- **Resource Usage**: CPU, memory, disk monitoring with alerting
- **Application Metrics**: Response times, error rates, throughput
- **Database Performance**: Query performance, index usage, storage metrics
- **Alert Thresholds**: Configurable performance thresholds with webhook notifications

#### Automated Log Analysis
- **Error Pattern Detection**: Automatic identification of common errors
- **Log Aggregation**: Daily log collection and analysis
- **Historical Trending**: Performance trend analysis over time
- **Alert Integration**: Webhook notifications for critical issues

### Database Management

#### Automated Database Operations
```bash
make backup-db                    # Create timestamped backup
make restore-db BACKUP=<name>    # Restore from specific backup
make db-shell                     # Interactive MongoDB shell
```

#### Database Optimization & Maintenance
```bash
./scripts/db/optimize.sh          # Index optimization and performance tuning
./scripts/db/cleanup.sh           # Remove old data based on retention policy
./scripts/db/migrate.sh up        # Run pending database migrations
./scripts/db/migrate.sh status    # Show migration status
```

#### Migration System
```bash
# Create new migration
./scripts/db/migrate.sh create add_user_preferences

# Run all pending migrations
./scripts/db/migrate.sh up

# Check migration status
./scripts/db/migrate.sh status
```

## 🚀 Deployment

### Automated Multi-Environment Deployment

#### Simple One-Command Deployments
```bash
make deploy-dev      # Deploy to development environment
make deploy-staging  # Deploy to staging environment  
make deploy-prod     # Deploy to production (with safety confirmations)
```

#### Advanced Deployment Scripts
```bash
# Multi-environment deployment script
./scripts/infra/deploy.sh dev      # Development
./scripts/infra/deploy.sh staging  # Staging
./scripts/infra/deploy.sh prod     # Production (with confirmations)
```

#### Build & Release Management
```bash
make build           # Build all services for deployment
make build-frontend  # Build optimized frontend bundle
make build-docker    # Build all Docker images
```

#### Production Deployment Features
- **Safety Checks**: Production deployments require explicit confirmation
- **Health Verification**: Automated health checks after deployment
- **Service Orchestration**: Proper startup order with dependency management
- **Zero-Downtime Updates**: Rolling updates with health monitoring
- **Rollback Capability**: Quick rollback to previous stable version

### Environment-Specific Configurations

#### Development Environment
- Hot reloading enabled for all services
- Debug logging and development tools included
- Volume mounts for live code changes
- Exposed ports for direct service access

#### Staging Environment  
- Production-like configuration for testing
- Performance monitoring enabled
- Integration test execution
- Security scanning validation

#### Production Environment
- Optimized builds with minimal attack surface
- SSL/TLS termination and security headers
- Multi-replica services for high availability
- Comprehensive monitoring and alerting
- Automated backup and maintenance

### Docker Orchestration
```bash
# Development with hot reloading
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d

# Production deployment
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d

# Service scaling
docker-compose up -d --scale api=3 --scale status-checker=2

# Rolling updates
make deploy-prod     # Automated rolling update with health checks
```

### Container Registry & CI/CD
- **Automated Builds**: CI/CD pipeline integration ready
- **Image Tagging**: Semantic versioning and environment tagging
- **Security Scanning**: Container vulnerability assessment
- **Multi-Architecture**: ARM64 and AMD64 support ready


## 🤝 Contributing

### Development Workflow
1. **Fork** the repository
2. **Clone** your fork locally
3. **Create** a feature branch (`git checkout -b feature/amazing-feature`)
4. **Make** your changes following the coding standards
5. **Test** your changes thoroughly
6. **Commit** with conventional commit messages
7. **Push** to your fork (`git push origin feature/amazing-feature`)
8. **Submit** a pull request with detailed description

### Code Standards
- **Go**: Follow `gofmt` and `golint` standards
- **JavaScript**: ES6+ with modern patterns, avoid jQuery
- **CSS**: Use custom properties, mobile-first approach
- **HTML**: Semantic HTML5 with proper accessibility

### Comprehensive Testing Automation

#### Complete Test Suite
```bash
make test           # Run all tests (Go, frontend, integration, security)
make ci-test        # CI test pipeline (includes security scans)
```

#### Individual Test Categories
```bash
make test-go        # Go tests with race detection and coverage
make test-frontend  # React/TypeScript tests with Jest
make test-integration # End-to-end integration tests
make test-e2e       # Browser-based end-to-end tests
```

#### Test Results & Coverage
- **Go Coverage**: HTML coverage reports generated automatically
- **Frontend Coverage**: Jest coverage with threshold enforcement
- **Integration Results**: Comprehensive API and service testing
- **Security Testing**: Automated vulnerability scanning

#### Quality Gates
- **Automated Linting**: Multi-language code quality checks
- **Security Scanning**: Go security analysis with gosec
- **Dependency Auditing**: NPM audit for frontend vulnerabilities
- **Performance Testing**: Load testing with configurable thresholds

#### Pre-commit Automation
All quality checks run automatically before commits via Git hooks:
- Go linting and testing
- Frontend linting and formatting
- Security scanning
- Code formatting validation

## 📄 License

MIT License - see LICENSE file for details.
