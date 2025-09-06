# Payments Service

## Overview
The Payments Service manages payment rules, fraud detection, gateway routing, and payment optimization. It handles intelligent payment method selection, fraud prevention, and payment processing optimization based on customer behavior and transaction characteristics.

## Domain Model

### Core Entities

#### Payment Rule
```go
type PaymentRule struct {
    ID                    string                 `json:"id" gorm:"primaryKey"`
    Name                  string                 `json:"name" gorm:"not null"`
    Description           string                 `json:"description"`
    RuleType              PaymentRuleType       `json:"rule_type" gorm:"not null"`
    
    // Rule Configuration
    Conditions            []PaymentCondition     `json:"conditions" gorm:"serializer:json"`
    Actions               []PaymentAction        `json:"actions" gorm:"serializer:json"`
    
    // Targeting
    CustomerSegments      []string               `json:"customer_segments" gorm:"serializer:json"`
    PaymentMethods        []string               `json:"payment_methods" gorm:"serializer:json"`
    GeographicRegions     []string               `json:"geographic_regions" gorm:"serializer:json"`
    Currencies            []string               `json:"currencies" gorm:"serializer:json"`
    
    // Thresholds
    MinAmount             *float64               `json:"min_amount,omitempty"`
    MaxAmount             *float64               `json:"max_amount,omitempty"`
    TransactionLimits     TransactionLimits      `json:"transaction_limits" gorm:"embedded"`
    
    // Status and Priority
    Status                RuleStatus             `json:"status" gorm:"default:'ACTIVE'"`
    Priority              int                    `json:"priority" gorm:"default:100"`
    
    // Effectiveness Period
    EffectiveDate         time.Time              `json:"effective_date"`
    ExpiryDate            *time.Time             `json:"expiry_date,omitempty"`
    
    // Metadata
    CreatedAt             time.Time              `json:"created_at"`
    UpdatedAt             time.Time              `json:"updated_at"`
    CreatedBy             string                 `json:"created_by"`
    Version               int                    `json:"version" gorm:"default:1"`
}

type PaymentRuleType string
const (
    RuleTypeRouting       PaymentRuleType = "ROUTING"
    RuleTypeFraud         PaymentRuleType = "FRAUD"
    RuleTypeOptimization  PaymentRuleType = "OPTIMIZATION"
    RuleTypeCompliance    PaymentRuleType = "COMPLIANCE"
    RuleTypeFeCalculation PaymentRuleType = "FEE_CALCULATION"
    RuleTypeRetry         PaymentRuleType = "RETRY"
)

type PaymentCondition struct {
    Field     string      `json:"field"`     // "amount", "currency", "customer_tier", "payment_method"
    Operator  string      `json:"operator"`  // "equals", "greater_than", "less_than", "in", "not_in"
    Value     interface{} `json:"value"`
    LogicOp   string      `json:"logic_op"`  // "AND", "OR"
}

type PaymentAction struct {
    ActionType  string                 `json:"action_type"` // "route_to_gateway", "block", "require_3ds", "add_fee"
    Parameters  map[string]interface{} `json:"parameters"`
    Priority    int                    `json:"priority"`
}

type TransactionLimits struct {
    DailyLimit    *float64 `json:"daily_limit,omitempty"`
    WeeklyLimit   *float64 `json:"weekly_limit,omitempty"`
    MonthlyLimit  *float64 `json:"monthly_limit,omitempty"`
    CountPerDay   *int     `json:"count_per_day,omitempty"`
    CountPerHour  *int     `json:"count_per_hour,omitempty"`
}
```

#### Payment Gateway
```go
type PaymentGateway struct {
    ID                    string                 `json:"id" gorm:"primaryKey"`
    Name                  string                 `json:"name" gorm:"not null"`
    Provider              string                 `json:"provider" gorm:"not null"`
    GatewayType           GatewayType           `json:"gateway_type" gorm:"not null"`
    
    // Configuration
    Configuration         GatewayConfiguration   `json:"configuration" gorm:"embedded"`
    Credentials           GatewayCredentials     `json:"credentials" gorm:"embedded"`
    
    // Supported Features
    SupportedMethods      []PaymentMethod        `json:"supported_methods" gorm:"serializer:json"`
    SupportedCurrencies   []string               `json:"supported_currencies" gorm:"serializer:json"`
    SupportedCountries    []string               `json:"supported_countries" gorm:"serializer:json"`
    
    // Capabilities
    Supports3DS           bool                   `json:"supports_3ds" gorm:"default:false"`
    SupportsRefunds       bool                   `json:"supports_refunds" gorm:"default:true"`
    SupportsPartialRefunds bool                  `json:"supports_partial_refunds" gorm:"default:true"`
    SupportsVoid          bool                   `json:"supports_void" gorm:"default:true"`
    SupportsRecurring     bool                   `json:"supports_recurring" gorm:"default:false"`
    
    // Performance and Limits
    PerformanceMetrics    GatewayMetrics         `json:"performance_metrics" gorm:"embedded"`
    RateLimits            RateLimits             `json:"rate_limits" gorm:"embedded"`
    
    // Fees
    FeeStructure          FeeStructure           `json:"fee_structure" gorm:"embedded"`
    
    // Status
    Status                GatewayStatus          `json:"status" gorm:"default:'ACTIVE'"`
    HealthStatus          HealthStatus           `json:"health_status" gorm:"default:'HEALTHY'"`
    
    // Metadata
    CreatedAt             time.Time              `json:"created_at"`
    UpdatedAt             time.Time              `json:"updated_at"`
    LastHealthCheck       *time.Time             `json:"last_health_check,omitempty"`
}

type GatewayType string
const (
    GatewayTypePrimary   GatewayType = "PRIMARY"
    GatewayTypeSecondary GatewayType = "SECONDARY"
    GatewayTypeBackup    GatewayType = "BACKUP"
    GatewayTypeSpecialty GatewayType = "SPECIALTY"
)

type GatewayStatus string
const (
    GatewayStatusActive      GatewayStatus = "ACTIVE"
    GatewayStatusInactive    GatewayStatus = "INACTIVE"
    GatewayStatusMaintenance GatewayStatus = "MAINTENANCE"
    GatewayStatusDeprecated  GatewayStatus = "DEPRECATED"
)

type HealthStatus string
const (
    HealthStatusHealthy   HealthStatus = "HEALTHY"
    HealthStatusDegraded  HealthStatus = "DEGRADED"
    HealthStatusUnhealthy HealthStatus = "UNHEALTHY"
    HealthStatusUnknown   HealthStatus = "UNKNOWN"
)

type GatewayConfiguration struct {
    BaseURL           string        `json:"base_url"`
    TimeoutSeconds    int           `json:"timeout_seconds" gorm:"default:30"`
    RetryAttempts     int           `json:"retry_attempts" gorm:"default:3"`
    RetryDelay        time.Duration `json:"retry_delay" gorm:"default:1000000000"` // 1 second in nanoseconds
    EnableLogging     bool          `json:"enable_logging" gorm:"default:true"`
    LogLevel          string        `json:"log_level" gorm:"default:'INFO'"`
}

type GatewayCredentials struct {
    APIKey            string `json:"api_key,omitempty"`
    SecretKey         string `json:"secret_key,omitempty"`
    MerchantID        string `json:"merchant_id,omitempty"`
    PartnerID         string `json:"partner_id,omitempty"`
    CertificatePath   string `json:"certificate_path,omitempty"`
    WebhookSecret     string `json:"webhook_secret,omitempty"`
}

type PaymentMethod struct {
    Type        PaymentMethodType `json:"type"`
    SubTypes    []string          `json:"sub_types,omitempty"` // e.g., ["visa", "mastercard"] for credit_card
    MinAmount   *float64          `json:"min_amount,omitempty"`
    MaxAmount   *float64          `json:"max_amount,omitempty"`
    IsEnabled   bool              `json:"is_enabled" gorm:"default:true"`
}

type PaymentMethodType string
const (
    PaymentMethodCreditCard    PaymentMethodType = "CREDIT_CARD"
    PaymentMethodDebitCard     PaymentMethodType = "DEBIT_CARD"
    PaymentMethodDigitalWallet PaymentMethodType = "DIGITAL_WALLET"
    PaymentMethodBankTransfer  PaymentMethodType = "BANK_TRANSFER"
    PaymentMethodCrypto        PaymentMethodType = "CRYPTO"
    PaymentMethodBuyNowPayLater PaymentMethodType = "BUY_NOW_PAY_LATER"
)

type GatewayMetrics struct {
    SuccessRate       float64   `json:"success_rate"`
    AverageResponseTime time.Duration `json:"average_response_time"`
    TotalTransactions int64     `json:"total_transactions"`
    FailedTransactions int64    `json:"failed_transactions"`
    LastUpdated       time.Time `json:"last_updated"`
}

type RateLimits struct {
    RequestsPerSecond int `json:"requests_per_second" gorm:"default:100"`
    RequestsPerMinute int `json:"requests_per_minute" gorm:"default:6000"`
    RequestsPerHour   int `json:"requests_per_hour" gorm:"default:360000"`
}

type FeeStructure struct {
    FixedFee       float64 `json:"fixed_fee"`
    PercentageFee  float64 `json:"percentage_fee"`
    MinFee         float64 `json:"min_fee"`
    MaxFee         *float64 `json:"max_fee,omitempty"`
    Currency       string  `json:"currency" gorm:"default:'USD'"`
}
```

#### Fraud Rule
```go
type FraudRule struct {
    ID                    string                 `json:"id" gorm:"primaryKey"`
    Name                  string                 `json:"name" gorm:"not null"`
    Description           string                 `json:"description"`
    RuleType              FraudRuleType         `json:"rule_type" gorm:"not null"`
    
    // Rule Logic
    Conditions            []FraudCondition       `json:"conditions" gorm:"serializer:json"`
    Actions               []FraudAction          `json:"actions" gorm:"serializer:json"`
    
    // Scoring
    RiskWeight            float64                `json:"risk_weight" gorm:"default:1.0"`
    MaxRiskScore          float64                `json:"max_risk_score" gorm:"default:100.0"`
    ThresholdScore        float64                `json:"threshold_score" gorm:"default:50.0"`
    
    // Machine Learning
    MLModelID             *string                `json:"ml_model_id,omitempty"`
    UseMLPrediction       bool                   `json:"use_ml_prediction" gorm:"default:false"`
    MLWeight              float64                `json:"ml_weight" gorm:"default:0.3"`
    
    // Status and Configuration
    Status                RuleStatus             `json:"status" gorm:"default:'ACTIVE'"`
    Severity              FraudSeverity          `json:"severity" gorm:"default:'MEDIUM'"`
    Priority              int                    `json:"priority" gorm:"default:100"`
    
    // Performance Tracking
    Metrics               FraudRuleMetrics       `json:"metrics" gorm:"embedded"`
    
    // Metadata
    CreatedAt             time.Time              `json:"created_at"`
    UpdatedAt             time.Time              `json:"updated_at"`
    CreatedBy             string                 `json:"created_by"`
}

type FraudRuleType string
const (
    FraudRuleVelocity    FraudRuleType = "VELOCITY"
    FraudRuleGeolocation FraudRuleType = "GEOLOCATION" 
    FraudRuleDevice      FraudRuleType = "DEVICE"
    FraudRuleBehavioral  FraudRuleType = "BEHAVIORAL"
    FraudRuleTransaction FraudRuleType = "TRANSACTION"
    FraudRuleIdentity    FraudRuleType = "IDENTITY"
)

type FraudSeverity string
const (
    SeverityLow      FraudSeverity = "LOW"
    SeverityMedium   FraudSeverity = "MEDIUM"
    SeverityHigh     FraudSeverity = "HIGH"
    SeverityCritical FraudSeverity = "CRITICAL"
)

type FraudCondition struct {
    Field         string      `json:"field"`        // "transaction_amount", "velocity_count", "device_fingerprint"
    Operator      string      `json:"operator"`     // "greater_than", "equals", "not_in", "matches_pattern"
    Value         interface{} `json:"value"`
    TimeWindow    *string     `json:"time_window,omitempty"` // "1h", "24h", "7d"
    ComparisonType string     `json:"comparison_type,omitempty"` // "absolute", "relative", "pattern"
}

type FraudAction struct {
    ActionType    string                 `json:"action_type"` // "block", "review", "challenge", "adjust_score"
    Parameters    map[string]interface{} `json:"parameters"`
    Severity      FraudSeverity          `json:"severity"`
    Automatic     bool                   `json:"automatic" gorm:"default:true"`
}

type FraudRuleMetrics struct {
    TotalTriggers       int     `json:"total_triggers"`
    TruePositives       int     `json:"true_positives"`
    FalsePositives      int     `json:"false_positives"`
    TrueNegatives       int     `json:"true_negatives"`
    FalseNegatives      int     `json:"false_negatives"`
    Precision           float64 `json:"precision"`
    Recall              float64 `json:"recall"`
    F1Score             float64 `json:"f1_score"`
    AccuracyRate        float64 `json:"accuracy_rate"`
    LastEvaluated       *time.Time `json:"last_evaluated,omitempty"`
}
```

## REST API Endpoints

### Payment Intelligence

#### Get Payment Recommendations
```
POST /api/v1/payments/recommendations
Content-Type: application/json

{
  "customer_id": "customer-123",
  "transaction": {
    "amount": 150.00,
    "currency": "USD",
    "country": "US"
  },
  "customer_profile": {
    "tier": "GOLD",
    "payment_history": {
      "preferred_methods": ["credit_card", "digital_wallet"],
      "success_rate": 0.98,
      "average_amount": 125.50
    }
  }
}

Response: 200 OK
{
  "recommendations": [
    {
      "payment_method": "CREDIT_CARD",
      "gateway_id": "gateway-stripe",
      "gateway_name": "Stripe",
      "confidence_score": 0.95,
      "expected_success_rate": 0.97,
      "estimated_fee": 4.65,
      "processing_time": "2-3 seconds",
      "reasons": [
        "High success rate for this customer segment",
        "Optimal fees for this transaction amount",
        "Gateway performance is excellent"
      ]
    },
    {
      "payment_method": "DIGITAL_WALLET",
      "gateway_id": "gateway-paypal",
      "gateway_name": "PayPal",
      "confidence_score": 0.87,
      "expected_success_rate": 0.94,
      "estimated_fee": 5.10,
      "processing_time": "1-2 seconds"
    }
  ]
}
```

#### Process Fraud Check
```
POST /api/v1/payments/fraud/check
Content-Type: application/json

{
  "transaction_id": "txn-456",
  "customer_id": "customer-123",
  "payment_details": {
    "amount": 500.00,
    "currency": "USD",
    "payment_method": "CREDIT_CARD",
    "card_number": "4111****1111",
    "billing_address": {
      "country": "US",
      "state": "CA",
      "city": "San Francisco"
    }
  },
  "context": {
    "ip_address": "192.168.1.1",
    "user_agent": "Mozilla/5.0...",
    "device_fingerprint": "abc123def456",
    "session_id": "session-789"
  }
}

Response: 200 OK
{
  "fraud_assessment": {
    "overall_risk_score": 25.5,
    "risk_level": "LOW",
    "decision": "APPROVE",
    "confidence": 0.92
  },
  "triggered_rules": [
    {
      "rule_id": "rule-velocity-check",
      "rule_name": "Transaction Velocity Check",
      "risk_score": 15.0,
      "severity": "LOW",
      "details": "2 transactions in last hour (threshold: 5)"
    },
    {
      "rule_id": "rule-amount-check",
      "rule_name": "High Amount Transaction",
      "risk_score": 10.5,
      "severity": "LOW",
      "details": "Amount above customer average but within normal range"
    }
  ],
  "recommendations": [
    {
      "action": "PROCEED",
      "reason": "Low risk score, customer has good payment history",
      "additional_checks": []
    }
  ],
  "ml_prediction": {
    "fraud_probability": 0.08,
    "model_version": "fraud-detection-v2.1",
    "feature_importance": {
      "transaction_amount": 0.35,
      "customer_history": 0.25,
      "device_reputation": 0.20,
      "geolocation": 0.20
    }
  }
}
```

### Gateway Management

#### Get Gateway Health
```
GET /api/v1/payments/gateways/health

Response: 200 OK
{
  "overall_health": "HEALTHY",
  "last_updated": "2024-01-15T10:30:00Z",
  "gateways": [
    {
      "id": "gateway-stripe",
      "name": "Stripe",
      "status": "ACTIVE",
      "health_status": "HEALTHY",
      "performance_metrics": {
        "success_rate": 0.982,
        "average_response_time": "1.2s",
        "uptime": 0.999
      },
      "last_health_check": "2024-01-15T10:29:45Z"
    },
    {
      "id": "gateway-paypal",
      "name": "PayPal",
      "status": "ACTIVE", 
      "health_status": "DEGRADED",
      "performance_metrics": {
        "success_rate": 0.945,
        "average_response_time": "3.1s",
        "uptime": 0.995
      },
      "issues": [
        "Elevated response times detected",
        "Success rate below optimal threshold"
      ]
    }
  ]
}
```

## Implementation Tasks

### Phase 1: Core Payment Rules Engine (3-4 days)
1. **Domain Model and Database**
   - Implement payment rule, gateway, and fraud rule entities
   - Create database migrations with proper indexes
   - Add repository implementations with GORM
   - Implement rule evaluation engine

2. **Payment Rules Processing**
   - Create rule condition evaluation logic
   - Implement action execution framework
   - Add rule priority and conflict resolution
   - Create rule performance tracking

### Phase 2: Gateway Management System (3-4 days)
1. **Gateway Integration Framework**
   - Implement gateway abstraction layer
   - Create gateway health monitoring
   - Add gateway performance tracking
   - Implement gateway failover logic

2. **Intelligent Routing**
   - Create gateway selection algorithms
   - Implement load balancing strategies
   - Add cost optimization logic
   - Create success rate optimization

### Phase 3: Fraud Detection Engine (3-4 days)
1. **Fraud Rule Engine**
   - Implement fraud condition evaluation
   - Create real-time fraud scoring
   - Add machine learning integration
   - Implement fraud action processing

2. **Risk Assessment**
   - Create risk scoring algorithms
   - Implement behavioral analysis
   - Add device fingerprinting
   - Create geolocation risk assessment

### Phase 4: Payment Optimization (2-3 days)
1. **Success Rate Optimization**
   - Implement retry logic
   - Create payment method optimization
   - Add customer-specific routing
   - Implement A/B testing framework

2. **Cost Optimization**
   - Create fee calculation engine
   - Implement cost-based routing
   - Add fee optimization algorithms
   - Create cost analysis reporting

### Phase 5: Analytics and Integration (2-3 days)
1. **Payment Analytics**
   - Create performance dashboards
   - Implement fraud analytics
   - Add payment success tracking
   - Create predictive analytics

2. **External Integration**
   - Integrate with payment gateways
   - Add fraud detection services
   - Create webhook handling
   - Implement compliance reporting

## Estimated Development Time: 13-18 days
