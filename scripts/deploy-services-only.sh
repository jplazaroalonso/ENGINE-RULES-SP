#!/bin/bash

# Rules Engine - Deploy Services Only Script
# This script deploys only the services (assumes infrastructure is already set up)

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
REGISTRY="localhost:5000"
NAMESPACE="rules-engine"
SERVICES=("rules-management-service" "rules-evaluation-service" "rules-calculator-service")

echo -e "${BLUE}🚀 Starting Rules Engine services deployment...${NC}"

# Function to print status
print_status() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Check if namespace exists
check_namespace() {
    echo -e "${BLUE}🔍 Checking if namespace exists...${NC}"
    
    if ! kubectl get namespace ${NAMESPACE} &> /dev/null; then
        print_error "Namespace ${NAMESPACE} does not exist. Please run deploy-from-scratch.sh first."
        exit 1
    fi
    
    print_status "Namespace ${NAMESPACE} exists"
}

# Build Docker images
build_images() {
    echo -e "${BLUE}🔨 Building Docker images...${NC}"
    
    for service in "${SERVICES[@]}"; do
        echo -e "${BLUE}Building ${service}...${NC}"
        
        # Build the image from service directory
        nerdctl build -t "${REGISTRY}/${service}:latest" "./${service}/"
        
        # Push to registry
        nerdctl push "${REGISTRY}/${service}:latest"
        
        print_status "Built and pushed ${service}"
    done
}

# Deploy services
deploy_services() {
    echo -e "${BLUE}🚀 Deploying services...${NC}"
    
    # Deploy all services
    kubectl apply -f k8s/rules-management-service.yaml
    kubectl apply -f k8s/rules-evaluation-service.yaml
    kubectl apply -f k8s/rules-calculator-service.yaml
    
    # Wait for deployments to be ready
    for service in "${SERVICES[@]}"; do
        echo -e "${BLUE}Waiting for ${service} to be ready...${NC}"
        kubectl wait --for=condition=available deployment/${service} -n ${NAMESPACE} --timeout=300s
        print_status "${service} is ready"
    done
}

# Verify deployment
verify_deployment() {
    echo -e "${BLUE}🔍 Verifying deployment...${NC}"
    
    # Check if all pods are running
    for service in "${SERVICES[@]}"; do
        if kubectl get pods -n ${NAMESPACE} -l app=${service} --field-selector=status.phase=Running | grep -q ${service}; then
            print_status "${service} pod is running"
        else
            print_error "${service} pod is not running"
            kubectl get pods -n ${NAMESPACE} -l app=${service}
        fi
    done
    
    # Check services
    if kubectl get svc -n ${NAMESPACE} | grep -q "rules-management-service"; then
        print_status "Services are created"
    else
        print_error "Services are not created"
    fi
}

# Print access information
print_access_info() {
    echo -e "${GREEN}🎉 Services deployment completed successfully!${NC}"
    echo -e "${BLUE}📋 Access Information:${NC}"
    echo -e "  • Rules Management API: https://rules-management.local.dev"
    echo -e "  • Rules Evaluation API: https://rules-evaluation.local.dev"
    echo -e "  • Rules Calculator API: https://rules-calculator.local.dev"
    echo -e "  • Rules Engine Gateway: https://rules-engine.local.dev"
    echo ""
    echo -e "${BLUE}🔧 Useful Commands:${NC}"
    echo -e "  • View pods: kubectl get pods -n ${NAMESPACE}"
    echo -e "  • View services: kubectl get svc -n ${NAMESPACE}"
    echo -e "  • View logs: kubectl logs -f deployment/rules-management-service -n ${NAMESPACE}"
    echo ""
    echo -e "${BLUE}🧪 Test the APIs:${NC}"
    echo -e "  • Health check: curl -k https://rules-management.local.dev/health"
    echo -e "  • List rules: curl -k https://rules-management.local.dev/v1/rules"
}

# Main execution
main() {
    check_namespace
    build_images
    deploy_services
    verify_deployment
    print_access_info
}

# Run main function
main "$@"
