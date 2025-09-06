# Behaviour Tests - Promotions Management

## Feature: Promotional Campaigns and Offers Management

### Background:
```gherkin
Given the promotional system is active
And the following campaign types are supported:
  | Type               | Description                          |
  | PERCENTAGE         | Percentage-based discounts          |
  | FIXED_AMOUNT       | Fixed dollar amount discounts       |
  | BUY_X_GET_Y        | Buy X items, get Y items free/discounted |
  | BUNDLE             | Product bundle special pricing       |
  | FLASH_SALE         | Time-limited promotional offers     |
And the following customer segments are configured:
  | Segment            | Criteria                            |
  | NEW_CUSTOMERS      | First-time purchasers               |
  | LOYAL_CUSTOMERS    | High lifetime value customers       |
  | TIER_BASED         | Based on loyalty tier classification |
  | GEOGRAPHIC         | Location-based targeting            |
  | BEHAVIORAL         | Purchase behavior patterns          |
```

## Scenario Outline: Campaign Creation and Configuration
```gherkin
Scenario: Marketing manager creates percentage discount campaign
Given a marketing manager wants to increase sales for Gold tier customers
When they create a new promotional campaign with the following details:
  | Field              | Value                    |
  | Campaign Name      | Gold Tier Spring Sale    |
  | Discount Type      | PERCENTAGE              |
  | Discount Value     | 20%                     |
  | Target Segment     | GOLD_TIER               |
  | Valid From         | 2024-03-01              |
  | Valid Until        | 2024-03-07              |
  | Budget Limit       | $10,000                 |
  | Product Categories | Electronics, Clothing    |
And they submit the campaign for approval
Then the campaign should be created with status "Pending Approval"
And the campaign should appear in the approval queue
And all configuration parameters should be saved correctly
And budget allocation should be recorded

Scenario: Create buy-X-get-Y promotional campaign
Given a marketing manager wants to increase average order value
When they configure a "Buy 2 Get 1 Free" promotion with:
  | Configuration      | Value                    |
  | Promotion Type     | BUY_X_GET_Y             |
  | Buy Quantity       | 2                       |
  | Get Quantity       | 1                       |
  | Get Discount       | 100%                    |
  | Eligible Categories| Electronics             |
  | Maximum Free Items | 3 per transaction       |
  | Campaign Duration  | 30 days                 |
Then the campaign should be configured successfully
And product eligibility rules should be applied
And quantity limits should be enforced
And the promotion should be ready for activation

Scenario: Configure product bundle promotion
Given a marketing manager wants to promote complementary products
When they create a bundle promotion with:
  | Bundle Component   | Individual Price | Bundle Role    |
  | Laptop Model X     | $999.99         | Primary Item   |
  | Wireless Mouse     | $29.99          | Accessory      |
  | USB Headphones     | $79.99          | Accessory      |
  | Bundle Total       | $899.99         | Special Price  |
And they set the bundle availability for "Back to School" season
Then the bundle should be created with correct pricing
And the savings calculation should show $109.98 discount
And the bundle should appear in the product catalog
And inventory requirements should be validated
```

## Scenario Outline: Customer Targeting and Eligibility
```gherkin
Scenario Outline: Customer segment targeting validation
Given a promotional campaign targeting "<target_segment>" customers
And a customer with segment classification "<customer_segment>"
When the customer views available promotions
Then the customer should <eligibility> the promotion
And the targeting reason should be "<reason>"

Examples:
| target_segment  | customer_segment | eligibility | reason                     |
| NEW_CUSTOMERS   | NEW_CUSTOMERS    | see         | Matches target segment     |
| NEW_CUSTOMERS   | LOYAL_CUSTOMERS  | not see     | Not in target segment      |
| LOYAL_CUSTOMERS | LOYAL_CUSTOMERS  | see         | Matches target segment     |
| GOLD_TIER       | GOLD_TIER        | see         | Matches loyalty tier       |
| GOLD_TIER       | SILVER_TIER      | not see     | Insufficient loyalty tier  |

Scenario: Geographic targeting enforcement
Given a promotion targeting customers in "California, New York, Florida"
And customers from different states:
  | Customer ID | State | Expected Eligibility |
  | CUST-001   | CA    | Eligible            |
  | CUST-002   | TX    | Not Eligible        |
  | CUST-003   | NY    | Eligible            |
  | CUST-004   | FL    | Eligible            |
  | CUST-005   | WA    | Not Eligible        |
When each customer checks promotion availability
Then only customers from target states should see the promotion
And customers from non-target states should see alternative promotions
And geographic targeting should be enforced accurately

Scenario: Behavioral targeting based on purchase history
Given a promotion targeting customers with:
  | Criteria               | Requirement    |
  | Minimum Orders         | 5 orders       |
  | Minimum Total Spend    | $500           |
  | Last Purchase Within   | 30 days        |
And customers with different purchase patterns:
  | Customer | Orders | Total Spend | Last Purchase Days | Expected |
  | CUST-A   | 8      | $750        | 15                | Eligible |
  | CUST-B   | 3      | $300        | 10                | Not Eligible |
  | CUST-C   | 6      | $600        | 45                | Not Eligible |
  | CUST-D   | 10     | $1200       | 5                 | Eligible |
When the behavioral targeting engine evaluates each customer
Then only customers meeting all criteria should be eligible
And behavioral metrics should be calculated accurately
And targeting should update dynamically with new purchases
```

## Scenario Outline: Discount Application and Calculation
```gherkin
Scenario Outline: Percentage discount calculation accuracy
Given a promotional campaign offering "<discount_percentage>%" discount
And a customer's cart contains items totaling $<cart_total>
And the promotion has a maximum discount limit of $<max_discount>
When the customer applies the promotion
Then the discount amount should be $<expected_discount>
And the final cart total should be $<expected_final_total>

Examples:
| discount_percentage | cart_total | max_discount | expected_discount | expected_final_total |
| 15                 | 100.00     | none         | 15.00            | 85.00               |
| 20                 | 250.00     | 40.00        | 40.00            | 210.00              |
| 25                 | 80.00      | none         | 20.00            | 60.00               |
| 10                 | 500.00     | 35.00        | 35.00            | 465.00              |

Scenario: Buy-X-Get-Y discount application
Given a "Buy 2 Get 1 Free" promotion for electronics
And a customer adds the following items to their cart:
  | Item          | Price  | Category    | Quantity |
  | Laptop        | $500   | Electronics | 1        |
  | Mouse         | $30    | Electronics | 1        |
  | Keyboard      | $50    | Electronics | 1        |
  | T-shirt       | $20    | Clothing    | 1        |
When the promotion is applied
Then the system should identify 1 qualifying set (Laptop + Mouse = 2 items)
And the lowest-priced eligible item (Mouse - $30) should be free
And the total discount should be $30
And the final cart total should be $570 ($500 + $50 + $20)
And the T-shirt should not be included in the promotion

Scenario: Bundle promotion application
Given a "Tech Starter Bundle" promotion with:
  | Component    | Individual Price | Bundle Position |
  | Laptop       | $999.99         | Required        |
  | Mouse        | $29.99          | Required        |
  | Headphones   | $79.99          | Required        |
And the bundle special price is $899.99
When a customer adds all bundle components to their cart
Then the bundle promotion should be automatically applied
And the customer should save $109.98 compared to individual prices
And the bundle discount should be clearly displayed
And removing any component should remove the bundle pricing
```

## Scenario Outline: Promotional Conflict Resolution
```gherkin
Scenario: Best promotion selection for customer benefit
Given a customer qualifies for multiple promotions:
  | Promotion ID | Type         | Value | Cart Total | Calculated Benefit |
  | PROMO-001   | PERCENTAGE   | 15%   | $200       | $30               |
  | PROMO-002   | FIXED_AMOUNT | $40   | $200       | $40               |
  | PROMO-003   | BUY_X_GET_Y  | B2G1F | $200       | $25               |
When the conflict resolution system evaluates all options
Then the system should automatically select "PROMO-002" (Fixed Amount $40)
And the customer should receive the maximum possible benefit
And the selection reason should be documented as "Highest customer benefit"
And alternative promotions should be shown for comparison

Scenario: Stackable promotions combination
Given the following stackable promotions are active:
  | Promotion        | Type      | Value | Stackable | Priority |
  | Category Bonus   | PERCENTAGE| 10%   | Yes       | 2        |
  | Loyalty Tier     | PERCENTAGE| 5%    | Yes       | 3        |
  | Free Shipping    | FIXED     | $15   | Yes       | 1        |
And a Gold tier customer's cart contains $200 of electronics
When all applicable promotions are evaluated
Then all three promotions should be applied
And the stacking order should be: Free Shipping, Category Bonus, Loyalty Tier
And the total benefit should be $46 ($15 + $20 + $11)
And the final cart total should be $154

Scenario: Priority-based conflict resolution
Given conflicting non-stackable promotions:
  | Promotion ID | Priority | Discount | Exclusive |
  | VIP-SALE     | 1        | $50      | Yes       |
  | FLASH-SALE   | 3        | $60      | Yes       |
  | CLEARANCE    | 5        | $70      | Yes       |
When a customer qualifies for all three promotions
Then the system should select "VIP-SALE" based on highest priority
And the customer should receive $50 discount (not the highest amount)
And the resolution should be documented as "Priority-based selection"
And the customer should be informed of the applied promotion
```

## Scenario Outline: Campaign Performance Analytics
```gherkin
Scenario: Real-time campaign performance tracking
Given an active promotional campaign "Spring Sale 2024"
And the campaign has been running for 24 hours
When the following activities occur:
  | Activity Type      | Count |
  | Campaign Views     | 1,000 |
  | Promotion Applications | 150   |
  | Successful Conversions | 120   |
  | Total Revenue Generated | $15,000 |
  | Promotion Costs    | $3,000 |
Then the campaign analytics should show:
  | Metric                | Value  | Calculation             |
  | Application Rate      | 15.0%  | 150/1000 * 100        |
  | Conversion Rate       | 80.0%  | 120/150 * 100         |
  | Return on Investment  | 400%   | (15000-3000)/3000 * 100 |
  | Cost per Conversion   | $25    | 3000/120               |
And all metrics should update in real-time
And performance dashboards should reflect current status

Scenario: Campaign ROI calculation and analysis
Given a promotional campaign with:
  | Cost Component        | Amount   |
  | Creative Development  | $5,000   |
  | Promotion Discounts   | $15,000  |
  | Platform Fees         | $1,000   |
  | Total Campaign Cost   | $21,000  |
And revenue attribution shows:
  | Revenue Type          | Amount   |
  | Gross Revenue         | $75,000  |
  | Attributed Revenue    | $60,000  |
  | Incremental Revenue   | $45,000  |
When ROI analysis is performed
Then the system should calculate:
  | ROI Metric           | Value   | Formula                    |
  | Gross ROI            | 257%    | (75000-21000)/21000 * 100  |
  | Attributed ROI       | 186%    | (60000-21000)/21000 * 100  |
  | Incremental ROI      | 114%    | (45000-21000)/21000 * 100  |
And detailed ROI breakdown should be available
And comparative analysis with previous campaigns should be provided
```

## Scenario Outline: Budget Management and Control
```gherkin
Scenario: Real-time budget tracking and alerts
Given a promotional campaign with a budget of $10,000
And budget alert thresholds at:
  | Threshold | Amount  | Alert Type |
  | 75%       | $7,500  | Warning    |
  | 90%       | $9,000  | Critical   |
  | 95%       | $9,500  | Urgent     |
When promotion applications consume the budget:
  | Application | Amount | Running Total | Expected Alert |
  | APP-001     | $3,000 | $3,000       | None          |
  | APP-002     | $2,000 | $5,000       | None          |
  | APP-003     | $3,000 | $8,000       | Warning       |
  | APP-004     | $1,500 | $9,500       | Urgent        |
Then budget tracking should be accurate to the cent
And appropriate alerts should be triggered at each threshold
And budget status should be updated in real-time
And campaign managers should receive notifications

Scenario: Budget exhaustion and campaign deactivation
Given a promotional campaign with $1,000 remaining budget
And the campaign is currently active
When a promotion application for $1,200 is attempted
Then the system should check budget availability
And the application should be rejected due to insufficient funds
And the campaign status should be updated to "Budget Exhausted"
And no further promotion applications should be accepted
And customers should see an "Offer no longer available" message
And campaign managers should be notified of budget exhaustion
```

## Scenario Outline: Mobile and Multi-Channel Experience
```gherkin
Scenario: Mobile promotion application and experience
Given a customer is using the mobile application
And a "Flash Sale" promotion is active with 4 hours remaining
When the customer browses eligible products
Then promotional pricing should be clearly displayed
And the countdown timer should show remaining time
And the mobile checkout process should apply discounts correctly
And the savings summary should be mobile-optimized
And push notifications should alert about expiring promotions

Scenario: Cross-channel promotion consistency
Given a multi-channel promotion active across:
  | Channel        | Platform    | Expected Behavior        |
  | Web Browser    | Desktop     | Full promotion details   |
  | Mobile App     | iOS/Android | Optimized for touch      |
  | In-Store       | POS System  | Automatic application    |
  | Phone Orders   | Call Center | Manual code entry        |
When customers access the promotion through different channels
Then promotion terms and benefits should be identical
And discount calculations should be consistent
And usage tracking should aggregate across channels
And customer experience should be seamless
```

## Scenario Outline: Error Handling and Edge Cases
```gherkin
Scenario: Invalid promotional code handling
Given a customer enters promotional code "INVALID2024"
And no active campaign uses this code
When the customer attempts to apply the code
Then the system should validate the code
And display error message "Invalid promotional code"
And suggest alternative active promotions
And log the invalid code attempt for analysis
And allow the customer to continue shopping without the code

Scenario: Expired promotion handling
Given a promotional campaign that expired 2 hours ago
And a customer still has the promotion page bookmarked
When the customer attempts to apply the expired promotion
Then the system should check promotion validity
And display message "This promotion has expired"
And show current active promotions as alternatives
And update the promotion page to reflect expiration
And remove expired promotions from search results

Scenario: Inventory insufficient for promotion fulfillment
Given a "Buy 2 Get 1 Free" promotion for limited stock items
And only 2 units remain in inventory
When a customer attempts to purchase 3 units (qualifying for 1 free)
Then the system should check inventory availability
And adjust the promotion to available stock
And inform customer of inventory limitation
And offer alternative promotions or products
And prevent overselling through promotion abuse

Scenario: System failure during promotion application
Given a customer is in the checkout process
And a valid promotion is being applied
When a system error occurs during discount calculation
Then the system should implement graceful error handling
And preserve the customer's cart contents
And display "Temporary issue applying discount"
And allow retry or manual customer service intervention
And log the error for technical investigation
And ensure no double-discount application on retry
```

## Feature: Campaign Approval Workflow
```gherkin
Scenario: Marketing campaign approval process
Given a marketing manager has created a new promotion campaign
And the campaign exceeds the auto-approval threshold of $5,000 budget
When the campaign is submitted for approval
Then it should enter the approval workflow
And campaign approvers should be notified
And the campaign status should be "Pending Approval"
And no customers should see the promotion until approved

Scenario: Multi-level approval for high-value campaigns
Given a campaign with budget exceeding $25,000
And the campaign targets all customer segments
When submitted for approval
Then it should require director-level approval
And legal review for terms and conditions
And finance approval for budget allocation
And all approvals must be completed before activation
And approval history should be maintained for audit

Scenario: Campaign rejection and feedback process
Given a campaign awaiting approval
And the approver identifies compliance issues
When the approver rejects the campaign
Then detailed rejection reasons should be provided
And the campaign should return to draft status
And the marketing manager should be notified
And the campaign should be available for revision
And resubmission should trigger a new approval cycle
```

## Feature: Fraud Prevention and Security
```gherkin
Scenario: Promotion abuse detection and prevention
Given a customer has applied the same promotion code multiple times
And the promotion has a "one per customer" limit
When the customer attempts another application
Then the system should detect the duplicate usage
And prevent the additional application
And display "Promotion limit reached" message
And log the attempt for fraud analysis
And not penalize legitimate future promotions

Scenario: Suspicious promotional activity monitoring
Given unusual promotion usage patterns are detected:
  | Pattern                    | Threshold | Detected Activity |
  | Rapid code attempts        | 10/minute | 15 attempts       |
  | Multiple account usage     | 3 accounts| 5 accounts        |
  | High-value applications    | $1000     | $2500            |
When the fraud detection system analyzes the activity
Then suspicious accounts should be flagged for review
And temporary holds should be placed on high-risk applications
And security team should be alerted
And legitimate customers should not be affected
And detailed audit trail should be maintained
```
