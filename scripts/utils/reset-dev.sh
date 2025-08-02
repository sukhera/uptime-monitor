#!/bin/bash
# Reset development environment

set -e

echo "Resetting development environment..."

# Stop all services
docker-compose down -v

# Remove development volumes
docker volume rm $(docker volume ls -q | grep status) 2>/dev/null || true

# Rebuild images
docker-compose build --no-cache

# Start fresh
docker-compose up -d

# Wait and seed
./scripts/wait-for-services.sh
./scripts/seed-db.sh

echo "✓ Development environment reset complete!"