// MongoDB seed script for status page
db = db.getSiblingDB('statuspage');

// Create collections
db.createCollection('services');
db.createCollection('status_logs');

// Insert sample services
db.services.insertMany([
  {
    name: "Auth API",
    slug: "auth-api",
    url: "https://httpstat.us/200",
    headers: {},
    expected_status: 200,
    enabled: true
  },
  {
    name: "Website",
    slug: "website", 
    url: "https://httpstat.us/500",
    headers: {},
    expected_status: 200,
    enabled: true
  },
  {
    name: "Database",
    slug: "database",
    url: "https://httpstat.us/200",
    headers: {},
    expected_status: 200,
    enabled: true
  }
]);

// Create indexes for better performance
db.status_logs.createIndex({ "service_name": 1, "timestamp": -1 });
db.services.createIndex({ "enabled": 1 });

print("Database seeded successfully!"); 