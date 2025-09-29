#!/bin/bash

# Script to create a new Chama group via the API
# Usage: ./create-group.sh [group_name] [description] [wallet_address]

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default API endpoint
API_URL="http://localhost:3000/api/groups"

# Function to display usage
show_usage() {
    echo -e "${BLUE}Usage: $0 [group_name] [description] [wallet_address]${NC}"
    echo -e "${BLUE}   or: $0 --interactive${NC}"
    echo ""
    echo -e "${YELLOW}Examples:${NC}"
    echo "  $0 \"Sar Queens\" \"Women-led savings group\" \"GA...WALLET\""
    echo "  $0 --interactive"
    echo ""
    exit 1
}

# Function to validate Stellar address
validate_stellar_address() {
    local address=$1
    if [[ ! $address =~ ^G[A-Z2-7]{55}$ ]]; then
        echo -e "${RED}❌ Invalid Stellar address format${NC}"
        echo -e "${YELLOW}Stellar addresses should start with 'G' and be 56 characters long${NC}"
        return 1
    fi
    return 0
}

# Function to create group via API
create_group() {
    local name="$1"
    local description="$2"
    local wallet="$3"

    echo -e "${BLUE}Creating group: ${name}${NC}"
    echo -e "${BLUE}Description: ${description}${NC}"
    echo -e "${BLUE}Wallet: ${wallet}${NC}"
    echo ""

    # Create JSON payload
    local json_payload=$(cat <<EOF
{
    "name": "$name",
    "description": "$description",
    "wallet": "$wallet"
}
EOF
)

    # Make API request
    echo -e "${YELLOW}Sending request to API...${NC}"
    local response=$(curl -s -X POST "$API_URL" \
        -H "Content-Type: application/json" \
        -d "$json_payload" \
        -w "\nHTTP_STATUS:%{http_code}")

    # Parse response and status
    local http_status=$(echo "$response" | grep "HTTP_STATUS" | cut -d: -f2)
    local json_response=$(echo "$response" | sed '/HTTP_STATUS/d')

    # Check if request was successful
    if [ "$http_status" = "200" ] || [ "$http_status" = "201" ]; then
        echo -e "${GREEN}✅ Group created successfully!${NC}"
        echo ""
        
        # Parse and display response using jq if available
        if command -v jq &> /dev/null; then
            echo -e "${GREEN}📋 Group Details:${NC}"
            echo "$json_response" | jq -r '
                "🆔 ID: " + .group.ID + 
                "\n📛 Name: " + .group.Name + 
                "\n📝 Description: " + .group.Description + 
                "\n💳 Wallet: " + .group.Wallet + 
                "\n🔗 Contract ID: " + .group.ContractID'
      