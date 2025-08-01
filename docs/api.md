# API Documentation

## Base URL
- Development: `http://localhost/api`
- Production: `https://your-domain.com/api`

## Endpoints

### GET /api/status
Returns the current status of all monitored services.

**Response:**
```json
[
  {
    "name": "Service Name",
    "status": "operational|degraded|down", 
    "latency_ms": 150,
    "updated_at": "2024-01-01T12:00:00Z"
  }
]
```

**Status Values:**
- `operational`: Service is working normally
- `degraded`: Service is working but with issues
- `down`: Service is not responding

### GET /api/health
Health check endpoint for the API service itself.

**Response:**
```json
{
  "status": "healthy"
}
```

## Error Responses

All endpoints may return standard HTTP error codes:

- `500 Internal Server Error`: Database connection issues
- `503 Service Unavailable`: Temporary service issues

Error responses include a message:
```json
{
  "error": "Database connection failed"
}
```

## CORS

The API includes CORS headers to allow cross-origin requests from web frontends.