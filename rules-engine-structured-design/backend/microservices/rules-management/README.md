# Rules Management Service

## Overview
The Rules Management Service handles the complete lifecycle of business rules including creation, editing, approval workflows, templating, and versioning. This service serves as the entry point for business users to manage rules.

## Domain Model

### Core Entities

#### Rule Aggregate
```go
type Rule struct {
    ID          string    `json:"id" gorm:"primaryKey"`
    Name        string    `json:"name" gorm:"not null"`
    Description string    `json:"description"`
    DSLContent  string    `json:"dsl_content" gorm:"type:text"`
    Status      RuleStatus `json:"status" gorm:"default:'DRAFT'"`
    Priority    Priority  `json:"priority" gorm:"default:'MEDIUM'"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    CreatedBy   string    `json:"created_by"`
    ApprovedBy  *string   `json:"approved_by,omitempty"`
    ApprovedAt  *time.Time `json:"approved_at,omitempty"`
    Version     int       `json:"version" gorm:"default:1"`
    
    // Relationships
    TemplateID  *string      `json:"template_id,omitempty"`
    Template    *RuleTemplate `json:"template,omitempty" gorm:"foreignKey:TemplateID"`
    Versions    []RuleVersion `json:"versions,omitempty" gorm:"foreignKey:RuleID"`
}

type RuleStatus string
const (
    StatusDraft      RuleStatus = "DRAFT"
    StatusUnderReview RuleStatus = "UNDER_REVIEW"
    StatusApproved   RuleStatus = "APPROVED"
    StatusActive     RuleStatus = "ACTIVE"
    StatusInactive   RuleStatus = "INACTIVE"
    StatusDeprecated RuleStatus = "DEPRECATED"
)

type Priority string
const (
    PriorityLow      Priority = "LOW"
    PriorityMedium   Priority = "MEDIUM"
    PriorityHigh     Priority = "HIGH"
    PriorityCritical Priority = "CRITICAL"
)
```

#### Rule Template
```go
type RuleTemplate struct {
    ID          string    `json:"id" gorm:"primaryKey"`
    Name        string    `json:"name" gorm:"not null"`
    Description string    `json:"description"`
    Category    string    `json:"category"`
    DSLTemplate string    `json:"dsl_template" gorm:"type:text"`
    Parameters  []TemplateParameter `json:"parameters" gorm:"foreignKey:TemplateID"`
    IsActive    bool      `json:"is_active" gorm:"default:true"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    CreatedBy   string    `json:"created_by"`
}

type TemplateParameter struct {
    ID          string `json:"id" gorm:"primaryKey"`
    TemplateID  string `json:"template_id"`
    Name        string `json:"name"`
    Type        string `json:"type"`
    Description string `json:"description"`
    Required    bool   `json:"required"`
    DefaultValue *string `json:"default_value,omitempty"`
}
```

#### Rule Version
```go
type RuleVersion struct {
    ID         string    `json:"id" gorm:"primaryKey"`
    RuleID     string    `json:"rule_id"`
    Version    int       `json:"version"`
    DSLContent string    `json:"dsl_content" gorm:"type:text"`
    ChangeLog  string    `json:"change_log"`
    CreatedAt  time.Time `json:"created_at"`
    CreatedBy  string    `json:"created_by"`
}
```

## REST API Endpoints

### Rule Management

#### Create Rule
```
POST /api/v1/rules
Content-Type: application/json

{
  "name": "Summer Promotion Rule",
  "description": "15% discount for summer products",
  "dsl_content": "IF product.category == 'summer' THEN discount = 0.15",
  "priority": "HIGH",
  "template_id": "promo-discount-template"
}

Response: 201 Created
{
  "id": "rule-123",
  "name": "Summer Promotion Rule",
  "status": "DRAFT",
  "version": 1,
  "created_at": "2024-01-15T10:00:00Z"
}
```

#### Get Rules (with filtering and pagination)
```
GET /api/v1/rules?status=ACTIVE&priority=HIGH&page=1&limit=20

Response: 200 OK
{
  "rules": [...],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 150,
    "total_pages": 8
  }
}
```

#### Update Rule
```
PUT /api/v1/rules/{rule_id}
Content-Type: application/json

{
  "name": "Updated Summer Promotion Rule",
  "dsl_content": "IF product.category == 'summer' THEN discount = 0.20"
}
```

#### Submit for Approval
```
POST /api/v1/rules/{rule_id}/submit-approval

Response: 200 OK
{
  "status": "UNDER_REVIEW",
  "workflow_id": "workflow-456"
}
```

### Template Management

#### Get Templates
```
GET /api/v1/templates?category=promotions

Response: 200 OK
{
  "templates": [
    {
      "id": "template-123",
      "name": "Percentage Discount Template",
      "category": "promotions",
      "parameters": [...]
    }
  ]
}
```

#### Create Rule from Template
```
POST /api/v1/templates/{template_id}/create-rule
Content-Type: application/json

{
  "name": "Black Friday Sale",
  "parameters": {
    "discount_percentage": "25",
    "product_category": "electronics"
  }
}
```

## gRPC Service Definition

```protobuf
syntax = "proto3";

package rules_management;

service RulesManagementService {
  rpc GetRule(GetRuleRequest) returns (GetRuleResponse);
  rpc CreateRule(CreateRuleRequest) returns (CreateRuleResponse);
  rpc UpdateRule(UpdateRuleRequest) returns (UpdateRuleResponse);
  rpc DeleteRule(DeleteRuleRequest) returns (DeleteRuleResponse);
  rpc ValidateRule(ValidateRuleRequest) returns (ValidateRuleResponse);
  rpc GetRulesByStatus(GetRulesByStatusRequest) returns (GetRulesByStatusResponse);
}

message Rule {
  string id = 1;
  string name = 2;
  string description = 3;
  string dsl_content = 4;
  string status = 5;
  string priority = 6;
  int32 version = 7;
  string created_by = 8;
  string created_at = 9;
}

message CreateRuleRequest {
  string name = 1;
  string description = 2;
  string dsl_content = 3;
  string priority = 4;
  string template_id = 5;
  string created_by = 6;
}

message ValidateRuleRequest {
  string dsl_content = 1;
}

message ValidateRuleResponse {
  bool is_valid = 1;
  repeated string errors = 2;
  repeated string warnings = 3;
}
```

## Domain Services

### Rule Service
```go
type RuleService interface {
    CreateRule(ctx context.Context, req CreateRuleRequest) (*Rule, error)
    UpdateRule(ctx context.Context, ruleID string, req UpdateRuleRequest) (*Rule, error)
    GetRule(ctx context.Context, ruleID string) (*Rule, error)
    DeleteRule(ctx context.Context, ruleID string) error
    ValidateRule(ctx context.Context, dslContent string) (*ValidationResult, error)
    SubmitForApproval(ctx context.Context, ruleID string, submitterID string) error
    ActivateRule(ctx context.Context, ruleID string) error
    DeactivateRule(ctx context.Context, ruleID string) error
}

type ValidationResult struct {
    IsValid  bool     `json:"is_valid"`
    Errors   []string `json:"errors,omitempty"`
    Warnings []string `json:"warnings,omitempty"`
}
```

### Template Service
```go
type TemplateService interface {
    GetTemplates(ctx context.Context, category string) ([]RuleTemplate, error)
    GetTemplate(ctx context.Context, templateID string) (*RuleTemplate, error)
    CreateRuleFromTemplate(ctx context.Context, templateID string, req CreateFromTemplateRequest) (*Rule, error)
    ValidateTemplateParameters(ctx context.Context, templateID string, parameters map[string]string) error
}
```

## Repository Interfaces

```go
type RuleRepository interface {
    Create(ctx context.Context, rule *Rule) error
    Update(ctx context.Context, rule *Rule) error
    GetByID(ctx context.Context, id string) (*Rule, error)
    GetByStatus(ctx context.Context, status RuleStatus) ([]Rule, error)
    Delete(ctx context.Context, id string) error
    Search(ctx context.Context, filters RuleFilters, pagination Pagination) ([]Rule, int64, error)
}

type TemplateRepository interface {
    GetAll(ctx context.Context, category string) ([]RuleTemplate, error)
    GetByID(ctx context.Context, id string) (*RuleTemplate, error)
    Create(ctx context.Context, template *RuleTemplate) error
    Update(ctx context.Context, template *RuleTemplate) error
}

type RuleVersionRepository interface {
    Create(ctx context.Context, version *RuleVersion) error
    GetByRuleID(ctx context.Context, ruleID string) ([]RuleVersion, error)
    GetVersion(ctx context.Context, ruleID string, version int) (*RuleVersion, error)
}
```

## Domain Events

```go
type RuleCreated struct {
    RuleID    string    `json:"rule_id"`
    Name      string    `json:"name"`
    CreatedBy string    `json:"created_by"`
    CreatedAt time.Time `json:"created_at"`
}

type RuleUpdated struct {
    RuleID      string `json:"rule_id"`
    Changes     map[string]interface{} `json:"changes"`
    UpdatedBy   string `json:"updated_by"`
    UpdatedAt   time.Time `json:"updated_at"`
    NewVersion  int    `json:"new_version"`
}

type RuleSubmittedForApproval struct {
    RuleID      string `json:"rule_id"`
    SubmittedBy string `json:"submitted_by"`
    SubmittedAt time.Time `json:"submitted_at"`
    WorkflowID  string `json:"workflow_id"`
}

type RuleApproved struct {
    RuleID     string `json:"rule_id"`
    ApprovedBy string `json:"approved_by"`
    ApprovedAt time.Time `json:"approved_at"`
}

type RuleActivated struct {
    RuleID      string `json:"rule_id"`
    ActivatedBy string `json:"activated_by"`
    ActivatedAt time.Time `json:"activated_at"`
}
```

## Implementation Tasks

### Phase 1: Core Infrastructure (2-3 days)
1. **Project Setup**
   - Initialize Go module with proper structure
   - Setup dependency injection container
   - Configure database connections (PostgreSQL)
   - Setup logging and configuration

2. **Database Layer**
   - Create database migrations for Rule, RuleTemplate, RuleVersion tables
   - Implement repository interfaces with GORM
   - Add database indexes for performance
   - Setup connection pooling and health checks

3. **Domain Layer Foundation**
   - Implement core domain entities (Rule, RuleTemplate, RuleVersion)
   - Create domain events structures
   - Implement basic domain services interfaces
   - Add domain validation logic

### Phase 2: REST API Implementation (3-4 days)
1. **HTTP Handlers**
   - Implement Gin HTTP handlers for all endpoints
   - Add request/response validation with struct tags
   - Implement proper error handling and status codes
   - Add OpenAPI/Swagger documentation

2. **Rule Management Endpoints**
   - POST /api/v1/rules (create rule)
   - GET /api/v1/rules (list rules with filtering)
   - GET /api/v1/rules/{id} (get rule details)
   - PUT /api/v1/rules/{id} (update rule)
   - DELETE /api/v1/rules/{id} (delete rule)
   - POST /api/v1/rules/{id}/submit-approval

3. **Template Management Endpoints**
   - GET /api/v1/templates (list templates)
   - GET /api/v1/templates/{id} (get template)
   - POST /api/v1/templates/{id}/create-rule (create from template)

### Phase 3: gRPC API Implementation (2-3 days)
1. **Protocol Buffer Definitions**
   - Define .proto files for all service methods
   - Generate Go code from protobuf definitions
   - Create gRPC service implementation
   - Add gRPC interceptors for logging and metrics

2. **Internal Service Methods**
   - Implement GetRule for internal service calls
   - Implement ValidateRule for rule validation
   - Implement GetRulesByStatus for bulk operations
   - Add streaming support for large result sets

### Phase 4: Business Logic & Validation (3-4 days)
1. **Rule Validation Engine**
   - Implement DSL syntax validation
   - Create rule semantic validation
   - Add business rule validation (conflicts, dependencies)
   - Implement rule testing with sample data

2. **Template Engine**
   - Implement template parameter substitution
   - Add template validation logic
   - Create rule generation from templates
   - Add template versioning support

3. **Version Management**
   - Implement rule versioning on updates
   - Add version comparison utilities
   - Create version rollback functionality
   - Add change tracking and audit logs

### Phase 5: Event Integration (2-3 days)
1. **NATS Integration**
   - Setup NATS client and connection management
   - Implement event publishing for domain events
   - Add event subscription handlers
   - Implement retry and dead letter queue handling

2. **Event Handlers**
   - Publish RuleCreated events
   - Publish RuleUpdated events with change tracking
   - Publish RuleSubmittedForApproval events
   - Handle ApprovalCompleted events from workflow service

### Phase 6: Testing & Quality (3-4 days)
1. **Unit Tests**
   - Write unit tests for all domain services (80% coverage)
   - Test repository implementations with test database
   - Test HTTP handlers with test requests
   - Test gRPC services with test clients

2. **Integration Tests**
   - Test end-to-end rule creation flow
   - Test rule approval workflow integration
   - Test template-based rule creation
   - Test event publishing and handling

3. **Performance Tests**
   - Load test rule creation endpoints
   - Performance test rule search and filtering
   - Test concurrent rule updates
   - Validate database query performance

### Phase 7: Production Readiness (2-3 days)
1. **Monitoring & Observability**
   - Add Prometheus metrics for all endpoints
   - Implement structured logging with correlation IDs
   - Add health check endpoints
   - Create service dashboards

2. **Security & Validation**
   - Implement input sanitization and validation
   - Add rate limiting and throttling
   - Implement proper error handling without data leakage
   - Add API authentication integration

3. **Documentation & Deployment**
   - Complete API documentation with examples
   - Create deployment configurations (Docker, K8s)
   - Add environment-specific configurations
   - Create runbooks and operational procedures

## Estimated Development Time: 17-24 days
