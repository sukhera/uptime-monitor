#!/bin/bash

# Backend Tasks Command Script
# This script helps assign the backend-expert agent to work on implementation tasks

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}🔧 Backend Expert Agent Task Assignment${NC}"
echo "============================================="

# Check if implementation doc exists
if [ ! -f "docs/IMPLMENTATION.md" ]; then
    echo -e "${YELLOW}⚠️  docs/IMPLMENTATION.md not found!${NC}"
    exit 1
fi

# Check if backend expert agent exists
if [ ! -f ".claude/agent/backend-expert.md" ]; then
    echo -e "${YELLOW}⚠️  .claude/agent/backend-expert.md not found!${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Found backend expert agent and implementation document${NC}"
echo ""

# Show available backend tasks from the implementation doc
echo -e "${BLUE}📋 Available Backend Tasks:${NC}"
echo "1. Integration Endpoints (Webhooks, Manual Status, Bulk Import)"
echo "2. Real-time Updates (SSE Implementation)" 
echo "3. YAML Configuration System"
echo "4. Data Model Extensions"
echo "5. Enhanced Checker Engine"
echo "6. Security & Middleware"
echo ""

# Command examples
echo -e "${BLUE}💡 Usage Examples:${NC}"
echo ""
echo "# Assign backend expert to work on webhook endpoints:"
echo "claude --agent .claude/agent/backend-expert.md 'Please implement the webhook endpoints listed in docs/IMPLMENTATION.md. Focus on POST /api/webhook/{service-id} and the webhook handler functionality.'"
echo ""
echo "# Assign backend expert to work on SSE real-time updates:"
echo "claude --agent .claude/agent/backend-expert.md 'Please implement the Server-Sent Events endpoints from docs/IMPLMENTATION.md. Create the SSE handler and real-time status broadcasting.'"
echo ""
echo "# Assign backend expert to work on YAML configuration:"
echo "claude --agent .claude/agent/backend-expert.md 'Please implement the YAML configuration system described in docs/IMPLMENTATION.md. Create the YAML parser and hot-reload functionality.'"
echo ""
echo "# Assign backend expert to work on data model extensions:"
echo "claude --agent .claude/agent/backend-expert.md 'Please extend the Service entity with the new fields listed in docs/IMPLMENTATION.md (webhook_url, webhook_secret, manual_status, etc.).'"
echo ""
echo "# General backend implementation:"
echo "claude --agent .claude/agent/backend-expert.md 'Please review docs/IMPLMENTATION.md and work on the backend checklist items. Start with the most critical integration endpoints.'"
echo ""

echo -e "${YELLOW}📝 Quick Copy Commands:${NC}"
echo "============================================="

# Create quick copy commands
cat << 'EOF'

# Copy and paste these commands:

# 1. Webhook Implementation
claude --agent .claude/agent/backend-expert.md "Please implement the webhook endpoints from docs/IMPLMENTATION.md. Create internal/application/handlers/webhook.go with POST /api/webhook/{service-id} endpoint and webhook authentication middleware."

# 2. SSE Real-time Updates  
claude --agent .claude/agent/backend-expert.md "Please implement Server-Sent Events from docs/IMPLMENTATION.md. Create internal/application/handlers/sse.go with GET /api/sse/status endpoint and real-time status broadcasting."

# 3. YAML Configuration System
claude --agent .claude/agent/backend-expert.md "Please implement YAML configuration system from docs/IMPLMENTATION.md. Create internal/config/yaml.go with parser, validator, and hot-reload functionality."

# 4. Data Model Extensions
claude --agent .claude/agent/backend-expert.md "Please extend the Service entity with fields from docs/IMPLMENTATION.md: webhook_url, webhook_secret, manual_status, manual_reason, service_type, integration_metadata."

# 5. Integration Handler
claude --agent .claude/agent/backend-expert.md "Please create internal/application/handlers/integration.go from docs/IMPLMENTATION.md. Implement manual status override, integration details, and bulk import endpoints."

EOF

echo ""
echo -e "${GREEN}🚀 Ready to assign backend tasks to the expert agent!${NC}"