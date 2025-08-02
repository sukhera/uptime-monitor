#!/bin/bash
# Performance monitoring and alerting

set -e

ALERT_THRESHOLD_CPU=${ALERT_THRESHOLD_CPU:-80}
ALERT_THRESHOLD_MEM=${ALERT_THRESHOLD_MEM:-85}
ALERT_THRESHOLD_DISK=${ALERT_THRESHOLD_DISK:-90}
WEBHOOK_URL=${WEBHOOK_URL:-""}

# Get metrics
CPU_USAGE=$(docker stats --no-stream --format "{{.CPUPerc}}" | sed 's/%//' | awk '{sum+=$1} END {print sum/NR}')
MEM_USAGE=$(docker stats --no-stream --format "{{.MemPerc}}" | sed 's/%//' | awk '{sum+=$1} END {print sum/NR}')
DISK_USAGE=$(df / | tail -1 | awk '{print $5}' | sed 's/%//')

echo "Performance Metrics:"
echo "CPU Usage: ${CPU_USAGE}%"
echo "Memory Usage: ${MEM_USAGE}%"
echo "Disk Usage: ${DISK_USAGE}%"

# Check thresholds and send alerts
send_alert() {
    local message="$1"
    echo "ALERT: $message"
    
    if [ -n "$WEBHOOK_URL" ]; then
        curl -X POST "$WEBHOOK_URL" \
            -H "Content-Type: application/json" \
            -d "{\"text\":\"Status Page Alert: $message\"}"
    fi
}

if (( $(echo "$CPU_USAGE > $ALERT_THRESHOLD_CPU" | bc -l) )); then
    send_alert "High CPU usage: ${CPU_USAGE}%"
fi

if (( $(echo "$MEM_USAGE > $ALERT_THRESHOLD_MEM" | bc -l) )); then
    send_alert "High memory usage: ${MEM_USAGE}%"
fi

if [ "$DISK_USAGE" -gt "$ALERT_THRESHOLD_DISK" ]; then
    send_alert "High disk usage: ${DISK_USAGE}%"
fi

# Application-specific checks
RESPONSE_TIME=$(curl -o /dev/null -s -w "%{time_total}" http://localhost/api/health)
if (( $(echo "$RESPONSE_TIME > 2" | bc -l) )); then
    send_alert "Slow API response time: ${RESPONSE_TIME}s"
fi