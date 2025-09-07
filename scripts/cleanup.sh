#!/bin/bash

# Rules Engine - Cleanup Script
# This script removes all resources created by the rules engine deployment

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
NAMESPACE="rules-engine"

echo -e "${BLUE}🧹 Starting Rules Engine cleanup...${NC}"

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

# Remove ingress resources
remove_ingress() {
    echo -e "${BLUE}🌐 Removing ingress resources...${NC}"
    
    if kubectl get mappings -n ${NAMESPACE} &> /dev/null; then
        kubectl delete mappings -n ${NAMESPACE} --all
        print_status "Ingress mappings removed"
    else
        print_warning "No ingress mappings found"
    fi
    
    if kubectl get certificates -n ${NAMESPACE} &> /dev/null; then
        kubectl delete certificates -n ${NAMESPACE} --all
        print_status "Certificates removed"
    else
        print_warning "No certificates found"
    fi
}

# Remove services
remove_services() {
    echo -e "${BLUE}🚀 Removing services...${NC}"
    
    if kubectl get deployments -n ${NAMESPACE} &> /dev/null; then
        kubectl delete deployments -n ${NAMESPACE} --all
        print_status "Deployments removed"
    else
        print_warning "No deployments found"
    fi
    
    if kubectl get svc -n ${NAMESPACE} &> /dev/null; then
        kubectl delete svc -n ${NAMESPACE} --all
        print_status "Services removed"
    else
        print_warning "No services found"
    fi
}

# Remove jobs
remove_jobs() {
    echo -e "${BLUE}📋 Removing jobs...${NC}"
    
    if kubectl get jobs -n ${NAMESPACE} &> /dev/null; then
        kubectl delete jobs -n ${NAMESPACE} --all
        print_status "Jobs removed"
    else
        print_warning "No jobs found"
    fi
}

# Remove configmaps and secrets
remove_configs() {
    echo -e "${BLUE}⚙️  Removing configmaps and secrets...${NC}"
    
    if kubectl get configmaps -n ${NAMESPACE} &> /dev/null; then
        kubectl delete configmaps -n ${NAMESPACE} --all
        print_status "ConfigMaps removed"
    else
        print_warning "No ConfigMaps found"
    fi
    
    if kubectl get secrets -n ${NAMESPACE} &> /dev/null; then
        kubectl delete secrets -n ${NAMESPACE} --all
        print_status "Secrets removed"
    else
        print_warning "No secrets found"
    fi
}

# Remove namespace
remove_namespace() {
    echo -e "${BLUE}📦 Removing namespace...${NC}"
    
    if kubectl get namespace ${NAMESPACE} &> /dev/null; then
        kubectl delete namespace ${NAMESPACE}
        print_status "Namespace ${NAMESPACE} removed"
    else
        print_warning "Namespace ${NAMESPACE} not found"
    fi
}

# Clean up /etc/hosts
cleanup_hosts() {
    echo -e "${BLUE}📝 Cleaning up /etc/hosts...${NC}"
    
    HOSTS_ENTRIES=(
        "rules-management.local.dev"
        "rules-evaluation.local.dev"
        "rules-calculator.local.dev"
        "rules-engine.local.dev"
    )
    
    for host in "${HOSTS_ENTRIES[@]}"; do
        if grep -q "${host}" /etc/hosts; then
            sudo sed -i '' "/${host}/d" /etc/hosts
            print_status "Removed ${host} from /etc/hosts"
        else
            print_warning "${host} not found in /etc/hosts"
        fi
    done
}

# Main execution
main() {
    echo -e "${YELLOW}⚠️  This will remove all rules engine resources. Are you sure? (y/N)${NC}"
    read -r response
    if [[ ! "$response" =~ ^[Yy]$ ]]; then
        echo -e "${BLUE}Cleanup cancelled.${NC}"
        exit 0
    fi
    
    remove_ingress
    remove_services
    remove_jobs
    remove_configs
    remove_namespace
    cleanup_hosts
    
    echo -e "${GREEN}🎉 Cleanup completed successfully!${NC}"
}

# Run main function
main "$@"
