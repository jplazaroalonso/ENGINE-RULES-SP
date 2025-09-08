# Settings Management Service

A microservice for managing configurations, feature flags, user preferences, and organization settings in the Rules Engine platform.

## Overview

The Settings Management Service provides a centralized way to manage various types of settings across the Rules Engine platform. It follows a hexagonal architecture pattern and implements Domain-Driven Design principles.

## Features

- **Configuration Management**: Store and manage application configurations with environment-specific values
- **Feature Flags**: Control feature rollouts and A/B testing with percentage-based targeting
- **User Preferences**: Manage user-specific settings and preferences
- **Organization Settings**: Handle organization-level configurations and settings
- **Event-Driven Architecture**: Publish domain events for settings changes
- **Caching**: Redis-based caching for improved performance
- **Security**: Encrypted storage for sensitive settings

## Architecture

The service follows a hexagonal architecture with the following layers:

```
├── cmd/                    # Application entry point
├── internal/
│   ├── domain/            # Domain layer (entities, value objects, business logic)
│   ├── application/       # Application layer (commands, queries, handlers)
│   ├── infrastructure/    # Infrastructure layer (repositories, external services)
│   └── interfaces/        # Interface layer (REST/gRPC handlers, DTOs)
├── api/                   # API contracts (OpenAPI, Protobuf)
├── tests/                 # Test suites
└── deployments/           # Deployment configurations
```

## Domain Models

### Configuration
- **Purpose**: Store application configurations with environment-specific values
- **Key Features**: Environment separation, encryption support, categorization
- **Use Cases**: Database connections, API endpoints, feature toggles

### Feature Flag
- **Purpose**: Control feature rollouts and A/B testing
- **Key Features**: Percentage-based rollout, target audience, conditions
- **Use Cases**: Feature toggles, A/B testing, gradual rollouts

### User Preference
- **Purpose**: Store user-specific settings and preferences
- **Key Features**: User isolation, public/private preferences, categorization
- **Use Cases**: UI preferences, notification settings, personalization

### Organization Setting
- **Purpose**: Manage organization-level configurations
- **Key Features**: Organization isolation, encryption support, categorization
- **Use Cases**: Billing settings, compliance configurations, custom branding

## API Endpoints

### Configurations
- `GET /v1/configurations` - List configurations
- `POST /v1/configurations` - Create configuration
- `GET /v1/configurations/{id}` - Get configuration
- `PUT /v1/configurations/{id}` - Update configuration
- `DELETE /v1/configurations/{id}` - Delete configuration

### Feature Flags
- `GET /v1/feature-flags` - List feature flags
- `POST /v1/feature-flags` - Create feature flag
- `GET /v1/feature-flags/{id}` - Get feature flag
- `PUT /v1/feature-flags/{id}` - Update feature flag
- `DELETE /v1/feature-flags/{id}` - Delete feature flag

### User Preferences
- `GET /v1/user-preferences` - List user preferences
- `POST /v1/user-preferences` - Create user preference
- `GET /v1/user-preferences/{id}` - Get user preference
- `PUT /v1/user-preferences/{id}` - Update user preference
- `DELETE /v1/user-preferences/{id}` - Delete user preference

### Organization Settings
- `GET /v1/organization-settings` - List organization settings
- `POST /v1/organization-settings` - Create organization setting
- `GET /v1/organization-settings/{id}` - Get organization setting
- `PUT /v1/organization-settings/{id}` - Update organization setting
- `DELETE /v1/organization-settings/{id}` - Delete organization setting

## Technology Stack

- **Language**: Go 1.21
- **Framework**: Gin (HTTP), gRPC
- **Database**: PostgreSQL with GORM
- **Caching**: Redis
- **Messaging**: NATS JetStream
- **Validation**: go-playground/validator
- **Testing**: Testify
- **Containerization**: Docker
- **Orchestration**: Kubernetes

## Getting Started

### Prerequisites

- Go 1.21 or later
- PostgreSQL 13 or later
- Redis 6 or later
- NATS 2.9 or later

### Local Development

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd settings-management-service
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up environment variables**
   ```bash
   export POSTGRES_HOST=localhost
   export POSTGRES_PORT=5432
   export POSTGRES_DB=settings_db
   export POSTGRES_USER=postgres
   export POSTGRES_PASSWORD=password
   export NATS_URL=nats://localhost:4222
   export REDIS_URL=redis://localhost:6379
   export PORT=8080
   ```

4. **Run the service**
   ```bash
   go run cmd/main.go
   ```

### Docker

1. **Build the image**
   ```bash
   docker build -t settings-management-service .
   ```

2. **Run the container**
   ```bash
   docker run -p 8080:8080 \
     -e POSTGRES_HOST=host.docker.internal \
     -e POSTGRES_PORT=5432 \
     -e POSTGRES_DB=settings_db \
     -e POSTGRES_USER=postgres \
     -e POSTGRES_PASSWORD=password \
     -e NATS_URL=nats://host.docker.internal:4222 \
     -e REDIS_URL=redis://host.docker.internal:6379 \
     settings-management-service
   ```

### Kubernetes

1. **Apply the configurations**
   ```bash
   kubectl apply -f deployments/k8s/
   ```

2. **Check the deployment**
   ```bash
   kubectl get pods -n rules-engine
   kubectl get services -n rules-engine
   kubectl get ingress -n rules-engine
   ```

## Testing

### Unit Tests
```bash
go test ./internal/domain/...
go test ./internal/application/...
go test ./internal/infrastructure/...
go test ./internal/interfaces/...
```

### Integration Tests
```bash
go test ./tests/integration/...
```

### Behavioral Tests
```bash
go test ./tests/behavioral/...
```

### Run All Tests
```bash
go test ./...
```

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | HTTP server port | `8080` |
| `POSTGRES_HOST` | PostgreSQL host | `localhost` |
| `POSTGRES_PORT` | PostgreSQL port | `5432` |
| `POSTGRES_DB` | PostgreSQL database name | `settings_db` |
| `POSTGRES_USER` | PostgreSQL username | `postgres` |
| `POSTGRES_PASSWORD` | PostgreSQL password | - |
| `NATS_URL` | NATS server URL | `nats://localhost:4222` |
| `REDIS_URL` | Redis server URL | `redis://localhost:6379` |

### Database Schema

The service automatically creates the following tables:
- `configurations` - Configuration storage
- `feature_flags` - Feature flag storage
- `user_preferences` - User preference storage
- `organization_settings` - Organization setting storage

## Monitoring and Observability

### Health Checks
- **Endpoint**: `GET /health`
- **Response**: Service health status

### Metrics
- **Endpoint**: `GET /metrics`
- **Format**: Prometheus metrics

### Logging
- **Format**: Structured JSON logging
- **Levels**: DEBUG, INFO, WARN, ERROR

## Security

### Data Encryption
- Sensitive settings can be encrypted at rest
- Encryption keys managed externally
- Support for different encryption algorithms

### Access Control
- JWT-based authentication
- Role-based access control (RBAC)
- API key authentication for service-to-service communication

### Data Privacy
- User preferences can be marked as public/private
- Organization settings are isolated by organization
- Audit logging for all changes

## Performance

### Caching Strategy
- Redis-based caching for frequently accessed settings
- Cache invalidation on updates
- Configurable cache TTL

### Database Optimization
- Indexed queries for common access patterns
- Connection pooling
- Query optimization

### Scalability
- Horizontal scaling support
- Stateless service design
- Load balancer ready

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For support and questions:
- Create an issue in the repository
- Contact the development team
- Check the documentation

## Changelog

### v1.0.0
- Initial release
- Configuration management
- Feature flag support
- User preference management
- Organization setting management
- REST API implementation
- gRPC API implementation
- Docker and Kubernetes support
