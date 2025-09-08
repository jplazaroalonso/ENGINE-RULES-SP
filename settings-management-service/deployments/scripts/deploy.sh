#!/bin/bash

# Settings Management Service Deployment Script
# This script deploys the Settings Management Service to Kubernetes

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
NAMESPACE="rules-engine"
SERVICE_NAME="settings-management-service"
IMAGE_TAG="${IMAGE_TAG:-latest}"
REPLICAS="${REPLICAS:-3}"

# Functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if kubectl is available
check_kubectl() {
    if ! command -v kubectl &> /dev/null; then
        log_error "kubectl is not installed or not in PATH"
        exit 1
    fi
    log_info "kubectl is available"
}

# Check if namespace exists
check_namespace() {
    if kubectl get namespace "$NAMESPACE" &> /dev/null; then
        log_info "Namespace '$NAMESPACE' exists"
    else
        log_warning "Namespace '$NAMESPACE' does not exist. Creating it..."
        kubectl create namespace "$NAMESPACE"
        log_success "Namespace '$NAMESPACE' created"
    fi
}

# Build Docker image
build_image() {
    log_info "Building Docker image for $SERVICE_NAME..."
    
    # Check if Docker is available
    if ! command -v docker &> /dev/null; then
        log_error "Docker is not installed or not in PATH"
        exit 1
    fi
    
    # Build the image
    docker build -t "$SERVICE_NAME:$IMAGE_TAG" .
    
    if [ $? -eq 0 ]; then
        log_success "Docker image built successfully: $SERVICE_NAME:$IMAGE_TAG"
    else
        log_error "Failed to build Docker image"
        exit 1
    fi
}

# Deploy to Kubernetes
deploy_to_k8s() {
    log_info "Deploying $SERVICE_NAME to Kubernetes..."
    
    # Update image tag in deployment
    sed -i.bak "s|image: $SERVICE_NAME:.*|image: $SERVICE_NAME:$IMAGE_TAG|g" k8s/deployment.yaml
    sed -i.bak "s|replicas: .*|replicas: $REPLICAS|g" k8s/deployment.yaml
    
    # Apply configurations in order
    log_info "Applying ConfigMap..."
    kubectl apply -f k8s/configmap.yaml -n "$NAMESPACE"
    
    log_info "Applying Secret..."
    kubectl apply -f k8s/secret.yaml -n "$NAMESPACE"
    
    log_info "Applying Deployment..."
    kubectl apply -f k8s/deployment.yaml -n "$NAMESPACE"
    
    log_info "Applying Service..."
    kubectl apply -f k8s/service.yaml -n "$NAMESPACE"
    
    log_info "Applying Ingress..."
    kubectl apply -f k8s/ingress.yaml -n "$NAMESPACE"
    
    log_info "Applying HPA..."
    kubectl apply -f k8s/hpa.yaml -n "$NAMESPACE"
    
    log_info "Applying PodDisruptionBudget..."
    kubectl apply -f k8s/pdb.yaml -n "$NAMESPACE"
    
    log_info "Applying NetworkPolicy..."
    kubectl apply -f k8s/networkpolicy.yaml -n "$NAMESPACE"
    
    # Apply monitoring resources if available
    if kubectl get crd servicemonitors.monitoring.coreos.com &> /dev/null; then
        log_info "Applying ServiceMonitor..."
        kubectl apply -f k8s/servicemonitor.yaml -n "$NAMESPACE"
    else
        log_warning "ServiceMonitor CRD not found. Skipping monitoring setup."
    fi
    
    if kubectl get crd prometheusrules.monitoring.coreos.com &> /dev/null; then
        log_info "Applying PrometheusRule..."
        kubectl apply -f k8s/prometheusrule.yaml -n "$NAMESPACE"
    else
        log_warning "PrometheusRule CRD not found. Skipping alerting setup."
    fi
    
    log_success "All resources applied successfully"
}

# Wait for deployment to be ready
wait_for_deployment() {
    log_info "Waiting for deployment to be ready..."
    
    kubectl wait --for=condition=available --timeout=300s deployment/"$SERVICE_NAME" -n "$NAMESPACE"
    
    if [ $? -eq 0 ]; then
        log_success "Deployment is ready"
    else
        log_error "Deployment failed to become ready within 5 minutes"
        exit 1
    fi
}

# Check deployment status
check_deployment_status() {
    log_info "Checking deployment status..."
    
    kubectl get deployment "$SERVICE_NAME" -n "$NAMESPACE"
    kubectl get pods -l app="$SERVICE_NAME" -n "$NAMESPACE"
    kubectl get service "$SERVICE_NAME" -n "$NAMESPACE"
    kubectl get ingress "$SERVICE_NAME" -n "$NAMESPACE"
}

# Show logs
show_logs() {
    log_info "Showing recent logs..."
    
    kubectl logs -l app="$SERVICE_NAME" -n "$NAMESPACE" --tail=50
}

# Cleanup function
cleanup() {
    log_info "Cleaning up temporary files..."
    rm -f k8s/deployment.yaml.bak
}

# Main deployment function
main() {
    log_info "Starting deployment of $SERVICE_NAME..."
    
    # Change to script directory
    cd "$(dirname "$0")"
    
    # Run deployment steps
    check_kubectl
    check_namespace
    build_image
    deploy_to_k8s
    wait_for_deployment
    check_deployment_status
    
    log_success "Deployment completed successfully!"
    
    # Show service information
    log_info "Service information:"
    echo "  Namespace: $NAMESPACE"
    echo "  Service: $SERVICE_NAME"
    echo "  Image: $SERVICE_NAME:$IMAGE_TAG"
    echo "  Replicas: $REPLICAS"
    echo "  Ingress: https://settings-api.rulesengine.local"
    
    # Cleanup
    cleanup
}

# Handle script interruption
trap cleanup EXIT

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --image-tag)
            IMAGE_TAG="$2"
            shift 2
            ;;
        --replicas)
            REPLICAS="$2"
            shift 2
            ;;
        --skip-build)
            SKIP_BUILD=true
            shift
            ;;
        --help)
            echo "Usage: $0 [OPTIONS]"
            echo "Options:"
            echo "  --image-tag TAG    Docker image tag (default: latest)"
            echo "  --replicas NUM     Number of replicas (default: 3)"
            echo "  --skip-build       Skip Docker image build"
            echo "  --help             Show this help message"
            exit 0
            ;;
        *)
            log_error "Unknown option: $1"
            exit 1
            ;;
    esac
done

# Run main function
main
