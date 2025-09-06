# Calculator Service

## Overview
The Calculator Service provides shared mathematical operations, algorithm optimization, and performance-critical calculations used across all other microservices. It serves as a centralized calculation engine for complex mathematical operations, formula evaluation, and high-performance computing tasks.

## Domain Model

### Core Entities

#### Calculation Request
```go
type CalculationRequest struct {
    ID                    string                 `json:"id" gorm:"primaryKey"`
    RequestID             string                 `json:"request_id" gorm:"uniqueIndex"`
    ServiceName           string                 `json:"service_name" gorm:"index"`
    CalculationType       CalculationType       `json:"calculation_type" gorm:"not null"`
    
    // Input Data
    InputData             map[string]interface{} `json:"input_data" gorm:"serializer:json"`
    FormulaID             *string                `json:"formula_id,omitempty"`
    AlgorithmID           *string                `json:"algorithm_id,omitempty"`
    
    // Configuration
    Precision             int                    `json:"precision" gorm:"default:2"`
    RoundingMode          RoundingMode          `json:"rounding_mode" gorm:"default:'ROUND_HALF_UP'"`
    Currency              string                 `json:"currency" gorm:"default:'USD'"`
    UseCache              bool                   `json:"use_cache" gorm:"default:true"`
    CacheTTL              time.Duration          `json:"cache_ttl" gorm:"default:3600000000000"` // 1 hour
    
    // Performance Requirements
    TimeoutMs             int                    `json:"timeout_ms" gorm:"default:5000"`
    MaxMemoryMB           int                    `json:"max_memory_mb" gorm:"default:100"`
    Priority              CalculationPriority   `json:"priority" gorm:"default:'NORMAL'"`
    
    // Status
    Status                RequestStatus          `json:"status" gorm:"default:'PENDING'"`
    
    // Metadata
    RequestedAt           time.Time              `json:"requested_at"`
    StartedAt             *time.Time             `json:"started_at,omitempty"`
    CompletedAt           *time.Time             `json:"completed_at,omitempty"`
    CreatedBy             string                 `json:"created_by,omitempty"`
}

type CalculationType string
const (
    CalculationTypeBasic       CalculationType = "BASIC"
    CalculationTypeFormula     CalculationType = "FORMULA"
    CalculationTypeStatistics  CalculationType = "STATISTICS"
    CalculationTypeFinancial   CalculationType = "FINANCIAL"
    CalculationTypePercentage  CalculationType = "PERCENTAGE"
    CalculationTypeCompound    CalculationType = "COMPOUND"
    CalculationTypeOptimization CalculationType = "OPTIMIZATION"
    CalculationTypeML          CalculationType = "MACHINE_LEARNING"
)

type CalculationPriority string
const (
    PriorityLow      CalculationPriority = "LOW"
    PriorityNormal   CalculationPriority = "NORMAL"
    PriorityHigh     CalculationPriority = "HIGH"
    PriorityCritical CalculationPriority = "CRITICAL"
)

type RoundingMode string
const (
    RoundingModeUp       RoundingMode = "ROUND_UP"
    RoundingModeDown     RoundingMode = "ROUND_DOWN"
    RoundingModeHalfUp   RoundingMode = "ROUND_HALF_UP"
    RoundingModeHalfDown RoundingMode = "ROUND_HALF_DOWN"
    RoundingModeHalfEven RoundingMode = "ROUND_HALF_EVEN"
    RoundingModeFloor    RoundingMode = "FLOOR"
    RoundingModeCeiling  RoundingMode = "CEILING"
)

type RequestStatus string
const (
    StatusPending   RequestStatus = "PENDING"
    StatusProcessing RequestStatus = "PROCESSING"
    StatusCompleted RequestStatus = "COMPLETED"
    StatusFailed    RequestStatus = "FAILED"
    StatusTimeout   RequestStatus = "TIMEOUT"
    StatusCancelled RequestStatus = "CANCELLED"
)
```

#### Calculation Result
```go
type CalculationResult struct {
    ID                    string                 `json:"id" gorm:"primaryKey"`
    RequestID             string                 `json:"request_id" gorm:"uniqueIndex"`
    CalculationRequestID  string                 `json:"calculation_request_id" gorm:"index"`
    
    // Results
    Success               bool                   `json:"success"`
    Result                interface{}            `json:"result" gorm:"serializer:json"`
    ResultType            string                 `json:"result_type"`
    
    // Multiple Results (for batch operations)
    Results               []CalculationOutput    `json:"results" gorm:"serializer:json"`
    
    // Calculation Details
    Formula               string                 `json:"formula,omitempty"`
    Steps                 []CalculationStep      `json:"steps" gorm:"serializer:json"`
    IntermediateValues    map[string]interface{} `json:"intermediate_values" gorm:"serializer:json"`
    
    // Performance Metrics
    ExecutionTime         time.Duration          `json:"execution_time"`
    MemoryUsed            int64                  `json:"memory_used"` // bytes
    CPUTime               time.Duration          `json:"cpu_time"`
    CacheHit              bool                   `json:"cache_hit"`
    
    // Error Information
    ErrorCode             string                 `json:"error_code,omitempty"`
    ErrorMessage          string                 `json:"error_message,omitempty"`
    ErrorDetails          map[string]interface{} `json:"error_details" gorm:"serializer:json"`
    
    // Validation
    ValidationResults     []ValidationResult     `json:"validation_results" gorm:"serializer:json"`
    ConfidenceScore       *float64               `json:"confidence_score,omitempty"`
    
    // Metadata
    CalculatedAt          time.Time              `json:"calculated_at"`
    ExpiresAt             *time.Time             `json:"expires_at,omitempty"`
    AlgorithmVersion      string                 `json:"algorithm_version,omitempty"`
}

type CalculationOutput struct {
    Key    string      `json:"key"`
    Value  interface{} `json:"value"`
    Type   string      `json:"type"`
    Unit   string      `json:"unit,omitempty"`
}

type CalculationStep struct {
    StepNumber   int         `json:"step_number"`
    Operation    string      `json:"operation"`
    Input        interface{} `json:"input"`
    Output       interface{} `json:"output"`
    Description  string      `json:"description,omitempty"`
    ExecutionTime time.Duration `json:"execution_time"`
}

type ValidationResult struct {
    ValidationType string      `json:"validation_type"`
    IsValid        bool        `json:"is_valid"`
    Message        string      `json:"message,omitempty"`
    ExpectedValue  interface{} `json:"expected_value,omitempty"`
    ActualValue    interface{} `json:"actual_value,omitempty"`
}
```

#### Mathematical Formula
```go
type MathematicalFormula struct {
    ID                    string                 `json:"id" gorm:"primaryKey"`
    Name                  string                 `json:"name" gorm:"not null"`
    Description           string                 `json:"description"`
    Category              FormulaCategory       `json:"category" gorm:"not null"`
    
    // Formula Definition
    Expression            string                 `json:"expression" gorm:"not null"`
    Variables             []FormulaVariable      `json:"variables" gorm:"serializer:json"`
    Constants             map[string]float64     `json:"constants" gorm:"serializer:json"`
    
    // Constraints and Validation
    Constraints           []FormulaConstraint    `json:"constraints" gorm:"serializer:json"`
    ValidationRules       []ValidationRule       `json:"validation_rules" gorm:"serializer:json"`
    
    // Output Configuration
    OutputType            string                 `json:"output_type" gorm:"default:'NUMERIC'"`
    OutputUnit            string                 `json:"output_unit,omitempty"`
    DefaultPrecision      int                    `json:"default_precision" gorm:"default:2"`
    
    // Performance
    Complexity            ComplexityLevel        `json:"complexity" gorm:"default:'MEDIUM'"`
    AverageExecutionTime  time.Duration          `json:"average_execution_time"`
    
    // Usage Tracking
    UsageCount            int64                  `json:"usage_count" gorm:"default:0"`
    LastUsed              *time.Time             `json:"last_used,omitempty"`
    
    // Status
    Status                FormulaStatus          `json:"status" gorm:"default:'ACTIVE'"`
    Version               string                 `json:"version" gorm:"default:'1.0'"`
    
    // Metadata
    CreatedAt             time.Time              `json:"created_at"`
    UpdatedAt             time.Time              `json:"updated_at"`
    CreatedBy             string                 `json:"created_by"`
}

type FormulaCategory string
const (
    CategoryFinancial    FormulaCategory = "FINANCIAL"
    CategoryStatistical  FormulaCategory = "STATISTICAL"
    CategoryGeometric    FormulaCategory = "GEOMETRIC"
    CategoryAlgebraic    FormulaCategory = "ALGEBRAIC"
    CategoryTrigonometric FormulaCategory = "TRIGONOMETRIC"
    CategoryLogic        FormulaCategory = "LOGIC"
    CategoryBusiness     FormulaCategory = "BUSINESS"
    CategoryCustom       FormulaCategory = "CUSTOM"
)

type ComplexityLevel string
const (
    ComplexityLow    ComplexityLevel = "LOW"
    ComplexityMedium ComplexityLevel = "MEDIUM"
    ComplexityHigh   ComplexityLevel = "HIGH"
    ComplexityCritical ComplexityLevel = "CRITICAL"
)

type FormulaStatus string
const (
    FormulaStatusActive     FormulaStatus = "ACTIVE"
    FormulaStatusInactive   FormulaStatus = "INACTIVE"
    FormulaStatusDeprecated FormulaStatus = "DEPRECATED"
    FormulaStatusDraft      FormulaStatus = "DRAFT"
)

type FormulaVariable struct {
    Name         string      `json:"name"`
    Type         string      `json:"type"` // "number", "boolean", "string", "array"
    Required     bool        `json:"required"`
    DefaultValue interface{} `json:"default_value,omitempty"`
    MinValue     *float64    `json:"min_value,omitempty"`
    MaxValue     *float64    `json:"max_value,omitempty"`
    Description  string      `json:"description,omitempty"`
}

type FormulaConstraint struct {
    ConstraintType string      `json:"constraint_type"` // "range", "enum", "pattern", "dependency"
    Field          string      `json:"field"`
    Condition      string      `json:"condition"`
    Value          interface{} `json:"value"`
    ErrorMessage   string      `json:"error_message"`
}

type ValidationRule struct {
    RuleType     string      `json:"rule_type"` // "range_check", "type_check", "business_logic"
    Expression   string      `json:"expression"`
    ErrorMessage string      `json:"error_message"`
    Severity     string      `json:"severity"` // "ERROR", "WARNING", "INFO"
}
```

#### Optimization Algorithm
```go
type OptimizationAlgorithm struct {
    ID                    string                 `json:"id" gorm:"primaryKey"`
    Name                  string                 `json:"name" gorm:"not null"`
    Description           string                 `json:"description"`
    AlgorithmType         AlgorithmType         `json:"algorithm_type" gorm:"not null"`
    
    // Algorithm Configuration
    Implementation        string                 `json:"implementation"` // Go code or reference
    Parameters            map[string]interface{} `json:"parameters" gorm:"serializer:json"`
    
    // Performance Characteristics
    TimeComplexity        string                 `json:"time_complexity"`
    SpaceComplexity       string                 `json:"space_complexity"`
    OptimalInputSize      *int                   `json:"optimal_input_size,omitempty"`
    MaxInputSize          *int                   `json:"max_input_size,omitempty"`
    
    // Quality Metrics
    Accuracy              float64                `json:"accuracy"`
    Precision             float64                `json:"precision"`
    Recall                float64                `json:"recall"`
    
    // Usage and Performance
    BenchmarkResults      []BenchmarkResult      `json:"benchmark_results" gorm:"serializer:json"`
    UsageStatistics       AlgorithmUsageStats    `json:"usage_statistics" gorm:"embedded"`
    
    // Status
    Status                AlgorithmStatus        `json:"status" gorm:"default:'ACTIVE'"`
    Version               string                 `json:"version" gorm:"default:'1.0'"`
    
    // Metadata
    CreatedAt             time.Time              `json:"created_at"`
    UpdatedAt             time.Time              `json:"updated_at"`
    CreatedBy             string                 `json:"created_by"`
}

type AlgorithmType string
const (
    AlgorithmTypeOptimization AlgorithmType = "OPTIMIZATION"
    AlgorithmTypeSorting      AlgorithmType = "SORTING"
    AlgorithmTypeSearch       AlgorithmType = "SEARCH"
    AlgorithmTypeML           AlgorithmType = "MACHINE_LEARNING"
    AlgorithmTypeStatistics   AlgorithmType = "STATISTICS"
    AlgorithmTypeNumerical    AlgorithmType = "NUMERICAL"
)

type AlgorithmStatus string
const (
    AlgorithmStatusActive      AlgorithmStatus = "ACTIVE"
    AlgorithmStatusInactive    AlgorithmStatus = "INACTIVE"
    AlgorithmStatusExperimental AlgorithmStatus = "EXPERIMENTAL"
    AlgorithmStatusDeprecated  AlgorithmStatus = "DEPRECATED"
)

type BenchmarkResult struct {
    InputSize       int           `json:"input_size"`
    ExecutionTime   time.Duration `json:"execution_time"`
    MemoryUsage     int64         `json:"memory_usage"`
    Accuracy        float64       `json:"accuracy"`
    Timestamp       time.Time     `json:"timestamp"`
}

type AlgorithmUsageStats struct {
    TotalExecutions      int64         `json:"total_executions"`
    AverageExecutionTime time.Duration `json:"average_execution_time"`
    SuccessRate          float64       `json:"success_rate"`
    LastExecuted         *time.Time    `json:"last_executed,omitempty"`
}
```

## gRPC Service Definition

```protobuf
syntax = "proto3";

package calculator;

service CalculatorService {
  // Basic calculation methods
  rpc Calculate(CalculateRequest) returns (CalculateResponse);
  rpc BatchCalculate(BatchCalculateRequest) returns (BatchCalculateResponse);
  rpc CalculateFormula(CalculateFormulaRequest) returns (CalculateFormulaResponse);
  
  // Optimization algorithms
  rpc RunOptimization(OptimizationRequest) returns (OptimizationResponse);
  rpc GetOptimizationStatus(GetOptimizationStatusRequest) returns (GetOptimizationStatusResponse);
  
  // Formula management
  rpc GetFormula(GetFormulaRequest) returns (GetFormulaResponse);
  rpc ValidateFormula(ValidateFormulaRequest) returns (ValidateFormulaResponse);
  rpc ListFormulas(ListFormulasRequest) returns (ListFormulasResponse);
  
  // Performance and metrics
  rpc GetPerformanceMetrics(GetPerformanceMetricsRequest) returns (GetPerformanceMetricsResponse);
  rpc GetCalculationHistory(GetCalculationHistoryRequest) returns (GetCalculationHistoryResponse);
}

message CalculateRequest {
  string request_id = 1;
  string calculation_type = 2;
  map<string, string> input_data = 3;
  int32 precision = 4;
  string rounding_mode = 5;
  bool use_cache = 6;
  int32 timeout_ms = 7;
}

message CalculateResponse {
  bool success = 1;
  string result = 2;
  string result_type = 3;
  repeated CalculationStep steps = 4;
  int64 execution_time_ms = 5;
  bool cache_hit = 6;
  string error_message = 7;
}

message BatchCalculateRequest {
  repeated CalculateRequest requests = 1;
  bool parallel_execution = 2;
  int32 max_concurrency = 3;
}

message BatchCalculateResponse {
  repeated CalculateResponse responses = 1;
  int64 total_execution_time_ms = 2;
  int32 successful_calculations = 3;
  int32 failed_calculations = 4;
}

message CalculateFormulaRequest {
  string formula_id = 1;
  map<string, string> variables = 2;
  int32 precision = 3;
  bool validate_inputs = 4;
  bool return_steps = 5;
}

message CalculateFormulaResponse {
  bool success = 1;
  string result = 2;
  repeated CalculationStep steps = 3;
  repeated ValidationResult validation_results = 4;
  int64 execution_time_ms = 5;
  string error_message = 6;
}
```

## Domain Services

### Calculation Engine Service
```go
type CalculationEngineService interface {
    Calculate(ctx context.Context, req *CalculationRequest) (*CalculationResult, error)
    BatchCalculate(ctx context.Context, requests []*CalculationRequest) ([]*CalculationResult, error)
    CalculateFormula(ctx context.Context, formulaID string, variables map[string]interface{}) (*CalculationResult, error)
    ValidateInputs(ctx context.Context, calculationType CalculationType, inputs map[string]interface{}) ([]ValidationResult, error)
    GetSupportedOperations(ctx context.Context) ([]string, error)
}
```

### Formula Engine Service
```go
type FormulaEngineService interface {
    EvaluateFormula(ctx context.Context, expression string, variables map[string]interface{}) (interface{}, error)
    ValidateFormula(ctx context.Context, expression string, variables []FormulaVariable) ([]ValidationResult, error)
    ParseFormula(ctx context.Context, expression string) (*ParsedFormula, error)
    OptimizeFormula(ctx context.Context, expression string) (string, error)
}

type ParsedFormula struct {
    Expression   string            `json:"expression"`
    Variables    []string          `json:"variables"`
    Constants    []string          `json:"constants"`
    Operations   []string          `json:"operations"`
    Complexity   ComplexityLevel   `json:"complexity"`
    Dependencies map[string]string `json:"dependencies"`
}
```

### Optimization Service
```go
type OptimizationService interface {
    RunOptimization(ctx context.Context, algorithmID string, parameters map[string]interface{}) (*OptimizationResult, error)
    GetOptimizationAlgorithms(ctx context.Context, algorithmType AlgorithmType) ([]OptimizationAlgorithm, error)
    BenchmarkAlgorithm(ctx context.Context, algorithmID string, testCases []interface{}) ([]BenchmarkResult, error)
    CompareAlgorithms(ctx context.Context, algorithmIDs []string, testCase interface{}) (*AlgorithmComparison, error)
}

type OptimizationResult struct {
    OptimalSolution   interface{}   `json:"optimal_solution"`
    OptimalValue      float64       `json:"optimal_value"`
    Iterations        int           `json:"iterations"`
    ExecutionTime     time.Duration `json:"execution_time"`
    ConvergenceStatus string        `json:"convergence_status"`
    Metadata          map[string]interface{} `json:"metadata"`
}

type AlgorithmComparison struct {
    Results       map[string]*OptimizationResult `json:"results"`
    BestAlgorithm string                         `json:"best_algorithm"`
    Metrics       map[string]interface{}         `json:"metrics"`
}
```

## Implementation Tasks

### Phase 1: Core Calculation Engine (3-4 days)
1. **Basic Mathematical Operations**
   - Implement arithmetic operations (add, subtract, multiply, divide)
   - Add advanced mathematical functions (power, root, logarithm)
   - Create percentage and proportion calculations
   - Implement rounding and precision handling

2. **Formula Evaluation Engine**
   - Create expression parser for mathematical formulas
   - Implement variable substitution and evaluation
   - Add function library (trigonometric, statistical, financial)
   - Create formula validation and error handling

### Phase 2: Advanced Calculations (3-4 days)
1. **Financial Calculations**
   - Implement compound interest calculations
   - Add present value and future value calculations
   - Create amortization and depreciation formulas
   - Implement ROI and performance metrics

2. **Statistical Operations**
   - Create basic statistics (mean, median, mode, standard deviation)
   - Implement correlation and regression analysis
   - Add probability distributions
   - Create data aggregation functions

### Phase 3: Optimization Algorithms (3-4 days)
1. **Optimization Framework**
   - Implement gradient descent algorithms
   - Add linear programming solvers
   - Create genetic algorithms for complex optimization
   - Implement constraint satisfaction solving

2. **Performance Optimization**
   - Add parallel processing for batch calculations
   - Implement result caching with TTL
   - Create memory-efficient algorithms
   - Add performance profiling and monitoring

### Phase 4: Formula Management (2-3 days)
1. **Formula Repository**
   - Create formula storage and retrieval system
   - Implement formula versioning and history
   - Add formula validation and testing
   - Create formula performance benchmarking

2. **Dynamic Formula Engine**
   - Implement runtime formula compilation
   - Add formula dependency management
   - Create formula optimization engine
   - Implement formula security validation

### Phase 5: Integration and Analytics (2-3 days)
1. **Service Integration**
   - Create gRPC service implementation
   - Add comprehensive error handling
   - Implement request/response caching
   - Create health monitoring and metrics

2. **Performance Analytics**
   - Implement calculation performance tracking
   - Add usage analytics and reporting
   - Create algorithm performance comparison
   - Implement predictive performance modeling

## Estimated Development Time: 13-18 days

## Performance Targets
- **Response Time**: <100ms for basic calculations, <1s for complex formulas
- **Throughput**: 10,000+ calculations per second
- **Accuracy**: 99.999% precision for financial calculations
- **Memory Usage**: <50MB per calculation request
- **Cache Hit Rate**: >90% for frequently used formulas
