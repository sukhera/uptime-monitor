#!/bin/bash
# Quick development setup

set -e

echo "Quick development setup..."

# Check requirements
./scripts/env-check.sh

# Install dependencies
go mod download
cd web && npm install && cd ..

# Start minimal services
docker-compose up -d mongo

# Wait for MongoDB
echo "Waiting for MongoDB..."
while ! docker-compose exec -T mongo mongosh --eval "db.adminCommand('ismaster')" >/dev/null 2>&1; do
  sleep 1
done

# Seed database
./scripts/seed-db.sh

echo "✓ Quick setup complete!"
echo "Run 'make dev-api' and 'make dev-frontend' in separate terminals"