#!/bin/bash
# Run all tests with coverage

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}Running comprehensive test suite...${NC}"

# Create test results directory
mkdir -p test-results

# Run Go tests
echo -e "${YELLOW}Running Go tests...${NC}"
go test -v -race -coverprofile=test-results/coverage.out ./... | tee test-results/go-test-results.txt
go tool cover -html=test-results/coverage.out -o test-results/go-coverage.html
go tool cover -func=test-results/coverage.out | tail -n 1 | awk '{print "Go Coverage: " $3}'

# Run frontend tests
if [ -d "web" ]; then
  echo -e "${YELLOW}Running frontend tests...${NC}"
  cd web
  npm test -- --coverage --watchAll=false --outputFile=../test-results/frontend-test-results.json --json
  cd ..
fi

# Run integration tests
echo -e "${YELLOW}Running integration tests...${NC}"
docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit
docker-compose -f docker-compose.test.yml down

# Run security tests
echo -e "${YELLOW}Running security tests...${NC}"
gosec -fmt json -out test-results/security-report.json ./... || true

echo -e "${GREEN}✓ All tests completed! Results in test-results/${NC}"