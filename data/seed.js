// MongoDB seed script for status page
db = db.getSiblingDB('statuspage');

// Create collections
db.createCollection('services');
db.createCollection('status_logs');

// Insert sample services
db.services.insertMany([
  {
    name: "API Server",
    slug: "api-server",
    url: "http://api:8080/api/health",
    headers: {},
    expected_status: 200,
    enabled: true
  },
  {
    name: "Web Dashboard",
    slug: "web-dashboard",
    url: "http://web/",
    headers: {},
    expected_status: 200,
    enabled: true
  },
  {
    name: "MongoDB",
    slug: "mongodb",
    url: "http://mongo:27017",
    headers: {},
    expected_status: 200,
    enabled: true
  }
]);

// Create indexes for better performance
db.status_logs.createIndex({ "service_name": 1, "timestamp": -1 });
db.services.createIndex({ "enabled": 1 });

print("Database seeded successfully!"); 