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
- **Authentication**: JWT token validation
- **Authorization**: Role-based access control
- **Routing**: Request routing to appropriate services
- **Rate Limiting**: API usage throttling
- **Monitoring**: Request/response logging and metrics

### Implementation
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

## Communication Patterns

### Synchronous Communication
- **REST APIs**: External client communication
- **gRPC**: Internal service-to-service communication
- **API Gateway**: Single entry point for external requests

### Asynchronous Communication
- **NATS**: Event-driven communication between services
- **Event Sourcing**: Domain events for audit and state reconstruction
- **Eventual Consistency**: Accepting eventual consistency between bounded contexts

## Data Management

### Database per Service
- Each microservice has its own database
- No direct database sharing between services
- Data consistency through events and eventual consistency

### Caching Strategy
- **Redis**: Shared cache for frequently accessed data
- **Local Cache**: In-memory caching for static data
- **Cache Invalidation**: Event-driven cache invalidation

### Data Synchronization
- **Event-Driven**: Domain events for data synchronization
- **Eventual Consistency**: Accept eventual consistency
- **Compensation**: Saga pattern for distributed transactions
