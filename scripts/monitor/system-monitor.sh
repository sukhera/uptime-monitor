#!/bin/bash
# Comprehensive system monitoring

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}=== System Monitoring Dashboard ===${NC}"
echo "Generated at: $(date)"
echo ""

# Docker containers status
echo -e "${YELLOW}Docker Containers:${NC}"
docker-compose ps --format "table {{.Name}}\t{{.State}}\t{{.Status}}\t{{.Ports}}"
echo ""

# Resource usage
echo -e "${YELLOW}Resource Usage:${NC}"
docker stats --no-stream --format "table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.MemPerc}}\t{{.NetIO}}\t{{.BlockIO}}"
echo ""

# Disk usage
echo -e "${YELLOW}Disk Usage:${NC}"
df -h | grep -E "(Filesystem|/dev/)"
echo ""

# Application health checks
echo -e "${YELLOW}Health Checks:${NC}"

# API health
if curl -f http://localhost/api/health >/dev/null 2>&1; then
    echo -e "${GREEN}✓ API: Healthy${NC}"
else
    echo -e "${RED}✗ API: Unhealthy${NC}"
fi

# Web health
if curl -f http://localhost >/dev/null 2>&1; then
    echo -e "${GREEN}✓ Web: Healthy${NC}"
else
    echo -e "${RED}✗ Web: Unhealthy${NC}"
fi

# Database health
if docker-compose exec -T mongo mongosh --eval "db.adminCommand('ismaster')" >/dev/null 2>&1; then
    echo -e "${GREEN}✓ Database: Healthy${NC}"
else
    echo -e "${RED}✗ Database: Unhealthy${NC}"
fi

echo ""

# Database statistics
echo -e "${YELLOW}Database Statistics:${NC}"
docker-compose exec -T mongo mongosh status_page --quiet --eval "
    print('Collections:');
    db.getCollectionNames().forEach(function(name) {
        var count = db[name].count();
        print('  ' + name + ': ' + count + ' documents');
    });
    
    print('\nDatabase size: ' + Math.round(db.stats().dataSize / 1024 / 1024 * 100) / 100 + ' MB');
" 2>/dev/null || echo "Database statistics unavailable"

echo ""

# Recent logs
echo -e "${YELLOW}Recent Error Logs (last 10):${NC}"
docker-compose logs --tail=10 2>&1 | grep -i error || echo "No recent errors"