#!/bin/bash
# Automated maintenance tasks

set -e

echo "Starting automated maintenance at $(date)"

# Database maintenance
echo "Running database maintenance..."
./scripts/db/optimize.sh
./scripts/db/cleanup.sh

# Log rotation
echo "Rotating logs..."
./scripts/monitor/log-aggregator.sh

# Docker maintenance
echo "Cleaning Docker resources..."
docker system prune -f
docker volume prune -f

# Update health check
echo "Running health checks..."
mkdir -p reports
./scripts/monitor/system-monitor.sh > "reports/maintenance-$(date +%Y%m%d).txt"

# Backup verification
echo "Verifying recent backups..."
if [ -d "backups" ]; then
    LATEST_BACKUP=$(find backups -name "*.tar.gz" -mtime -1 | head -1)
    if [ -z "$LATEST_BACKUP" ]; then
        echo "WARNING: No recent backup found!"
    else
        echo "✓ Recent backup found: $LATEST_BACKUP"
    fi
fi

echo "✓ Automated maintenance completed at $(date)"