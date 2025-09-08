# Analytics Dashboard Service

The Analytics Dashboard Service provides comprehensive analytics, reporting, and insights for the rules engine ecosystem. It aggregates data from all services to provide real-time dashboards, custom reports, and business intelligence capabilities.

## Features

### Core Capabilities
- **Real-time Dashboards**: Interactive dashboards with customizable widgets
- **Custom Reports**: Generate reports in multiple formats (PDF, Excel, CSV, JSON, HTML)
- **Metric Management**: Create and manage custom metrics with various aggregation types
- **Data Aggregation**: Aggregate data from all rules engine services
- **Performance Monitoring**: Real-time system performance metrics
- **Business Intelligence**: Business metrics and insights
- **Compliance Reporting**: Automated compliance and audit reports

### Dashboard Features
- Drag-and-drop widget interface
- Multiple widget types (charts, tables, KPIs, gauges, heatmaps, maps)
- Real-time data refresh
- Public/private dashboard sharing
- Responsive layouts
- Custom filters and parameters

### Report Features
- Scheduled report generation
- Multiple output formats
- Email distribution
- Custom templates
- Parameterized reports
- Compliance and audit reports

### Metric Features
- Custom metric creation
- Multiple metric types (counter, gauge, histogram, summary)
- Dimension-based analysis
- Calculated metrics with formulas
- Real-time and historical data
- Performance optimization

## Architecture

The service follows Domain-Driven Design (DDD) principles with a clean architecture:

```
analytics-dashboard-service/
├── cmd/                    # Application entry point
├── internal/
│   ├── domain/            # Domain models and business logic
│   │   ├── analytics/     # Core domain aggregates
│   │   └── shared/        # Shared domain components
│   ├── application/       # Application layer (CQRS)
│   │   ├── commands/      # Command handlers
│   │   └── queries/       # Query handlers
│   ├── infrastructure/    # Infrastructure layer
│   │   ├── persistence/   # Database repositories
│   │   ├── messaging/     # Event publishing
│   │   └── config/        # Configuration
│   └── interfaces/        # Interface layer
│       ├── rest/          # REST API handlers
│       └── grpc/          # gRPC service
├── api/                   # API specifications
│   ├── openapi/          # OpenAPI/Swagger specs
│   └── proto/            # Protocol buffer definitions
├── tests/                # Test suites
└── deployments/          # Deployment configurations
```

## API Endpoints

### Dashboards
- `GET /api/v1/dashboards` - List dashboards
- `POST /api/v1/dashboards` - Create dashboard
- `GET /api/v1/dashboards/{id}` - Get dashboard
- `PUT /api/v1/dashboards/{id}` - Update dashboard
- `DELETE /api/v1/dashboards/{id}` - Delete dashboard

### Reports
- `GET /api/v1/reports` - List reports
- `POST /api/v1/reports` - Create report
- `GET /api/v1/reports/{id}` - Get report
- `PUT /api/v1/reports/{id}` - Update report
- `DELETE /api/v1/reports/{id}` - Delete report
- `POST /api/v1/reports/{id}/generate` - Generate report

### Metrics
- `GET /api/v1/metrics` - List metrics
- `POST /api/v1/metrics` - Create metric
- `GET /api/v1/metrics/{id}` - Get metric
- `PUT /api/v1/metrics/{id}` - Update metric
- `DELETE /api/v1/metrics/{id}` - Delete metric
- `GET /api/v1/metrics/{id}/data` - Get metric data

### Analytics
- `GET /api/v1/analytics/real-time` - Real-time analytics
- `GET /api/v1/analytics/performance` - Performance metrics
- `GET /api/v1/analytics/business` - Business metrics
- `GET /api/v1/analytics/compliance` - Compliance metrics

## Database Schema

### Core Tables
- `dashboards` - Dashboard configurations and metadata
- `reports` - Report definitions and schedules
- `metrics` - Metric definitions and configurations
- `metric_data` - Time-series metric data
- `domain_events` - Event sourcing for audit trails
- `audit_logs` - Audit logging for compliance

### Performance Optimizations
- Partitioned metric data tables by month
- Indexed queries for fast data retrieval
- Connection pooling and query optimization
- Caching for frequently accessed data

## Configuration

### Environment Variables
- `PORT` - Server port (default: 8080)
- `DATABASE_DSN` - PostgreSQL connection string
- `NATS_URL` - NATS messaging URL
- `REDIS_URL` - Redis cache URL
- `RULES_SERVICE_URL` - Rules service endpoint
- `CUSTOMER_SERVICE_URL` - Customer service endpoint
- `CAMPAIGN_SERVICE_URL` - Campaign service endpoint
- `PROMOTION_SERVICE_URL` - Promotion service endpoint

### Performance Settings
- `DB_MAX_OPEN_CONNS` - Database connection pool size
- `DB_MAX_IDLE_CONNS` - Database idle connections
- `CACHE_TTL` - Cache time-to-live
- `READ_TIMEOUT` - HTTP read timeout
- `WRITE_TIMEOUT` - HTTP write timeout

## Development

### Prerequisites
- Go 1.23+
- PostgreSQL 13+
- NATS Server
- Redis (optional, for caching)

### Running Locally
```bash
# Clone the repository
git clone <repository-url>
cd analytics-dashboard-service

# Install dependencies
go mod download

# Set environment variables
export DATABASE_DSN="postgres://user:password@localhost:5432/analytics_db?sslmode=disable"
export NATS_URL="nats://localhost:4222"

# Run the service
go run cmd/main.go
```

### Building
```bash
# Build the application
go build -o analytics-dashboard-service cmd/main.go

# Build Docker image
docker build -t analytics-dashboard-service:latest .
```

### Testing
```bash
# Run unit tests
go test ./...

# Run integration tests
go test -tags=integration ./...

# Run behavioral tests
go test -tags=behavioral ./...
```

## Deployment

### Docker
```bash
# Build and run with Docker
docker build -t analytics-dashboard-service .
docker run -p 8080:8080 analytics-dashboard-service
```

### Kubernetes
```bash
# Apply Kubernetes manifests
kubectl apply -f deployments/k8s/
```

### Production Considerations
- Use external PostgreSQL database
- Configure NATS clustering for high availability
- Set up Redis for caching
- Configure proper resource limits
- Enable monitoring and alerting
- Set up log aggregation
- Configure backup strategies

## Monitoring

### Health Checks
- `GET /health` - Service health status
- `GET /metrics` - Prometheus metrics

### Key Metrics
- Request latency and throughput
- Database connection pool status
- Cache hit/miss ratios
- Error rates by endpoint
- Resource utilization

## Security

### Authentication & Authorization
- JWT token validation
- Role-based access control
- API key authentication
- Rate limiting

### Data Protection
- Encryption in transit (TLS)
- Encryption at rest
- PII data masking
- Audit logging
- Secure configuration management

## Performance Requirements

### Response Time Targets
- Dashboard Load: < 2 seconds
- Widget Refresh: < 1 second
- Report Generation: < 30 seconds
- Metric Query: < 500ms
- Real-time Data: < 100ms

### Scalability Targets
- 1,000+ concurrent users
- 1TB+ analytics data
- 100,000+ metrics per minute
- 1,000+ reports per hour
- 10,000+ dashboard views per hour

## Contributing

1. Follow the established code structure and patterns
2. Write comprehensive tests for new features
3. Update documentation for API changes
4. Follow the coding standards and conventions
5. Ensure all tests pass before submitting PRs

## License

This project is proprietary and confidential. All rights reserved.
