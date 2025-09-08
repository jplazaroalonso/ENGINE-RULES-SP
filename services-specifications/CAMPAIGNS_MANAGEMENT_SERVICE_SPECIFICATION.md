# Campaigns Management Service Specification

## Executive Summary

The Campaigns Management Service is a comprehensive backend service that handles marketing campaigns, targeting, performance tracking, and integration with the rules engine. This service enables the creation, management, and optimization of marketing campaigns using business rules for targeting and personalization.

## Service Overview

### Purpose
- Manage marketing campaigns with rule-based targeting
- Track campaign performance and ROI
- Integrate with rules engine for dynamic targeting
- Support A/B testing and campaign optimization
- Provide comprehensive campaign analytics

### Business Value
- **Targeted Marketing**: Use business rules for precise customer targeting
- **Performance Optimization**: Real-time metrics for campaign optimization
- **Automated Campaigns**: Rule-driven campaign activation and management
- **ROI Tracking**: Comprehensive performance and financial metrics
- **Scalable Operations**: Support for multiple concurrent campaigns

## Technical Architecture

### Service Structure
```
campaigns-management-service/
├── cmd/
│   └── main.go
├── internal/
│   ├── domain/
│   │   ├── campaign/
│   │   │   ├── campaign.go
│   │   │   ├── campaign_metrics.go
│   │   │   ├── campaign_settings.go
│   │   │   ├── repository.go
│   │   │   └── service.go
│   │   └── shared/
│   │       ├── errors.go
│   │       └── events.go
│   ├── application/
│   │   ├── commands/
│   │   │   ├── create_campaign.go
│   │   │   ├── update_campaign.go
│   │   │   ├── activate_campaign.go
│   │   │   ├── pause_campaign.go
│   │   │   └── delete_campaign.go
│   │   └── queries/
│   │       ├── get_campaign.go
│   │       ├── list_campaigns.go
│   │       └── get_campaign_metrics.go
│   ├── infrastructure/
│   │   ├── persistence/
│   │   │   └── postgres/
│   │   │       ├── campaign_repository.go
│   │   │       └── migrations/
│   │   ├── messaging/
│   │   │   └── nats/
│   │   │       └── event_publisher.go
│   │   └── external/
│   │       ├── rules_client.go
│   │       └── analytics_client.go
│   └── interfaces/
│       ├── rest/
│       │   ├── handlers/
│       │   │   └── campaign_handler.go
│       │   └── dto/
│       │       └── campaign_dto.go
│       └── grpc/
│           └── campaign_service.go
├── api/
│   ├── openapi/
│   │   └── campaigns-api.v1.yaml
│   └── proto/
│       └── campaigns.proto
├── tests/
│   ├── unit/
│   ├── integration/
│   └── behavioral/
├── deployments/
│   └── k8s/
├── Dockerfile
├── go.mod
└── go.sum
```

## Domain Model

### Core Entities

#### Campaign Aggregate
```go
type Campaign struct {
    ID              CampaignID       `json:"id"`
    Name            string           `json:"name"`
    Description     string           `json:"description"`
    Status          CampaignStatus   `json:"status"`
    Type            CampaignType     `json:"type"`
    TargetingRules  []RuleID         `json:"targetingRules"`
    StartDate       time.Time        `json:"startDate"`
    EndDate         *time.Time       `json:"endDate,omitempty"`
    Budget          *Money           `json:"budget,omitempty"`
    CreatedBy       UserID           `json:"createdBy"`
    CreatedAt       time.Time        `json:"createdAt"`
    UpdatedAt       time.Time        `json:"updatedAt"`
    Settings        CampaignSettings `json:"settings"`
    Metrics         CampaignMetrics  `json:"metrics"`
    Version         int              `json:"version"`
    Events          []DomainEvent    `json:"-"`
}

type CampaignStatus string

const (
    CampaignStatusDraft     CampaignStatus = "DRAFT"
    CampaignStatusActive    CampaignStatus = "ACTIVE"
    CampaignStatusPaused    CampaignStatus = "PAUSED"
    CampaignStatusCompleted CampaignStatus = "COMPLETED"
    CampaignStatusCancelled CampaignStatus = "CANCELLED"
)

type CampaignType string

const (
    CampaignTypePromotion     CampaignType = "PROMOTION"
    CampaignTypeLoyalty       CampaignType = "LOYALTY"
    CampaignTypeCoupon        CampaignType = "COUPON"
    CampaignTypeSegmentation  CampaignType = "SEGMENTATION"
    CampaignTypeRetargeting   CampaignType = "RETARGETING"
)
```

#### Campaign Settings Value Object
```go
type CampaignSettings struct {
    TargetAudience    []string           `json:"targetAudience"`
    Channels          []Channel          `json:"channels"`
    Frequency         Frequency          `json:"frequency"`
    MaxImpressions    *int               `json:"maxImpressions,omitempty"`
    BudgetLimit       *Money             `json:"budgetLimit,omitempty"`
    ABTestSettings    *ABTestSettings    `json:"abTestSettings,omitempty"`
    SchedulingRules   []SchedulingRule   `json:"schedulingRules"`
    Personalization   PersonalizationConfig `json:"personalization"`
}

type Channel string

const (
    ChannelEmail      Channel = "EMAIL"
    ChannelSMS        Channel = "SMS"
    ChannelPush       Channel = "PUSH"
    ChannelWeb        Channel = "WEB"
    ChannelSocial     Channel = "SOCIAL"
    ChannelDisplay    Channel = "DISPLAY"
)

type Frequency string

const (
    FrequencyOnce     Frequency = "ONCE"
    FrequencyDaily    Frequency = "DAILY"
    FrequencyWeekly   Frequency = "WEEKLY"
    FrequencyMonthly  Frequency = "MONTHLY"
)
```

#### Campaign Metrics Value Object
```go
type CampaignMetrics struct {
    Impressions       int64   `json:"impressions"`
    Clicks            int64   `json:"clicks"`
    Conversions       int64   `json:"conversions"`
    Revenue           Money   `json:"revenue"`
    Cost              Money   `json:"cost"`
    CTR               float64 `json:"ctr"`               // Click-through rate
    ConversionRate    float64 `json:"conversionRate"`    // Conversion rate
    CostPerClick      Money   `json:"costPerClick"`      // CPC
    CostPerConversion Money   `json:"costPerConversion"` // CPA
    ROAS              float64 `json:"roas"`              // Return on ad spend
    ROI               float64 `json:"roi"`               // Return on investment
    LastUpdated       time.Time `json:"lastUpdated"`
}
```

### Domain Services

#### Campaign Targeting Service
```go
type CampaignTargetingService interface {
    EvaluateTargeting(ctx context.Context, campaignID CampaignID, customerID CustomerID) (bool, error)
    GetTargetAudience(ctx context.Context, campaignID CampaignID) ([]CustomerID, error)
    UpdateTargetingRules(ctx context.Context, campaignID CampaignID, rules []RuleID) error
}
```

#### Campaign Performance Service
```go
type CampaignPerformanceService interface {
    CalculateMetrics(ctx context.Context, campaignID CampaignID) (*CampaignMetrics, error)
    TrackEvent(ctx context.Context, campaignID CampaignID, event CampaignEvent) error
    GetPerformanceReport(ctx context.Context, campaignID CampaignID, period TimePeriod) (*PerformanceReport, error)
}
```

## API Specification

### REST API Endpoints

#### Campaign Management
```
GET    /api/v1/campaigns                     # List campaigns with pagination/filtering
POST   /api/v1/campaigns                     # Create new campaign
GET    /api/v1/campaigns/:id                 # Get campaign details
PUT    /api/v1/campaigns/:id                 # Update campaign
DELETE /api/v1/campaigns/:id                 # Delete campaign
```

#### Campaign Lifecycle
```
POST   /api/v1/campaigns/:id/activate        # Activate campaign
POST   /api/v1/campaigns/:id/pause           # Pause campaign
POST   /api/v1/campaigns/:id/stop            # Stop campaign
POST   /api/v1/campaigns/:id/duplicate       # Duplicate campaign
```

#### Campaign Analytics
```
GET    /api/v1/campaigns/:id/metrics         # Get campaign metrics
GET    /api/v1/campaigns/:id/performance     # Get performance report
GET    /api/v1/campaigns/:id/audience        # Get target audience
POST   /api/v1/campaigns/:id/track           # Track campaign event
```

#### Bulk Operations
```
POST   /api/v1/campaigns/bulk/activate       # Bulk activate campaigns
POST   /api/v1/campaigns/bulk/pause          # Bulk pause campaigns
POST   /api/v1/campaigns/bulk/delete         # Bulk delete campaigns
```

#### Import/Export
```
GET    /api/v1/campaigns/export              # Export campaigns
POST   /api/v1/campaigns/import              # Import campaigns
```

### Request/Response Models

#### Create Campaign Request
```json
{
  "name": "Summer Sale 2024",
  "description": "Summer promotion campaign with 20% discount",
  "type": "PROMOTION",
  "targetingRules": ["rule-123", "rule-456"],
  "startDate": "2024-06-01T00:00:00Z",
  "endDate": "2024-08-31T23:59:59Z",
  "budget": {
    "amount": 10000.00,
    "currency": "EUR"
  },
  "settings": {
    "targetAudience": ["segment-1", "segment-2"],
    "channels": ["EMAIL", "WEB", "PUSH"],
    "frequency": "WEEKLY",
    "maxImpressions": 100000,
    "personalization": {
      "enabled": true,
      "rules": ["personalization-rule-1"]
    }
  }
}
```

#### Campaign Response
```json
{
  "id": "campaign-123",
  "name": "Summer Sale 2024",
  "description": "Summer promotion campaign with 20% discount",
  "status": "ACTIVE",
  "type": "PROMOTION",
  "targetingRules": ["rule-123", "rule-456"],
  "startDate": "2024-06-01T00:00:00Z",
  "endDate": "2024-08-31T23:59:59Z",
  "budget": {
    "amount": 10000.00,
    "currency": "EUR"
  },
  "createdBy": "user-123",
  "createdAt": "2024-05-15T10:30:00Z",
  "updatedAt": "2024-05-15T10:30:00Z",
  "settings": {
    "targetAudience": ["segment-1", "segment-2"],
    "channels": ["EMAIL", "WEB", "PUSH"],
    "frequency": "WEEKLY",
    "maxImpressions": 100000,
    "personalization": {
      "enabled": true,
      "rules": ["personalization-rule-1"]
    }
  },
  "metrics": {
    "impressions": 15420,
    "clicks": 1234,
    "conversions": 89,
    "revenue": {
      "amount": 4450.00,
      "currency": "EUR"
    },
    "cost": {
      "amount": 890.00,
      "currency": "EUR"
    },
    "ctr": 8.0,
    "conversionRate": 7.2,
    "roi": 400.0,
    "lastUpdated": "2024-05-20T15:45:00Z"
  }
}
```

## Business Rules and Invariants

### Campaign Creation Rules
1. **Name Uniqueness**: Campaign names must be unique within the organization
2. **Date Validation**: Start date must be in the future, end date must be after start date
3. **Budget Validation**: Budget must be positive and within organization limits
4. **Targeting Rules**: At least one targeting rule must be specified
5. **Channel Validation**: At least one channel must be selected

### Campaign Lifecycle Rules
1. **Status Transitions**: 
   - DRAFT → ACTIVE (requires approval if budget > threshold)
   - ACTIVE → PAUSED (allowed anytime)
   - PAUSED → ACTIVE (allowed if within date range)
   - ACTIVE/PAUSED → COMPLETED (automatic when end date reached)
   - Any status → CANCELLED (requires justification)

2. **Budget Management**:
   - Cannot exceed budget limit during execution
   - Automatic pause when 90% of budget consumed
   - Automatic stop when 100% of budget consumed

3. **Date Constraints**:
   - Cannot activate campaign with past start date
   - Cannot modify active campaign's targeting rules
   - Cannot delete campaign with active metrics

### Performance Rules
1. **Metrics Calculation**: Real-time calculation of CTR, conversion rate, ROI
2. **Event Tracking**: All campaign events must be tracked for analytics
3. **Audience Updates**: Target audience recalculated when rules change
4. **Performance Alerts**: Automatic alerts for performance degradation

## Integration Requirements

### Rules Engine Integration
- **Targeting Rules**: Evaluate customer eligibility using rules engine
- **Personalization Rules**: Apply dynamic content based on customer profile
- **Performance Rules**: Trigger actions based on campaign performance

### Analytics Integration
- **Metrics Collection**: Send campaign events to analytics service
- **Performance Reporting**: Generate comprehensive performance reports
- **Trend Analysis**: Provide historical performance trends

### External Services
- **Email Service**: Send campaign emails through external provider
- **SMS Service**: Send SMS campaigns through external provider
- **Push Notification Service**: Send push notifications
- **Social Media APIs**: Manage social media campaigns

## Database Schema

### Campaigns Table
```sql
CREATE TABLE campaigns (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'DRAFT',
    type VARCHAR(50) NOT NULL,
    targeting_rules JSONB NOT NULL DEFAULT '[]',
    start_date TIMESTAMP WITH TIME ZONE NOT NULL,
    end_date TIMESTAMP WITH TIME ZONE,
    budget_amount DECIMAL(15,2),
    budget_currency VARCHAR(3),
    created_by UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    settings JSONB NOT NULL DEFAULT '{}',
    version INTEGER NOT NULL DEFAULT 1,
    
    CONSTRAINT campaigns_name_unique UNIQUE (name),
    CONSTRAINT campaigns_status_check CHECK (status IN ('DRAFT', 'ACTIVE', 'PAUSED', 'COMPLETED', 'CANCELLED')),
    CONSTRAINT campaigns_type_check CHECK (type IN ('PROMOTION', 'LOYALTY', 'COUPON', 'SEGMENTATION', 'RETARGETING')),
    CONSTRAINT campaigns_dates_check CHECK (end_date IS NULL OR end_date > start_date),
    CONSTRAINT campaigns_budget_check CHECK (budget_amount IS NULL OR budget_amount > 0)
);

CREATE INDEX idx_campaigns_status ON campaigns(status);
CREATE INDEX idx_campaigns_type ON campaigns(type);
CREATE INDEX idx_campaigns_dates ON campaigns(start_date, end_date);
CREATE INDEX idx_campaigns_created_by ON campaigns(created_by);
```

### Campaign Metrics Table
```sql
CREATE TABLE campaign_metrics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    campaign_id UUID NOT NULL REFERENCES campaigns(id) ON DELETE CASCADE,
    impressions BIGINT NOT NULL DEFAULT 0,
    clicks BIGINT NOT NULL DEFAULT 0,
    conversions BIGINT NOT NULL DEFAULT 0,
    revenue_amount DECIMAL(15,2) NOT NULL DEFAULT 0,
    revenue_currency VARCHAR(3) NOT NULL DEFAULT 'EUR',
    cost_amount DECIMAL(15,2) NOT NULL DEFAULT 0,
    cost_currency VARCHAR(3) NOT NULL DEFAULT 'EUR',
    last_updated TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT campaign_metrics_positive CHECK (
        impressions >= 0 AND clicks >= 0 AND conversions >= 0 AND
        revenue_amount >= 0 AND cost_amount >= 0
    )
);

CREATE INDEX idx_campaign_metrics_campaign_id ON campaign_metrics(campaign_id);
CREATE INDEX idx_campaign_metrics_last_updated ON campaign_metrics(last_updated);
```

### Campaign Events Table
```sql
CREATE TABLE campaign_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    campaign_id UUID NOT NULL REFERENCES campaigns(id) ON DELETE CASCADE,
    event_type VARCHAR(50) NOT NULL,
    customer_id UUID,
    event_data JSONB,
    occurred_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT campaign_events_type_check CHECK (
        event_type IN ('IMPRESSION', 'CLICK', 'CONVERSION', 'BOUNCE', 'UNSUBSCRIBE')
    )
);

CREATE INDEX idx_campaign_events_campaign_id ON campaign_events(campaign_id);
CREATE INDEX idx_campaign_events_type ON campaign_events(event_type);
CREATE INDEX idx_campaign_events_occurred_at ON campaign_events(occurred_at);
```

## Performance Requirements

### Response Time Targets
- **List Campaigns**: < 200ms for 1000 campaigns
- **Get Campaign**: < 100ms
- **Create Campaign**: < 500ms
- **Update Campaign**: < 300ms
- **Get Metrics**: < 150ms
- **Track Event**: < 50ms

### Scalability Requirements
- **Concurrent Campaigns**: Support 10,000+ active campaigns
- **Event Throughput**: Handle 100,000+ events per minute
- **User Concurrency**: Support 1,000+ concurrent users
- **Data Volume**: Store 1TB+ of campaign and metrics data

### Availability Requirements
- **Uptime**: 99.9% availability
- **Recovery Time**: < 5 minutes for service recovery
- **Data Durability**: 99.999% data durability
- **Backup**: Daily automated backups with point-in-time recovery

## Security Requirements

### Authentication & Authorization
- **JWT Token Validation**: All endpoints require valid JWT tokens
- **Role-Based Access**: Different permissions for campaign management
- **API Key Support**: For external service integrations
- **Rate Limiting**: Prevent abuse with configurable rate limits

### Data Protection
- **Encryption at Rest**: All sensitive data encrypted in database
- **Encryption in Transit**: TLS 1.3 for all API communications
- **PII Handling**: Proper handling of customer personal information
- **Audit Logging**: Complete audit trail for all operations

### Input Validation
- **Request Validation**: Comprehensive input validation and sanitization
- **SQL Injection Prevention**: Parameterized queries and ORM usage
- **XSS Prevention**: Output encoding and CSP headers
- **File Upload Security**: Secure handling of campaign assets

## Testing Strategy

### Unit Testing
- **Domain Logic**: 95%+ coverage for business rules and invariants
- **Repository Layer**: Mock database interactions
- **Service Layer**: Test business logic with mocked dependencies
- **Handler Layer**: Test HTTP request/response handling

### Integration Testing
- **Database Integration**: Test with real PostgreSQL database
- **External Service Integration**: Mock external API calls
- **Message Queue Integration**: Test NATS event publishing
- **End-to-End API Testing**: Complete request/response cycles

### Performance Testing
- **Load Testing**: Test with expected production load
- **Stress Testing**: Test system limits and failure scenarios
- **Database Performance**: Test query performance with large datasets
- **Memory Usage**: Monitor memory consumption under load

### Security Testing
- **Authentication Testing**: Test JWT validation and authorization
- **Input Validation Testing**: Test malicious input handling
- **SQL Injection Testing**: Test database security
- **Rate Limiting Testing**: Test abuse prevention

## Deployment Configuration

### Docker Configuration
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o campaigns-service ./cmd

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /app/campaigns-service .
EXPOSE 8080 9090
CMD ["./campaigns-service"]
```

### Kubernetes Deployment
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: campaigns-management-service
  namespace: rules-engine
spec:
  replicas: 3
  selector:
    matchLabels:
      app: campaigns-management-service
  template:
    metadata:
      labels:
        app: campaigns-management-service
    spec:
      containers:
      - name: campaigns-service
        image: localhost:5000/campaigns-management:latest
        ports:
        - containerPort: 8080
        - containerPort: 9090
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: campaigns-db-secret
              key: url
        - name: NATS_URL
          valueFrom:
            configMapKeyRef:
              name: campaigns-config
              key: nats-url
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
```

## Monitoring and Observability

### Metrics Collection
- **Business Metrics**: Campaign performance, conversion rates, ROI
- **Technical Metrics**: Response times, error rates, throughput
- **System Metrics**: CPU, memory, disk usage
- **Custom Metrics**: Campaign-specific KPIs and alerts

### Logging Strategy
- **Structured Logging**: JSON format with consistent fields
- **Log Levels**: DEBUG, INFO, WARN, ERROR with appropriate usage
- **Correlation IDs**: Track requests across service boundaries
- **Sensitive Data**: Proper handling of PII and sensitive information

### Alerting Rules
- **Service Health**: Alert on service downtime or high error rates
- **Performance**: Alert on response time degradation
- **Business Metrics**: Alert on campaign performance issues
- **Resource Usage**: Alert on high resource consumption

## Success Criteria

### Functional Requirements
- ✅ **Campaign CRUD**: Complete create, read, update, delete operations
- ✅ **Lifecycle Management**: Proper campaign status transitions
- ✅ **Performance Tracking**: Real-time metrics and analytics
- ✅ **Rule Integration**: Seamless integration with rules engine
- ✅ **Bulk Operations**: Efficient bulk campaign management

### Non-Functional Requirements
- ✅ **Performance**: Meet all response time targets
- ✅ **Scalability**: Support required concurrent users and campaigns
- ✅ **Availability**: Achieve 99.9% uptime
- ✅ **Security**: Pass all security requirements
- ✅ **Reliability**: Handle failures gracefully with proper recovery

### Business Value
- ✅ **ROI Tracking**: Accurate campaign ROI calculation
- ✅ **Targeting**: Effective rule-based customer targeting
- ✅ **Automation**: Automated campaign management workflows
- ✅ **Analytics**: Comprehensive campaign performance insights
- ✅ **Integration**: Seamless integration with existing systems

## Implementation Timeline

### Phase 1: Core Service (2 weeks)
- Domain model and business logic
- Basic CRUD operations
- Database schema and migrations
- Unit tests

### Phase 2: API Implementation (1 week)
- REST API endpoints
- Request/response handling
- Input validation
- Integration tests

### Phase 3: Advanced Features (2 weeks)
- Campaign lifecycle management
- Performance tracking
- Rule engine integration
- Bulk operations

### Phase 4: Production Readiness (1 week)
- Security implementation
- Performance optimization
- Monitoring and logging
- Deployment configuration

**Total Estimated Effort**: 6 weeks
**Team Size**: 2-3 developers
**Dependencies**: Rules Management Service, Analytics Service
