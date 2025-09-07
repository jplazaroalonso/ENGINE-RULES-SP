# Ubiquitous Language Dictionary - Rules Engine Domain

**Extracted from PRD Sources**: `04-functional-requirements/`, `stories.md`, `domain/model.md`, `acceptance.md` files from all features

## Core Concepts

### Rule
- **Definition**: A business rule that defines conditional logic and corresponding actions to be executed when specific conditions are met during transaction processing
- **Context**: Central concept across all bounded contexts, representing encapsulated business logic
- **Synonyms**: Business Rule, Logic Rule, Decision Rule, Promotional Rule
- **Related Concepts**: Rule Template, DSL Content, Rule Evaluation, Rule Execution, Rule Status
- **Business Rules**: 
  - Must have valid DSL syntax and semantic consistency
  - Must pass through approval workflow before activation (AC-05: DRAFT → UNDER_REVIEW → APPROVED → ACTIVE)
  - Cannot be modified while in ACTIVE status (extracted from AC-05 status transitions)
  - Must maintain complete audit trail of all changes
- **Examples**: 
  - "10% discount for electronics purchases over $200" (from US-0001)
  - "Double loyalty points for gold tier customers" (from US-0002)
  - "Free shipping for orders over $50"

### Rule Engine
- **Definition**: The comprehensive system that manages the complete lifecycle of business rules from creation to execution, including validation, approval, and real-time evaluation
- **Context**: System-wide concept encompassing all rule-related capabilities
- **Synonyms**: Rules Platform, Business Rules Management System (BRMS), Rules Processing Platform
- **Related Concepts**: Rule Evaluation Engine, Rule Management, Rule Repository, Rule Lifecycle
- **Business Rules**: 
  - Must maintain 99.9% availability (from success criteria)
  - Must process rules within 500ms response time (from AC-04: testing within 3 seconds)
  - Must support 1000+ transactions per second (from executive summary)
- **Examples**: The entire Rules Engine platform processing promotional rules during Black Friday sales

### DSL (Domain-Specific Language)
- **Definition**: A specialized language designed specifically for expressing business rules in a syntax that is both human-readable for business users and machine-executable
- **Context**: Used throughout rule creation and management processes (AC-02: DSL rule definition)
- **Synonyms**: Rule Language, Business Rule Syntax, Rule Expression Language, Rule Definition Language
- **Related Concepts**: Rule Content, Syntax Validation, Template Parameters, Rule Logic
- **Business Rules**: 
  - Must be syntactically valid (AC-02: real-time validation feedback)
  - Must pass semantic validation (AC-03: validation within 2 seconds)
  - Must not exceed maximum complexity limits (from domain model constraints)
  - Must provide syntax highlighting and auto-completion (AC-02)
- **Examples**: 
  - `customer.tier = 'GOLD' AND purchase.category = 'electronics' THEN discount.percentage = 15`
  - `order.amount > 100 AND customer.segment = 'VIP' THEN shipping.cost = 0`

### Transaction
- **Definition**: A business transaction that triggers rule evaluation, containing customer information, product details, and contextual data necessary for rule processing
- **Context**: Input to rule evaluation process across evaluation and calculation contexts
- **Synonyms**: Business Transaction, Purchase Transaction, Commerce Event, Order Transaction
- **Related Concepts**: Evaluation Context, Customer Data, Transaction Data, Context Data
- **Business Rules**: 
  - Must contain valid customer and transaction identifiers
  - Must include all required data for rule evaluation (from AC-04: sample transaction data)
  - Must be processed within system performance limits (<500ms evaluation)
- **Examples**: 
  - Customer purchasing electronics items totaling $250 from e-commerce website
  - In-store purchase with loyalty card scan and multiple product categories
  - Mobile app transaction with location-based promotions

### Customer
- **Definition**: The entity representing an individual or organization that participates in transactions and is subject to rule evaluation for benefits, discounts, or restrictions
- **Context**: Key entity in rule conditions and evaluation criteria (from user stories US-0001, US-0002)
- **Synonyms**: Client, Buyer, Account Holder, Shopper
- **Related Concepts**: Customer Segment, Customer Tier, Customer Profile, Customer Behavior
- **Business Rules**: 
  - Must have unique identifier across all systems
  - Must belong to defined customer segments for rule targeting
  - Must have valid profile information for rule condition evaluation
- **Examples**: 
  - Gold tier customer with ID "CUST-12345" making a purchase
  - New customer eligible for first-time buyer promotions
  - VIP customer with special pricing privileges

## Domain Entities

### Rule Aggregate
- **Purpose**: Central aggregate that encapsulates all rule-related data and behavior, serving as the consistency boundary for rule operations
- **Identity**: Unique RuleId (UUID format) ensuring global uniqueness
- **Lifecycle**: 
  - **Creation**: DRAFT status with basic validation (US-0001, US-0002: create rules)
  - **Modification**: Version-controlled changes with approval workflow (US-0004: edit existing rules)
  - **Testing**: Rule testing with sample data (US-0006: test rule with sample data)
  - **Approval**: Multi-level approval process (US-0008: submit rule for approval)
  - **Activation**: Transition to ACTIVE status after approval
  - **Deactivation**: Transition to INACTIVE or ARCHIVED status
- **Business Rules**: 
  - Only valid status transitions allowed: DRAFT → UNDER_REVIEW → APPROVED → ACTIVE (AC-05)
  - Cannot modify content while in UNDER_REVIEW or ACTIVE status (from rule modification policy)
  - Must maintain complete audit trail of all changes (US-0009: view rule history)
  - Must pass validation before status transitions (AC-03: validation check)
- **Relationships**: 
  - Uses Rule Templates for creation (US-0003: use rule templates)
  - References Rule Categories for organization
  - Produces Domain Events for lifecycle changes
  - Can be cloned to create new rules (US-0010: clone existing rules)

### Rule Template
- **Purpose**: Reusable pattern that provides structure and default values for creating new rules of specific types
- **Identity**: Unique TemplateId with category-based organization
- **Lifecycle**: 
  - **Creation**: Template design with parameter definitions
  - **Application**: Applied during rule creation with parameter substitution (AC-06: template application)
  - **Customization**: User can modify template content while preserving original (AC-06)
  - **Versioning**: Template versions for backward compatibility
- **Business Rules**: 
  - Must define valid parameter structure for customization
  - Must generate syntactically correct DSL when applied (AC-06)
  - Original template must remain unchanged when customized (AC-06)
  - Must be organized by categories (AC-01: templates organized by category)
- **Relationships**: 
  - Used by Rule Creation process (US-0003)
  - Categorized by business domain (Promotions, Loyalty, Coupons)
  - Referenced in Rule metadata for traceability

### Evaluation Context
- **Purpose**: Container that holds all necessary data and state for rule evaluation, including transaction details, customer information, and environmental context
- **Identity**: Unique EvaluationId for each evaluation session
- **Lifecycle**: 
  - **Creation**: Initialized with transaction and customer data
  - **Processing**: Populated with applicable rules and evaluation state
  - **Execution**: Rules executed in priority order with conflict resolution
  - **Completion**: Contains final results and execution metrics
- **Business Rules**: 
  - Must contain valid transaction and customer identifiers
  - Must include all data required by applicable rules (AC-04: sample transaction data)
  - Must maintain evaluation state throughout processing
  - Must complete within performance thresholds (<500ms for 95% of requests)
- **Relationships**: 
  - References Customer and Transaction entities
  - Contains Evaluation Results
  - Links to Applied Rules
  - Produces Performance Metrics

### Approval Workflow
- **Purpose**: Manages the multi-step approval process for rules ensuring proper governance and compliance
- **Identity**: Unique WorkflowId linked to specific Rule
- **Lifecycle**: 
  - **Initiation**: Started when rule submitted for approval (US-0008)
  - **Processing**: Multi-level approval with role-based permissions
  - **Decision**: Approved or rejected with comments and audit trail
  - **Completion**: Rule status updated based on approval decision
- **Business Rules**: 
  - Must enforce role-based approval permissions
  - Must maintain complete audit trail of approval decisions
  - Must validate rule before approval (conflict check, impact analysis)
  - Must support rejection with feedback for rule improvement
- **Relationships**: 
  - Linked to specific Rule instance
  - References User roles and permissions
  - Produces Approval Events for audit trail

## Value Objects

### RuleId
- **Purpose**: Unique identifier for rules ensuring global uniqueness and referential integrity
- **Attributes**: String value in UUID v4 format
- **Equality**: Based on string value comparison
- **Immutability**: Immutable once created
- **Validation Rules**: 
  - Must be valid UUID v4 format (from domain model constraints)
  - Must be unique across entire system
  - Cannot be null or empty

### RuleName
- **Purpose**: Human-readable identifier for rules that provides meaningful business context
- **Attributes**: String value with specific constraints
- **Equality**: Case-insensitive string comparison within same category
- **Immutability**: Can be modified but must maintain uniqueness
- **Validation Rules**: 
  - 3-100 characters in length (from domain model constraints)
  - Alphanumeric and spaces only
  - Must be unique within the same category (from business rule invariants)
  - Cannot contain special characters except spaces and hyphens

### DSLContent
- **Purpose**: Encapsulates the business logic of a rule in domain-specific language format
- **Attributes**: String containing DSL syntax with metadata
- **Equality**: Based on normalized DSL content comparison
- **Immutability**: Content is immutable but can be versioned
- **Validation Rules**: 
  - Must pass syntax validation (AC-02: syntax highlighting and validation)
  - Must pass semantic validation (AC-03: business logic consistency)
  - Must not exceed maximum complexity limits (from domain model)
  - Must reference only valid attributes and functions

### Priority
- **Purpose**: Defines the execution order and conflict resolution precedence for rules
- **Attributes**: Enumerated value (CRITICAL, HIGH, MEDIUM, LOW)
- **Equality**: Based on enumeration value
- **Immutability**: Immutable value object
- **Validation Rules**: 
  - Must be one of the defined priority levels
  - Within same category, no two active rules can have identical priority (from business rule invariants)
  - Cannot be null

### Money
- **Purpose**: Represents monetary amounts with currency information for accurate financial calculations
- **Attributes**: Amount (BigDecimal) and Currency Code (ISO 4217)
- **Equality**: Based on amount and currency comparison
- **Immutability**: Immutable once created
- **Validation Rules**: 
  - Amount must be non-negative for discounts
  - Currency must be valid ISO 4217 code
  - Precision must not exceed currency decimal places

### PercentageValue
- **Purpose**: Represents percentage values with proper range validation for discount calculations
- **Attributes**: Numeric value with percentage constraints
- **Equality**: Based on numeric value comparison
- **Immutability**: Immutable once created
- **Validation Rules**: 
  - Must be between 0 and 100 for standard percentages
  - May exceed 100 for special loyalty multipliers
  - Must have appropriate decimal precision

## Domain Events

### RuleCreated
- **Trigger**: When a new rule is successfully created in the system (US-0001, US-0002, US-0003)
- **Payload**: 
  - RuleId: Unique identifier of created rule
  - RuleName: Name of the created rule
  - Category: Business category of the rule
  - CreatedBy: User who created the rule
  - CreatedAt: Timestamp of creation
  - TemplateUsed: Template used for creation (if any)
- **Consequences**: 
  - Rule appears in management interfaces
  - Validation processes are triggered
  - Audit log entry is created
  - Template usage statistics updated
- **Timing**: Immediately after successful rule persistence
- **Subscribers**: 
  - Rule Management UI for display updates
  - Audit Service for compliance logging
  - Analytics Service for usage metrics
  - Template Service for usage tracking

### RuleValidated
- **Trigger**: When rule validation is completed (AC-03: validation within 2 seconds)
- **Payload**: 
  - RuleId: Identifier of validated rule
  - ValidationResults: Syntax and semantic validation results
  - Errors: List of validation errors with line numbers
  - Warnings: List of warnings including conflict warnings
  - ValidatedAt: Timestamp of validation
- **Consequences**: 
  - Validation feedback displayed to user
  - Rule status updated based on validation results
  - Performance impact assessment available
- **Timing**: After validation process completion
- **Subscribers**: 
  - Rule Management UI for feedback display
  - Performance Monitoring for impact analysis
  - Conflict Detection Service for warnings

### RuleApproved
- **Trigger**: When a rule successfully completes the approval workflow (US-0008)
- **Payload**: 
  - RuleId: Identifier of approved rule
  - ApprovedBy: User who approved the rule
  - ApprovedAt: Timestamp of approval
  - ApprovalComments: Comments from approver
  - PreviousStatus: Previous status before approval
- **Consequences**: 
  - Rule becomes eligible for activation
  - Stakeholders are notified of approval
  - Rule can be scheduled for deployment
  - Approval audit trail updated
- **Timing**: Immediately after approval decision is recorded
- **Subscribers**: 
  - Notification Service for stakeholder alerts
  - Deployment Service for activation scheduling
  - Rule Management UI for status updates
  - Audit Service for compliance tracking

### RuleActivated
- **Trigger**: When a rule transitions to ACTIVE status and becomes available for evaluation
- **Payload**: 
  - RuleId: Identifier of activated rule
  - RuleName: Name of activated rule
  - DSLContent: Rule logic for evaluation
  - Priority: Rule priority for execution order
  - ActivatedAt: Timestamp of activation
  - EffectiveFrom: When rule becomes effective
  - ActivatedBy: User or system that activated the rule
- **Consequences**: 
  - Rule becomes available in evaluation engine
  - Cache updates to include new rule
  - Monitoring begins for rule performance
  - Conflict detection runs against active rules
- **Timing**: When rule status changes to ACTIVE
- **Subscribers**: 
  - Rule Evaluation Engine for rule loading
  - Cache Management Service for cache updates
  - Performance Monitoring Service for tracking
  - Conflict Detection Service for analysis

### RuleEvaluationCompleted
- **Trigger**: When rule evaluation process completes for a transaction (AC-04: test results within 3 seconds)
- **Payload**: 
  - EvaluationId: Unique evaluation session identifier
  - TransactionId: Transaction that was evaluated
  - AppliedRules: List of rules that were applied
  - TotalExecutionTime: Total time for evaluation
  - Results: Final calculation results
  - ConflictsResolved: Any conflicts that were resolved
- **Consequences**: 
  - Results are returned to requesting system
  - Performance metrics are updated
  - Business analytics are updated
  - Rule effectiveness data collected
- **Timing**: After all applicable rules have been evaluated
- **Subscribers**: 
  - Requesting System for result processing
  - Analytics Service for business metrics
  - Performance Monitoring for system metrics
  - Rule Optimization Service for effectiveness analysis

### ConflictDetected
- **Trigger**: When conflicting rules are identified during evaluation or validation
- **Payload**: 
  - ConflictId: Unique conflict identifier
  - ConflictingRuleIds: Rules that are in conflict
  - ConflictType: Type of conflict (overlap, contradiction, priority)
  - TransactionId: Transaction where conflict occurred (if applicable)
  - DetectedAt: Timestamp of conflict detection
  - ResolutionStrategy: Suggested resolution approach
- **Consequences**: 
  - Conflict resolution process is triggered
  - Conflict is logged for analysis
  - Alert may be sent to administrators
  - Rule authors notified for resolution
- **Timing**: During rule evaluation or validation when conflicts are detected
- **Subscribers**: 
  - Conflict Resolution Service for automatic resolution
  - Alert Service for administrator notification
  - Analytics Service for conflict trend analysis
  - Rule Management UI for user notification

## Aggregates

### Rule Aggregate
- **Purpose**: Manages the complete lifecycle and behavior of business rules while maintaining consistency and enforcing business invariants
- **Root Entity**: Rule
- **Boundaries**: 
  - **Included**: Rule metadata, DSL content, status, version history, approval information, template references
  - **Excluded**: Rule execution results, performance metrics, external system data, evaluation contexts
- **Invariants**: 
  - Rule status transitions must follow defined workflow (AC-05: DRAFT → UNDER_REVIEW → APPROVED → ACTIVE)
  - DSL content must be syntactically and semantically valid (AC-02, AC-03)
  - Only approved rules can be activated (from approval policy)
  - Rule modifications must maintain version history (US-0009: view rule history)
  - Priority must be unique within category for active rules (from business rule invariants)
- **Commands**: 
  - CreateRule: Create new rule with validation (US-0001, US-0002)
  - UpdateRule: Modify rule content with version control (US-0004)
  - ValidateRule: Perform syntax and semantic validation (AC-03)
  - TestRule: Execute rule with sample data (US-0006)
  - SubmitForApproval: Submit rule for approval workflow (US-0008)
  - ApproveRule: Approve rule for activation
  - ActivateRule: Make rule available for evaluation
  - DeactivateRule: Remove rule from active evaluation
  - CloneRule: Create copy of existing rule (US-0010)

### Evaluation Context Aggregate
- **Purpose**: Encapsulates all data and state required for rule evaluation while ensuring consistency throughout the evaluation process
- **Root Entity**: EvaluationContext
- **Boundaries**: 
  - **Included**: Transaction data, customer data, rule execution state, evaluation results, performance metrics
  - **Excluded**: Rule definitions, approval workflow, template management, user interface state
- **Invariants**: 
  - Context must contain valid transaction and customer data
  - All applicable rules must be evaluated before completion
  - Conflicts must be resolved before final results
  - Evaluation state must be consistent throughout process
  - Performance constraints must be maintained (<500ms for 95% of requests)
- **Commands**: 
  - CreateEvaluationContext: Initialize context with transaction data
  - AddRulesToContext: Add applicable rules for evaluation
  - ExecuteRules: Process rules in priority order
  - ResolveConflicts: Apply conflict resolution strategies
  - CompleteEvaluation: Finalize results and metrics

### Template Aggregate
- **Purpose**: Manages reusable rule patterns and their application to create new rules efficiently
- **Root Entity**: RuleTemplate
- **Boundaries**: 
  - **Included**: Template structure, parameters, default values, usage instructions, version history
  - **Excluded**: Specific rule instances, approval workflows, evaluation results, user preferences
- **Invariants**: 
  - Template must generate valid DSL when applied (AC-06)
  - Parameters must have proper type definitions
  - Template versioning must maintain backward compatibility
  - Template categories must align with business domains
- **Commands**: 
  - CreateTemplate: Design new template with parameters
  - UpdateTemplate: Modify template structure
  - ApplyTemplate: Create rule from template with parameter values (US-0003)
  - ValidateTemplate: Verify template generates valid rules

## Business Processes

### Rule Creation Process
- **Goal**: Enable business users to create new rules efficiently with proper validation and guidance (US-0001, US-0002, US-0003)
- **Trigger**: Business user initiates rule creation
- **Steps**: 
  1. Select appropriate template or start from scratch (AC-01: template selection)
  2. Define rule conditions using DSL with auto-completion (AC-02: DSL rule definition)
  3. Configure rule actions and parameters with validation feedback
  4. Validate syntax and business logic within 2 seconds (AC-03: rule validation)
  5. Test rule with sample data within 3 seconds (AC-04: rule testing)
  6. Save rule in DRAFT status (AC-05: rule status management)
- **Participants**: Business Analyst, Rule Template Service, Validation Service, Testing Service
- **Outcomes**: 
  - **Success**: Rule created and ready for approval (US-0008)
  - **Validation Failure**: Rule returned for corrections with specific error messages
  - **Technical Error**: Error logged and user notified with recovery options

### Rule Approval Process
- **Goal**: Ensure rule changes are properly reviewed and approved before activation (US-0008)
- **Trigger**: Rule submitted for approval
- **Steps**: 
  1. Compliance review for regulatory requirements
  2. Business impact analysis and risk assessment
  3. Technical review for performance implications
  4. Conflict detection against existing active rules
  5. Stakeholder approval based on impact level
  6. Final approval and preparation for activation
- **Participants**: Rule Creator, Compliance Officer, Business Approver, Technical Reviewer
- **Outcomes**: 
  - **Approved**: Rule ready for activation
  - **Rejected**: Rule returned with feedback for improvement
  - **Conditional Approval**: Rule approved with specific conditions

### Rule Evaluation Process
- **Goal**: Evaluate applicable rules against transaction data to determine benefits and actions (AC-04)
- **Trigger**: Transaction requiring rule evaluation
- **Steps**: 
  1. Receive transaction and customer data
  2. Identify applicable rules based on criteria and effective dates
  3. Execute rules in priority order (CRITICAL → HIGH → MEDIUM → LOW)
  4. Detect and resolve conflicts using resolution strategies
  5. Calculate final results and benefits
  6. Return results with audit information within 500ms
- **Participants**: External System, Evaluation Engine, Conflict Resolution Service
- **Outcomes**: 
  - **Success**: Results calculated and returned within performance targets
  - **Performance Issue**: Results with degraded performance but within acceptable limits
  - **Failure**: Error response with fallback values and detailed error information

## Anti-Patterns and Forbidden Terms

### Forbidden Terms and Replacements
- **"Function"**: Use **"Rule"** or **"Business Rule"** instead
  - *Reason*: "Function" is too technical and doesn't convey business meaning
- **"Code"**: Use **"DSL Content"** or **"Rule Logic"** instead
  - *Reason*: "Code" implies programming rather than business rule definition
- **"Script"**: Use **"Rule Definition"** or **"Business Logic"** instead
  - *Reason*: "Script" suggests technical scripting rather than business rules
- **"Procedure"**: Use **"Business Process"** or **"Rule Workflow"** instead
  - *Reason*: "Procedure" is too generic and lacks business context
- **"System"**: Use **"Rules Engine"** or **"Platform"** instead when referring to the domain
  - *Reason*: "System" is too generic and doesn't specify the rules domain

### Anti-Patterns to Avoid
- **Technical Implementation Details in Business Conversations**: 
  - *Problem*: Discussing database schemas or technical architecture when talking about business rules
  - *Solution*: Focus on business capabilities and rule behavior
- **Generic Terms for Specific Concepts**: 
  - *Problem*: Using "data" instead of "transaction data" or "customer data"
  - *Solution*: Be specific about the type and context of data
- **Mixing Lifecycle Stages**: 
  - *Problem*: Confusing rule creation with rule evaluation or approval
  - *Solution*: Clearly identify which stage of the rule lifecycle you're discussing
- **Ambiguous Status References**: 
  - *Problem*: Saying "active" without specifying rule status vs. system status
  - *Solution*: Use full terms like "ACTIVE rule status" or "system operational status"

### Domain Boundary Clarifications
- **Rule vs. Rule Template**: 
  - *Rule*: Specific instance with actual business logic for execution
  - *Rule Template*: Reusable pattern for creating multiple similar rules
- **Rule Creation vs. Rule Evaluation**: 
  - *Rule Creation*: Process of defining and configuring new business rules
  - *Rule Evaluation*: Process of executing rules against transaction data
- **Validation vs. Testing**: 
  - *Validation*: Checking rule syntax and semantic correctness
  - *Testing*: Executing rule with sample data to verify behavior
- **Approval vs. Activation**: 
  - *Approval*: Business approval for rule to be used
  - *Activation*: Technical process of making rule available for evaluation

This ubiquitous language serves as the foundation for all communication within the Rules Engine domain, ensuring consistent understanding across business and technical stakeholders while maintaining complete traceability to the original PRD sources.