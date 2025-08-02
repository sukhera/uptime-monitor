#!/bin/bash

set -e

ENVIRONMENT=${1:-dev}

echo "Deploying to $ENVIRONMENT environment..."

if [ "$ENVIRONMENT" = "prod" ]; then
    docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d --build
else
    docker-compose up -d --build
fi

echo "Deployment completed!"