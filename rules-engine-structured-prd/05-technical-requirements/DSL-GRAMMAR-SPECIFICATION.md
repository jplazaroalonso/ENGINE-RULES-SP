# DSL Grammar Specification for Rules Engine

**Version:** 1.0.0  
**Last Updated:** 2024-12-19  
**Document Type:** Technical Specification  
**Target Audience:** Development Team, Rules Authors, System Architects  
**Status:** Draft

## Table of Contents

1. [Overview](#overview)
2. [Grammar Definition](#grammar-definition)
3. [Language Constructs](#language-constructs)
4. [Data Types and Operators](#data-types-and-operators)
5. [Built-in Functions](#built-in-functions)
6. [Rule Examples](#rule-examples)
7. [Grammar Implementation](#grammar-implementation)
8. [Validation Rules](#validation-rules)
9. [Error Handling](#error-handling)
10. [Performance Considerations](#performance-considerations)

## Overview

The Rules Engine DSL (Domain-Specific Language) provides a human-readable, business-friendly syntax for expressing complex business rules. The language is designed to be interpreted by business users while being precise enough for automated execution.

### Design Principles

- **Business Readability**: Natural language constructs that business users can understand
- **Precise Semantics**: Unambiguous execution behavior
- **Extensibility**: Support for domain-specific functions and operators
- **Performance**: Optimizable for high-frequency evaluation
- **Type Safety**: Strong typing with compile-time validation

### Language Features

- Conditional logic with IF-THEN-ELSE constructs
- Arithmetic and logical expressions
- String manipulation and pattern matching
- Date and time calculations
- List and set operations
- Function calls and custom domain functions
- Variable binding and scoping
- Comments and documentation

## Grammar Definition

### ANTLR4 Grammar Specification

```antlr
grammar RulesEngine;

// Parser Rules
rule: IF condition THEN action (ELSE action)? ;

condition: expression ;

action: actionStatement+ ;

actionStatement
    : assignment
    | functionCall
    | returnStatement
    ;

expression
    : expression AND expression                             # andExpression
    | expression OR expression                              # orExpression
    | NOT expression                                        # notExpression
    | expression comparisonOperator expression              # comparisonExpression
    | expression arithmeticOperator expression              # arithmeticExpression
    | functionCall                                          # functionCallExpression
    | fieldAccess                                           # fieldAccessExpression
    | literal                                               # literalExpression
    | IDENTIFIER                                            # identifierExpression
    | LPAREN expression RPAREN                              # parenthesizedExpression
    ;

comparisonOperator
    : EQ | NE | LT | LE | GT | GE | IN | NOT_IN | CONTAINS | MATCHES
    ;

arithmeticOperator
    : PLUS | MINUS | MULTIPLY | DIVIDE | MODULO
    ;

functionCall
    : IDENTIFIER LPAREN (expression (COMMA expression)*)? RPAREN
    ;

fieldAccess
    : IDENTIFIER (DOT IDENTIFIER)*
    ;

assignment
    : IDENTIFIER ASSIGN expression
    ;

returnStatement
    : RETURN expression
    ;

literal
    : numberLiteral
    | stringLiteral
    | booleanLiteral
    | dateLiteral
    | listLiteral
    ;

numberLiteral
    : INTEGER
    | DECIMAL
    | PERCENTAGE
    ;

stringLiteral
    : STRING
    ;

booleanLiteral
    : TRUE | FALSE
    ;

dateLiteral
    : DATE_LITERAL
    ;

listLiteral
    : LBRACKET (expression (COMMA expression)*)? RBRACKET
    ;

// Lexer Rules
IF: 'IF' | 'if' ;
THEN: 'THEN' | 'then' ;
ELSE: 'ELSE' | 'else' ;
AND: 'AND' | 'and' | '&&' ;
OR: 'OR' | 'or' | '||' ;
NOT: 'NOT' | 'not' | '!' ;
RETURN: 'RETURN' | 'return' ;
TRUE: 'TRUE' | 'true' ;
FALSE: 'FALSE' | 'false' ;
IN: 'IN' | 'in' ;
NOT_IN: 'NOT IN' | 'not in' ;
CONTAINS: 'CONTAINS' | 'contains' ;
MATCHES: 'MATCHES' | 'matches' ;

// Operators
EQ: '=' | '==' ;
NE: '!=' | '<>' ;
LT: '<' ;
LE: '<=' ;
GT: '>' ;
GE: '>=' ;
ASSIGN: ':=' ;
PLUS: '+' ;
MINUS: '-' ;
MULTIPLY: '*' ;
DIVIDE: '/' ;
MODULO: '%' ;

// Delimiters
LPAREN: '(' ;
RPAREN: ')' ;
LBRACKET: '[' ;
RBRACKET: ']' ;
DOT: '.' ;
COMMA: ',' ;

// Literals
INTEGER: [0-9]+ ;
DECIMAL: [0-9]+ '.' [0-9]+ ;
PERCENTAGE: [0-9]+ ('.' [0-9]+)? '%' ;
STRING: '"' (~["\r\n])* '"' | '\'' (~['\r\n])* '\'' ;
DATE_LITERAL: '@' [0-9]{4} '-' [0-9]{2} '-' [0-9]{2} ('@' [0-9]{2} ':' [0-9]{2} (':' [0-9]{2})?)? ;

// Identifiers
IDENTIFIER: [a-zA-Z_][a-zA-Z0-9_]* ;

// Whitespace and Comments
WS: [ \t\r\n]+ -> skip ;
LINE_COMMENT: '//' ~[\r\n]* -> skip ;
BLOCK_COMMENT: '/*' .*? '*/' -> skip ;
```

## Language Constructs

### Basic Rule Structure

```
IF <condition> THEN <action> [ELSE <action>]
```

### Conditional Expressions

```dsl
// Simple condition
IF customer.tier = "GOLD" THEN discount := 15%

// Complex condition with logical operators
IF customer.tier = "GOLD" AND purchase.amount > 100 AND product.category IN ["Electronics", "Books"] 
THEN discount := 20%

// Nested conditions
IF customer.tier = "PLATINUM" THEN
    IF purchase.amount > 500 THEN
        discount := 25%
    ELSE
        discount := 15%
    END
END
```

### Variable Assignment

```dsl
// Simple assignment
baseDiscount := 10%

// Expression assignment
finalDiscount := baseDiscount + tierBonus

// Conditional assignment
adjustedAmount := IF customer.isVIP THEN purchase.amount * 0.9 ELSE purchase.amount
```

### Function Calls

```dsl
// Built-in functions
discount := MAX(tierDiscount, promotionalDiscount)
isEligible := BETWEEN(customer.age, 18, 65)
formattedName := UPPER(customer.firstName + " " + customer.lastName)

// Domain-specific functions
taxAmount := CALCULATE_TAX(purchase.amount, customer.jurisdiction)
loyaltyPoints := CALCULATE_LOYALTY_POINTS(purchase.amount, customer.tier)
```

## Data Types and Operators

### Primitive Data Types

#### Numbers
```dsl
// Integers
quantity := 5
maxItems := 100

// Decimals
price := 29.99
taxRate := 0.08

// Percentages (treated as decimals internally)
discountRate := 15%  // Equivalent to 0.15
```

#### Strings
```dsl
// String literals
productName := "iPhone 15"
category := 'Electronics'

// String concatenation
fullName := customer.firstName + " " + customer.lastName

// String interpolation
message := "Welcome back, " + customer.firstName + "!"
```

#### Booleans
```dsl
isVIP := true
isExpired := false
isEligible := customer.age >= 18 AND customer.verified = true
```

#### Dates and Times
```dsl
// Date literals
startDate := @2024-01-01
endDate := @2024-12-31@23:59:59

// Date calculations
daysSincePurchase := DAYS_BETWEEN(customer.lastPurchase, TODAY())
isWithinPeriod := BETWEEN(TODAY(), startDate, endDate)
```

#### Lists and Sets
```dsl
// List literals
validCategories := ["Electronics", "Books", "Clothing"]
discountTiers := [5%, 10%, 15%, 20%]

// List operations
isValidCategory := product.category IN validCategories
categoryCount := LENGTH(validCategories)
firstCategory := validCategories[0]
```

### Comparison Operators

```dsl
// Equality
customer.tier = "GOLD"
customer.tier != "BRONZE"

// Numeric comparisons
purchase.amount > 100
customer.age >= 18
product.stock < 10

// String matching
customer.email MATCHES ".*@company\.com"
product.name CONTAINS "iPhone"

// List membership
product.category IN ["Electronics", "Books"]
customer.preferences NOT_IN restrictedCategories
```

### Arithmetic Operators

```dsl
// Basic arithmetic
totalAmount := subtotal + tax - discount
averageOrderValue := totalRevenue / orderCount
discountAmount := baseAmount * discountRate

// Modulo for cycling
tierIndex := customer.loyaltyPoints % 4
isEvenMonth := MONTH(TODAY()) % 2 = 0
```

### Logical Operators

```dsl
// AND, OR, NOT
isEligible := customer.active AND customer.verified AND NOT customer.suspended

// Short-circuit evaluation
canProceed := customer.exists AND customer.balance > 0

// Grouped conditions
complexCondition := (customer.tier = "GOLD" OR customer.tier = "PLATINUM") 
                   AND (purchase.amount > 50 OR customer.loyaltyMember = true)
```

## Built-in Functions

### Mathematical Functions

```dsl
// Basic math
MAX(value1, value2, ...)           // Maximum value
MIN(value1, value2, ...)           // Minimum value
ABS(number)                        // Absolute value
ROUND(number, precision)           // Round to precision
CEILING(number)                    // Round up
FLOOR(number)                      // Round down
SQRT(number)                       // Square root
POW(base, exponent)               // Power

// Statistical
SUM(list)                         // Sum of list elements
AVERAGE(list)                     // Average of list elements
COUNT(list)                       // Count of list elements
```

### String Functions

```dsl
// Case conversion
UPPER(string)                     // Convert to uppercase
LOWER(string)                     // Convert to lowercase
PROPER(string)                    // Proper case (title case)

// String manipulation
LENGTH(string)                    // String length
SUBSTRING(string, start, length)  // Extract substring
TRIM(string)                      // Remove whitespace
REPLACE(string, find, replace)    // Replace text

// Pattern matching
MATCHES(string, pattern)          // Regex matching
CONTAINS(string, substring)       // Contains check
STARTS_WITH(string, prefix)       // Prefix check
ENDS_WITH(string, suffix)         // Suffix check
```

### Date and Time Functions

```dsl
// Current date/time
TODAY()                           // Current date
NOW()                            // Current date and time
CURRENT_TIME()                   // Current time

// Date arithmetic
ADD_DAYS(date, days)             // Add days to date
ADD_MONTHS(date, months)         // Add months to date
ADD_YEARS(date, years)           // Add years to date

// Date differences
DAYS_BETWEEN(date1, date2)       // Days between dates
MONTHS_BETWEEN(date1, date2)     // Months between dates
YEARS_BETWEEN(date1, date2)      // Years between dates

// Date components
YEAR(date)                       // Year component
MONTH(date)                      // Month component
DAY(date)                        // Day component
DAY_OF_WEEK(date)               // Day of week (1-7)
```

### List Functions

```dsl
// List operations
LENGTH(list)                     // List length
CONTAINS(list, item)             // Check if list contains item
FIRST(list)                      // First element
LAST(list)                       // Last element
NTH(list, index)                 // Nth element (0-based)

// List transformations
FILTER(list, condition)          // Filter list by condition
MAP(list, expression)            // Transform list elements
SORT(list)                       // Sort list
REVERSE(list)                    // Reverse list
UNIQUE(list)                     // Remove duplicates
```

### Domain-Specific Functions

#### Customer Functions
```dsl
GET_CUSTOMER_TIER(customerId)                    // Get customer tier
GET_CUSTOMER_SEGMENT(customerId)                 // Get customer segment
GET_PURCHASE_HISTORY(customerId, days)           // Get recent purchases
IS_FIRST_TIME_CUSTOMER(customerId)               // Check if first purchase
GET_CUSTOMER_PREFERENCES(customerId)             // Get customer preferences
```

#### Product Functions
```dsl
GET_PRODUCT_CATEGORY(productId)                  // Get product category
GET_PRODUCT_PRICE(productId)                     // Get current price
IS_PRODUCT_AVAILABLE(productId)                  // Check availability
GET_PRODUCT_ATTRIBUTES(productId)                // Get product attributes
GET_RELATED_PRODUCTS(productId)                  // Get related products
```

#### Financial Functions
```dsl
CALCULATE_TAX(amount, jurisdiction)               // Calculate tax
CALCULATE_SHIPPING(weight, destination)          // Calculate shipping
CONVERT_CURRENCY(amount, fromCurrency, toCurrency) // Convert currency
APPLY_DISCOUNT(amount, discountType, discountValue) // Apply discount
CALCULATE_LOYALTY_POINTS(amount, tier)           // Calculate points
```

#### Validation Functions
```dsl
IS_VALID_EMAIL(email)                           // Validate email format
IS_VALID_PHONE(phone)                           // Validate phone format
IS_VALID_ZIP_CODE(zipCode, country)             // Validate postal code
IS_BUSINESS_DAY(date)                           // Check if business day
IS_HOLIDAY(date, country)                       // Check if holiday
```

## Rule Examples

### Coupons Rules

```dsl
// Basic coupon validation
IF coupon.code = inputCode AND coupon.isActive = true AND coupon.expiryDate >= TODAY() 
THEN 
    discount := coupon.discountValue
    isValid := true
ELSE
    isValid := false
END

// Complex coupon with customer eligibility
IF coupon.code = inputCode AND coupon.isActive = true THEN
    IF customer.tier IN coupon.eligibleTiers AND 
       purchase.amount >= coupon.minimumPurchase AND
       customer.usageCount < coupon.maxUsesPerCustomer THEN
        discount := MIN(coupon.discountValue, coupon.maxDiscount)
        isValid := true
    ELSE
        isValid := false
        errorMessage := "Customer not eligible for this coupon"
    END
ELSE
    isValid := false
    errorMessage := "Invalid or expired coupon"
END
```

### Loyalty Rules

```dsl
// Tier calculation
IF customer.annualSpending >= 10000 THEN
    tier := "PLATINUM"
    multiplier := 3.0
ELSE IF customer.annualSpending >= 5000 THEN
    tier := "GOLD"
    multiplier := 2.0
ELSE IF customer.annualSpending >= 1000 THEN
    tier := "SILVER"
    multiplier := 1.5
ELSE
    tier := "BRONZE"
    multiplier := 1.0
END

// Points calculation
basePoints := FLOOR(purchase.amount)
bonusPoints := IF product.category = "Electronics" THEN basePoints * 0.5 ELSE 0
totalPoints := (basePoints + bonusPoints) * multiplier

// Tier benefits
IF customer.tier = "PLATINUM" THEN
    freeShipping := true
    prioritySupport := true
    exclusiveOffers := true
ELSE IF customer.tier = "GOLD" THEN
    freeShipping := purchase.amount > 50
    prioritySupport := true
    exclusiveOffers := false
END
```

### Promotions Rules

```dsl
// Buy-X-Get-Y promotion
IF promotion.type = "BUY_X_GET_Y" THEN
    eligibleItems := FILTER(cart.items, item.category = promotion.targetCategory)
    buyQuantity := promotion.buyQuantity
    getQuantity := promotion.getQuantity
    
    IF LENGTH(eligibleItems) >= buyQuantity THEN
        freeItems := FLOOR(LENGTH(eligibleItems) / buyQuantity) * getQuantity
        discount := SUM(FIRST(SORT(eligibleItems, item.price), freeItems), item.price)
    ELSE
        discount := 0
    END
END

// Percentage discount with conditions
IF promotion.type = "PERCENTAGE_DISCOUNT" AND
   customer.segment IN promotion.targetSegments AND
   purchase.amount >= promotion.minimumPurchase THEN
    
    eligibleAmount := SUM(FILTER(cart.items, item.category IN promotion.categories), item.price)
    discount := eligibleAmount * promotion.discountPercentage
    maxDiscount := promotion.maxDiscount
    
    finalDiscount := IF maxDiscount > 0 THEN MIN(discount, maxDiscount) ELSE discount
END
```

### Tax Calculation Rules

```dsl
// Tax calculation by jurisdiction
IF customer.jurisdiction = "CA" THEN
    // California tax rules
    baseTaxRate := 0.0725
    localTaxRate := GET_LOCAL_TAX_RATE(customer.zipCode)
    totalTaxRate := baseTaxRate + localTaxRate
    
    // Tax exemptions
    IF customer.taxExempt = true OR product.category IN ["Food", "Medicine"] THEN
        taxAmount := 0
    ELSE
        taxableAmount := purchase.amount - discounts
        taxAmount := taxableAmount * totalTaxRate
    END
ELSE IF customer.jurisdiction = "NY" THEN
    // New York tax rules
    taxRate := 0.08
    taxAmount := IF customer.taxExempt THEN 0 ELSE purchase.amount * taxRate
END

// International tax handling
IF customer.country != "US" THEN
    vatRate := GET_VAT_RATE(customer.country)
    vatAmount := purchase.amount * vatRate / (1 + vatRate)
    taxAmount := vatAmount
END
```

### Payment Rules

```dsl
// Payment method selection
IF customer.region = "EU" AND purchase.amount > 1000 THEN
    preferredMethods := ["SEPA", "BANK_TRANSFER"]
ELSE IF customer.hasStoredCard = true AND customer.fraudScore < 0.3 THEN
    preferredMethods := ["STORED_CARD", "CREDIT_CARD"]
ELSE
    preferredMethods := ["CREDIT_CARD", "PAYPAL"]
END

// Fraud detection
riskScore := 0
IF customer.isNewCustomer THEN riskScore := riskScore + 0.2
IF purchase.amount > customer.averageOrderValue * 3 THEN riskScore := riskScore + 0.3
IF customer.ipCountry != customer.billingCountry THEN riskScore := riskScore + 0.4

IF riskScore > 0.7 THEN
    requiresVerification := true
    recommendedAction := "MANUAL_REVIEW"
ELSE IF riskScore > 0.5 THEN
    requiresVerification := true
    recommendedAction := "ADDITIONAL_AUTH"
ELSE
    requiresVerification := false
    recommendedAction := "APPROVE"
END
```

## Grammar Implementation

### Parser Implementation Strategy

```java
// ANTLR4-generated parser usage
public class RuleEvaluator {
    
    public EvaluationResult evaluate(String ruleText, EvaluationContext context) {
        // Lexical analysis
        ANTLRInputStream input = new ANTLRInputStream(ruleText);
        RulesEngineLexer lexer = new RulesEngineLexer(input);
        CommonTokenStream tokens = new CommonTokenStream(lexer);
        
        // Parsing
        RulesEngineParser parser = new RulesEngineParser(tokens);
        ParseTree tree = parser.rule();
        
        // Evaluation
        RuleEvaluationVisitor visitor = new RuleEvaluationVisitor(context);
        return visitor.visit(tree);
    }
}

// Custom visitor for rule evaluation
public class RuleEvaluationVisitor extends RulesEngineBaseVisitor<Object> {
    
    private final EvaluationContext context;
    
    @Override
    public Object visitRule(RulesEngineParser.RuleContext ctx) {
        Object conditionResult = visit(ctx.condition());
        
        if (Boolean.TRUE.equals(conditionResult)) {
            return visit(ctx.action(0)); // THEN action
        } else if (ctx.action().size() > 1) {
            return visit(ctx.action(1)); // ELSE action
        }
        
        return null;
    }
    
    @Override
    public Object visitComparisonExpression(RulesEngineParser.ComparisonExpressionContext ctx) {
        Object left = visit(ctx.expression(0));
        Object right = visit(ctx.expression(1));
        String operator = ctx.comparisonOperator().getText();
        
        return ComparisonOperators.evaluate(operator, left, right);
    }
    
    // Additional visitor methods...
}
```

### Type System Implementation

```java
// Type system for DSL values
public abstract class DSLValue {
    public abstract Object getValue();
    public abstract DSLType getType();
    
    public static class NumberValue extends DSLValue {
        private final BigDecimal value;
        
        @Override
        public BigDecimal getValue() { return value; }
        
        @Override
        public DSLType getType() { return DSLType.NUMBER; }
    }
    
    public static class StringValue extends DSLValue {
        private final String value;
        
        @Override
        public String getValue() { return value; }
        
        @Override
        public DSLType getType() { return DSLType.STRING; }
    }
    
    // Additional value types...
}

// Type checking during compilation
public class TypeChecker extends RulesEngineBaseVisitor<DSLType> {
    
    @Override
    public DSLType visitComparisonExpression(RulesEngineParser.ComparisonExpressionContext ctx) {
        DSLType leftType = visit(ctx.expression(0));
        DSLType rightType = visit(ctx.expression(1));
        String operator = ctx.comparisonOperator().getText();
        
        if (!TypeCompatibility.areCompatible(leftType, rightType, operator)) {
            throw new TypeMismatchException(
                String.format("Cannot apply %s to %s and %s", operator, leftType, rightType)
            );
        }
        
        return DSLType.BOOLEAN;
    }
}
```

## Validation Rules

### Syntax Validation

1. **Grammar Compliance**: All rules must conform to the defined ANTLR grammar
2. **Token Recognition**: All tokens must be recognized by the lexer
3. **Parse Tree**: Must produce a valid parse tree without syntax errors

### Semantic Validation

1. **Type Checking**: All expressions must have compatible types
2. **Variable Declaration**: All variables must be declared before use
3. **Function Signatures**: Function calls must match defined signatures
4. **Field Access**: Object field access must be valid for the object type

### Business Rule Validation

1. **Domain Constraints**: Values must conform to business domain constraints
2. **Performance Limits**: Rules must not exceed performance thresholds
3. **Security Rules**: Rules must not access unauthorized data or functions
4. **Completeness**: Rules must handle all required scenarios

### Example Validation Implementation

```java
public class RuleValidator {
    
    public ValidationResult validate(String ruleText) {
        ValidationResult result = new ValidationResult();
        
        try {
            // Syntax validation
            ParseTree tree = parseRule(ruleText);
            result.addInfo("Syntax validation passed");
            
            // Type checking
            TypeChecker typeChecker = new TypeChecker();
            typeChecker.visit(tree);
            result.addInfo("Type checking passed");
            
            // Business validation
            BusinessRuleValidator businessValidator = new BusinessRuleValidator();
            businessValidator.visit(tree);
            result.addInfo("Business validation passed");
            
            result.setValid(true);
            
        } catch (SyntaxException e) {
            result.addError("Syntax error: " + e.getMessage());
        } catch (TypeMismatchException e) {
            result.addError("Type error: " + e.getMessage());
        } catch (BusinessRuleException e) {
            result.addWarning("Business rule warning: " + e.getMessage());
        }
        
        return result;
    }
}
```

## Error Handling

### Error Types and Recovery

1. **Syntax Errors**: Clear messages with line and column information
2. **Type Errors**: Specific type mismatch information with suggestions
3. **Runtime Errors**: Graceful handling with fallback values
4. **Performance Errors**: Timeout and resource limit handling

### Error Message Examples

```java
public class ErrorMessages {
    
    // Syntax errors
    public static final String UNEXPECTED_TOKEN = 
        "Unexpected token '%s' at line %d, column %d. Expected: %s";
    
    public static final String MISSING_THEN = 
        "Missing THEN keyword after IF condition at line %d";
    
    // Type errors
    public static final String TYPE_MISMATCH = 
        "Cannot compare %s with %s at line %d, column %d";
    
    public static final String INVALID_FUNCTION_ARGS = 
        "Function %s expects %d arguments but got %d at line %d";
    
    // Runtime errors
    public static final String DIVISION_BY_ZERO = 
        "Division by zero in expression at line %d, column %d";
    
    public static final String NULL_REFERENCE = 
        "Null reference access: %s at line %d, column %d";
}
```

## Performance Considerations

### Optimization Strategies

1. **Rule Compilation**: Pre-compile rules to bytecode for faster execution
2. **Caching**: Cache compiled rules and frequently accessed data
3. **Lazy Evaluation**: Evaluate expressions only when needed
4. **Parallel Execution**: Execute independent rules in parallel
5. **Index Optimization**: Optimize data access patterns

### Performance Targets

- **Compilation Time**: <100ms for typical business rules
- **Execution Time**: <10ms for simple rules, <100ms for complex rules
- **Memory Usage**: <1MB per compiled rule
- **Throughput**: 10,000+ rule evaluations per second per core

### Example Performance Monitoring

```java
public class PerformanceMonitor {
    
    public EvaluationResult evaluateWithMonitoring(CompiledRule rule, EvaluationContext context) {
        long startTime = System.nanoTime();
        long memoryBefore = getUsedMemory();
        
        try {
            EvaluationResult result = rule.evaluate(context);
            
            long executionTime = System.nanoTime() - startTime;
            long memoryUsed = getUsedMemory() - memoryBefore;
            
            result.setPerformanceMetrics(new PerformanceMetrics(
                executionTime / 1_000_000.0, // Convert to milliseconds
                memoryUsed
            ));
            
            return result;
            
        } catch (Exception e) {
            long executionTime = System.nanoTime() - startTime;
            throw new RuleExecutionException("Rule execution failed after " + 
                executionTime / 1_000_000.0 + "ms", e);
        }
    }
}
```

## Integration with Rules Engine

### Rule Storage and Retrieval

```java
public interface RuleRepository {
    
    CompiledRule getCompiledRule(String ruleId);
    
    void storeCompiledRule(String ruleId, CompiledRule rule);
    
    List<CompiledRule> getRulesForCategory(String category);
    
    void invalidateCache(String ruleId);
}
```

### Rule Execution Context

```java
public class EvaluationContext {
    
    private final Map<String, Object> variables = new HashMap<>();
    private final FunctionRegistry functionRegistry;
    private final SecurityContext securityContext;
    
    public Object getVariable(String name) {
        return variables.get(name);
    }
    
    public void setVariable(String name, Object value) {
        variables.put(name, value);
    }
    
    public Object callFunction(String name, Object... args) {
        return functionRegistry.call(name, args);
    }
}
```

This DSL grammar specification provides a comprehensive foundation for implementing a business-rules engine that is both powerful and accessible to business users. The grammar supports all the required bounded contexts (Coupons, Loyalty, Promotions, Taxes, Payments) while maintaining performance and extensibility requirements.
