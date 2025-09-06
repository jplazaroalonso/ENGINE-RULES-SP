# Unit Tests - Rule Evaluation Engine

## Test Suite: Rule Evaluation Core Components

### UT-001: EvaluationContext Unit Tests

#### Test: Create evaluation context with transaction data
```java
@Test
public void testCreateEvaluationContext() {
    // Arrange
    TransactionData transaction = TransactionData.builder()
        .customerId("CUST-123")
        .amount(new BigDecimal("100.00"))
        .currency("USD")
        .items(Arrays.asList(
            new Item("ITEM-1", new BigDecimal("60.00")),
            new Item("ITEM-2", new BigDecimal("40.00"))
        ))
        .build();
    
    // Act
    EvaluationContext context = EvaluationContext.create(transaction);
    
    // Assert
    assertThat(context.getCustomerId()).isEqualTo("CUST-123");
    assertThat(context.getTransactionAmount()).isEqualTo(new BigDecimal("100.00"));
    assertThat(context.getItemCount()).isEqualTo(2);
    assertThat(context.isValid()).isTrue();
}
```

#### Test: Add applicable rules to context
```java
@Test
public void testAddApplicableRules() {
    // Arrange
    EvaluationContext context = createTestContext();
    Rule rule1 = createTestRule("RULE-001", Priority.HIGH);
    Rule rule2 = createTestRule("RULE-002", Priority.MEDIUM);
    
    // Act
    context.addApplicableRule(rule1);
    context.addApplicableRule(rule2);
    
    // Assert
    assertThat(context.getApplicableRules()).hasSize(2);
    assertThat(context.getApplicableRules().get(0)).isEqualTo(rule1); // Higher priority first
    assertThat(context.getApplicableRules().get(1)).isEqualTo(rule2);
}
```

#### Test: Validate context completeness
```java
@Test
public void testValidateContextCompleteness() {
    // Arrange
    EvaluationContext incompleteContext = EvaluationContext.builder()
        .customerId("CUST-123")
        // Missing required fields
        .build();
    
    // Act
    ValidationResult result = incompleteContext.validate();
    
    // Assert
    assertThat(result.isValid()).isFalse();
    assertThat(result.getErrors()).contains("Transaction amount is required");
    assertThat(result.getErrors()).contains("Currency is required");
}
```

---

### UT-002: RuleEvaluator Unit Tests

#### Test: Evaluate single rule successfully
```java
@Test
public void testEvaluateSingleRule() {
    // Arrange
    RuleEvaluator evaluator = new RuleEvaluator();
    Rule rule = Rule.builder()
        .id("RULE-001")
        .condition("customer.tier = 'GOLD' AND purchase.amount > 50")
        .action("discount := 10%")
        .build();
    
    EvaluationContext context = EvaluationContext.builder()
        .customerId("CUST-123")
        .customerTier("GOLD")
        .transactionAmount(new BigDecimal("75.00"))
        .build();
    
    // Act
    RuleEvaluationResult result = evaluator.evaluate(rule, context);
    
    // Assert
    assertThat(result.isSuccessful()).isTrue();
    assertThat(result.getAppliedActions()).hasSize(1);
    assertThat(result.getAppliedActions().get(0).getType()).isEqualTo(ActionType.DISCOUNT);
    assertThat(result.getAppliedActions().get(0).getValue()).isEqualTo(new BigDecimal("0.10"));
}
```

#### Test: Rule condition not met
```java
@Test
public void testRuleConditionNotMet() {
    // Arrange
    RuleEvaluator evaluator = new RuleEvaluator();
    Rule rule = Rule.builder()
        .condition("customer.tier = 'GOLD'")
        .action("discount := 10%")
        .build();
    
    EvaluationContext context = EvaluationContext.builder()
        .customerTier("SILVER") // Does not match condition
        .transactionAmount(new BigDecimal("75.00"))
        .build();
    
    // Act
    RuleEvaluationResult result = evaluator.evaluate(rule, context);
    
    // Assert
    assertThat(result.isSuccessful()).isTrue();
    assertThat(result.wasApplied()).isFalse();
    assertThat(result.getAppliedActions()).isEmpty();
    assertThat(result.getSkipReason()).isEqualTo("Condition not met");
}
```

#### Test: Handle rule evaluation error
```java
@Test
public void testHandleRuleEvaluationError() {
    // Arrange
    RuleEvaluator evaluator = new RuleEvaluator();
    Rule invalidRule = Rule.builder()
        .condition("invalid.syntax.here") // Invalid DSL
        .action("discount := 10%")
        .build();
    
    EvaluationContext context = createTestContext();
    
    // Act
    RuleEvaluationResult result = evaluator.evaluate(invalidRule, context);
    
    // Assert
    assertThat(result.isSuccessful()).isFalse();
    assertThat(result.getError()).isNotNull();
    assertThat(result.getError().getType()).isEqualTo(ErrorType.SYNTAX_ERROR);
    assertThat(result.getError().getMessage()).contains("invalid syntax");
}
```

---

### UT-003: ConflictResolver Unit Tests

#### Test: Resolve priority-based conflicts
```java
@Test
public void testResolvePriorityBasedConflicts() {
    // Arrange
    ConflictResolver resolver = new ConflictResolver();
    
    RuleEvaluationResult highPriorityResult = RuleEvaluationResult.builder()
        .rule(createRuleWithPriority("RULE-HIGH", Priority.HIGH))
        .appliedActions(Arrays.asList(new DiscountAction(new BigDecimal("0.15"))))
        .build();
    
    RuleEvaluationResult mediumPriorityResult = RuleEvaluationResult.builder()
        .rule(createRuleWithPriority("RULE-MED", Priority.MEDIUM))
        .appliedActions(Arrays.asList(new DiscountAction(new BigDecimal("0.10"))))
        .build();
    
    List<RuleEvaluationResult> conflictingResults = Arrays.asList(
        highPriorityResult, mediumPriorityResult
    );
    
    // Act
    ConflictResolution resolution = resolver.resolveConflicts(conflictingResults);
    
    // Assert
    assertThat(resolution.getWinningResult()).isEqualTo(highPriorityResult);
    assertThat(resolution.getRejectedResults()).containsExactly(mediumPriorityResult);
    assertThat(resolution.getResolutionStrategy()).isEqualTo(ResolutionStrategy.PRIORITY_BASED);
}
```

#### Test: Detect mutually exclusive conflicts
```java
@Test
public void testDetectMutuallyExclusiveConflicts() {
    // Arrange
    ConflictResolver resolver = new ConflictResolver();
    
    RuleEvaluationResult freeShippingResult = createResultWithAction(
        new FreeShippingAction()
    );
    RuleEvaluationResult expressShippingResult = createResultWithAction(
        new ExpressShippingAction()
    );
    
    List<RuleEvaluationResult> results = Arrays.asList(
        freeShippingResult, expressShippingResult
    );
    
    // Act
    ConflictDetection detection = resolver.detectConflicts(results);
    
    // Assert
    assertThat(detection.hasConflicts()).isTrue();
    assertThat(detection.getConflictType()).isEqualTo(ConflictType.MUTUALLY_EXCLUSIVE);
    assertThat(detection.getConflictingResults()).containsExactlyInAnyOrder(
        freeShippingResult, expressShippingResult
    );
}
```

#### Test: Resolve stackable promotions
```java
@Test
public void testResolveStackablePromotions() {
    // Arrange
    ConflictResolver resolver = new ConflictResolver();
    
    RuleEvaluationResult loyaltyDiscount = createResultWithStackableAction(
        new DiscountAction(new BigDecimal("0.05")), true
    );
    RuleEvaluationResult couponDiscount = createResultWithStackableAction(
        new DiscountAction(new BigDecimal("0.10")), true
    );
    
    List<RuleEvaluationResult> results = Arrays.asList(
        loyaltyDiscount, couponDiscount
    );
    
    // Act
    ConflictResolution resolution = resolver.resolveConflicts(results);
    
    // Assert
    assertThat(resolution.isStackingAllowed()).isTrue();
    assertThat(resolution.getStackedResults()).containsExactlyInAnyOrder(
        loyaltyDiscount, couponDiscount
    );
    assertThat(resolution.getCombinedDiscount()).isEqualTo(new BigDecimal("0.15"));
}
```

---

### UT-004: PerformanceMonitor Unit Tests

#### Test: Track evaluation performance
```java
@Test
public void testTrackEvaluationPerformance() {
    // Arrange
    PerformanceMonitor monitor = new PerformanceMonitor();
    String evaluationId = "EVAL-123";
    
    // Act
    monitor.startEvaluation(evaluationId);
    
    // Simulate some processing time
    try { Thread.sleep(50); } catch (InterruptedException e) {}
    
    monitor.endEvaluation(evaluationId);
    
    // Assert
    PerformanceMetrics metrics = monitor.getMetrics(evaluationId);
    assertThat(metrics.getExecutionTime()).isGreaterThan(Duration.ofMillis(45));
    assertThat(metrics.getExecutionTime()).isLessThan(Duration.ofMillis(100));
    assertThat(metrics.isWithinSLA()).isTrue(); // <500ms SLA
}
```

#### Test: Monitor cache performance
```java
@Test
public void testMonitorCachePerformance() {
    // Arrange
    PerformanceMonitor monitor = new PerformanceMonitor();
    
    // Act
    monitor.recordCacheHit("RULE-001");
    monitor.recordCacheHit("RULE-002");
    monitor.recordCacheMiss("RULE-003");
    
    // Assert
    CacheMetrics metrics = monitor.getCacheMetrics();
    assertThat(metrics.getHitRate()).isEqualTo(new BigDecimal("0.67")); // 2/3
    assertThat(metrics.getTotalRequests()).isEqualTo(3);
    assertThat(metrics.getHits()).isEqualTo(2);
    assertThat(metrics.getMisses()).isEqualTo(1);
}
```

#### Test: Alert on SLA violations
```java
@Test
public void testAlertOnSLAViolations() {
    // Arrange
    PerformanceMonitor monitor = new PerformanceMonitor();
    String evaluationId = "EVAL-SLOW";
    
    // Act
    monitor.startEvaluation(evaluationId);
    
    // Simulate slow processing (>500ms SLA)
    try { Thread.sleep(600); } catch (InterruptedException e) {}
    
    monitor.endEvaluation(evaluationId);
    
    // Assert
    PerformanceMetrics metrics = monitor.getMetrics(evaluationId);
    assertThat(metrics.isWithinSLA()).isFalse();
    
    List<Alert> alerts = monitor.getAlerts();
    assertThat(alerts).hasSize(1);
    assertThat(alerts.get(0).getType()).isEqualTo(AlertType.SLA_VIOLATION);
    assertThat(alerts.get(0).getMessage()).contains("Evaluation exceeded 500ms SLA");
}
```

---

### UT-005: EvaluationResult Unit Tests

#### Test: Aggregate multiple rule results
```java
@Test
public void testAggregateMultipleRuleResults() {
    // Arrange
    List<RuleEvaluationResult> results = Arrays.asList(
        createDiscountResult("RULE-001", new BigDecimal("10.00")),
        createShippingResult("RULE-002", ShippingType.FREE),
        createLoyaltyResult("RULE-003", 100)
    );
    
    // Act
    AggregatedEvaluationResult aggregated = AggregatedEvaluationResult.from(results);
    
    // Assert
    assertThat(aggregated.getAppliedRulesCount()).isEqualTo(3);
    assertThat(aggregated.getTotalDiscount()).isEqualTo(new BigDecimal("10.00"));
    assertThat(aggregated.getShippingType()).isEqualTo(ShippingType.FREE);
    assertThat(aggregated.getLoyaltyPointsEarned()).isEqualTo(100);
}
```

#### Test: Calculate final transaction amount
```java
@Test
public void testCalculateFinalTransactionAmount() {
    // Arrange
    EvaluationResult result = EvaluationResult.builder()
        .originalAmount(new BigDecimal("100.00"))
        .discountAmount(new BigDecimal("15.00"))
        .shippingCost(new BigDecimal("5.00"))
        .taxAmount(new BigDecimal("8.50"))
        .build();
    
    // Act
    BigDecimal finalAmount = result.calculateFinalAmount();
    
    // Assert
    // 100.00 - 15.00 + 5.00 + 8.50 = 98.50
    assertThat(finalAmount).isEqualTo(new BigDecimal("98.50"));
}
```

#### Test: Generate evaluation summary
```java
@Test
public void testGenerateEvaluationSummary() {
    // Arrange
    EvaluationResult result = EvaluationResult.builder()
        .evaluationId("EVAL-123")
        .customerId("CUST-456")
        .appliedRules(Arrays.asList("RULE-001", "RULE-002"))
        .executionTime(Duration.ofMillis(250))
        .cacheHitRate(new BigDecimal("0.80"))
        .build();
    
    // Act
    EvaluationSummary summary = result.generateSummary();
    
    // Assert
    assertThat(summary.getEvaluationId()).isEqualTo("EVAL-123");
    assertThat(summary.getCustomerId()).isEqualTo("CUST-456");
    assertThat(summary.getRulesAppliedCount()).isEqualTo(2);
    assertThat(summary.getExecutionTime()).isEqualTo(Duration.ofMillis(250));
    assertThat(summary.isWithinPerformanceSLA()).isTrue();
    assertThat(summary.getCacheEfficiency()).isEqualTo("Good"); // >75%
}
```

---

### UT-006: Error Handling Unit Tests

#### Test: Handle DSL syntax errors
```java
@Test
public void testHandleDSLSyntaxErrors() {
    // Arrange
    DSLParser parser = new DSLParser();
    String invalidDSL = "IF customer.tier = GOLD THEN"; // Missing action
    
    // Act
    DSLParseResult result = parser.parse(invalidDSL);
    
    // Assert
    assertThat(result.isSuccessful()).isFalse();
    assertThat(result.getError().getType()).isEqualTo(ErrorType.SYNTAX_ERROR);
    assertThat(result.getError().getLineNumber()).isEqualTo(1);
    assertThat(result.getError().getMessage()).contains("Missing action after THEN");
}
```

#### Test: Handle runtime evaluation errors
```java
@Test
public void testHandleRuntimeEvaluationErrors() {
    // Arrange
    RuleEvaluator evaluator = new RuleEvaluator();
    Rule rule = Rule.builder()
        .condition("customer.tier = 'GOLD'")
        .action("discount := purchase.amount / 0") // Division by zero
        .build();
    
    EvaluationContext context = createTestContext();
    
    // Act
    RuleEvaluationResult result = evaluator.evaluate(rule, context);
    
    // Assert
    assertThat(result.isSuccessful()).isFalse();
    assertThat(result.getError().getType()).isEqualTo(ErrorType.RUNTIME_ERROR);
    assertThat(result.getError().getMessage()).contains("Division by zero");
    assertThat(result.getFallbackValue()).isNotNull();
}
```

#### Test: Graceful degradation on service failures
```java
@Test
public void testGracefulDegradationOnServiceFailures() {
    // Arrange
    RuleEvaluationService service = new RuleEvaluationService();
    service.setCustomerService(createFailingCustomerService());
    service.setFallbackMode(true);
    
    EvaluationRequest request = createTestEvaluationRequest();
    
    // Act
    EvaluationResult result = service.evaluate(request);
    
    // Assert
    assertThat(result.isSuccessful()).isTrue();
    assertThat(result.isDegradedMode()).isTrue();
    assertThat(result.getWarnings()).contains("Customer service unavailable, using cached data");
}
```

---

## Test Utilities and Builders

### Test Data Builders
```java
public class EvaluationTestDataBuilder {
    
    public static EvaluationContext createTestContext() {
        return EvaluationContext.builder()
            .customerId("CUST-123")
            .customerTier("GOLD")
            .transactionAmount(new BigDecimal("100.00"))
            .currency("USD")
            .channel(Channel.ONLINE)
            .build();
    }
    
    public static Rule createTestRule(String ruleId, Priority priority) {
        return Rule.builder()
            .id(ruleId)
            .priority(priority)
            .condition("customer.tier = 'GOLD'")
            .action("discount := 10%")
            .isActive(true)
            .build();
    }
    
    public static RuleEvaluationResult createDiscountResult(String ruleId, BigDecimal amount) {
        return RuleEvaluationResult.builder()
            .ruleId(ruleId)
            .wasApplied(true)
            .appliedActions(Arrays.asList(new DiscountAction(amount)))
            .executionTime(Duration.ofMillis(50))
            .build();
    }
}
```

### Mock Services
```java
@ExtendWith(MockitoExtension.class)
class RuleEvaluationServiceTest {
    
    @Mock
    private CustomerService customerService;
    
    @Mock
    private RuleRepository ruleRepository;
    
    @Mock
    private PerformanceMonitor performanceMonitor;
    
    @InjectMocks
    private RuleEvaluationService evaluationService;
    
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
- All evaluation paths must have 100% test coverage
- All error handling scenarios must be tested
- All performance-critical code must be covered
- All conflict resolution strategies must be tested

### Performance Test Targets
- Unit tests should complete within 100ms each
- Test suite should complete within 10 minutes
- No external dependencies in unit tests
- All tests should be deterministic and repeatable
