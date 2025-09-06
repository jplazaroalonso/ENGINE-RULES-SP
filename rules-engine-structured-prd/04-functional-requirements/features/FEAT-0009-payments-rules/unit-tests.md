# Unit Tests - Payment Rules and Processing

## Test Suite: Payment Domain Components

### UT-001: PaymentTransaction Aggregate Unit Tests

#### Test: Create payment transaction with valid details
```java
@Test
public void testCreatePaymentTransaction() {
    // Arrange
    TransactionId transactionId = new TransactionId("TXN-PAY-001");
    PaymentAmount amount = new PaymentAmount(
        new Money(BigDecimal.valueOf(250.00), Currency.USD)
    );
    CustomerId customerId = new CustomerId("CUST-123");
    PaymentMethod paymentMethod = PaymentMethod.creditCard(
        "4111111111111111", "12/25", "123", "John Doe"
    );
    
    // Act
    PaymentTransaction transaction = PaymentTransaction.create(
        transactionId, amount, customerId, paymentMethod
    );
    
    // Assert
    assertThat(transaction.getId()).isEqualTo(transactionId);
    assertThat(transaction.getAmount()).isEqualTo(amount);
    assertThat(transaction.getCustomerId()).isEqualTo(customerId);
    assertThat(transaction.getStatus()).isEqualTo(PaymentStatus.INITIATED);
    
    // Verify domain event
    List<DomainEvent> events = transaction.getUncommittedEvents();
    assertThat(events).hasSize(1);
    assertThat(events.get(0)).isInstanceOf(PaymentTransactionInitiated.class);
}
```

#### Test: Select optimal payment gateway
```java
@Test
public void testSelectOptimalPaymentGateway() {
    // Arrange
    PaymentTransaction transaction = createTestPaymentTransaction();
    List<PaymentGateway> availableGateways = Arrays.asList(
        createGateway("GATEWAY_A", 97.0, BigDecimal.valueOf(2.9)),
        createGateway("GATEWAY_B", 95.0, BigDecimal.valueOf(2.1)),
        createGateway("GATEWAY_C", 99.0, BigDecimal.valueOf(3.2))
    );
    
    // Act
    GatewaySelection selection = transaction.selectOptimalGateway(availableGateways);
    
    // Assert
    assertThat(selection.getSelectedGateway().getId()).isEqualTo("GATEWAY_C");
    assertThat(selection.getSelectionReason()).contains("highest success rate");
    assertThat(transaction.getGateway()).isEqualTo(selection.getSelectedGateway());
    
    // Verify domain event
    List<DomainEvent> events = transaction.getUncommittedEvents();
    assertThat(events.stream()
        .anyMatch(e -> e instanceof PaymentGatewaySelected))
        .isTrue();
}
```

#### Test: Assess fraud risk for transaction
```java
@Test
public void testAssessFraudRisk() {
    // Arrange
    PaymentTransaction transaction = createTestPaymentTransaction();
    CustomerProfile customerProfile = createTrustedCustomerProfile();
    TransactionContext context = createNormalTransactionContext();
    
    // Act
    FraudRiskScore riskScore = transaction.assessFraudRisk(customerProfile, context);
    
    // Assert
    assertThat(riskScore.getLevel()).isEqualTo(RiskLevel.VERY_LOW);
    assertThat(riskScore.getScore()).isBetween(0, 20);
    assertThat(riskScore.getFactors()).doesNotContain(RiskFactor.AMOUNT_ANOMALY);
    assertThat(transaction.getRiskScore()).isEqualTo(riskScore);
    
    // Verify domain event
    List<DomainEvent> events = transaction.getUncommittedEvents();
    assertThat(events.stream()
        .anyMatch(e -> e instanceof FraudRiskAssessed))
        .isTrue();
}
```

#### Test: Process payment successfully
```java
@Test
public void testProcessPaymentSuccessfully() {
    // Arrange
    PaymentTransaction transaction = createReadyTransaction();
    AuthorizationCode authCode = new AuthorizationCode("AUTH-12345");
    LocalDateTime completedAt = LocalDateTime.now();
    
    // Act
    transaction.completePayment(authCode, completedAt);
    
    // Assert
    assertThat(transaction.getStatus()).isEqualTo(PaymentStatus.COMPLETED);
    assertThat(transaction.getAuthorizationCode()).isEqualTo(authCode);
    assertThat(transaction.getCompletedAt()).isEqualTo(completedAt);
    
    // Verify domain event
    List<DomainEvent> events = transaction.getUncommittedEvents();
    assertThat(events.stream()
        .anyMatch(e -> e instanceof PaymentProcessingCompleted))
        .isTrue();
}
```

#### Test: Handle payment failure with retry
```java
@Test
public void testHandlePaymentFailureWithRetry() {
    // Arrange
    PaymentTransaction transaction = createProcessingTransaction();
    PaymentFailureReason failureReason = PaymentFailureReason.TEMPORARY_GATEWAY_ERROR;
    LocalDateTime failedAt = LocalDateTime.now();
    
    // Act
    transaction.handlePaymentFailure(failureReason, failedAt);
    
    // Assert
    assertThat(transaction.getStatus()).isEqualTo(PaymentStatus.RETRY_PENDING);
    assertThat(transaction.getFailureReason()).isEqualTo(failureReason);
    assertThat(transaction.canRetry()).isTrue();
    
    // Verify domain event
    List<DomainEvent> events = transaction.getUncommittedEvents();
    assertThat(events.stream()
        .anyMatch(e -> e instanceof PaymentFailureDetected))
        .isTrue();
}
```

### UT-002: PaymentAmount Value Object Unit Tests

#### Test: Create payment amount with currency conversion
```java
@Test
public void testCreatePaymentAmountWithCurrencyConversion() {
    // Arrange
    Money baseAmount = new Money(BigDecimal.valueOf(100.00), Currency.EUR);
    Currency targetCurrency = Currency.USD;
    ExchangeRate exchangeRate = new ExchangeRate(
        Currency.EUR, Currency.USD, BigDecimal.valueOf(1.0850)
    );
    
    // Act
    PaymentAmount paymentAmount = PaymentAmount.createWithConversion(
        baseAmount, targetCurrency, exchangeRate
    );
    
    // Assert
    assertThat(paymentAmount.getBaseAmount()).isEqualTo(baseAmount);
    assertThat(paymentAmount.getTotalAmount().getCurrency()).isEqualTo(Currency.USD);
    assertThat(paymentAmount.getTotalAmount().getAmount())
        .isEqualTo(BigDecimal.valueOf(108.50));
    assertThat(paymentAmount.getExchangeRate()).isEqualTo(exchangeRate);
}
```

#### Test: Calculate processing fees
```java
@Test
public void testCalculateProcessingFees() {
    // Arrange
    Money baseAmount = new Money(BigDecimal.valueOf(1000.00), Currency.USD);
    PaymentAmount paymentAmount = new PaymentAmount(baseAmount);
    FeeRules feeRules = new FeeRules(
        BigDecimal.valueOf(2.9),  // 2.9% processing fee
        new Money(BigDecimal.valueOf(0.30), Currency.USD)  // $0.30 fixed fee
    );
    
    // Act
    paymentAmount.calculateFees(feeRules);
    
    // Assert
    assertThat(paymentAmount.getProcessingFees())
        .isEqualTo(new Money(BigDecimal.valueOf(29.30), Currency.USD));  // $29.00 + $0.30
    assertThat(paymentAmount.getTotalAmount())
        .isEqualTo(new Money(BigDecimal.valueOf(1029.30), Currency.USD));
}
```

#### Test: Validate amount limits
```java
@Test
public void testValidateAmountLimits() {
    // Arrange
    PaymentAmount smallAmount = new PaymentAmount(
        new Money(BigDecimal.valueOf(5.00), Currency.USD)
    );
    PaymentAmount largeAmount = new PaymentAmount(
        new Money(BigDecimal.valueOf(15000.00), Currency.USD)
    );
    AmountLimits limits = new AmountLimits(
        new Money(BigDecimal.valueOf(10.00), Currency.USD),    // minimum
        new Money(BigDecimal.valueOf(10000.00), Currency.USD)  // maximum
    );
    
    // Act & Assert
    assertThrows(AmountBelowMinimumException.class, () -> {
        smallAmount.validateAmount(limits);
    });
    
    assertThrows(AmountExceedsMaximumException.class, () -> {
        largeAmount.validateAmount(limits);
    });
    
    // Valid amount should not throw
    PaymentAmount validAmount = new PaymentAmount(
        new Money(BigDecimal.valueOf(500.00), Currency.USD)
    );
    assertDoesNotThrow(() -> validAmount.validateAmount(limits));
}
```

### UT-003: PaymentMethod Value Object Unit Tests

#### Test: Create credit card payment method
```java
@Test
public void testCreateCreditCardPaymentMethod() {
    // Arrange & Act
    PaymentMethod creditCard = PaymentMethod.creditCard(
        "4111111111111111",
        "12/25",
        "123",
        "John Doe"
    );
    
    // Assert
    assertThat(creditCard.getType()).isEqualTo(PaymentMethodType.CREDIT_CARD);
    assertThat(creditCard.getInstrument().getMaskedNumber()).isEqualTo("****-****-****-1111");
    assertThat(creditCard.getInstrument().getExpiryDate()).isEqualTo("12/25");
    assertThat(creditCard.getStatus()).isEqualTo(PaymentMethodStatus.VALID);
}
```

#### Test: Validate payment method capabilities
```java
@Test
public void testValidatePaymentMethodCapabilities() {
    // Arrange
    PaymentMethod creditCard = PaymentMethod.creditCard(
        "4111111111111111", "12/25", "123", "John Doe"
    );
    TransactionType transactionType = TransactionType.PURCHASE;
    Money transactionAmount = new Money(BigDecimal.valueOf(5000.00), Currency.USD);
    
    // Act
    CapabilityAssessment assessment = creditCard.assessCapabilities(
        transactionType, transactionAmount
    );
    
    // Assert
    assertThat(assessment.isCapable()).isTrue();
    assertThat(assessment.getProcessingTime()).isEqualTo(ProcessingTime.IMMEDIATE);
    assertThat(assessment.getSecurityLevel()).isEqualTo(SecurityLevel.HIGH);
    assertThat(assessment.getLimitationsReason()).isNull();
}
```

#### Test: Calculate processing cost for payment method
```java
@Test
public void testCalculateProcessingCostForPaymentMethod() {
    // Arrange
    PaymentMethod creditCard = PaymentMethod.creditCard(
        "4111111111111111", "12/25", "123", "John Doe"
    );
    PaymentMethod bankTransfer = PaymentMethod.bankTransfer(
        "123456789", "ROUTING123", "John Doe"
    );
    Money transactionAmount = new Money(BigDecimal.valueOf(1000.00), Currency.USD);
    
    // Act
    ProcessingCost creditCardCost = creditCard.calculateProcessingCost(transactionAmount);
    ProcessingCost bankTransferCost = bankTransfer.calculateProcessingCost(transactionAmount);
    
    // Assert
    assertThat(creditCardCost.getAmount())
        .isEqualTo(new Money(BigDecimal.valueOf(29.30), Currency.USD));  // 2.9% + $0.30
    assertThat(bankTransferCost.getAmount())
        .isEqualTo(new Money(BigDecimal.valueOf(5.00), Currency.USD));   // Flat $5.00 fee
}
```

### UT-004: FraudRiskScore Value Object Unit Tests

#### Test: Calculate risk score from factors
```java
@Test
public void testCalculateRiskScoreFromFactors() {
    // Arrange
    List<RiskFactor> riskFactors = Arrays.asList(
        RiskFactor.AMOUNT_ANOMALY.withWeight(25),
        RiskFactor.LOCATION_ANOMALY.withWeight(20),
        RiskFactor.VELOCITY_ANOMALY.withWeight(15)
    );
    
    // Act
    FraudRiskScore riskScore = FraudRiskScore.calculateFromFactors(riskFactors);
    
    // Assert
    assertThat(riskScore.getScore()).isEqualTo(60);  // 25 + 20 + 15
    assertThat(riskScore.getLevel()).isEqualTo(RiskLevel.MEDIUM);
    assertThat(riskScore.getFactors()).containsExactlyElementsOf(riskFactors);
    assertThat(riskScore.getConfidenceLevel()).isEqualTo(ConfidenceLevel.HIGH);
}
```

#### Test: Risk score threshold validation
```java
@Test
public void testRiskScoreThresholdValidation() {
    // Arrange
    FraudRiskScore lowRisk = new FraudRiskScore(15, RiskLevel.VERY_LOW);
    FraudRiskScore mediumRisk = new FraudRiskScore(50, RiskLevel.MEDIUM);
    FraudRiskScore highRisk = new FraudRiskScore(85, RiskLevel.HIGH);
    
    RiskThreshold autoApprovalThreshold = new RiskThreshold(40);
    RiskThreshold manualReviewThreshold = new RiskThreshold(70);
    
    // Act & Assert
    assertThat(lowRisk.exceedsThreshold(autoApprovalThreshold)).isFalse();
    assertThat(mediumRisk.exceedsThreshold(autoApprovalThreshold)).isTrue();
    assertThat(mediumRisk.exceedsThreshold(manualReviewThreshold)).isFalse();
    assertThat(highRisk.exceedsThreshold(manualReviewThreshold)).isTrue();
}
```

#### Test: Generate risk assessment report
```java
@Test
public void testGenerateRiskAssessmentReport() {
    // Arrange
    List<RiskFactor> riskFactors = Arrays.asList(
        RiskFactor.DEVICE_ANOMALY.withWeight(30),
        RiskFactor.TIME_ANOMALY.withWeight(10)
    );
    FraudRiskScore riskScore = FraudRiskScore.calculateFromFactors(riskFactors);
    
    // Act
    RiskAssessmentReport report = riskScore.generateRiskReport();
    
    // Assert
    assertThat(report.getOverallScore()).isEqualTo(40);
    assertThat(report.getRiskLevel()).isEqualTo(RiskLevel.MEDIUM);
    assertThat(report.getFactorBreakdown()).hasSize(2);
    assertThat(report.getRecommendedAction()).isEqualTo("Additional verification required");
}
```

### UT-005: PaymentGatewaySelectionService Domain Service Unit Tests

#### Test: Select optimal gateway based on success rate and cost
```java
@Test
public void testSelectOptimalGatewayBasedOnSuccessRateAndCost() {
    // Arrange
    PaymentGatewaySelectionService service = new PaymentGatewaySelectionService();
    PaymentTransaction transaction = createStandardTransaction();
    
    List<PaymentGateway> availableGateways = Arrays.asList(
        createGateway("GATEWAY_A", 97.0, BigDecimal.valueOf(2.9)),
        createGateway("GATEWAY_B", 95.0, BigDecimal.valueOf(2.1)),
        createGateway("GATEWAY_C", 99.0, BigDecimal.valueOf(3.2))
    );
    
    // Act
    GatewaySelection selection = service.selectOptimalGateway(transaction, availableGateways);
    
    // Assert
    assertThat(selection.getSelectedGateway().getId()).isEqualTo("GATEWAY_C");
    assertThat(selection.getSelectionCriteria()).contains("success_rate");
    assertThat(selection.getExpectedSuccessRate()).isEqualTo(99.0);
    assertThat(selection.getEstimatedCost()).isEqualTo(BigDecimal.valueOf(3.2));
}
```

#### Test: Gateway capability evaluation
```java
@Test
public void testGatewayCapabilityEvaluation() {
    // Arrange
    PaymentGatewaySelectionService service = new PaymentGatewaySelectionService();
    PaymentGateway gateway = createGatewayWithCapabilities(
        Arrays.asList("VISA", "MASTERCARD"),
        Arrays.asList("USD", "EUR"),
        new Money(BigDecimal.valueOf(10000.00), Currency.USD)
    );
    PaymentTransaction transaction = createTransactionWithAmex(
        new Money(BigDecimal.valueOf(500.00), Currency.USD)
    );
    
    // Act
    CapabilityAssessment assessment = service.evaluateGatewayCapabilities(
        gateway, transaction
    );
    
    // Assert
    assertThat(assessment.isCapable()).isFalse();
    assertThat(assessment.getLimitationsReason()).contains("AMEX not supported");
    assertThat(assessment.getSupportedMethods()).containsExactly("VISA", "MASTERCARD");
}
```

#### Test: Calculate expected cost for gateway
```java
@Test
public void testCalculateExpectedCostForGateway() {
    // Arrange
    PaymentGatewaySelectionService service = new PaymentGatewaySelectionService();
    PaymentGateway gateway = createGatewayWithFeeStructure(
        BigDecimal.valueOf(2.9),  // 2.9% rate
        new Money(BigDecimal.valueOf(0.30), Currency.USD)  // $0.30 fixed
    );
    PaymentTransaction transaction = createTransactionWithAmount(
        new Money(BigDecimal.valueOf(1000.00), Currency.USD)
    );
    
    // Act
    CostEstimate estimate = service.calculateExpectedCost(gateway, transaction);
    
    // Assert
    assertThat(estimate.getProcessingFee())
        .isEqualTo(new Money(BigDecimal.valueOf(29.30), Currency.USD));
    assertThat(estimate.getPercentageFee()).isEqualTo(BigDecimal.valueOf(2.9));
    assertThat(estimate.getFixedFee())
        .isEqualTo(new Money(BigDecimal.valueOf(0.30), Currency.USD));
}
```

### UT-006: FraudDetectionService Domain Service Unit Tests

#### Test: Assess fraud risk for normal transaction
```java
@Test
public void testAssessFraudRiskForNormalTransaction() {
    // Arrange
    FraudDetectionService service = new FraudDetectionService();
    PaymentTransaction transaction = createNormalTransaction();
    CustomerProfile customerProfile = createEstablishedCustomer();
    
    // Act
    FraudRiskAssessment assessment = service.assessFraudRisk(transaction, customerProfile);
    
    // Assert
    assertThat(assessment.getRiskScore().getLevel()).isEqualTo(RiskLevel.VERY_LOW);
    assertThat(assessment.getRiskScore().getScore()).isLessThan(20);
    assertThat(assessment.getRecommendedAction()).isEqualTo(SecurityAction.AUTO_APPROVE);
    assertThat(assessment.getDetectedAnomalies()).isEmpty();
}
```

#### Test: Detect suspicious transaction patterns
```java
@Test
public void testDetectSuspiciousTransactionPatterns() {
    // Arrange
    FraudDetectionService service = new FraudDetectionService();
    PaymentTransaction transaction = createSuspiciousTransaction();
    CustomerProfile customerProfile = createNewCustomer();
    
    // Act
    FraudRiskAssessment assessment = service.assessFraudRisk(transaction, customerProfile);
    
    // Assert
    assertThat(assessment.getRiskScore().getLevel()).isEqualTo(RiskLevel.HIGH);
    assertThat(assessment.getRiskScore().getScore()).isGreaterThan(70);
    assertThat(assessment.getRecommendedAction()).isEqualTo(SecurityAction.MANUAL_REVIEW);
    assertThat(assessment.getDetectedAnomalies()).isNotEmpty();
}
```

#### Test: Validate payment instrument security
```java
@Test
public void testValidatePaymentInstrumentSecurity() {
    // Arrange
    FraudDetectionService service = new FraudDetectionService();
    PaymentMethod validCard = PaymentMethod.creditCard(
        "4111111111111111", "12/25", "123", "John Doe"
    );
    PaymentMethod invalidCard = PaymentMethod.creditCard(
        "1234567890123456", "01/20", "000", "Test User"  // Invalid/expired
    );
    
    // Act
    ValidationResult validResult = service.validatePaymentInstrument(validCard);
    ValidationResult invalidResult = service.validatePaymentInstrument(invalidCard);
    
    // Assert
    assertThat(validResult.isValid()).isTrue();
    assertThat(validResult.getValidationIssues()).isEmpty();
    
    assertThat(invalidResult.isValid()).isFalse();
    assertThat(invalidResult.getValidationIssues()).contains("Expired card", "Invalid CVV");
}
```

### UT-007: Domain Events Unit Tests

#### Test: PaymentTransactionInitiated domain event
```java
@Test
public void testPaymentTransactionInitiatedDomainEvent() {
    // Arrange
    TransactionId transactionId = new TransactionId("TXN-EVENT-001");
    CustomerId customerId = new CustomerId("CUST-456");
    PaymentAmount amount = new PaymentAmount(
        new Money(BigDecimal.valueOf(199.99), Currency.USD)
    );
    PaymentMethod paymentMethod = PaymentMethod.creditCard(
        "4111111111111111", "12/25", "123", "Test Customer"
    );
    LocalDateTime initiatedAt = LocalDateTime.now();
    
    // Act
    PaymentTransactionInitiated event = new PaymentTransactionInitiated(
        transactionId, customerId, amount, paymentMethod, initiatedAt
    );
    
    // Assert
    assertThat(event.getTransactionId()).isEqualTo(transactionId);
    assertThat(event.getCustomerId()).isEqualTo(customerId);
    assertThat(event.getPaymentAmount()).isEqualTo(amount);
    assertThat(event.getPaymentMethod()).isEqualTo(paymentMethod);
    assertThat(event.getInitiatedAt()).isEqualTo(initiatedAt);
    assertThat(event.getEventType()).isEqualTo("payment.transaction.initiated");
}
```

#### Test: PaymentProcessingCompleted domain event
```java
@Test
public void testPaymentProcessingCompletedDomainEvent() {
    // Arrange
    TransactionId transactionId = new TransactionId("TXN-COMPLETE-001");
    PaymentResult result = PaymentResult.success(
        new AuthorizationCode("AUTH-789"),
        new Money(BigDecimal.valueOf(299.99), Currency.USD)
    );
    Duration processingTime = Duration.ofSeconds(2);
    LocalDateTime completedAt = LocalDateTime.now();
    
    // Act
    PaymentProcessingCompleted event = new PaymentProcessingCompleted(
        transactionId, result, processingTime, completedAt
    );
    
    // Assert
    assertThat(event.getTransactionId()).isEqualTo(transactionId);
    assertThat(event.getPaymentResult()).isEqualTo(result);
    assertThat(event.getProcessingTime()).isEqualTo(processingTime);
    assertThat(event.getCompletedAt()).isEqualTo(completedAt);
    assertThat(event.getEventType()).isEqualTo("payment.processing.completed");
}
```

## Test Utilities and Builders

### Test Data Builders
```java
public class PaymentTestDataBuilder {
    
    public static PaymentTransaction createTestPaymentTransaction() {
        return PaymentTransaction.create(
            new TransactionId("TEST-TXN-001"),
            new PaymentAmount(new Money(BigDecimal.valueOf(100.00), Currency.USD)),
            new CustomerId("TEST-CUST-001"),
            PaymentMethod.creditCard("4111111111111111", "12/25", "123", "Test User")
        );
    }
    
    public static PaymentGateway createGateway(String id, double successRate, BigDecimal cost) {
        return PaymentGateway.builder()
            .id(new GatewayId(id))
            .provider(GatewayProvider.fromId(id))
            .capabilities(GatewayCapabilities.standard())
            .costs(ProcessingCosts.percentage(cost))
            .performance(PerformanceMetrics.withSuccessRate(successRate))
            .build();
    }
    
    public static CustomerProfile createTrustedCustomerProfile() {
        return CustomerProfile.builder()
            .customerId(new CustomerId("TRUSTED-CUST-001"))
            .accountAge(Duration.ofDays(365))
            .transactionCount(150)
            .averageTransactionAmount(new Money(BigDecimal.valueOf(85.00), Currency.USD))
            .riskLevel(CustomerRiskLevel.LOW)
            .verificationStatus(VerificationStatus.VERIFIED)
            .build();
    }
}
```

### Mock Services and Repositories
```java
@ExtendWith(MockitoExtension.class)
class PaymentServiceTest {
    
    @Mock
    private PaymentGatewayRepository gatewayRepository;
    
    @Mock
    private FraudDetectionService fraudDetectionService;
    
    @Mock
    private CurrencyConversionService currencyService;
    
    @Mock
    private ComplianceValidationService complianceService;
    
    @InjectMocks
    private PaymentService paymentService;
    
    @Test
    public void testPaymentServiceCoordination() {
        // Test service methods that coordinate multiple domain components
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
- All payment processing algorithms must have 100% test coverage
- All fraud detection rules and scoring must be fully tested
- All gateway selection logic must have complete test coverage
- All currency conversion calculations must be tested
- All compliance validation rules must be covered

### Security Test Coverage
- Payment data encryption and tokenization scenarios
- Fraud detection accuracy and false positive testing
- Compliance validation for all regulatory requirements
- Security breach simulation and response testing

### Performance Test Targets
- Unit tests should complete within 20ms each
- Test suite should complete within 3 minutes
- No external dependencies in unit tests
- All tests should be deterministic and repeatable
- Comprehensive test data builders for complex payment scenarios
