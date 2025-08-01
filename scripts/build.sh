#!/bin/bash

set -e

echo "Building Status Page Application..."

# Build Go modules
echo "Building API..."
go build -o bin/api ./cmd/api

echo "Building Status Checker..."
go build -o bin/status-checker ./cmd/status-checker

# Build frontend (if you add a build process later)
echo "Building Frontend..."
mkdir -p web/dist
cp -r web/src/* web/dist/

echo "Build completed successfully!"