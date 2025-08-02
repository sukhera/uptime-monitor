#!/bin/bash
# Database optimization and maintenance

set -e

MONGO_URI=${MONGO_URI:-"mongodb://localhost:27017/status_page"}

echo "Starting database optimization..."

# Create optimal indexes
mongosh "$MONGO_URI" << 'EOF'
use status_page;

// Services collection indexes
db.services.createIndex({ "slug": 1 }, { unique: true });
db.services.createIndex({ "enabled": 1 });
db.services.createIndex({ "name": "text" });

// Status logs indexes
db.status_logs.createIndex({ "service_name": 1, "timestamp": -1 });
db.status_logs.createIndex({ "timestamp": -1 });
db.status_logs.createIndex({ "status": 1, "timestamp": -1 });

// Incidents collection indexes
db.incidents.createIndex({ "created_at": -1 });
db.incidents.createIndex({ "severity": 1, "created_at": -1 });
db.incidents.createIndex({ "affected_services": 1 });

// Maintenance collection indexes
db.maintenance.createIndex({ "scheduled_start": 1 });
db.maintenance.createIndex({ "status": 1, "scheduled_start": 1 });

print("✓ Indexes created/updated");

// Analyze collection statistics
print("\n=== Collection Statistics ===");
["services", "status_logs", "incidents", "maintenance"].forEach(function(collection) {
    var stats = db[collection].stats();
    print(collection + ": " + stats.count + " documents, " + 
          Math.round(stats.storageSize / 1024 / 1024 * 100) / 100 + " MB");
});

// Check index usage
print("\n=== Index Usage ===");
db.runCommand({ "collStats": "status_logs", "indexDetails": true });

EOF

echo "✓ Database optimization completed"