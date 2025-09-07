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

# Counters
RULES_CREATED=0
RULES_FAILED=0

# Function to print status
print_status() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✅ SUCCESS${NC}: $2"
        ((RULES_CREATED++))
    else
        echo -e "${RED}❌ FAILED${NC}: $2"
        ((RULES_FAILED++))
    fi
}

# Function to create rule
create_rule() {
    local name=$1
    local description=$2
    local dsl_content=$3
    local category=$4
    local priority=$5
    local tags=$6
    
    echo -e "${BLUE}Creating rule:${NC} $name"
    echo -e "${CYAN}Category:${NC} $category"
    echo -e "${CYAN}DSL:${NC} $dsl_content"
    
    response=$(curl -s -w "\n%{http_code}" -k -X POST -H "Content-Type: application/json" \
        -d "{
            \"name\": \"$name\",
            \"description\": \"$description\",
            \"dsl_content\": \"$dsl_content\",
            \"category\": \"$category\",
            \"priority\": \"$priority\",
            \"tags\": $tags
        }" "$MANAGEMENT_URL/v1/rules")
    
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n -1)
    
    if [ "$http_code" = "201" ]; then
        rule_id=$(echo "$body" | grep -o '"rule_id":"[^"]*"' | cut -d'"' -f4)
        echo -e "${GREEN}Rule ID:${NC} $rule_id"
        print_status 0 "$name"
    else
        echo -e "${RED}Error:${NC} $body"
        print_status 1 "$name"
    fi
    echo "---"
}

# Function to list all rules
list_rules() {
    echo -e "${YELLOW}📋 Listing all rules in the database:${NC}"
    echo "=================================================="
    
    response=$(curl -s -k "$MANAGEMENT_URL/v1/rules")
    echo "$response" | jq '.' 2>/dev/null || echo "$response"
}

echo -e "${YELLOW}🚀 Starting Rules Engine Database Population${NC}"
echo "=================================================="

# PROMOTIONS RULES
echo -e "${PURPLE}🎯 Creating PROMOTIONS Rules${NC}"
echo "=================================================="

create_rule \
    "Three for Two Promotion" \
    "When you buy three products, you pay for two. The promotion is three for two." \
    "IF quantity >= 3 THEN discount = (quantity DIV 3) * price" \
    "PROMOTIONS" \
    "HIGH" \
    '["promotion", "quantity", "three-for-two"]'

create_rule \
    "Second Unit 40% Discount" \
    "A 40% discount is applied to the second unit of a product family" \
    "IF quantity >= 2 AND product_family = context.product_family THEN discount = price * 0.4" \
    "PROMOTIONS" \
    "HIGH" \
    '["promotion", "second-unit", "product-family"]'

create_rule \
    "Canned Beer 30% Discount" \
    "30% on all canned beer products" \
    "IF product_type = 'canned_beer' THEN discount = price * 0.3" \
    "PROMOTIONS" \
    "MEDIUM" \
    '["promotion", "beer", "canned", "alcohol"]'

create_rule \
    "Meat Thursday 10% Discount" \
    "10% on meat products on Thursdays" \
    "IF product_category = 'meat' AND day_of_week = 'Thursday' THEN discount = price * 0.1" \
    "PROMOTIONS" \
    "MEDIUM" \
    '["promotion", "meat", "thursday", "weekly"]'

create_rule \
    "Senior Citizen Discount" \
    "Discounts for people over 65" \
    "IF customer_age > 65 THEN discount = price * 0.15" \
    "PROMOTIONS" \
    "HIGH" \
    '["promotion", "senior", "age-based", "social"]'

create_rule \
    "Online Tuesday Thursday 5% Discount" \
    "When you shop online on Tuesdays and Thursdays, you get a 5% discount." \
    "IF channel = 'online' AND (day_of_week = 'Tuesday' OR day_of_week = 'Thursday') THEN discount = price * 0.05" \
    "PROMOTIONS" \
    "MEDIUM" \
    '["promotion", "online", "tuesday", "thursday", "channel"]'

# COUPONS RULES
echo -e "${PURPLE}🎫 Creating COUPONS Rules${NC}"
echo "=================================================="

create_rule \
    "5 Euro Discount Over 40" \
    "€5 discount on purchases over €40." \
    "IF total_amount > 40 THEN discount = 5" \
    "COUPONS" \
    "HIGH" \
    '["coupon", "fixed-amount", "minimum-purchase"]'

create_rule \
    "Product Specific Coupon" \
    "Discount coupon when you buy a specific product." \
    "IF product_id = context.specific_product_id THEN discount = context.coupon_value" \
    "COUPONS" \
    "MEDIUM" \
    '["coupon", "product-specific", "conditional"]'

create_rule \
    "Product Set Coupon" \
    "€5 discount coupon if you buy a specific set of products." \
    "IF contains_all_products(context.required_products) THEN discount = 5" \
    "COUPONS" \
    "MEDIUM" \
    '["coupon", "product-set", "bundle"]'

# LOYALTY RULES
echo -e "${PURPLE}💎 Creating LOYALTY Rules${NC}"
echo "=================================================="

create_rule \
    "10% Purchase Points" \
    "You accumulate 10% of your purchase in points, which you can redeem the following month." \
    "IF customer_id IS NOT NULL THEN points_earned = total_amount * 0.1 AND points_valid_from = next_month" \
    "LOYALTY" \
    "HIGH" \
    '["loyalty", "points", "monthly", "redemption"]'

create_rule \
    "8% Fuel Points" \
    "8% fuel discount accumulated for the following month" \
    "IF product_category = 'fuel' THEN points_earned = total_amount * 0.08 AND points_valid_from = next_month" \
    "LOYALTY" \
    "MEDIUM" \
    '["loyalty", "fuel", "points", "monthly"]'

create_rule \
    "2% Credit Card Points" \
    "2% discount accumulated the following month on purchases made in our stores and paid for with our credit card." \
    "IF channel = 'store' AND payment_method = 'our_credit_card' THEN points_earned = total_amount * 0.02 AND points_valid_from = next_month" \
    "LOYALTY" \
    "LOW" \
    '["loyalty", "credit-card", "store", "monthly"]'

# TAXES RULES
echo -e "${PURPLE}💰 Creating TAXES Rules${NC}"
echo "=================================================="

create_rule \
    "VAT 21% Standard Rate" \
    "VAT-21 21% applies to a list of product families" \
    "IF product_family IN context.vat_21_families THEN tax_rate = 0.21 AND tax_amount = price * 0.21" \
    "TAXES" \
    "HIGH" \
    '["tax", "vat", "21%", "standard-rate"]'

create_rule \
    "VAT 10% Reduced Rate" \
    "VAT-10 10% applies to a list of product families" \
    "IF product_family IN context.vat_10_families THEN tax_rate = 0.10 AND tax_amount = price * 0.10" \
    "TAXES" \
    "HIGH" \
    '["tax", "vat", "10%", "reduced-rate"]'

create_rule \
    "VAT 4% Super Reduced Rate" \
    "VAT-5 4% applies to a list of product families" \
    "IF product_family IN context.vat_4_families THEN tax_rate = 0.04 AND tax_amount = price * 0.04" \
    "TAXES" \
    "HIGH" \
    '["tax", "vat", "4%", "super-reduced-rate"]'

create_rule \
    "VAT 0% Zero Rate" \
    "VAT-0 0% applies to a list of product families" \
    "IF product_family IN context.vat_0_families THEN tax_rate = 0.00 AND tax_amount = 0" \
    "TAXES" \
    "HIGH" \
    '["tax", "vat", "0%", "zero-rate"]'

# PAYMENTS RULES
echo -e "${PURPLE}💳 Creating PAYMENTS Rules${NC}"
echo "=================================================="

create_rule \
    "Credit Card Processing Fee" \
    "2% processing fee for credit card payments" \
    "IF payment_method = 'credit_card' THEN processing_fee = total_amount * 0.02" \
    "PAYMENTS" \
    "MEDIUM" \
    '["payment", "credit-card", "processing-fee"]'

create_rule \
    "Cash Payment Discount" \
    "1% discount for cash payments" \
    "IF payment_method = 'cash' THEN discount = total_amount * 0.01" \
    "PAYMENTS" \
    "LOW" \
    '["payment", "cash", "discount"]'

create_rule \
    "Installment Payment Fee" \
    "€5 fee for installment payments" \
    "IF payment_method = 'installment' THEN processing_fee = 5" \
    "PAYMENTS" \
    "MEDIUM" \
    '["payment", "installment", "fee"]'

# Summary
echo "=================================================="
echo -e "${YELLOW}📊 Population Summary${NC}"
echo "=================================================="
echo -e "${GREEN}Rules Created Successfully: $RULES_CREATED${NC}"
echo -e "${RED}Rules Failed: $RULES_FAILED${NC}"
echo -e "${BLUE}Total Rules Attempted: $((RULES_CREATED + RULES_FAILED))${NC}"

if [ $RULES_FAILED -eq 0 ]; then
    echo -e "${GREEN}🎉 All rules created successfully!${NC}"
else
    echo -e "${YELLOW}⚠️  Some rules failed to create. Check the output above.${NC}"
fi

echo ""
echo -e "${YELLOW}📋 Listing all rules in the database:${NC}"
echo "=================================================="

# List all rules
list_rules

echo ""
echo -e "${GREEN}✅ Database population completed!${NC}"
echo -e "${BLUE}You can now test the rules using the evaluation and calculator services.${NC}"
