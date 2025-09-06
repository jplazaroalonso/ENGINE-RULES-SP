# Coupons Service

## Overview
The Coupons Service manages coupon lifecycle, validation, redemption, and fraud prevention. It handles coupon campaigns, distribution strategies, usage tracking, and security measures to prevent abuse.

## Domain Model

### Core Entities

#### Coupon
```go
type Coupon struct {
    ID                    string                 `json:"id" gorm:"primaryKey"`
    Code                  string                 `json:"code" gorm:"uniqueIndex;not null"`
    Name                  string                 `json:"name" gorm:"not null"`
    Description           string                 `json:"description"`
    CouponType            CouponType            `json:"coupon_type" gorm:"not null"`
    
    // Discount Configuration
    DiscountType          DiscountType          `json:"discount_type" gorm:"not null"`
    DiscountValue         float64               `json:"discount_value" gorm:"not null"`
    MaxDiscountAmount     *float64              `json:"max_discount_amount,omitempty"`
    MinOrderAmount        *float64              `json:"min_order_amount,omitempty"`
    
    // Validity Period
    StartDate             time.Time             `json:"start_date"`
    EndDate               time.Time             `json:"end_date"`
    TimeZone              string                `json:"time_zone" gorm:"default:'UTC'"`
    
    // Usage Limits
    UsageLimit            *int                  `json:"usage_limit,omitempty"`
    UsageLimitPerCustomer *int                  `json:"usage_limit_per_customer,omitempty"`
    CurrentUsage          int                   `json:"current_usage" gorm:"default:0"`
    
    // Targeting
    ApplicableProducts    []ProductCriteria     `json:"applicable_products" gorm:"serializer:json"`
    ExcludedProducts      []ProductCriteria     `json:"excluded_products" gorm:"serializer:json"`
    CustomerSegments      []string              `json:"customer_segments" gorm:"serializer:json"`
    ChannelRestrictions   []string              `json:"channel_restrictions" gorm:"serializer:json"`
    
    // Distribution
    CampaignID            *string               `json:"campaign_id,omitempty"`
    DistributionMethod    DistributionMethod    `json:"distribution_method" gorm:"not null"`
    IsPublic              bool                  `json:"is_public" gorm:"default:false"`
    
    // Security
    SecuritySettings      SecuritySettings      `json:"security_settings" gorm:"embedded"`
    
    // Status
    Status                CouponStatus          `json:"status" gorm:"default:'ACTIVE'"`
    
    // Metadata
    CreatedAt             time.Time             `json:"created_at"`
    UpdatedAt             time.Time             `json:"updated_at"`
    CreatedBy             string                `json:"created_by"`
    
    // Analytics
    Metrics               CouponMetrics         `json:"metrics" gorm:"embedded"`
}

type CouponType string
const (
    CouponTypeSingleUse   CouponType = "SINGLE_USE"
    CouponTypeMultiUse    CouponType = "MULTI_USE"
    CouponTypeUnlimited   CouponType = "UNLIMITED"
    CouponTypeGeneric     CouponType = "GENERIC"
    CouponTypePersonalized CouponType = "PERSONALIZED"
)

type DiscountType string
const (
    DiscountTypePercentage    DiscountType = "PERCENTAGE"
    DiscountTypeFixed         DiscountType = "FIXED"
    DiscountTypeFreeShipping  DiscountType = "FREE_SHIPPING"
    DiscountTypeBuyXGetY      DiscountType = "BUY_X_GET_Y"
    DiscountTypeProduct       DiscountType = "PRODUCT"
)

type DistributionMethod string
const (
    DistributionMethodManual     DistributionMethod = "MANUAL"
    DistributionMethodEmail      DistributionMethod = "EMAIL"
    DistributionMethodSMS        DistributionMethod = "SMS"
    DistributionMethodApp        DistributionMethod = "APP"
    DistributionMethodWebsite    DistributionMethod = "WEBSITE"
    DistributionMethodPartner    DistributionMethod = "PARTNER"
)

type CouponStatus string
const (
    StatusActive    CouponStatus = "ACTIVE"
    StatusInactive  CouponStatus = "INACTIVE"
    StatusExpired   CouponStatus = "EXPIRED"
    StatusExhausted CouponStatus = "EXHAUSTED"
    StatusSuspended CouponStatus = "SUSPENDED"
)

type ProductCriteria struct {
    CriteriaType string   `json:"criteria_type"` // "category", "brand", "sku", "tag", "price_range"
    Values       []string `json:"values"`
    Operator     string   `json:"operator"`      // "in", "not_in", "equals", "greater_than", "less_than"
}

type SecuritySettings struct {
    IPRestrictions        []string `json:"ip_restrictions,omitempty" gorm:"serializer:json"`
    DeviceFingerprinting  bool     `json:"device_fingerprinting" gorm:"default:false"`
    VelocityChecks        bool     `json:"velocity_checks" gorm:"default:true"`
    MaxRedemptionsPerHour int      `json:"max_redemptions_per_hour" gorm:"default:10"`
    RequireEmailVerification bool  `json:"require_email_verification" gorm:"default:false"`
    BlockSuspiciousPatterns bool   `json:"block_suspicious_patterns" gorm:"default:true"`
}

type CouponMetrics struct {
    TotalViews           int     `json:"total_views"`
    TotalRedemptions     int     `json:"total_redemptions"`
    TotalRevenue         float64 `json:"total_revenue"`
    TotalDiscount        float64 `json:"total_discount"`
    ConversionRate       float64 `json:"conversion_rate"`
    AverageOrderValue    float64 `json:"average_order_value"`
    FraudAttempts        int     `json:"fraud_attempts"`
    BlockedRedemptions   int     `json:"blocked_redemptions"`
    LastRedemptionAt     *time.Time `json:"last_redemption_at,omitempty"`
    LastUpdated          time.Time `json:"last_updated"`
}
```

#### Coupon Redemption
```go
type CouponRedemption struct {
    ID                    string                 `json:"id" gorm:"primaryKey"`
    CouponID              string                 `json:"coupon_id" gorm:"index"`
    CouponCode            string                 `json:"coupon_code" gorm:"index"`
    CustomerID            string                 `json:"customer_id" gorm:"index"`
    
    // Transaction Details
    OrderID               string                 `json:"order_id" gorm:"index"`
    TransactionID         string                 `json:"transaction_id,omitempty"`
    
    // Redemption Details
    OriginalAmount        float64                `json:"original_amount"`
    DiscountAmount        float64                `json:"discount_amount"`
    FinalAmount           float64                `json:"final_amount"`
    Currency              string                 `json:"currency" gorm:"default:'USD'"`
    
    // Context
    Channel               string                 `json:"channel"`
    DeviceInfo            DeviceInfo             `json:"device_info" gorm:"embedded"`
    IPAddress             string                 `json:"ip_address"`
    UserAgent             string                 `json:"user_agent"`
    Location              *Location              `json:"location,omitempty" gorm:"embedded"`
    
    // Security Validation
    SecurityChecks        []SecurityCheck        `json:"security_checks" gorm:"serializer:json"`
    FraudScore            float64                `json:"fraud_score"`
    RiskLevel             RiskLevel              `json:"risk_level" gorm:"default:'LOW'"`
    
    // Status
    Status                RedemptionStatus       `json:"status" gorm:"default:'PENDING'"`
    ValidatedAt           *time.Time             `json:"validated_at,omitempty"`
    ProcessedAt           *time.Time             `json:"processed_at,omitempty"`
    
    // Metadata
    CreatedAt             time.Time              `json:"created_at"`
    UpdatedAt             time.Time              `json:"updated_at"`
}

type RedemptionStatus string
const (
    RedemptionStatusPending   RedemptionStatus = "PENDING"
    RedemptionStatusApproved  RedemptionStatus = "APPROVED"
    RedemptionStatusRejected  RedemptionStatus = "REJECTED"
    RedemptionStatusProcessed RedemptionStatus = "PROCESSED"
    RedemptionStatusReversed  RedemptionStatus = "REVERSED"
)

type RiskLevel string
const (
    RiskLevelLow      RiskLevel = "LOW"
    RiskLevelMedium   RiskLevel = "MEDIUM"
    RiskLevelHigh     RiskLevel = "HIGH"
    RiskLevelCritical RiskLevel = "CRITICAL"
)

type DeviceInfo struct {
    DeviceType    string `json:"device_type"`
    Platform      string `json:"platform"`
    Browser       string `json:"browser"`
    Version       string `json:"version"`
    Fingerprint   string `json:"fingerprint,omitempty"`
}

type Location struct {
    Country   string  `json:"country"`
    State     string  `json:"state,omitempty"`
    City      string  `json:"city,omitempty"`
    Latitude  float64 `json:"latitude,omitempty"`
    Longitude float64 `json:"longitude,omitempty"`
}

type SecurityCheck struct {
    CheckType   string    `json:"check_type"`
    Result      string    `json:"result"`
    Score       float64   `json:"score,omitempty"`
    Details     string    `json:"details,omitempty"`
    CheckedAt   time.Time `json:"checked_at"`
}
```

#### Coupon Campaign
```go
type CouponCampaign struct {
    ID                    string                 `json:"id" gorm:"primaryKey"`
    Name                  string                 `json:"name" gorm:"not null"`
    Description           string                 `json:"description"`
    CampaignType          CampaignType          `json:"campaign_type" gorm:"not null"`
    
    // Campaign Configuration
    CouponTemplate        CouponTemplate         `json:"coupon_template" gorm:"embedded"`
    GenerationRules       GenerationRules        `json:"generation_rules" gorm:"embedded"`
    
    // Distribution
    DistributionStrategy  DistributionStrategy   `json:"distribution_strategy" gorm:"embedded"`
    TargetAudience        TargetAudience         `json:"target_audience" gorm:"embedded"`
    
    // Campaign Limits
    BudgetLimit           *float64               `json:"budget_limit,omitempty"`
    CouponLimit           *int                   `json:"coupon_limit,omitempty"`
    TotalCouponsGenerated int                    `json:"total_coupons_generated" gorm:"default:0"`
    
    // Timing
    StartDate             time.Time              `json:"start_date"`
    EndDate               time.Time              `json:"end_date"`
    
    // Status
    Status                CampaignStatus         `json:"status" gorm:"default:'DRAFT'"`
    
    // Analytics
    Metrics               CampaignMetrics        `json:"metrics" gorm:"embedded"`
    
    // Metadata
    CreatedAt             time.Time              `json:"created_at"`
    UpdatedAt             time.Time              `json:"updated_at"`
    CreatedBy             string                 `json:"created_by"`
}

type CampaignType string
const (
    CampaignTypeBulk         CampaignType = "BULK"
    CampaignTypeTriggered    CampaignType = "TRIGGERED"
    CampaignTypePersonalized CampaignType = "PERSONALIZED"
    CampaignTypeReferral     CampaignType = "REFERRAL"
)

type CampaignStatus string
const (
    CampaignStatusDraft     CampaignStatus = "DRAFT"
    CampaignStatusScheduled CampaignStatus = "SCHEDULED"
    CampaignStatusActive    CampaignStatus = "ACTIVE"
    CampaignStatusPaused    CampaignStatus = "PAUSED"
    CampaignStatusCompleted CampaignStatus = "COMPLETED"
    CampaignStatusCancelled CampaignStatus = "CANCELLED"
)

type CouponTemplate struct {
    NamePattern       string       `json:"name_pattern"`
    CodePattern       string       `json:"code_pattern"`
    DiscountType      DiscountType `json:"discount_type"`
    DiscountValue     float64      `json:"discount_value"`
    ValidityDays      int          `json:"validity_days"`
    UsageLimit        *int         `json:"usage_limit,omitempty"`
}

type GenerationRules struct {
    CodeLength        int    `json:"code_length" gorm:"default:8"`
    CodePrefix        string `json:"code_prefix,omitempty"`
    CodeSuffix        string `json:"code_suffix,omitempty"`
    IncludeNumbers    bool   `json:"include_numbers" gorm:"default:true"`
    IncludeLetters    bool   `json:"include_letters" gorm:"default:true"`
    IncludeSpecialChars bool `json:"include_special_chars" gorm:"default:false"`
    ExcludeAmbiguous  bool   `json:"exclude_ambiguous" gorm:"default:true"`
}

type DistributionStrategy struct {
    Method          DistributionMethod `json:"method"`
    BatchSize       int               `json:"batch_size" gorm:"default:1000"`
    DistributionRate int              `json:"distribution_rate" gorm:"default:100"`
    AutoDistribute  bool              `json:"auto_distribute" gorm:"default:false"`
}

type TargetAudience struct {
    Segments          []string `json:"segments" gorm:"serializer:json"`
    CustomerIDs       []string `json:"customer_ids" gorm:"serializer:json"`
    IncludeNewCustomers bool   `json:"include_new_customers" gorm:"default:false"`
    ExcludeRecentUsers  bool   `json:"exclude_recent_users" gorm:"default:false"`
}

type CampaignMetrics struct {
    CouponsGenerated  int     `json:"coupons_generated"`
    CouponsDistributed int    `json:"coupons_distributed"`
    CouponsRedeemed   int     `json:"coupons_redeemed"`
    TotalRevenue      float64 `json:"total_revenue"`
    TotalDiscount     float64 `json:"total_discount"`
    ConversionRate    float64 `json:"conversion_rate"`
    RedemptionRate    float64 `json:"redemption_rate"`
    ROI               float64 `json:"roi"`
    LastUpdated       time.Time `json:"last_updated"`
}
```

## REST API Endpoints

### Coupon Validation and Redemption

#### Validate Coupon
```
POST /api/v1/coupons/validate
Content-Type: application/json

{
  "code": "SUMMER20",
  "customer_id": "customer-123",
  "order": {
    "amount": 150.00,
    "currency": "USD",
    "channel": "web",
    "products": [
      {
        "id": "product-456",
        "sku": "DRESS-001",
        "category": "clothing",
        "price": 75.00,
        "quantity": 2
      }
    ]
  },
  "context": {
    "ip_address": "192.168.1.1",
    "user_agent": "Mozilla/5.0...",
    "device_fingerprint": "abc123def456"
  }
}

Response: 200 OK
{
  "valid": true,
  "coupon": {
    "id": "coupon-789",
    "code": "SUMMER20",
    "name": "Summer Sale 20% Off",
    "discount_type": "PERCENTAGE",
    "discount_value": 20.0,
    "max_discount_amount": 50.0
  },
  "discount_calculation": {
    "applicable_amount": 150.00,
    "discount_amount": 30.00,
    "final_amount": 120.00
  },
  "security_assessment": {
    "risk_level": "LOW",
    "fraud_score": 0.15,
    "checks_passed": [
      "velocity_check",
      "device_fingerprint",
      "ip_reputation"
    ]
  },
  "usage_info": {
    "remaining_uses": 95,
    "customer_usage_count": 0,
    "expires_at": "2024-08-31T23:59:59Z"
  }
}
```

#### Redeem Coupon
```
POST /api/v1/coupons/redeem
Content-Type: application/json

{
  "code": "SUMMER20",
  "customer_id": "customer-123",
  "order_id": "order-456",
  "transaction_details": {
    "amount": 150.00,
    "currency": "USD",
    "channel": "web",
    "products": [...]
  },
  "context": {
    "ip_address": "192.168.1.1",
    "user_agent": "Mozilla/5.0...",
    "device_fingerprint": "abc123def456"
  }
}

Response: 200 OK
{
  "redemption_id": "redemption-789",
  "success": true,
  "coupon_id": "coupon-789",
  "discount_applied": 30.00,
  "final_amount": 120.00,
  "redemption_details": {
    "redeemed_at": "2024-01-15T10:30:00Z",
    "remaining_uses": 94,
    "security_score": 0.15
  }
}
```

### Coupon Management

#### Create Coupon
```
POST /api/v1/coupons
Content-Type: application/json

{
  "code": "WELCOME10",
  "name": "Welcome 10% Discount",
  "description": "New customer welcome discount",
  "coupon_type": "SINGLE_USE",
  "discount_type": "PERCENTAGE",
  "discount_value": 10.0,
  "start_date": "2024-01-01T00:00:00Z",
  "end_date": "2024-12-31T23:59:59Z",
  "usage_limit": 1000,
  "usage_limit_per_customer": 1,
  "customer_segments": ["new_customers"],
  "security_settings": {
    "velocity_checks": true,
    "max_redemptions_per_hour": 5,
    "block_suspicious_patterns": true
  }
}

Response: 201 Created
{
  "id": "coupon-123",
  "code": "WELCOME10",
  "status": "ACTIVE",
  "created_at": "2024-01-15T10:00:00Z"
}
```

#### Get Coupon Analytics
```
GET /api/v1/coupons/{coupon_id}/analytics

Response: 200 OK
{
  "coupon_id": "coupon-123",
  "metrics": {
    "total_views": 5420,
    "total_redemptions": 847,
    "total_revenue": 42350.00,
    "total_discount": 4235.00,
    "conversion_rate": 0.156,
    "average_order_value": 50.00,
    "fraud_attempts": 23,
    "blocked_redemptions": 12
  },
  "usage_timeline": [
    {
      "date": "2024-01-15",
      "redemptions": 45,
      "revenue": 2250.00
    }
  ],
  "top_products": [
    {
      "product_id": "product-123",
      "redemptions": 156,
      "revenue": 7800.00
    }
  ]
}
```

### Fraud Prevention

#### Get Fraud Report
```
GET /api/v1/coupons/fraud/report?start_date=2024-01-01&end_date=2024-01-31

Response: 200 OK
{
  "summary": {
    "total_attempts": 156,
    "blocked_attempts": 34,
    "fraud_rate": 0.218,
    "estimated_savings": 1250.00
  },
  "fraud_patterns": [
    {
      "pattern_type": "velocity_abuse",
      "occurrences": 12,
      "blocked": 10,
      "description": "Multiple redemption attempts in short time"
    },
    {
      "pattern_type": "device_spoofing",
      "occurrences": 8,
      "blocked": 7,
      "description": "Suspicious device fingerprint changes"
    }
  ],
  "high_risk_coupons": [
    {
      "coupon_id": "coupon-456",
      "code": "FLASH50",
      "fraud_attempts": 15,
      "risk_score": 0.85
    }
  ]
}
```

## Implementation Tasks

### Phase 1: Core Coupon Management (3-4 days)
1. **Domain Model and Database**
   - Implement coupon, redemption, and campaign entities
   - Create database migrations with proper indexes
   - Add repository implementations with GORM
   - Implement domain validation and business rules

2. **Coupon Lifecycle Management**
   - Create coupon CRUD operations
   - Implement coupon status management
   - Add expiry and usage limit tracking
   - Create bulk coupon operations

### Phase 2: Validation and Redemption Engine (3-4 days)
1. **Validation Engine**
   - Implement coupon code validation
   - Create eligibility checking logic
   - Add product applicability validation
   - Implement customer restriction validation

2. **Redemption Processing**
   - Create redemption transaction processing
   - Implement discount calculation logic
   - Add redemption status tracking
   - Create redemption reversal capability

### Phase 3: Security and Fraud Prevention (3-4 days)
1. **Security Framework**
   - Implement device fingerprinting
   - Add IP-based restrictions
   - Create velocity checking mechanisms
   - Implement suspicious pattern detection

2. **Fraud Detection**
   - Create fraud scoring algorithms
   - Implement risk assessment logic
   - Add machine learning-based fraud detection
   - Create fraud reporting and analytics

### Phase 4: Campaign Management (2-3 days)
1. **Campaign System**
   - Implement campaign creation and management
   - Create bulk coupon generation
   - Add distribution strategy implementation
   - Create campaign analytics and reporting

2. **Distribution Management**
   - Implement email and SMS distribution
   - Create API-based distribution
   - Add partner distribution channels
   - Create distribution tracking and analytics

### Phase 5: Analytics and Integration (2-3 days)
1. **Analytics Engine**
   - Create real-time usage analytics
   - Implement performance metrics tracking
   - Add coupon effectiveness analysis
   - Create predictive analytics for fraud

2. **External Integration**
   - Integrate with email/SMS services
   - Add payment gateway integration
   - Create partner API integrations
   - Implement analytics platform integration

## Estimated Development Time: 13-18 days
