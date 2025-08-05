# React Status Page

A modern, responsive status page built with React, Tailwind CSS, and Vite.

## Features

- 🎨 Modern UI with Tailwind CSS
- 🌙 Dark mode support
- 📱 Responsive design
- 🔄 Real-time status updates
- ⚡ Fast development with Vite
- 🚀 Production-ready build

## Development

### Prerequisites

- Node.js 18+ 
- npm or yarn

### Installation

```bash
npm install
```

### Development Server

For local development with the API running on localhost:8080:

```bash
npm run dev:local
```

For development with the default API proxy (requires Docker setup):

```bash
npm run dev
```

### Building for Production

```bash
npm run build
```

### Preview Production Build

```bash
npm run preview
```

## Environment Variables

- `VITE_API_URL`: API base URL (defaults to `/api` for Docker setup)

## API Integration

The frontend expects the following API endpoints:

- `GET /api/status` - Returns array of service statuses
- `GET /api/health` - Health check endpoint
- `GET /api/incidents` - Incidents list
- `GET /api/maintenance` - Maintenance schedule

### Service Status Format

```json
{
  "name": "Service Name",
  "status": "operational|degraded|down",
  "latency_ms": 150,
  "updated_at": "2024-01-01T12:00:00Z",
  "error": "Optional error message"
}
```

## Docker

The frontend is containerized with nginx for production deployment.

```bash
docker build -t status-page-frontend .
```

## Project Structure

```
src/
├── components/          # React components
│   ├── Dashboard/      # Status dashboard components
│   ├── Incidents/      # Incident management components
│   ├── Layout/         # Layout components
│   └── common/         # Shared components
├── hooks/              # Custom React hooks
├── contexts/           # React contexts
├── services/           # API service functions
└── utils/              # Utility functions
```
