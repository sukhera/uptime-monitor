#!/bin/bash
# Run all linters

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}Running all linters...${NC}"

# Create lint results directory
mkdir -p lint-results

# Run Go linting
echo -e "${YELLOW}Running Go linter...${NC}"
golangci-lint run --out-format json ./... > lint-results/go-lint.json || true
golangci-lint run ./... | tee lint-results/go-lint.txt

# Run frontend linting
if [ -d "web" ]; then
  echo -e "${YELLOW}Running frontend linter...${NC}"
  cd web
  npm run lint -- --format json --output-file ../lint-results/frontend-lint.json || true
  npm run lint | tee ../lint-results/frontend-lint.txt || true
  cd ..
fi

# Run Dockerfile linting
if command -v hadolint &> /dev/null; then
  echo -e "${YELLOW}Running Dockerfile linter...${NC}"
  find . -name "Dockerfile*" -exec hadolint {} \; | tee lint-results/dockerfile-lint.txt || true
fi

# Run shell script linting
if command -v shellcheck &> /dev/null; then
  echo -e "${YELLOW}Running shell script linter...${NC}"
  find . -name "*.sh" -exec shellcheck {} \; | tee lint-results/shellcheck.txt || true
fi

# Run YAML linting
if command -v yamllint &> /dev/null; then
  echo -e "${YELLOW}Running YAML linter...${NC}"
  yamllint . | tee lint-results/yaml-lint.txt || true
fi

echo -e "${GREEN}✓ All linting completed! Results in lint-results/${NC}"