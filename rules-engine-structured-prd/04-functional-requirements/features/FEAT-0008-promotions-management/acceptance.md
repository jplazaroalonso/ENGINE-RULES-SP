# Acceptance Criteria - Promotions Management

## AC-001: Campaign Creation and Configuration

### Scenario: Create percentage discount campaign
**Given** a marketing manager wants to create a new promotional campaign  
**When** they configure a 20% discount for Gold tier customers  
**Then** the campaign should be created with correct targeting and discount rules  
**And** the campaign should require approval before activation

### Scenario: Configure buy-X-get-Y promotion
**Given** a marketing manager creates a "Buy 2 Get 1 Free" campaign  
**When** they configure the promotion for specific product categories  
**Then** the campaign should automatically apply to qualifying cart combinations  
**And** the lowest-priced item should be discounted

### Scenario: Set campaign budget limits
**Given** a campaign with $10,000 budget allocation  
**When** promotional applications approach the budget limit  
**Then** the system should track spending in real-time  
**And** automatically pause the campaign when budget is exhausted

## AC-002: Discount Application and Calculation

### Scenario: Apply percentage discount correctly
**Given** a customer qualifies for a 15% discount promotion  
**When** they purchase $100 worth of eligible items  
**Then** they should receive $15 discount  
**And** the final amount should be $85

### Scenario: Apply fixed amount discount with minimum purchase
**Given** a "$10 off orders over $50" promotion  
**When** a customer's cart totals $60  
**Then** they should receive $10 discount  
**And** pay $50 final amount

### Scenario: Enforce maximum discount limits
**Given** a "20% off, max $25" promotion  
**When** a customer purchases $200 worth of items  
**Then** they should receive $25 discount (not $40)  
**And** pay $175 final amount

## AC-003: Customer Targeting and Eligibility

### Scenario: Target new customers only
**Given** a "20% off first purchase" campaign for new customers  
**When** an existing customer attempts to use the promotion  
**Then** they should be ineligible for the discount  
**And** receive a clear explanation

### Scenario: Geographic targeting
**Given** a promotion limited to California customers  
**When** a customer from New York tries to apply the promotion  
**Then** they should be ineligible based on location  
**And** see location-specific messaging

### Scenario: Tier-based promotion eligibility
**Given** a "Platinum tier exclusive" promotion  
**When** a Gold tier customer views the promotion  
**Then** they should see tier upgrade messaging  
**And** understand how to qualify

## AC-004: Promotional Conflict Resolution

### Scenario: Apply best promotion for customer
**Given** a customer qualifies for both "15% off" and "$20 off $100+"  
**When** their cart total is $120  
**Then** they should receive the "$20 off" promotion  
**And** see explanation of applied discount

### Scenario: Handle stackable promotions
**Given** stackable promotions: "10% off" + "Free shipping"  
**When** a customer qualifies for both  
**Then** both promotions should be applied  
**And** total benefit should be clearly shown

### Scenario: Priority-based conflict resolution
**Given** conflicting promotions with different priorities  
**When** both promotions apply to the same transaction  
**Then** the higher priority promotion should be applied  
**And** the reason should be documented

## AC-005: Performance and Scalability

### Scenario: High-volume promotion processing
**Given** 1,000 simultaneous customers applying promotions  
**When** promotions are evaluated and applied  
**Then** all applications should complete within 5 seconds  
**And** calculation accuracy should be 100%

### Scenario: Real-time budget tracking
**Given** multiple concurrent promotion applications  
**When** budget spending is tracked across applications  
**Then** budget should be accurately maintained in real-time  
**And** no over-spending should occur

## AC-006: Campaign Analytics and Reporting

### Scenario: Track campaign performance metrics
**Given** an active promotional campaign  
**When** customers interact with the promotion  
**Then** metrics should track application rate, conversion rate, and revenue impact  
**And** data should be available in real-time dashboards

### Scenario: Calculate campaign ROI
**Given** a campaign with $5,000 budget generating $25,000 additional revenue  
**When** ROI is calculated  
**Then** ROI should be 400% ($20,000 net / $5,000 investment)  
**And** all costs should be included in calculation

## AC-007: Mobile and Multi-Channel Experience

### Scenario: Mobile promotion application
**Given** a customer using the mobile app  
**When** they apply a promotional code  
**Then** the promotion should work seamlessly on mobile  
**And** savings should be clearly displayed

### Scenario: Cross-channel promotion consistency
**Given** a promotion available across web and mobile  
**When** customers use different channels  
**Then** promotion terms and benefits should be identical  
**And** usage should be tracked across channels

## AC-008: Error Handling and Edge Cases

### Scenario: Invalid promotional code
**Given** a customer enters an invalid promotional code  
**When** they attempt to apply it  
**Then** they should receive a clear error message  
**And** suggestions for valid promotions should be offered

### Scenario: Expired promotion handling
**Given** a promotion that has expired  
**When** a customer tries to apply it  
**Then** they should receive expiration notification  
**And** alternative active promotions should be suggested

### Scenario: Insufficient inventory for promotion
**Given** a promotion requiring specific products  
**When** those products are out of stock  
**Then** the promotion should be marked as unavailable  
**And** alternative promotions should be suggested
