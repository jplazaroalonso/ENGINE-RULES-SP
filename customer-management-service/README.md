# Customer Management Service

A comprehensive backend service that handles customer data, segmentation, analytics, and integration with the rules engine for personalized customer experiences. This service provides GDPR-compliant customer data management with advanced segmentation and analytics capabilities.

## Features

### Core Functionality
- **Customer Management**: Complete CRUD operations for customer data
- **Customer Segmentation**: Rule-based customer segmentation with dynamic criteria
- **Customer Analytics**: Comprehensive analytics and insights generation
- **GDPR Compliance**: Full compliance with privacy regulations including data export, deletion, and anonymization
- **Event Tracking**: Real-time customer event tracking and analytics
- **Privacy Management**: Granular consent management and privacy controls

### Technical Features
- **Hexagonal Architecture**: Clean separation of concerns with domain-driven design
- **Event-Driven**: NATS JetStream for asynchronous event publishing
- **Database**: PostgreSQL with GORM for robust data persistence
- **Validation**: Comprehensive input validation with custom validators
- **Security**: Non-root containers, read-only filesystems, and capability dropping
- **Observability**: Health checks, structured logging, and metrics endpoints
- **Scalability**: Horizontal scaling with stateless design

## Architecture

### Domain Model

#### Customer Aggregate
- **CustomerID**: Unique identifier for customers
- **Email**: Email address with validation
- **Personal Information**: Name, age, gender, location
- **Preferences**: Language, currency, timezone, notification settings, privacy settings
- **Segments**: Dynamic segment membership
- **Metadata**: Purchase history, interaction history, device information
- **Status**: Active, Inactive, Suspended, Deleted

#### Customer Segment Aggregate
- **SegmentID**: Unique identifier for segments
- **Criteria**: Demographic, behavioral, geographic, and custom criteria
- **Rule Integration**: Integration with rules engine for dynamic evaluation
- **Membership Management**: Automatic customer assignment and removal

#### Customer Analytics
- **Metrics**: Customer lifetime value, engagement scores, churn risk
- **Insights**: Behavior patterns, recommendations, risk factors
- **Event Tracking**: Real-time event capture and analysis
- **Reporting**: Comprehensive reporting and analytics

### API Endpoints

#### Customer Management
```
GET    /api/v1/customers                     # List customers with pagination/filtering
POST   /api/v1/customers                     # Create new customer
GET    /api/v1/customers/:id                 # Get customer details
PUT    /api/v1/customers/:id                 # Update customer
DELETE /api/v1/customers/:id                 # Delete customer (GDPR compliant)
```

#### Customer Analytics
```
GET    /api/v1/customers/:id/analytics       # Get customer analytics
GET    /api/v1/customers/:id/insights        # Get customer insights
POST   /api/v1/customers/:id/track           # Track customer event
GET    /api/v1/customers/:id/segments        # Get customer segments
```

#### Customer Segmentation
```
GET    /api/v1/customers/segments            # List customer segments
POST   /api/v1/customers/segments            # Create new segment
GET    /api/v1/customers/segments/:id        # Get segment details
PUT    /api/v1/customers/segments/:id        # Update segment
DELETE /api/v1/customers/segments/:id        # Delete segment
POST   /api/v1/customers/segments/:id/calculate # Calculate segment membership
GET    /api/v1/customers/segments/:id/customers # Get segment customers
```

#### Customer Privacy & GDPR
```
GET    /api/v1/customers/:id/data            # Export customer data (GDPR)
DELETE /api/v1/customers/:id/data            # Delete customer data (GDPR)
PUT    /api/v1/customers/:id/consent         # Update privacy consent
GET    /api/v1/customers/:id/consent         # Get privacy consent status
POST   /api/v1/customers/:id/anonymize       # Anonymize customer data
```

#### Bulk Operations
```
POST   /api/v1/customers/bulk/update         # Bulk update customers
POST   /api/v1/customers/bulk/delete         # Bulk delete customers
POST   /api/v1/customers/bulk/segments       # Bulk assign segments
```

## Getting Started

### Prerequisites
- Go 1.21 or later
- PostgreSQL 13 or later
- NATS Server with JetStream enabled
- Docker and Kubernetes (for deployment)

### Environment Variables
```bash
PORT=8080                                    # HTTP server port
DATABASE_DSN=postgres://user:pass@host/db   # PostgreSQL connection string
NATS_URL=nats://localhost:4222              # NATS server URL
RULES_ENGINE_URL=http://localhost:8081      # Rules engine service URL
RULES_ENGINE_API_KEY=your-api-key           # Rules engine API key
```

### Local Development

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd customer-management-service
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up database**
   ```bash
   # Create PostgreSQL database
   createdb customer_management
   
   # Run migrations (handled automatically on startup)
   ```

4. **Start NATS server**
   ```bash
   nats-server -js
   ```

5. **Run the service**
   ```bash
   go run cmd/main.go
   ```

### Docker Deployment

1. **Build the image**
   ```bash
   docker build -t customer-management-service .
   ```

2. **Run the container**
   ```bash
   docker run -p 8080:8080 \
     -e DATABASE_DSN="postgres://user:pass@host/db" \
     -e NATS_URL="nats://host:4222" \
     customer-management-service
   ```

### Kubernetes Deployment

1. **Apply Kubernetes manifests**
   ```bash
   kubectl apply -f deployments/k8s/
   ```

2. **Check deployment status**
   ```bash
   kubectl get pods -n rules-engine
   kubectl get services -n rules-engine
   kubectl get ingress -n rules-engine
   ```

## Database Schema

### Customers Table
```sql
CREATE TABLE customers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    age INTEGER CHECK (age >= 0 AND age <= 150),
    gender VARCHAR(20) CHECK (gender IN ('MALE', 'FEMALE', 'OTHER', 'UNKNOWN')),
    location JSONB,
    preferences JSONB NOT NULL DEFAULT '{}',
    segments JSONB NOT NULL DEFAULT '[]',
    tags JSONB NOT NULL DEFAULT '[]',
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_activity TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}',
    version INTEGER NOT NULL DEFAULT 1
);
```

### Customer Segments Table
```sql
CREATE TABLE customer_segments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    rule_id UUID NOT NULL,
    customer_count INTEGER NOT NULL DEFAULT 0,
    criteria JSONB NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
    created_by UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_calculated TIMESTAMP WITH TIME ZONE,
    version INTEGER NOT NULL DEFAULT 1
);
```

### Customer Segment Membership Table
```sql
CREATE TABLE customer_segment_membership (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    segment_id UUID NOT NULL REFERENCES customer_segments(id) ON DELETE CASCADE,
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    left_at TIMESTAMP WITH TIME ZONE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);
```

### Customer Events Table
```sql
CREATE TABLE customer_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    event_type VARCHAR(50) NOT NULL,
    event_data JSONB,
    occurred_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    session_id UUID,
    device_info JSONB
);
```

## Business Rules

### Customer Creation Rules
1. **Email Uniqueness**: Email addresses must be unique across all customers
2. **Email Validation**: Email must be in valid format
3. **Age Validation**: Age must be between 0 and 150 if provided
4. **Location Validation**: Country and timezone must be valid if provided
5. **Consent Requirements**: Marketing and data processing consent must be explicitly given

### Customer Update Rules
1. **Email Change**: Email changes require verification
2. **Status Transitions**: 
   - ACTIVE → INACTIVE (allowed)
   - INACTIVE → ACTIVE (allowed)
   - Any status → SUSPENDED (requires admin approval)
   - Any status → DELETED (GDPR deletion process)

### Segmentation Rules
1. **Segment Creation**: 
   - Name must be unique within organization
   - At least one criteria must be specified
   - Rule ID must reference valid rule in rules engine

2. **Segment Calculation**:
   - Automatic recalculation when criteria change
   - Background processing for large segments
   - Error handling for invalid criteria

### Privacy and GDPR Rules
1. **Data Export**: Complete customer data export within 30 days
2. **Data Deletion**: Complete data deletion within 30 days
3. **Consent Management**: Consent can be withdrawn at any time
4. **Data Anonymization**: Anonymize instead of delete when legally required
5. **Audit Trail**: All privacy operations must be logged

## Performance Requirements

### Response Time Targets
- **List Customers**: < 300ms for 10,000 customers
- **Get Customer**: < 100ms
- **Create Customer**: < 200ms
- **Update Customer**: < 150ms
- **Get Segments**: < 200ms
- **Calculate Segment**: < 5 seconds for 100,000 customers
- **Track Event**: < 50ms

### Scalability Requirements
- **Customer Volume**: Support 10,000,000+ customers
- **Event Throughput**: Handle 1,000,000+ events per minute
- **Segment Calculation**: Process 1,000,000+ customers per segment
- **Concurrent Users**: Support 5,000+ concurrent users
- **Data Volume**: Store 10TB+ of customer and event data

## Security

### Authentication & Authorization
- **JWT Token Validation**: All endpoints require valid JWT tokens
- **Role-Based Access**: Different permissions for customer data access
- **Data Access Control**: Field-level access control for sensitive data
- **API Key Support**: For external service integrations
- **Rate Limiting**: Prevent abuse with configurable rate limits

### Data Protection
- **Encryption at Rest**: All sensitive data encrypted in database
- **Encryption in Transit**: TLS 1.3 for all API communications
- **PII Protection**: Special handling for personal identifiable information
- **GDPR Compliance**: Full compliance with GDPR requirements
- **Audit Logging**: Complete audit trail for all operations

## Monitoring and Observability

### Health Checks
- **Liveness Probe**: `/health` endpoint for Kubernetes health checks
- **Readiness Probe**: Service readiness validation
- **Database Health**: Database connection and query health
- **External Dependencies**: NATS and Rules Engine connectivity

### Metrics
- **Business Metrics**: Customer acquisition, retention, lifetime value
- **Technical Metrics**: Response times, error rates, throughput
- **System Metrics**: CPU, memory, disk usage
- **Privacy Metrics**: GDPR requests, consent changes, data exports

### Logging
- **Structured Logging**: JSON format with consistent fields
- **Log Levels**: DEBUG, INFO, WARN, ERROR with appropriate usage
- **Correlation IDs**: Track requests across service boundaries
- **Privacy Logging**: Special handling for PII in logs

## Testing

### Unit Tests
- **Domain Logic**: 95%+ coverage for business rules and invariants
- **Repository Layer**: Mock database interactions
- **Service Layer**: Test business logic with mocked dependencies
- **Handler Layer**: Test HTTP request/response handling

### Integration Tests
- **Database Integration**: Test with real PostgreSQL database
- **External Service Integration**: Mock external API calls
- **Message Queue Integration**: Test NATS event publishing
- **End-to-End API Testing**: Complete request/response cycles

### Performance Tests
- **Load Testing**: Test with expected production load
- **Stress Testing**: Test system limits and failure scenarios
- **Database Performance**: Test query performance with large datasets
- **Memory Usage**: Monitor memory consumption under load

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

For support and questions:
- Create an issue in the repository
- Contact the development team
- Check the documentation and examples

## Roadmap

### Phase 1: Core Service (3 weeks)
- Domain model and business logic
- Basic CRUD operations
- Database schema and migrations
- Unit tests

### Phase 2: API Implementation (2 weeks)
- REST API endpoints
- Request/response handling
- Input validation
- Integration tests

### Phase 3: Advanced Features (3 weeks)
- Customer segmentation
- Analytics and insights
- GDPR compliance features
- Rule engine integration

### Phase 4: Production Readiness (2 weeks)
- Security implementation
- Performance optimization
- Monitoring and logging
- Deployment configuration

**Total Estimated Effort**: 10 weeks
**Team Size**: 3-4 developers
**Dependencies**: Rules Management Service, Analytics Service
