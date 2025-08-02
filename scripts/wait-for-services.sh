#!/bin/bash
# Wait for services to be ready

set -e

echo "Waiting for MongoDB to be ready..."
while ! docker-compose exec -T mongo mongosh --eval "db.adminCommand('ismaster')" >/dev/null 2>&1; do
  sleep 2
done
echo "✓ MongoDB is ready"

echo "Waiting for API to be ready..."
while ! curl -f http://localhost/api/health >/dev/null 2>&1; do
  sleep 2
done
echo "✓ API is ready"

echo "Waiting for Web server to be ready..."
while ! curl -f http://localhost >/dev/null 2>&1; do
  sleep 2
done
echo "✓ Web server is ready"

echo "All services are ready!"