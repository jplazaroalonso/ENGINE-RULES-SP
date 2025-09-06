---
description: API Contract Implementation Rule
globs:
alwaysApply: false
---
# API Contract Implementation Rule

## Overview

This rule provides comprehensive guidance for designing and implementing REST APIs with OpenAPI 3.0 specifications and gRPC services with Protobuf definitions. Based on successful Rules Engine API implementations achieving production-ready contract-first development.

## 1. Feedback Section

### Identified Issues in Generic API Development:
- **Inconsistent API Design**: Mixing REST conventions, inconsistent response formats
- **Missing Documentation**: APIs without proper OpenAPI specifications or examples
- **Poor Error Handling**: Generic error responses without proper HTTP status codes
- **No Versioning Strategy**: Breaking changes without backward compatibility
- **Weak Validation**: Missing input validation and business rule enforcement
- **Unclear Contracts**: gRPC services without proper Protobuf documentation

### Recommendations:
- Design APIs following REST conventions with consistent resource modeling
- Create comprehensive OpenAPI 3.0 specifications with examples and validation
- Implement structured error responses with proper HTTP status codes
- Use semantic versioning with backward compatibility guarantees
- Add comprehensive input validation with business rule enforcement
- Define clear gRPC contracts with well-documented Protobuf schemas

## 2. Role and Context Definition

### Target Role: API Developer / Integration Architect
### Background Context:
- **API Styles**: REST (external), gRPC (internal high-performance)
- **Documentation**: OpenAPI 3.0, Protobuf with comprehensive examples
- **Validation**: JSON Schema validation, business rule enforcement
- **Error Handling**: RFC 7807 Problem Details, structured error responses
- **Versioning**: Semantic versioning with deprecation policies

## 3. Objective and Goals

### Primary Objective:
Design and implement production-ready APIs with comprehensive contracts, proper error handling, validation, and documentation that enable reliable service integration.

### Success Criteria:
- **Contract Completeness**: 100% API coverage with OpenAPI 3.0 and Protobuf specs
- **Documentation Quality**: Complete examples, descriptions, and integration guides
- **Error Handling**: Structured error responses with proper HTTP status codes
- **Validation**: Comprehensive input validation with business rule enforcement
- **Performance**: <200ms response time for simple operations, <500ms for complex
- **Compatibility**: Backward compatibility with semantic versioning

## 4. Key Terms and Definitions

### Technical Terminology:
- **OpenAPI 3.0**: API specification standard for REST APIs with validation schemas
- **Protobuf**: Protocol buffer language for defining gRPC service contracts
- **JSON Schema**: Validation schema for JSON request/response payloads
- **RFC 7807**: Problem Details standard for HTTP error responses
- **HATEOAS**: Hypermedia as the Engine of Application State for REST APIs
- **Service Contract**: Formal specification defining API behavior and data structures
- **API Gateway**: Centralized entry point for API routing, authentication, and rate limiting

## 5. Task Decomposition (Chain-of-Thought)

### Step 1: API Design and Resource Modeling
- **Input**: Domain requirements, business entities, user workflows
- **Process**: Design REST resources and gRPC services following best practices
- **Output**: Complete API design with resource hierarchy and operations
- **Human Validation Point**: Review API design follows REST conventions and business requirements

### Step 2: OpenAPI Specification Creation
- **Input**: REST API design, data models, error scenarios
- **Process**: Create comprehensive OpenAPI 3.0 specifications with validation
- **Output**: Complete OpenAPI specs with examples and documentation
- **Human Validation Point**: Verify specifications are complete and accurate

### Step 3: Protobuf Contract Definition
- **Input**: gRPC service requirements, performance needs, data models
- **Process**: Define Protobuf schemas with proper service and message definitions
- **Output**: Complete .proto files with documentation and validation
- **Human Validation Point**: Confirm gRPC contracts meet performance and usability requirements

### Step 4: Validation and Error Handling Implementation
- **Input**: Business rules, validation requirements, error scenarios
- **Process**: Implement comprehensive validation and structured error responses
- **Output**: Validation middleware and error handling with proper status codes
- **Human Validation Point**: Validate error handling covers all scenarios properly

### Step 5: Implementation and Testing
- **Input**: API contracts, validation rules, integration requirements
- **Process**: Implement APIs following contracts with comprehensive testing
- **Output**: Working APIs with contract compliance and integration tests
- **Human Validation Point**: Confirm implementation matches contracts exactly

### Step 6: Documentation and Integration Guides
- **Input**: API specifications, usage patterns, integration examples
- **Process**: Create comprehensive documentation and integration guides
- **Output**: Complete API documentation with examples and troubleshooting
- **Human Validation Point**: Verify documentation enables successful integration

## 6. Context and Constraints

### Technical Context:
- **REST Framework**: Gin/Echo (Go), Express (Node.js), FastAPI (Python)
- **gRPC Framework**: grpc-go, grpc-node, grpcio (Python)
- **Validation**: JSON Schema, Protobuf validation, custom business rules
- **Documentation**: OpenAPI 3.0, Protobuf comments, Swagger UI
- **Testing**: Contract testing, integration testing, load testing

### Business Context:
- **Performance**: <200ms simple operations, <500ms complex operations
- **Availability**: 99.9% uptime with graceful degradation
- **Security**: Authentication, authorization, input validation
- **Scalability**: Support for 1000+ concurrent requests
- **Compatibility**: Backward compatibility with deprecation policies

### Negative Constraints:
- **Do NOT** break existing API contracts without proper versioning
- **Do NOT** return generic error messages without context
- **Do NOT** skip input validation and business rule enforcement
- **Do NOT** create APIs without comprehensive documentation
- **Do NOT** ignore performance requirements and optimization

## 7. Examples and Illustrations (Few-Shot)

### Example 1: OpenAPI 3.0 Specification for Rules API

```yaml
# api/openapi/rules-api.yaml
openapi: 3.0.3
info:
  title: Rules Management API
  description: |
    Comprehensive API for managing business rules, templates, and workflows.
    
    ## Authentication
    All endpoints require Bearer token authentication except for health checks.
    
    ## Rate Limiting
    - 1000 requests per minute for authenticated users
    - 100 requests per minute for unauthenticated health checks
    
    ## Error Handling
    All errors follow RFC 7807 Problem Details format with structured responses.
    
  version: 1.0.0
  contact:
    name: API Support
    url: https://api-docs.rules-engine.com
    email: api-support@rules-engine.com
  license:
    name: MIT
    url: https://opensource.org/licenses/MIT

servers:
  - url: https://api.rules-engine.com/v1
    description: Production server
  - url: https://staging-api.rules-engine.com/v1
    description: Staging server
  - url: http://localhost:8080/v1
    description: Development server

security:
  - BearerAuth: []

tags:
  - name: Rules
    description: Business rule management operations
    externalDocs:
      description: Find out more
      url: https://docs.rules-engine.com/rules
  - name: Templates
    description: Rule template operations
  - name: Evaluation
    description: Rule evaluation and testing
  - name: Health
    description: Health check endpoints

paths:
  /rules:
    get:
      tags: [Rules]
      summary: List rules
      description: |
        Retrieve a paginated list of rules with optional filtering and sorting.
        
        ## Filtering
        - Filter by status, priority, category, or creation date
        - Search by name or description using the `search` parameter
        - Combine multiple filters for precise results
        
        ## Sorting
        - Sort by any field using `sortBy` parameter
        - Use `sortOrder` to specify ascending or descending order
        
      operationId: listRules
      parameters:
        - $ref: '#/components/parameters/PageParam'
        - $ref: '#/components/parameters/LimitParam'
        - $ref: '#/components/parameters/SortByParam'
        - $ref: '#/components/parameters/SortOrderParam'
        - name: status
          in: query
          description: Filter by rule status
          schema:
            type: array
            items:
              $ref: '#/components/schemas/RuleStatus'
          example: ["ACTIVE", "DRAFT"]
        - name: priority
          in: query
          description: Filter by rule priority
          schema:
            type: array
            items:
              $ref: '#/components/schemas/RulePriority'
          example: ["HIGH", "CRITICAL"]
        - name: category
          in: query
          description: Filter by rule category
          schema:
            type: string
          example: "customer-loyalty"
        - name: createdBy
          in: query
          description: Filter by rule creator
          schema:
            type: string
          example: "john.doe"
        - name: search
          in: query
          description: Search in rule name and description
          schema:
            type: string
          example: "discount"
        - name: dateFrom
          in: query
          description: Filter rules created after this date
          schema:
            type: string
            format: date-time
          example: "2024-01-01T00:00:00Z"
        - name: dateTo
          in: query
          description: Filter rules created before this date
          schema:
            type: string
            format: date-time
          example: "2024-12-31T23:59:59Z"
      responses:
        '200':
          description: Rules retrieved successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/RuleListResponse'
              examples:
                success:
                  summary: Successful response with rules
                  value:
                    data:
                      - id: "rule-123"
                        name: "VIP Customer Discount"
                        description: "Apply 20% discount for VIP customers"
                        status: "ACTIVE"
                        priority: "HIGH"
                        category: "customer-loyalty"
                        dslContent: "customer.tier == 'VIP'"
                        version: 1
                        createdAt: "2024-01-15T10:30:00Z"
                        updatedAt: "2024-01-15T10:30:00Z"
                        createdBy: "john.doe"
                        tags: ["discount", "vip", "loyalty"]
                    pagination:
                      page: 1
                      limit: 20
                      total: 45
                      totalPages: 3
                    links:
                      self: "/v1/rules?page=1&limit=20"
                      next: "/v1/rules?page=2&limit=20"
                      last: "/v1/rules?page=3&limit=20"
        '400':
          $ref: '#/components/responses/BadRequest'
        '401':
          $ref: '#/components/responses/Unauthorized'
        '403':
          $ref: '#/components/responses/Forbidden'
        '429':
          $ref: '#/components/responses/TooManyRequests'
        '500':
          $ref: '#/components/responses/InternalServerError'
    
    post:
      tags: [Rules]
      summary: Create a new rule
      description: |
        Create a new business rule with validation and optional template usage.
        
        ## Validation
        - Rule name must be unique within the organization
        - DSL content must pass syntax validation
        - Priority and status must be valid enum values
        
        ## Business Rules
        - Only users with 'rule_creator' role can create rules
        - Rules start in 'DRAFT' status by default
        - DSL content is validated against the rule engine grammar
        
      operationId: createRule
      requestBody:
        description: Rule creation data
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/CreateRuleRequest'
            examples:
              basic_rule:
                summary: Basic rule creation
                value:
                  name: "New Customer Welcome Discount"
                  description: "Apply 10% discount for new customers"
                  dslContent: "customer.isNew == true"
                  priority: "MEDIUM"
                  category: "customer-acquisition"
                  tags: ["welcome", "discount", "new-customer"]
              template_based:
                summary: Rule from template
                value:
                  name: "Seasonal Promotion Rule"
                  description: "Summer sale discount rule"
                  templateId: "template-456"
                  parameters:
                    discountPercentage: 25
                    startDate: "2024-06-01"
                    endDate: "2024-08-31"
                  priority: "HIGH"
                  category: "promotions"
      responses:
        '201':
          description: Rule created successfully
          headers:
            Location:
              description: URL of the created rule
              schema:
                type: string
              example: "/v1/rules/rule-789"
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/RuleResponse'
              examples:
                success:
                  summary: Successfully created rule
                  value:
                    data:
                      id: "rule-789"
                      name: "New Customer Welcome Discount"
                      description: "Apply 10% discount for new customers"
                      status: "DRAFT"
                      priority: "MEDIUM"
                      category: "customer-acquisition"
                      dslContent: "customer.isNew == true"
                      version: 1
                      createdAt: "2024-01-20T14:30:00Z"
                      updatedAt: "2024-01-20T14:30:00Z"
                      createdBy: "jane.smith"
                      tags: ["welcome", "discount", "new-customer"]
                    links:
                      self: "/v1/rules/rule-789"
                      edit: "/v1/rules/rule-789"
                      validate: "/v1/rules/rule-789/validate"
        '400':
          $ref: '#/components/responses/BadRequest'
        '401':
          $ref: '#/components/responses/Unauthorized'
        '403':
          $ref: '#/components/responses/Forbidden'
        '409':
          description: Rule name already exists
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ProblemDetails'
              examples:
                duplicate_name:
                  summary: Duplicate rule name
                  value:
                    type: "https://api.rules-engine.com/problems/duplicate-rule-name"
                    title: "Rule Name Already Exists"
                    status: 409
                    detail: "A rule with the name 'New Customer Welcome Discount' already exists"
                    instance: "/v1/rules"
                    errors:
                      name: ["Rule name must be unique"]
        '422':
          description: Validation failed
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ValidationProblemDetails'
              examples:
                validation_errors:
                  summary: Multiple validation errors
                  value:
                    type: "https://api.rules-engine.com/problems/validation-error"
                    title: "Validation Failed"
                    status: 422
                    detail: "The request contains invalid data"
                    instance: "/v1/rules"
                    errors:
                      name: ["Name is required", "Name must be at least 3 characters"]
                      dslContent: ["DSL syntax error at line 1: unexpected token 'invalid'"]
                      priority: ["Priority must be one of: LOW, MEDIUM, HIGH, CRITICAL"]
        '429':
          $ref: '#/components/responses/TooManyRequests'
        '500':
          $ref: '#/components/responses/InternalServerError'

  /rules/{ruleId}:
    get:
      tags: [Rules]
      summary: Get rule by ID
      description: |
        Retrieve a specific rule by its unique identifier.
        
        ## Response Details
        - Includes complete rule information
        - Contains metadata about creation and modifications
        - Provides links to related operations
        
      operationId: getRule
      parameters:
        - $ref: '#/components/parameters/RuleIdParam'
      responses:
        '200':
          description: Rule retrieved successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/RuleResponse'
        '400':
          $ref: '#/components/responses/BadRequest'
        '401':
          $ref: '#/components/responses/Unauthorized'
        '403':
          $ref: '#/components/responses/Forbidden'
        '404':
          $ref: '#/components/responses/NotFound'
        '500':
          $ref: '#/components/responses/InternalServerError'
    
    put:
      tags: [Rules]
      summary: Update rule
      description: |
        Update an existing rule with validation and version control.
        
        ## Business Rules
        - Only draft rules can be fully updated
        - Active rules require approval workflow for changes
        - Version is automatically incremented on successful update
        
      operationId: updateRule
      parameters:
        - $ref: '#/components/parameters/RuleIdParam'
      requestBody:
        description: Rule update data
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/UpdateRuleRequest'
      responses:
        '200':
          description: Rule updated successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/RuleResponse'
        '400':
          $ref: '#/components/responses/BadRequest'
        '401':
          $ref: '#/components/responses/Unauthorized'
        '403':
          $ref: '#/components/responses/Forbidden'
        '404':
          $ref: '#/components/responses/NotFound'
        '409':
          description: Rule cannot be updated in current state
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ProblemDetails'
        '422':
          $ref: '#/components/responses/ValidationError'
        '500':
          $ref: '#/components/responses/InternalServerError'
    
    delete:
      tags: [Rules]
      summary: Delete rule
      description: |
        Delete a rule (soft delete with audit trail).
        
        ## Business Rules
        - Only draft rules can be deleted
        - Active rules must be deactivated first
        - Deletion creates audit log entry
        
      operationId: deleteRule
      parameters:
        - $ref: '#/components/parameters/RuleIdParam'
      responses:
        '204':
          description: Rule deleted successfully
        '400':
          $ref: '#/components/responses/BadRequest'
        '401':
          $ref: '#/components/responses/Unauthorized'
        '403':
          $ref: '#/components/responses/Forbidden'
        '404':
          $ref: '#/components/responses/NotFound'
        '409':
          description: Rule cannot be deleted in current state
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ProblemDetails'
        '500':
          $ref: '#/components/responses/InternalServerError'

  /rules/{ruleId}/validate:
    post:
      tags: [Rules]
      summary: Validate rule DSL
      description: |
        Validate rule DSL syntax and optionally test with sample data.
        
        ## Validation Types
        - Syntax validation: Check DSL grammar and structure
        - Semantic validation: Verify variable references and types
        - Test execution: Run rule with provided test data
        
      operationId: validateRule
      parameters:
        - $ref: '#/components/parameters/RuleIdParam'
      requestBody:
        description: Validation request with optional test data
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/ValidationRequest'
      responses:
        '200':
          description: Validation completed
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ValidationResponse'
        '400':
          $ref: '#/components/responses/BadRequest'
        '401':
          $ref: '#/components/responses/Unauthorized'
        '403':
          $ref: '#/components/responses/Forbidden'
        '404':
          $ref: '#/components/responses/NotFound'
        '500':
          $ref: '#/components/responses/InternalServerError'

  /health:
    get:
      tags: [Health]
      summary: Health check
      description: Check API health and dependencies
      operationId: healthCheck
      security: []  # No authentication required
      responses:
        '200':
          description: Service is healthy
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/HealthResponse'

components:
  securitySchemes:
    BearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT
      description: |
        JWT token authentication. Include the token in the Authorization header:
        `Authorization: Bearer <token>`

  parameters:
    RuleIdParam:
      name: ruleId
      in: path
      description: Unique rule identifier
      required: true
      schema:
        type: string
        format: uuid
        pattern: '^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$'
      example: "123e4567-e89b-12d3-a456-426614174000"
    
    PageParam:
      name: page
      in: query
      description: Page number (1-based)
      schema:
        type: integer
        minimum: 1
        default: 1
      example: 1
    
    LimitParam:
      name: limit
      in: query
      description: Number of items per page
      schema:
        type: integer
        minimum: 1
        maximum: 100
        default: 20
      example: 20
    
    SortByParam:
      name: sortBy
      in: query
      description: Field to sort by
      schema:
        type: string
        enum: [name, status, priority, createdAt, updatedAt]
        default: createdAt
      example: "name"
    
    SortOrderParam:
      name: sortOrder
      in: query
      description: Sort order
      schema:
        type: string
        enum: [asc, desc]
        default: desc
      example: "asc"

  schemas:
    Rule:
      type: object
      description: Business rule entity
      required:
        - id
        - name
        - status
        - priority
        - dslContent
        - version
        - createdAt
        - updatedAt
        - createdBy
      properties:
        id:
          type: string
          format: uuid
          description: Unique rule identifier
          example: "123e4567-e89b-12d3-a456-426614174000"
        name:
          type: string
          minLength: 3
          maxLength: 100
          description: Human-readable rule name
          example: "VIP Customer Discount"
        description:
          type: string
          maxLength: 500
          description: Detailed rule description
          example: "Apply 20% discount for VIP customers on all purchases"
        status:
          $ref: '#/components/schemas/RuleStatus'
        priority:
          $ref: '#/components/schemas/RulePriority'
        category:
          type: string
          maxLength: 50
          description: Rule category for organization
          example: "customer-loyalty"
        dslContent:
          type: string
          description: Domain-specific language rule definition
          example: "customer.tier == 'VIP' && order.amount > 100"
        version:
          type: integer
          minimum: 1
          description: Rule version number
          example: 1
        templateId:
          type: string
          format: uuid
          description: Template used to create this rule
          example: "456e7890-e89b-12d3-a456-426614174111"
        tags:
          type: array
          items:
            type: string
            maxLength: 30
          maxItems: 10
          description: Tags for categorization and search
          example: ["discount", "vip", "loyalty"]
        createdAt:
          type: string
          format: date-time
          description: Rule creation timestamp
          example: "2024-01-15T10:30:00Z"
        updatedAt:
          type: string
          format: date-time
          description: Last modification timestamp
          example: "2024-01-15T10:30:00Z"
        createdBy:
          type: string
          description: User who created the rule
          example: "john.doe"
        approvedBy:
          type: string
          description: User who approved the rule
          example: "jane.manager"
        approvedAt:
          type: string
          format: date-time
          description: Rule approval timestamp
          example: "2024-01-16T09:15:00Z"

    RuleStatus:
      type: string
      enum: [DRAFT, UNDER_REVIEW, APPROVED, ACTIVE, INACTIVE, DEPRECATED]
      description: Rule lifecycle status
      example: "ACTIVE"

    RulePriority:
      type: string
      enum: [LOW, MEDIUM, HIGH, CRITICAL]
      description: Rule execution priority
      example: "HIGH"

    CreateRuleRequest:
      type: object
      description: Request to create a new rule
      required:
        - name
        - dslContent
        - priority
      properties:
        name:
          type: string
          minLength: 3
          maxLength: 100
          description: Rule name (must be unique)
          example: "New Customer Welcome Discount"
        description:
          type: string
          maxLength: 500
          description: Rule description
          example: "Apply 10% discount for new customers"
        dslContent:
          type: string
          description: DSL rule definition
          example: "customer.isNew == true"
        priority:
          $ref: '#/components/schemas/RulePriority'
        category:
          type: string
          maxLength: 50
          description: Rule category
          example: "customer-acquisition"
        templateId:
          type: string
          format: uuid
          description: Template to use for rule creation
          example: "456e7890-e89b-12d3-a456-426614174111"
        parameters:
          type: object
          description: Template parameters (if using template)
          additionalProperties: true
          example:
            discountPercentage: 10
            validDays: 30
        tags:
          type: array
          items:
            type: string
            maxLength: 30
          maxItems: 10
          description: Rule tags
          example: ["welcome", "discount", "new-customer"]

    UpdateRuleRequest:
      type: object
      description: Request to update an existing rule
      properties:
        name:
          type: string
          minLength: 3
          maxLength: 100
          description: Updated rule name
          example: "Updated Customer Welcome Discount"
        description:
          type: string
          maxLength: 500
          description: Updated rule description
          example: "Apply 15% discount for new customers"
        dslContent:
          type: string
          description: Updated DSL rule definition
          example: "customer.isNew == true && order.amount > 50"
        priority:
          $ref: '#/components/schemas/RulePriority'
        category:
          type: string
          maxLength: 50
          description: Updated rule category
          example: "customer-acquisition"
        tags:
          type: array
          items:
            type: string
            maxLength: 30
          maxItems: 10
          description: Updated rule tags
          example: ["welcome", "discount", "new-customer", "updated"]

    RuleResponse:
      type: object
      description: Rule response with data and links
      required:
        - data
        - links
      properties:
        data:
          $ref: '#/components/schemas/Rule'
        links:
          type: object
          description: HATEOAS links for related operations
          properties:
            self:
              type: string
              format: uri
              example: "/v1/rules/123e4567-e89b-12d3-a456-426614174000"
            edit:
              type: string
              format: uri
              example: "/v1/rules/123e4567-e89b-12d3-a456-426614174000"
            validate:
              type: string
              format: uri
              example: "/v1/rules/123e4567-e89b-12d3-a456-426614174000/validate"
            approve:
              type: string
              format: uri
              example: "/v1/rules/123e4567-e89b-12d3-a456-426614174000/approve"

    RuleListResponse:
      type: object
      description: Paginated list of rules
      required:
        - data
        - pagination
        - links
      properties:
        data:
          type: array
          items:
            $ref: '#/components/schemas/Rule'
        pagination:
          $ref: '#/components/schemas/PaginationInfo'
        links:
          $ref: '#/components/schemas/PaginationLinks'

    PaginationInfo:
      type: object
      description: Pagination metadata
      required:
        - page
        - limit
        - total
        - totalPages
      properties:
        page:
          type: integer
          minimum: 1
          description: Current page number
          example: 1
        limit:
          type: integer
          minimum: 1
          description: Items per page
          example: 20
        total:
          type: integer
          minimum: 0
          description: Total number of items
          example: 45
        totalPages:
          type: integer
          minimum: 0
          description: Total number of pages
          example: 3

    PaginationLinks:
      type: object
      description: HATEOAS pagination links
      required:
        - self
      properties:
        self:
          type: string
          format: uri
          description: Current page link
          example: "/v1/rules?page=1&limit=20"
        first:
          type: string
          format: uri
          description: First page link
          example: "/v1/rules?page=1&limit=20"
        prev:
          type: string
          format: uri
          description: Previous page link
          example: "/v1/rules?page=1&limit=20"
        next:
          type: string
          format: uri
          description: Next page link
          example: "/v1/rules?page=2&limit=20"
        last:
          type: string
          format: uri
          description: Last page link
          example: "/v1/rules?page=3&limit=20"

    ValidationRequest:
      type: object
      description: Rule validation request
      properties:
        testData:
          type: object
          description: Optional test data for rule execution
          additionalProperties: true
          example:
            customer:
              id: "cust-123"
              tier: "VIP"
              isNew: false
            order:
              id: "order-456"
              amount: 150.00
              currency: "USD"

    ValidationResponse:
      type: object
      description: Rule validation result
      required:
        - valid
        - errors
        - warnings
      properties:
        valid:
          type: boolean
          description: Whether the rule is valid
          example: true
        errors:
          type: array
          items:
            $ref: '#/components/schemas/ValidationError'
          description: Validation errors
        warnings:
          type: array
          items:
            $ref: '#/components/schemas/ValidationWarning'
          description: Validation warnings
        testResult:
          type: object
          description: Test execution result (if test data provided)
          properties:
            executed:
              type: boolean
              description: Whether rule was executed
              example: true
            result:
              type: boolean
              description: Rule evaluation result
              example: true
            executionTime:
              type: number
              description: Execution time in milliseconds
              example: 12.5
            details:
              type: object
              description: Detailed execution information
              additionalProperties: true

    ValidationError:
      type: object
      description: Validation error details
      required:
        - code
        - message
      properties:
        code:
          type: string
          description: Error code
          example: "SYNTAX_ERROR"
        message:
          type: string
          description: Human-readable error message
          example: "Unexpected token 'invalid' at line 1, column 15"
        line:
          type: integer
          description: Line number (if applicable)
          example: 1
        column:
          type: integer
          description: Column number (if applicable)
          example: 15
        context:
          type: string
          description: Error context
          example: "customer.tier == 'invalid'"

    ValidationWarning:
      type: object
      description: Validation warning details
      required:
        - code
        - message
      properties:
        code:
          type: string
          description: Warning code
          example: "UNUSED_VARIABLE"
        message:
          type: string
          description: Human-readable warning message
          example: "Variable 'order.discount' is referenced but may not be available"
        line:
          type: integer
          description: Line number (if applicable)
          example: 2

    HealthResponse:
      type: object
      description: Health check response
      required:
        - status
        - timestamp
        - version
        - dependencies
      properties:
        status:
          type: string
          enum: [UP, DOWN, DEGRADED]
          description: Overall service status
          example: "UP"
        timestamp:
          type: string
          format: date-time
          description: Health check timestamp
          example: "2024-01-20T15:30:00Z"
        version:
          type: string
          description: API version
          example: "1.0.0"
        uptime:
          type: integer
          description: Service uptime in seconds
          example: 86400
        dependencies:
          type: object
          description: Dependency health status
          properties:
            database:
              $ref: '#/components/schemas/DependencyHealth'
            cache:
              $ref: '#/components/schemas/DependencyHealth'
            messaging:
              $ref: '#/components/schemas/DependencyHealth'

    DependencyHealth:
      type: object
      description: Individual dependency health
      required:
        - status
        - responseTime
      properties:
        status:
          type: string
          enum: [UP, DOWN, DEGRADED]
          description: Dependency status
          example: "UP"
        responseTime:
          type: number
          description: Response time in milliseconds
          example: 15.2
        lastCheck:
          type: string
          format: date-time
          description: Last health check timestamp
          example: "2024-01-20T15:30:00Z"
        error:
          type: string
          description: Error message (if status is DOWN)
          example: "Connection timeout"

    ProblemDetails:
      type: object
      description: RFC 7807 Problem Details for HTTP APIs
      required:
        - type
        - title
        - status
      properties:
        type:
          type: string
          format: uri
          description: Problem type identifier
          example: "https://api.rules-engine.com/problems/duplicate-rule-name"
        title:
          type: string
          description: Short, human-readable summary
          example: "Rule Name Already Exists"
        status:
          type: integer
          description: HTTP status code
          example: 409
        detail:
          type: string
          description: Human-readable explanation
          example: "A rule with the name 'VIP Customer Discount' already exists"
        instance:
          type: string
          format: uri
          description: URI reference to specific occurrence
          example: "/v1/rules"
        timestamp:
          type: string
          format: date-time
          description: Error timestamp
          example: "2024-01-20T15:30:00Z"
        traceId:
          type: string
          description: Request trace identifier
          example: "123e4567-e89b-12d3-a456-426614174000"

    ValidationProblemDetails:
      allOf:
        - $ref: '#/components/schemas/ProblemDetails'
        - type: object
          properties:
            errors:
              type: object
              description: Field-specific validation errors
              additionalProperties:
                type: array
                items:
                  type: string
              example:
                name: ["Name is required", "Name must be unique"]
                dslContent: ["DSL syntax error: unexpected token"]

  responses:
    BadRequest:
      description: Bad request - invalid input
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/ProblemDetails'
          examples:
            invalid_parameter:
              summary: Invalid parameter value
              value:
                type: "https://api.rules-engine.com/problems/invalid-parameter"
                title: "Invalid Parameter"
                status: 400
                detail: "The parameter 'limit' must be between 1 and 100"
                instance: "/v1/rules"

    Unauthorized:
      description: Unauthorized - authentication required
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/ProblemDetails'
          examples:
            missing_token:
              summary: Missing authentication token
              value:
                type: "https://api.rules-engine.com/problems/unauthorized"
                title: "Authentication Required"
                status: 401
                detail: "Access token is required"
                instance: "/v1/rules"

    Forbidden:
      description: Forbidden - insufficient permissions
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/ProblemDetails'
          examples:
            insufficient_permissions:
              summary: Insufficient permissions
              value:
                type: "https://api.rules-engine.com/problems/forbidden"
                title: "Insufficient Permissions"
                status: 403
                detail: "User does not have permission to perform this action"
                instance: "/v1/rules"

    NotFound:
      description: Resource not found
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/ProblemDetails'
          examples:
            rule_not_found:
              summary: Rule not found
              value:
                type: "https://api.rules-engine.com/problems/resource-not-found"
                title: "Rule Not Found"
                status: 404
                detail: "Rule with ID '123e4567-e89b-12d3-a456-426614174000' was not found"
                instance: "/v1/rules/123e4567-e89b-12d3-a456-426614174000"

    ValidationError:
      description: Validation error - request data is invalid
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/ValidationProblemDetails'

    TooManyRequests:
      description: Too many requests - rate limit exceeded
      headers:
        Retry-After:
          description: Seconds to wait before retrying
          schema:
            type: integer
          example: 60
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/ProblemDetails'
          examples:
            rate_limit_exceeded:
              summary: Rate limit exceeded
              value:
                type: "https://api.rules-engine.com/problems/rate-limit-exceeded"
                title: "Rate Limit Exceeded"
                status: 429
                detail: "Rate limit of 1000 requests per minute exceeded"
                instance: "/v1/rules"

    InternalServerError:
      description: Internal server error
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/ProblemDetails'
          examples:
            internal_error:
              summary: Internal server error
              value:
                type: "https://api.rules-engine.com/problems/internal-error"
                title: "Internal Server Error"
                status: 500
                detail: "An unexpected error occurred while processing the request"
                instance: "/v1/rules"
```

### Example 2: gRPC Protobuf Service Definition

```protobuf
// api/proto/rules/v1/rules_service.proto
syntax = "proto3";

package rules.v1;

import "google/api/annotations.proto";
import "google/protobuf/timestamp.proto";
import "google/protobuf/empty.proto";
import "validate/validate.proto";
import "rules/v1/rule.proto";
import "common/v1/pagination.proto";

option go_package = "github.com/rules-engine/api/gen/go/rules/v1;rulesv1";
option java_package = "com.rulesengine.api.rules.v1";
option java_multiple_files = true;

// RulesService provides high-performance rule management operations
// for internal service-to-service communication.
//
// This service is optimized for:
// - High throughput rule evaluation (>10,000 TPS)
// - Low latency operations (<50ms P95)
// - Reliable rule lifecycle management
// - Efficient batch operations
service RulesService {
  // GetRule retrieves a single rule by ID with caching optimization
  rpc GetRule(GetRuleRequest) returns (GetRuleResponse) {
    option (google.api.http) = {
      get: "/v1/rules/{rule_id}"
    };
  }

  // ListRules retrieves paginated rules with advanced filtering
  rpc ListRules(ListRulesRequest) returns (ListRulesResponse) {
    option (google.api.http) = {
      get: "/v1/rules"
    };
  }

  // CreateRule creates a new rule with validation
  rpc CreateRule(CreateRuleRequest) returns (CreateRuleResponse) {
    option (google.api.http) = {
      post: "/v1/rules"
      body: "*"
    };
  }

  // UpdateRule updates an existing rule with version control
  rpc UpdateRule(UpdateRuleRequest) returns (UpdateRuleResponse) {
    option (google.api.http) = {
      put: "/v1/rules/{rule_id}"
      body: "*"
    };
  }

  // DeleteRule soft-deletes a rule with audit trail
  rpc DeleteRule(DeleteRuleRequest) returns (google.protobuf.Empty) {
    option (google.api.http) = {
      delete: "/v1/rules/{rule_id}"
    };
  }

  // EvaluateRule executes rule evaluation with performance metrics
  rpc EvaluateRule(EvaluateRuleRequest) returns (EvaluateRuleResponse) {
    option (google.api.http) = {
      post: "/v1/rules/{rule_id}/evaluate"
      body: "*"
    };
  }

  // BatchEvaluateRules evaluates multiple rules in a single request
  // for improved performance and reduced network overhead
  rpc BatchEvaluateRules(BatchEvaluateRulesRequest) returns (BatchEvaluateRulesResponse) {
    option (google.api.http) = {
      post: "/v1/rules/batch-evaluate"
      body: "*"
    };
  }

  // ValidateRule performs DSL syntax and semantic validation
  rpc ValidateRule(ValidateRuleRequest) returns (ValidateRuleResponse) {
    option (google.api.http) = {
      post: "/v1/rules/{rule_id}/validate"
      body: "*"
    };
  }

  // SubscribeRuleChanges provides real-time rule change notifications
  // using server-side streaming for event-driven architectures
  rpc SubscribeRuleChanges(SubscribeRuleChangesRequest) returns (stream RuleChangeEvent);

  // GetRuleMetrics retrieves performance and usage metrics for a rule
  rpc GetRuleMetrics(GetRuleMetricsRequest) returns (GetRuleMetricsResponse) {
    option (google.api.http) = {
      get: "/v1/rules/{rule_id}/metrics"
    };
  }
}

// GetRuleRequest specifies the rule ID to retrieve
message GetRuleRequest {
  // Unique rule identifier
  string rule_id = 1 [(validate.rules).string.uuid = true];

  // Include related entities in response
  bool include_template = 2;
  bool include_metrics = 3;
}

// GetRuleResponse contains the requested rule data
message GetRuleResponse {
  // The requested rule
  Rule rule = 1;

  // Rule template if requested and available
  RuleTemplate template = 2;

  // Rule metrics if requested
  RuleMetrics metrics = 3;
}

// ListRulesRequest defines filtering and pagination options
message ListRulesRequest {
  // Pagination parameters
  common.v1.PaginationRequest pagination = 1;

  // Filtering options
  RuleFilter filter = 2;

  // Sorting options
  RuleSorting sorting = 3;

  // Fields to include in response for performance optimization
  repeated string include_fields = 4;
}

// RuleFilter defines available filtering criteria
message RuleFilter {
  // Filter by rule status
  repeated RuleStatus status = 1;

  // Filter by rule priority
  repeated RulePriority priority = 2;

  // Filter by category
  repeated string category = 3;

  // Filter by creator
  repeated string created_by = 4;

  // Filter by creation date range
  google.protobuf.Timestamp created_after = 5;
  google.protobuf.Timestamp created_before = 6;

  // Filter by tags
  repeated string tags = 7;

  // Search in name and description
  string search_query = 8 [(validate.rules).string.max_len = 100];

  // Filter by template ID
  repeated string template_id = 9;
}

// RuleSorting defines sorting options
message RuleSorting {
  // Field to sort by
  enum SortField {
    SORT_FIELD_UNSPECIFIED = 0;
    SORT_FIELD_NAME = 1;
    SORT_FIELD_STATUS = 2;
    SORT_FIELD_PRIORITY = 3;
    SORT_FIELD_CREATED_AT = 4;
    SORT_FIELD_UPDATED_AT = 5;
  }

  // Sort direction
  enum SortDirection {
    SORT_DIRECTION_UNSPECIFIED = 0;
    SORT_DIRECTION_ASC = 1;
    SORT_DIRECTION_DESC = 2;
  }

  SortField field = 1;
  SortDirection direction = 2;
}

// ListRulesResponse contains paginated rule results
message ListRulesResponse {
  // List of rules matching the criteria
  repeated Rule rules = 1;

  // Pagination information
  common.v1.PaginationResponse pagination = 2;
}

// CreateRuleRequest defines data for creating a new rule
message CreateRuleRequest {
  // Rule name (must be unique within organization)
  string name = 1 [
    (validate.rules).string.min_len = 3,
    (validate.rules).string.max_len = 100
  ];

  // Rule description
  string description = 2 [(validate.rules).string.max_len = 500];

  // DSL content defining the rule logic
  string dsl_content = 3 [(validate.rules).string.min_len = 1];

  // Rule priority for execution ordering
  RulePriority priority = 4 [(validate.rules).enum.defined_only = true];

  // Rule category for organization
  string category = 5 [(validate.rules).string.max_len = 50];

  // Template ID if creating from template
  string template_id = 6 [(validate.rules).string.uuid = true];

  // Template parameters if using template
  map<string, string> template_parameters = 7;

  // Tags for categorization and search
  repeated string tags = 8 [
    (validate.rules).repeated.max_items = 10,
    (validate.rules).repeated.items.string.max_len = 30
  ];

  // User creating the rule
  string created_by = 9 [
    (validate.rules).string.min_len = 1,
    (validate.rules).string.max_len = 100
  ];
}

// CreateRuleResponse contains the created rule data
message CreateRuleResponse {
  // The newly created rule
  Rule rule = 1;

  // Validation results if any warnings were generated
  ValidationResult validation_result = 2;
}

// UpdateRuleRequest defines data for updating an existing rule
message UpdateRuleRequest {
  // Rule ID to update
  string rule_id = 1 [(validate.rules).string.uuid = true];

  // Updated rule name
  string name = 2 [(validate.rules).string.max_len = 100];

  // Updated rule description
  string description = 3 [(validate.rules).string.max_len = 500];

  // Updated DSL content
  string dsl_content = 4;

  // Updated rule priority
  RulePriority priority = 5;

  // Updated rule category
  string category = 6 [(validate.rules).string.max_len = 50];

  // Updated tags
  repeated string tags = 7 [
    (validate.rules).repeated.max_items = 10,
    (validate.rules).repeated.items.string.max_len = 30
  ];

  // User performing the update
  string updated_by = 8 [
    (validate.rules).string.min_len = 1,
    (validate.rules).string.max_len = 100
  ];
}

// UpdateRuleResponse contains the updated rule data
message UpdateRuleResponse {
  // The updated rule
  Rule rule = 1;

  // Validation results if any warnings were generated
  ValidationResult validation_result = 2;
}

// DeleteRuleRequest specifies the rule to delete
message DeleteRuleRequest {
  // Rule ID to delete
  string rule_id = 1 [(validate.rules).string.uuid = true];

  // User performing the deletion
  string deleted_by = 2 [
    (validate.rules).string.min_len = 1,
    (validate.rules).string.max_len = 100
  ];

  // Reason for deletion (optional)
  string reason = 3 [(validate.rules).string.max_len = 200];
}

// EvaluateRuleRequest defines data for rule evaluation
message EvaluateRuleRequest {
  // Rule ID to evaluate
  string rule_id = 1 [(validate.rules).string.uuid = true];

  // Context data for rule evaluation
  map<string, string> context = 2;

  // Include performance metrics in response
  bool include_metrics = 3;

  // Request ID for tracing and debugging
  string request_id = 4 [(validate.rules).string.uuid = true];
}

// EvaluateRuleResponse contains evaluation results
message EvaluateRuleResponse {
  // Evaluation result (true/false)
  bool result = 1;

  // Execution time in milliseconds
  double execution_time_ms = 2;

  // Detailed execution information
  ExecutionDetails execution_details = 3;

  // Performance metrics if requested
  EvaluationMetrics metrics = 4;

  // Any warnings or errors during evaluation
  repeated ExecutionWarning warnings = 5;
}

// ExecutionDetails provides detailed information about rule execution
message ExecutionDetails {
  // Variables accessed during execution
  repeated string accessed_variables = 1;

  // Intermediate calculation results
  map<string, string> intermediate_results = 2;

  // Execution path taken through the rule
  repeated string execution_path = 3;

  // Memory usage during execution
  int64 memory_bytes = 4;
}

// EvaluationMetrics contains performance metrics for the evaluation
message EvaluationMetrics {
  // Parse time in milliseconds
  double parse_time_ms = 1;

  // Validation time in milliseconds
  double validation_time_ms = 2;

  // Execution time in milliseconds
  double execution_time_ms = 3;

  // Total processing time in milliseconds
  double total_time_ms = 4;

  // Number of operations performed
  int32 operation_count = 5;
}

// ExecutionWarning represents a warning during rule execution
message ExecutionWarning {
  // Warning code
  string code = 1;

  // Human-readable warning message
  string message = 2;

  // Context where the warning occurred
  string context = 3;
}

// BatchEvaluateRulesRequest defines batch evaluation parameters
message BatchEvaluateRulesRequest {
  // Rules to evaluate with their contexts
  repeated RuleEvaluationRequest rule_requests = 1 [
    (validate.rules).repeated.min_items = 1,
    (validate.rules).repeated.max_items = 100
  ];

  // Include performance metrics for all evaluations
  bool include_metrics = 2;

  // Execution timeout in milliseconds
  int32 timeout_ms = 3 [(validate.rules).int32.gte = 100];

  // Request ID for tracing
  string request_id = 4 [(validate.rules).string.uuid = true];
}

// RuleEvaluationRequest represents a single rule evaluation in a batch
message RuleEvaluationRequest {
  // Rule ID to evaluate
  string rule_id = 1 [(validate.rules).string.uuid = true];

  // Context data for this specific rule evaluation
  map<string, string> context = 2;
}

// BatchEvaluateRulesResponse contains batch evaluation results
message BatchEvaluateRulesResponse {
  // Individual evaluation results
  repeated RuleEvaluationResult results = 1;

  // Overall batch execution metrics
  BatchExecutionMetrics batch_metrics = 2;

  // Any batch-level warnings or errors
  repeated ExecutionWarning warnings = 3;
}

// RuleEvaluationResult represents the result of a single rule evaluation
message RuleEvaluationResult {
  // Rule ID that was evaluated
  string rule_id = 1;

  // Evaluation result (true/false)
  bool result = 2;

  // Execution time for this specific rule
  double execution_time_ms = 3;

  // Any errors that occurred during evaluation
  repeated ExecutionError errors = 4;

  // Execution details if available
  ExecutionDetails execution_details = 5;
}

// BatchExecutionMetrics contains metrics for the entire batch operation
message BatchExecutionMetrics {
  // Total batch processing time
  double total_time_ms = 1;

  // Average execution time per rule
  double average_time_ms = 2;

  // Number of successful evaluations
  int32 successful_count = 3;

  // Number of failed evaluations
  int32 failed_count = 4;

  // Total memory used for batch processing
  int64 total_memory_bytes = 5;
}

// ExecutionError represents an error during rule execution
message ExecutionError {
  // Error code
  string code = 1;

  // Human-readable error message
  string message = 2;

  // Error context and location
  string context = 3;

  // Stack trace if available
  string stack_trace = 4;
}

// ValidateRuleRequest defines rule validation parameters
message ValidateRuleRequest {
  // Rule ID to validate (if validating existing rule)
  string rule_id = 1 [(validate.rules).string.uuid = true];

  // DSL content to validate (if validating new content)
  string dsl_content = 2;

  // Include semantic validation (variable references, types)
  bool include_semantic_validation = 3;

  // Test data for execution validation
  map<string, string> test_data = 4;
}

// ValidateRuleResponse contains validation results
message ValidateRuleResponse {
  // Overall validation result
  ValidationResult validation_result = 1;

  // Test execution result if test data was provided
  EvaluateRuleResponse test_result = 2;
}

// ValidationResult contains detailed validation information
message ValidationResult {
  // Whether the rule is valid
  bool is_valid = 1;

  // Validation errors
  repeated ValidationError errors = 2;

  // Validation warnings
  repeated ValidationWarning warnings = 3;

  // Validation performance metrics
  ValidationMetrics metrics = 4;
}

// ValidationError represents a validation error
message ValidationError {
  // Error code
  string code = 1;

  // Human-readable error message
  string message = 2;

  // Line number where error occurred
  int32 line = 3;

  // Column number where error occurred
  int32 column = 4;

  // Error context
  string context = 5;
}

// ValidationWarning represents a validation warning
message ValidationWarning {
  // Warning code
  string code = 1;

  // Human-readable warning message
  string message = 2;

  // Line number where warning occurred
  int32 line = 3;

  // Context information
  string context = 4;
}

// ValidationMetrics contains performance metrics for validation
message ValidationMetrics {
  // Syntax validation time in milliseconds
  double syntax_validation_time_ms = 1;

  // Semantic validation time in milliseconds
  double semantic_validation_time_ms = 2;

  // Total validation time in milliseconds
  double total_validation_time_ms = 3;
}

// SubscribeRuleChangesRequest defines subscription parameters
message SubscribeRuleChangesRequest {
  // Filter for specific rule IDs (empty for all rules)
  repeated string rule_ids = 1;

  // Filter for specific event types
  repeated RuleChangeEventType event_types = 2;

  // Client identifier for connection tracking
  string client_id = 3 [(validate.rules).string.min_len = 1];
}

// RuleChangeEvent represents a rule change notification
message RuleChangeEvent {
  // Type of change event
  RuleChangeEventType event_type = 1;

  // Rule that changed
  Rule rule = 2;

  // Previous rule state (for update events)
  Rule previous_rule = 3;

  // Timestamp when the change occurred
  google.protobuf.Timestamp timestamp = 4;

  // User who made the change
  string changed_by = 5;

  // Change reason or description
  string change_reason = 6;
}

// RuleChangeEventType defines the types of rule change events
enum RuleChangeEventType {
  RULE_CHANGE_EVENT_TYPE_UNSPECIFIED = 0;
  RULE_CHANGE_EVENT_TYPE_CREATED = 1;
  RULE_CHANGE_EVENT_TYPE_UPDATED = 2;
  RULE_CHANGE_EVENT_TYPE_DELETED = 3;
  RULE_CHANGE_EVENT_TYPE_ACTIVATED = 4;
  RULE_CHANGE_EVENT_TYPE_DEACTIVATED = 5;
  RULE_CHANGE_EVENT_TYPE_APPROVED = 6;
  RULE_CHANGE_EVENT_TYPE_REJECTED = 7;
}

// GetRuleMetricsRequest specifies the rule and time range for metrics
message GetRuleMetricsRequest {
  // Rule ID to get metrics for
  string rule_id = 1 [(validate.rules).string.uuid = true];

  // Start time for metrics range
  google.protobuf.Timestamp start_time = 2;

  // End time for metrics range
  google.protobuf.Timestamp end_time = 3;

  // Metrics granularity (e.g., hourly, daily)
  MetricsGranularity granularity = 4;
}

// MetricsGranularity defines the granularity for metrics aggregation
enum MetricsGranularity {
  METRICS_GRANULARITY_UNSPECIFIED = 0;
  METRICS_GRANULARITY_MINUTE = 1;
  METRICS_GRANULARITY_HOUR = 2;
  METRICS_GRANULARITY_DAY = 3;
  METRICS_GRANULARITY_WEEK = 4;
  METRICS_GRANULARITY_MONTH = 5;
}

// GetRuleMetricsResponse contains rule performance and usage metrics
message GetRuleMetricsResponse {
  // Rule metrics data
  RuleMetrics metrics = 1;

  // Time series data points
  repeated MetricsDataPoint time_series = 2;
}

// RuleMetrics contains comprehensive metrics for a rule
message RuleMetrics {
  // Total number of evaluations
  int64 total_evaluations = 1;

  // Number of successful evaluations
  int64 successful_evaluations = 2;

  // Number of failed evaluations
  int64 failed_evaluations = 3;

  // Success rate percentage
  double success_rate = 4;

  // Average execution time in milliseconds
  double average_execution_time_ms = 5;

  // 95th percentile execution time
  double p95_execution_time_ms = 6;

  // 99th percentile execution time
  double p99_execution_time_ms = 7;

  // Total execution time across all evaluations
  double total_execution_time_ms = 8;

  // Last evaluation timestamp
  google.protobuf.Timestamp last_evaluation = 9;

  // Most common errors
  repeated ErrorFrequency common_errors = 10;
}

// MetricsDataPoint represents a single data point in time series metrics
message MetricsDataPoint {
  // Timestamp for this data point
  google.protobuf.Timestamp timestamp = 1;

  // Number of evaluations in this time period
  int64 evaluation_count = 2;

  // Average execution time for this period
  double average_execution_time_ms = 3;

  // Success rate for this period
  double success_rate = 4;

  // Error count for this period
  int64 error_count = 5;
}

// ErrorFrequency represents the frequency of a specific error
message ErrorFrequency {
  // Error code
  string error_code = 1;

  // Error message
  string error_message = 2;

  // Number of occurrences
  int64 count = 3;

  // Percentage of total errors
  double percentage = 4;
}
```

## 8. Output Specifications

### Format Requirements:
- **OpenAPI 3.0**: Complete specifications with examples, validation schemas, and error responses
- **Protobuf**: Well-documented .proto files with validation rules and service documentation
- **Error Handling**: RFC 7807 Problem Details format with structured error responses
- **Documentation**: Comprehensive API documentation with integration guides

### Quality Criteria:
- **Contract Completeness**: 100% API operation coverage with examples
- **Validation**: Comprehensive input validation with business rule enforcement
- **Error Handling**: Structured error responses with appropriate HTTP status codes
- **Performance**: <200ms simple operations, <500ms complex operations
- **Backward Compatibility**: Semantic versioning with deprecation policies

## 9. Validation Checkpoints

### Pre-execution Validation:
- [ ] API design follows REST conventions and resource modeling
- [ ] gRPC service contracts meet performance requirements
- [ ] Error scenarios identified with proper status codes
- [ ] Validation rules cover all business requirements

### Mid-execution Validation:
- [ ] OpenAPI specifications complete with examples and validation
- [ ] Protobuf definitions include proper documentation and validation
- [ ] Error handling covers all identified scenarios
- [ ] API implementations match contract specifications exactly

### Post-execution Validation:
- [ ] Contract testing validates implementation compliance
- [ ] Performance testing meets SLA requirements
- [ ] Error scenarios tested with proper responses
- [ ] Documentation enables successful client integration
- [ ] Backward compatibility maintained across versions

## Implementation Tasks Breakdown

### Phase 1: API Design and Contract Definition (Days 1-3)
1. **REST API Design**
   - Define resource hierarchy and operations
   - Create OpenAPI 3.0 specifications
   - Design error response formats
   - Define validation schemas

2. **gRPC Service Design**
   - Define service interfaces and operations
   - Create Protobuf message definitions
   - Design streaming operations
   - Add performance optimization strategies

### Phase 2: Implementation and Validation (Days 4-6)
1. **API Implementation**
   - Implement REST endpoints with validation
   - Create gRPC service handlers
   - Add comprehensive error handling
   - Implement business rule validation

2. **Testing and Validation**
   - Create contract tests for all operations
   - Add integration tests with error scenarios
   - Implement performance testing
   - Validate backward compatibility

### Phase 3: Documentation and Integration (Days 7-8)
1. **Documentation Creation**
   - Generate API documentation with examples
   - Create integration guides and tutorials
   - Add troubleshooting and FAQ sections
   - Create client SDK documentation

2. **Integration Support**
   - Create client libraries/SDKs
   - Add monitoring and observability
   - Implement rate limiting and throttling
   - Create deployment and versioning guides

## Best Practices

### API Design:
- Follow REST conventions for resource modeling
- Use proper HTTP status codes for all scenarios
- Implement comprehensive input validation
- Design for backward compatibility

### Error Handling:
- Use RFC 7807 Problem Details format
- Provide actionable error messages
- Include correlation IDs for debugging
- Implement proper error logging

### Performance Optimization:
- Use appropriate response formats (JSON vs Protobuf)
- Implement caching strategies
- Add request timeout handling
- Monitor and optimize response times

### Documentation Quality:
- Provide complete examples for all operations
- Include integration guides and tutorials
- Document error scenarios and responses
- Maintain up-to-date API documentation
