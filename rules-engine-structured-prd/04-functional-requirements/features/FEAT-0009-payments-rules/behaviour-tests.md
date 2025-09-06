# Behaviour Tests - Payment Rules and Processing

## Feature: Intelligent Payment Processing System

### Background:
```gherkin
Given the payment processing system is active
And the following payment methods are supported:
  | Method Type    | Processing Time | Security Level | Cost Level |
  | CREDIT_CARD    | Immediate      | High          | Medium     |
  | DEBIT_CARD     | Immediate      | High          | Low        |
  | DIGITAL_WALLET | Immediate      | Very High     | Medium     |
  | BANK_TRANSFER  | 1-3 days       | Medium        | Very Low   |
  | BUY_NOW_PAY_LATER | Immediate   | Medium        | High       |
And the following risk levels are configured:
  | Risk Level | Score Range | Action Required    |
  | VERY_LOW   | 0-20       | Auto-approve       |
  | LOW        | 21-40      | Standard processing |
  | MEDIUM     | 41-70      | Additional verification |
  | HIGH       | 71-90      | Manual review      |
  | VERY_HIGH  | 91-100     | Block transaction  |
```

## Scenario Outline: Smart Payment Method Selection
```gherkin
Scenario: System recommends optimal payment method based on customer history
Given a customer with the following payment history:
  | Payment Method | Success Rate | Last Used  | Average Amount |
  | Visa Card *1234| 98%         | 2024-12-15 | $150          |
  | PayPal Account | 95%         | 2024-12-10 | $200          |
  | Amex Card *5678| 85%         | 2024-12-05 | $300          |
When the customer initiates a payment for $250
Then the system should recommend "Visa Card *1234" as the primary method
And "PayPal Account" should be offered as the first alternative
And "Amex Card *5678" should be offered as the second alternative
And the success probability should be displayed for each method

Scenario: Payment method validation against transaction requirements
Given a transaction requiring immediate processing for $5,000
And the customer has the following payment methods:
  | Method         | Max Amount | Processing Time | Status |
  | Credit Card    | $3,000     | Immediate      | Valid  |
  | Bank Transfer  | $50,000    | 1-3 days       | Valid  |
  | Digital Wallet | $10,000    | Immediate      | Valid  |
When the payment method validation is performed
Then "Digital Wallet" should be recommended as it meets all requirements
And "Credit Card" should be marked as invalid due to amount limit
And "Bank Transfer" should be marked as unsuitable due to processing time
And clear validation messages should be provided for each method
```

## Scenario Outline: Fraud Detection and Risk Assessment
```gherkin
Scenario Outline: Risk assessment based on transaction characteristics
Given a customer with "<customer_profile>" transaction history
And a transaction for $<amount> from "<location>" using "<device>"
When the fraud risk assessment is performed
Then the risk score should be "<risk_score>"
And the risk level should be "<risk_level>"
And the system should "<action>"

Examples:
| customer_profile | amount | location        | device         | risk_score | risk_level | action           |
| trusted_customer | 100    | usual_location  | trusted_device | 15         | VERY_LOW   | auto-approve     |
| new_customer     | 500    | new_location    | new_device     | 45         | MEDIUM     | verify_identity  |
| suspicious_user  | 2000   | foreign_country | unknown_device | 85         | HIGH       | manual_review    |
| flagged_account  | 5000   | blacklist_ip    | tor_browser    | 95         | VERY_HIGH  | block_transaction|

Scenario: Multi-factor fraud detection for suspicious patterns
Given a customer account with normal spending pattern of $50-100 per transaction
And typical transaction frequency of 2-3 times per week
And usual geographic location of "New York, USA"
When the customer attempts 5 transactions of $1,000 each within 1 hour
And the transactions originate from "Lagos, Nigeria"
And the device fingerprint is unrecognized
Then the fraud detection system should identify multiple risk factors:
  | Risk Factor         | Weight | Detected |
  | Amount anomaly      | High   | Yes      |
  | Velocity anomaly    | High   | Yes      |
  | Geographic anomaly  | Medium | Yes      |
  | Device anomaly      | Medium | Yes      |
And the composite risk score should exceed 90
And all transactions should be automatically blocked
And the security team should be immediately alerted
```

## Scenario Outline: Intelligent Gateway Routing
```gherkin
Scenario: Optimal gateway selection based on multiple factors
Given the following payment gateways are available:
  | Gateway   | Success Rate | Cost | Response Time | Card Support        |
  | Gateway A | 97%         | 2.9% | 1.2s         | Visa, MasterCard    |
  | Gateway B | 95%         | 2.1% | 0.8s         | Visa, MasterCard, Amex |
  | Gateway C | 99%         | 3.2% | 1.5s         | All Cards           |
When a customer pays $150 with a Visa card for a standard priority transaction
Then the system should evaluate gateway options
And select "Gateway C" for its highest success rate
And document the selection reason as "Optimal success rate for standard transaction"
And configure "Gateway A" as the backup option

Scenario: Automatic gateway failover during processing
Given a payment is routed to "Primary Gateway" 
And the payment amount is $200 with estimated processing time of 2 seconds
When the primary gateway responds with "TIMEOUT" after 3 seconds
Then the system should detect the gateway failure
And automatically failover to "Backup Gateway" within 500ms
And retry the payment processing without customer intervention
And notify the operations team of the gateway failure
And complete the transaction successfully via backup gateway
And the customer should experience minimal delay
```

## Scenario Outline: Multi-Currency Payment Processing
```gherkin
Scenario Outline: Currency conversion with transparent fees
Given a customer located in "<customer_country>"
And they want to purchase an item priced at <original_amount> <original_currency>
When they proceed to payment
Then the system should detect their location as "<customer_country>"
And offer to display the price in "<customer_currency>"
And show the converted amount as "<converted_amount>"
And clearly display the conversion fee as "<conversion_fee>"
And show the total charge as "<total_amount>"

Examples:
| customer_country | original_amount | original_currency | customer_currency | converted_amount | conversion_fee | total_amount |
| Canada          | 100             | USD              | CAD              | 135.00          | 3.38          | 138.38      |
| United Kingdom  | 100             | USD              | GBP              | 79.50           | 1.99          | 81.49       |
| European Union  | 100             | USD              | EUR              | 92.30           | 2.31          | 94.61       |

Scenario: Real-time exchange rate updates and accuracy
Given a customer is viewing a product priced at €500 EUR
And the current EUR/USD exchange rate is 1.0850
When the customer switches to USD pricing
Then the converted amount should be calculated as $542.50
And the conversion should use real-time exchange rates updated within the last minute
And if the exchange rate changes during checkout
Then the customer should be notified of the rate change
And given the option to proceed with new rate or cancel
```

## Scenario Outline: Payment Retry and Recovery
```gherkin
Scenario Outline: Intelligent retry strategies based on failure type
Given a payment transaction fails with reason "<failure_reason>"
When the retry logic is triggered
Then the system should apply "<retry_strategy>" strategy
And attempt "<max_retries>" retries
And use "<backoff_pattern>" timing between attempts

Examples:
| failure_reason           | retry_strategy      | max_retries | backoff_pattern |
| TEMPORARY_GATEWAY_ERROR  | IMMEDIATE_RETRY     | 3          | EXPONENTIAL     |
| NETWORK_TIMEOUT         | GATEWAY_SWITCH      | 2          | IMMEDIATE       |
| INSUFFICIENT_FUNDS      | ALTERNATIVE_METHOD  | 0          | NONE           |
| INVALID_CARD           | CUSTOMER_ACTION     | 0          | NONE           |
| BANK_DECLINE           | ALTERNATIVE_GATEWAY | 2          | LINEAR         |

Scenario: Payment method fallback and customer guidance
Given a customer's primary payment method (Credit Card) fails with "CARD_EXPIRED"
When the payment failure is detected
Then the system should identify the failure as requiring customer action
And suggest alternative payment methods from the customer's wallet:
  | Alternative Method | Availability | Recommendation |
  | PayPal Account    | Available   | Primary        |
  | Bank Transfer     | Available   | Secondary      |
  | Add New Card      | Always      | Tertiary       |
And preserve the shopping cart and session information
And provide clear instructions for updating the expired card
And offer one-click retry with alternative methods
```

## Scenario Outline: Compliance and Security Validation
```gherkin
Scenario: PCI DSS compliance for credit card processing
Given a customer enters credit card information:
  | Card Number      | Expiry | CVV | Name         |
  | 4111111111111111 | 12/25  | 123 | Test Customer|
When the card data is submitted for processing
Then all card data should be encrypted in transit using TLS 1.3
And the card number should be tokenized within 100ms
And the tokenized version should replace the card number: "tok_4111xxxxxxxx1111"
And no plain text card data should appear in logs
And the CVV should be discarded immediately after validation
And a PCI DSS compliance audit entry should be created

Scenario: AML screening for high-value transactions
Given a customer attempts a transaction for $15,000 USD
And the AML screening threshold is set to $10,000
When the transaction is processed
Then the AML screening should be automatically triggered
And the customer should be checked against sanctions lists
And the transaction patterns should be analyzed for suspicious activity:
  | Check Type           | Result    | Action Required |
  | Sanctions List       | Clear     | None           |
  | Suspicious Patterns  | Clear     | None           |
  | Customer Verification| Valid     | None           |
And the transaction should be approved for processing
And an AML compliance report should be generated
And the compliance officer should be notified for high-value transaction tracking

Scenario: KYC verification for new high-value customers
Given a new customer who registered 2 days ago
And they attempt their first transaction for $8,000
And the KYC requirement threshold is $5,000 for new customers
When the transaction is initiated
Then the KYC verification should be triggered
And the customer should be prompted to provide:
  | Document Type        | Required | Status   |
  | Government ID        | Yes      | Pending  |
  | Proof of Address     | Yes      | Pending  |
  | Income Verification  | No       | Optional |
And the transaction should be held pending verification
And clear instructions should be provided for document upload
And the estimated verification time should be communicated
```

## Scenario Outline: Mobile Payment Optimization
```gherkin
Scenario: Mobile wallet integration and biometric authentication
Given a customer is using an iPhone with Apple Pay enabled
And they have Touch ID configured
When they select payment method at checkout
Then Apple Pay should be offered as the primary payment option
And the payment sheet should display:
  | Field           | Value                    |
  | Merchant Name   | Rules Engine Store       |
  | Amount          | $150.00                 |
  | Payment Method  | Apple Pay (Touch ID)    |
When they authorize with Touch ID
Then the payment should be processed using the device's secure element
And the transaction should complete within 3 seconds
And the customer should receive immediate confirmation

Scenario: Mobile-specific fraud detection considerations
Given a customer making a mobile payment
When the fraud assessment is performed
Then the system should consider mobile-specific factors:
  | Factor                | Weight | Description                    |
  | Device Fingerprint    | High   | Hardware and software profile  |
  | Location Consistency  | Medium | GPS vs. billing address        |
  | App Behavior         | Medium | Usage patterns and navigation  |
  | Biometric Usage      | Low    | Biometric authentication used  |
And adjust the risk score based on mobile context
And apply mobile-optimized security measures
```

## Feature: Recurring Payment Management
```gherkin
Scenario: Automatic recurring payment processing
Given a customer has a monthly subscription for $29.99
And the next billing date is today
And their saved payment method is valid
When the recurring payment job runs
Then the payment should be automatically processed
And the customer should be charged $29.99
And a receipt should be sent via email
And the next billing date should be updated to next month

Scenario: Failed recurring payment retry with dunning management
Given a recurring payment fails due to "INSUFFICIENT_FUNDS"
When the retry logic for recurring payments is triggered
Then the system should apply dunning management strategy:
  | Retry Attempt | Delay | Action              |
  | 1            | 3 days| Retry same method   |
  | 2            | 7 days| Retry + notify customer |
  | 3            | 14 days| Final retry + warning |
And send appropriate communications at each stage
And provide easy payment update options
And only cancel subscription after final failure
```

## Feature: Performance and Reliability
```gherkin
Scenario: High-volume payment processing during peak times
Given the system is processing 1000 concurrent payments
When payment processing load increases to peak capacity
Then all payment operations should complete within SLA:
  | Operation              | SLA Target | Actual Performance |
  | Payment Method Selection| <100ms    | 85ms              |
  | Fraud Risk Assessment  | <200ms    | 175ms             |
  | Gateway Selection      | <150ms    | 120ms             |
  | Total Processing Time  | <500ms    | 425ms             |
And the system should maintain 99.5% success rate
And no transactions should be lost or corrupted
And error rates should remain below 0.1%

Scenario: System resilience during partial service failures
Given the primary fraud detection service becomes unavailable
When payment processing continues
Then the system should switch to backup fraud detection
And apply conservative fraud rules as fallback
And continue processing low-risk transactions
And queue medium-risk transactions for review
And automatically resume normal operations when service recovers
And maintain audit trail of all decisions made during outage
```
