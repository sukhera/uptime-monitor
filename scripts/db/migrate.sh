#!/bin/bash
# Database migration system

set -e

MONGO_URI=${MONGO_URI:-"mongodb://localhost:27017/status_page"}
MIGRATIONS_DIR="scripts/db/migrations"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Create migrations directory
mkdir -p "$MIGRATIONS_DIR"

case "$1" in
    create)
        if [ -z "$2" ]; then
            echo -e "${RED}Usage: $0 create <migration_name>${NC}"
            exit 1
        fi
        
        TIMESTAMP=$(date +%Y%m%d_%H%M%S)
        MIGRATION_FILE="$MIGRATIONS_DIR/${TIMESTAMP}_$2.js"
        
        cat > "$MIGRATION_FILE" << EOF
// Migration: $2
// Created: $(date)

db = db.getSiblingDB('status_page');

// Up migration
function up() {
    // Add your migration code here
    print("Running migration: $2");
    
    // Example:
    // db.services.createIndex({ "slug": 1 }, { unique: true });
}

// Down migration  
function down() {
    // Add your rollback code here
    print("Rolling back migration: $2");
    
    // Example:
    // db.services.dropIndex({ "slug": 1 });
}

// Run the migration
up();
EOF
        
        echo -e "${GREEN}✓ Created migration: $MIGRATION_FILE${NC}"
        ;;
        
    up)
        echo -e "${BLUE}Running database migrations...${NC}"
        
        # Create migrations tracking collection
        mongosh "$MONGO_URI" --eval "
            db.migrations.createIndex({ filename: 1 }, { unique: true });
        " >/dev/null 2>&1
        
        for migration in "$MIGRATIONS_DIR"/*.js; do
            if [ -f "$migration" ]; then
                filename=$(basename "$migration")
                
                # Check if migration already ran
                if mongosh "$MONGO_URI" --quiet --eval "
                    db.migrations.findOne({ filename: '$filename' })
                " | grep -q "null"; then
                    echo -e "${YELLOW}Running migration: $filename${NC}"
                    mongosh "$MONGO_URI" "$migration"
                    
                    # Mark as completed
                    mongosh "$MONGO_URI" --eval "
                        db.migrations.insertOne({
                            filename: '$filename',
                            executedAt: new Date()
                        });
                    " >/dev/null 2>&1
                    
                    echo -e "${GREEN}✓ Completed: $filename${NC}"
                else
                    echo -e "${BLUE}Skipping (already ran): $filename${NC}"
                fi
            fi
        done
        ;;
        
    status)
        echo -e "${BLUE}Migration Status:${NC}"
        mongosh "$MONGO_URI" --eval "
            db.migrations.find().sort({ executedAt: 1 }).forEach(
                function(doc) {
                    print(doc.filename + ' - ' + doc.executedAt);
                }
            );
        "
        ;;
        
    *)
        echo "Usage: $0 {create|up|status} [migration_name]"
        echo ""
        echo "Commands:"
        echo "  create <name>  Create a new migration file"
        echo "  up             Run pending migrations"
        echo "  status         Show migration status"
        ;;
esac