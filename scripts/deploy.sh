#!/bin/bash

set -e

ENVIRONMENT=${1:-dev}
COMPOSE_FILE="configs/docker/docker-compose.${ENVIRONMENT}.yml"

echo "Deploying to ${ENVIRONMENT} environment..."

if [ ! -f "$COMPOSE_FILE" ]; then
    echo "Error: Compose file ${COMPOSE_FILE} not found"
    exit 1
fi

# Stop existing containers
echo "Stopping existing containers..."
docker-compose -f "$COMPOSE_FILE" down

# Build and start new containers
echo "Building and starting containers..."
docker-compose -f "$COMPOSE_FILE" up -d --build

# Wait for services to be ready
echo "Waiting for services to be ready..."
sleep 10

# Health check
echo "Checking service health..."
if curl -f http://localhost/api/health; then
    echo "Deployment successful!"
else
    echo "Deployment failed - health check failed"
    exit 1
fi