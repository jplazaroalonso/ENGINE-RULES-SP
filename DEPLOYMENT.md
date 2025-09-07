# Rules Engine - Kubernetes Deployment Guide

This guide provides instructions for deploying the Rules Engine microservices to a local Rancher Desktop cluster using containerd/nerdctl, Emissary Ingress, and cert-manager.

## Prerequisites

- Rancher Desktop with containerd/nerdctl
- Kubernetes cluster running
- Emissary Ingress (Ambassador) deployed
- cert-manager deployed
- PostgreSQL deployed
- kubectl configured to access the cluster

## Architecture

The Rules Engine consists of three microservices:

1. **Rules Management Service** (Port 8080) - Manages business rules
2. **Rules Evaluation Service** (Port 8081) - Evaluates rules using different strategies
3. **Rules Calculator Service** (Port 8082) - Calculates results from rule evaluations

## Deployment Options

### Option 1: Deploy from Scratch (Recommended for first-time setup)

This script builds Docker images, creates all Kubernetes resources, initializes the database, and sets up ingress with TLS certificates.

```bash
./scripts/deploy-from-scratch.sh
```

### Option 2: Deploy Services Only

This script only deploys the services (assumes infrastructure is already set up).

```bash
./scripts/deploy-services-only.sh
```

### Option 3: Cleanup

This script removes all resources created by the deployment.

```bash
./scripts/cleanup.sh
```

## Manual Deployment Steps

If you prefer to deploy manually, follow these steps:

### 1. Build Docker Images

```bash
# Build and push images to local registry
nerdctl build -t localhost:5000/rules-management-service:latest ./rules-management-service/
nerdctl build -t localhost:5000/rules-evaluation-service:latest ./rules-evaluation-service/
nerdctl build -t localhost:5000/rules-calculator-service:latest ./rules-calculator-service/

nerdctl push localhost:5000/rules-management-service:latest
nerdctl push localhost:5000/rules-evaluation-service:latest
nerdctl push localhost:5000/rules-calculator-service:latest
```

### 2. Create Namespace and Basic Resources

```bash
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/secrets.yaml
```

### 3. Initialize Database

```bash
kubectl apply -f k8s/db-init-job.yaml
kubectl wait --for=condition=complete job/db-init-job -n rules-engine --timeout=300s
```

### 4. Deploy Services

```bash
kubectl apply -f k8s/rules-management-service.yaml
kubectl apply -f k8s/rules-evaluation-service.yaml
kubectl apply -f k8s/rules-calculator-service.yaml
```

### 5. Setup Ingress

```bash
kubectl apply -f k8s/ingress.yaml
kubectl wait --for=condition=ready certificate/rules-engine-tls -n rules-engine --timeout=300s
```

### 6. Update /etc/hosts

Add the following entries to your `/etc/hosts` file:

```
127.0.0.1 rules-management.local.dev
127.0.0.1 rules-evaluation.local.dev
127.0.0.1 rules-calculator.local.dev
127.0.0.1 rules-engine.local.dev
```

## Access URLs

After deployment, the services will be available at:

- **Rules Management API**: https://rules-management.local.dev
- **Rules Evaluation API**: https://rules-evaluation.local.dev
- **Rules Calculator API**: https://rules-calculator.local.dev
- **Rules Engine Gateway**: https://rules-engine.local.dev

## API Endpoints

### Rules Management Service

- `GET /health` - Health check
- `GET /v1/rules` - List rules
- `POST /v1/rules` - Create rule
- `GET /v1/rules/{id}` - Get rule by ID
- `PUT /v1/rules/{id}` - Update rule
- `DELETE /v1/rules/{id}` - Delete rule
- `POST /v1/rules/{id}/validate` - Validate rule DSL

### Rules Evaluation Service

- `GET /health` - Health check
- `POST /v1/evaluate` - Evaluate a rule

### Rules Calculator Service

- `GET /health` - Health check
- `POST /v1/calculate` - Calculate rules

## Monitoring and Observability

### Prometheus Metrics

All services expose Prometheus metrics on port 9090:

- `GET /metrics` - Prometheus metrics endpoint

### Health Checks

All services provide health check endpoints:

- `GET /health` - Application health status

### Logging

View logs using kubectl:

```bash
# View logs for all services
kubectl logs -f deployment/rules-management-service -n rules-engine
kubectl logs -f deployment/rules-evaluation-service -n rules-engine
kubectl logs -f deployment/rules-calculator-service -n rules-engine
```

## Troubleshooting

### Check Pod Status

```bash
kubectl get pods -n rules-engine
kubectl describe pod <pod-name> -n rules-engine
```

### Check Service Status

```bash
kubectl get svc -n rules-engine
kubectl describe svc <service-name> -n rules-engine
```

### Check Ingress Status

```bash
kubectl get mappings -n rules-engine
kubectl describe mapping <mapping-name> -n rules-engine
```

### Check Certificate Status

```bash
kubectl get certificates -n rules-engine
kubectl describe certificate rules-engine-tls -n rules-engine
```

### Database Connection Issues

```bash
# Check database initialization job
kubectl get jobs -n rules-engine
kubectl logs job/db-init-job -n rules-engine

# Test database connectivity
kubectl run -it --rm --restart=Never db-test --image=postgres:15-alpine -- psql -h postgresql.default.svc.cluster.local -U rules_user -d rules_engine
```

### Common Issues

1. **Certificate not ready**: Wait for cert-manager to issue the certificate
2. **Database connection failed**: Ensure PostgreSQL is running and accessible
3. **Image pull failed**: Check if the local registry is accessible
4. **Ingress not working**: Verify Emissary Ingress is running and configured

## Configuration

### Environment Variables

The services are configured using environment variables from ConfigMaps and Secrets:

- Database connection settings
- Service ports
- NATS configuration
- Telemetry settings

### Resource Limits

Each service is configured with:
- CPU: 100m request, 500m limit
- Memory: 256Mi request, 512Mi limit

### Security

- Services run as non-root user (UID 1000)
- Read-only root filesystem
- No privilege escalation
- All capabilities dropped

## Development

### Local Development

For local development, you can run the services directly:

```bash
# Rules Management Service
cd rules-management-service
go run ./cmd

# Rules Evaluation Service
cd rules-evaluation-service
go run ./cmd

# Rules Calculator Service
cd rules-calculator-service
go run ./cmd
```

### Testing

Test the deployed services:

```bash
# Health check
curl -k https://rules-management.local.dev/health

# List rules
curl -k https://rules-management.local.dev/v1/rules

# Create a rule
curl -k -X POST https://rules-management.local.dev/v1/rules \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Rule",
    "description": "A test rule",
    "dslContent": "IF customer.tier == \"VIP\" THEN discount.percentage = 20",
    "priority": "HIGH",
    "category": "test"
  }'
```

## Cleanup

To remove all resources:

```bash
./scripts/cleanup.sh
```

Or manually:

```bash
kubectl delete namespace rules-engine
```

## Support

For issues or questions, please check the troubleshooting section above or review the service logs.
