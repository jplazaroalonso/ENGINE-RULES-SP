# Rules Calculation Service

## Overview
The Rules Calculation Service is the high-performance core engine responsible for rule evaluation, conflict resolution, and result calculation. This service processes evaluation requests with sub-500ms response times and handles 1000+ TPS.

## Domain Model

### Core Entities

#### Calculation Engine
```go
type CalculationEngine struct {
    ID              string                `json:"id" gorm:"primaryKey"`
    Version         string                `json:"version"`
    Status          EngineStatus          `json:"status"`
    Configuration   EngineConfiguration   `json:"configuration" gorm:"embedded"`
    PerformanceMetrics PerformanceMetrics `json:"metrics" gorm:"embedded"`
    CreatedAt       time.Time             `json:"created_at"`
    UpdatedAt       time.Time             `json:"updated_at"`
}

type EngineStatus string
const (
    EngineStatusActive     EngineStatus = "ACTIVE"
    EngineStatusMaintenance EngineStatus = "MAINTENANCE"
    EngineStatusStopped    EngineStatus = "STOPPED"
)

type EngineConfiguration struct {
    MaxConcurrentEvaluations int           `json:"max_concurrent_evaluations"`
    CacheEnabled            bool          `json:"cache_enabled"`
    CacheTTL               time.Duration `json:"cache_ttl"`
    ConflictResolutionStrategy string    `json:"conflict_resolution_strategy"`
    PerformanceMode        string        `json:"performance_mode"`
}

type PerformanceMetrics struct {
    AverageResponseTime time.Duration `json:"average_response_time"`
    ThroughputTPS      float64       `json:"throughput_tps"`
    ErrorRate          float64       `json:"error_rate"`
    CacheHitRate       float64       `json:"cache_hit_rate"`
    LastUpdated        time.Time     `json:"last_updated"`
}
```

#### Calculation Context
```go
type CalculationContext struct {
    ID              string                 `json:"id" gorm:"primaryKey"`
    TransactionID   string                 `json:"transaction_id"`
    CustomerID      string                 `json:"customer_id"`
    SessionID       string                 `json:"session_id,omitempty"`
    EvaluationTime  time.Time              `json:"evaluation_time"`
    
    // Transaction Data
    Transaction     TransactionData        `json:"transaction" gorm:"embedded"`
    Customer        CustomerData           `json:"customer" gorm:"embedded"`
    Products        []ProductData          `json:"products" gorm:"serializer:json"`
    
    // Additional Context
    AdditionalData  map[string]interface{} `json:"additional_data" gorm:"serializer:json"`
    Metadata        ContextMetadata        `json:"metadata" gorm:"embedded"`
    
    CreatedAt       time.Time              `json:"created_at"`
}

type TransactionData struct {
    Amount          float64   `json:"amount"`
    Currency        string    `json:"currency"`
    Type            string    `json:"type"`
    Channel         string    `json:"channel"`
    Location        string    `json:"location,omitempty"`
    Timestamp       time.Time `json:"timestamp"`
}

type CustomerData struct {
    ID              string                 `json:"id"`
    Tier            string                 `json:"tier,omitempty"`
    Segment         string                 `json:"segment,omitempty"`
    LoyaltyPoints   int                    `json:"loyalty_points,omitempty"`
    Attributes      map[string]interface{} `json:"attributes,omitempty"`
}

type ProductData struct {
    ID              string  `json:"id"`
    SKU             string  `json:"sku"`
    Category        string  `json:"category"`
    Price           float64 `json:"price"`
    Quantity        int     `json:"quantity"`
    Attributes      map[string]interface{} `json:"attributes,omitempty"`
}
```

#### Rule Execution Result
```go
type RuleExecutionResult struct {
    ID              string                 `json:"id" gorm:"primaryKey"`
    ContextID       string                 `json:"context_id"`
    EngineVersion   string                 `json:"engine_version"`
    
    // Execution Details
    ExecutedRules   []RuleExecution        `json:"executed_rules" gorm:"serializer:json"`
    ConflictResults ConflictResolution     `json:"conflict_results" gorm:"embedded"`
    FinalResult     AggregatedResult       `json:"final_result" gorm:"embedded"`
    
    // Performance
    ExecutionTime   time.Duration          `json:"execution_time"`
    CacheHit        bool                   `json:"cache_hit"`
    
    // Metadata
    CreatedAt       time.Time              `json:"created_at"`
    ProcessedAt     time.Time              `json:"processed_at"`
}

type RuleExecution struct {
    RuleID          string        `json:"rule_id"`
    RuleName        string        `json:"rule_name"`
    RuleVersion     int           `json:"rule_version"`
    ExecutionOrder  int           `json:"execution_order"`
    Result          RuleResult    `json:"result"`
    ExecutionTime   time.Duration `json:"execution_time"`
    CacheUsed       bool          `json:"cache_used"`
}

type RuleResult struct {
    Success         bool                   `json:"success"`
    Applied         bool                   `json:"applied"`
    Value           interface{}            `json:"value,omitempty"`
    DiscountAmount  float64               `json:"discount_amount,omitempty"`
    PointsAwarded   int                   `json:"points_awarded,omitempty"`
    Reason          string                `json:"reason,omitempty"`
    AdditionalData  map[string]interface{} `json:"additional_data,omitempty"`
}

type ConflictResolution struct {
    ConflictsDetected bool               `json:"conflicts_detected"`
    ConflictCount     int               `json:"conflict_count"`
    ResolutionStrategy string           `json:"resolution_strategy"`
    ResolvedConflicts []ResolvedConflict `json:"resolved_conflicts,omitempty"`
}

type ResolvedConflict struct {
    ConflictingRules []string `json:"conflicting_rules"`
    WinningRule      string   `json:"winning_rule"`
    ResolutionReason string   `json:"resolution_reason"`
}

type AggregatedResult struct {
    TotalDiscount    float64 `json:"total_discount"`
    TotalPoints      int     `json:"total_points"`
    AppliedTaxes     float64 `json:"applied_taxes"`
    FinalAmount      float64 `json:"final_amount"`
    AppliedRuleCount int     `json:"applied_rule_count"`
}
```

## gRPC Service Definition

```protobuf
syntax = "proto3";

package rules_calculation;

service RulesCalculationService {
  // Core calculation methods
  rpc EvaluateRules(EvaluateRulesRequest) returns (EvaluateRulesResponse);
  rpc EvaluateRulesBatch(stream EvaluateRulesRequest) returns (stream EvaluateRulesResponse);
  
  // Performance and optimization
  rpc PreloadRules(PreloadRulesRequest) returns (PreloadRulesResponse);
  rpc ClearCache(ClearCacheRequest) returns (ClearCacheResponse);
  rpc GetPerformanceMetrics(GetPerformanceMetricsRequest) returns (GetPerformanceMetricsResponse);
  
  // Engine management
  rpc GetEngineStatus(GetEngineStatusRequest) returns (GetEngineStatusResponse);
  rpc UpdateEngineConfiguration(UpdateEngineConfigRequest) returns (UpdateEngineConfigResponse);
}

message EvaluateRulesRequest {
  string context_id = 1;
  string transaction_id = 2;
  string customer_id = 3;
  TransactionData transaction = 4;
  CustomerData customer = 5;
  repeated ProductData products = 6;
  map<string, string> additional_data = 7;
  repeated string rule_ids = 8; // Optional: specific rules to evaluate
  bool use_cache = 9;
}

message EvaluateRulesResponse {
  string context_id = 1;
  bool success = 2;
  string error_message = 3;
  repeated RuleExecution executed_rules = 4;
  ConflictResolution conflict_resolution = 5;
  AggregatedResult final_result = 6;
  int64 execution_time_ms = 7;
  bool cache_hit = 8;
}

message TransactionData {
  double amount = 1;
  string currency = 2;
  string type = 3;
  string channel = 4;
  string location = 5;
  int64 timestamp = 6;
}

message CustomerData {
  string id = 1;
  string tier = 2;
  string segment = 3;
  int32 loyalty_points = 4;
  map<string, string> attributes = 5;
}

message ProductData {
  string id = 1;
  string sku = 2;
  string category = 3;
  double price = 4;
  int32 quantity = 5;
  map<string, string> attributes = 6;
}
```

## Domain Services

### Calculation Engine Service
```go
type CalculationEngineService interface {
    EvaluateRules(ctx context.Context, req *EvaluationRequest) (*RuleExecutionResult, error)
    EvaluateRulesBatch(ctx context.Context, requests []*EvaluationRequest) ([]*RuleExecutionResult, error)
    PreloadRules(ctx context.Context, ruleIDs []string) error
    ClearCache(ctx context.Context, pattern string) error
    GetPerformanceMetrics(ctx context.Context) (*PerformanceMetrics, error)
    UpdateConfiguration(ctx context.Context, config *EngineConfiguration) error
}

type EvaluationRequest struct {
    ContextID      string
    TransactionID  string
    CustomerID     string
    Transaction    TransactionData
    Customer       CustomerData
    Products       []ProductData
    AdditionalData map[string]interface{}
    RuleIDs        []string // Optional: specific rules to evaluate
    UseCache       bool
}
```

### Conflict Resolution Service
```go
type ConflictResolutionService interface {
    DetectConflicts(ctx context.Context, results []RuleExecution) ([]Conflict, error)
    ResolveConflicts(ctx context.Context, conflicts []Conflict, strategy ConflictStrategy) (*ConflictResolution, error)
    GetResolutionStrategies(ctx context.Context) ([]ConflictStrategy, error)
}

type Conflict struct {
    RuleIDs          []string
    ConflictType     ConflictType
    ConflictSeverity ConflictSeverity
    Description      string
}

type ConflictType string
const (
    ConflictTypePriority     ConflictType = "PRIORITY"
    ConflictTypeExclusive    ConflictType = "EXCLUSIVE"
    ConflictTypeQuantitative ConflictType = "QUANTITATIVE"
)

type ConflictStrategy string
const (
    StrategyHighestPriority ConflictStrategy = "HIGHEST_PRIORITY"
    StrategyFirstWins       ConflictStrategy = "FIRST_WINS"
    StrategyLastWins        ConflictStrategy = "LAST_WINS"
    StrategySum             ConflictStrategy = "SUM"
    StrategyMaximum         ConflictStrategy = "MAXIMUM"
    StrategyMinimum         ConflictStrategy = "MINIMUM"
)
```

### Cache Service
```go
type CacheService interface {
    GetRuleExecutionResult(ctx context.Context, key string) (*RuleExecutionResult, error)
    SetRuleExecutionResult(ctx context.Context, key string, result *RuleExecutionResult, ttl time.Duration) error
    GetCompiledRule(ctx context.Context, ruleID string) (*CompiledRule, error)
    SetCompiledRule(ctx context.Context, ruleID string, rule *CompiledRule, ttl time.Duration) error
    InvalidateCache(ctx context.Context, pattern string) error
    GetCacheStats(ctx context.Context) (*CacheStats, error)
}

type CompiledRule struct {
    ID              string
    Version         int
    CompiledCode    []byte
    Dependencies    []string
    CompilationTime time.Time
}

type CacheStats struct {
    HitRate         float64
    MissRate        float64
    EvictionRate    float64
    TotalRequests   int64
    MemoryUsage     int64
}
```

## Repository Interfaces

```go
type CalculationContextRepository interface {
    Create(ctx context.Context, context *CalculationContext) error
    GetByID(ctx context.Context, id string) (*CalculationContext, error)
    GetByTransactionID(ctx context.Context, transactionID string) (*CalculationContext, error)
    Update(ctx context.Context, context *CalculationContext) error
    Delete(ctx context.Context, id string) error
}

type RuleExecutionResultRepository interface {
    Create(ctx context.Context, result *RuleExecutionResult) error
    GetByContextID(ctx context.Context, contextID string) (*RuleExecutionResult, error)
    GetByTransactionID(ctx context.Context, transactionID string) ([]*RuleExecutionResult, error)
    GetPerformanceStats(ctx context.Context, from, to time.Time) (*PerformanceStats, error)
    Archive(ctx context.Context, olderThan time.Time) error
}

type EngineConfigurationRepository interface {
    GetConfiguration(ctx context.Context) (*EngineConfiguration, error)
    UpdateConfiguration(ctx context.Context, config *EngineConfiguration) error
    GetPerformanceMetrics(ctx context.Context) (*PerformanceMetrics, error)
    UpdatePerformanceMetrics(ctx context.Context, metrics *PerformanceMetrics) error
}
```

## Domain Events

```go
type RuleEvaluationStarted struct {
    ContextID     string    `json:"context_id"`
    TransactionID string    `json:"transaction_id"`
    CustomerID    string    `json:"customer_id"`
    RuleCount     int       `json:"rule_count"`
    StartedAt     time.Time `json:"started_at"`
}

type RuleEvaluationCompleted struct {
    ContextID       string        `json:"context_id"`
    TransactionID   string        `json:"transaction_id"`
    Success         bool          `json:"success"`
    ExecutionTime   time.Duration `json:"execution_time"`
    RulesApplied    int           `json:"rules_applied"`
    ConflictsFound  int           `json:"conflicts_found"`
    CompletedAt     time.Time     `json:"completed_at"`
    FinalResult     AggregatedResult `json:"final_result"`
}

type RuleConflictDetected struct {
    ContextID       string   `json:"context_id"`
    ConflictingRules []string `json:"conflicting_rules"`
    ConflictType    string   `json:"conflict_type"`
    ResolutionUsed  string   `json:"resolution_used"`
    DetectedAt      time.Time `json:"detected_at"`
}

type PerformanceThresholdExceeded struct {
    MetricType      string    `json:"metric_type"`
    CurrentValue    float64   `json:"current_value"`
    ThresholdValue  float64   `json:"threshold_value"`
    EngineVersion   string    `json:"engine_version"`
    DetectedAt      time.Time `json:"detected_at"`
}
```

## Implementation Tasks

### Phase 1: Core Engine Foundation (3-4 days)
1. **Project Setup and Dependencies**
   - Initialize Go module with high-performance dependencies
   - Setup Redis client for caching (go-redis)
   - Configure PostgreSQL with connection pooling
   - Setup monitoring with Prometheus metrics

2. **Domain Model Implementation**
   - Implement CalculationEngine, CalculationContext entities
   - Create RuleExecutionResult and related value objects
   - Implement performance metrics tracking
   - Add domain validation and business rules

3. **Repository Layer**
   - Implement PostgreSQL repositories with optimized queries
   - Add database indexes for performance-critical queries
   - Implement Redis cache layer for rule results
   - Add repository performance monitoring

### Phase 2: Core Calculation Engine (4-5 days)
1. **Rule Evaluation Engine**
   - Implement high-performance rule evaluation logic
   - Create rule compilation and caching mechanism
   - Add parallel rule execution for independent rules
   - Implement execution context management

2. **DSL Interpreter Integration**
   - Integrate with DSL parsing and compilation
   - Implement rule expression evaluation
   - Add variable binding and context injection
   - Create rule result extraction and mapping

3. **Performance Optimization**
   - Implement rule result caching with intelligent invalidation
   - Add compiled rule caching for repeated evaluations
   - Optimize database queries with prepared statements
   - Implement connection pooling and resource management

### Phase 3: Conflict Resolution System (3-4 days)
1. **Conflict Detection**
   - Implement conflict detection algorithms
   - Create rule priority and exclusivity checking
   - Add quantitative conflict detection (overlapping discounts)
   - Implement conflict severity classification

2. **Resolution Strategies**
   - Implement priority-based conflict resolution
   - Add quantitative resolution strategies (sum, max, min)
   - Create custom resolution rule support
   - Add resolution audit trail and logging

3. **Resolution Engine**
   - Create conflict resolution orchestrator
   - Implement strategy pattern for different resolution types
   - Add resolution result tracking and reporting
   - Implement resolution performance optimization

### Phase 4: gRPC API Implementation (3-4 days)
1. **Protocol Buffer Definitions**
   - Define comprehensive .proto files for all operations
   - Generate Go code from protobuf definitions
   - Implement streaming support for batch operations
   - Add performance-optimized message serialization

2. **gRPC Service Implementation**
   - Implement EvaluateRules with sub-500ms target
   - Add EvaluateRulesBatch for high-throughput operations
   - Implement engine management operations
   - Add performance metrics and health check endpoints

3. **Performance Optimization**
   - Implement gRPC connection pooling
   - Add request batching and optimization
   - Implement response compression
   - Add load balancing and circuit breaker patterns

### Phase 5: Caching and Performance (3-4 days)
1. **Multi-Level Caching**
   - Implement L1 cache (in-memory) for frequently used rules
   - Add L2 cache (Redis) for rule results and compiled rules
   - Create intelligent cache warming strategies
   - Implement cache invalidation on rule updates

2. **Performance Monitoring**
   - Add comprehensive metrics collection (response time, throughput)
   - Implement performance alerting and thresholds
   - Create performance dashboard and reporting
   - Add execution profiling and bottleneck detection

3. **Optimization Features**
   - Implement rule preloading for frequently used rules
   - Add query optimization and execution plan analysis
   - Create dynamic performance tuning
   - Implement adaptive caching based on usage patterns

### Phase 6: Event Integration and Messaging (2-3 days)
1. **NATS Integration**
   - Setup NATS client with high-performance configuration
   - Implement event publishing for all domain events
   - Add event subscription for rule updates
   - Create retry and dead letter queue handling

2. **Event Handlers**
   - Publish RuleEvaluationStarted events
   - Publish RuleEvaluationCompleted with detailed metrics
   - Handle RuleUpdated events for cache invalidation
   - Implement PerformanceThresholdExceeded alerting

### Phase 7: Testing and Quality Assurance (4-5 days)
1. **Unit Testing**
   - Write comprehensive unit tests (90% coverage target)
   - Test calculation engine with various rule combinations
   - Test conflict resolution with complex scenarios
   - Test caching mechanisms and invalidation

2. **Performance Testing**
   - Load test with 1000+ TPS target
   - Benchmark response time (<500ms target)
   - Test concurrent rule evaluation scenarios
   - Validate memory usage and garbage collection

3. **Integration Testing**
   - Test end-to-end rule evaluation flows
   - Test event publishing and handling
   - Test cache behavior under load
   - Validate database performance under stress

### Phase 8: Production Readiness (3-4 days)
1. **Monitoring and Observability**
   - Implement comprehensive Prometheus metrics
   - Add distributed tracing with correlation IDs
   - Create performance dashboards and alerts
   - Implement health checks for Kubernetes

2. **Security and Reliability**
   - Add input validation and sanitization
   - Implement rate limiting and circuit breakers
   - Add graceful shutdown and resource cleanup
   - Implement backup and recovery procedures

3. **Documentation and Deployment**
   - Complete gRPC API documentation with examples
   - Create performance tuning guides
   - Add troubleshooting runbooks
   - Create deployment configurations for production

## Estimated Development Time: 25-32 days

## Performance Targets
- **Response Time**: <500ms (95th percentile)
- **Throughput**: 1000+ TPS sustained, 3000+ TPS peak
- **Accuracy**: 99.99% calculation accuracy
- **Availability**: 99.9% uptime
- **Cache Hit Rate**: >80%
