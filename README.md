# Status Page Starter

A small, production‑minded starter for a public status page and uptime monitor. Backend is Go + MongoDB; frontend is React + Tailwind; everything runs via Docker Compose.

[![CI](https://github.com/sukhera/uptime-monitor/actions/workflows/ci.yml/badge.svg)](https://github.com/sukhera/uptime-monitor/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org/)
[![Coverage](https://img.shields.io/badge/coverage-check%20CI-blue.svg)](https://github.com/sukhera/uptime-monitor/actions/workflows/ci.yml)
[![Security](https://img.shields.io/badge/security-gosec%20%7C%20trivy-green.svg)](https://github.com/sukhera/uptime-monitor/actions/workflows/security.yml)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Docker](https://img.shields.io/badge/docker-ready-blue.svg)](docker-compose.yml)

## What This Is

A lightweight status page system with automated health checks and real-time monitoring. Built with enterprise-grade Go backend patterns, bounded worker pools, and production-ready infrastructure automation.

## Features

- 🔄 **Real-time Service Monitoring** with bounded worker pools and retry logic
- 📊 **Modern Dashboard** with dark mode and responsive design  
- ⚡ **High Performance** Go backend with optimized MongoDB indexes
- 🐳 **Production Ready** with graceful shutdown and security hardening
- 📈 **Observability** with structured logging and context propagation
- 🛠️ **Infrastructure Automation** with comprehensive Makefile commands

## Quick Start

```bash
git clone https://github.com/sukhera/uptime-monitor.git
cd uptime-monitor
cp .env.example .env
make dev
```

### Access Points

- **Status Dashboard**: http://localhost (configurable in `.env`)
- **API Health**: http://localhost/api/health
- **API Status**: http://localhost/api/status  
- **MongoDB**: mongodb://localhost:27017

## Configuration

Configuration is env-first; `.env` sets ports/URIs for dev, compose wires containers; prod uses environment/secret manager.

- **Environment**: Copy `.env.example` to `.env` and customize
- **YAML Config**: See `config.example.yaml` for structured configuration
- **Precedence**: Environment variables > flags > config file > defaults

Key settings:
```bash
PORT=8080                           # Server port
MONGO_URI=mongodb://localhost:27017 # Database connection
READ_TIMEOUT=15s                    # HTTP server timeouts
LOG_LEVEL=info                      # Logging verbosity
CHECK_INTERVAL=2m                   # Health check frequency
```

## Project Structure

```
├── cmd/                    # Application entrypoints
│   ├── api/               # REST API server
│   └── status-checker/    # Health check daemon
├── internal/              # Private application code
│   ├── application/       # HTTP handlers, middleware, routes
│   ├── checker/           # Health check engine (Command pattern)
│   ├── container/         # Dependency injection container
│   ├── infrastructure/    # Database, external services
│   ├── server/           # HTTP server with graceful shutdown
│   └── shared/           # Config, logging, utilities
├── web/                   # Frontend React application
├── deployments/           # Kubernetes manifests
└── docs/                  # Architecture and design docs
```

→ **Full structure**: See [docs/architecture.md](docs/architecture.md)

## API

### Health Check
```bash
GET /api/health
```
```json
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### Service Status
```bash
GET /api/status
```
```json
{
  "overall_status": "operational",
  "services": [
    {
      "name": "API Server",
      "status": "operational",
      "last_check": "2024-01-15T10:29:30Z",
      "response_time": 145
    }
  ]
}
```

## Health Check Architecture

The system uses **Command** and **Observer** patterns for modular, extensible health checking:

- **Bounded Worker Pools**: Configurable concurrency (default: 10 workers, max 5 per service)
- **Timeout Management**: Per-probe (30s) and global (5m) budgets with context cancellation
- **Retry Logic**: Exponential backoff with jitter (3 attempts, 2x multiplier)
- **Event Processing**: Decoupled observers for logging, metrics, and alerting

Design patterns used here are documented under `/docs`—the README stays high-level so you can find what you need fast.

→ **Implementation details**: [docs/design-patterns.md](docs/design-patterns.md)

## MongoDB Persistence & Data Retention

### Collections & Indexes

- **`services`**: Service definitions
  - `services.slug` (unique) - Fast service lookups
  - `services.name`, `services.enabled` - Query optimization
  
- **`status_logs`**: Health check results with 30-day TTL
  - `status_logs(service_id, created_at)` - Efficient time-series queries
  - `status_logs.created_at` (TTL: 30 days) - Automated retention

All operations use context deadlines for resource protection.

## Observability

- **Structured Logging**: Context-aware with trace correlation
- **Metrics Collection**: Health check latency, success rates, system metrics
- **Health Checks**: Database connectivity, dependency status
- **Error Tracking**: Comprehensive error wrapping with stack traces

## Deployment

### Docker Compose
Production-ready setup with `docker-compose.prod.yml` featuring web app behind Nginx reverse proxy, API server, MongoDB with persistence, and automated health checks.

```bash
make prod    # Start production environment
```

### Kubernetes
Deployment includes: Deployment + Service configs, ConfigMap/Secret management, liveness/readiness probes, HPA for autoscaling, and Ingress configuration.

→ **Manifests**: [deployments/kubernetes/](deployments/kubernetes/)

## Development Commands

Run `make help` to see the full toolbox; this README only shows the commands you'll use daily.

```bash
make setup     # Complete project setup
make dev       # Start development environment  
make test      # Run all tests
make lint      # Run all linters
make security  # Security scanning (gosec, npm audit)
make build     # Build all components
make clean     # Clean build artifacts
make docs      # Generate documentation
```

## Project Status

### ✅ Implemented
- Production-ready HTTP server with security hardening
- Bounded worker pool health checking with retry logic
- MongoDB persistence with optimized indexes and TTL
- Graceful shutdown with context propagation
- Comprehensive configuration management
- Docker Compose deployment setup
- Structured logging and error handling
- Modern React dashboard with dark mode

### 🚧 Roadmap
- Incident management and tracking
- Maintenance scheduling system
- Kubernetes deployment automation
- Prometheus metrics integration
- Advanced alerting system
- API rate limiting
- OAuth2 authentication

## Contributing

1. **Conventional Commits**: Use conventional commit format
2. **Quality Gates**: Run `make lint test` before submitting
3. **Small PRs**: Keep changes focused and reviewable
4. **Tests Required**: Add tests for new functionality
5. **Documentation**: Update docs for API/config changes

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---
