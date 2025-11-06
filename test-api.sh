#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

# Base URLs
APP_URL="http://localhost:8080"
STORAGE_URL="http://localhost:8081"

# Test file details
FILE="test.png"
FILE_SIZE=$(stat -f%z "$FILE")
CONTENT_TYPE="image/png"

echo -e "${BLUE}Starting API Tests${NC}"
echo "----------------------------------------"
echo "File: $FILE"
echo "Size: $FILE_SIZE bytes"
echo "Content-Type: $CONTENT_TYPE"
echo "----------------------------------------"

# Function to check if services are up
check_services() {
    echo -e "\n${BLUE}Checking if services are up...${NC}"
    timeout=30
    count=0
    
    while [ $count -lt $timeout ]; do
        if curl -s "http://localhost:8080" > /dev/null; then
            if curl -s "http://localhost:8081" > /dev/null; then
                echo -e "${GREEN}Both services are up!${NC}"
                return 0
            fi
        fi
        echo "Waiting for services to start... ($count/$timeout)"
        sleep 1
        count=$((count + 1))
    done
    
    echo -e "${RED}Timeout waiting for services${NC}"
    return 1
}

# Wait for services
check_services


# Get a test JWT token from the app service
echo -e "\n${BLUE}Requesting test JWT token from /auth/token${NC}"
TOKEN_RESPONSE=$(curl -s -X POST "${APP_URL}/auth/token")
JWT_TOKEN=$(echo "$TOKEN_RESPONSE" | jq -r '.token')
if [ -z "$JWT_TOKEN" ] || [ "$JWT_TOKEN" = "null" ]; then
    echo -e "${RED}Failed to obtain JWT token: $TOKEN_RESPONSE${NC}"
    exit 1
fi
echo -e "${GREEN}Got token${NC}"

# Test 1: Get Upload URL
echo -e "\n${BLUE}Test 1: Requesting Upload URL${NC}"
UPLOAD_RESPONSE=$(curl -s -X POST "${APP_URL}/upload-url" \
    -H "Authorization: Bearer ${JWT_TOKEN}" \
    -H "Content-Type: application/json" \
    -d "{
        \"filename\": \"${FILE}\",
        \"size\": ${FILE_SIZE},
        \"content_type\": \"${CONTENT_TYPE}\"
    }")

echo "Upload URL Response: $UPLOAD_RESPONSE"

# Extract values from response
PRESIGNED_URL=$(echo $UPLOAD_RESPONSE | jq -r '.upload_url')
IMAGE_ID=$(echo $UPLOAD_RESPONSE | jq -r '.image_id')

if [ "$PRESIGNED_URL" != "null" ] && [ "$PRESIGNED_URL" != "" ]; then
    echo -e "${GREEN}✓ Got presigned URL successfully${NC}"
    
    # Test 2: Upload File
    echo -e "\n${BLUE}Test 2: Uploading File${NC}"
    UPLOAD_RESULT=$(curl -s -X PUT "$PRESIGNED_URL" \
        -H "Content-Type: ${CONTENT_TYPE}" \
        --upload-file "${FILE}")
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ File uploaded successfully${NC}"
        
        # Wait for processing
        echo "Waiting for upload processing..."
        sleep 3
        
        # Test 3: Verify Upload
        echo -e "\n${BLUE}Test 3: Verifying Upload${NC}"
        VERIFY_RESPONSE=$(curl -s "${STORAGE_URL}/verify/${IMAGE_ID}")
        echo "Verify Response: $VERIFY_RESPONSE"
        
        # Check boolean `valid` field in verify response
        VERIFY_OK=$(echo "$VERIFY_RESPONSE" | jq -r '.valid')
        if [ "$VERIFY_OK" = "true" ]; then
            echo -e "${GREEN}✓ Upload verified successfully${NC}"
            
            # Test 4: Create Product
            echo -e "\n${BLUE}Test 4: Creating Product${NC}"
            PRODUCT_RESPONSE=$(curl -s -X POST "${APP_URL}/products" \
                -H "Content-Type: application/json" \
                -d "{
                    \"name\": \"Test Product\",
                    \"description\": \"A test product\",
                    \"price\": 99.99,
                    \"image_id\": \"${IMAGE_ID}\"
                }")
            
            echo "Product Response: $PRODUCT_RESPONSE"
            if echo "$PRODUCT_RESPONSE" | grep -q "id"; then
                echo -e "${GREEN}✓ Product created successfully${NC}"
                echo -e "\n${GREEN}All tests passed successfully!${NC}"
            else
                echo -e "${RED}✗ Failed to create product${NC}"
            fi
        else
            echo -e "${RED}✗ Upload verification failed${NC}"
        fi
    else
        echo -e "${RED}✗ File upload failed${NC}"
    fi
else
    echo -e "${RED}✗ Failed to get presigned URL${NC}"
fi