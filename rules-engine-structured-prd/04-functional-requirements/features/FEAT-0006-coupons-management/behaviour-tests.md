# Behaviour Tests - Coupons Management

## Feature: Coupon Creation and Management

### Scenario: Create a percentage discount coupon
```gherkin
Given I am a Marketing Manager with coupon creation privileges
When I create a coupon with the following details:
  | Field               | Value                    |
  | Discount Type       | Percentage               |
  | Discount Value      | 15%                      |
  | Maximum Discount    | $100.00                  |
  | Minimum Purchase    | $50.00                   |
  | Usage Type          | Multi-use                |
  | Maximum Uses        | 1000                     |
  | Valid From          | 2024-01-01T00:00:00Z     |
  | Valid To            | 2024-01-31T23:59:59Z     |
  | Target Category     | Electronics              |
Then the coupon should be created successfully
And the coupon code should follow the pattern "XXXX-XXXX-XXXX"
And the coupon status should be "DRAFT"
And an audit log entry should be created
```

### Scenario: Prevent duplicate coupon codes
```gherkin
Given a coupon with code "SUMMER2024" already exists
When I attempt to create another coupon with code "SUMMER2024"
Then the system should reject the creation
And display error message "Coupon code already exists"
And no new coupon should be saved
```

### Scenario: Generate bulk coupons for campaign
```gherkin
Given I have a campaign template "Summer Sale 2024"
When I generate 5000 coupons using the template
Then all 5000 coupons should be created with unique codes
And each coupon should follow the pattern "SUMMER-XXXX"
And the generation should complete within 5 minutes
And all coupons should have status "DRAFT"
```

---

## Feature: Coupon Validation and Redemption

### Scenario: Apply valid coupon at checkout
```gherkin
Given I have a cart with items totaling $75.00
And there is an active coupon "SAVE10NOW" offering $10 off purchases over $50
When I apply the coupon code "SAVE10NOW"
Then the coupon should be validated within 200ms
And a discount of $10.00 should be applied
And my order total should become $65.00
And the coupon usage count should be incremented
```

### Scenario: Reject expired coupon
```gherkin
Given there is a coupon "EXPIRED123" that expired yesterday
And I have eligible items worth $100 in my cart
When I attempt to apply coupon "EXPIRED123"
Then the system should validate the coupon within 200ms
And display error message "This coupon has expired"
And no discount should be applied to my order
And my order total should remain $100.00
```

### Scenario: Prevent multiple uses of single-use coupon
```gherkin
Given I am customer "john.doe@email.com"
And I have already used single-use coupon "WELCOME50"
And I have eligible items worth $200 in my cart
When I attempt to apply coupon "WELCOME50" again
Then the system should check my redemption history
And display error message "This coupon has already been used"
And no discount should be applied
```

### Scenario: Apply coupon with minimum purchase requirement
```gherkin
Given there is a coupon "SAVE20" offering $20 off purchases over $100
When I apply coupon "SAVE20" to a cart worth $75
Then the system should validate the minimum purchase requirement
And display error message "Minimum purchase of $100.00 required"
And no discount should be applied
But when I add items to reach $100 total
And apply the same coupon
Then the discount should be applied successfully
```

---

## Feature: Fraud Detection and Security

### Scenario: Detect rapid redemption pattern
```gherkin
Given there is a coupon "PROMO25" configured with fraud monitoring
And the fraud threshold is set to 50 uses per 5 minutes
When 60 different customers attempt to use "PROMO25" within 3 minutes
Then the fraud detection system should trigger after the 50th attempt
And the coupon should be automatically suspended
And a security alert should be sent to administrators
And subsequent attempts should receive message "Coupon temporarily unavailable"
```

### Scenario: Block systematic code guessing attempts
```gherkin
Given I am attempting to guess coupon codes systematically
When I try 25 different codes from the same IP address within 1 minute
Then the system should detect the brute force pattern
And block my IP address for 1 hour
And log the security incident
And return generic error message for all attempts
```

### Scenario: Detect suspicious usage velocity
```gherkin
Given coupon "FLASH50" has normal usage pattern of 10 redemptions per hour
When suddenly 100 redemptions occur within 10 minutes
Then the fraud detection should flag unusual velocity
And trigger investigation workflow
And alert security team
But continue processing legitimate redemptions
```

---

## Feature: Multi-Channel Integration

### Scenario: Validate coupon in POS system
```gherkin
Given I am at a physical store checkout
And there is an active coupon "STORE15" for 15% off in-store purchases
And my purchase total is $120.00
When the cashier enters coupon code "STORE15"
Then the POS system should validate the coupon within 300ms
And apply 15% discount ($18.00)
And update my total to $102.00
And print coupon details on the receipt
And update central redemption tracking
```

### Scenario: Auto-apply activated mobile coupon
```gherkin
Given I have activated coupon "MOBILE10" in the mobile app
And the coupon offers $10 off mobile purchases over $50
And I have $60 worth of items in my mobile cart
When I proceed to checkout in the mobile app
Then the activated coupon should be automatically detected
And the $10 discount should be applied without manual entry
And I should see the discount in my order summary
And I should have the option to remove the auto-applied coupon
```

### Scenario: Cross-channel coupon usage tracking
```gherkin
Given coupon "OMNI20" is valid across all channels
And I use it once in the mobile app for a $80 purchase
When I later attempt to use the same coupon in-store
And the coupon is configured as single-use per customer
Then the POS system should detect my previous usage
And prevent duplicate redemption
And display message "Coupon already used by this customer"
```

---

## Feature: Campaign Management and Analytics

### Scenario: Track campaign performance in real-time
```gherkin
Given I have a campaign "SPRING2024" with 1000 coupons
And the campaign has been active for 24 hours
When I view the campaign dashboard
Then I should see current redemption count
And redemption rate percentage
And revenue impact calculation
And channel breakdown (web, mobile, in-store)
And customer acquisition metrics
And all data should be updated in real-time
```

### Scenario: Alert on underperforming campaign
```gherkin
Given campaign "FLASH48" has target redemption rate of 10%
And it has been active for 12 hours with only 2% redemption rate
When the performance monitoring system evaluates the campaign
Then an underperformance alert should be triggered
And sent to the campaign manager
And include optimization suggestions
And provide comparison with similar campaigns
```

### Scenario: Automatically deactivate exhausted coupons
```gherkin
Given coupon "LIMITED500" has maximum 500 uses
And 499 redemptions have already occurred
When the 500th customer successfully redeems the coupon
Then the coupon should be automatically deactivated
And removed from all distribution channels
And any pending validations should be rejected
And the marketing team should be notified
```

---

## Feature: Advanced Coupon Features

### Scenario: Stack compatible coupons
```gherkin
Given I have a 10% loyalty discount as a Gold member
And there is a stackable coupon "EXTRA5" offering additional 5% off
And both promotions are configured as stackable
And I have $100 worth of eligible items
When I apply both promotions
Then the loyalty discount should be applied first: $100 → $90
And the coupon discount should be applied to result: $90 → $85.50
And I should see breakdown of both discounts
And both promotions should be recorded in transaction history
```

### Scenario: Resolve conflicting non-stackable promotions
```gherkin
Given there is a site-wide promotion offering 15% off electronics
And I apply coupon "TECH20" offering 20% off electronics
And promotions are configured as non-stackable
When the system processes my order
Then it should identify the promotion conflict
And apply the better discount for me (20%)
And display message "Applied best available discount: 20% off"
And record the conflict resolution in audit log
```

### Scenario: Apply dynamic coupon value based on customer tier
```gherkin
Given I am a Platinum tier customer
And there is a dynamic coupon "TIER_DISCOUNT" that adjusts based on customer tier
And the tier multipliers are: Bronze 5%, Silver 10%, Gold 15%, Platinum 20%
When I apply the dynamic coupon to a $100 purchase
Then the system should calculate my tier-based discount: 20%
And apply $20 discount to my order
And the final total should be $80
And the transaction should record the dynamic calculation
```

---

## Feature: Security and Compliance

### Scenario: Handle GDPR data deletion request
```gherkin
Given customer "jane.doe@email.com" has coupon redemption history
And she requests account deletion under GDPR
When the data deletion process is executed
Then her personal data should be anonymized in coupon records
And redemption events should maintain anonymous transaction data
And aggregated analytics should preserve business metrics
And customer-specific coupon allocations should be removed
And compliance audit log should record the deletion
```

### Scenario: Maintain comprehensive audit trail
```gherkin
Given coupon "AUDIT_TEST" goes through complete lifecycle
When an audit review is performed
Then the audit trail should show:
  | Event Type          | Details                           |
  | Coupon Created      | Creator, timestamp, configuration |
  | Configuration Edit  | Changes made, approver, timestamp |
  | Campaign Launch     | Distribution details, channels    |
  | Customer Redemption | Customer ID, transaction, amount  |
  | Fraud Detection     | Security events, responses        |
  | Coupon Expiration   | Automatic deactivation, unused   |
And all events should have immutable timestamps
And user attribution for all manual actions
```

---

## Feature: Error Handling and Recovery

### Scenario: Handle payment system failure during redemption
```gherkin
Given I have applied coupon "SAVE15" to my $100 order
And the discount has been calculated as $15
When the payment system fails during checkout
Then the coupon redemption should not be recorded
And the coupon should remain available for use
And I should be able to retry checkout with same coupon
And no phantom redemption should be logged
```

### Scenario: Recover from temporary service outage
```gherkin
Given the coupon validation service experiences temporary outage
And I attempt to apply coupon "EMERGENCY" during the outage
When the service comes back online within 30 seconds
Then my coupon application should be processed automatically
And the discount should be applied correctly
And the user experience should be seamless
And no manual intervention should be required
```

### Scenario: Handle database connection failure
```gherkin
Given the coupon database connection is lost
When customers attempt to apply coupons
Then the system should fail gracefully
And display helpful error message "Service temporarily unavailable"
And not process any redemptions that could cause data inconsistency
And automatically retry when connection is restored
And log all failed attempts for later processing
```
