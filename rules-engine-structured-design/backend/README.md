# Backend Services Design - Golang Microservices

## Overview

The backend consists of 9 independent microservices implemented in Golang, following Domain-Driven Design principles. Each service represents a bounded context with clear domain boundaries, implementing both REST and gRPC APIs.

## Microservices Architecture

### Core Domain Services

#### 1. Rules Management Service
- **Purpose**: Rule lifecycle management, approval workflows, templates
- **Domain**: Rule creation, editing, versioning, approval
- **APIs**: REST (external), gRPC (internal)
- **Database**: PostgreSQL (rules, templates, versions)
- **Key Entities**: Rule, RuleTemplate, RuleVersion, ApprovalWorkflow

#### 2. Rules Calculation Service  
- **Purpose**: High-performance rule evaluation engine
- **Domain**: Rule execution, conflict resolution, optimization
- **APIs**: gRPC (high-performance internal)
- **Database**: PostgreSQL (rules cache), Redis (performance cache)
- **Key Entities**: CalculationEngine, RuleExecution, ConflictResolver

#### 3. Rules Evaluation Service
- **Purpose**: External API gateway for rule evaluation
- **Domain**: API orchestration, load balancing, result aggregation
- **APIs**: REST (external clients), gRPC (internal services)
- **Database**: Redis (response cache), Elasticsearch (audit logs)
- **Key Entities**: EvaluationRequest, EvaluationResponse, ResultAggregator

### Supporting Domain Services

#### 4. Promotions Service
- **Purpose**: Promotional campaigns and discount rules
- **Domain**: Campaign management, discount calculation, targeting
- **APIs**: REST, gRPC
- **Database**: PostgreSQL
- **Key Entities**: PromotionalCampaign, DiscountRule, CustomerSegment

#### 5. Loyalty Service
- **Purpose**: Customer loyalty programs and points management
- **Domain**: Points earning/redemption, tier management, rewards
- **APIs**: REST, gRPC
- **Database**: PostgreSQL
- **Key Entities**: LoyaltyProgram, CustomerTier, PointsTransaction

#### 6. Coupons Service
- **Purpose**: Coupon validation and redemption
- **Domain**: Coupon lifecycle, fraud prevention, usage tracking
- **APIs**: REST, gRPC
- **Database**: PostgreSQL, Redis (validation cache)
- **Key Entities**: Coupon, CouponCampaign, RedemptionRecord

#### 7. Taxes Service
- **Purpose**: Tax calculation and compliance
- **Domain**: Tax rules, jurisdiction management, compliance reporting
- **APIs**: REST, gRPC
- **Database**: PostgreSQL
- **Key Entities**: TaxRule, TaxJurisdiction, TaxCalculation

#### 8. Payments Service
- **Purpose**: Payment rules and fraud detection
- **Domain**: Payment routing, fraud prevention, gateway management
- **APIs**: REST, gRPC
- **Database**: PostgreSQL
- **Key Entities**: PaymentRule, PaymentGateway, FraudRule

#### 9. Calculator Service
- **Purpose**: Shared calculation engine and mathematical operations
- **Domain**: Mathematical calculations, algorithm optimization, shared utilities
- **APIs**: gRPC (internal only)
- **Database**: Redis (calculation cache)
- **Key Entities**: CalculationRequest, CalculationResult, Algorithm

## Shared Components

### Domain Layer
- **Common Value Objects**: Money, Date ranges, Identifiers
- **Common Interfaces**: Repository patterns, Domain events
- **Common Enums**: Status types, Priority levels
- **Domain Events**: Base event types and publishing interfaces

### Infrastructure Layer
- **Database Connections**: PostgreSQL, Redis, Elasticsearch clients
- **Messaging**: NATS client and event publishing
- **Monitoring**: Prometheus metrics, logging utilities
- **Configuration**: Environment-based configuration management

### Protocols Layer
- **gRPC Definitions**: Service contracts and message types
- **REST Contracts**: OpenAPI specifications
- **Event Schemas**: Message schemas for NATS
- **Authentication**: JWT validation and authorization

## API Gateway

### Responsibilities
- **Authentication**: JWT token validation with MFA support
- **Authorization**: Enhanced RBAC with attribute-based access control
- **Routing**: Request routing to appropriate services
- **Rate Limiting**: Advanced API usage throttling
- **Monitoring**: Request/response logging and metrics
- **Circuit Breaking**: Fault tolerance and resilience patterns

### Implementation - Emissary-Ingress

```yaml
# API gateway configuration (https://emissary-ingress.dev/)
api_gateway:
  rate_limiting:
    default: "1000 requests/minute/user"
    premium: "5000 requests/minute/user"
    service_specific:
      rules_evaluation: "2000 requests/minute/user"
      rules_calculation: "500 requests/minute/user"
  
  authentication:
    jwt:
      issuer: "rules-engine-auth"
      algorithm: "RS256"
      expiry: "24 hours"
      refresh: "7 days"
  
  load_balancing:
    algorithm: "round_robin"
    health_checks:
      interval: "30 seconds"
      timeout: "5 seconds"
      unhealthy_threshold: 3
    
  circuit_breaker:
    failure_threshold: "50%"
    timeout: "30 seconds"
    half_open_max_calls: 10
```

### Legacy Implementation (Fallback)
- **Framework**: Gin (Go HTTP framework)
- **Authentication**: JWT with RS256
- **Documentation**: Swagger/OpenAPI 3.0
- **Health Checks**: Service health monitoring

## Development Standards

### Code Structure (per service)
```
service-name/
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── domain/                 # Domain entities and business logic
│   │   ├── entities/           # Domain entities
│   │   ├── repositories/       # Repository interfaces
│   │   ├── services/          # Domain services
│   │   └── events/            # Domain events
│   ├── application/           # Application services and use cases
│   │   ├── commands/          # Command handlers
│   │   ├── queries/           # Query handlers
│   │   └── handlers/          # HTTP/gRPC handlers
│   ├── infrastructure/        # External concerns
│   │   ├── persistence/       # Database implementation
│   │   ├── messaging/         # NATS implementation
│   │   └── external/          # External API clients
│   └── interfaces/            # API definitions
│       ├── rest/              # REST endpoints
│       └── grpc/              # gRPC services
├── migrations/                # Database migrations
├── proto/                     # Protocol buffer definitions
├── docker/                    # Docker configuration
└── tests/                     # Test files
```

### Dependencies
- **Framework**: Gin (REST), gRPC-Go
- **Database**: GORM (PostgreSQL), go-redis
- **Messaging**: NATS client
- **Monitoring**: Prometheus client
- **Testing**: Testify, Mockery
- **Configuration**: Viper
- **Logging**: Logrus

### Quality Standards
- **Test Coverage**: Minimum 80% unit test coverage
- **API Documentation**: Complete OpenAPI/gRPC documentation
- **Error Handling**: Structured error responses
- **Logging**: Structured JSON logging
- **Metrics**: Prometheus metrics for all endpoints
- **Health Checks**: Kubernetes-ready health endpoints

## Performance Specifications

### Service-Level Performance Targets

| Service | Response Time | Throughput | Availability | Error Rate | Resource Limits |
|---------|---------------|------------|--------------|------------|-----------------|
| Rule Management | <200ms (CRUD) | 500 RPS | 99.9% | <0.1% | 512MB RAM |
| Rule Calculation | <500ms (complex) | 1000+ TPS | 99.95% | <0.01% | 2GB RAM |
| Rule Evaluation | <300ms (95th) | 2000+ RPS | 99.95% | <0.05% | 1GB RAM |
| Promotions | <150ms (queries) | 800 RPS | 99.9% | <0.1% | 512MB RAM |
| Loyalty | <200ms (updates) | 600 RPS | 99.9% | <0.1% | 512MB RAM |
| Coupons | <100ms (validation) | 1500 RPS | 99.95% | <0.01% | 512MB RAM |
| Taxes | <300ms (calculation) | 800 RPS | 99.9% | <0.1% | 1GB RAM |
| Payments | <250ms (processing) | 500 RPS | 99.99% | <0.001% | 1GB RAM |
| Calculator | <50ms (basic) | 10000+ ops/s | 99.9% | <0.01% | 256MB RAM |

### SLI/SLO/SLA Framework
- **Service Level Indicators (SLIs)**: Response time, throughput, availability, error rate
- **Service Level Objectives (SLOs)**: Internal targets for service quality
- **Service Level Agreements (SLAs)**: External commitments to customers

## Communication Patterns

### Service Communication Matrix

```go
// Service communication standards
type CommunicationPattern string
const (
    Synchronous  CommunicationPattern = "sync"   // <300ms SLA
    Asynchronous CommunicationPattern = "async"  // Event-driven
    Hybrid       CommunicationPattern = "hybrid" // Both patterns
)

// Service communication matrix
var ServiceCommunication = map[string]map[string]CommunicationPattern{
    "rules-management": {
        "rules-calculation": Asynchronous,
        "audit-service":     Asynchronous,
        "cache-service":     Asynchronous,
    },
    "rules-evaluation": {
        "rules-calculation": Synchronous,
        "analytics":         Asynchronous,
        "audit-service":     Asynchronous,
    },
}
```

### Synchronous Communication
- **REST APIs**: External client communication
- **gRPC**: Internal service-to-service communication
- **API Gateway**: Single entry point for external requests

### Asynchronous Communication
- **NATS**: Event-driven communication between services
- **Event Sourcing**: Domain events for audit and state reconstruction
- **Eventual Consistency**: Accepting eventual consistency between bounded contexts

## Data Management

### Database Architecture & Scaling Strategy

#### PostgreSQL Cluster Configuration
- **Development**: Local PostgreSQL installation using Rancher Desktop for testing
- **Production**: GCP Cloud SQL-based solution with high availability

```yaml
# Database Scaling Configuration
postgresql_cluster:
  primary_node:
    instance_type: "db.r5.xlarge"
    storage: "1TB SSD"
    max_connections: 200
  
  read_replicas:
    count: 3
    instance_type: "db.r5.large"
    lag_threshold: "100ms"
    auto_failover: true
  
  sharding_strategy:
    rules_table:
      shard_key: "tenant_id"
      shards: 4
      rebalancing: "automatic"
    
    audit_logs:
      shard_key: "created_date"
      retention: "7_years"
      archival: "s3_glacier"
```

#### Connection Pooling & Performance
- **PgBouncer**: Connection pooling for all services
- **Performance Monitoring**: Custom metrics with PostgreSQL stats
- **Database Optimization**: Automated index optimization and query performance analysis

### Multi-Level Caching Strategy

```yaml
# Multi-level caching configuration
caching_strategy:
  level_1_memory:
    technology: "Go sync.Map + TTL"
    size_limit: "128 MB per service"
    ttl: "5 minutes"
    use_cases: ["frequently_accessed_rules", "user_sessions"]
    
  level_2_redis:
    technology: 'Redis Cluster'
    nodes: 3
    memory_per_node: '4 GB'
    ttl: '1 hour'
    use_cases: ['calculation_results', "user_profiles", 'promotion_data']
  
  level_3_database:
    technology: 'PostgreSQL with indexes'
    query_cache: 'enabled'
    materialised_views: 'for_complex_aggregations'
    use_cases: ['master_data', "audit_logs", 'historical_analytics']

# Cache invalidation patterns
invalidation_patterns:
  rule_updates:
    trigger: 'RuleUpdated event'
    scope: 'tenant_specific'
    cascade: ['calculation_cache', 'evaluation_cache']
  promotion_changes:
    trigger: 'PromotionModified event'
    scope: 'campaign_specific'
    cascade: ["promotion_cache", 'discount_cache']
```

### Database per Service
- Each microservice has its own database
- No direct database sharing between services
- Data consistency through events and eventual consistency

### Data Synchronization
- **Event-Driven**: Domain events for data synchronization
- **Eventual Consistency**: Accept eventual consistency
- **Compensation**: Saga pattern for distributed transactions
