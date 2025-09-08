# Campaigns Management Service

A microservice for managing marketing campaigns, built with Go, following Domain-Driven Design (DDD) principles and hexagonal architecture.

## Overview

The Campaigns Management Service is responsible for:
- Creating and managing marketing campaigns
- Campaign lifecycle management (draft, active, paused, completed, cancelled)
- Campaign targeting and personalization
- Performance metrics tracking
- Integration with the Rules Engine for targeting rules

## Architecture

This service follows:
- **Domain-Driven Design (DDD)**: Rich domain models with business logic
- **Hexagonal Architecture**: Clean separation between domain, application, and infrastructure layers
- **CQRS Pattern**: Separate command and query handlers
- **Event-Driven Architecture**: Domain events for loose coupling

## Project Structure

```
campaigns-management-service/
├── cmd/                           # Application entry point
│   └── main.go
├── internal/
│   ├── domain/                    # Domain layer (business logic)
│   │   ├── campaign/             # Campaign aggregate
│   │   └── shared/               # Shared domain concepts
│   ├── application/              # Application layer (use cases)
│   │   ├── commands/             # Command handlers
│   │   └── queries/              # Query handlers
│   ├── infrastructure/           # Infrastructure layer
│   │   ├── persistence/          # Database implementations
│   │   ├── messaging/            # Event bus (NATS)
│   │   ├── external/             # External service clients
│   │   └── validation/           # Input validation
│   └── interfaces/               # Interface layer
│       └── rest/                 # REST API handlers and DTOs
├── api/                          # API specifications
│   ├── openapi/                  # OpenAPI 3.0 specs
│   └── proto/                    # gRPC proto files
├── tests/                        # Test suites
│   ├── unit/                     # Unit tests
│   ├── integration/              # Integration tests
│   └── behavioral/               # BDD tests (Gherkin)
├── deployments/                  # Deployment configurations
│   └── k8s/                      # Kubernetes manifests
├── Dockerfile                    # Container image
├── go.mod                        # Go module definition
└── README.md                     # This file
```

## Features

### Campaign Management
- **Create Campaigns**: Define campaigns with targeting rules, budget, and settings
- **Update Campaigns**: Modify campaign details (only in draft/paused state)
- **Campaign Lifecycle**: Draft → Active → Paused/Completed/Cancelled
- **Campaign Types**: Promotion, Loyalty, Coupon, Segmentation, Retargeting

### Targeting & Personalization
- **Rule-Based Targeting**: Integration with Rules Engine for complex targeting
- **Multi-Channel Support**: Email, SMS, Push, Web, Social, Display
- **A/B Testing**: Built-in A/B testing capabilities
- **Scheduling Rules**: Time-based and rule-based scheduling

### Performance Tracking
- **Real-time Metrics**: Impressions, clicks, conversions, revenue
- **Performance Indicators**: CTR, conversion rate, ROAS, ROI
- **Cost Tracking**: Cost per click, cost per conversion

### Event-Driven Architecture
- **Domain Events**: Campaign lifecycle events
- **NATS Integration**: Asynchronous event publishing
- **Event Sourcing**: Audit trail of all campaign changes

## API Endpoints

### Campaign Management
- `POST /api/v1/campaigns` - Create a new campaign
- `GET /api/v1/campaigns` - List campaigns with filtering and pagination
- `GET /api/v1/campaigns/:id` - Get campaign details
- `PUT /api/v1/campaigns/:id` - Update campaign
- `DELETE /api/v1/campaigns/:id` - Delete campaign

### Campaign Actions
- `POST /api/v1/campaigns/:id/activate` - Activate campaign
- `POST /api/v1/campaigns/:id/pause` - Pause campaign

### Analytics
- `GET /api/v1/campaigns/:id/metrics` - Get campaign performance metrics

### Health Check
- `GET /health` - Service health status

## Configuration

The service can be configured using environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | HTTP server port | `8080` |
| `DB_HOST` | PostgreSQL host | `localhost` |
| `DB_PORT` | PostgreSQL port | `5432` |
| `DB_USER` | Database user | `postgres` |
| `DB_PASSWORD` | Database password | `password` |
| `DB_NAME` | Database name | `campaigns_db` |
| `DB_SSLMODE` | SSL mode | `disable` |
| `NATS_URL` | NATS server URL | `nats://localhost:4222` |
| `NATS_USERNAME` | NATS username | (empty) |
| `NATS_PASSWORD` | NATS password | (empty) |
| `NATS_TOKEN` | NATS token | (empty) |
| `RULES_SERVICE_URL` | Rules Engine service URL | `http://localhost:8081` |

## Database Schema

The service uses PostgreSQL with the following main tables:

### campaigns
- `id` (UUID, Primary Key)
- `name` (VARCHAR, Unique)
- `description` (TEXT)
- `status` (VARCHAR) - DRAFT, ACTIVE, PAUSED, COMPLETED, CANCELLED
- `type` (VARCHAR) - PROMOTION, LOYALTY, COUPON, SEGMENTATION, RETARGETING
- `targeting_rules` (JSONB) - Array of rule IDs
- `start_date` (TIMESTAMP)
- `end_date` (TIMESTAMP, nullable)
- `budget_amount` (DECIMAL)
- `budget_currency` (VARCHAR)
- `created_by` (UUID)
- `created_at` (TIMESTAMP)
- `updated_at` (TIMESTAMP)
- `settings` (JSONB) - Campaign settings
- `metrics` (JSONB) - Performance metrics
- `version` (INTEGER) - Optimistic locking

## Domain Model

### Campaign Aggregate
The `Campaign` is the main aggregate root containing:
- **Identity**: CampaignID (UUID)
- **Basic Info**: Name, description, type
- **Lifecycle**: Status and state transitions
- **Targeting**: Rules from the Rules Engine
- **Settings**: Channels, frequency, A/B testing, personalization
- **Metrics**: Performance indicators
- **Audit**: Created/updated timestamps, version

### Value Objects
- **CampaignID**: Typed UUID for campaign identity
- **UserID**: Typed UUID for user identity
- **RuleID**: Typed UUID for rule identity
- **Money**: Monetary value with currency
- **CampaignSettings**: Complex settings configuration
- **CampaignMetrics**: Performance metrics

### Domain Events
- `CampaignCreated`
- `CampaignUpdated`
- `CampaignActivated`
- `CampaignPaused`
- `CampaignCompleted`
- `CampaignCancelled`
- `CampaignDeleted`

## Development

### Prerequisites
- Go 1.21+
- PostgreSQL 13+
- NATS Server
- Docker (for containerization)

### Local Development

1. **Clone and setup**:
   ```bash
   cd campaigns-management-service
   go mod download
   ```

2. **Start dependencies**:
   ```bash
   # Start PostgreSQL
   docker run -d --name postgres -e POSTGRES_PASSWORD=password -p 5432:5432 postgres:13
   
   # Start NATS
   docker run -d --name nats -p 4222:4222 nats:latest
   ```

3. **Run the service**:
   ```bash
   go run cmd/main.go
   ```

### Testing

```bash
# Run unit tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run integration tests
go test -tags=integration ./tests/integration/...

# Run behavioral tests
go test -tags=behavioral ./tests/behavioral/...
```

### Building

```bash
# Build for current platform
go build -o campaigns-management-service cmd/main.go

# Build for Linux (for Docker)
GOOS=linux GOARCH=amd64 go build -o campaigns-management-service cmd/main.go
```

### Docker

```bash
# Build image
docker build -t campaigns-management:latest .

# Run container
docker run -p 8080:8080 \
  -e DB_HOST=host.docker.internal \
  -e NATS_URL=nats://host.docker.internal:4222 \
  campaigns-management:latest
```

## Deployment

### Kubernetes

The service includes Kubernetes manifests for deployment:

```bash
# Deploy to Kubernetes
kubectl apply -f deployments/k8s/

# Check deployment status
kubectl get pods -n rules-engine -l app=campaigns-management-service

# View logs
kubectl logs -n rules-engine -l app=campaigns-management-service
```

### Environment Variables for Production

```yaml
env:
- name: PORT
  value: "8080"
- name: DB_HOST
  value: "postgres-service"
- name: DB_USER
  value: "postgres"
- name: DB_PASSWORD
  valueFrom:
    secretKeyRef:
      name: postgres-secret
      key: password
- name: DB_NAME
  value: "campaigns_db"
- name: NATS_URL
  value: "nats://nats-service:4222"
- name: RULES_SERVICE_URL
  value: "http://rules-management-service:8080"
```

## Monitoring & Observability

### Health Checks
- **Liveness Probe**: `/health` endpoint
- **Readiness Probe**: `/health` endpoint

### Metrics (Planned)
- Campaign creation rate
- Campaign activation rate
- API response times
- Database connection pool status
- NATS connection status

### Logging
- Structured logging with correlation IDs
- Request/response logging
- Domain event logging
- Error logging with stack traces

## Integration

### Rules Engine Integration
The service integrates with the Rules Engine for:
- **Targeting Rules**: Complex customer segmentation
- **Personalization Rules**: Dynamic content customization
- **Scheduling Rules**: Time-based campaign execution

### Event Publishing
Domain events are published to NATS for:
- **Analytics Processing**: Campaign performance tracking
- **Notification Services**: Campaign status updates
- **Audit Logging**: Compliance and debugging
- **Integration**: Third-party system notifications

## Security

### Authentication & Authorization
- JWT token validation (planned)
- Role-based access control (planned)
- API key authentication (planned)

### Data Protection
- Input validation and sanitization
- SQL injection prevention (GORM)
- XSS protection
- CSRF protection (planned)

### Container Security
- Non-root user execution
- Read-only root filesystem
- Minimal base image (Alpine)
- Security context constraints

## Performance

### Optimization Strategies
- **Database Indexing**: Optimized queries for campaign listing
- **Connection Pooling**: Efficient database connections
- **Caching**: Redis integration (planned)
- **Async Processing**: NATS for non-blocking operations

### Scalability
- **Horizontal Scaling**: Stateless service design
- **Database Sharding**: Campaign-based partitioning (planned)
- **Load Balancing**: Kubernetes service load balancing
- **Auto-scaling**: HPA based on CPU/memory usage (planned)

## Contributing

1. Follow Go coding standards
2. Write comprehensive tests
3. Update documentation
4. Follow DDD principles
5. Maintain backward compatibility

## License

This project is part of the ENGINE-RULES-SP system and follows the same licensing terms.
