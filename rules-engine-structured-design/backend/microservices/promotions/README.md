# Promotions Service

## Overview
The Promotions Service manages promotional campaigns, discount rules, and customer targeting. It handles campaign lifecycle management, discount calculation, budget tracking, and performance analytics for marketing campaigns.

## Domain Model

### Core Entities

#### Promotional Campaign
```go
type PromotionalCampaign struct {
    ID              string                 `json:"id" gorm:"primaryKey"`
    Name            string                 `json:"name" gorm:"not null"`
    Description     string                 `json:"description"`
    CampaignType    CampaignType          `json:"campaign_type" gorm:"not null"`
    
    // Campaign Timing
    StartDate       time.Time             `json:"start_date"`
    EndDate         time.Time             `json:"end_date"`
    TimeZone        string                `json:"time_zone" gorm:"default:'UTC'"`
    
    // Targeting
    TargetSegments  []CustomerSegment     `json:"target_segments" gorm:"many2many:campaign_segments;"`
    TargetChannels  []string              `json:"target_channels" gorm:"serializer:json"`
    TargetRegions   []string              `json:"target_regions" gorm:"serializer:json"`
    
    // Budget and Limits
    BudgetLimit     *float64              `json:"budget_limit,omitempty"`
    UsageLimit      *int                  `json:"usage_limit,omitempty"`
    UsageLimitPerCustomer *int            `json:"usage_limit_per_customer,omitempty"`
    CurrentUsage    int                   `json:"current_usage" gorm:"default:0"`
    CurrentSpend    float64               `json:"current_spend" gorm:"default:0"`
    
    // Campaign Status
    Status          CampaignStatus        `json:"status" gorm:"default:'DRAFT'"`
    Priority        Priority              `json:"priority" gorm:"default:'MEDIUM'"`
    
    // Discount Configuration
    DiscountRules   []DiscountRule        `json:"discount_rules" gorm:"foreignKey:CampaignID"`
    
    // Analytics
    Metrics         CampaignMetrics       `json:"metrics" gorm:"embedded"`
    
    // Metadata
    CreatedAt       time.Time             `json:"created_at"`
    UpdatedAt       time.Time             `json:"updated_at"`
    CreatedBy       string                `json:"created_by"`
    ApprovedBy      *string               `json:"approved_by,omitempty"`
    ApprovedAt      *time.Time            `json:"approved_at,omitempty"`
}

type CampaignType string
const (
    CampaignTypePercentageDiscount CampaignType = "PERCENTAGE_DISCOUNT"
    CampaignTypeFixedDiscount      CampaignType = "FIXED_DISCOUNT"
    CampaignTypeBuyXGetY           CampaignType = "BUY_X_GET_Y"
    CampaignTypeBundleDiscount     CampaignType = "BUNDLE_DISCOUNT"
    CampaignTypeFreeShipping       CampaignType = "FREE_SHIPPING"
    CampaignTypeVolumeDiscount     CampaignType = "VOLUME_DISCOUNT"
)

type CampaignStatus string
const (
    StatusDraft      CampaignStatus = "DRAFT"
    StatusScheduled  CampaignStatus = "SCHEDULED"
    StatusActive     CampaignStatus = "ACTIVE"
    StatusPaused     CampaignStatus = "PAUSED"
    StatusCompleted  CampaignStatus = "COMPLETED"
    StatusCancelled  CampaignStatus = "CANCELLED"
)

type CampaignMetrics struct {
    TotalViews       int     `json:"total_views"`
    TotalClicks      int     `json:"total_clicks"`
    TotalRedemptions int     `json:"total_redemptions"`
    TotalRevenue     float64 `json:"total_revenue"`
    ConversionRate   float64 `json:"conversion_rate"`
    ClickThroughRate float64 `json:"click_through_rate"`
    AverageOrderValue float64 `json:"average_order_value"`
    ROI              float64 `json:"roi"`
    LastUpdated      time.Time `json:"last_updated"`
}
```

#### Discount Rule
```go
type DiscountRule struct {
    ID                string                 `json:"id" gorm:"primaryKey"`
    CampaignID        string                 `json:"campaign_id" gorm:"index"`
    Name              string                 `json:"name" gorm:"not null"`
    Description       string                 `json:"description"`
    
    // Discount Configuration
    DiscountType      DiscountType          `json:"discount_type" gorm:"not null"`
    DiscountValue     float64               `json:"discount_value" gorm:"not null"`
    MaxDiscountAmount *float64              `json:"max_discount_amount,omitempty"`
    MinOrderAmount    *float64              `json:"min_order_amount,omitempty"`
    
    // Product Targeting
    ApplicableProducts []ProductFilter      `json:"applicable_products" gorm:"serializer:json"`
    ExcludedProducts   []ProductFilter      `json:"excluded_products" gorm:"serializer:json"`
    
    // Buy X Get Y Configuration
    BuyQuantity       *int                  `json:"buy_quantity,omitempty"`
    GetQuantity       *int                  `json:"get_quantity,omitempty"`
    GetDiscountValue  *float64              `json:"get_discount_value,omitempty"`
    
    // Bundle Configuration
    RequiredProducts  []ProductBundle       `json:"required_products" gorm:"serializer:json"`
    
    // Volume Discount Tiers
    VolumeTiers       []VolumeTier          `json:"volume_tiers" gorm:"serializer:json"`
    
    // Rule Status
    IsActive          bool                  `json:"is_active" gorm:"default:true"`
    Priority          int                   `json:"priority" gorm:"default:100"`
    
    // Metadata
    CreatedAt         time.Time             `json:"created_at"`
    UpdatedAt         time.Time             `json:"updated_at"`
}

type DiscountType string
const (
    DiscountTypePercentage DiscountType = "PERCENTAGE"
    DiscountTypeFixed      DiscountType = "FIXED"
    DiscountTypeFreeShipping DiscountType = "FREE_SHIPPING"
    DiscountTypeBuyXGetY   DiscountType = "BUY_X_GET_Y"
    DiscountTypeBundle     DiscountType = "BUNDLE"
    DiscountTypeVolume     DiscountType = "VOLUME"
)

type ProductFilter struct {
    FilterType string      `json:"filter_type"` // "category", "brand", "sku", "tag"
    Values     []string    `json:"values"`
    Operator   string      `json:"operator"`    // "in", "not_in", "equals", "contains"
}

type ProductBundle struct {
    ProductID string `json:"product_id"`
    Quantity  int    `json:"quantity"`
    Required  bool   `json:"required"`
}

type VolumeTier struct {
    MinQuantity   int     `json:"min_quantity"`
    MaxQuantity   *int    `json:"max_quantity,omitempty"`
    DiscountValue float64 `json:"discount_value"`
}
```

#### Customer Segment
```go
type CustomerSegment struct {
    ID              string                 `json:"id" gorm:"primaryKey"`
    Name            string                 `json:"name" gorm:"not null"`
    Description     string                 `json:"description"`
    
    // Segmentation Criteria
    Criteria        []SegmentCriteria     `json:"criteria" gorm:"serializer:json"`
    
    // Customer Lists
    IncludedCustomers []string            `json:"included_customers" gorm:"serializer:json"`
    ExcludedCustomers []string            `json:"excluded_customers" gorm:"serializer:json"`
    
    // Segment Statistics
    EstimatedSize     int                 `json:"estimated_size"`
    ActualSize        int                 `json:"actual_size"`
    LastCalculated    *time.Time          `json:"last_calculated,omitempty"`
    
    // Status
    IsActive          bool                `json:"is_active" gorm:"default:true"`
    AutoUpdate        bool                `json:"auto_update" gorm:"default:false"`
    
    // Metadata
    CreatedAt         time.Time           `json:"created_at"`
    UpdatedAt         time.Time           `json:"updated_at"`
    CreatedBy         string              `json:"created_by"`
}

type SegmentCriteria struct {
    Field     string      `json:"field"`     // "age", "tier", "total_spent", "last_purchase"
    Operator  string      `json:"operator"`  // "equals", "greater_than", "less_than", "in", "between"
    Value     interface{} `json:"value"`
    LogicOp   string      `json:"logic_op"`  // "AND", "OR"
}
```

## REST API Endpoints

### Campaign Management

#### Create Campaign
```
POST /api/v1/campaigns
Content-Type: application/json

{
  "name": "Summer Sale 2024",
  "description": "25% off summer collection",
  "campaign_type": "PERCENTAGE_DISCOUNT",
  "start_date": "2024-06-01T00:00:00Z",
  "end_date": "2024-08-31T23:59:59Z",
  "target_segments": ["summer-shoppers", "loyal-customers"],
  "target_channels": ["web", "mobile"],
  "budget_limit": 50000.00,
  "usage_limit": 10000,
  "discount_rules": [
    {
      "name": "25% Summer Discount",
      "discount_type": "PERCENTAGE",
      "discount_value": 25.0,
      "min_order_amount": 100.0,
      "applicable_products": [
        {
          "filter_type": "category",
          "values": ["summer-clothing", "swimwear"],
          "operator": "in"
        }
      ]
    }
  ]
}

Response: 201 Created
{
  "id": "campaign-123",
  "name": "Summer Sale 2024",
  "status": "DRAFT",
  "created_at": "2024-01-15T10:00:00Z"
}
```

#### Get Campaigns
```
GET /api/v1/campaigns?status=ACTIVE&limit=20&page=1

Response: 200 OK
{
  "campaigns": [
    {
      "id": "campaign-123",
      "name": "Summer Sale 2024",
      "campaign_type": "PERCENTAGE_DISCOUNT",
      "status": "ACTIVE",
      "current_usage": 1250,
      "current_spend": 31250.00,
      "metrics": {
        "total_redemptions": 1250,
        "total_revenue": 156250.00,
        "conversion_rate": 0.18,
        "roi": 3.125
      }
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 45,
    "total_pages": 3
  }
}
```

#### Activate Campaign
```
POST /api/v1/campaigns/{campaign_id}/activate

Response: 200 OK
{
  "id": "campaign-123",
  "status": "ACTIVE",
  "activated_at": "2024-01-15T10:00:00Z"
}
```

### Discount Calculation

#### Calculate Discount
```
POST /api/v1/campaigns/calculate-discount
Content-Type: application/json

{
  "customer_id": "customer-123",
  "transaction": {
    "amount": 250.00,
    "currency": "USD",
    "channel": "web"
  },
  "products": [
    {
      "id": "product-456",
      "sku": "SUMMER-DRESS-001",
      "category": "summer-clothing",
      "price": 125.00,
      "quantity": 2
    }
  ],
  "campaign_ids": ["campaign-123"]
}

Response: 200 OK
{
  "applicable_campaigns": [
    {
      "campaign_id": "campaign-123",
      "campaign_name": "Summer Sale 2024",
      "applicable_rules": [
        {
          "rule_id": "rule-789",
          "rule_name": "25% Summer Discount",
          "discount_amount": 62.50,
          "discount_percentage": 25.0,
          "applicable_products": [
            {
              "product_id": "product-456",
              "original_price": 125.00,
              "discounted_price": 93.75,
              "discount_amount": 31.25
            }
          ]
        }
      ],
      "total_discount": 62.50
    }
  ],
  "total_discount": 62.50,
  "final_amount": 187.50
}
```

### Customer Segmentation

#### Create Segment
```
POST /api/v1/segments
Content-Type: application/json

{
  "name": "High Value Customers",
  "description": "Customers with high lifetime value",
  "criteria": [
    {
      "field": "total_spent",
      "operator": "greater_than",
      "value": 1000.0,
      "logic_op": "AND"
    },
    {
      "field": "tier",
      "operator": "in",
      "value": ["GOLD", "PLATINUM"],
      "logic_op": "AND"
    }
  ],
  "auto_update": true
}

Response: 201 Created
{
  "id": "segment-456",
  "name": "High Value Customers",
  "estimated_size": 2450,
  "created_at": "2024-01-15T10:00:00Z"
}
```

## Implementation Tasks

### Phase 1: Core Campaign Management (3-4 days)
1. **Domain Model and Database**
   - Implement campaign, discount rule, and segment entities
   - Create database migrations with proper indexes
   - Add repository implementations with GORM
   - Implement domain validation logic

2. **Campaign CRUD Operations**
   - Create campaign management endpoints
   - Implement campaign status lifecycle
   - Add campaign approval workflow
   - Create campaign scheduling functionality

### Phase 2: Discount Calculation Engine (3-4 days)
1. **Discount Calculation Logic**
   - Implement percentage and fixed discount calculations
   - Create Buy-X-Get-Y discount logic
   - Add bundle discount calculations
   - Implement volume discount tiers

2. **Product Targeting**
   - Create product filter evaluation engine
   - Implement category and brand targeting
   - Add SKU-based targeting
   - Create product exclusion logic

### Phase 3: Customer Segmentation (2-3 days)
1. **Segmentation Engine**
   - Implement criteria evaluation engine
   - Create customer matching algorithms
   - Add segment size calculation
   - Implement auto-update functionality

2. **Campaign Targeting**
   - Create segment-based campaign targeting
   - Implement channel and region targeting
   - Add customer inclusion/exclusion lists
   - Create targeting validation

### Phase 4: Analytics and Metrics (2-3 days)
1. **Campaign Analytics**
   - Implement real-time metrics tracking
   - Create campaign performance calculations
   - Add conversion rate tracking
   - Implement ROI calculations

2. **Reporting and Insights**
   - Create campaign performance reports
   - Add trend analysis and forecasting
   - Implement A/B testing support
   - Create performance dashboards

### Phase 5: Integration and Testing (2-3 days)
1. **Service Integration**
   - Integrate with Rules Calculation service
   - Add NATS event publishing
   - Create cache integration for performance
   - Implement external analytics integration

2. **Testing and Optimization**
   - Write comprehensive unit and integration tests
   - Implement performance testing
   - Add load testing for high-traffic campaigns
   - Create monitoring and alerting

## Estimated Development Time: 12-17 days
