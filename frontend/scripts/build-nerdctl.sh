#!/bin/bash

# Build script for Rules Engine Frontend using nerdctl
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
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

# Check if nerdctl is available
check_nerdctl() {
    print_info "Checking nerdctl..."
    if ! command -v nerdctl &> /dev/null; then
        print_status 1 "nerdctl is not installed or not in PATH"
    fi
    print_status 0 "nerdctl is available"
}

# Install dependencies
install_dependencies() {
    print_info "Installing dependencies..."
    if [ ! -d "node_modules" ]; then
        npm ci
        print_status 0 "Dependencies installed"
    else
        print_info "Dependencies already installed, skipping..."
    fi
}

# Run tests
run_tests() {
    print_info "Running tests..."
    if npm run test:run; then
        print_status 0 "Tests passed"
    else
        print_status 1 "Tests failed"
    fi
}

# Run linting
run_lint() {
    print_info "Running linter..."
    if npm run lint; then
        print_status 0 "Linting passed"
    else
        print_warning "Linting issues found, but continuing..."
    fi
}

# Build the application
build_app() {
    print_info "Building application..."
    if npm run build; then
        print_status 0 "Application built successfully"
    else
        print_status 1 "Build failed"
    fi
}

# Build container image with nerdctl
build_container_image() {
    print_info "Building container image with nerdctl..."
    local full_image_name="${REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG}"
    
    if nerdctl build -t "$full_image_name" .; then
        print_status 0 "Container image built: $full_image_name"
    else
        print_status 1 "Container image build failed"
    fi
}

# Push container image
push_container_image() {
    print_info "Pushing container image to registry..."
    local full_image_name="${REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG}"
    
    if nerdctl push "$full_image_name"; then
        print_status 0 "Container image pushed: $full_image_name"
    else
        print_status 1 "Failed to push container image"
    fi
}

# Clean up
cleanup() {
    print_info "Cleaning up..."
    rm -rf dist/
    print_status 0 "Cleanup completed"
}

# Main execution
main() {
    echo -e "${BLUE}🚀 Building Rules Engine Frontend with nerdctl${NC}"
    echo "============================================="
    
    # Parse command line arguments
    SKIP_TESTS=false
    SKIP_LINT=false
    PUSH_IMAGE=false
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            --skip-tests)
                SKIP_TESTS=true
                shift
                ;;
            --skip-lint)
                SKIP_LINT=true
                shift
                ;;
            --push)
                PUSH_IMAGE=true
                shift
                ;;
            --help)
                echo "Usage: $0 [OPTIONS]"
                echo "Options:"
                echo "  --skip-tests    Skip running tests"
                echo "  --skip-lint     Skip running linter"
                echo "  --push          Push image to registry after building"
                echo "  --help          Show this help message"
                exit 0
                ;;
            *)
                echo "Unknown option: $1"
                exit 1
                ;;
        esac
    done
    
    # Execute build steps
    check_nerdctl
    install_dependencies
    
    if [ "$SKIP_TESTS" = false ]; then
        run_tests
    else
        print_warning "Skipping tests"
    fi
    
    if [ "$SKIP_LINT" = false ]; then
        run_lint
    else
        print_warning "Skipping linting"
    fi
    
    build_app
    build_container_image
    
    if [ "$PUSH_IMAGE" = true ]; then
        push_container_image
    fi
    
    cleanup
    
    echo ""
    echo -e "${GREEN}🎉 Build completed successfully!${NC}"
    echo ""
    echo "Next steps:"
    echo "1. Deploy to Kubernetes: ./scripts/deploy-nerdctl.sh"
    echo "2. Or run locally: nerdctl run -p 3000:3000 ${REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG}"
}

# Run main function
main "$@"
