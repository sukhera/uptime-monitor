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
├── cmd/                    # Application entry points
│   ├── api/               # API server
│   └── status-checker/    # Health checker service
├── internal/              # Private application code
│   ├── api/              # API handlers, middleware, routes
│   ├── checker/          # Health checking logic
│   ├── database/         # Database connections
│   ├── models/           # Data models
│   └── config/           # Configuration management
├── configs/              # Configuration files
│   ├── docker/          # Docker configurations
│   ├── nginx/           # Nginx configurations
│   └── env/             # Environment templates
├── web/                  # Frontend application
│   ├── src/             # Source files
│   └── dist/            # Built files
├── scripts/              # Deployment and utility scripts
├── docs/                 # Documentation
├── tests/                # Test files
└── deployments/          # Deployment configurations
```