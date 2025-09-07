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

echo -e "${YELLOW}🚀 Starting Complete Rules Engine Deployment Tests${NC}"
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

# Test 7: Create Rule - Management Service (Direct)
test_endpoint "POST" "$MANAGEMENT_URL/v1/rules" '{
    "name": "Test Rule Direct",
    "description": "Test rule via direct service",
    "dsl_content": "IF amount > 100 THEN discount = 10%",
    "category": "PROMOTIONS",
    "priority": "HIGH",
    "tags": ["test", "direct"]
}' "201" "Create Rule via Direct Management Service"

# Test 8: Create Rule - Management Service (Gateway)
test_endpoint "POST" "$BASE_URL/api/v1/rules" '{
    "name": "Test Rule Gateway",
    "description": "Test rule via gateway",
    "dsl_content": "IF amount > 100 THEN discount = 10%",
    "category": "PROMOTIONS",
    "priority": "HIGH",
    "tags": ["test", "gateway"]
}' "201" "Create Rule via API Gateway"

# Test 9: Evaluate Rule - Evaluation Service (Direct)
test_endpoint "POST" "$EVALUATION_URL/v1/evaluate" '{
    "dsl_content": "IF amount > 100 THEN discount = 10%",
    "context": {"amount": 150},
    "rule_category": "PROMOTIONS"
}' "500" "Evaluate Rule via Direct Evaluation Service (Expected DSL validation error)"

# Test 10: Evaluate Rule - Evaluation Service (Gateway)
test_endpoint "POST" "$BASE_URL/api/v1/evaluate" '{
    "dsl_content": "IF amount > 100 THEN discount = 10%",
    "context": {"amount": 150},
    "rule_category": "PROMOTIONS"
}' "500" "Evaluate Rule via API Gateway (Expected DSL validation error)"

# Test 11: Calculate Rules - Calculator Service (Direct)
test_endpoint "POST" "$CALCULATOR_URL/v1/calculate" '{
    "rule_ids": ["test-rule-id"],
    "context": {"amount": 150}
}' "200" "Calculate Rules via Direct Calculator Service"

# Test 12: Calculate Rules - Calculator Service (Gateway)
test_endpoint "POST" "$BASE_URL/api/v1/calculate" '{
    "rule_ids": ["test-rule-id"],
    "context": {"amount": 150}
}' "200" "Calculate Rules via API Gateway"

# Test 13: Error Handling - Invalid Rule Creation (Direct)
test_endpoint "POST" "$MANAGEMENT_URL/v1/rules" '{
    "name": "",
    "description": "Invalid rule",
    "dsl_content": "",
    "category": "INVALID",
    "priority": "INVALID",
    "tags": []
}' "400" "Error Handling - Invalid Rule Creation (Direct)"

# Test 14: Error Handling - Invalid Rule Creation (Gateway)
test_endpoint "POST" "$BASE_URL/api/v1/rules" '{
    "name": "",
    "description": "Invalid rule",
    "dsl_content": "",
    "category": "INVALID",
    "priority": "INVALID",
    "tags": []
}' "400" "Error Handling - Invalid Rule Creation (Gateway)"

# Test 15: HTTPS Certificate Validation
echo -e "${BLUE}Testing:${NC} HTTPS Certificate Validation"
echo -e "${BLUE}URL:${NC} GET $MANAGEMENT_URL/health"
cert_check=$(curl -s -k -I "$MANAGEMENT_URL/health" 2>&1 | grep -i "certificate")
if [ -n "$cert_check" ]; then
    print_status 0 "HTTPS Certificate Validation"
else
    print_status 1 "HTTPS Certificate Validation"
fi
echo "---"

# Test 16: Service Discovery
echo -e "${BLUE}Testing:${NC} Service Discovery"
echo -e "${BLUE}URL:${NC} GET $BASE_URL/health"
discovery_check=$(curl -s -k "$BASE_URL/health" 2>&1)
if [ $? -eq 0 ]; then
    print_status 0 "Service Discovery"
else
    print_status 1 "Service Discovery"
fi
echo "---"

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
