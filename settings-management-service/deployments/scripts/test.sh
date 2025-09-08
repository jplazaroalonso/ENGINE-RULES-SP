#!/bin/bash

# Settings Management Service Test Script
# This script tests the deployed Settings Management Service

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
BASE_URL="${BASE_URL:-https://settings-api.rulesengine.local}"
LOCAL_URL="${LOCAL_URL:-http://localhost:8080}"

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

# Check if curl is available
check_curl() {
    if ! command -v curl &> /dev/null; then
        log_error "curl is not installed or not in PATH"
        exit 1
    fi
    log_info "curl is available"
}

# Check if jq is available
check_jq() {
    if ! command -v jq &> /dev/null; then
        log_warning "jq is not installed. JSON responses will not be formatted."
        JQ_AVAILABLE=false
    else
        log_info "jq is available"
        JQ_AVAILABLE=true
    fi
}

# Test health endpoint
test_health() {
    log_info "Testing health endpoint..."
    
    local url="$1/health"
    local response=$(curl -s -w "%{http_code}" -o /tmp/health_response.json "$url")
    
    if [ "$response" = "200" ]; then
        log_success "Health endpoint is responding"
        if [ "$JQ_AVAILABLE" = true ]; then
            cat /tmp/health_response.json | jq .
        else
            cat /tmp/health_response.json
        fi
    else
        log_error "Health endpoint failed with status: $response"
        return 1
    fi
}

# Test metrics endpoint
test_metrics() {
    log_info "Testing metrics endpoint..."
    
    local url="$1/metrics"
    local response=$(curl -s -w "%{http_code}" -o /tmp/metrics_response.txt "$url")
    
    if [ "$response" = "200" ]; then
        log_success "Metrics endpoint is responding"
        # Check if response contains Prometheus metrics
        if grep -q "# HELP" /tmp/metrics_response.txt; then
            log_success "Prometheus metrics are available"
        else
            log_warning "Prometheus metrics format not detected"
        fi
    else
        log_error "Metrics endpoint failed with status: $response"
        return 1
    fi
}

# Test configuration endpoints
test_configurations() {
    log_info "Testing configuration endpoints..."
    
    local url="$1"
    
    # Test list configurations
    log_info "Testing list configurations..."
    local list_response=$(curl -s -w "%{http_code}" -o /tmp/list_configs.json "$url/v1/configurations")
    
    if [ "$list_response" = "200" ]; then
        log_success "List configurations endpoint is responding"
        if [ "$JQ_AVAILABLE" = true ]; then
            cat /tmp/list_configs.json | jq .
        fi
    else
        log_error "List configurations endpoint failed with status: $list_response"
        return 1
    fi
    
    # Test create configuration
    log_info "Testing create configuration..."
    local create_response=$(curl -s -w "%{http_code}" -o /tmp/create_config.json \
        -X POST \
        -H "Content-Type: application/json" \
        -d '{
            "key": "test.config",
            "value": {"host": "localhost", "port": 5432},
            "category": "test",
            "environment": "development",
            "description": "Test configuration",
            "is_encrypted": false
        }' \
        "$url/v1/configurations")
    
    if [ "$create_response" = "201" ]; then
        log_success "Create configuration endpoint is responding"
        if [ "$JQ_AVAILABLE" = true ]; then
            cat /tmp/create_config.json | jq .
        fi
        
        # Extract configuration ID for further tests
        if [ "$JQ_AVAILABLE" = true ]; then
            CONFIG_ID=$(cat /tmp/create_config.json | jq -r '.data.id')
        else
            # Fallback: extract ID using grep and sed
            CONFIG_ID=$(grep -o '"id":"[^"]*"' /tmp/create_config.json | sed 's/"id":"\([^"]*\)"/\1/')
        fi
        
        if [ -n "$CONFIG_ID" ] && [ "$CONFIG_ID" != "null" ]; then
            log_info "Created configuration with ID: $CONFIG_ID"
            
            # Test get configuration
            log_info "Testing get configuration..."
            local get_response=$(curl -s -w "%{http_code}" -o /tmp/get_config.json "$url/v1/configurations/$CONFIG_ID")
            
            if [ "$get_response" = "200" ]; then
                log_success "Get configuration endpoint is responding"
            else
                log_error "Get configuration endpoint failed with status: $get_response"
            fi
            
            # Test update configuration
            log_info "Testing update configuration..."
            local update_response=$(curl -s -w "%{http_code}" -o /tmp/update_config.json \
                -X PUT \
                -H "Content-Type: application/json" \
                -d '{
                    "value": {"host": "updated-host", "port": 3306},
                    "description": "Updated test configuration"
                }' \
                "$url/v1/configurations/$CONFIG_ID")
            
            if [ "$update_response" = "200" ]; then
                log_success "Update configuration endpoint is responding"
            else
                log_error "Update configuration endpoint failed with status: $update_response"
            fi
            
            # Test delete configuration
            log_info "Testing delete configuration..."
            local delete_response=$(curl -s -w "%{http_code}" -o /tmp/delete_config.json \
                -X DELETE \
                "$url/v1/configurations/$CONFIG_ID")
            
            if [ "$delete_response" = "204" ]; then
                log_success "Delete configuration endpoint is responding"
            else
                log_error "Delete configuration endpoint failed with status: $delete_response"
            fi
        else
            log_warning "Could not extract configuration ID for further tests"
        fi
    else
        log_error "Create configuration endpoint failed with status: $create_response"
        return 1
    fi
}

# Test feature flag endpoints
test_feature_flags() {
    log_info "Testing feature flag endpoints..."
    
    local url="$1"
    
    # Test list feature flags
    log_info "Testing list feature flags..."
    local list_response=$(curl -s -w "%{http_code}" -o /tmp/list_flags.json "$url/v1/feature-flags")
    
    if [ "$list_response" = "200" ]; then
        log_success "List feature flags endpoint is responding"
    else
        log_error "List feature flags endpoint failed with status: $list_response"
        return 1
    fi
    
    # Test create feature flag
    log_info "Testing create feature flag..."
    local create_response=$(curl -s -w "%{http_code}" -o /tmp/create_flag.json \
        -X POST \
        -H "Content-Type: application/json" \
        -d '{
            "name": "Test Feature",
            "key": "test-feature",
            "description": "Test feature flag",
            "status": "active",
            "rollout_percentage": 50,
            "target_audience": {"users": ["user1", "user2"]},
            "conditions": {"country": "US"}
        }' \
        "$url/v1/feature-flags")
    
    if [ "$create_response" = "201" ]; then
        log_success "Create feature flag endpoint is responding"
        if [ "$JQ_AVAILABLE" = true ]; then
            cat /tmp/create_flag.json | jq .
        fi
    else
        log_error "Create feature flag endpoint failed with status: $create_response"
        return 1
    fi
}

# Test user preference endpoints
test_user_preferences() {
    log_info "Testing user preference endpoints..."
    
    local url="$1"
    
    # Test list user preferences
    log_info "Testing list user preferences..."
    local list_response=$(curl -s -w "%{http_code}" -o /tmp/list_prefs.json "$url/v1/user-preferences")
    
    if [ "$list_response" = "200" ]; then
        log_success "List user preferences endpoint is responding"
    else
        log_error "List user preferences endpoint failed with status: $list_response"
        return 1
    fi
    
    # Test create user preference
    log_info "Testing create user preference..."
    local create_response=$(curl -s -w "%{http_code}" -o /tmp/create_pref.json \
        -X POST \
        -H "Content-Type: application/json" \
        -d '{
            "user_id": "123e4567-e89b-12d3-a456-426614174000",
            "key": "theme",
            "value": {"theme": "dark", "font_size": "medium"},
            "category": "ui",
            "description": "User interface preferences",
            "is_public": false
        }' \
        "$url/v1/user-preferences")
    
    if [ "$create_response" = "201" ]; then
        log_success "Create user preference endpoint is responding"
        if [ "$JQ_AVAILABLE" = true ]; then
            cat /tmp/create_pref.json | jq .
        fi
    else
        log_error "Create user preference endpoint failed with status: $create_response"
        return 1
    fi
}

# Test organization setting endpoints
test_organization_settings() {
    log_info "Testing organization setting endpoints..."
    
    local url="$1"
    
    # Test list organization settings
    log_info "Testing list organization settings..."
    local list_response=$(curl -s -w "%{http_code}" -o /tmp/list_org_settings.json "$url/v1/organization-settings")
    
    if [ "$list_response" = "200" ]; then
        log_success "List organization settings endpoint is responding"
    else
        log_error "List organization settings endpoint failed with status: $list_response"
        return 1
    fi
    
    # Test create organization setting
    log_info "Testing create organization setting..."
    local create_response=$(curl -s -w "%{http_code}" -o /tmp/create_org_setting.json \
        -X POST \
        -H "Content-Type: application/json" \
        -d '{
            "organization_id": "123e4567-e89b-12d3-a456-426614174000",
            "key": "billing.plan",
            "value": {"plan": "enterprise", "features": ["advanced_analytics", "custom_branding"]},
            "category": "billing",
            "description": "Organization billing plan",
            "is_encrypted": false
        }' \
        "$url/v1/organization-settings")
    
    if [ "$create_response" = "201" ]; then
        log_success "Create organization setting endpoint is responding"
        if [ "$JQ_AVAILABLE" = true ]; then
            cat /tmp/create_org_setting.json | jq .
        fi
    else
        log_error "Create organization setting endpoint failed with status: $create_response"
        return 1
    fi
}

# Test load balancing
test_load_balancing() {
    log_info "Testing load balancing..."
    
    local url="$1"
    local requests=10
    local success_count=0
    
    for i in $(seq 1 $requests); do
        local response=$(curl -s -w "%{http_code}" -o /dev/null "$url/health")
        if [ "$response" = "200" ]; then
            ((success_count++))
        fi
    done
    
    local success_rate=$((success_count * 100 / requests))
    log_info "Load balancing test: $success_count/$requests requests successful ($success_rate%)"
    
    if [ $success_rate -ge 90 ]; then
        log_success "Load balancing is working well"
    else
        log_warning "Load balancing may have issues"
    fi
}

# Test error handling
test_error_handling() {
    log_info "Testing error handling..."
    
    local url="$1"
    
    # Test 404 error
    log_info "Testing 404 error..."
    local not_found_response=$(curl -s -w "%{http_code}" -o /tmp/not_found.json "$url/v1/configurations/non-existent-id")
    
    if [ "$not_found_response" = "404" ]; then
        log_success "404 error handling is working"
    else
        log_error "404 error handling failed with status: $not_found_response"
    fi
    
    # Test 400 error (invalid JSON)
    log_info "Testing 400 error..."
    local bad_request_response=$(curl -s -w "%{http_code}" -o /tmp/bad_request.json \
        -X POST \
        -H "Content-Type: application/json" \
        -d '{"invalid": json}' \
        "$url/v1/configurations")
    
    if [ "$bad_request_response" = "400" ]; then
        log_success "400 error handling is working"
    else
        log_error "400 error handling failed with status: $bad_request_response"
    fi
}

# Cleanup function
cleanup() {
    log_info "Cleaning up temporary files..."
    rm -f /tmp/health_response.json
    rm -f /tmp/metrics_response.txt
    rm -f /tmp/list_configs.json
    rm -f /tmp/create_config.json
    rm -f /tmp/get_config.json
    rm -f /tmp/update_config.json
    rm -f /tmp/delete_config.json
    rm -f /tmp/list_flags.json
    rm -f /tmp/create_flag.json
    rm -f /tmp/list_prefs.json
    rm -f /tmp/create_pref.json
    rm -f /tmp/list_org_settings.json
    rm -f /tmp/create_org_setting.json
    rm -f /tmp/not_found.json
    rm -f /tmp/bad_request.json
}

# Main test function
main() {
    log_info "Starting tests for $SERVICE_NAME..."
    
    # Change to script directory
    cd "$(dirname "$0")"
    
    # Check prerequisites
    check_curl
    check_jq
    
    # Determine which URL to use
    local test_url
    if [ "$1" = "local" ]; then
        test_url="$LOCAL_URL"
        log_info "Testing against local URL: $test_url"
    else
        test_url="$BASE_URL"
        log_info "Testing against production URL: $test_url"
    fi
    
    # Run tests
    test_health "$test_url"
    test_metrics "$test_url"
    test_configurations "$test_url"
    test_feature_flags "$test_url"
    test_user_preferences "$test_url"
    test_organization_settings "$test_url"
    test_load_balancing "$test_url"
    test_error_handling "$test_url"
    
    log_success "All tests completed!"
}

# Handle script interruption
trap cleanup EXIT

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --base-url)
            BASE_URL="$2"
            shift 2
            ;;
        --local-url)
            LOCAL_URL="$2"
            shift 2
            ;;
        --help)
            echo "Usage: $0 [OPTIONS] [local]"
            echo "Options:"
            echo "  --base-url URL     Base URL for testing (default: https://settings-api.rulesengine.local)"
            echo "  --local-url URL    Local URL for testing (default: http://localhost:8080)"
            echo "  --help             Show this help message"
            echo ""
            echo "Arguments:"
            echo "  local              Test against local URL instead of production URL"
            exit 0
            ;;
        *)
            if [ "$1" = "local" ]; then
                shift
            else
                log_error "Unknown option: $1"
                exit 1
            fi
            ;;
    esac
done

# Run main function
main "$@"
