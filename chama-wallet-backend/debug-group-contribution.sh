#!/bin/bash

# Debug script for group contribution issues
# Usage: ./debug-group-contribution.sh

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
API_URL="http://localhost:3000"
TEST_USER_EMAIL="test@example.com"
TEST_USER_PASSWORD="password123"
TEST_USER_NAME="Test User"

echo -e "${BLUE}🔧 Debugging Group Contribution Issues${NC}"
echo "======================================"
echo ""

