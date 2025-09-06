# Behaviour Tests - Loyalty Management

## Feature: Customer Loyalty Program Management

### Background:
```gherkin
Given the loyalty program is active
And the following tier thresholds are configured:
  | Tier     | Minimum Spending | Multiplier |
  | Bronze   | $0              | 1.0x       |
  | Silver   | $1,000          | 1.5x       |
  | Gold     | $5,000          | 2.0x       |
  | Platinum | $10,000         | 3.0x       |
```

## Scenario Outline: Customer Tier Assignment
```gherkin
Scenario: New customer gets default Bronze tier
Given a new customer registers with the system
When the customer account is created
Then the customer should be assigned to "Bronze" tier
And the customer should see Bronze tier benefits
And the tier effective date should be the account creation date

Scenario: Customer qualifies for tier upgrade
Given a customer with "Silver" tier
And the customer has spent $4,800 in the last 12 months
When the customer makes a purchase of $500
And the tier calculation process runs
Then the customer should be upgraded to "Gold" tier
And the customer should receive a tier upgrade notification
And Gold tier benefits should be immediately available

Scenario: Customer tier downgrade after grace period
Given a customer with "Gold" tier
And the customer's 12-month spending drops to $2,000
And the spending has been below the Gold threshold for 91 days
When the monthly tier review process runs
Then the customer should be downgraded to "Silver" tier
And the customer should receive a tier downgrade notification
And Gold tier benefits should be removed
```

## Scenario Outline: Points Earning and Calculation
```gherkin
Scenario: Standard points earning based on tier
Given a customer with "<tier>" tier
And the tier has a "<multiplier>" points multiplier
When the customer makes a purchase of $<amount>
Then the customer should earn <expected_points> points
And the points should be credited within 24 hours

Examples:
| tier     | multiplier | amount | expected_points |
| Bronze   | 1.0x       | 100    | 100            |
| Silver   | 1.5x       | 100    | 150            |
| Gold     | 2.0x       | 100    | 200            |
| Platinum | 3.0x       | 100    | 300            |

Scenario: Category bonus points earning
Given a customer with "Gold" tier (2.0x multiplier)
And "Electronics" category has a 2.0x bonus multiplier
When the customer purchases $200 worth of electronics
Then the customer should earn 800 points
And the breakdown should show:
  | Component | Points |
  | Base      | 200    |
  | Tier      | 200    |
  | Category  | 400    |
  | Total     | 800    |

Scenario: Promotional points multiplier
Given a customer with "Bronze" tier
And there is an active "Double Points Weekend" promotion with 2.0x multiplier
When the customer makes a purchase of $75 during the promotion
Then the customer should earn 150 points
And the promotional bonus should be tracked separately
And the points breakdown should indicate promotional earning
```

## Scenario Outline: Points Redemption Process
```gherkin
Scenario: Successful reward redemption
Given a customer has 5,000 points in their account
And a "$25 Gift Card" reward costs 2,500 points
And the reward is available in inventory
When the customer redeems the reward
Then 2,500 points should be deducted from their balance
And the customer should have 2,500 points remaining
And the gift card should be issued to the customer
And the redemption should appear in transaction history

Scenario: Insufficient points for redemption
Given a customer has 1,500 points in their account
And a "$25 Gift Card" reward costs 2,500 points
When the customer attempts to redeem the reward
Then the redemption should be rejected
And the customer should see an error message "Insufficient points"
And alternative rewards within their points range should be suggested
And no points should be deducted

Scenario: Reward unavailable due to inventory
Given a customer has 5,000 points in their account
And a "$25 Gift Card" reward costs 2,500 points
But the reward is out of stock
When the customer attempts to redeem the reward
Then the redemption should be rejected
And the customer should see a message "Reward currently unavailable"
And alternative similar rewards should be suggested
And no points should be deducted
```

## Scenario Outline: Points Expiration Management
```gherkin
Scenario: Points expiration warning notification
Given a customer has 1,200 points earned 22 months ago
And points expire after 24 months of inactivity
And the customer has not made any purchases recently
When the daily expiration check runs
And the points are 30 days from expiring
Then the customer should receive a 30-day expiration warning
And the notification should include:
  | Information              | Value                    |
  | Points expiring          | 1,200                   |
  | Expiration date          | [30 days from today]    |
  | Suggested actions        | Make purchase or redeem |

Scenario: Points expiration processing
Given a customer has 800 points that expire today
And the customer has not been active in the last 24 months
When the daily expiration processing runs
Then the 800 points should be removed from the customer's account
And an expiration transaction should be recorded
And the customer should receive an expiration notification

Scenario: Activity-based expiration extension
Given a customer has 1,500 points expiring in 15 days
When the customer makes any purchase
Then all the customer's points should have their expiration extended to 24 months from the purchase date
And the customer should receive a notification about the extension
```

## Scenario Outline: Partner Integration
```gherkin
Scenario: Earning points at partner merchant
Given a customer shops at "Partner Coffee Shop"
And the partner has a 2.0 points per dollar earning rate
When the customer makes a $15 purchase at the partner location
Then the customer should earn 30 points
And the points should be synchronized to the main loyalty account within 4 hours
And the partner transaction should appear in the customer's activity history

Scenario: Redeeming points for partner rewards
Given a customer has 3,000 points in their account
And "Partner Restaurant" offers a "$15 meal voucher" for 1,500 points
When the customer redeems points for the partner reward
Then 1,500 points should be deducted from the main account
And the partner should be notified of the redemption
And the partner should receive fulfillment instructions
And the customer should receive confirmation with redemption details
```

## Scenario Outline: Tier Benefits and Privileges
```gherkin
Scenario: Gold tier benefits activation
Given a customer is upgraded to "Gold" tier
When the tier upgrade is processed
Then the customer should receive the following benefits:
  | Benefit              | Description                    |
  | Free shipping        | On orders over $50            |
  | Priority support     | Dedicated customer service   |
  | Exclusive offers     | Early access to sales         |
  | Birthday bonus       | 2x points on birthday month   |

Scenario: Tier benefits inheritance
Given a customer with "Platinum" tier
When the customer makes a purchase
Then they should receive all lower tier benefits plus:
  | Additional Benefit   | Description                   |
  | Free returns         | Unlimited free return shipping|
  | Personal shopper     | Dedicated personal assistant  |
  | VIP events          | Exclusive event invitations   |
```

## Scenario Outline: Customer Service and Support
```gherkin
Scenario: Customer service representative accesses loyalty information
Given a customer service representative is assisting a customer
And the customer ID is "CUST-12345"
When the representative looks up the customer's loyalty information
Then they should see:
  | Information          | Value    |
  | Current tier         | Gold     |
  | Points balance       | 3,250    |
  | Next tier progress   | 65%      |
  | Last activity        | 5 days   |
  | Expiring points      | 500      |

Scenario: Manual points adjustment by customer service
Given a customer service representative needs to adjust points due to a technical issue
And the representative has proper authorization
And the customer should receive 500 points compensation
When the representative processes the manual adjustment
Then 500 points should be added to the customer's account
And the adjustment should be recorded with:
  | Field       | Value                          |
  | Reason      | Technical issue compensation   |
  | Authorized  | Rep ID and supervisor approval |
  | Timestamp   | Current date and time         |
And the customer should receive a notification about the adjustment
```

## Scenario Outline: Analytics and Reporting
```gherkin
Scenario: Loyalty program performance analytics
Given the loyalty program has been running for 12 months
When the monthly analytics report is generated
Then the report should include:
  | Metric                    | Description                           |
  | Active members           | Customers with activity in last 90 days |
  | Tier distribution        | Percentage of customers in each tier     |
  | Points liability         | Total outstanding points value           |
  | Redemption rate          | Percentage of earned points redeemed     |
  | Customer lifetime value  | Average CLV by tier                      |

Scenario: Customer engagement insights
Given customer loyalty data for analysis
When engagement metrics are calculated
Then insights should show:
  | Insight                  | Measurement                        |
  | Purchase frequency       | Average days between purchases     |
  | Tier progression rate    | Percentage advancing tiers yearly  |
  | Points hoarding behavior | Average points balance vs earning  |
  | Preferred reward types   | Most popular redemption categories |
```

## Scenario Outline: Error Handling and Edge Cases
```gherkin
Scenario: System failure during points redemption
Given a customer initiates a points redemption for a $50 gift card
And the system experiences a failure after deducting points but before issuing the reward
When the error recovery process runs
Then the points should be restored to the customer's account
And the customer should be notified of the technical issue
And the redemption should be marked as failed for retry

Scenario: Duplicate transaction prevention
Given a customer makes a purchase
And due to network issues, the same transaction is submitted twice
When the duplicate detection system processes the transactions
Then points should only be awarded once
And the duplicate transaction should be logged
And the customer should receive points for only one transaction

Scenario: Points calculation during tier transition
Given a customer is about to be upgraded from Silver to Gold tier
And the customer makes a purchase during the tier calculation process
When the points are calculated for that purchase
Then the points should be calculated using the new Gold tier rate
And the tier upgrade should be applied before points calculation
And the customer should receive notifications for both tier upgrade and points earned
```

## Feature: Mobile Loyalty Experience
```gherkin
Scenario: Mobile app loyalty dashboard
Given a customer opens the mobile loyalty app
When they navigate to the loyalty dashboard
Then they should see:
  | Information           | Description                    |
  | Current points        | Real-time points balance      |
  | Tier status          | Current tier and progress bar  |
  | Nearby offers        | Location-based promotions      |
  | Recent activity      | Last 5 transactions           |
  | Expiring points      | Points expiring within 60 days |

Scenario: Mobile reward redemption
Given a customer browses rewards in the mobile app
When they select a reward and confirm redemption
Then the redemption should process immediately
And a QR code or digital voucher should be generated
And the reward should be available for immediate use
And the transaction should sync with the main loyalty system
```

## Feature: Fraud Prevention and Security
```gherkin
Scenario: Suspicious points activity detection
Given a customer account shows unusual activity patterns:
  | Pattern               | Description                    |
  | Multiple large purchases | 5 purchases over $1000 in 1 day |
  | Rapid tier progression  | Bronze to Platinum in 1 week    |
  | Unusual redemption      | 50,000 points redeemed at once  |
When the fraud detection system analyzes the activity
Then the account should be flagged for review
And suspicious transactions should be temporarily held
And the customer should be contacted for verification

Scenario: Account security verification
Given a customer attempts to redeem a high-value reward (>10,000 points)
When the redemption is submitted
Then additional security verification should be required
And the customer should receive a verification code via their preferred method
And the redemption should only proceed after successful verification
```
