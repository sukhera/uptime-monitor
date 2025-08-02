#!/bin/bash
# Automated backup system

set -e

BACKUP_DIR="/backups"
RETENTION_DAYS=${BACKUP_RETENTION_DAYS:-30}
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

echo "Starting automated backup at $(date)"

# Create backup directory
mkdir -p "$BACKUP_DIR/mongodb/$TIMESTAMP"

# MongoDB backup
if [ -n "$MONGO_URI" ]; then
    mongodump --uri="$MONGO_URI" --out="$BACKUP_DIR/mongodb/$TIMESTAMP"
else
    echo "ERROR: MONGO_URI environment variable not set"
    exit 1
fi

# Compress backup
cd "$BACKUP_DIR/mongodb"
tar -czf "$TIMESTAMP.tar.gz" "$TIMESTAMP"
rm -rf "$TIMESTAMP"

# Clean old backups
find "$BACKUP_DIR/mongodb" -name "*.tar.gz" -mtime +$RETENTION_DAYS -delete

# Upload to cloud storage (if configured)
if [ -n "$AWS_S3_BUCKET" ] && [ -n "$AWS_ACCESS_KEY_ID" ] && [ -n "$AWS_SECRET_ACCESS_KEY" ]; then
    echo "Uploading backup to S3..."
    aws s3 cp "$BACKUP_DIR/mongodb/$TIMESTAMP.tar.gz" "s3://$AWS_S3_BUCKET/backups/mongodb/" \
        --region "${AWS_REGION:-us-east-1}"
fi

echo "Backup completed successfully at $(date)"

# Set up cron job to run daily at 2 AM
if [ ! -f /etc/cron.d/status-page-backup ]; then
    echo "0 2 * * * root /scripts/backup-cron.sh >> /var/log/backup.log 2>&1" > /etc/cron.d/status-page-backup
fi