# Configuration Guide

## Environment Setup

### Quick Setup
```bash
make setup-env     # Creates .env from .env.example
```

### Manual Setup
```bash
cp .env.example .env
# Edit .env with your configuration
```

## Environment Variables Reference

### Database Configuration
```bash
# MongoDB connection (required)
MONGO_URI=mongodb://mongo:27017/status_page

# MongoDB authentication (production)
MONGO_ROOT_USERNAME=admin
MONGO_ROOT_PASSWORD=secure_password_here
```

### API Configuration
```bash
# API server port
PORT=8080

# Environment mode
GO_ENV=development|production

# JWT secret for authentication (required in production)
JWT_SECRET=your_secure_jwt_secret_key
```

### Frontend Configuration
```bash
# Node.js environment
NODE_ENV=development|production

# API endpoint URL for frontend
VITE_API_URL=http://localhost/api
```

### Monitoring & Alerting
```bash
# Webhook URL for alerts (optional)
WEBHOOK_URL=https://hooks.slack.com/your/webhook/url

# Performance alert thresholds (optional)
ALERT_THRESHOLD_CPU=80
ALERT_THRESHOLD_MEM=85
ALERT_THRESHOLD_DISK=90
```

### Data Management
```bash
# Data retention in days (default: 90)
RETENTION_DAYS=90

# Backup retention in days (default: 30)
BACKUP_RETENTION_DAYS=30
```

### Development Tools
```bash
# Enable hot reloading in development
HOT_RELOAD=true
```

### Cloud Storage (Optional)
```bash
# AWS S3 for backups
AWS_S3_BUCKET=your-backup-bucket
AWS_ACCESS_KEY_ID=your-access-key
AWS_SECRET_ACCESS_KEY=your-secret-key
AWS_REGION=us-east-1
```

## Environment-Specific Configurations

### Development Environment
- Uses `docker-compose.yml` + `docker-compose.dev.yml`
- Hot reloading enabled
- Debug logging
- Exposed database ports
- Default passwords for convenience

### Production Environment
- Uses `docker-compose.yml` + `docker-compose.prod.yml`
- SSL/TLS termination
- No exposed database ports
- Authentication required
- Performance monitoring enabled
- Automated backups

### Security Best Practices

#### Development
- Use default credentials only in development
- Keep `.env` files out of version control
- Use separate `.env` files for different environments

#### Production
- Use strong, unique passwords
- Generate secure JWT secrets (32+ characters)
- Enable SSL/TLS with valid certificates
- Use environment-specific database credentials
- Enable MongoDB authentication
- Set up proper firewall rules

## Configuration Examples

### Local Development
```bash
MONGO_URI=mongodb://mongo:27017/status_page
MONGO_ROOT_USERNAME=admin
MONGO_ROOT_PASSWORD=devpassword
PORT=8080
GO_ENV=development
JWT_SECRET=dev_jwt_secret_not_for_production
NODE_ENV=development
VITE_API_URL=http://localhost/api
HOT_RELOAD=true
RETENTION_DAYS=30
```

### Production
```bash
MONGO_URI=mongodb://mongo:27017/status_page
MONGO_ROOT_USERNAME=admin
MONGO_ROOT_PASSWORD=your_secure_production_password
PORT=8080
GO_ENV=production
JWT_SECRET=your_super_secure_jwt_secret_key_32_chars_min
NODE_ENV=production
VITE_API_URL=/api
WEBHOOK_URL=https://hooks.slack.com/your/production/webhook
RETENTION_DAYS=90
BACKUP_RETENTION_DAYS=30
AWS_S3_BUCKET=your-production-backup-bucket
```

## Validation

### Environment Check
```bash
make env-check     # Validates all required tools and dependencies
```

### Configuration Validation
The application will validate required environment variables on startup and fail fast with clear error messages if configuration is invalid.

## Troubleshooting

### Common Issues

1. **Database Connection Failed**
   - Check `MONGO_URI` format
   - Ensure MongoDB is running
   - Verify authentication credentials

2. **JWT Authentication Errors**
   - Ensure `JWT_SECRET` is set and secure
   - Check secret length (minimum 32 characters recommended)

3. **Frontend API Connection Issues**
   - Verify `VITE_API_URL` points to correct API endpoint
   - Check CORS configuration in production

4. **Alert Webhook Not Working**
   - Validate `WEBHOOK_URL` format
   - Test webhook endpoint manually
   - Check network connectivity

### Debug Mode
Enable debug logging by setting:
```bash
GO_ENV=development
```

This will provide detailed logging for troubleshooting.