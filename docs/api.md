# API Documentation

## Overview

The Status Page API provides endpoints for retrieving service status information and health checks. All endpoints return JSON responses and support CORS for cross-origin requests.

## Base URL

- **Development**: `http://localhost:8080`
- **Production**: `https://your-domain.com`

## Authentication

Currently, the API does not require authentication. However, rate limiting is applied to prevent abuse.

## Rate Limiting

- **Limit**: 10 requests per second per IP address
- **Headers**: Rate limit information is included in response headers
- **Status Code**: 429 (Too Many Requests) when limit is exceeded

## Endpoints

### GET /api/status

Retrieves the current status of all monitored services.

#### Request

```http
GET /api/status
Accept: application/json
```

#### Response

**Success (200 OK)**

```json
{
  "status": "success",
  "data": [
    {
      "name": "API Service",
      "status": "operational",
      "latency_ms": 150,
      "updated_at": "2024-01-15T10:30:00Z"
    },
    {
      "name": "Database Service",
      "status": "degraded",
      "latency_ms": 2500,
      "updated_at": "2024-01-15T10:29:45Z"
    },
    {
      "name": "Web Service",
      "status": "down",
      "latency_ms": 0,
      "updated_at": "2024-01-15T10:28:30Z"
    }
  ],
  "count": 3,
  "timestamp": "2024-01-15T10:30:00Z"
}
```

**Error Responses**

- `500 Internal Server Error`: Database connection issues
- `429 Too Many Requests`: Rate limit exceeded
- `503 Service Unavailable`: Service temporarily unavailable

#### Status Values

- `operational`: Service is functioning normally
- `degraded`: Service is experiencing performance issues
- `down`: Service is completely unavailable

### GET /api/health

Provides a health check for the API service and database connectivity.

#### Request

```http
GET /api/health
Accept: application/json
```

#### Response

**Success (200 OK)**

```json
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z",
  "version": "1.0.0"
}
```

**Error (503 Service Unavailable)**

```json
{
  "status": "unhealthy",
  "error": "database connection failed",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

## Error Handling

All endpoints follow a consistent error response format:

```json
{
  "error": "Error message",
  "status": "error",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### Common HTTP Status Codes

- `200 OK`: Request successful
- `400 Bad Request`: Invalid request parameters
- `404 Not Found`: Endpoint not found
- `429 Too Many Requests`: Rate limit exceeded
- `500 Internal Server Error`: Server error
- `503 Service Unavailable`: Service temporarily unavailable

## Response Headers

All responses include the following headers:

```
Content-Type: application/json
Cache-Control: no-cache, no-store, must-revalidate
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
X-XSS-Protection: 1; mode=block
```

## CORS Support

The API supports Cross-Origin Resource Sharing (CORS) with the following configuration:

- **Allowed Origins**: Configurable (default: `*`)
- **Allowed Methods**: GET, POST, PUT, DELETE, OPTIONS
- **Allowed Headers**: Content-Type, Authorization
- **Max Age**: 86400 seconds (24 hours)

## Client Examples

### JavaScript (Fetch API)

```javascript
// Fetch status data
const response = await fetch('/api/status');
const data = await response.json();

if (data.status === 'success') {
  console.log('Services:', data.data);
} else {
  console.error('Error:', data.error);
}
```

### JavaScript (Axios)

```javascript
import axios from 'axios';

const api = axios.create({
  baseURL: '/api',
  timeout: 10000,
});

// Fetch status data
try {
  const response = await api.get('/status');
  console.log('Services:', response.data.data);
} catch (error) {
  console.error('API Error:', error.response?.data);
}
```

### cURL

```bash
# Get service status
curl -X GET http://localhost:8080/api/status \
  -H "Accept: application/json"

# Health check
curl -X GET http://localhost:8080/api/health \
  -H "Accept: application/json"
```

### Python

```python
import requests

# Get service status
response = requests.get('http://localhost:8080/api/status')
if response.status_code == 200:
    data = response.json()
    print(f"Services: {data['data']}")
else:
    print(f"Error: {response.status_code}")
```

## Monitoring and Alerts

### Health Check Monitoring

Monitor the `/api/health` endpoint to ensure the service is running:

```bash
# Check every 30 seconds
while true; do
  curl -f http://localhost:8080/api/health || echo "Service down"
  sleep 30
done
```

### Status Monitoring

Monitor service status changes:

```javascript
// Poll status every 30 seconds
setInterval(async () => {
  try {
    const response = await fetch('/api/status');
    const data = await response.json();
    
    // Check for any down services
    const downServices = data.data.filter(service => service.status === 'down');
    if (downServices.length > 0) {
      console.warn('Down services:', downServices);
    }
  } catch (error) {
    console.error('Status check failed:', error);
  }
}, 30000);
```

## Rate Limiting

The API implements rate limiting to prevent abuse:

- **Limit**: 10 requests per second per IP
- **Window**: 1 second sliding window
- **Headers**: Rate limit info in response headers
- **Response**: 429 status code when limit exceeded

### Rate Limit Headers

```
X-RateLimit-Limit: 10
X-RateLimit-Remaining: 5
X-RateLimit-Reset: 1642234567
```

## Security

### Security Headers

All responses include security headers:

- `X-Content-Type-Options: nosniff`
- `X-Frame-Options: DENY`
- `X-XSS-Protection: 1; mode=block`
- `Referrer-Policy: strict-origin-when-cross-origin`
- `Content-Security-Policy: default-src 'self'`

### Input Validation

- All endpoints validate input parameters
- Malformed requests return 400 Bad Request
- Content-Type validation for POST/PUT requests

## Versioning

API versioning is handled through the URL path. Current version is v1 (default).

- Current: `/api/status`
- Future: `/api/v2/status`

## Changelog

### v1.0.0 (Current)
- Initial API release
- Status endpoint
- Health check endpoint
- Rate limiting
- Security headers
- CORS support