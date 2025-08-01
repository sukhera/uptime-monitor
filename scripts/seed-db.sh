#!/bin/bash

set -e

MONGO_URI=${MONGO_URI:-mongodb://localhost:27017}
DB_NAME=${DB_NAME:-statuspage}

echo "Seeding database ${DB_NAME}..."

# Run the seed script
docker run --rm -v $(pwd)/data:/data --network host mongo:6 \
    mongosh "$MONGO_URI/$DB_NAME" /data/seed.js

echo "Database seeded successfully!"