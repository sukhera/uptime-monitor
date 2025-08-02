#!/bin/bash
# Automated deployment script

set -e

ENVIRONMENT=${1:-dev}
VERSION=${2:-latest}

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}Deploying to $ENVIRONMENT environment (version: $VERSION)...${NC}"

case $ENVIRONMENT in
  dev)
    echo -e "${YELLOW}Development deployment${NC}"
    docker-compose -f docker-compose.yml up -d --build
    ;;
  staging)
    echo -e "${YELLOW}Staging deployment${NC}"
    docker-compose -f docker-compose.yml -f docker-compose.staging.yml up -d --build
    ;;
  prod)
    echo -e "${YELLOW}Production deployment${NC}"
    # Add production safety checks
    read -p "Are you sure you want to deploy to production? (y/N): " confirm
    if [ "$confirm" != "y" ]; then
      echo "Deployment cancelled"
      exit 1
    fi
    docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d --build
    ;;
  *)
    echo -e "${RED}Unknown environment: $ENVIRONMENT${NC}"
    echo "Valid environments: dev, staging, prod"
    exit 1
    ;;
esac

# Wait for services
echo -e "${BLUE}Waiting for services to be ready...${NC}"
./scripts/wait-for-services.sh

# Run health checks
echo -e "${BLUE}Running health checks...${NC}"
curl -f http://localhost/api/health
curl -f http://localhost

echo -e "${GREEN}✓ Deployment to $ENVIRONMENT completed successfully!${NC}"