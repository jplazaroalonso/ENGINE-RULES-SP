# Unit Tests - Rule Creation and Management

## Domain Entity Tests

### Rule Entity Tests

#### TC-UT-RULE-01: Rule Creation and Validation
**Unit Under Test**: `Rule` entity  
**Test Scope**: Rule creation, validation, and invariants  
**Coverage Target**: 100% of Rule entity methods

**Test Cases**:
```java
// Valid rule creation
@Test
public void shouldCreateValidRule() {
    Rule rule = new Rule(
        RuleId.generate(),
        "Summer Discount",
        "Electronics discount for summer",
        RuleStatus.DRAFT,
        Priority.HIGH,
        "category = 'electronics' AND amount > 100"
    );
    
    assertThat(rule.isValid()).isTrue();
    assertThat(rule.getStatus()).isEqualTo(RuleStatus.DRAFT);
    assertThat(rule.getPriority()).isEqualTo(Priority.HIGH);
}

// Rule validation with invalid DSL
@Test
public void shouldFailValidationWithInvalidDSL() {
    assertThrows(InvalidRuleException.class, () -> {
        new Rule(
            RuleId.generate(),
            "Invalid Rule",
            "Test rule",
            RuleStatus.DRAFT,
            Priority.LOW,
            "invalid DSL syntax here"
        );
    });
}

// Rule status transitions
@Test
public void shouldAllowValidStatusTransitions() {
    Rule rule = createValidDraftRule();
    
    rule.submitForReview();
    assertThat(rule.getStatus()).isEqualTo(RuleStatus.UNDER_REVIEW);
    
    rule.approve("admin@company.com");
    assertThat(rule.getStatus()).isEqualTo(RuleStatus.APPROVED);
    assertThat(rule.getApprovedBy()).isEqualTo("admin@company.com");
    
    rule.activate();
    assertThat(rule.getStatus()).isEqualTo(RuleStatus.ACTIVE);
}

// Invalid status transitions
@Test
public void shouldRejectInvalidStatusTransitions() {
    Rule rule = createValidDraftRule();
    
    assertThrows(IllegalStateTransitionException.class, () -> {
        rule.activate(); // Cannot activate from DRAFT
    });
}
```

#### TC-UT-RULE-02: Rule Business Logic
**Unit Under Test**: Rule business logic methods  
**Test Scope**: Rule operations and business invariants

**Test Cases**:
```java
// Rule cloning
@Test
public void shouldCloneRuleWithNewIdentity() {
    Rule original = createValidActiveRule();
    Rule cloned = original.clone();
    
    assertThat(cloned.getId()).isNotEqualTo(original.getId());
    assertThat(cloned.getName()).isEqualTo(original.getName() + " (Copy)");
    assertThat(cloned.getStatus()).isEqualTo(RuleStatus.DRAFT);
    assertThat(cloned.getDslContent()).isEqualTo(original.getDslContent());
}

// Rule priority comparison
@Test
public void shouldCompareRulePriorityCorrectly() {
    Rule highPriorityRule = createRuleWithPriority(Priority.HIGH);
    Rule lowPriorityRule = createRuleWithPriority(Priority.LOW);
    
    assertThat(highPriorityRule.compareTo(lowPriorityRule)).isLessThan(0);
    assertThat(lowPriorityRule.compareTo(highPriorityRule)).isGreaterThan(0);
}

// Rule equality
@Test
public void shouldTestRuleEqualityCorrectly() {
    RuleId id = RuleId.generate();
    Rule rule1 = createRuleWithId(id);
    Rule rule2 = createRuleWithId(id);
    Rule rule3 = createRuleWithId(RuleId.generate());
    
    assertThat(rule1).isEqualTo(rule2);
    assertThat(rule1).isNotEqualTo(rule3);
    assertThat(rule1.hashCode()).isEqualTo(rule2.hashCode());
}
```

### Value Object Tests

#### TC-UT-VO-01: RuleId Value Object
**Unit Under Test**: `RuleId` value object  
**Test Scope**: Immutability, equality, validation

**Test Cases**:
```java
@Test
public void shouldGenerateUniqueRuleIds() {
    RuleId id1 = RuleId.generate();
    RuleId id2 = RuleId.generate();
    
    assertThat(id1).isNotEqualTo(id2);
    assertThat(id1.getValue()).isNotEmpty();
    assertThat(id2.getValue()).isNotEmpty();
}

@Test
public void shouldValidateRuleIdFormat() {
    assertThrows(IllegalArgumentException.class, () -> {
        RuleId.of(null);
    });
    
    assertThrows(IllegalArgumentException.class, () -> {
        RuleId.of("");
    });
    
    assertThrows(IllegalArgumentException.class, () -> {
        RuleId.of("invalid-format");
    });
}

@Test
public void shouldCreateValidRuleIdFromString() {
    String validId = "RULE-12345678-9ABC-DEF0-1234-567890ABCDEF";
    RuleId ruleId = RuleId.of(validId);
    
    assertThat(ruleId.getValue()).isEqualTo(validId);
}
```

#### TC-UT-VO-02: DSL Content Value Object
**Unit Under Test**: `DSLContent` value object  
**Test Scope**: DSL validation, parsing, optimization

**Test Cases**:
```java
@Test
public void shouldValidateDSLSyntax() {
    DSLContent validDSL = DSLContent.of("category = 'electronics' AND amount > 100");
    assertThat(validDSL.isValid()).isTrue();
    
    assertThrows(InvalidDSLException.class, () -> {
        DSLContent.of("invalid syntax here");
    });
}

@Test
public void shouldOptimizeDSLExpressions() {
    DSLContent dsl = DSLContent.of("amount > 100 AND amount > 50");
    DSLContent optimized = dsl.optimize();
    
    assertThat(optimized.getContent()).isEqualTo("amount > 100");
}

@Test
public void shouldExtractDSLVariables() {
    DSLContent dsl = DSLContent.of("category = 'electronics' AND customer.tier = 'gold'");
    Set<String> variables = dsl.getVariables();
    
    assertThat(variables).containsExactly("category", "customer.tier");
}
```

## Domain Service Tests

### RuleValidationService Tests

#### TC-UT-VAL-01: Syntax Validation
**Unit Under Test**: `RuleValidationService`  
**Test Scope**: DSL syntax validation, error reporting

**Test Cases**:
```java
@Test
public void shouldValidateCorrectDSLSyntax() {
    RuleValidationService validator = new RuleValidationService();
    Rule rule = createRuleWithDSL("category = 'electronics' AND amount > 100");
    
    ValidationResult result = validator.validateSyntax(rule);
    
    assertThat(result.isValid()).isTrue();
    assertThat(result.getErrors()).isEmpty();
}

@Test
public void shouldDetectSyntaxErrors() {
    RuleValidationService validator = new RuleValidationService();
    Rule rule = createRuleWithDSL("category == 'electronics' AND amount >");
    
    ValidationResult result = validator.validateSyntax(rule);
    
    assertThat(result.isValid()).isFalse();
    assertThat(result.getErrors()).hasSize(2);
    assertThat(result.getErrors().get(0).getMessage()).contains("Invalid operator");
    assertThat(result.getErrors().get(1).getPosition()).isEqualTo(45);
}

@Test
public void shouldProvideHelpfulErrorMessages() {
    RuleValidationService validator = new RuleValidationService();
    Rule rule = createRuleWithDSL("amount >> 100");
    
    ValidationResult result = validator.validateSyntax(rule);
    
    assertThat(result.getErrors().get(0).getMessage())
        .contains("Invalid operator '>>'")
        .contains("Did you mean '>' or '>='?");
}
```

#### TC-UT-VAL-02: Business Logic Validation
**Unit Under Test**: `RuleValidationService`  
**Test Scope**: Business rule validation, conflict detection

**Test Cases**:
```java
@Test
public void shouldDetectRuleConflicts() {
    RuleValidationService validator = new RuleValidationService();
    Rule existingRule = createActiveRule("category = 'electronics'", "discount = 10%");
    Rule newRule = createDraftRule("category = 'electronics'", "discount = 15%");
    
    when(ruleRepository.findActiveRules()).thenReturn(List.of(existingRule));
    
    ValidationResult result = validator.validateBusinessLogic(newRule);
    
    assertThat(result.hasConflicts()).isTrue();
    assertThat(result.getConflicts()).hasSize(1);
    assertThat(result.getConflicts().get(0).getType()).isEqualTo(ConflictType.OVERLAPPING_CONDITIONS);
}

@Test
public void shouldValidateRuleComplexity() {
    RuleValidationService validator = new RuleValidationService();
    Rule complexRule = createRuleWithComplexDSL(50); // Very complex DSL
    
    ValidationResult result = validator.validateBusinessLogic(complexRule);
    
    assertThat(result.hasWarnings()).isTrue();
    assertThat(result.getWarnings()).anyMatch(w -> 
        w.getType() == WarningType.HIGH_COMPLEXITY);
}
```

### RuleTemplateService Tests

#### TC-UT-TEMP-01: Template Management
**Unit Under Test**: `RuleTemplateService`  
**Test Scope**: Template CRUD operations, categorization

**Test Cases**:
```java
@Test
public void shouldCreateTemplateSuccessfully() {
    RuleTemplateService service = new RuleTemplateService(templateRepository);
    RuleTemplate template = createDiscountTemplate();
    
    RuleTemplate saved = service.saveTemplate(template);
    
    assertThat(saved.getId()).isNotNull();
    assertThat(saved.getName()).isEqualTo("Percentage Discount");
    assertThat(saved.getCategory()).isEqualTo(TemplateCategory.PROMOTIONS);
    verify(templateRepository).save(template);
}

@Test
public void shouldRetrieveTemplatesByCategory() {
    RuleTemplateService service = new RuleTemplateService(templateRepository);
    List<RuleTemplate> promotionTemplates = List.of(
        createDiscountTemplate(),
        createBuyOneGetOneTemplate()
    );
    
    when(templateRepository.findByCategory(TemplateCategory.PROMOTIONS))
        .thenReturn(promotionTemplates);
    
    List<RuleTemplate> result = service.getTemplatesByCategory(TemplateCategory.PROMOTIONS);
    
    assertThat(result).hasSize(2);
    assertThat(result).extracting(RuleTemplate::getName)
        .containsExactly("Percentage Discount", "Buy One Get One");
}
```

#### TC-UT-TEMP-02: Template Application
**Unit Under Test**: `RuleTemplateService`  
**Test Scope**: Template application to rules, variable substitution

**Test Cases**:
```java
@Test
public void shouldApplyTemplateToRule() {
    RuleTemplateService service = new RuleTemplateService(templateRepository);
    RuleTemplate template = createLoyaltyPointsTemplate();
    Map<String, Object> variables = Map.of(
        "points_multiplier", 2.0,
        "customer_tier", "gold"
    );
    
    Rule rule = service.applyTemplate(template, "Gold Customer Points", variables);
    
    assertThat(rule.getName()).isEqualTo("Gold Customer Points");
    assertThat(rule.getDslContent().getContent())
        .contains("points = amount * 2.0")
        .contains("customer_tier = 'gold'");
    assertThat(rule.getStatus()).isEqualTo(RuleStatus.DRAFT);
}

@Test
public void shouldValidateTemplateVariables() {
    RuleTemplateService service = new RuleTemplateService(templateRepository);
    RuleTemplate template = createTemplateWithRequiredVariables("discount_percentage");
    Map<String, Object> incompleteVariables = Map.of("customer_tier", "gold");
    
    assertThrows(MissingTemplateVariableException.class, () -> {
        service.applyTemplate(template, "Incomplete Rule", incompleteVariables);
    });
}
```

## Application Service Tests

### RuleCreationService Tests

#### TC-UT-CREATE-01: Rule Creation Workflow
**Unit Under Test**: `RuleCreationService`  
**Test Scope**: Complete rule creation workflow

**Test Cases**:
```java
@Test
public void shouldCreateRuleFromTemplate() {
    RuleCreationService service = new RuleCreationService(
        ruleRepository, templateService, validationService, eventPublisher);
    
    CreateRuleCommand command = CreateRuleCommand.builder()
        .templateId(DISCOUNT_TEMPLATE_ID)
        .ruleName("Summer Electronics Discount")
        .variables(Map.of("discount_percentage", 20))
        .createdBy("analyst@company.com")
        .build();
    
    RuleId ruleId = service.createRuleFromTemplate(command);
    
    assertThat(ruleId).isNotNull();
    verify(ruleRepository).save(any(Rule.class));
    verify(eventPublisher).publish(any(RuleCreatedEvent.class));
}

@Test
public void shouldCreateRuleFromScratch() {
    RuleCreationService service = new RuleCreationService(
        ruleRepository, templateService, validationService, eventPublisher);
    
    CreateRuleCommand command = CreateRuleCommand.builder()
        .ruleName("Custom Loyalty Rule")
        .description("Custom rule for special customers")
        .dslContent("customer.vip = true AND purchase_amount > 500")
        .priority(Priority.HIGH)
        .createdBy("analyst@company.com")
        .build();
    
    RuleId ruleId = service.createRuleFromScratch(command);
    
    assertThat(ruleId).isNotNull();
    verify(validationService).validateSyntax(any(Rule.class));
    verify(ruleRepository).save(any(Rule.class));
}
```

#### TC-UT-CREATE-02: Rule Testing
**Unit Under Test**: `RuleCreationService`  
**Test Scope**: Rule testing functionality

**Test Cases**:
```java
@Test
public void shouldTestRuleWithSampleData() {
    RuleCreationService service = new RuleCreationService(
        ruleRepository, templateService, validationService, eventPublisher);
    
    Rule rule = createValidRule();
    TestData testData = TestData.builder()
        .category("electronics")
        .amount(BigDecimal.valueOf(150))
        .customerId("CUST001")
        .build();
    
    TestResult result = service.testRule(rule, testData);
    
    assertThat(result.isSuccessful()).isTrue();
    assertThat(result.getExecutionTime()).isLessThan(Duration.ofMillis(100));
    assertThat(result.getAppliedActions()).hasSize(1);
    assertThat(result.getAppliedActions().get(0).getType()).isEqualTo(ActionType.DISCOUNT);
}

@Test
public void shouldHandleRuleTestFailures() {
    RuleCreationService service = new RuleCreationService(
        ruleRepository, templateService, validationService, eventPublisher);
    
    Rule invalidRule = createRuleWithInvalidDSL();
    TestData testData = createValidTestData();
    
    TestResult result = service.testRule(invalidRule, testData);
    
    assertThat(result.isSuccessful()).isFalse();
    assertThat(result.getErrorMessage()).isNotEmpty();
    assertThat(result.getAppliedActions()).isEmpty();
}
```

## Repository Tests

### RuleRepository Tests

#### TC-UT-REPO-01: Rule Persistence
**Unit Under Test**: `RuleRepository`  
**Test Scope**: Rule CRUD operations

**Test Cases**:
```java
@Test
@Transactional
public void shouldSaveAndRetrieveRule() {
    Rule rule = createValidRule();
    
    Rule saved = ruleRepository.save(rule);
    Optional<Rule> retrieved = ruleRepository.findById(saved.getId());
    
    assertThat(retrieved).isPresent();
    assertThat(retrieved.get().getName()).isEqualTo(rule.getName());
    assertThat(retrieved.get().getDslContent()).isEqualTo(rule.getDslContent());
}

@Test
@Transactional
public void shouldFindRulesByStatus() {
    Rule draftRule = createRuleWithStatus(RuleStatus.DRAFT);
    Rule activeRule = createRuleWithStatus(RuleStatus.ACTIVE);
    ruleRepository.saveAll(List.of(draftRule, activeRule));
    
    List<Rule> draftRules = ruleRepository.findByStatus(RuleStatus.DRAFT);
    List<Rule> activeRules = ruleRepository.findByStatus(RuleStatus.ACTIVE);
    
    assertThat(draftRules).hasSize(1);
    assertThat(activeRules).hasSize(1);
    assertThat(draftRules.get(0).getId()).isEqualTo(draftRule.getId());
    assertThat(activeRules.get(0).getId()).isEqualTo(activeRule.getId());
}
```

#### TC-UT-REPO-02: Query Performance
**Unit Under Test**: `RuleRepository`  
**Test Scope**: Query optimization and performance

**Test Cases**:
```java
@Test
public void shouldFindActiveRulesEfficiently() {
    // Create large dataset
    List<Rule> rules = createRulesWithMixedStatuses(1000);
    ruleRepository.saveAll(rules);
    
    long startTime = System.currentTimeMillis();
    List<Rule> activeRules = ruleRepository.findActiveRulesForEvaluation();
    long endTime = System.currentTimeMillis();
    
    assertThat(endTime - startTime).isLessThan(100); // Under 100ms
    assertThat(activeRules).allMatch(rule -> rule.getStatus() == RuleStatus.ACTIVE);
}
```

## Coverage Requirements

### Statement Coverage: ≥80%
- All domain entities and value objects: 90%
- Domain services: 85%
- Application services: 80%
- Repository implementations: 75%

### Branch Coverage: ≥80%
- Validation logic: 90%
- Status transition logic: 95%
- Error handling paths: 80%
- Business rule conditions: 85%

### Method Coverage: 100%
- All public methods must have tests
- All business logic methods covered
- All error scenarios tested

## DDD Unit Testing Scope

### Entities
- **Rule**: Business logic, invariants, state transitions
- **RuleTemplate**: Template application, variable substitution
- **TestResult**: Result calculation, execution tracking

### Value Objects
- **RuleId**: Immutability, validation, uniqueness
- **DSLContent**: Syntax validation, optimization, variable extraction
- **Priority**: Ordering, comparison, business rules

### Aggregates
- **Rule Aggregate**: Consistency boundaries, transaction integrity
- **Template Aggregate**: Template management, versioning

### Domain Services
- **RuleValidationService**: Complex business logic, conflict detection
- **RuleTemplateService**: Cross-aggregate operations, template management
- **ConflictDetectionService**: Business rule validation

### Domain Events
- **RuleCreatedEvent**: Event creation, payload validation
- **RuleUpdatedEvent**: Event handling, state consistency
- **RuleApprovedEvent**: Workflow integration

### Specifications
- **ActiveRuleSpecification**: Business rule validation
- **ConflictingRuleSpecification**: Complex rule queries
- **TemplateMatchSpecification**: Template selection logic

## Test Data Management

### Test Fixtures
```java
public class RuleTestFixtures {
    public static Rule createValidDraftRule() {
        return Rule.builder()
            .id(RuleId.generate())
            .name("Test Discount Rule")
            .description("Test rule for unit testing")
            .status(RuleStatus.DRAFT)
            .priority(Priority.MEDIUM)
            .dslContent(DSLContent.of("category = 'electronics' AND amount > 100"))
            .createdBy("test@company.com")
            .createdAt(LocalDateTime.now())
            .build();
    }
    
    public static RuleTemplate createDiscountTemplate() {
        return RuleTemplate.builder()
            .id(TemplateId.generate())
            .name("Percentage Discount")
            .category(TemplateCategory.PROMOTIONS)
            .dslTemplate("category = '${category}' AND amount > ${min_amount}")
            .variables(Set.of("category", "min_amount", "discount_percentage"))
            .build();
    }
}
```

### Mock Configuration
```java
@ExtendWith(MockitoExtension.class)
class RuleCreationServiceTest {
    @Mock private RuleRepository ruleRepository;
    @Mock private RuleTemplateService templateService;
    @Mock private RuleValidationService validationService;
    @Mock private DomainEventPublisher eventPublisher;
    
    @InjectMocks private RuleCreationService ruleCreationService;
    
    @BeforeEach
    void setUp() {
        when(validationService.validateSyntax(any())).thenReturn(ValidationResult.success());
        when(ruleRepository.save(any())).thenAnswer(i -> i.getArgument(0));
    }
}
```

## Performance Testing

### Unit Test Performance Requirements
- Individual test execution: <50ms
- Test suite execution: <5 minutes
- Memory usage per test: <10MB
- Database tests: <100ms per test

### Load Testing for Repository
- Concurrent read operations: 100 threads
- Bulk insert operations: 1000 records in <1 second
- Query performance: Complex queries <100ms
