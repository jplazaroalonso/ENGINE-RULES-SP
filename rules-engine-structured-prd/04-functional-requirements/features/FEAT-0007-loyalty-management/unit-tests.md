# Unit Tests - Loyalty Management

## Test Suite: Loyalty Domain Components

### UT-001: LoyaltyAccount Aggregate Unit Tests

#### Test: Create loyalty account with valid customer
```java
@Test
public void testCreateLoyaltyAccount() {
    // Arrange
    CustomerId customerId = new CustomerId("CUST-123");
    LocalDateTime creationDate = LocalDateTime.now();
    
    // Act
    LoyaltyAccount account = LoyaltyAccount.create(customerId, creationDate);
    
    // Assert
    assertThat(account.getCustomerId()).isEqualTo(customerId);
    assertThat(account.getCurrentTier()).isEqualTo(CustomerTier.BRONZE);
    assertThat(account.getPointsBalance().getAvailablePoints()).isEqualTo(0);
    assertThat(account.getTierEffectiveDate()).isEqualTo(creationDate);
    
    // Verify domain event
    List<DomainEvent> events = account.getUncommittedEvents();
    assertThat(events).hasSize(1);
    assertThat(events.get(0)).isInstanceOf(LoyaltyAccountCreated.class);
}
```

#### Test: Add points to loyalty account
```java
@Test
public void testAddPointsToAccount() {
    // Arrange
    LoyaltyAccount account = createTestAccount();
    TransactionId transactionId = new TransactionId("TXN-456");
    Money purchaseAmount = new Money(BigDecimal.valueOf(100.00), Currency.USD);
    LocalDateTime earnedDate = LocalDateTime.now();
    
    // Act
    account.earnPoints(transactionId, purchaseAmount, earnedDate);
    
    // Assert
    assertThat(account.getPointsBalance().getAvailablePoints()).isEqualTo(100);
    assertThat(account.getPointsBalance().getLifetimeEarned()).isEqualTo(100);
    
    // Verify domain event
    List<DomainEvent> events = account.getUncommittedEvents();
    assertThat(events.stream()
        .anyMatch(e -> e instanceof PointsEarned))
        .isTrue();
}
```

#### Test: Redeem points for reward
```java
@Test
public void testRedeemPointsForReward() {
    // Arrange
    LoyaltyAccount account = createTestAccountWithPoints(5000);
    RewardId rewardId = new RewardId("REWARD-123");
    int pointsCost = 2500;
    
    // Act
    RedemptionResult result = account.redeemPoints(rewardId, pointsCost);
    
    // Assert
    assertThat(result.isSuccessful()).isTrue();
    assertThat(account.getPointsBalance().getAvailablePoints()).isEqualTo(2500);
    assertThat(account.getPointsBalance().getLifetimeRedeemed()).isEqualTo(2500);
    
    // Verify domain event
    List<DomainEvent> events = account.getUncommittedEvents();
    assertThat(events.stream()
        .anyMatch(e -> e instanceof PointsRedeemed))
        .isTrue();
}
```

#### Test: Insufficient points for redemption
```java
@Test
public void testInsufficientPointsRedemption() {
    // Arrange
    LoyaltyAccount account = createTestAccountWithPoints(1000);
    RewardId rewardId = new RewardId("REWARD-123");
    int pointsCost = 2500;
    
    // Act
    RedemptionResult result = account.redeemPoints(rewardId, pointsCost);
    
    // Assert
    assertThat(result.isSuccessful()).isFalse();
    assertThat(result.getErrorReason()).isEqualTo("Insufficient points balance");
    assertThat(account.getPointsBalance().getAvailablePoints()).isEqualTo(1000); // Unchanged
    
    // Verify no redemption event
    List<DomainEvent> events = account.getUncommittedEvents();
    assertThat(events.stream()
        .noneMatch(e -> e instanceof PointsRedeemed))
        .isTrue();
}
```

#### Test: Tier upgrade based on spending
```java
@Test
public void testTierUpgradeBasedOnSpending() {
    // Arrange
    LoyaltyAccount account = createTestAccount();
    Money annualSpending = new Money(BigDecimal.valueOf(5500.00), Currency.USD);
    LocalDateTime calculationDate = LocalDateTime.now();
    
    // Act
    TierCalculationResult result = account.calculateTier(annualSpending, calculationDate);
    
    // Assert
    assertThat(result.getNewTier()).isEqualTo(CustomerTier.GOLD);
    assertThat(result.getTierChanged()).isTrue();
    assertThat(account.getCurrentTier()).isEqualTo(CustomerTier.GOLD);
    
    // Verify domain event
    List<DomainEvent> events = account.getUncommittedEvents();
    assertThat(events.stream()
        .anyMatch(e -> e instanceof CustomerTierUpgraded))
        .isTrue();
}
```

### UT-002: PointsBalance Value Object Unit Tests

#### Test: Create points balance with valid values
```java
@Test
public void testCreatePointsBalance() {
    // Arrange & Act
    PointsBalance balance = new PointsBalance(
        1500,  // available
        3000,  // lifetime earned
        1500   // lifetime redeemed
    );
    
    // Assert
    assertThat(balance.getAvailablePoints()).isEqualTo(1500);
    assertThat(balance.getLifetimeEarned()).isEqualTo(3000);
    assertThat(balance.getLifetimeRedeemed()).isEqualTo(1500);
    assertThat(balance.isValid()).isTrue();
}
```

#### Test: Invalid points balance invariant
```java
@Test
public void testInvalidPointsBalanceInvariant() {
    // Act & Assert
    assertThrows(IllegalArgumentException.class, () -> {
        new PointsBalance(
            2000,  // available
            1500,  // lifetime earned (invalid: earned < available + redeemed)
            500    // lifetime redeemed
        );
    });
}
```

#### Test: Add points to balance
```java
@Test
public void testAddPointsToBalance() {
    // Arrange
    PointsBalance originalBalance = new PointsBalance(1000, 2000, 1000);
    
    // Act
    PointsBalance newBalance = originalBalance.addPoints(500);
    
    // Assert
    assertThat(newBalance.getAvailablePoints()).isEqualTo(1500);
    assertThat(newBalance.getLifetimeEarned()).isEqualTo(2500);
    assertThat(newBalance.getLifetimeRedeemed()).isEqualTo(1000); // Unchanged
}
```

#### Test: Deduct points from balance
```java
@Test
public void testDeductPointsFromBalance() {
    // Arrange
    PointsBalance originalBalance = new PointsBalance(1500, 2500, 1000);
    
    // Act
    PointsBalance newBalance = originalBalance.deductPoints(500);
    
    // Assert
    assertThat(newBalance.getAvailablePoints()).isEqualTo(1000);
    assertThat(newBalance.getLifetimeEarned()).isEqualTo(2500); // Unchanged
    assertThat(newBalance.getLifetimeRedeemed()).isEqualTo(1500);
}
```

### UT-003: CustomerTier Value Object Unit Tests

#### Test: Tier qualification validation
```java
@Test
public void testTierQualificationValidation() {
    // Test cases for tier qualification
    assertThat(CustomerTier.fromAnnualSpending(new Money(BigDecimal.valueOf(500), Currency.USD)))
        .isEqualTo(CustomerTier.BRONZE);
    
    assertThat(CustomerTier.fromAnnualSpending(new Money(BigDecimal.valueOf(2500), Currency.USD)))
        .isEqualTo(CustomerTier.SILVER);
    
    assertThat(CustomerTier.fromAnnualSpending(new Money(BigDecimal.valueOf(7500), Currency.USD)))
        .isEqualTo(CustomerTier.GOLD);
    
    assertThat(CustomerTier.fromAnnualSpending(new Money(BigDecimal.valueOf(15000), Currency.USD)))
        .isEqualTo(CustomerTier.PLATINUM);
}
```

#### Test: Points multiplier by tier
```java
@Test
public void testPointsMultiplierByTier() {
    // Arrange & Act & Assert
    assertThat(CustomerTier.BRONZE.getPointsMultiplier()).isEqualTo(BigDecimal.valueOf(1.0));
    assertThat(CustomerTier.SILVER.getPointsMultiplier()).isEqualTo(BigDecimal.valueOf(1.5));
    assertThat(CustomerTier.GOLD.getPointsMultiplier()).isEqualTo(BigDecimal.valueOf(2.0));
    assertThat(CustomerTier.PLATINUM.getPointsMultiplier()).isEqualTo(BigDecimal.valueOf(3.0));
}
```

#### Test: Tier comparison and ordering
```java
@Test
public void testTierComparisonAndOrdering() {
    // Test tier ordering
    assertThat(CustomerTier.GOLD.compareTo(CustomerTier.SILVER)).isGreaterThan(0);
    assertThat(CustomerTier.SILVER.compareTo(CustomerTier.GOLD)).isLessThan(0);
    assertThat(CustomerTier.GOLD.compareTo(CustomerTier.GOLD)).isEqualTo(0);
    
    // Test tier hierarchy
    assertThat(CustomerTier.PLATINUM.isHigherThan(CustomerTier.GOLD)).isTrue();
    assertThat(CustomerTier.BRONZE.isHigherThan(CustomerTier.SILVER)).isFalse();
}
```

### UT-004: TierCalculationService Domain Service Unit Tests

#### Test: Calculate tier upgrade
```java
@Test
public void testCalculateTierUpgrade() {
    // Arrange
    TierCalculationService service = new TierCalculationService();
    LoyaltyAccount account = createTestAccountWithTier(CustomerTier.SILVER);
    Money annualSpending = new Money(BigDecimal.valueOf(6000), Currency.USD);
    
    // Act
    TierCalculationResult result = service.calculateTier(account, annualSpending);
    
    // Assert
    assertThat(result.getNewTier()).isEqualTo(CustomerTier.GOLD);
    assertThat(result.getTierChanged()).isTrue();
    assertThat(result.getQualificationAmount()).isEqualTo(annualSpending);
    assertThat(result.isUpgrade()).isTrue();
}
```

#### Test: Tier downgrade with grace period
```java
@Test
public void testTierDowngradeWithGracePeriod() {
    // Arrange
    TierCalculationService service = new TierCalculationService();
    LoyaltyAccount account = createTestAccountWithTier(CustomerTier.GOLD);
    Money currentSpending = new Money(BigDecimal.valueOf(3000), Currency.USD);
    int daysBelowThreshold = 95; // Exceeds 90-day grace period
    
    // Act
    TierCalculationResult result = service.calculateTierWithGracePeriod(
        account, currentSpending, daysBelowThreshold);
    
    // Assert
    assertThat(result.getNewTier()).isEqualTo(CustomerTier.SILVER);
    assertThat(result.getTierChanged()).isTrue();
    assertThat(result.isDowngrade()).isTrue();
    assertThat(result.getGracePeriodExpired()).isTrue();
}
```

#### Test: Tier maintained during grace period
```java
@Test
public void testTierMaintainedDuringGracePeriod() {
    // Arrange
    TierCalculationService service = new TierCalculationService();
    LoyaltyAccount account = createTestAccountWithTier(CustomerTier.GOLD);
    Money currentSpending = new Money(BigDecimal.valueOf(3000), Currency.USD);
    int daysBelowThreshold = 60; // Within 90-day grace period
    
    // Act
    TierCalculationResult result = service.calculateTierWithGracePeriod(
        account, currentSpending, daysBelowThreshold);
    
    // Assert
    assertThat(result.getNewTier()).isEqualTo(CustomerTier.GOLD);
    assertThat(result.getTierChanged()).isFalse();
    assertThat(result.getGracePeriodActive()).isTrue();
}
```

### UT-005: PointsEarningService Domain Service Unit Tests

#### Test: Calculate standard points earning
```java
@Test
public void testCalculateStandardPointsEarning() {
    // Arrange
    PointsEarningService service = new PointsEarningService();
    Money purchaseAmount = new Money(BigDecimal.valueOf(150), Currency.USD);
    CustomerTier tier = CustomerTier.GOLD;
    ProductCategory category = ProductCategory.GENERAL;
    
    // Act
    PointsEarningResult result = service.calculateEarning(purchaseAmount, tier, category);
    
    // Assert
    assertThat(result.getBasePoints()).isEqualTo(150);
    assertThat(result.getTierMultiplier()).isEqualTo(BigDecimal.valueOf(2.0));
    assertThat(result.getCategoryMultiplier()).isEqualTo(BigDecimal.valueOf(1.0));
    assertThat(result.getTotalPoints()).isEqualTo(300);
}
```

#### Test: Calculate points with category bonus
```java
@Test
public void testCalculatePointsWithCategoryBonus() {
    // Arrange
    PointsEarningService service = new PointsEarningService();
    Money purchaseAmount = new Money(BigDecimal.valueOf(100), Currency.USD);
    CustomerTier tier = CustomerTier.SILVER;
    ProductCategory category = ProductCategory.ELECTRONICS; // 2x bonus
    
    // Act
    PointsEarningResult result = service.calculateEarning(purchaseAmount, tier, category);
    
    // Assert
    assertThat(result.getBasePoints()).isEqualTo(100);
    assertThat(result.getTierMultiplier()).isEqualTo(BigDecimal.valueOf(1.5));
    assertThat(result.getCategoryMultiplier()).isEqualTo(BigDecimal.valueOf(2.0));
    assertThat(result.getTotalPoints()).isEqualTo(300); // 100 * 1.5 * 2.0
}
```

#### Test: Calculate points with promotional multiplier
```java
@Test
public void testCalculatePointsWithPromotionalMultiplier() {
    // Arrange
    PointsEarningService service = new PointsEarningService();
    Money purchaseAmount = new Money(BigDecimal.valueOf(50), Currency.USD);
    CustomerTier tier = CustomerTier.BRONZE;
    ProductCategory category = ProductCategory.GENERAL;
    PromotionalMultiplier promotion = new PromotionalMultiplier(BigDecimal.valueOf(3.0));
    
    // Act
    PointsEarningResult result = service.calculateEarningWithPromotion(
        purchaseAmount, tier, category, promotion);
    
    // Assert
    assertThat(result.getBasePoints()).isEqualTo(50);
    assertThat(result.getTierMultiplier()).isEqualTo(BigDecimal.valueOf(1.0));
    assertThat(result.getPromotionalMultiplier()).isEqualTo(BigDecimal.valueOf(3.0));
    assertThat(result.getTotalPoints()).isEqualTo(150); // 50 * 1.0 * 3.0
}
```

### UT-006: ExpirationManagementService Domain Service Unit Tests

#### Test: Calculate points expiration schedule
```java
@Test
public void testCalculatePointsExpirationSchedule() {
    // Arrange
    ExpirationManagementService service = new ExpirationManagementService();
    LocalDateTime earnedDate = LocalDateTime.of(2024, 1, 15, 10, 0);
    int pointsAmount = 1000;
    
    // Act
    ExpirationSchedule schedule = service.calculateExpirationSchedule(
        pointsAmount, earnedDate);
    
    // Assert
    assertThat(schedule.getPointsAmount()).isEqualTo(1000);
    assertThat(schedule.getEarnedDate()).isEqualTo(earnedDate);
    assertThat(schedule.getExpirationDate()).isEqualTo(
        earnedDate.plusMonths(24)); // 24-month expiration
    assertThat(schedule.getWarning30DaysDate()).isEqualTo(
        earnedDate.plusMonths(24).minusDays(30));
    assertThat(schedule.getWarning7DaysDate()).isEqualTo(
        earnedDate.plusMonths(24).minusDays(7));
}
```

#### Test: Extend expiration due to activity
```java
@Test
public void testExtendExpirationDueToActivity() {
    // Arrange
    ExpirationManagementService service = new ExpirationManagementService();
    LocalDateTime originalExpiration = LocalDateTime.of(2024, 6, 15, 0, 0);
    LocalDateTime activityDate = LocalDateTime.of(2024, 6, 1, 14, 30);
    
    ExpirationSchedule originalSchedule = new ExpirationSchedule(
        1500, LocalDateTime.of(2022, 6, 15, 0, 0), originalExpiration);
    
    // Act
    ExpirationSchedule extendedSchedule = service.extendExpirationDueToActivity(
        originalSchedule, activityDate);
    
    // Assert
    assertThat(extendedSchedule.getExpirationDate()).isEqualTo(
        activityDate.plusMonths(24)); // Extended to 24 months from activity
    assertThat(extendedSchedule.getPointsAmount()).isEqualTo(1500); // Same amount
}
```

#### Test: Process points expiration
```java
@Test
public void testProcessPointsExpiration() {
    // Arrange
    ExpirationManagementService service = new ExpirationManagementService();
    LocalDateTime currentDate = LocalDateTime.of(2024, 6, 15, 0, 0);
    
    List<ExpirationSchedule> expiringSchedules = Arrays.asList(
        new ExpirationSchedule(500, currentDate.minusMonths(24), currentDate.minusDays(1)),
        new ExpirationSchedule(300, currentDate.minusMonths(24), currentDate),
        new ExpirationSchedule(200, currentDate.minusMonths(24), currentDate.plusDays(1))
    );
    
    // Act
    List<PointsExpirationResult> results = service.processExpiration(
        expiringSchedules, currentDate);
    
    // Assert
    assertThat(results).hasSize(2); // Only first two expire today or before
    assertThat(results.get(0).getExpiredPoints()).isEqualTo(500);
    assertThat(results.get(1).getExpiredPoints()).isEqualTo(300);
    assertThat(results.stream().mapToInt(PointsExpirationResult::getExpiredPoints).sum())
        .isEqualTo(800);
}
```

### UT-007: Domain Events Unit Tests

#### Test: PointsEarned domain event
```java
@Test
public void testPointsEarnedDomainEvent() {
    // Arrange
    CustomerId customerId = new CustomerId("CUST-123");
    TransactionId transactionId = new TransactionId("TXN-456");
    int pointsEarned = 250;
    LocalDateTime earnedAt = LocalDateTime.now();
    
    // Act
    PointsEarned event = new PointsEarned(
        customerId, transactionId, pointsEarned, earnedAt);
    
    // Assert
    assertThat(event.getCustomerId()).isEqualTo(customerId);
    assertThat(event.getTransactionId()).isEqualTo(transactionId);
    assertThat(event.getPointsEarned()).isEqualTo(250);
    assertThat(event.getEarnedAt()).isEqualTo(earnedAt);
    assertThat(event.getEventType()).isEqualTo("loyalty.points.earned");
    assertThat(event.getOccurredAt()).isNotNull();
}
```

#### Test: CustomerTierUpgraded domain event
```java
@Test
public void testCustomerTierUpgradedDomainEvent() {
    // Arrange
    CustomerId customerId = new CustomerId("CUST-789");
    CustomerTier previousTier = CustomerTier.SILVER;
    CustomerTier newTier = CustomerTier.GOLD;
    Money qualificationSpending = new Money(BigDecimal.valueOf(5500), Currency.USD);
    LocalDateTime upgradedAt = LocalDateTime.now();
    
    // Act
    CustomerTierUpgraded event = new CustomerTierUpgraded(
        customerId, previousTier, newTier, qualificationSpending, upgradedAt);
    
    // Assert
    assertThat(event.getCustomerId()).isEqualTo(customerId);
    assertThat(event.getPreviousTier()).isEqualTo(CustomerTier.SILVER);
    assertThat(event.getNewTier()).isEqualTo(CustomerTier.GOLD);
    assertThat(event.getQualificationSpending()).isEqualTo(qualificationSpending);
    assertThat(event.getUpgradedAt()).isEqualTo(upgradedAt);
    assertThat(event.getEventType()).isEqualTo("loyalty.tier.upgraded");
}
```

## Test Utilities and Builders

### Test Data Builders
```java
public class LoyaltyTestDataBuilder {
    
    public static LoyaltyAccount createTestAccount() {
        return LoyaltyAccount.create(
            new CustomerId("TEST-CUST-001"),
            LocalDateTime.of(2024, 1, 1, 0, 0)
        );
    }
    
    public static LoyaltyAccount createTestAccountWithPoints(int points) {
        LoyaltyAccount account = createTestAccount();
        account.earnPoints(
            new TransactionId("TEST-TXN-001"),
            new Money(BigDecimal.valueOf(points), Currency.USD),
            LocalDateTime.now()
        );
        account.markEventsAsCommitted(); // Clear events for clean testing
        return account;
    }
    
    public static LoyaltyAccount createTestAccountWithTier(CustomerTier tier) {
        LoyaltyAccount account = createTestAccount();
        // Simulate spending to achieve tier
        Money requiredSpending = getMinimumSpendingForTier(tier);
        account.calculateTier(requiredSpending, LocalDateTime.now());
        account.markEventsAsCommitted();
        return account;
    }
    
    private static Money getMinimumSpendingForTier(CustomerTier tier) {
        return switch (tier) {
            case BRONZE -> new Money(BigDecimal.ZERO, Currency.USD);
            case SILVER -> new Money(BigDecimal.valueOf(1000), Currency.USD);
            case GOLD -> new Money(BigDecimal.valueOf(5000), Currency.USD);
            case PLATINUM -> new Money(BigDecimal.valueOf(10000), Currency.USD);
        };
    }
}
```

### Mock Repositories
```java
@ExtendWith(MockitoExtension.class)
class LoyaltyServiceTest {
    
    @Mock
    private LoyaltyAccountRepository loyaltyAccountRepository;
    
    @Mock
    private RewardsCatalogRepository rewardsCatalogRepository;
    
    @Mock
    private PartnerRepository partnerRepository;
    
    @InjectMocks
    private LoyaltyService loyaltyService;
    
    @Test
    public void testServiceMethod() {
        // Test service methods that coordinate multiple aggregates
    }
}
```

## Test Coverage Requirements

### Minimum Coverage Targets
- **Line Coverage**: 95%
- **Branch Coverage**: 90%
- **Method Coverage**: 98%
- **Class Coverage**: 100%

### Critical Path Coverage
- All business rule validations must have 100% test coverage
- All domain event publishing scenarios must be tested
- All error conditions and edge cases must be covered
- All value object invariant violations must be tested

### Performance Test Targets
- Unit tests should complete within 50ms each
- Test suite should complete within 5 minutes
- No external dependencies in unit tests
- All tests should be deterministic and repeatable
