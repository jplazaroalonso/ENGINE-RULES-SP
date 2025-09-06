# NATS Messaging Patterns and Event Schemas

## Overview
This document defines the complete messaging patterns, event schemas, and communication protocols used across the Rules Engine microservices using NATS as the messaging backbone.

## Event Schema Standards

### Base Event Structure
All events follow a standardized structure to ensure consistency and enable proper event handling across services.

```go
// BaseEvent represents the common structure for all domain events
type BaseEvent struct {
    // Event Identity
    ID              string                 `json:"id"`
    Type            string                 `json:"type"`
    Source          string                 `json:"source"`
    Subject         string                 `json:"subject"`
    
    // Temporal Information
    Time            time.Time              `json:"time"`
    Version         string                 `json:"version"`
    
    // Event Correlation
    CorrelationID   string                 `json:"correlation_id,omitempty"`
    CausationID     string                 `json:"causation_id,omitempty"`
    
    // Event Chain
    EventChainID    string                 `json:"event_chain_id,omitempty"`
    SequenceNumber  int64                  `json:"sequence_number,omitempty"`
    
    // Metadata
    Metadata        EventMetadata          `json:"metadata,omitempty"`
    
    // Payload
    Data            interface{}            `json:"data"`
}

type EventMetadata struct {
    UserID          string                 `json:"user_id,omitempty"`
    SessionID       string                 `json:"session_id,omitempty"`
    RequestID       string                 `json:"request_id,omitempty"`
    Source          string                 `json:"source"`
    Environment     string                 `json:"environment"`
    Region          string                 `json:"region,omitempty"`
    
    // Technical Metadata
    SchemaVersion   string                 `json:"schema_version"`
    ContentType     string                 `json:"content_type"`
    Encoding        string                 `json:"encoding"`
    
    // Retry Information
    RetryCount      int                    `json:"retry_count,omitempty"`
    MaxRetries      int                    `json:"max_retries,omitempty"`
    
    // Processing Hints
    Priority        EventPriority          `json:"priority,omitempty"`
    TTL             time.Duration          `json:"ttl,omitempty"`
    ProcessAfter    *time.Time             `json:"process_after,omitempty"`
}

type EventPriority string
const (
    PriorityLow      EventPriority = "LOW"
    PriorityNormal   EventPriority = "NORMAL"
    PriorityHigh     EventPriority = "HIGH"
    PriorityCritical EventPriority = "CRITICAL"
)
```

## Domain Event Schemas

### Rules Management Events

#### RuleCreated
```go
type RuleCreatedEvent struct {
    BaseEvent
    Data RuleCreatedData `json:"data"`
}

type RuleCreatedData struct {
    RuleID          string                 `json:"rule_id"`
    Name            string                 `json:"name"`
    Description     string                 `json:"description"`
    RuleType        string                 `json:"rule_type"`
    Status          string                 `json:"status"`
    Priority        string                 `json:"priority"`
    DSLContent      string                 `json:"dsl_content"`
    CreatedBy       string                 `json:"created_by"`
    CreatedAt       time.Time              `json:"created_at"`
    Version         int                    `json:"version"`
    TemplateID      *string                `json:"template_id,omitempty"`
    Tags            []string               `json:"tags,omitempty"`
    Dependencies    []string               `json:"dependencies,omitempty"`
    AffectedServices []string              `json:"affected_services,omitempty"`
}

// Subject: rules.domain.rule.created
// Consumers: rules-calculation, rules-evaluation, promotions, loyalty, coupons, taxes, payments
```

#### RuleUpdated
```go
type RuleUpdatedEvent struct {
    BaseEvent
    Data RuleUpdatedData `json:"data"`
}

type RuleUpdatedData struct {
    RuleID           string                 `json:"rule_id"`
    Name             string                 `json:"name"`
    Changes          map[string]interface{} `json:"changes"`
    PreviousVersion  int                    `json:"previous_version"`
    NewVersion       int                    `json:"new_version"`
    UpdatedBy        string                 `json:"updated_by"`
    UpdatedAt        time.Time              `json:"updated_at"`
    ChangeReason     string                 `json:"change_reason,omitempty"`
    BackwardCompatible bool                 `json:"backward_compatible"`
    ImpactAssessment []ImpactArea           `json:"impact_assessment"`
}

type ImpactArea struct {
    Service         string   `json:"service"`
    Components      []string `json:"components"`
    SeverityLevel   string   `json:"severity_level"`
    ActionRequired  bool     `json:"action_required"`
}

// Subject: rules.domain.rule.updated
// Consumers: rules-calculation, rules-evaluation, cache-invalidation-service
```

#### RuleActivated / RuleDeactivated
```go
type RuleStatusChangedEvent struct {
    BaseEvent
    Data RuleStatusChangedData `json:"data"`
}

type RuleStatusChangedData struct {
    RuleID          string    `json:"rule_id"`
    Name            string    `json:"name"`
    PreviousStatus  string    `json:"previous_status"`
    NewStatus       string    `json:"new_status"`
    StatusChangedBy string    `json:"status_changed_by"`
    StatusChangedAt time.Time `json:"status_changed_at"`
    EffectiveAt     time.Time `json:"effective_at"`
    Reason          string    `json:"reason,omitempty"`
    AutomatedChange bool      `json:"automated_change"`
}

// Subject: rules.domain.rule.activated / rules.domain.rule.deactivated
// Consumers: rules-calculation, rules-evaluation, analytics-service
```

### Rules Evaluation Events

#### EvaluationStarted
```go
type EvaluationStartedEvent struct {
    BaseEvent
    Data EvaluationStartedData `json:"data"`
}

type EvaluationStartedData struct {
    ContextID       string                 `json:"context_id"`
    TransactionID   string                 `json:"transaction_id"`
    CustomerID      string                 `json:"customer_id"`
    EvaluationType  string                 `json:"evaluation_type"`
    RuleCount       int                    `json:"rule_count"`
    StartedAt       time.Time              `json:"started_at"`
    ExpectedDuration time.Duration         `json:"expected_duration,omitempty"`
    Priority        string                 `json:"priority"`
    RequestedServices []string             `json:"requested_services"`
    EvaluationConfig map[string]interface{} `json:"evaluation_config"`
}

// Subject: rules.domain.evaluation.started
// Consumers: monitoring-service, analytics-service, timeout-manager
```

#### EvaluationCompleted
```go
type EvaluationCompletedEvent struct {
    BaseEvent
    Data EvaluationCompletedData `json:"data"`
}

type EvaluationCompletedData struct {
    ContextID         string              `json:"context_id"`
    TransactionID     string              `json:"transaction_id"`
    CustomerID        string              `json:"customer_id"`
    Success           bool                `json:"success"`
    ExecutionTime     time.Duration       `json:"execution_time"`
    RulesEvaluated    int                 `json:"rules_evaluated"`
    RulesApplied      int                 `json:"rules_applied"`
    ConflictsFound    int                 `json:"conflicts_found"`
    CompletedAt       time.Time           `json:"completed_at"`
    FinalResult       AggregatedResult    `json:"final_result"`
    ServiceResults    []ServiceResult     `json:"service_results"`
    PerformanceMetrics PerformanceMetrics `json:"performance_metrics"`
    ErrorMessage      string              `json:"error_message,omitempty"`
}

type AggregatedResult struct {
    TotalDiscount    float64 `json:"total_discount"`
    TotalPoints      int     `json:"total_points"`
    AppliedTaxes     float64 `json:"applied_taxes"`
    AppliedFees      float64 `json:"applied_fees"`
    FinalAmount      float64 `json:"final_amount"`
    AppliedRuleCount int     `json:"applied_rule_count"`
    Savings          float64 `json:"savings"`
}

type ServiceResult struct {
    ServiceName    string        `json:"service_name"`
    ExecutionTime  time.Duration `json:"execution_time"`
    Success        bool          `json:"success"`
    RulesApplied   int           `json:"rules_applied"`
    Result         interface{}   `json:"result"`
    ErrorMessage   string        `json:"error_message,omitempty"`
}

type PerformanceMetrics struct {
    TotalExecutionTime   time.Duration `json:"total_execution_time"`
    DatabaseQueryTime    time.Duration `json:"database_query_time"`
    ExternalServiceTime  time.Duration `json:"external_service_time"`
    CacheHitRate         float64       `json:"cache_hit_rate"`
    MemoryUsage          int64         `json:"memory_usage"`
    CPUUsage             float64       `json:"cpu_usage"`
}

// Subject: rules.domain.evaluation.completed
// Consumers: analytics-service, monitoring-service, customer-service, billing-service
```

### Business Domain Events

#### Promotion Events
```go
type PromotionAppliedEvent struct {
    BaseEvent
    Data PromotionAppliedData `json:"data"`
}

type PromotionAppliedData struct {
    PromotionID      string    `json:"promotion_id"`
    PromotionName    string    `json:"promotion_name"`
    CampaignID       string    `json:"campaign_id"`
    CustomerID       string    `json:"customer_id"`
    TransactionID    string    `json:"transaction_id"`
    DiscountAmount   float64   `json:"discount_amount"`
    DiscountType     string    `json:"discount_type"`
    OriginalAmount   float64   `json:"original_amount"`
    FinalAmount      float64   `json:"final_amount"`
    AppliedAt        time.Time `json:"applied_at"`
    ValidUntil       *time.Time `json:"valid_until,omitempty"`
    UsageCount       int       `json:"usage_count"`
    MaxUsage         *int      `json:"max_usage,omitempty"`
}

// Subject: rules.domain.promotion.applied
// Consumers: analytics-service, customer-service, campaign-manager
```

#### Loyalty Events
```go
type PointsAwardedEvent struct {
    BaseEvent
    Data PointsAwardedData `json:"data"`
}

type PointsAwardedData struct {
    CustomerID       string    `json:"customer_id"`
    AccountID        string    `json:"account_id"`
    TransactionID    string    `json:"transaction_id"`
    PointsAwarded    int64     `json:"points_awarded"`
    EarningRuleID    string    `json:"earning_rule_id"`
    EarningRate      float64   `json:"earning_rate"`
    BaseAmount       float64   `json:"base_amount"`
    Multiplier       float64   `json:"multiplier"`
    BonusPoints      int64     `json:"bonus_points"`
    NewBalance       int64     `json:"new_balance"`
    AwardedAt        time.Time `json:"awarded_at"`
    ExpiresAt        *time.Time `json:"expires_at,omitempty"`
    Source           string    `json:"source"`
    SourceID         string    `json:"source_id"`
}

// Subject: rules.domain.loyalty.points.awarded
// Consumers: customer-service, analytics-service, notification-service
```

#### Coupon Events
```go
type CouponRedeemedEvent struct {
    BaseEvent
    Data CouponRedeemedData `json:"data"`
}

type CouponRedeemedData struct {
    CouponID         string          `json:"coupon_id"`
    CouponCode       string          `json:"coupon_code"`
    CustomerID       string          `json:"customer_id"`
    TransactionID    string          `json:"transaction_id"`
    DiscountAmount   float64         `json:"discount_amount"`
    OriginalAmount   float64         `json:"original_amount"`
    FinalAmount      float64         `json:"final_amount"`
    RedeemedAt       time.Time       `json:"redeemed_at"`
    UsageCount       int             `json:"usage_count"`
    RemainingUses    *int            `json:"remaining_uses,omitempty"`
    FraudScore       float64         `json:"fraud_score"`
    SecurityChecks   []SecurityCheck `json:"security_checks"`
}

type SecurityCheck struct {
    CheckType   string    `json:"check_type"`
    Result      string    `json:"result"`
    Score       float64   `json:"score"`
    Details     string    `json:"details,omitempty"`
    CheckedAt   time.Time `json:"checked_at"`
}

// Subject: rules.domain.coupon.redeemed
// Consumers: analytics-service, fraud-detection-service, campaign-manager
```

## Integration Event Schemas

### External System Events

#### ExternalServiceCallInitiated
```go
type ExternalServiceCallEvent struct {
    BaseEvent
    Data ExternalServiceCallData `json:"data"`
}

type ExternalServiceCallData struct {
    ServiceName      string                 `json:"service_name"`
    Endpoint         string                 `json:"endpoint"`
    Method           string                 `json:"method"`
    RequestID        string                 `json:"request_id"`
    InitiatedAt      time.Time              `json:"initiated_at"`
    TimeoutDuration  time.Duration          `json:"timeout_duration"`
    RetryPolicy      RetryPolicy            `json:"retry_policy"`
    RequestSize      int64                  `json:"request_size,omitempty"`
    RequestHeaders   map[string]string      `json:"request_headers,omitempty"`
}

type RetryPolicy struct {
    MaxRetries    int           `json:"max_retries"`
    InitialDelay  time.Duration `json:"initial_delay"`
    MaxDelay      time.Duration `json:"max_delay"`
    Multiplier    float64       `json:"multiplier"`
    RetryOn       []string      `json:"retry_on"`
}

// Subject: rules.integration.external.call.initiated
// Consumers: monitoring-service, timeout-manager, rate-limiter
```

#### WebhookReceived
```go
type WebhookReceivedEvent struct {
    BaseEvent
    Data WebhookReceivedData `json:"data"`
}

type WebhookReceivedData struct {
    WebhookID       string            `json:"webhook_id"`
    Source          string            `json:"source"`
    EventType       string            `json:"event_type"`
    ReceivedAt      time.Time         `json:"received_at"`
    Headers         map[string]string `json:"headers"`
    PayloadSize     int64             `json:"payload_size"`
    Signature       string            `json:"signature,omitempty"`
    VerificationStatus string         `json:"verification_status"`
    ProcessingStatus string           `json:"processing_status"`
    Payload         interface{}       `json:"payload"`
}

// Subject: rules.integration.webhook.received
// Consumers: webhook-processor, security-validator, analytics-service
```

## Message Routing and Filtering

### Subject Hierarchy
```
rules.>
├── rules.domain.>                    # Domain events from business logic
│   ├── rules.domain.rule.>          # Rule lifecycle events
│   │   ├── rules.domain.rule.created
│   │   ├── rules.domain.rule.updated
│   │   ├── rules.domain.rule.activated
│   │   ├── rules.domain.rule.deactivated
│   │   └── rules.domain.rule.deleted
│   ├── rules.domain.evaluation.>    # Evaluation lifecycle events
│   │   ├── rules.domain.evaluation.started
│   │   ├── rules.domain.evaluation.completed
│   │   ├── rules.domain.evaluation.failed
│   │   └── rules.domain.evaluation.timeout
│   ├── rules.domain.promotion.>     # Promotion domain events
│   │   ├── rules.domain.promotion.applied
│   │   ├── rules.domain.promotion.expired
│   │   └── rules.domain.promotion.conflict
│   ├── rules.domain.loyalty.>       # Loyalty domain events
│   │   ├── rules.domain.loyalty.points.awarded
│   │   ├── rules.domain.loyalty.points.redeemed
│   │   ├── rules.domain.loyalty.tier.upgraded
│   │   └── rules.domain.loyalty.tier.downgraded
│   ├── rules.domain.coupon.>        # Coupon domain events
│   │   ├── rules.domain.coupon.created
│   │   ├── rules.domain.coupon.redeemed
│   │   ├── rules.domain.coupon.expired
│   │   └── rules.domain.coupon.fraud.detected
│   ├── rules.domain.tax.>           # Tax calculation events
│   │   ├── rules.domain.tax.calculated
│   │   ├── rules.domain.tax.jurisdiction.updated
│   │   └── rules.domain.tax.rate.changed
│   └── rules.domain.payment.>       # Payment rule events
│       ├── rules.domain.payment.route.selected
│       ├── rules.domain.payment.fraud.detected
│       └── rules.domain.payment.optimization.applied
├── rules.integration.>              # Integration and external events
│   ├── rules.integration.external.> # External system interactions
│   │   ├── rules.integration.external.call.initiated
│   │   ├── rules.integration.external.call.completed
│   │   ├── rules.integration.external.call.failed
│   │   └── rules.integration.external.call.timeout
│   ├── rules.integration.webhook.>  # Webhook events
│   │   ├── rules.integration.webhook.received
│   │   ├── rules.integration.webhook.processed
│   │   └── rules.integration.webhook.failed
│   └── rules.integration.sync.>     # Data synchronization events
│       ├── rules.integration.sync.started
│       ├── rules.integration.sync.completed
│       └── rules.integration.sync.failed
├── rules.command.>                  # Command requests (CQRS pattern)
│   ├── rules.command.rule.>         # Rule commands
│   │   ├── rules.command.rule.create
│   │   ├── rules.command.rule.update
│   │   ├── rules.command.rule.activate
│   │   └── rules.command.rule.deactivate
│   ├── rules.command.evaluation.>   # Evaluation commands
│   │   ├── rules.command.evaluation.execute
│   │   ├── rules.command.evaluation.cancel
│   │   └── rules.command.evaluation.retry
│   └── rules.command.cache.>        # Cache management commands
│       ├── rules.command.cache.invalidate
│       ├── rules.command.cache.warm
│       └── rules.command.cache.clear
└── rules.query.>                    # Query requests (CQRS pattern)
    ├── rules.query.rule.>           # Rule queries
    │   ├── rules.query.rule.get
    │   ├── rules.query.rule.list
    │   └── rules.query.rule.search
    ├── rules.query.metrics.>        # Metrics queries
    │   ├── rules.query.metrics.performance
    │   ├── rules.query.metrics.usage
    │   └── rules.query.metrics.health
    └── rules.query.audit.>          # Audit queries
        ├── rules.query.audit.trail
        ├── rules.query.audit.compliance
        └── rules.query.audit.changes
```

### Event Filtering Patterns

#### Service-Specific Subscriptions
```go
// Rules Calculation Service Subscriptions
var calculationServiceSubscriptions = []string{
    "rules.domain.rule.>",           // All rule lifecycle events
    "rules.command.evaluation.>",    // Evaluation commands
    "rules.command.cache.>",         // Cache management
}

// Analytics Service Subscriptions  
var analyticsServiceSubscriptions = []string{
    "rules.domain.evaluation.completed", // Completed evaluations only
    "rules.domain.*.applied",            // All applied events (promotion, coupon, etc.)
    "rules.integration.external.call.>", // External service performance
}

// Monitoring Service Subscriptions
var monitoringServiceSubscriptions = []string{
    "rules.domain.evaluation.>",         // All evaluation events
    "rules.integration.external.call.>", // External service calls
    "rules.domain.*.failed",             // All failure events
}
```

#### Content-Based Filtering
```go
// Example: Filter events by priority
func HighPriorityEventFilter(event *BaseEvent) bool {
    return event.Metadata.Priority == PriorityHigh || 
           event.Metadata.Priority == PriorityCritical
}

// Example: Filter events by customer tier
func VIPCustomerEventFilter(event *BaseEvent) bool {
    if customerData, ok := event.Data.(map[string]interface{}); ok {
        if tier, exists := customerData["customer_tier"]; exists {
            return tier == "GOLD" || tier == "PLATINUM"
        }
    }
    return false
}

// Example: Filter events by service impact
func ServiceImpactFilter(serviceName string) func(*BaseEvent) bool {
    return func(event *BaseEvent) bool {
        if metadata, ok := event.Metadata.(EventMetadata); ok {
            for _, affected := range metadata.AffectedServices {
                if affected == serviceName {
                    return true
                }
            }
        }
        return false
    }
}
```

## Message Delivery Guarantees

### Delivery Semantics
- **At-least-once delivery**: Default for all domain events
- **Exactly-once processing**: For financial transactions and critical updates
- **At-most-once delivery**: For monitoring and analytics events (acceptable data loss)

### Durability Configuration
```yaml
# JetStream Stream Configuration for Domain Events
domain_events_stream:
  name: "RULES_DOMAIN_EVENTS"
  subjects: ["rules.domain.>"]
  storage: file
  retention: limits
  max_age: 30d
  max_bytes: 10GB
  max_msgs: 1000000
  replicas: 3
  duplicates: 2m

# JetStream Stream Configuration for Integration Events  
integration_events_stream:
  name: "RULES_INTEGRATION_EVENTS"
  subjects: ["rules.integration.>"]
  storage: file
  retention: limits
  max_age: 7d
  max_bytes: 5GB
  max_msgs: 500000
  replicas: 3
  duplicates: 1m

# JetStream Stream Configuration for Commands/Queries
cqrs_stream:
  name: "RULES_CQRS_EVENTS"
  subjects: ["rules.command.>", "rules.query.>"]
  storage: memory
  retention: limits
  max_age: 1h
  max_bytes: 1GB
  max_msgs: 100000
  replicas: 1
  duplicates: 30s
```

### Consumer Configuration Examples
```yaml
# Durable Consumer for Rules Calculation Service
rules_calculation_consumer:
  durable_name: "rules-calculation-service"
  deliver_policy: all
  ack_policy: explicit
  ack_wait: 30s
  max_deliver: 3
  max_ack_pending: 1000
  replay_policy: instant
  filter_subject: "rules.domain.rule.>, rules.command.evaluation.>"

# Ephemeral Consumer for Real-time Monitoring
monitoring_consumer:
  deliver_policy: new
  ack_policy: none
  replay_policy: instant
  filter_subject: "rules.domain.evaluation.>, rules.integration.external.call.>"
```

## Implementation Tasks

### Phase 1: Event Schema Definition (2-3 days)
1. **Define Base Event Structures**
   - Implement BaseEvent and EventMetadata types
   - Create event validation framework
   - Add event serialization/deserialization
   - Implement event versioning support

2. **Domain Event Schemas**
   - Define all rule management event schemas
   - Create evaluation lifecycle event schemas
   - Implement business domain event schemas
   - Add integration event schemas

### Phase 2: Subject Hierarchy and Routing (2-3 days)
1. **Subject Design**
   - Implement hierarchical subject structure
   - Create subject naming conventions
   - Add subject validation and enforcement
   - Implement subject-based security policies

2. **Message Routing**
   - Create content-based filtering framework
   - Implement service-specific subscription patterns
   - Add dynamic routing capabilities
   - Create routing performance optimization

### Phase 3: Stream and Consumer Management (3-4 days)
1. **Stream Configuration**
   - Implement JetStream stream setup
   - Create stream management utilities
   - Add stream monitoring and alerting
   - Implement stream backup and recovery

2. **Consumer Management**
   - Create consumer configuration templates
   - Implement consumer lifecycle management
   - Add consumer health monitoring
   - Create consumer scaling policies

### Phase 4: Event Processing Framework (3-4 days)
1. **Event Publishing**
   - Implement event publishing utilities
   - Create batch publishing capabilities
   - Add publishing retry and error handling
   - Implement publishing performance optimization

2. **Event Consumption**
   - Create event consumption framework
   - Implement event handler registration
   - Add event processing middleware
   - Create event replay capabilities

### Phase 5: Testing and Documentation (2-3 days)
1. **Testing Framework**
   - Create event schema validation tests
   - Implement integration testing utilities
   - Add performance testing for message throughput
   - Create chaos testing for resilience

2. **Documentation and Tooling**
   - Create event schema documentation
   - Implement event schema registry
   - Add development and debugging tools
   - Create operational runbooks

## Estimated Development Time: 12-17 days
