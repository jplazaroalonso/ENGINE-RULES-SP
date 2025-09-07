#!/bin/bash

# Rules Engine - Deploy from Scratch Script
# This script builds and deploys the entire rules engine system from scratch

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

echo -e "${BLUE}🚀 Starting Rules Engine deployment from scratch...${NC}"

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

# Check if required tools are installed
check_prerequisites() {
    echo -e "${BLUE}🔍 Checking prerequisites...${NC}"
    
    if ! command -v nerdctl &> /dev/null; then
        print_error "nerdctl is not installed. Please install nerdctl."
        exit 1
    fi
    
    if ! command -v kubectl &> /dev/null; then
        print_error "kubectl is not installed. Please install kubectl."
        exit 1
    fi
    
    # Check if kubectl can connect to cluster
    if ! kubectl cluster-info &> /dev/null; then
        print_error "Cannot connect to Kubernetes cluster. Please check your kubeconfig."
        exit 1
    fi
    
    print_status "Prerequisites check passed"
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

# Create namespace and basic resources
setup_namespace() {
    echo -e "${BLUE}📦 Setting up namespace and basic resources...${NC}"
    
    # Apply namespace and basic resources
    kubectl apply -f k8s/namespace.yaml
    kubectl apply -f k8s/configmap.yaml
    kubectl apply -f k8s/secrets.yaml
    
    print_status "Namespace and basic resources created"
}

# Initialize database
init_database() {
    echo -e "${BLUE}🗄️  Initializing database...${NC}"
    
    # Apply database initialization job
    kubectl apply -f k8s/db-init-job.yaml
    
    # Wait for job to complete
    echo -e "${BLUE}Waiting for database initialization to complete...${NC}"
    kubectl wait --for=condition=complete job/db-init-job -n ${NAMESPACE} --timeout=300s
    
    print_status "Database initialized"
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

# Setup ingress
setup_ingress() {
    echo -e "${BLUE}🌐 Setting up ingress...${NC}"
    
    # Apply ingress configuration
    kubectl apply -f k8s/ingress.yaml
    
    # Wait for certificate to be ready
    echo -e "${BLUE}Waiting for TLS certificate to be ready...${NC}"
    kubectl wait --for=condition=ready certificate/rules-engine-tls -n ${NAMESPACE} --timeout=300s
    
    print_status "Ingress and TLS certificates configured"
}

# Update /etc/hosts
update_hosts() {
    echo -e "${BLUE}📝 Updating /etc/hosts...${NC}"
    
    # Get the ingress IP
    INGRESS_IP=$(kubectl get svc -n ambassador ambassador --output jsonpath='{.status.loadBalancer.ingress[0].ip}')
    
    if [ -z "$INGRESS_IP" ]; then
        INGRESS_IP="127.0.0.1"
        print_warning "Could not get ingress IP, using localhost"
    fi
    
    # Add entries to /etc/hosts
    HOSTS_ENTRIES=(
        "rules-management.local.dev"
        "rules-evaluation.local.dev"
        "rules-calculator.local.dev"
        "rules-engine.local.dev"
    )
    
    for host in "${HOSTS_ENTRIES[@]}"; do
        if ! grep -q "${host}" /etc/hosts; then
            echo "${INGRESS_IP} ${host}" | sudo tee -a /etc/hosts
            print_status "Added ${host} to /etc/hosts"
        else
            print_warning "${host} already exists in /etc/hosts"
        fi
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
    
    # Check ingress
    if kubectl get mappings -n ${NAMESPACE} | grep -q "rules-management-mapping"; then
        print_status "Ingress mappings are created"
    else
        print_error "Ingress mappings are not created"
    fi
}

# Print access information
print_access_info() {
    echo -e "${GREEN}🎉 Deployment completed successfully!${NC}"
    echo -e "${BLUE}📋 Access Information:${NC}"
    echo -e "  • Rules Management API: https://rules-management.local.dev"
    echo -e "  • Rules Evaluation API: https://rules-evaluation.local.dev"
    echo -e "  • Rules Calculator API: https://rules-calculator.local.dev"
    echo -e "  • Rules Engine Gateway: https://rules-engine.local.dev"
    echo ""
    echo -e "${BLUE}🔧 Useful Commands:${NC}"
    echo -e "  • View pods: kubectl get pods -n ${NAMESPACE}"
    echo -e "  • View services: kubectl get svc -n ${NAMESPACE}"
    echo -e "  • View ingress: kubectl get mappings -n ${NAMESPACE}"
    echo -e "  • View logs: kubectl logs -f deployment/rules-management-service -n ${NAMESPACE}"
    echo ""
    echo -e "${BLUE}🧪 Test the APIs:${NC}"
    echo -e "  • Health check: curl -k https://rules-management.local.dev/health"
    echo -e "  • List rules: curl -k https://rules-management.local.dev/v1/rules"
}

# Main execution
main() {
    check_prerequisites
    build_images
    setup_namespace
    init_database
    deploy_services
    setup_ingress
    update_hosts
    verify_deployment
    print_access_info
}

# Run main function
main "$@"
