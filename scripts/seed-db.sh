#!/bin/bash
# Seed database with sample data

set -e

MONGO_URI=${MONGO_URI:-"mongodb://localhost:27017/status_page"}

echo "Seeding database with sample data..."

mongosh "$MONGO_URI" << 'EOF'
use status_page;

// Clear existing data
db.services.deleteMany({});
db.status_logs.deleteMany({});
db.incidents.deleteMany({});
db.maintenance.deleteMany({});

print("Seeding sample data...");

// Insert services
const services = [
    {
        name: "API Gateway",
        slug: "api-gateway", 
        url: "https://api.example.com/health",
        headers: { "Authorization": "Bearer test-token" },
        expected_status: 200,
        enabled: true,
        created_at: new Date(),
        updated_at: new Date()
    },
    {
        name: "Web Application",
        slug: "web-app",
        url: "https://app.example.com",
        expected_status: 200,
        enabled: true,
        created_at: new Date(),
        updated_at: new Date()
    },
    {
        name: "Database Service",
        slug: "database",
        url: "mongodb://localhost:27017",
        enabled: true,
        created_at: new Date(),
        updated_at: new Date()
    }
];

db.services.insertMany(services);

// Generate sample status logs
const statusOptions = ['operational', 'degraded', 'down'];
const logs = [];

services.forEach(service => {
    for (let i = 0; i < 50; i++) {
        const timestamp = new Date();
        timestamp.setHours(timestamp.getHours() - i);
        
        logs.push({
            service_name: service.name,
            status: statusOptions[Math.floor(Math.random() * statusOptions.length)],
            latency_ms: Math.floor(Math.random() * 1000) + 50,
            status_code: 200,
            timestamp: timestamp,
            created_at: timestamp
        });
    }
});

db.status_logs.insertMany(logs);

// Insert sample incident
db.incidents.insertOne({
    id: "INC-001",
    title: "API Gateway Intermittent Timeouts",
    description: "Users experiencing timeout errors when accessing the API",
    severity: "major",
    status: "investigating",
    affected_services: ["API Gateway"],
    created_at: new Date(Date.now() - 2 * 60 * 60 * 1000),
    updated_at: new Date()
});

// Insert sample maintenance
db.maintenance.insertOne({
    id: "MAINT-001",
    title: "Database Upgrade",
    description: "Upgrading database to latest version",
    scheduled_start: new Date(Date.now() + 24 * 60 * 60 * 1000),
    scheduled_end: new Date(Date.now() + 26 * 60 * 60 * 1000),
    impact: "Some API endpoints may experience brief interruptions",
    status: "scheduled",
    created_at: new Date(),
    updated_at: new Date()
});

print("✓ Seeding completed");
print("Services: " + db.services.count());
print("Status logs: " + db.status_logs.count());
print("Incidents: " + db.incidents.count());
print("Maintenance: " + db.maintenance.count());
EOF

echo "✓ Database seeded successfully"