#!/bin/bash

# Deployment script for Rules Engine Frontend using nerdctl
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
NAMESPACE="rules-engine-frontend"
IMAGE_NAME="rules-engine-frontend"
IMAGE_TAG="latest"
REGISTRY="localhost:5000"

# Functions
print_status() {
    local status=$1
    local message=$2
    if [ $status -eq 0 ]; then
        echo -e "${GREEN}✓${NC} $message"
    else
        echo -e "${RED}✗${NC} $message"
        exit 1
    fi
}

print_info() {
    echo -e "${BLUE}ℹ${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

# Check if kubectl is available
check_kubectl() {
    print_info "Checking kubectl..."
    if ! command -v kubectl &> /dev/null; then
        print_status 1 "kubectl is not installed or not in PATH"
    fi
    
    if ! kubectl cluster-info &> /dev/null; then
        print_status 1 "Cannot connect to Kubernetes cluster"
    fi
    
    print_status 0 "kubectl is available and cluster is accessible"
}

# Check if nerdctl is available
check_nerdctl() {
    print_info "Checking nerdctl..."
    if ! command -v nerdctl &> /dev/null; then
        print_status 1 "nerdctl is not installed or not in PATH"
    fi
    print_status 0 "nerdctl is available"
}

# Check if namespace exists
check_namespace() {
    print_info "Checking namespace..."
    if kubectl get namespace "$NAMESPACE" &> /dev/null; then
        print_info "Namespace $NAMESPACE already exists"
    else
        print_info "Creating namespace $NAMESPACE..."
        kubectl apply -f k8s/namespace.yaml
        print_status 0 "Namespace created"
    fi
}

# Deploy ConfigMap
deploy_configmap() {
    print_info "Deploying ConfigMap..."
    kubectl apply -f k8s/configmap.yaml
    print_status 0 "ConfigMap deployed"
}

# Deploy Service
deploy_service() {
    print_info "Deploying Service..."
    kubectl apply -f k8s/service.yaml
    print_status 0 "Service deployed"
}

# Deploy Deployment
deploy_deployment() {
    print_info "Deploying Deployment..."
    kubectl apply -f k8s/deployment.yaml
    print_status 0 "Deployment deployed"
}

# Deploy Ingress
deploy_ingress() {
    print_info "Deploying Ingress..."
    kubectl apply -f k8s/ingress.yaml
    print_status 0 "Ingress deployed"
}

# Wait for deployment to be ready
wait_for_deployment() {
    print_info "Waiting for deployment to be ready..."
    if kubectl wait --for=condition=available --timeout=300s deployment/rules-engine-frontend -n "$NAMESPACE"; then
        print_status 0 "Deployment is ready"
    else
        print_status 1 "Deployment failed to become ready"
    fi
}

# Check deployment status
check_deployment_status() {
    print_info "Checking deployment status..."
    
    # Get pod status
    local pods=$(kubectl get pods -n "$NAMESPACE" -l app.kubernetes.io/name=rules-engine-frontend --no-headers)
    if [ -z "$pods" ]; then
        print_status 1 "No pods found for the deployment"
    fi
    
    echo "$pods" | while read -r pod_name ready status restarts age; do
        if [ "$status" = "Running" ] && [[ "$ready" =~ ^[0-9]+/[0-9]+$ ]]; then
            print_status 0 "Pod $pod_name is running"
        else
            print_warning "Pod $pod_name is not ready: $status"
        fi
    done
}

# Get service URLs
get_service_urls() {
    print_info "Service URLs:"
    echo ""
    
    # Get ingress URL
    local ingress_host=$(kubectl get ingress rules-engine-frontend -n "$NAMESPACE" -o jsonpath='{.spec.rules[0].host}' 2>/dev/null || echo "N/A")
    if [ "$ingress_host" != "N/A" ]; then
        echo -e "  ${GREEN}Frontend URL:${NC} https://$ingress_host"
    fi
    
    # Get service URL
    local service_ip=$(kubectl get service rules-engine-frontend -n "$NAMESPACE" -o jsonpath='{.spec.clusterIP}' 2>/dev/null || echo "N/A")
    if [ "$service_ip" != "N/A" ]; then
        echo -e "  ${GREEN}Service IP:${NC} $service_ip"
    fi
    
    echo ""
}

# Show logs
show_logs() {
    print_info "Showing recent logs..."
    kubectl logs -n "$NAMESPACE" -l app.kubernetes.io/name=rules-engine-frontend --tail=20
}

# Clean up deployment
cleanup() {
    print_info "Cleaning up deployment..."
    kubectl delete -f k8s/ingress.yaml --ignore-not-found=true
    kubectl delete -f k8s/deployment.yaml --ignore-not-found=true
    kubectl delete -f k8s/service.yaml --ignore-not-found=true
    kubectl delete -f k8s/configmap.yaml --ignore-not-found=true
    kubectl delete -f k8s/namespace.yaml --ignore-not-found=true
    print_status 0 "Cleanup completed"
}

# Main execution
main() {
    echo -e "${BLUE}🚀 Deploying Rules Engine Frontend with nerdctl${NC}"
    echo "============================================="
    
    # Parse command line arguments
    CLEANUP=false
    SHOW_LOGS=false
    SKIP_WAIT=false
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            --cleanup)
                CLEANUP=true
                shift
                ;;
            --logs)
                SHOW_LOGS=true
                shift
                ;;
            --skip-wait)
                SKIP_WAIT=true
                shift
                ;;
            --help)
                echo "Usage: $0 [OPTIONS]"
                echo "Options:"
                echo "  --cleanup       Clean up existing deployment"
                echo "  --logs          Show logs after deployment"
                echo "  --skip-wait     Skip waiting for deployment to be ready"
                echo "  --help          Show this help message"
                exit 0
                ;;
            *)
                echo "Unknown option: $1"
                exit 1
                ;;
        esac
    done
    
    if [ "$CLEANUP" = true ]; then
        cleanup
        exit 0
    fi
    
    # Execute deployment steps
    check_kubectl
    check_nerdctl
    check_namespace
    deploy_configmap
    deploy_service
    deploy_deployment
    deploy_ingress
    
    if [ "$SKIP_WAIT" = false ]; then
        wait_for_deployment
    fi
    
    check_deployment_status
    get_service_urls
    
    if [ "$SHOW_LOGS" = true ]; then
        show_logs
    fi
    
    echo ""
    echo -e "${GREEN}🎉 Deployment completed successfully!${NC}"
    echo ""
    echo "Access the frontend at: https://rules-engine.local.dev"
    echo ""
    echo "Useful commands:"
    echo "  kubectl get pods -n $NAMESPACE"
    echo "  kubectl logs -n $NAMESPACE -l app.kubernetes.io/name=rules-engine-frontend"
    echo "  kubectl describe ingress rules-engine-frontend -n $NAMESPACE"
}

# Run main function
main "$@"
