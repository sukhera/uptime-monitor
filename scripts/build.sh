#!/bin/bash

set -e

echo "Building Status Page Application..."

echo "Building React frontend..."
cd web/react-status-page
npm install
npm run build
cd ../..

echo "Building Go services..."
docker-compose build

echo "Build completed successfully!"