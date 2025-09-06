# Rules Evaluation Service

## Overview
The Rules Evaluation Service acts as the external API gateway for rule evaluation, orchestrating requests across multiple domain services and aggregating results. This service handles high-throughput external requests and provides optimized response times through caching and load balancing.

## Domain Model

### Core Entities

#### Evaluation Request
```go
type EvaluationRequest struct {
    ID                string                 `json:"id" gorm:"primaryKey"`
    ExternalRequestID string                 `json:"external_request_id" gorm:"uniqueIndex"`
    CustomerID        string                 `json:"customer_id" gorm:"index"`
    SessionID         string                 `json:"session_id,omitempty"`
    
    // Request Data
    TransactionData   TransactionData        `json:"transaction_data" gorm:"embedded"`
    CustomerData      CustomerData           `json:"customer_data" gorm:"embedded"`
    Products          []ProductData          `json:"products" gorm:"serializer:json"`
    Context           map[string]interface{} `json:"context" gorm:"serializer:json"`
    
    // Configuration
    RuleFilters       []string               `json:"rule_filters,omitempty" gorm:"serializer:json"`
    ResponseFormat    ResponseFormat         `json:"response_format" gorm:"default:'STANDARD'"`
    IncludeDetails    bool                   `json:"include_details" gorm:"default:false"`
    
    // Status and Timing
    Status            RequestStatus          `json:"status" gorm:"default:'PENDING'"`
    Priority          RequestPriority        `json:"priority" gorm:"default:'NORMAL'"`
    ReceivedAt        time.Time              `json:"received_at"`
    ProcessedAt       *time.Time             `json:"processed_at,omitempty"`
    CompletedAt       *time.Time             `json:"completed_at,omitempty"`
    
    // Results
    Response          *EvaluationResponse    `json:"response,omitempty" gorm:"foreignKey:RequestID"`
}

type RequestStatus string
const (
    StatusPending    RequestStatus = "PENDING"
    StatusProcessing RequestStatus = "PROCESSING"
    StatusCompleted  RequestStatus = "COMPLETED"
    StatusFailed     RequestStatus = "FAILED"
    StatusTimeout    RequestStatus = "TIMEOUT"
)

type RequestPriority string
const (
    PriorityLow    RequestPriority = "LOW"
    PriorityNormal RequestPriority = "NORMAL"
    PriorityHigh   RequestPriority = "HIGH"
    PriorityUrgent RequestPriority = "URGENT"
)

type ResponseFormat string
const (
    ResponseFormatStandard  ResponseFormat = "STANDARD"
    ResponseFormatDetailed  ResponseFormat = "DETAILED"
    ResponseFormatSummary   ResponseFormat = "SUMMARY"
    ResponseFormatCustom    ResponseFormat = "CUSTOM"
)
```

#### Evaluation Response
```go
type EvaluationResponse struct {
    ID                string                 `json:"id" gorm:"primaryKey"`
    RequestID         string                 `json:"request_id" gorm:"index"`
    ExternalRequestID string                 `json:"external_request_id"`
    
    // Results
    Success           bool                   `json:"success"`
    ErrorCode         string                 `json:"error_code,omitempty"`
    ErrorMessage      string                 `json:"error_message,omitempty"`
    
    // Aggregated Results
    TotalDiscount     float64                `json:"total_discount"`
    TotalPoints       int                    `json:"total_points"`
    AppliedTaxes      float64                `json:"applied_taxes"`
    AppliedFees       float64                `json:"applied_fees"`
    FinalAmount       float64                `json:"final_amount"`
    
    // Applied Rules
    AppliedRules      []AppliedRule          `json:"applied_rules" gorm:"serializer:json"`
    ConflictsResolved []ResolvedConflict     `json:"conflicts_resolved" gorm:"serializer:json"`
    
    // Performance Metrics
    ProcessingTime    time.Duration          `json:"processing_time"`
    ServiceCalls      []ServiceCall          `json:"service_calls" gorm:"serializer:json"`
    CacheHits         int                    `json:"cache_hits"`
    CacheMisses       int                    `json:"cache_misses"`
    
    // Metadata
    ResponseFormat    ResponseFormat         `json:"response_format"`
    IncludeDetails    bool                   `json:"include_details"`
    GeneratedAt       time.Time              `json:"generated_at"`
    ExpiresAt         *time.Time             `json:"expires_at,omitempty"`
}

type AppliedRule struct {
    RuleID          string      `json:"rule_id"`
    RuleName        string      `json:"rule_name"`
    RuleType        string      `json:"rule_type"`
    Service         string      `json:"service"`
    Applied         bool        `json:"applied"`
    Result          interface{} `json:"result"`
    ExecutionTime   time.Duration `json:"execution_time"`
    CacheUsed       bool        `json:"cache_used"`
}

type ServiceCall struct {
    Service       string        `json:"service"`
    Method        string        `json:"method"`
    Duration      time.Duration `json:"duration"`
    Success       bool          `json:"success"`
    ErrorMessage  string        `json:"error_message,omitempty"`
    CacheHit      bool          `json:"cache_hit"`
}
```

#### Load Balancer Configuration
```go
type LoadBalancerConfig struct {
    ID                string                 `json:"id" gorm:"primaryKey"`
    ServiceName       string                 `json:"service_name" gorm:"uniqueIndex"`
    Strategy          LoadBalancingStrategy  `json:"strategy" gorm:"default:'ROUND_ROBIN'"`
    Endpoints         []ServiceEndpoint      `json:"endpoints" gorm:"serializer:json"`
    HealthCheck       HealthCheckConfig      `json:"health_check" gorm:"embedded"`
    CircuitBreaker    CircuitBreakerConfig   `json:"circuit_breaker" gorm:"embedded"`
    RetryPolicy       RetryPolicyConfig      `json:"retry_policy" gorm:"embedded"`
    IsActive          bool                   `json:"is_active" gorm:"default:true"`
    CreatedAt         time.Time              `json:"created_at"`
    UpdatedAt         time.Time              `json:"updated_at"`
}

type LoadBalancingStrategy string
const (
    StrategyRoundRobin     LoadBalancingStrategy = "ROUND_ROBIN"
    StrategyWeightedRound  LoadBalancingStrategy = "WEIGHTED_ROUND_ROBIN"
    StrategyLeastConnection LoadBalancingStrategy = "LEAST_CONNECTION"
    StrategyWeightedRandom LoadBalancingStrategy = "WEIGHTED_RANDOM"
    StrategyIPHash         LoadBalancingStrategy = "IP_HASH"
)

type ServiceEndpoint struct {
    URL           string  `json:"url"`
    Weight        int     `json:"weight"`
    IsHealthy     bool    `json:"is_healthy"`
    LastChecked   time.Time `json:"last_checked"`
    ResponseTime  time.Duration `json:"response_time"`
    ErrorCount    int     `json:"error_count"`
    MaxErrors     int     `json:"max_errors"`
}

type HealthCheckConfig struct {
    Enabled       bool          `json:"enabled"`
    Interval      time.Duration `json:"interval"`
    Timeout       time.Duration `json:"timeout"`
    HealthyThreshold   int      `json:"healthy_threshold"`
    UnhealthyThreshold int      `json:"unhealthy_threshold"`
    Path          string        `json:"path"`
    Method        string        `json:"method"`
}

type CircuitBreakerConfig struct {
    Enabled           bool          `json:"enabled"`
    FailureThreshold  int           `json:"failure_threshold"`
    SuccessThreshold  int           `json:"success_threshold"`
    Timeout           time.Duration `json:"timeout"`
    HalfOpenRequests  int           `json:"half_open_requests"`
}

type RetryPolicyConfig struct {
    Enabled       bool          `json:"enabled"`
    MaxRetries    int           `json:"max_retries"`
    InitialDelay  time.Duration `json:"initial_delay"`
    MaxDelay      time.Duration `json:"max_delay"`
    Multiplier    float64       `json:"multiplier"`
    RetryOn       []string      `json:"retry_on"`
}
```

## REST API Endpoints

### Evaluation API

#### Evaluate Rules
```
POST /api/v1/evaluate
Content-Type: application/json

{
  "customer_id": "customer-123",
  "session_id": "session-456",
  "transaction": {
    "amount": 150.00,
    "currency": "USD",
    "type": "purchase",
    "channel": "web"
  },
  "customer": {
    "id": "customer-123",
    "tier": "GOLD",
    "loyalty_points": 1250
  },
  "products": [
    {
      "id": "product-789",
      "sku": "LAPTOP-001",
      "category": "electronics",
      "price": 150.00,
      "quantity": 1
    }
  ],
  "options": {
    "include_details": true,
    "response_format": "DETAILED",
    "rule_filters": ["promotions", "loyalty"]
  }
}

Response: 200 OK
{
  "request_id": "req-123456",
  "success": true,
  "results": {
    "total_discount": 15.00,
    "total_points": 150,
    "applied_taxes": 12.00,
    "applied_fees": 2.50,
    "final_amount": 149.50
  },
  "applied_rules": [
    {
      "rule_id": "rule-001",
      "rule_name": "10% Electronics Discount",
      "rule_type": "promotion",
      "service": "promotions",
      "applied": true,
      "result": {
        "discount_amount": 15.00,
        "discount_percentage": 10
      }
    }
  ],
  "performance": {
    "processing_time": "245ms",
    "service_calls": [
      {
        "service": "rules-calculation",
        "duration": "123ms",
        "cache_hit": false
      }
    ]
  }
}
```

#### Batch Evaluation
```
POST /api/v1/evaluate/batch
Content-Type: application/json

{
  "requests": [
    {
      "id": "batch-1",
      "customer_id": "customer-123",
      "transaction": {...}
    },
    {
      "id": "batch-2", 
      "customer_id": "customer-456",
      "transaction": {...}
    }
  ],
  "options": {
    "parallel_processing": true,
    "max_concurrency": 10
  }
}

Response: 200 OK
{
  "batch_id": "batch-789",
  "total_requests": 2,
  "successful": 2,
  "failed": 0,
  "results": [
    {
      "request_id": "batch-1",
      "success": true,
      "results": {...}
    }
  ]
}
```

#### Get Evaluation Status
```
GET /api/v1/evaluate/{request_id}/status

Response: 200 OK
{
  "request_id": "req-123456",
  "status": "COMPLETED",
  "progress": 100,
  "estimated_completion": null,
  "created_at": "2024-01-15T10:00:00Z",
  "completed_at": "2024-01-15T10:00:00.245Z"
}
```

### Health and Monitoring

#### Service Health
```
GET /api/v1/health

Response: 200 OK
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:00:00Z",
  "version": "1.0.0",
  "dependencies": {
    "rules-calculation": "healthy",
    "rules-management": "healthy",
    "promotions": "healthy",
    "loyalty": "healthy",
    "coupons": "healthy",
    "taxes": "healthy",
    "payments": "healthy",
    "calculator": "healthy"
  },
  "metrics": {
    "requests_per_second": 1250,
    "average_response_time": "298ms",
    "error_rate": 0.02,
    "cache_hit_rate": 0.78
  }
}
```

#### Performance Metrics
```
GET /api/v1/metrics

Response: 200 OK
{
  "current_load": {
    "requests_per_second": 1250,
    "concurrent_requests": 45,
    "queue_length": 12
  },
  "response_times": {
    "p50": "180ms",
    "p95": "420ms",
    "p99": "650ms"
  },
  "service_performance": {
    "rules-calculation": {
      "average_response_time": "145ms",
      "error_rate": 0.01
    }
  },
  "cache_statistics": {
    "hit_rate": 0.78,
    "miss_rate": 0.22,
    "eviction_rate": 0.05
  }
}
```

## gRPC Service Definition

```protobuf
syntax = "proto3";

package rules_evaluation;

service RulesEvaluationService {
  // Core evaluation methods
  rpc EvaluateRules(EvaluateRulesRequest) returns (EvaluateRulesResponse);
  rpc EvaluateRulesBatch(EvaluateRulesBatchRequest) returns (EvaluateRulesBatchResponse);
  rpc GetEvaluationStatus(GetEvaluationStatusRequest) returns (GetEvaluationStatusResponse);
  
  // Service management
  rpc GetServiceHealth(GetServiceHealthRequest) returns (GetServiceHealthResponse);
  rpc GetPerformanceMetrics(GetPerformanceMetricsRequest) returns (GetPerformanceMetricsResponse);
  rpc UpdateLoadBalancerConfig(UpdateLoadBalancerConfigRequest) returns (UpdateLoadBalancerConfigResponse);
}

message EvaluateRulesRequest {
  string customer_id = 1;
  string session_id = 2;
  TransactionData transaction = 3;
  CustomerData customer = 4;
  repeated ProductData products = 5;
  map<string, string> context = 6;
  EvaluationOptions options = 7;
}

message EvaluateRulesResponse {
  string request_id = 1;
  bool success = 2;
  string error_message = 3;
  EvaluationResults results = 4;
  repeated AppliedRule applied_rules = 5;
  PerformanceMetrics performance = 6;
}

message EvaluationOptions {
  bool include_details = 1;
  string response_format = 2;
  repeated string rule_filters = 3;
  int32 timeout_seconds = 4;
  bool use_cache = 5;
}

message EvaluationResults {
  double total_discount = 1;
  int32 total_points = 2;
  double applied_taxes = 3;
  double applied_fees = 4;
  double final_amount = 5;
}
```

## Domain Services

### Evaluation Orchestrator Service
```go
type EvaluationOrchestratorService interface {
    EvaluateRules(ctx context.Context, req *EvaluationRequest) (*EvaluationResponse, error)
    EvaluateRulesBatch(ctx context.Context, requests []*EvaluationRequest) ([]*EvaluationResponse, error)
    GetEvaluationStatus(ctx context.Context, requestID string) (*EvaluationStatus, error)
    CancelEvaluation(ctx context.Context, requestID string) error
}

type EvaluationStatus struct {
    RequestID           string    `json:"request_id"`
    Status              RequestStatus `json:"status"`
    Progress            int       `json:"progress"`
    EstimatedCompletion *time.Time `json:"estimated_completion,omitempty"`
    CreatedAt           time.Time `json:"created_at"`
    CompletedAt         *time.Time `json:"completed_at,omitempty"`
}
```

### Load Balancer Service
```go
type LoadBalancerService interface {
    GetHealthyEndpoint(ctx context.Context, serviceName string) (*ServiceEndpoint, error)
    UpdateEndpointHealth(ctx context.Context, serviceName, url string, isHealthy bool) error
    GetServiceMetrics(ctx context.Context, serviceName string) (*ServiceMetrics, error)
    UpdateLoadBalancerConfig(ctx context.Context, config *LoadBalancerConfig) error
}

type ServiceMetrics struct {
    ServiceName      string        `json:"service_name"`
    TotalRequests    int64         `json:"total_requests"`
    SuccessfulReqs   int64         `json:"successful_requests"`
    FailedRequests   int64         `json:"failed_requests"`
    AverageResponseTime time.Duration `json:"average_response_time"`
    ErrorRate        float64       `json:"error_rate"`
    LastUpdated      time.Time     `json:"last_updated"`
}
```

### Cache Service
```go
type CacheService interface {
    GetEvaluationResult(ctx context.Context, key string) (*EvaluationResponse, error)
    SetEvaluationResult(ctx context.Context, key string, result *EvaluationResponse, ttl time.Duration) error
    InvalidateCache(ctx context.Context, pattern string) error
    GetCacheStats(ctx context.Context) (*CacheStats, error)
    WarmupCache(ctx context.Context, keys []string) error
}

type CacheStats struct {
    HitRate         float64   `json:"hit_rate"`
    MissRate        float64   `json:"miss_rate"`
    EvictionRate    float64   `json:"eviction_rate"`
    TotalRequests   int64     `json:"total_requests"`
    MemoryUsage     int64     `json:"memory_usage"`
    KeyCount        int64     `json:"key_count"`
    LastUpdated     time.Time `json:"last_updated"`
}
```

## Implementation Tasks

### Phase 1: Core Gateway Infrastructure (3-4 days)
1. **Project Setup and Gateway Foundation**
   - Initialize Go project with Gin framework
   - Setup request routing and middleware
   - Implement request/response logging and metrics
   - Configure CORS and security headers

2. **Load Balancer Implementation**
   - Implement service discovery and endpoint management
   - Create load balancing strategies (round-robin, weighted, etc.)
   - Add health checking for backend services
   - Implement circuit breaker patterns

3. **Request/Response Handling**
   - Create request validation and sanitization
   - Implement response transformation and aggregation
   - Add request timeout and cancellation handling
   - Create error handling and status code mapping

### Phase 2: Service Integration Layer (3-4 days)
1. **Backend Service Clients**
   - Implement gRPC clients for all backend services
   - Add connection pooling and management
   - Create service-specific request/response mapping
   - Implement retry logic and fallback mechanisms

2. **Orchestration Engine**
   - Create evaluation orchestrator for coordinating service calls
   - Implement parallel and sequential execution strategies
   - Add result aggregation and conflict resolution
   - Create dependency management between service calls

3. **Cache Integration**
   - Implement Redis cache for evaluation results
   - Add cache key generation and invalidation strategies
   - Create cache warming and preloading mechanisms
   - Implement cache statistics and monitoring

### Phase 3: Performance Optimization (2-3 days)
1. **Response Time Optimization**
   - Implement request batching and bulk operations
   - Add response compression and optimization
   - Create intelligent caching strategies
   - Optimize JSON serialization and deserialization

2. **Throughput Enhancement**
   - Implement connection pooling for all services
   - Add request queuing and priority handling
   - Create horizontal scaling preparations
   - Optimize memory usage and garbage collection

3. **Monitoring and Observability**
   - Add Prometheus metrics for all operations
   - Implement distributed tracing with correlation IDs
   - Create performance dashboards and alerts
   - Add detailed logging for debugging

### Phase 4: Advanced Features (2-3 days)
1. **Batch Processing**
   - Implement batch evaluation endpoints
   - Add parallel processing for batch requests
   - Create batch result aggregation
   - Implement batch status tracking and reporting

2. **Advanced Caching**
   - Implement multi-level caching (L1: memory, L2: Redis)
   - Add intelligent cache warming based on usage patterns
   - Create cache invalidation based on rule updates
   - Implement cache analytics and optimization

3. **Circuit Breakers and Resilience**
   - Implement circuit breakers for all backend services
   - Add fallback mechanisms and degraded responses
   - Create service health monitoring and alerting
   - Implement graceful degradation strategies

### Phase 5: Testing and Production Readiness (3-4 days)
1. **Comprehensive Testing**
   - Write unit tests for all service components (80% coverage)
   - Create integration tests with backend services
   - Implement load testing and performance validation
   - Add chaos engineering tests for resilience

2. **Security and Compliance**
   - Implement input validation and sanitization
   - Add rate limiting and DDoS protection
   - Create audit logging for compliance
   - Implement security headers and HTTPS enforcement

3. **Production Deployment**
   - Create Docker containers and Kubernetes configurations
   - Setup monitoring and alerting for production
   - Implement backup and disaster recovery procedures
   - Create operational runbooks and troubleshooting guides

## Estimated Development Time: 13-18 days

## Performance Targets
- **Response Time**: <300ms (95th percentile) for standard evaluations
- **Throughput**: 2000+ RPS sustained, 5000+ RPS peak
- **Availability**: 99.95% uptime with proper failover
- **Cache Hit Rate**: >85% for frequently evaluated rules
- **Service Integration**: <50ms overhead for service orchestration
