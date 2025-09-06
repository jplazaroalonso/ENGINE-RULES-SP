# Unit Tests - Rule Evaluator/Calculator

## TC-UNIT-01: Calculation Engine Core Logic
**Test Objective**: Verify core calculation engine functionality
**Test Scope**: CalculationEngine.evaluateRules method

**Test Cases**:
```java
@Test
public void testBasicRuleEvaluation() {
    // Given
    CalculationContext context = createTestContext();
    List<Rule> rules = createTestRules();
    
    // When
    RuleEvaluationResult result = calculationEngine.evaluateRules(context, rules);
    
    // Then
    assertNotNull(result);
    assertTrue(result.getExecutions().size() > 0);
    assertNotNull(result.getFinalResult());
}

@Test
public void testPerformanceRequirement() {
    // Given
    CalculationContext context = createComplexContext();
    List<Rule> rules = createComplexRules();
    
    // When
    long startTime = System.currentTimeMillis();
    RuleEvaluationResult result = calculationEngine.evaluateRules(context, rules);
    long endTime = System.currentTimeMillis();
    
    // Then
    assertTrue("Evaluation should complete within 500ms", (endTime - startTime) < 500);
    assertNotNull(result);
}
```

## TC-UNIT-02: Conflict Resolution Logic
**Test Objective**: Verify conflict detection and resolution
**Test Scope**: ConflictResolver.resolveConflicts method

**Test Cases**:
```java
@Test
public void testConflictDetection() {
    // Given
    List<RuleExecution> conflictingExecutions = createConflictingExecutions();
    
    // When
    List<RuleConflict> conflicts = conflictResolver.detectConflicts(conflictingExecutions);
    
    // Then
    assertFalse(conflicts.isEmpty());
    assertEquals(2, conflicts.size());
}

@Test
public void testPriorityBasedResolution() {
    // Given
    List<RuleConflict> conflicts = createPriorityConflicts();
    
    // When
    ResolutionResult resolution = conflictResolver.resolveConflicts(conflicts);
    
    // Then
    assertNotNull(resolution);
    assertEquals("HIGH_PRIORITY_RULE", resolution.getSelectedRule().getId());
}
```

## TC-UNIT-03: Performance Optimization
**Test Objective**: Verify performance optimization logic
**Test Scope**: PerformanceOptimizationService.optimizeExecution method

**Test Cases**:
```java
@Test
public void testExecutionOptimization() {
    // Given
    List<Rule> rules = createFrequentlyUsedRules();
    
    // When
    OptimizedRuleSet optimizedRules = optimizationService.optimizeRuleExecution(rules);
    
    // Then
    assertNotNull(optimizedRules);
    assertTrue(optimizedRules.isOptimized());
    assertTrue(optimizedRules.hasCompiledRules());
}
```
