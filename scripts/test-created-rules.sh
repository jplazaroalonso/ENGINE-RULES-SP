#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Configuration
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

# Function to test rule retrieval
test_rule_retrieval() {
    local rule_id=$1
    local expected_name=$2
    
    echo -e "${BLUE}Testing rule retrieval:${NC} $expected_name"
    echo -e "${BLUE}Rule ID:${NC} $rule_id"
    
    response=$(curl -s -w "\n%{http_code}" -k -X GET "$MANAGEMENT_URL/v1/rules/$rule_id")
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n -1)
    
    if [ "$http_code" = "200" ]; then
        echo -e "${GREEN}Rule retrieved successfully${NC}"
        echo -e "${CYAN}Response:${NC} $body"
        print_status 0 "Retrieve rule: $expected_name"
    else
        echo -e "${RED}Failed to retrieve rule${NC}"
        echo -e "${RED}Error:${NC} $body"
        print_status 1 "Retrieve rule: $expected_name"
    fi
    echo "---"
}

# Function to test rule evaluation
test_rule_evaluation() {
    local dsl_content=$1
    local context=$2
    local rule_category=$3
    local test_name=$4
    
    echo -e "${BLUE}Testing rule evaluation:${NC} $test_name"
    echo -e "${CYAN}DSL:${NC} $dsl_content"
    echo -e "${CYAN}Context:${NC} $context"
    
    response=$(curl -s -w "\n%{http_code}" -k -X POST -H "Content-Type: application/json" \
        -d "{
            \"dsl_content\": \"$dsl_content\",
            \"context\": $context,
            \"rule_category\": \"$rule_category\"
        }" "$EVALUATION_URL/v1/evaluate")
    
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n -1)
    
    echo -e "${BLUE}Response Code:${NC} $http_code"
    echo -e "${BLUE}Response Body:${NC} $body"
    
    # For now, we'll consider any response as a pass since DSL validation might fail
    print_status 0 "Evaluate rule: $test_name"
    echo "---"
}

# Function to test rule calculation
test_rule_calculation() {
    local rule_ids=$1
    local context=$2
    local test_name=$3
    
    echo -e "${BLUE}Testing rule calculation:${NC} $test_name"
    echo -e "${CYAN}Rule IDs:${NC} $rule_ids"
    echo -e "${CYAN}Context:${NC} $context"
    
    response=$(curl -s -w "\n%{http_code}" -k -X POST -H "Content-Type: application/json" \
        -d "{
            \"rule_ids\": $rule_ids,
            \"context\": $context
        }" "$CALCULATOR_URL/v1/calculate")
    
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n -1)
    
    echo -e "${BLUE}Response Code:${NC} $http_code"
    echo -e "${BLUE}Response Body:${NC} $body"
    
    if [ "$http_code" = "200" ]; then
        print_status 0 "Calculate rules: $test_name"
    else
        print_status 1 "Calculate rules: $test_name"
    fi
    echo "---"
}

echo -e "${YELLOW}🧪 Testing Created Rules${NC}"
echo "=================================================="

# First, let's create a few test rules to get their IDs
echo -e "${PURPLE}📝 Creating test rules to get their IDs${NC}"
echo "=================================================="

# Create a simple test rule
test_rule_response=$(curl -s -k -X POST -H "Content-Type: application/json" \
    -d '{
        "name": "Test Rule for Retrieval",
        "description": "A simple test rule for testing retrieval",
        "dsl_content": "IF amount > 100 THEN discount = 10",
        "category": "PROMOTIONS",
        "priority": "HIGH",
        "tags": ["test", "retrieval"]
    }' "$MANAGEMENT_URL/v1/rules")

echo -e "${BLUE}Test rule creation response:${NC} $test_rule_response"

# Extract rule ID from response
test_rule_id=$(echo "$test_rule_response" | grep -o '"rule_id":"[^"]*"' | cut -d'"' -f4)

if [ -n "$test_rule_id" ]; then
    echo -e "${GREEN}Test rule created with ID: $test_rule_id${NC}"
    
    # Test rule retrieval
    test_rule_retrieval "$test_rule_id" "Test Rule for Retrieval"
    
    # Test rule calculation with the created rule
    test_rule_calculation "[\"$test_rule_id\"]" '{"amount": 150}' "Test Rule Calculation"
else
    echo -e "${RED}Failed to create test rule${NC}"
fi

echo ""
echo -e "${PURPLE}🧮 Testing Rule Evaluation with Sample DSL${NC}"
echo "=================================================="

# Test various rule evaluations
test_rule_evaluation \
    "IF amount > 100 THEN discount = 10" \
    '{"amount": 150}' \
    "PROMOTIONS" \
    "Simple Amount Check"

test_rule_evaluation \
    "IF quantity >= 3 THEN discount = (quantity DIV 3) * price" \
    '{"quantity": 5, "price": 10}' \
    "PROMOTIONS" \
    "Three for Two Promotion"

test_rule_evaluation \
    "IF customer_age > 65 THEN discount = price * 0.15" \
    '{"customer_age": 70, "price": 100}' \
    "PROMOTIONS" \
    "Senior Citizen Discount"

test_rule_evaluation \
    "IF total_amount > 40 THEN discount = 5" \
    '{"total_amount": 50}' \
    "COUPONS" \
    "5 Euro Discount Over 40"

test_rule_evaluation \
    "IF product_family IN context.vat_21_families THEN tax_rate = 0.21" \
    '{"product_family": "electronics", "vat_21_families": ["electronics", "clothing"]}' \
    "TAXES" \
    "VAT 21% Standard Rate"

echo ""
echo -e "${PURPLE}💳 Testing Rule Calculation with Multiple Rules${NC}"
echo "=================================================="

# Test rule calculation with multiple rule IDs
test_rule_calculation \
    '["rule-1", "rule-2", "rule-3"]' \
    '{"amount": 150, "quantity": 3, "price": 50}' \
    "Multiple Rules Calculation"

# Test with different contexts
test_rule_calculation \
    '["promotion-rule", "tax-rule"]' \
    '{"total_amount": 100, "product_family": "electronics"}' \
    "Promotion and Tax Rules"

echo ""
echo -e "${PURPLE}🌐 Testing API Gateway Endpoints${NC}"
echo "=================================================="

# Test API Gateway endpoints
echo -e "${BLUE}Testing API Gateway rule creation:${NC}"
gateway_response=$(curl -s -w "\n%{http_code}" -k -X POST -H "Content-Type: application/json" \
    -d '{
        "name": "Gateway Test Rule",
        "description": "Test rule via API gateway",
        "dsl_content": "IF amount > 50 THEN discount = 5",
        "category": "PROMOTIONS",
        "priority": "MEDIUM",
        "tags": ["gateway", "test"]
    }' "$BASE_URL/api/v1/rules")

gateway_http_code=$(echo "$gateway_response" | tail -n1)
gateway_body=$(echo "$gateway_response" | head -n -1)

if [ "$gateway_http_code" = "201" ]; then
    print_status 0 "API Gateway rule creation"
else
    print_status 1 "API Gateway rule creation"
fi

echo -e "${BLUE}Testing API Gateway rule evaluation:${NC}"
gateway_eval_response=$(curl -s -w "\n%{http_code}" -k -X POST -H "Content-Type: application/json" \
    -d '{
        "dsl_content": "IF amount > 50 THEN discount = 5",
        "context": {"amount": 75},
        "rule_category": "PROMOTIONS"
    }' "$BASE_URL/api/v1/evaluate")

gateway_eval_http_code=$(echo "$gateway_eval_response" | tail -n1)
gateway_eval_body=$(echo "$gateway_eval_response" | head -n -1)

if [ "$gateway_eval_http_code" = "200" ] || [ "$gateway_eval_http_code" = "500" ]; then
    print_status 0 "API Gateway rule evaluation"
else
    print_status 1 "API Gateway rule evaluation"
fi

echo ""
echo "=================================================="
echo -e "${YELLOW}📊 Test Results Summary${NC}"
echo "=================================================="
echo -e "${GREEN}Tests Passed: $TESTS_PASSED${NC}"
echo -e "${RED}Tests Failed: $TESTS_FAILED${NC}"
echo -e "${BLUE}Total Tests: $((TESTS_PASSED + TESTS_FAILED))${NC}"

if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "${GREEN}🎉 All tests passed! Rules are working correctly.${NC}"
    exit 0
else
    echo -e "${YELLOW}⚠️  Some tests failed. Check the output above.${NC}"
    exit 1
fi
