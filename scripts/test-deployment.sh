#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test configuration
BASE_URL="https://rules-engine.local.dev"
MANAGEMENT_URL="https://rules-management.local.dev"
EVALUATION_URL="https://rules-evaluation.local.dev"
CALCULATOR_URL="https://rules-calculator.local.dev"

# Test results
TESTS_PASSED=0
TESTS_FAILED=0

# Function to print status
print_status() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✅ PASS${NC}: $2"
        ((TESTS_PASSED++))
    else
        echo -e "${RED}❌ FAIL${NC}: $2"
        ((TESTS_FAILED++))
    fi
}

# Function to test endpoint
test_endpoint() {
    local method=$1
    local url=$2
    local data=$3
    local expected_status=$4
    local description=$5
    
    echo -e "${BLUE}Testing:${NC} $description"
    echo -e "${BLUE}URL:${NC} $method $url"
    
    if [ -n "$data" ]; then
        response=$(curl -s -w "\n%{http_code}" -k -X $method -H "Content-Type: application/json" -d "$data" "$url")
    else
        response=$(curl -s -w "\n%{http_code}" -k -X $method "$url")
    fi
    
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n -1)
    
    echo -e "${BLUE}Response Code:${NC} $http_code"
    echo -e "${BLUE}Response Body:${NC} $body"
    
    if [ "$http_code" = "$expected_status" ]; then
        print_status 0 "$description"
    else
        print_status 1 "$description (Expected: $expected_status, Got: $http_code)"
    fi
    echo "---"
}

echo -e "${YELLOW}🚀 Starting Rules Engine Deployment Tests${NC}"
echo "=================================================="

# Test 1: Health Check - Management Service
test_endpoint "GET" "$MANAGEMENT_URL/health" "" "200" "Management Service Health Check"

# Test 2: Health Check - Evaluation Service  
test_endpoint "GET" "$EVALUATION_URL/health" "" "200" "Evaluation Service Health Check"

# Test 3: Health Check - Calculator Service
test_endpoint "GET" "$CALCULATOR_URL/health" "" "200" "Calculator Service Health Check"

# Test 4: Metrics - Management Service
test_endpoint "GET" "$MANAGEMENT_URL/metrics" "" "200" "Management Service Metrics"

# Test 5: Metrics - Evaluation Service
test_endpoint "GET" "$EVALUATION_URL/metrics" "" "200" "Evaluation Service Metrics"

# Test 6: Metrics - Calculator Service
test_endpoint "GET" "$CALCULATOR_URL/metrics" "" "200" "Calculator Service Metrics"

# Test 7: Create Rule - Management Service
test_endpoint "POST" "$MANAGEMENT_URL/v1/rules" '{
    "name": "Test Promotion Rule",
    "description": "Test rule for promotions",
    "dsl_content": "IF amount > 100 THEN discount = 10%",
    "category": "PROMOTIONS",
    "tags": ["test", "promotion"]
}' "201" "Create Rule via Management Service"

# Test 8: Get Rule - Management Service (using rule ID from previous test)
test_endpoint "GET" "$MANAGEMENT_URL/v1/rules/test-rule-id" "" "200" "Get Rule via Management Service"

# Test 9: Validate Rule - Management Service
test_endpoint "POST" "$MANAGEMENT_URL/v1/validate" '{
    "dsl_content": "IF amount > 100 THEN discount = 10%"
}' "200" "Validate Rule DSL"

# Test 10: Evaluate Rule - Evaluation Service
test_endpoint "POST" "$EVALUATION_URL/v1/evaluate" '{
    "dsl_content": "IF amount > 100 THEN discount = 10%",
    "context": {"amount": 150},
    "rule_category": "PROMOTIONS"
}' "200" "Evaluate Rule via Evaluation Service"

# Test 11: Calculate Rules - Calculator Service
test_endpoint "POST" "$CALCULATOR_URL/v1/calculate" '{
    "rule_ids": ["test-rule-id"],
    "context": {"amount": 150}
}' "200" "Calculate Rules via Calculator Service"

# Test 12: Gateway Endpoints - Management
test_endpoint "POST" "$BASE_URL/api/v1/rules" '{
    "name": "Gateway Test Rule",
    "description": "Test rule via gateway",
    "dsl_content": "IF amount > 50 THEN discount = 5%",
    "category": "PROMOTIONS",
    "tags": ["gateway", "test"]
}' "201" "Create Rule via Gateway"

# Test 13: Gateway Endpoints - Evaluation
test_endpoint "POST" "$BASE_URL/api/v1/evaluate" '{
    "dsl_content": "IF amount > 50 THEN discount = 5%",
    "context": {"amount": 75},
    "rule_category": "PROMOTIONS"
}' "200" "Evaluate Rule via Gateway"

# Test 14: Gateway Endpoints - Calculator
test_endpoint "POST" "$BASE_URL/api/v1/calculate" '{
    "rule_ids": ["gateway-test-rule-id"],
    "context": {"amount": 75}
}' "200" "Calculate Rules via Gateway"

# Test 15: Error Handling - Invalid Rule Creation
test_endpoint "POST" "$MANAGEMENT_URL/v1/rules" '{
    "name": "",
    "description": "Invalid rule",
    "dsl_content": "",
    "category": "INVALID"
}' "400" "Error Handling - Invalid Rule Creation"

# Test 16: Error Handling - Invalid Evaluation
test_endpoint "POST" "$EVALUATION_URL/v1/evaluate" '{
    "dsl_content": "INVALID DSL",
    "context": {},
    "rule_category": "INVALID"
}' "400" "Error Handling - Invalid Evaluation"

echo "=================================================="
echo -e "${YELLOW}📊 Test Results Summary${NC}"
echo "=================================================="
echo -e "${GREEN}Tests Passed: $TESTS_PASSED${NC}"
echo -e "${RED}Tests Failed: $TESTS_FAILED${NC}"
echo -e "${BLUE}Total Tests: $((TESTS_PASSED + TESTS_FAILED))${NC}"

if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "${GREEN}🎉 All tests passed! Deployment is successful.${NC}"
    exit 0
else
    echo -e "${RED}⚠️  Some tests failed. Please check the deployment.${NC}"
    exit 1
fi
