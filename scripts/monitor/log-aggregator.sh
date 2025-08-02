#!/bin/bash
# Log aggregation and analysis

set -e

LOG_DIR="logs/$(date +%Y%m%d)"
mkdir -p "$LOG_DIR"

echo "Aggregating logs for $(date +%Y-%m-%d)..."

# Collect container logs
docker-compose logs --since="24h" api > "$LOG_DIR/api.log" 2>&1
docker-compose logs --since="24h" status-checker > "$LOG_DIR/status-checker.log" 2>&1
docker-compose logs --since="24h" web > "$LOG_DIR/web.log" 2>&1
docker-compose logs --since="24h" mongo > "$LOG_DIR/mongo.log" 2>&1

# Analyze logs for patterns
echo "Log Analysis Summary:" > "$LOG_DIR/analysis.txt"
echo "=====================" >> "$LOG_DIR/analysis.txt"
echo "" >> "$LOG_DIR/analysis.txt"

# Error count analysis
echo "Error Counts:" >> "$LOG_DIR/analysis.txt"
grep -c -i "error" "$LOG_DIR"/*.log | sed 's/.*\///' >> "$LOG_DIR/analysis.txt" 2>/dev/null || true
echo "" >> "$LOG_DIR/analysis.txt"

# Warning count analysis
echo "Warning Counts:" >> "$LOG_DIR/analysis.txt"
grep -c -i "warn" "$LOG_DIR"/*.log | sed 's/.*\///' >> "$LOG_DIR/analysis.txt" 2>/dev/null || true
echo "" >> "$LOG_DIR/analysis.txt"

# Most common errors
echo "Most Common Errors:" >> "$LOG_DIR/analysis.txt"
grep -h -i "error" "$LOG_DIR"/*.log | sort | uniq -c | sort -rn | head -10 >> "$LOG_DIR/analysis.txt" 2>/dev/null || true

# Compress old logs
find logs -name "*.log" -mtime +7 -exec gzip {} \;

echo "✓ Log aggregation completed: $LOG_DIR"