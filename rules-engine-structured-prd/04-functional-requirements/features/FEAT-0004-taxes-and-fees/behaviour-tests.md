# Behaviour Tests - Taxes and Fees Rules

## TC-BEH-01: Standard Tax Calculation
**Scenario**: Customer makes a purchase requiring tax calculation
**Given** I am processing a transaction for a customer in California
**And** the transaction total is $100.00
**And** the California state tax rate is 7.25%
**When** I request tax calculation for the transaction
**Then** the system should calculate state tax as $7.25
**And** the total transaction amount should be $107.25
**And** the calculation should complete within 200ms

## TC-BEH-02: Multi-Jurisdiction Tax Calculation
**Scenario**: Complex transaction spans multiple tax jurisdictions
**Given** I have a transaction with items shipping to different states
**And** items worth $50 ship to California (7.25% tax)
**And** items worth $75 ship to New York (8.00% tax)
**When** I calculate taxes for the entire transaction
**Then** California items should be taxed at 7.25% = $3.63
**And** New York items should be taxed at 8.00% = $6.00
**And** total tax should be $9.63
**And** each jurisdiction should be clearly itemized

## TC-BEH-03: Fee Calculation with Exemptions
**Scenario**: VIP customer receives fee exemptions
**Given** I have a VIP customer with fee exemption status
**And** the customer makes a $300 purchase
**And** standard processing fees would be $7.50
**When** I calculate fees with exemption applied
**Then** the customer should pay $0 in processing fees
**And** the exemption should be applied automatically
**And** the exemption usage should be tracked