# Functional Tests - Taxes and Fees Rules

## TC-FUNC-01: Basic Tax Calculation
**Prerequisites**: Valid jurisdiction data exists, tax rates are configured
**Steps**: 
1. Submit transaction with customer location and items
2. Request tax calculation for the transaction
3. Verify calculated tax amounts for each jurisdiction
4. Validate tax line item details
5. Confirm total tax calculation accuracy

**Test Data**: 
- Transaction amount: $100.00
- Location: California, USA
- Expected state tax: 7.25%
- Expected total tax: $7.25

**Expected Results**: 
- Tax calculation completes within 200ms
- Correct tax rates applied per jurisdiction
- Tax amounts calculated accurately to the cent
- Complete tax line item breakdown provided

## TC-FUNC-02: Multi-Jurisdiction Tax Calculation
**Prerequisites**: Multi-jurisdiction tax rates configured, address validation active
**Steps**:
1. Create transaction with items shipping to multiple states
2. Request tax calculation for complex transaction
3. Verify jurisdiction-specific tax calculations
4. Validate shipping tax calculations
5. Confirm total tax aggregation

**Test Data**:
- Items shipping to CA, NY, TX
- Item values: $50, $75, $25 respectively
- Different tax rates per state

**Expected Results**:
- Each jurisdiction's taxes calculated separately
- Shipping taxes applied correctly
- Total tax equals sum of jurisdiction taxes
- Jurisdiction precedence properly applied

## TC-FUNC-03: Fee Calculation with Tiers
**Prerequisites**: Tiered fee schedules configured, fee calculation engine active
**Steps**:
1. Submit transactions at different tier levels
2. Calculate fees for each tier
3. Verify correct tier selection
4. Validate tier-specific fee rates
5. Confirm progressive fee calculation

**Test Data**:
- Tier 1 (0-$100): 1.5%
- Tier 2 ($100-$500): 2.0%
- Tier 3 ($500+): 2.5%
- Test amounts: $50, $300, $750

**Expected Results**:
- $50 transaction: 1.5% fee
- $300 transaction: 2.0% fee
- $750 transaction: 2.5% fee
- Correct tier applied for each amount