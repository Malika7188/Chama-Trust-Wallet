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

# Check if API is running
echo -e "${BLUE}🔍 Checking API availability...${NC}"
if ! curl -s --max-time 3 "$API_URL" > /dev/null; then
    echo -e "${RED}❌ API is not responding. Please start your backend server.${NC}"
    exit 1
fi
echo -e "${GREEN}✅ API is running${NC}"
echo ""

# Step 1: Register/Login test user
echo -e "${YELLOW}Step 1: User Authentication${NC}"
echo "Attempting to register test user..."

REGISTER_DATA="{
    \"name\": \"$TEST_USER_NAME\",
    \"email\": \"$TEST_USER_EMAIL\",
    \"password\": \"$TEST_USER_PASSWORD\"
}"

REGISTER_RESPONSE=$(curl -s -X POST "$API_URL/auth/register" \
    -H "Content-Type: application/json" \
    -d "$REGISTER_DATA" \
    -w "\nHTTP_STATUS:%{http_code}")

REGISTER_STATUS=$(echo "$REGISTER_RESPONSE" | grep "HTTP_STATUS" | cut -d: -f2)
REGISTER_JSON=$(echo "$REGISTER_RESPONSE" | sed '/HTTP_STATUS/d')

if [[ "$REGISTER_STATUS" == "409" ]] || [[ "$REGISTER_STATUS" == "400" ]]; then
    echo "User already exists, attempting login..."
    
    LOGIN_DATA="{
        \"email\": \"$TEST_USER_EMAIL\",
        \"password\": \"$TEST_USER_PASSWORD\"
    }"
    
    LOGIN_RESPONSE=$(curl -s -X POST "$API_URL/auth/login" \
        -H "Content-Type: application/json" \
        -d "$LOGIN_DATA")
    
    TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.token // empty')
    USER_ID=$(echo "$LOGIN_RESPONSE" | jq -r '.user.id // empty')
    USER_WALLET=$(echo "$LOGIN_RESPONSE" | jq -r '.user.wallet // empty')
    USER_SECRET=$(echo "$LOGIN_RESPONSE" | jq -r '.user.secret_key // empty')
elif [[ "$REGISTER_STATUS" =~ ^2[0-9]{2}$ ]]; then
    echo -e "${GREEN}✅ User registered successfully${NC}"
    TOKEN=$(echo "$REGISTER_JSON" | jq -r '.token // empty')
    USER_ID=$(echo "$REGISTER_JSON" | jq -r '.user.id // empty')
    USER_WALLET=$(echo "$REGISTER_JSON" | jq -r '.user.wallet // empty')
    USER_SECRET=$(echo "$REGISTER_JSON" | jq -r '.user.secret_key // empty')
else
    echo -e "${RED}❌ Failed to register/login user${NC}"
    echo "$REGISTER_JSON"
    exit 1
fi

if [[ -z "$TOKEN" ]]; then
    echo -e "${RED}❌ Failed to get authentication token${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Authentication successful${NC}"
echo "User ID: $USER_ID"
echo "Wallet: $USER_WALLET"
echo ""

# Step 2: Fund user account
echo -e "${YELLOW}Step 2: Fund User Account${NC}"
FUND_RESPONSE=$(curl -s -X POST "$API_URL/fund/$USER_WALLET" \
    -H "Authorization: Bearer $TOKEN")
echo "$FUND_RESPONSE" | jq .
echo ""

# Wait for funding
sleep 3

# Step 3: Check user balance
echo -e "${YELLOW}Step 3: Check User Balance${NC}"
USER_BALANCE_RESPONSE=$(curl -s "$API_URL/balance/$USER_WALLET" \
