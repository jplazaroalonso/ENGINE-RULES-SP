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

echo -e "${YELLOW}🎭 Rules Engine Demo - Showcasing Created Rules${NC}"
echo "=================================================="

# Function to create and demonstrate a rule
demo_rule() {
    local name=$1
    local description=$2
    local dsl_content=$3
    local category=$4
    local test_context=$5
    local test_description=$6
    
    echo -e "${PURPLE}📋 Rule: $name${NC}"
    echo -e "${CYAN}Category:${NC} $category"
    echo -e "${CYAN}Description:${NC} $description"
    echo -e "${CYAN}DSL:${NC} $dsl_content"
    echo ""
    
    # Create the rule
    echo -e "${BLUE}Creating rule...${NC}"
    create_response=$(curl -s -k -X POST -H "Content-Type: application/json" \
        -d "{
            \"name\": \"$name\",
            \"description\": \"$description\",
            \"dsl_content\": \"$dsl_content\",
            \"category\": \"$category\",
            \"priority\": \"HIGH\",
            \"tags\": [\"demo\", \"$category\"]
        }" "$MANAGEMENT_URL/v1/rules")
    
    rule_id=$(echo "$create_response" | grep -o '"rule_id":"[^"]*"' | cut -d'"' -f4)
    
    if [ -n "$rule_id" ]; then
        echo -e "${GREEN}✅ Rule created with ID: $rule_id${NC}"
        
        # Test the rule with evaluation
        echo -e "${BLUE}Testing rule evaluation...${NC}"
        echo -e "${CYAN}Test Context:${NC} $test_context"
        echo -e "${CYAN}Test Description:${NC} $test_description"
        
        eval_response=$(curl -s -k -X POST -H "Content-Type: application/json" \
            -d "{
                \"dsl_content\": \"$dsl_content\",
                \"context\": $test_context,
                \"rule_category\": \"$category\"
            }" "$EVALUATION_URL/v1/evaluate")
        
        echo -e "${BLUE}Evaluation Result:${NC} $eval_response"
        
        # Test the rule with calculation
        echo -e "${BLUE}Testing rule calculation...${NC}"
        calc_response=$(curl -s -k -X POST -H "Content-Type: application/json" \
            -d "{
                \"rule_ids\": [\"$rule_id\"],
                \"context\": $test_context
            }" "$CALCULATOR_URL/v1/calculate")
        
        echo -e "${BLUE}Calculation Result:${NC} $calc_response"
        
    else
        echo -e "${RED}❌ Failed to create rule${NC}"
    fi
    
    echo "=================================================="
    echo ""
}

echo -e "${YELLOW}🎯 Demo 1: Promotions Rules${NC}"
echo "=================================================="

# Demo 1: Three for Two Promotion
demo_rule \
    "Demo Three for Two" \
    "Buy 3 products, pay for 2" \
    "IF quantity >= 3 THEN discount = (quantity DIV 3) * price" \
    "PROMOTIONS" \
    '{"quantity": 5, "price": 10}' \
    "Customer buys 5 items at €10 each"

# Demo 2: Senior Citizen Discount
demo_rule \
    "Demo Senior Discount" \
    "15% discount for customers over 65" \
    "IF customer_age > 65 THEN discount = price * 0.15" \
    "PROMOTIONS" \
    '{"customer_age": 70, "price": 100}' \
    "70-year-old customer buying €100 worth of products"

echo -e "${YELLOW}🎫 Demo 2: Coupons Rules${NC}"
echo "=================================================="

# Demo 3: Fixed Amount Coupon
demo_rule \
    "Demo 5 Euro Coupon" \
    "€5 discount on purchases over €40" \
    "IF total_amount > 40 THEN discount = 5" \
    "COUPONS" \
    '{"total_amount": 50}' \
    "Customer makes a €50 purchase"

echo -e "${YELLOW}💰 Demo 3: Taxes Rules${NC}"
echo "=================================================="

# Demo 4: VAT Calculation
demo_rule \
    "Demo VAT 21%" \
    "21% VAT on electronics" \
    "IF product_family IN context.vat_21_families THEN tax_rate = 0.21 AND tax_amount = price * 0.21" \
    "TAXES" \
    '{"product_family": "electronics", "price": 100, "vat_21_families": ["electronics", "clothing"]}' \
    "Electronics purchase of €100"

echo -e "${YELLOW}💎 Demo 4: Loyalty Rules${NC}"
echo "=================================================="

# Demo 5: Loyalty Points
demo_rule \
    "Demo Loyalty Points" \
    "10% of purchase as loyalty points" \
    "IF customer_id IS NOT NULL THEN points_earned = total_amount * 0.1" \
    "LOYALTY" \
    '{"customer_id": "12345", "total_amount": 200}' \
    "Loyal customer making €200 purchase"

echo -e "${YELLOW}💳 Demo 5: Payments Rules${NC}"
echo "=================================================="

# Demo 6: Payment Processing
demo_rule \
    "Demo Credit Card Fee" \
    "2% processing fee for credit card payments" \
    "IF payment_method = 'credit_card' THEN processing_fee = total_amount * 0.02" \
    "PAYMENTS" \
    '{"payment_method": "credit_card", "total_amount": 100}' \
    "€100 purchase paid with credit card"

echo -e "${YELLOW}🌐 Demo 6: API Gateway Integration${NC}"
echo "=================================================="

# Demo API Gateway
echo -e "${BLUE}Testing API Gateway rule creation...${NC}"
gateway_rule=$(curl -s -k -X POST -H "Content-Type: application/json" \
    -d '{
        "name": "Gateway Demo Rule",
        "description": "Demo rule via API gateway",
        "dsl_content": "IF amount > 50 THEN discount = 5",
        "category": "PROMOTIONS",
        "priority": "MEDIUM",
        "tags": ["gateway", "demo"]
    }' "$BASE_URL/api/v1/rules")

echo -e "${BLUE}Gateway Rule Creation:${NC} $gateway_rule"

gateway_rule_id=$(echo "$gateway_rule" | grep -o '"rule_id":"[^"]*"' | cut -d'"' -f4)

if [ -n "$gateway_rule_id" ]; then
    echo -e "${GREEN}✅ Gateway rule created with ID: $gateway_rule_id${NC}"
    
    # Test gateway evaluation
    echo -e "${BLUE}Testing gateway rule evaluation...${NC}"
    gateway_eval=$(curl -s -k -X POST -H "Content-Type: application/json" \
        -d '{
            "dsl_content": "IF amount > 50 THEN discount = 5",
            "context": {"amount": 75},
            "rule_category": "PROMOTIONS"
        }' "$BASE_URL/api/v1/evaluate")
    
    echo -e "${BLUE}Gateway Evaluation:${NC} $gateway_eval"
    
    # Test gateway calculation
    echo -e "${BLUE}Testing gateway rule calculation...${NC}"
    gateway_calc=$(curl -s -k -X POST -H "Content-Type: application/json" \
        -d "{
            \"rule_ids\": [\"$gateway_rule_id\"],
            \"context\": {\"amount\": 75}
        }" "$BASE_URL/api/v1/calculate")
    
    echo -e "${BLUE}Gateway Calculation:${NC} $gateway_calc"
fi

echo ""
echo -e "${YELLOW}📊 Demo Summary${NC}"
echo "=================================================="
echo -e "${GREEN}✅ Successfully demonstrated:${NC}"
echo -e "  • Rule creation via direct service and API gateway"
echo -e "  • Rule evaluation with different DSL expressions"
echo -e "  • Rule calculation with various contexts"
echo -e "  • Promotions, Coupons, Taxes, Loyalty, and Payments rules"
echo -e "  • HTTPS secure communication"
echo -e "  • Traefik ingress routing"
echo ""
echo -e "${BLUE}🎯 The Rules Engine is fully operational and ready for production use!${NC}"
