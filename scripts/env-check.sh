#!/bin/bash
# Check environment requirements

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo "Checking environment requirements..."

# Check Go
if command -v go &> /dev/null; then
  echo -e "${GREEN}✓ Go is installed: $(go version)${NC}"
else
  echo -e "${RED}✗ Go is not installed${NC}"
  exit 1
fi

# Check Node.js
if command -v node &> /dev/null; then
  echo -e "${GREEN}✓ Node.js is installed: $(node --version)${NC}"
else
  echo -e "${RED}✗ Node.js is not installed${NC}"
  exit 1
fi

# Check npm
if command -v npm &> /dev/null; then
  echo -e "${GREEN}✓ npm is installed: $(npm --version)${NC}"
else
  echo -e "${RED}✗ npm is not installed${NC}"
  exit 1
fi

# Check Docker
if command -v docker &> /dev/null; then
  echo -e "${GREEN}✓ Docker is installed: $(docker --version)${NC}"
else
  echo -e "${RED}✗ Docker is not installed${NC}"
  exit 1
fi

# Check Docker Compose
if command -v docker-compose &> /dev/null; then
  echo -e "${GREEN}✓ Docker Compose is installed: $(docker-compose --version)${NC}"
else
  echo -e "${RED}✗ Docker Compose is not installed${NC}"
  exit 1
fi

# Check Docker daemon
if docker info &> /dev/null; then
  echo -e "${GREEN}✓ Docker daemon is running${NC}"
else
  echo -e "${RED}✗ Docker daemon is not running${NC}"
  exit 1
fi

echo -e "${GREEN}✓ All environment requirements met!${NC}"