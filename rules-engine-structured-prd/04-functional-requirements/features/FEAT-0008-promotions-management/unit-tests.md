# Unit Tests - Promotions Management

## Test Suite: Promotional Campaign Domain Components

### UT-001: PromotionalCampaign Aggregate Unit Tests

#### Test: Create promotional campaign with valid configuration
```java
@Test
public void testCreatePromotionalCampaign() {
    // Arrange
    CampaignId campaignId = new CampaignId("PROMO-SPRING-2024");
    CampaignName campaignName = new CampaignName("Spring Sale 2024");
    CampaignType campaignType = CampaignType.PERCENTAGE_DISCOUNT;
    TargetingCriteria targeting = new TargetingCriteria(
        CustomerSegments.of(CustomerSegment.GOLD_TIER),
        ProductCategories.of("ELECTRONICS", "CLOTHING")
    );
    ValidityPeriod validity = new ValidityPeriod(
        LocalDateTime.of(2024, 3, 1, 0, 0),
        LocalDateTime.of(2024, 3, 31, 23, 59)
    );
    Budget budget = new Budget(new Money(BigDecimal.valueOf(10000), Currency.USD));
    
    // Act
    PromotionalCampaign campaign = PromotionalCampaign.create(
        campaignId, campaignName, campaignType, targeting, validity, budget
    );
    
    // Assert
    assertThat(campaign.getId()).isEqualTo(campaignId);
    assertThat(campaign.getName()).isEqualTo(campaignName);
    assertThat(campaign.getType()).isEqualTo(campaignType);
    assertThat(campaign.getStatus()).isEqualTo(CampaignStatus.DRAFT);
    assertThat(campaign.getBudget().getRemainingBudget()).isEqualTo(budget.getTotalBudget());
    
    // Verify domain event
    List<DomainEvent> events = campaign.getUncommittedEvents();
    assertThat(events).hasSize(1);
    assertThat(events.get(0)).isInstanceOf(PromotionalCampaignCreated.class);
}
```

#### Test: Activate campaign with proper validation
```java
@Test
public void testActivateCampaign() {
    // Arrange
    PromotionalCampaign campaign = createTestCampaign();
    campaign.approve(new ApprovalDecision(true, "Approved by marketing director"));
    LocalDateTime activationTime = LocalDateTime.now();
    
    // Act
    campaign.activate(activationTime);
    
    // Assert
    assertThat(campaign.getStatus()).isEqualTo(CampaignStatus.ACTIVE);
    assertThat(campaign.getActivatedAt()).isEqualTo(activationTime);
    
    // Verify domain event
    List<DomainEvent> events = campaign.getUncommittedEvents();
    assertThat(events.stream()
        .anyMatch(e -> e instanceof PromotionalCampaignActivated))
        .isTrue();
}
```

#### Test: Campaign activation fails without approval
```java
@Test
public void testActivateCampaignWithoutApproval() {
    // Arrange
    PromotionalCampaign campaign = createTestCampaign();
    // Note: No approval given
    
    // Act & Assert
    assertThrows(CampaignNotApprovedException.class, () -> {
        campaign.activate(LocalDateTime.now());
    });
    
    assertThat(campaign.getStatus()).isEqualTo(CampaignStatus.DRAFT);
}
```

#### Test: Apply promotion and track budget consumption
```java
@Test
public void testApplyPromotionAndTrackBudget() {
    // Arrange
    PromotionalCampaign campaign = createActiveCampaign();
    CustomerId customerId = new CustomerId("CUST-123");
    TransactionId transactionId = new TransactionId("TXN-456");
    Money discountAmount = new Money(BigDecimal.valueOf(150.00), Currency.USD);
    
    // Act
    PromotionApplication application = campaign.applyPromotion(
        customerId, transactionId, discountAmount
    );
    
    // Assert
    assertThat(application.getApplicationId()).isNotNull();
    assertThat(application.getCampaignId()).isEqualTo(campaign.getId());
    assertThat(application.getAppliedBenefit()).isEqualTo(discountAmount);
    assertThat(campaign.getBudget().getSpentAmount()).isEqualTo(discountAmount);
    assertThat(campaign.getBudget().getRemainingBudget())
        .isEqualTo(new Money(BigDecimal.valueOf(9850.00), Currency.USD));
    
    // Verify domain event
    List<DomainEvent> events = campaign.getUncommittedEvents();
    assertThat(events.stream()
        .anyMatch(e -> e instanceof PromotionApplied))
        .isTrue();
}
```

#### Test: Campaign deactivates when budget exhausted
```java
@Test
public void testCampaignDeactivatesWhenBudgetExhausted() {
    // Arrange
    Budget smallBudget = new Budget(new Money(BigDecimal.valueOf(100.00), Currency.USD));
    PromotionalCampaign campaign = createActiveCampaignWithBudget(smallBudget);
    Money largeDiscount = new Money(BigDecimal.valueOf(150.00), Currency.USD);
    
    // Act & Assert
    assertThrows(InsufficientBudgetException.class, () -> {
        campaign.applyPromotion(
            new CustomerId("CUST-123"),
            new TransactionId("TXN-456"),
            largeDiscount
        );
    });
    
    // Verify budget remains unchanged
    assertThat(campaign.getBudget().getSpentAmount())
        .isEqualTo(new Money(BigDecimal.ZERO, Currency.USD));
    assertThat(campaign.getStatus()).isEqualTo(CampaignStatus.ACTIVE);
}
```

### UT-002: DiscountConfiguration Value Object Unit Tests

#### Test: Create percentage discount configuration
```java
@Test
public void testCreatePercentageDiscountConfiguration() {
    // Arrange & Act
    DiscountConfiguration config = DiscountConfiguration.percentage(
        BigDecimal.valueOf(20.0),           // 20% discount
        new Money(BigDecimal.valueOf(50.00), Currency.USD),  // minimum purchase
        new Money(BigDecimal.valueOf(100.00), Currency.USD)  // maximum discount
    );
    
    // Assert
    assertThat(config.getDiscountType()).isEqualTo(DiscountType.PERCENTAGE);
    assertThat(config.getDiscountValue()).isEqualTo(BigDecimal.valueOf(20.0));
    assertThat(config.getMinimumPurchase()).isEqualTo(new Money(BigDecimal.valueOf(50.00), Currency.USD));
    assertThat(config.getMaximumDiscount()).isEqualTo(new Money(BigDecimal.valueOf(100.00), Currency.USD));
}
```

#### Test: Create fixed amount discount configuration
```java
@Test
public void testCreateFixedAmountDiscountConfiguration() {
    // Arrange & Act
    DiscountConfiguration config = DiscountConfiguration.fixedAmount(
        new Money(BigDecimal.valueOf(25.00), Currency.USD),  // $25 discount
        new Money(BigDecimal.valueOf(100.00), Currency.USD)  // minimum purchase
    );
    
    // Assert
    assertThat(config.getDiscountType()).isEqualTo(DiscountType.FIXED_AMOUNT);
    assertThat(config.getFixedDiscountAmount()).isEqualTo(new Money(BigDecimal.valueOf(25.00), Currency.USD));
    assertThat(config.getMinimumPurchase()).isEqualTo(new Money(BigDecimal.valueOf(100.00), Currency.USD));
    assertThat(config.getMaximumDiscount()).isNull(); // No max for fixed amount
}
```

#### Test: Create buy-X-get-Y discount configuration
```java
@Test
public void testCreateBuyXGetYDiscountConfiguration() {
    // Arrange & Act
    DiscountConfiguration config = DiscountConfiguration.buyXGetY(
        2,  // buy quantity
        1,  // get quantity
        BigDecimal.valueOf(100.0),  // 100% discount (free)
        3   // maximum free items
    );
    
    // Assert
    assertThat(config.getDiscountType()).isEqualTo(DiscountType.BUY_X_GET_Y);
    assertThat(config.getBuyQuantity()).isEqualTo(2);
    assertThat(config.getGetQuantity()).isEqualTo(1);
    assertThat(config.getGetDiscountPercentage()).isEqualTo(BigDecimal.valueOf(100.0));
    assertThat(config.getMaximumFreeItems()).isEqualTo(3);
}
```

#### Test: Invalid discount configuration validation
```java
@Test
public void testInvalidDiscountConfigurationValidation() {
    // Test: Percentage over 100%
    assertThrows(IllegalArgumentException.class, () -> {
        DiscountConfiguration.percentage(
            BigDecimal.valueOf(150.0),  // Invalid: over 100%
            Money.zero(Currency.USD),
            Money.zero(Currency.USD)
        );
    });
    
    // Test: Negative fixed amount
    assertThrows(IllegalArgumentException.class, () -> {
        DiscountConfiguration.fixedAmount(
            new Money(BigDecimal.valueOf(-10.00), Currency.USD),  // Invalid: negative
            Money.zero(Currency.USD)
        );
    });
    
    // Test: Invalid buy-X-get-Y quantities
    assertThrows(IllegalArgumentException.class, () -> {
        DiscountConfiguration.buyXGetY(0, 1, BigDecimal.valueOf(50.0), 1);  // Invalid: buy quantity = 0
    });
}
```

### UT-003: CustomerSegments Value Object Unit Tests

#### Test: Create customer segments with multiple criteria
```java
@Test
public void testCreateCustomerSegmentsWithMultipleCriteria() {
    // Arrange & Act
    CustomerSegments segments = CustomerSegments.builder()
        .addSegment(CustomerSegment.NEW_CUSTOMERS)
        .addSegment(CustomerSegment.LOYAL_CUSTOMERS)
        .addTierBasedSegment(LoyaltyTier.GOLD)
        .addTierBasedSegment(LoyaltyTier.PLATINUM)
        .addGeographicSegment("CA", "NY", "FL")
        .build();
    
    // Assert
    assertThat(segments.includes(CustomerSegment.NEW_CUSTOMERS)).isTrue();
    assertThat(segments.includes(CustomerSegment.LOYAL_CUSTOMERS)).isTrue();
    assertThat(segments.includesTier(LoyaltyTier.GOLD)).isTrue();
    assertThat(segments.includesTier(LoyaltyTier.SILVER)).isFalse();
    assertThat(segments.includesGeography("CA")).isTrue();
    assertThat(segments.includesGeography("TX")).isFalse();
}
```

#### Test: Customer segment matching logic
```java
@Test
public void testCustomerSegmentMatching() {
    // Arrange
    CustomerSegments segments = CustomerSegments.of(
        CustomerSegment.LOYAL_CUSTOMERS,
        LoyaltyTier.GOLD
    );
    
    Customer loyalGoldCustomer = Customer.builder()
        .customerId(new CustomerId("CUST-001"))
        .segment(CustomerSegment.LOYAL_CUSTOMERS)
        .loyaltyTier(LoyaltyTier.GOLD)
        .geography("CA")
        .build();
    
    Customer newSilverCustomer = Customer.builder()
        .customerId(new CustomerId("CUST-002"))
        .segment(CustomerSegment.NEW_CUSTOMERS)
        .loyaltyTier(LoyaltyTier.SILVER)
        .geography("NY")
        .build();
    
    // Act & Assert
    assertThat(segments.matches(loyalGoldCustomer)).isTrue();
    assertThat(segments.matches(newSilverCustomer)).isFalse();
}
```

#### Test: Segment exclusion rules
```java
@Test
public void testSegmentExclusionRules() {
    // Arrange
    CustomerSegments segments = CustomerSegments.builder()
        .addSegment(CustomerSegment.ALL_CUSTOMERS)
        .excludeSegment(CustomerSegment.VIP_CUSTOMERS)
        .excludeTier(LoyaltyTier.PLATINUM)
        .build();
    
    Customer regularCustomer = Customer.builder()
        .segment(CustomerSegment.LOYAL_CUSTOMERS)
        .loyaltyTier(LoyaltyTier.GOLD)
        .build();
    
    Customer vipCustomer = Customer.builder()
        .segment(CustomerSegment.VIP_CUSTOMERS)
        .loyaltyTier(LoyaltyTier.PLATINUM)
        .build();
    
    // Act & Assert
    assertThat(segments.matches(regularCustomer)).isTrue();
    assertThat(segments.matches(vipCustomer)).isFalse();
}
```

### UT-004: CampaignEligibilityService Domain Service Unit Tests

#### Test: Evaluate customer eligibility for percentage discount
```java
@Test
public void testEvaluateCustomerEligibilityForPercentageDiscount() {
    // Arrange
    CampaignEligibilityService service = new CampaignEligibilityService();
    
    PromotionalCampaign campaign = createCampaignWithTargeting(
        CustomerSegments.of(CustomerSegment.GOLD_TIER)
    );
    
    Customer eligibleCustomer = Customer.builder()
        .segment(CustomerSegment.LOYAL_CUSTOMERS)
        .loyaltyTier(LoyaltyTier.GOLD)
        .build();
    
    Transaction transaction = Transaction.builder()
        .amount(new Money(BigDecimal.valueOf(150.00), Currency.USD))
        .items(createEligibleItems())
        .build();
    
    // Act
    EligibilityResult result = service.evaluateEligibility(eligibleCustomer, transaction, campaign);
    
    // Assert
    assertThat(result.isEligible()).isTrue();
    assertThat(result.getEligibilityReason()).isEqualTo("Customer matches target segment and meets all criteria");
    assertThat(result.getEligibleAmount()).isEqualTo(transaction.getAmount());
}
```

#### Test: Evaluate ineligible customer for targeted promotion
```java
@Test
public void testEvaluateIneligibleCustomerForTargetedPromotion() {
    // Arrange
    CampaignEligibilityService service = new CampaignEligibilityService();
    
    PromotionalCampaign campaign = createCampaignWithTargeting(
        CustomerSegments.of(CustomerSegment.NEW_CUSTOMERS)
    );
    
    Customer ineligibleCustomer = Customer.builder()
        .segment(CustomerSegment.LOYAL_CUSTOMERS)
        .loyaltyTier(LoyaltyTier.GOLD)
        .build();
    
    Transaction transaction = createStandardTransaction();
    
    // Act
    EligibilityResult result = service.evaluateEligibility(ineligibleCustomer, transaction, campaign);
    
    // Assert
    assertThat(result.isEligible()).isFalse();
    assertThat(result.getIneligibilityReason()).isEqualTo("Customer does not match target segment");
    assertThat(result.getEligibleAmount()).isEqualTo(Money.zero(Currency.USD));
}
```

#### Test: Geographic eligibility validation
```java
@Test
public void testGeographicEligibilityValidation() {
    // Arrange
    CampaignEligibilityService service = new CampaignEligibilityService();
    
    PromotionalCampaign campaign = createCampaignWithGeographicTargeting("CA", "NY", "FL");
    
    Customer californiaCustomer = Customer.builder()
        .geography("CA")
        .build();
    
    Customer texasCustomer = Customer.builder()
        .geography("TX")
        .build();
    
    Transaction transaction = createStandardTransaction();
    
    // Act
    EligibilityResult caResult = service.evaluateEligibility(californiaCustomer, transaction, campaign);
    EligibilityResult txResult = service.evaluateEligibility(texasCustomer, transaction, campaign);
    
    // Assert
    assertThat(caResult.isEligible()).isTrue();
    assertThat(txResult.isEligible()).isFalse();
    assertThat(txResult.getIneligibilityReason()).isEqualTo("Customer location not in target geography");
}
```

#### Test: Minimum purchase requirement validation
```java
@Test
public void testMinimumPurchaseRequirementValidation() {
    // Arrange
    CampaignEligibilityService service = new CampaignEligibilityService();
    
    DiscountConfiguration discountConfig = DiscountConfiguration.percentage(
        BigDecimal.valueOf(20.0),
        new Money(BigDecimal.valueOf(100.00), Currency.USD),  // $100 minimum
        null
    );
    
    PromotionalCampaign campaign = createCampaignWithDiscount(discountConfig);
    Customer customer = createEligibleCustomer();
    
    Transaction qualifyingTransaction = Transaction.builder()
        .amount(new Money(BigDecimal.valueOf(150.00), Currency.USD))
        .build();
    
    Transaction nonQualifyingTransaction = Transaction.builder()
        .amount(new Money(BigDecimal.valueOf(75.00), Currency.USD))
        .build();
    
    // Act
    EligibilityResult qualifyingResult = service.evaluateEligibility(customer, qualifyingTransaction, campaign);
    EligibilityResult nonQualifyingResult = service.evaluateEligibility(customer, nonQualifyingTransaction, campaign);
    
    // Assert
    assertThat(qualifyingResult.isEligible()).isTrue();
    assertThat(nonQualifyingResult.isEligible()).isFalse();
    assertThat(nonQualifyingResult.getIneligibilityReason()).isEqualTo("Transaction amount below minimum purchase requirement");
}
```

### UT-005: PromotionCalculationService Domain Service Unit Tests

#### Test: Calculate percentage discount with maximum limit
```java
@Test
public void testCalculatePercentageDiscountWithMaximumLimit() {
    // Arrange
    PromotionCalculationService service = new PromotionCalculationService();
    
    DiscountConfiguration config = DiscountConfiguration.percentage(
        BigDecimal.valueOf(25.0),  // 25% discount
        Money.zero(Currency.USD),  // no minimum
        new Money(BigDecimal.valueOf(50.00), Currency.USD)  // $50 maximum
    );
    
    Transaction largeTransaction = Transaction.builder()
        .amount(new Money(BigDecimal.valueOf(300.00), Currency.USD))
        .build();
    
    Transaction smallTransaction = Transaction.builder()
        .amount(new Money(BigDecimal.valueOf(100.00), Currency.USD))
        .build();
    
    // Act
    DiscountAmount largeDiscount = service.calculateDiscount(largeTransaction, config);
    DiscountAmount smallDiscount = service.calculateDiscount(smallTransaction, config);
    
    // Assert
    assertThat(largeDiscount.getAmount()).isEqualTo(new Money(BigDecimal.valueOf(50.00), Currency.USD));  // Capped at max
    assertThat(largeDiscount.isMaximumApplied()).isTrue();
    
    assertThat(smallDiscount.getAmount()).isEqualTo(new Money(BigDecimal.valueOf(25.00), Currency.USD));  // 25% of $100
    assertThat(smallDiscount.isMaximumApplied()).isFalse();
}
```

#### Test: Calculate buy-X-get-Y discount with multiple sets
```java
@Test
public void testCalculateBuyXGetYDiscountWithMultipleSets() {
    // Arrange
    PromotionCalculationService service = new PromotionCalculationService();
    
    DiscountConfiguration config = DiscountConfiguration.buyXGetY(
        2,  // buy 2
        1,  // get 1
        BigDecimal.valueOf(100.0),  // 100% off (free)
        5   // max 5 free items
    );
    
    List<CartItem> cartItems = Arrays.asList(
        new CartItem("ITEM-1", new Money(BigDecimal.valueOf(30.00), Currency.USD)),
        new CartItem("ITEM-2", new Money(BigDecimal.valueOf(25.00), Currency.USD)),
        new CartItem("ITEM-3", new Money(BigDecimal.valueOf(20.00), Currency.USD)),
        new CartItem("ITEM-4", new Money(BigDecimal.valueOf(15.00), Currency.USD)),
        new CartItem("ITEM-5", new Money(BigDecimal.valueOf(10.00), Currency.USD))
    );
    
    // Act
    BuyXGetYResult result = service.calculateBuyXGetY(cartItems, config);
    
    // Assert
    assertThat(result.getQualifyingSets()).isEqualTo(2);  // 5 items = 2 complete sets (4 items) + 1 remainder
    assertThat(result.getFreeItems()).hasSize(2);
    assertThat(result.getFreeItems()).containsExactly(
        cartItems.get(2),  // $20 item (lowest in first set)
        cartItems.get(4)   // $10 item (lowest in second set)
    );
    assertThat(result.getTotalDiscount()).isEqualTo(new Money(BigDecimal.valueOf(30.00), Currency.USD));
}
```

#### Test: Calculate bundle discount for complete bundle
```java
@Test
public void testCalculateBundleDiscountForCompleteBundle() {
    // Arrange
    PromotionCalculationService service = new PromotionCalculationService();
    
    BundleConfiguration bundleConfig = BundleConfiguration.builder()
        .addRequiredItem("LAPTOP-001", new Money(BigDecimal.valueOf(999.99), Currency.USD))
        .addRequiredItem("MOUSE-001", new Money(BigDecimal.valueOf(29.99), Currency.USD))
        .addRequiredItem("HEADPHONES-001", new Money(BigDecimal.valueOf(79.99), Currency.USD))
        .setBundlePrice(new Money(BigDecimal.valueOf(899.99), Currency.USD))
        .build();
    
    List<CartItem> cartItems = Arrays.asList(
        new CartItem("LAPTOP-001", new Money(BigDecimal.valueOf(999.99), Currency.USD)),
        new CartItem("MOUSE-001", new Money(BigDecimal.valueOf(29.99), Currency.USD)),
        new CartItem("HEADPHONES-001", new Money(BigDecimal.valueOf(79.99), Currency.USD)),
        new CartItem("KEYBOARD-001", new Money(BigDecimal.valueOf(49.99), Currency.USD))  // Not in bundle
    );
    
    // Act
    BundleDiscountResult result = service.calculateBundleDiscount(cartItems, bundleConfig);
    
    // Assert
    assertThat(result.isBundleComplete()).isTrue();
    assertThat(result.getIndividualTotal()).isEqualTo(new Money(BigDecimal.valueOf(1109.97), Currency.USD));
    assertThat(result.getBundlePrice()).isEqualTo(new Money(BigDecimal.valueOf(899.99), Currency.USD));
    assertThat(result.getBundleSavings()).isEqualTo(new Money(BigDecimal.valueOf(209.98), Currency.USD));
    assertThat(result.getNonBundleItems()).hasSize(1);
    assertThat(result.getNonBundleItems().get(0).getProductId()).isEqualTo("KEYBOARD-001");
}
```

#### Test: Calculate bundle discount for incomplete bundle
```java
@Test
public void testCalculateBundleDiscountForIncompleteBundle() {
    // Arrange
    PromotionCalculationService service = new PromotionCalculationService();
    
    BundleConfiguration bundleConfig = BundleConfiguration.builder()
        .addRequiredItem("LAPTOP-001", new Money(BigDecimal.valueOf(999.99), Currency.USD))
        .addRequiredItem("MOUSE-001", new Money(BigDecimal.valueOf(29.99), Currency.USD))
        .addRequiredItem("HEADPHONES-001", new Money(BigDecimal.valueOf(79.99), Currency.USD))
        .setBundlePrice(new Money(BigDecimal.valueOf(899.99), Currency.USD))
        .build();
    
    List<CartItem> cartItems = Arrays.asList(
        new CartItem("LAPTOP-001", new Money(BigDecimal.valueOf(999.99), Currency.USD)),
        new CartItem("MOUSE-001", new Money(BigDecimal.valueOf(29.99), Currency.USD))
        // Missing HEADPHONES-001
    );
    
    // Act
    BundleDiscountResult result = service.calculateBundleDiscount(cartItems, bundleConfig);
    
    // Assert
    assertThat(result.isBundleComplete()).isFalse();
    assertThat(result.getBundleSavings()).isEqualTo(Money.zero(Currency.USD));
    assertThat(result.getMissingItems()).hasSize(1);
    assertThat(result.getMissingItems().get(0)).isEqualTo("HEADPHONES-001");
}
```

### UT-006: ConflictResolutionService Domain Service Unit Tests

#### Test: Resolve conflict by selecting best customer benefit
```java
@Test
public void testResolveConflictBySelectingBestCustomerBenefit() {
    // Arrange
    ConflictResolutionService service = new ConflictResolutionService();
    
    List<PromotionalCampaign> conflictingCampaigns = Arrays.asList(
        createCampaignWithPercentageDiscount("PROMO-1", 15.0),  // 15% = $30 on $200
        createCampaignWithFixedDiscount("PROMO-2", 40.00),      // $40 fixed
        createCampaignWithPercentageDiscount("PROMO-3", 10.0)   // 10% = $20 on $200
    );
    
    Customer customer = createStandardCustomer();
    Transaction transaction = createTransactionWithAmount(200.00);
    
    // Act
    ConflictResolutionResult result = service.resolveBestForCustomer(
        conflictingCampaigns, customer, transaction
    );
    
    // Assert
    assertThat(result.getSelectedCampaign().getId().getValue()).isEqualTo("PROMO-2");
    assertThat(result.getResolutionStrategy()).isEqualTo(ResolutionStrategy.BEST_FOR_CUSTOMER);
    assertThat(result.getCustomerBenefit()).isEqualTo(new Money(BigDecimal.valueOf(40.00), Currency.USD));
    assertThat(result.getResolutionReason()).isEqualTo("Selected promotion with highest customer benefit");
}
```

#### Test: Resolve conflict by campaign priority
```java
@Test
public void testResolveConflictByCampaignPriority() {
    // Arrange
    ConflictResolutionService service = new ConflictResolutionService();
    
    List<PromotionalCampaign> conflictingCampaigns = Arrays.asList(
        createCampaignWithPriority("PROMO-HIGH", 1, 30.00),    // High priority, lower benefit
        createCampaignWithPriority("PROMO-LOW", 5, 50.00),     // Low priority, higher benefit
        createCampaignWithPriority("PROMO-MED", 3, 40.00)      // Medium priority, medium benefit
    );
    
    Customer customer = createStandardCustomer();
    Transaction transaction = createStandardTransaction();
    
    // Act
    ConflictResolutionResult result = service.applyPriorityBasedResolution(
        conflictingCampaigns, customer, transaction
    );
    
    // Assert
    assertThat(result.getSelectedCampaign().getId().getValue()).isEqualTo("PROMO-HIGH");
    assertThat(result.getResolutionStrategy()).isEqualTo(ResolutionStrategy.PRIORITY_BASED);
    assertThat(result.getResolutionReason()).isEqualTo("Selected promotion with highest priority");
}
```

#### Test: Calculate optimal combination of stackable promotions
```java
@Test
public void testCalculateOptimalCombinationOfStackablePromotions() {
    // Arrange
    ConflictResolutionService service = new ConflictResolutionService();
    
    List<PromotionalCampaign> stackablePromotions = Arrays.asList(
        createStackableCampaign("CATEGORY-BONUS", 10.0, true),   // 10% category bonus
        createStackableCampaign("LOYALTY-TIER", 5.0, true),      // 5% loyalty tier bonus
        createStackableCampaign("FREE-SHIPPING", 15.00, true)    // $15 free shipping
    );
    
    Customer goldTierCustomer = createCustomerWithTier(LoyaltyTier.GOLD);
    Transaction electronicsTransaction = createTransactionWithCategory("ELECTRONICS", 200.00);
    
    // Act
    StackingResult result = service.calculateOptimalCombination(
        stackablePromotions, goldTierCustomer, electronicsTransaction
    );
    
    // Assert
    assertThat(result.getAppliedPromotions()).hasSize(3);
    assertThat(result.getTotalBenefit()).isEqualTo(new Money(BigDecimal.valueOf(46.00), Currency.USD));
    // $200 * 10% = $20 (category), $200 * 5% = $10 (loyalty), $15 (shipping) = $45 total
    assertThat(result.getStackingOrder()).containsExactly("FREE-SHIPPING", "CATEGORY-BONUS", "LOYALTY-TIER");
}
```

### UT-007: Domain Events Unit Tests

#### Test: PromotionalCampaignCreated domain event
```java
@Test
public void testPromotionalCampaignCreatedDomainEvent() {
    // Arrange
    CampaignId campaignId = new CampaignId("PROMO-TEST-001");
    CampaignName campaignName = new CampaignName("Test Campaign");
    CampaignType campaignType = CampaignType.PERCENTAGE_DISCOUNT;
    String createdBy = "marketing.manager@company.com";
    LocalDateTime createdAt = LocalDateTime.now();
    
    // Act
    PromotionalCampaignCreated event = new PromotionalCampaignCreated(
        campaignId, campaignName, campaignType, createdBy, createdAt
    );
    
    // Assert
    assertThat(event.getCampaignId()).isEqualTo(campaignId);
    assertThat(event.getCampaignName()).isEqualTo(campaignName);
    assertThat(event.getCampaignType()).isEqualTo(campaignType);
    assertThat(event.getCreatedBy()).isEqualTo(createdBy);
    assertThat(event.getCreatedAt()).isEqualTo(createdAt);
    assertThat(event.getEventType()).isEqualTo("promotion.campaign.created");
    assertThat(event.getOccurredAt()).isNotNull();
}
```

#### Test: PromotionApplied domain event
```java
@Test
public void testPromotionAppliedDomainEvent() {
    // Arrange
    ApplicationId applicationId = new ApplicationId("APP-001");
    CampaignId campaignId = new CampaignId("PROMO-001");
    CustomerId customerId = new CustomerId("CUST-123");
    TransactionId transactionId = new TransactionId("TXN-456");
    Money appliedBenefit = new Money(BigDecimal.valueOf(25.00), Currency.USD);
    LocalDateTime appliedAt = LocalDateTime.now();
    
    // Act
    PromotionApplied event = new PromotionApplied(
        applicationId, campaignId, customerId, transactionId, appliedBenefit, appliedAt
    );
    
    // Assert
    assertThat(event.getApplicationId()).isEqualTo(applicationId);
    assertThat(event.getCampaignId()).isEqualTo(campaignId);
    assertThat(event.getCustomerId()).isEqualTo(customerId);
    assertThat(event.getTransactionId()).isEqualTo(transactionId);
    assertThat(event.getAppliedBenefit()).isEqualTo(appliedBenefit);
    assertThat(event.getAppliedAt()).isEqualTo(appliedAt);
    assertThat(event.getEventType()).isEqualTo("promotion.applied");
}
```

#### Test: CampaignBudgetExhausted domain event
```java
@Test
public void testCampaignBudgetExhaustedDomainEvent() {
    // Arrange
    CampaignId campaignId = new CampaignId("PROMO-BUDGET-001");
    Money totalBudget = new Money(BigDecimal.valueOf(5000.00), Currency.USD);
    Money totalSpent = new Money(BigDecimal.valueOf(5000.00), Currency.USD);
    LocalDateTime exhaustedAt = LocalDateTime.now();
    ApplicationId lastApplicationId = new ApplicationId("APP-FINAL-001");
    
    // Act
    CampaignBudgetExhausted event = new CampaignBudgetExhausted(
        campaignId, totalBudget, totalSpent, exhaustedAt, lastApplicationId
    );
    
    // Assert
    assertThat(event.getCampaignId()).isEqualTo(campaignId);
    assertThat(event.getTotalBudget()).isEqualTo(totalBudget);
    assertThat(event.getTotalSpent()).isEqualTo(totalSpent);
    assertThat(event.getExhaustedAt()).isEqualTo(exhaustedAt);
    assertThat(event.getLastApplicationId()).isEqualTo(lastApplicationId);
    assertThat(event.getEventType()).isEqualTo("promotion.budget.exhausted");
}
```

## Test Utilities and Builders

### Test Data Builders
```java
public class PromotionTestDataBuilder {
    
    public static PromotionalCampaign createTestCampaign() {
        return PromotionalCampaign.create(
            new CampaignId("TEST-CAMPAIGN-001"),
            new CampaignName("Test Campaign"),
            CampaignType.PERCENTAGE_DISCOUNT,
            createStandardTargeting(),
            createStandardValidity(),
            createStandardBudget()
        );
    }
    
    public static PromotionalCampaign createActiveCampaign() {
        PromotionalCampaign campaign = createTestCampaign();
        campaign.approve(new ApprovalDecision(true, "Auto-approved for testing"));
        campaign.activate(LocalDateTime.now());
        campaign.markEventsAsCommitted();
        return campaign;
    }
    
    public static TargetingCriteria createStandardTargeting() {
        return new TargetingCriteria(
            CustomerSegments.of(CustomerSegment.ALL_CUSTOMERS),
            ProductCategories.all()
        );
    }
    
    public static ValidityPeriod createStandardValidity() {
        return new ValidityPeriod(
            LocalDateTime.now(),
            LocalDateTime.now().plusDays(30)
        );
    }
    
    public static Budget createStandardBudget() {
        return new Budget(new Money(BigDecimal.valueOf(10000.00), Currency.USD));
    }
    
    public static Customer createEligibleCustomer() {
        return Customer.builder()
            .customerId(new CustomerId("TEST-CUSTOMER-001"))
            .segment(CustomerSegment.LOYAL_CUSTOMERS)
            .loyaltyTier(LoyaltyTier.GOLD)
            .geography("CA")
            .build();
    }
    
    public static Transaction createStandardTransaction() {
        return Transaction.builder()
            .transactionId(new TransactionId("TEST-TXN-001"))
            .amount(new Money(BigDecimal.valueOf(200.00), Currency.USD))
            .items(createStandardCartItems())
            .build();
    }
    
    public static List<CartItem> createStandardCartItems() {
        return Arrays.asList(
            new CartItem("ITEM-1", new Money(BigDecimal.valueOf(120.00), Currency.USD)),
            new CartItem("ITEM-2", new Money(BigDecimal.valueOf(80.00), Currency.USD))
        );
    }
}
```

### Mock Repositories and Services
```java
@ExtendWith(MockitoExtension.class)
class PromotionServiceTest {
    
    @Mock
    private PromotionalCampaignRepository campaignRepository;
    
    @Mock
    private CustomerRepository customerRepository;
    
    @Mock
    private ProductCatalogService productCatalogService;
    
    @Mock
    private BudgetTrackingService budgetTrackingService;
    
    @InjectMocks
    private PromotionService promotionService;
    
    @Test
    public void testServiceCoordinationMethods() {
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
- All discount calculation algorithms must have 100% test coverage
- All customer targeting and eligibility rules must be fully tested
- All conflict resolution strategies must have complete test coverage
- All budget tracking and validation logic must be tested
- All domain event publishing scenarios must be covered

### Edge Case Coverage
- Invalid discount configurations and boundary conditions
- Budget exhaustion and overflow scenarios
- Complex promotion stacking and conflict resolution
- Invalid customer targeting and geography edge cases
- Error conditions and exception handling paths

### Performance Test Targets
- Unit tests should complete within 10ms each
- Test suite should complete within 2 minutes
- No external dependencies in unit tests
- All tests should be deterministic and repeatable
- Comprehensive test data builders for complex scenarios
