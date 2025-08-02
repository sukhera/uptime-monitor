# Uptime Monitor

A complete uptime monitoring solution with real-time service monitoring, built with Go, MongoDB, and modern web technologies. Features a modern, responsive dashboard with dark mode, real-time updates, incident tracking, and maintenance scheduling.


## 🚀 Quick Start

### Prerequisites
- Docker and Docker Compose
- Git

### Setup
1. Clone the repository:
```bash
git clone https://github.com/sukhera/uptime-monitor.git
cd uptime-monitor
```

2. Start all services:
```bash
docker-compose up -d
```

3. Access the status page:
- **Status Dashboard**: http://localhost
- **API Endpoint**: http://localhost/api/status
- **API Health**: http://localhost/api/health


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
- **Kubernetes**: Production deployment (optional)
- **GitHub Actions**: CI/CD pipeline ready

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

### Environment Variables
- `MONGO_URI`: MongoDB connection string (default: `mongodb://localhost:27017`)
- `PORT`: API server port (default: `8080`)

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
│   │   ├── Dockerfile.api
│   │   ├── Dockerfile.status-checker
│   │   └── docker-compose.*.yml
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
│   └── tsconfig.json        # TypeScript configuration
├── data/                    # Data and seed files
├── scripts/                 # Deployment and utility scripts
├── docs/                    # API and architecture documentation
├── tests/                   # Test files
├── deployments/             # Deployment configurations (K8s, Helm)
├── docker-compose.yml       # Main service orchestration
├── go.mod                   # Go module definition
└── README.md               # This comprehensive guide
```


## 🛠️ Development

### Local Development
1. Start MongoDB:
```bash
docker-compose up mongo -d
```

2. Run status checker:
```bash
go run ./cmd/status-checker
```

3. Run API server:
```bash
go run ./cmd/api
```

4. Serve web frontend:
```bash
cd web/dist
python -m http.server 8000
```

### Using Scripts

Build the project:
```bash
./scripts/build.sh
```

Deploy to different environments:
```bash
./scripts/deploy.sh dev    # Development
./scripts/deploy.sh prod   # Production
```

Seed the database:
```bash
./scripts/seed-db.sh
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

## 🔍 Monitoring

### Logs
View service logs:
```bash
# Status checker logs
docker-compose logs status-checker

# API logs
docker-compose logs api

# All logs
docker-compose logs -f
```

### Database
Connect to MongoDB:
```bash
docker-compose exec mongo mongosh
```

## 🚀 Deployment



### Docker Deployment
```bash
# Build and start all services
docker-compose up -d --build

# Scale specific services
docker-compose up -d --scale api=3 --scale status-checker=2

# Update services with zero downtime
docker-compose pull
docker-compose up -d --no-deps api status-checker

# Production deployment
docker-compose -f docker-compose.prod.yml up -d
```

### Kubernetes Deployment
```bash
# Apply Kubernetes manifests
kubectl apply -f deployments/kubernetes/

# Check deployment status
kubectl get pods -n status-page

# Scale deployment
kubectl scale deployment api --replicas=5
```


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

### Testing
```bash
# Run Go tests
go test ./...

# Frontend testing (if implemented)
npm test

# Integration tests
docker-compose -f docker-compose.test.yml up --abort-on-container-exit
```

## 📄 License

MIT License - see LICENSE file for details.
