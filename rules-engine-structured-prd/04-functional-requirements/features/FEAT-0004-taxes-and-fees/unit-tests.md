# Unit Tests - Taxes and Fees Rules

## TC-UNIT-01: Tax Amount Calculation
**Test Objective**: Verify accurate tax amount calculation logic
**Test Scope**: TaxCalculationService.calculateTaxAmount method

**Test Cases**:
```java
@Test
public void testBasicTaxCalculation() {
    // Given
    BigDecimal taxableAmount = new BigDecimal("100.00");
    BigDecimal taxRate = new BigDecimal("0.0825");
    
    // When
    BigDecimal taxAmount = taxCalculationService.calculateTaxAmount(taxableAmount, taxRate);
    
    // Then
    assertEquals(new BigDecimal("8.25"), taxAmount);
}

@Test
public void testZeroTaxRate() {
    // Given
    BigDecimal taxableAmount = new BigDecimal("100.00");
    BigDecimal taxRate = BigDecimal.ZERO;
    
    // When
    BigDecimal taxAmount = taxCalculationService.calculateTaxAmount(taxableAmount, taxRate);
    
    // Then
    assertEquals(BigDecimal.ZERO, taxAmount);
}
```

## TC-UNIT-02: Fee Calculation Logic
**Test Objective**: Verify fee calculation accuracy
**Test Scope**: FeeCalculationService.calculateFees method

**Test Cases**:
```java
@Test
public void testPercentageFeeCalculation() {
    // Given
    BigDecimal feeableAmount = new BigDecimal("500.00");
    BigDecimal feeRate = new BigDecimal("0.025");
    
    // When
    BigDecimal feeAmount = feeCalculationService.calculateFee(feeableAmount, feeRate);
    
    // Then
    assertEquals(new BigDecimal("12.50"), feeAmount);
}
```