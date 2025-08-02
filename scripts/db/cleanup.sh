#!/bin/bash
# Clean up old data

set -e

MONGO_URI=${MONGO_URI:-"mongodb://localhost:27017/status_page"}
RETENTION_DAYS=${RETENTION_DAYS:-90}

echo "Cleaning up data older than $RETENTION_DAYS days..."

mongosh "$MONGO_URI" << EOF
use status_page;

var cutoffDate = new Date();
cutoffDate.setDate(cutoffDate.getDate() - $RETENTION_DAYS);

print("Cutoff date: " + cutoffDate);

// Clean old status logs
var result = db.status_logs.deleteMany({
    "timestamp": { "\$lt": cutoffDate }
});
print("Deleted " + result.deletedCount + " old status logs");

// Clean resolved incidents older than retention period
var incidentResult = db.incidents.deleteMany({
    "resolved_at": { "\$lt": cutoffDate },
    "status": "resolved"
});
print("Deleted " + incidentResult.deletedCount + " old resolved incidents");

// Clean completed maintenance
var maintenanceResult = db.maintenance.deleteMany({
    "completed_at": { "\$lt": cutoffDate },
    "status": "completed"
});
print("Deleted " + maintenanceResult.deletedCount + " old maintenance records");

// Optimize collections
db.runCommand({ "compact": "status_logs" });
db.runCommand({ "compact": "incidents" });
db.runCommand({ "compact": "maintenance" });

print("✓ Cleanup completed");
EOF

echo "✓ Data cleanup completed"