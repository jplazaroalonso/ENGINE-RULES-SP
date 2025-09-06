# Unit Tests - Coupons Management

## Test Suite: Coupon Aggregate Unit Tests

### UT-001: CouponCode Value Object Tests

#### Test: Valid coupon code creation
```java
@Test
public void testValidCouponCodeCreation() {
    // Arrange
    String validCode = "SAVE20-2024";
    
    // Act
    CouponCode couponCode = CouponCode.create(validCode);
    
    // Assert
    assertThat(couponCode.getValue()).isEqualTo("SAVE20-2024");
    assertThat(couponCode.isValid()).isTrue();
}
```

#### Test: Invalid coupon code rejection
```java
@Test
public void testInvalidCouponCodeRejection() {
    // Arrange
    String invalidCode = ""; // Empty string
    
    // Act & Assert
    assertThrows(IllegalArgumentException.class, () -> {
        CouponCode.create(invalidCode);
    });
}
```

#### Test: Coupon code pattern validation
```java
@Test
public void testCouponCodePatternValidation() {
    // Arrange
    CodePattern pattern = CodePattern.create("SAVE-XXXX");
    
    // Act
    CouponCode code1 = CouponCode.generate(pattern);
    CouponCode code2 = CouponCode.generate(pattern);
    
    // Assert
    assertThat(code1.getValue()).startsWith("SAVE-");
    assertThat(code2.getValue()).startsWith("SAVE-");
    assertThat(code1.getValue()).isNotEqualTo(code2.getValue()); // Uniqueness
}
```

#### Test: Checksum validation
```java
@Test
public void testChecksumValidation() {
    // Arrange
    String codeWithValidChecksum = "SAVE20-ABC3";
    String codeWithInvalidChecksum = "SAVE20-ABC4";
    
    // Act
    CouponCode validCode = CouponCode.create(codeWithValidChecksum);
    
    // Assert
    assertThat(validCode.hasValidChecksum()).isTrue();
    assertThrows(IllegalArgumentException.class, () -> {
        CouponCode.create(codeWithInvalidChecksum);
    });
}
```

---

### UT-002: DiscountValue Value Object Tests

#### Test: Fixed amount discount creation
```java
@Test
public void testFixedAmountDiscountCreation() {
    // Arrange
    BigDecimal amount = new BigDecimal("10.00");
    Currency currency = Currency.getInstance("USD");
    
    // Act
    DiscountValue discount = DiscountValue.fixedAmount(amount, currency);
    
    // Assert
    assertThat(discount.getType()).isEqualTo(DiscountType.FIXED_AMOUNT);
    assertThat(discount.getAmount()).isEqualTo(amount);
    assertThat(discount.getCurrency()).isEqualTo(currency);
}
```

#### Test: Percentage discount creation
```java
@Test
public void testPercentageDiscountCreation() {
    // Arrange
    BigDecimal percentage = new BigDecimal("15.5");
    
    // Act
    DiscountValue discount = DiscountValue.percentage(percentage);
    
    // Assert
    assertThat(discount.getType()).isEqualTo(DiscountType.PERCENTAGE);
    assertThat(discount.getPercentage()).isEqualTo(percentage);
}
```

#### Test: Invalid percentage rejection
```java
@Test
public void testInvalidPercentageRejection() {
    // Arrange
    BigDecimal invalidPercentage = new BigDecimal("150.0"); // Over 100%
    
    // Act & Assert
    assertThrows(IllegalArgumentException.class, () -> {
        DiscountValue.percentage(invalidPercentage);
    });
}
```

#### Test: Discount calculation for transaction
```java
@Test
public void testDiscountCalculationForTransaction() {
    // Arrange
    DiscountValue percentageDiscount = DiscountValue.percentage(new BigDecimal("20"));
    TransactionContext context = TransactionContext.builder()
        .subtotal(new Money(new BigDecimal("100.00"), Currency.getInstance("USD")))
        .build();
    
    // Act
    Money calculatedDiscount = percentageDiscount.calculateDiscount(context);
    
    // Assert
    assertThat(calculatedDiscount.getAmount()).isEqualTo(new BigDecimal("20.00"));
}
```

---

### UT-003: Coupon Aggregate Tests

#### Test: Create new coupon
```java
@Test
public void testCreateNewCoupon() {
    // Arrange
    CouponSpecification spec = CouponSpecification.builder()
        .name(CouponName.create("Summer Sale"))
        .discountValue(DiscountValue.percentage(new BigDecimal("15")))
        .validityPeriod(ValidityPeriod.create(
            LocalDateTime.of(2024, 1, 1, 0, 0),
            LocalDateTime.of(2024, 1, 31, 23, 59)
        ))
        .usageRules(UsageRules.singleUse())
        .build();
    
    // Act
    Coupon coupon = Coupon.create(spec);
    
    // Assert
    assertThat(coupon.getId()).isNotNull();
    assertThat(coupon.getStatus()).isEqualTo(CouponStatus.DRAFT);
    assertThat(coupon.getName().getValue()).isEqualTo("Summer Sale");
    assertThat(coupon.getDiscountValue().getPercentage()).isEqualTo(new BigDecimal("15"));
}
```

#### Test: Activate approved coupon
```java
@Test
public void testActivateApprovedCoupon() {
    // Arrange
    Coupon coupon = createTestCoupon();
    coupon.approve(UserId.create("approver123"));
    
    // Act
    coupon.activate();
    
    // Assert
    assertThat(coupon.getStatus()).isEqualTo(CouponStatus.ACTIVE);
    List<DomainEvent> events = coupon.getUncommittedEvents();
    assertThat(events).hasSize(1);
    assertThat(events.get(0)).isInstanceOf(CouponActivated.class);
}
```

#### Test: Prevent activation of unapproved coupon
```java
@Test
public void testPreventActivationOfUnapprovedCoupon() {
    // Arrange
    Coupon coupon = createTestCoupon(); // Status is DRAFT
    
    // Act & Assert
    assertThrows(IllegalStateException.class, () -> {
        coupon.activate();
    });
    assertThat(coupon.getStatus()).isEqualTo(CouponStatus.DRAFT);
}
```

#### Test: Validate coupon usage eligibility
```java
@Test
public void testValidateCouponUsageEligibility() {
    // Arrange
    Coupon activeCoupon = createActiveCoupon();
    UsageContext context = UsageContext.builder()
        .customerId(CustomerId.create("customer123"))
        .transactionAmount(new Money(new BigDecimal("75.00"), Currency.getInstance("USD")))
        .transactionDate(LocalDateTime.now())
        .build();
    
    // Act
    ValidationResult result = activeCoupon.validateUsage(context);
    
    // Assert
    assertThat(result.isValid()).isTrue();
    assertThat(result.getErrors()).isEmpty();
}
```

#### Test: Reject usage when minimum purchase not met
```java
@Test
public void testRejectUsageWhenMinimumPurchaseNotMet() {
    // Arrange
    Coupon coupon = createCouponWithMinimumPurchase(new BigDecimal("50.00"));
    UsageContext context = UsageContext.builder()
        .customerId(CustomerId.create("customer123"))
        .transactionAmount(new Money(new BigDecimal("25.00"), Currency.getInstance("USD")))
        .build();
    
    // Act
    ValidationResult result = coupon.validateUsage(context);
    
    // Assert
    assertThat(result.isValid()).isFalse();
    assertThat(result.getErrors()).containsExactly("Minimum purchase amount not met");
}
```

---

### UT-004: UsageRules Value Object Tests

#### Test: Single-use coupon rules
```java
@Test
public void testSingleUseCouponRules() {
    // Arrange
    UsageRules rules = UsageRules.singleUse();
    CustomerId customerId = CustomerId.create("customer123");
    
    // Act
    boolean canUseFirst = rules.canCustomerUse(customerId, 0); // No previous uses
    boolean canUseSecond = rules.canCustomerUse(customerId, 1); // One previous use
    
    // Assert
    assertThat(canUseFirst).isTrue();
    assertThat(canUseSecond).isFalse();
}
```

#### Test: Multi-use coupon rules
```java
@Test
public void testMultiUseCouponRules() {
    // Arrange
    UsageRules rules = UsageRules.multiUse(3); // Max 3 uses per customer
    CustomerId customerId = CustomerId.create("customer123");
    
    // Act
    boolean canUseThird = rules.canCustomerUse(customerId, 2); // Two previous uses
    boolean canUseFourth = rules.canCustomerUse(customerId, 3); // Three previous uses
    
    // Assert
    assertThat(canUseThird).isTrue();
    assertThat(canUseFourth).isFalse();
}
```

#### Test: Channel restrictions
```java
@Test
public void testChannelRestrictions() {
    // Arrange
    UsageRules rules = UsageRules.builder()
        .allowedChannels(Set.of(Channel.ONLINE, Channel.MOBILE))
        .build();
    
    // Act
    boolean canUseOnline = rules.isChannelAllowed(Channel.ONLINE);
    boolean canUseInStore = rules.isChannelAllowed(Channel.IN_STORE);
    
    // Assert
    assertThat(canUseOnline).isTrue();
    assertThat(canUseInStore).isFalse();
}
```

---

### UT-005: RedemptionTracking Tests

#### Test: Record successful redemption
```java
@Test
public void testRecordSuccessfulRedemption() {
    // Arrange
    RedemptionTracking tracking = RedemptionTracking.create();
    RedemptionDetails details = RedemptionDetails.builder()
        .customerId(CustomerId.create("customer123"))
        .transactionId(TransactionId.create("txn456"))
        .discountApplied(new Money(new BigDecimal("15.00"), Currency.getInstance("USD")))
        .channel(Channel.ONLINE)
        .redemptionTime(LocalDateTime.now())
        .build();
    
    // Act
    tracking.recordRedemption(details);
    
    // Assert
    assertThat(tracking.getTotalRedemptions()).isEqualTo(1);
    assertThat(tracking.getRedemptionHistory()).hasSize(1);
    assertThat(tracking.getRedemptionHistory().get(0).getCustomerId())
        .isEqualTo(CustomerId.create("customer123"));
}
```

#### Test: Calculate usage rate
```java
@Test
public void testCalculateUsageRate() {
    // Arrange
    RedemptionTracking tracking = RedemptionTracking.create();
    tracking.setMaximumUses(100);
    
    // Simulate 25 redemptions
    for (int i = 0; i < 25; i++) {
        tracking.recordRedemption(createTestRedemptionDetails());
    }
    
    // Act
    BigDecimal usageRate = tracking.calculateUsageRate();
    
    // Assert
    assertThat(usageRate).isEqualTo(new BigDecimal("25.00")); // 25%
}
```

#### Test: Detect anomalous usage patterns
```java
@Test
public void testDetectAnomalousUsagePatterns() {
    // Arrange
    RedemptionTracking tracking = RedemptionTracking.create();
    LocalDateTime baseTime = LocalDateTime.now();
    
    // Simulate rapid redemptions (10 within 1 minute)
    for (int i = 0; i < 10; i++) {
        RedemptionDetails details = createRedemptionDetails(baseTime.plusSeconds(i * 5));
        tracking.recordRedemption(details);
    }
    
    // Act
    List<AnomalyDetection> anomalies = tracking.detectAnomalies();
    
    // Assert
    assertThat(anomalies).isNotEmpty();
    assertThat(anomalies.get(0).getType()).isEqualTo(AnomalyType.RAPID_USAGE);
}
```

---

### UT-006: Domain Services Unit Tests

#### Test: CouponValidationService - Valid coupon validation
```java
@Test
public void testValidCouponValidation() {
    // Arrange
    CouponValidationService service = new CouponValidationService();
    Coupon validCoupon = createActiveCoupon();
    UsageContext context = createValidUsageContext();
    
    // Act
    ValidationResult result = service.validateCoupon(validCoupon.getCode(), context);
    
    // Assert
    assertThat(result.isValid()).isTrue();
    assertThat(result.getValidatedCoupon()).isEqualTo(validCoupon);
    assertThat(result.getCalculatedDiscount()).isNotNull();
}
```

#### Test: FraudDetectionService - Detect suspicious patterns
```java
@Test
public void testDetectSuspiciousPatterns() {
    // Arrange
    FraudDetectionService service = new FraudDetectionService();
    CouponUsage usage = CouponUsage.builder()
        .couponCode(CouponCode.create("SUSPICIOUS"))
        .ipAddress("192.168.1.100")
        .usageCount(50)
        .timeWindow(Duration.ofMinutes(5))
        .build();
    
    // Act
    FraudRisk risk = service.analyzeSuspiciousActivity(usage);
    
    // Assert
    assertThat(risk.getRiskLevel()).isEqualTo(RiskLevel.HIGH);
    assertThat(risk.getDetectedPatterns()).contains(FraudPattern.RAPID_USAGE);
}
```

#### Test: CouponDistributionService - Distribute to segment
```java
@Test
public void testDistributeToSegment() {
    // Arrange
    CouponDistributionService service = new CouponDistributionService();
    CouponCampaign campaign = createTestCampaign();
    CustomerSegment segment = CustomerSegment.goldTierCustomers();
    
    // Act
    DistributionResult result = service.distributeToSegment(campaign, segment);
    
    // Assert
    assertThat(result.isSuccessful()).isTrue();
    assertThat(result.getDistributedCount()).isGreaterThan(0);
    assertThat(result.getFailedDeliveries()).isEmpty();
}
```

---

### UT-007: Event Tests

#### Test: CouponCreated event
```java
@Test
public void testCouponCreatedEvent() {
    // Arrange
    CouponId couponId = CouponId.create();
    CouponCode couponCode = CouponCode.create("TEST123");
    UserId createdBy = UserId.create("user123");
    LocalDateTime createdAt = LocalDateTime.now();
    
    // Act
    CouponCreated event = new CouponCreated(couponId, couponCode, createdBy, createdAt);
    
    // Assert
    assertThat(event.getCouponId()).isEqualTo(couponId);
    assertThat(event.getCouponCode()).isEqualTo(couponCode);
    assertThat(event.getCreatedBy()).isEqualTo(createdBy);
    assertThat(event.getOccurredAt()).isEqualTo(createdAt);
    assertThat(event.getEventType()).isEqualTo("coupon.created.v1");
}
```

#### Test: CouponRedeemed event
```java
@Test
public void testCouponRedeemedEvent() {
    // Arrange
    CouponId couponId = CouponId.create();
    CustomerId customerId = CustomerId.create("customer123");
    TransactionId transactionId = TransactionId.create("txn456");
    Money discountApplied = new Money(new BigDecimal("15.00"), Currency.getInstance("USD"));
    Channel channel = Channel.ONLINE;
    LocalDateTime redeemedAt = LocalDateTime.now();
    
    // Act
    CouponRedeemed event = new CouponRedeemed(
        couponId, customerId, transactionId, discountApplied, channel, redeemedAt);
    
    // Assert
    assertThat(event.getCouponId()).isEqualTo(couponId);
    assertThat(event.getCustomerId()).isEqualTo(customerId);
    assertThat(event.getTransactionId()).isEqualTo(transactionId);
    assertThat(event.getDiscountApplied()).isEqualTo(discountApplied);
    assertThat(event.getChannel()).isEqualTo(channel);
    assertThat(event.getEventType()).isEqualTo("coupon.redeemed.v1");
}
```

---

## Test Utilities and Builders

### Test Data Builders
```java
public class CouponTestDataBuilder {
    
    public static Coupon createTestCoupon() {
        return Coupon.create(CouponSpecification.builder()
            .name(CouponName.create("Test Coupon"))
            .discountValue(DiscountValue.percentage(new BigDecimal("10")))
            .validityPeriod(ValidityPeriod.create(
                LocalDateTime.now(),
                LocalDateTime.now().plusDays(30)
            ))
            .usageRules(UsageRules.singleUse())
            .build());
    }
    
    public static Coupon createActiveCoupon() {
        Coupon coupon = createTestCoupon();
        coupon.approve(UserId.create("approver"));
        coupon.activate();
        return coupon;
    }
    
    public static UsageContext createValidUsageContext() {
        return UsageContext.builder()
            .customerId(CustomerId.create("customer123"))
            .transactionAmount(new Money(new BigDecimal("100.00"), Currency.getInstance("USD")))
            .transactionDate(LocalDateTime.now())
            .channel(Channel.ONLINE)
            .build();
    }
}
```

### Mock Objects
```java
@ExtendWith(MockitoExtension.class)
class CouponAggregateTest {
    
    @Mock
    private CouponRepository couponRepository;
    
    @Mock
    private CustomerRepository customerRepository;
    
    @Mock
    private FraudDetectionService fraudDetectionService;
    
    // Test methods using mocks...
}
```

## Test Coverage Requirements

### Minimum Coverage Targets
- **Line Coverage**: 90%
- **Branch Coverage**: 85%
- **Method Coverage**: 95%
- **Class Coverage**: 100%

### Critical Path Coverage
- All domain invariants must have 100% test coverage
- All business rule validations must be thoroughly tested
- All error conditions and edge cases must be covered
- All domain events must have corresponding tests

### Performance Test Targets
- Unit tests should complete within 50ms each
- Test suite should complete within 5 minutes
- No external dependencies in unit tests
- All tests should be deterministic and repeatable
