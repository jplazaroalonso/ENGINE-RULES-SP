# Loyalty Service

## Overview
The Loyalty Service manages customer loyalty programs, points earning and redemption, tier management, and rewards calculation. It handles loyalty program lifecycle, customer tier progression, and reward fulfillment.

## Domain Model

### Core Entities

#### Loyalty Program
```go
type LoyaltyProgram struct {
    ID                    string                 `json:"id" gorm:"primaryKey"`
    Name                  string                 `json:"name" gorm:"not null"`
    Description           string                 `json:"description"`
    ProgramType           ProgramType           `json:"program_type" gorm:"not null"`
    
    // Program Configuration
    Currency              string                 `json:"currency" gorm:"default:'POINTS'"`
    PointsExpiryDays      *int                  `json:"points_expiry_days,omitempty"`
    MinimumRedemption     int                   `json:"minimum_redemption" gorm:"default:100"`
    MaximumRedemption     *int                  `json:"maximum_redemption,omitempty"`
    
    // Earning Rules
    EarningRules          []PointsEarningRule   `json:"earning_rules" gorm:"foreignKey:ProgramID"`
    
    // Tier Configuration
    Tiers                 []LoyaltyTier         `json:"tiers" gorm:"foreignKey:ProgramID"`
    TierCalculationMethod TierCalculationMethod `json:"tier_calculation_method" gorm:"default:'POINTS_BASED'"`
    TierReviewPeriod      TierReviewPeriod      `json:"tier_review_period" gorm:"default:'ANNUAL'"`
    
    // Redemption Options
    RedemptionOptions     []RedemptionOption    `json:"redemption_options" gorm:"foreignKey:ProgramID"`
    
    // Program Status
    Status                ProgramStatus         `json:"status" gorm:"default:'DRAFT'"`
    StartDate             *time.Time            `json:"start_date,omitempty"`
    EndDate               *time.Time            `json:"end_date,omitempty"`
    
    // Statistics
    Metrics               ProgramMetrics        `json:"metrics" gorm:"embedded"`
    
    // Metadata
    CreatedAt             time.Time             `json:"created_at"`
    UpdatedAt             time.Time             `json:"updated_at"`
    CreatedBy             string                `json:"created_by"`
}

type ProgramType string
const (
    ProgramTypePoints     ProgramType = "POINTS"
    ProgramTypeCashback   ProgramType = "CASHBACK"
    ProgramTypeTiered     ProgramType = "TIERED"
    ProgramTypeCoalition  ProgramType = "COALITION"
)

type ProgramStatus string
const (
    StatusDraft    ProgramStatus = "DRAFT"
    StatusActive   ProgramStatus = "ACTIVE"
    StatusInactive ProgramStatus = "INACTIVE"
    StatusArchived ProgramStatus = "ARCHIVED"
)

type TierCalculationMethod string
const (
    TierMethodPointsBased    TierCalculationMethod = "POINTS_BASED"
    TierMethodSpendBased     TierCalculationMethod = "SPEND_BASED"
    TierMethodActivityBased  TierCalculationMethod = "ACTIVITY_BASED"
    TierMethodHybrid         TierCalculationMethod = "HYBRID"
)

type TierReviewPeriod string
const (
    ReviewPeriodMonthly    TierReviewPeriod = "MONTHLY"
    ReviewPeriodQuarterly  TierReviewPeriod = "QUARTERLY"
    ReviewPeriodAnnual     TierReviewPeriod = "ANNUAL"
    ReviewPeriodLifetime   TierReviewPeriod = "LIFETIME"
)

type ProgramMetrics struct {
    TotalMembers         int     `json:"total_members"`
    ActiveMembers        int     `json:"active_members"`
    TotalPointsIssued    int64   `json:"total_points_issued"`
    TotalPointsRedeemed  int64   `json:"total_points_redeemed"`
    TotalPointsBalance   int64   `json:"total_points_balance"`
    AveragePointsPerMember int   `json:"average_points_per_member"`
    RedemptionRate       float64 `json:"redemption_rate"`
    EngagementRate       float64 `json:"engagement_rate"`
    LastUpdated          time.Time `json:"last_updated"`
}
```

#### Loyalty Tier
```go
type LoyaltyTier struct {
    ID                string                 `json:"id" gorm:"primaryKey"`
    ProgramID         string                 `json:"program_id" gorm:"index"`
    Name              string                 `json:"name" gorm:"not null"`
    Description       string                 `json:"description"`
    Level             int                    `json:"level" gorm:"not null"`
    Color             string                 `json:"color"`
    Icon              string                 `json:"icon,omitempty"`
    
    // Qualification Requirements
    RequiredPoints    *int                   `json:"required_points,omitempty"`
    RequiredSpend     *float64               `json:"required_spend,omitempty"`
    RequiredPurchases *int                   `json:"required_purchases,omitempty"`
    QualificationPeriod QualificationPeriod  `json:"qualification_period" gorm:"default:'ANNUAL'"`
    
    // Benefits
    PointsMultiplier  float64                `json:"points_multiplier" gorm:"default:1.0"`
    BonusPointsRate   float64                `json:"bonus_points_rate" gorm:"default:0.0"`
    Benefits          []TierBenefit          `json:"benefits" gorm:"serializer:json"`
    
    // Tier Maintenance
    RetentionPeriod   int                    `json:"retention_period" gorm:"default:365"`
    GracePeriod       int                    `json:"grace_period" gorm:"default:30"`
    
    // Status
    IsActive          bool                   `json:"is_active" gorm:"default:true"`
    
    // Metadata
    CreatedAt         time.Time              `json:"created_at"`
    UpdatedAt         time.Time              `json:"updated_at"`
}

type QualificationPeriod string
const (
    QualificationMonthly    QualificationPeriod = "MONTHLY"
    QualificationQuarterly  QualificationPeriod = "QUARTERLY"
    QualificationAnnual     QualificationPeriod = "ANNUAL"
    QualificationLifetime   QualificationPeriod = "LIFETIME"
)

type TierBenefit struct {
    Type        string      `json:"type"`        // "discount", "free_shipping", "priority_support"
    Value       interface{} `json:"value"`
    Description string      `json:"description"`
    IsActive    bool        `json:"is_active"`
}
```

#### Customer Loyalty Account
```go
type CustomerLoyaltyAccount struct {
    ID                    string                 `json:"id" gorm:"primaryKey"`
    CustomerID            string                 `json:"customer_id" gorm:"uniqueIndex"`
    ProgramID             string                 `json:"program_id" gorm:"index"`
    
    // Account Status
    Status                AccountStatus          `json:"status" gorm:"default:'ACTIVE'"`
    JoinedAt              time.Time              `json:"joined_at"`
    LastActivityAt        *time.Time             `json:"last_activity_at,omitempty"`
    
    // Points Balance
    TotalPointsEarned     int64                  `json:"total_points_earned" gorm:"default:0"`
    TotalPointsRedeemed   int64                  `json:"total_points_redeemed" gorm:"default:0"`
    CurrentPointsBalance  int64                  `json:"current_points_balance" gorm:"default:0"`
    PendingPoints         int64                  `json:"pending_points" gorm:"default:0"`
    ExpiringPoints        int64                  `json:"expiring_points" gorm:"default:0"`
    NextExpiryDate        *time.Time             `json:"next_expiry_date,omitempty"`
    
    // Tier Information
    CurrentTierID         *string                `json:"current_tier_id,omitempty"`
    CurrentTier           *LoyaltyTier          `json:"current_tier,omitempty" gorm:"foreignKey:CurrentTierID"`
    TierAchievedAt        *time.Time             `json:"tier_achieved_at,omitempty"`
    TierExpiresAt         *time.Time             `json:"tier_expires_at,omitempty"`
    NextTierID            *string                `json:"next_tier_id,omitempty"`
    NextTier              *LoyaltyTier          `json:"next_tier,omitempty" gorm:"foreignKey:NextTierID"`
    PointsToNextTier      int                    `json:"points_to_next_tier"`
    
    // Tier Qualification Progress
    CurrentPeriodSpend    float64                `json:"current_period_spend" gorm:"default:0"`
    CurrentPeriodPurchases int                   `json:"current_period_purchases" gorm:"default:0"`
    QualificationPeriodStart *time.Time          `json:"qualification_period_start,omitempty"`
    QualificationPeriodEnd   *time.Time          `json:"qualification_period_end,omitempty"`
    
    // Account Preferences
    Preferences           AccountPreferences     `json:"preferences" gorm:"embedded"`
    
    // Metadata
    CreatedAt             time.Time              `json:"created_at"`
    UpdatedAt             time.Time              `json:"updated_at"`
}

type AccountStatus string
const (
    AccountStatusActive    AccountStatus = "ACTIVE"
    AccountStatusInactive  AccountStatus = "INACTIVE"
    AccountStatusSuspended AccountStatus = "SUSPENDED"
    AccountStatusClosed    AccountStatus = "CLOSED"
)

type AccountPreferences struct {
    EmailNotifications    bool `json:"email_notifications" gorm:"default:true"`
    SMSNotifications      bool `json:"sms_notifications" gorm:"default:false"`
    PushNotifications     bool `json:"push_notifications" gorm:"default:true"`
    PointsExpiryReminders bool `json:"points_expiry_reminders" gorm:"default:true"`
    TierProgressUpdates   bool `json:"tier_progress_updates" gorm:"default:true"`
    RewardRecommendations bool `json:"reward_recommendations" gorm:"default:true"`
}
```

#### Points Transaction
```go
type PointsTransaction struct {
    ID                string                 `json:"id" gorm:"primaryKey"`
    CustomerID        string                 `json:"customer_id" gorm:"index"`
    AccountID         string                 `json:"account_id" gorm:"index"`
    
    // Transaction Details
    TransactionType   TransactionType        `json:"transaction_type" gorm:"not null"`
    Points            int64                  `json:"points" gorm:"not null"`
    PointsValue       float64                `json:"points_value,omitempty"`
    
    // Source Information
    SourceType        SourceType             `json:"source_type" gorm:"not null"`
    SourceID          string                 `json:"source_id,omitempty"`
    SourceDescription string                 `json:"source_description"`
    
    // Related Transaction
    RelatedTransactionID *string             `json:"related_transaction_id,omitempty"`
    OrderID           *string                `json:"order_id,omitempty"`
    
    // Earning Details (for earned points)
    EarningRuleID     *string                `json:"earning_rule_id,omitempty"`
    EarningRate       *float64               `json:"earning_rate,omitempty"`
    BaseAmount        *float64               `json:"base_amount,omitempty"`
    Multiplier        *float64               `json:"multiplier,omitempty"`
    
    // Redemption Details (for redeemed points)
    RedemptionID      *string                `json:"redemption_id,omitempty"`
    RewardID          *string                `json:"reward_id,omitempty"`
    RewardDescription *string                `json:"reward_description,omitempty"`
    
    // Expiry Information
    ExpiresAt         *time.Time             `json:"expires_at,omitempty"`
    IsExpired         bool                   `json:"is_expired" gorm:"default:false"`
    
    // Status
    Status            TransactionStatus      `json:"status" gorm:"default:'COMPLETED'"`
    
    // Metadata
    ProcessedAt       time.Time              `json:"processed_at"`
    CreatedAt         time.Time              `json:"created_at"`
}

type TransactionType string
const (
    TransactionTypeEarn      TransactionType = "EARN"
    TransactionTypeRedeem    TransactionType = "REDEEM"
    TransactionTypeAdjust    TransactionType = "ADJUST"
    TransactionTypeExpire    TransactionType = "EXPIRE"
    TransactionTypeTransfer  TransactionType = "TRANSFER"
    TransactionTypeBonus     TransactionType = "BONUS"
)

type SourceType string
const (
    SourceTypePurchase    SourceType = "PURCHASE"
    SourceTypeSignup      SourceType = "SIGNUP"
    SourceTypeReferral    SourceType = "REFERRAL"
    SourceTypeReview      SourceType = "REVIEW"
    SourceTypeSocial      SourceType = "SOCIAL"
    SourceTypePromotion   SourceType = "PROMOTION"
    SourceTypeAdmin       SourceType = "ADMIN"
    SourceTypeExpiry      SourceType = "EXPIRY"
    SourceTypeRedemption  SourceType = "REDEMPTION"
)

type TransactionStatus string
const (
    TransactionStatusPending   TransactionStatus = "PENDING"
    TransactionStatusCompleted TransactionStatus = "COMPLETED"
    TransactionStatusCancelled TransactionStatus = "CANCELLED"
    TransactionStatusReversed  TransactionStatus = "REVERSED"
)
```

## REST API Endpoints

### Loyalty Account Management

#### Get Customer Loyalty Account
```
GET /api/v1/customers/{customer_id}/loyalty

Response: 200 OK
{
  "id": "account-123",
  "customer_id": "customer-456",
  "program_id": "program-789",
  "status": "ACTIVE",
  "current_points_balance": 2450,
  "pending_points": 150,
  "expiring_points": 200,
  "next_expiry_date": "2024-03-15T00:00:00Z",
  "current_tier": {
    "id": "tier-gold",
    "name": "Gold",
    "level": 2,
    "points_multiplier": 1.5,
    "benefits": [
      {
        "type": "discount",
        "value": 10,
        "description": "10% discount on all purchases"
      }
    ]
  },
  "next_tier": {
    "id": "tier-platinum",
    "name": "Platinum",
    "level": 3
  },
  "points_to_next_tier": 550
}
```

#### Award Points
```
POST /api/v1/customers/{customer_id}/loyalty/points/award
Content-Type: application/json

{
  "points": 150,
  "source_type": "PURCHASE",
  "source_id": "order-789",
  "source_description": "Purchase reward for order #789",
  "earning_rule_id": "rule-123",
  "base_amount": 150.00,
  "earning_rate": 1.0,
  "multiplier": 1.5
}

Response: 201 Created
{
  "transaction_id": "txn-456",
  "points_awarded": 150,
  "new_balance": 2600,
  "tier_progression": {
    "tier_changed": false,
    "current_tier": "Gold",
    "points_to_next_tier": 400
  }
}
```

#### Redeem Points
```
POST /api/v1/customers/{customer_id}/loyalty/points/redeem
Content-Type: application/json

{
  "points": 500,
  "reward_id": "reward-123",
  "reward_description": "$5 discount coupon"
}

Response: 200 OK
{
  "transaction_id": "txn-789",
  "points_redeemed": 500,
  "new_balance": 2100,
  "redemption_id": "redemption-456",
  "reward_details": {
    "type": "coupon",
    "value": 5.00,
    "code": "LOYAL5OFF",
    "expires_at": "2024-06-15T23:59:59Z"
  }
}
```

### Points Transactions

#### Get Transaction History
```
GET /api/v1/customers/{customer_id}/loyalty/transactions?type=EARN&limit=20&page=1

Response: 200 OK
{
  "transactions": [
    {
      "id": "txn-123",
      "transaction_type": "EARN",
      "points": 150,
      "source_type": "PURCHASE",
      "source_description": "Purchase reward for order #789",
      "processed_at": "2024-01-15T10:30:00Z",
      "status": "COMPLETED"
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

### Tier Management

#### Check Tier Eligibility
```
POST /api/v1/customers/{customer_id}/loyalty/tier/check-eligibility

Response: 200 OK
{
  "current_tier": {
    "id": "tier-gold",
    "name": "Gold",
    "level": 2
  },
  "eligible_for_upgrade": true,
  "next_tier": {
    "id": "tier-platinum",
    "name": "Platinum",
    "level": 3
  },
  "qualification_progress": {
    "points_progress": {
      "current": 2600,
      "required": 3000,
      "percentage": 86.7
    },
    "spend_progress": {
      "current": 2400.00,
      "required": 2500.00,
      "percentage": 96.0
    }
  },
  "qualification_period": {
    "start": "2024-01-01T00:00:00Z",
    "end": "2024-12-31T23:59:59Z"
  }
}
```

## Implementation Tasks

### Phase 1: Core Loyalty Account Management (3-4 days)
1. **Domain Model and Database**
   - Implement loyalty program, tier, and account entities
   - Create database migrations with proper indexes
   - Add repository implementations with GORM
   - Implement domain validation and business rules

2. **Account Management**
   - Create customer loyalty account endpoints
   - Implement account status lifecycle management
   - Add tier assignment and progression logic
   - Create account preferences management

### Phase 2: Points Engine (3-4 days)
1. **Points Transaction System**
   - Implement points earning calculations
   - Create points redemption logic
   - Add points expiry management
   - Implement transaction history tracking

2. **Earning Rules Engine**
   - Create configurable earning rules
   - Implement multipliers and bonus calculations
   - Add activity-based earning logic
   - Create promotional earning campaigns

### Phase 3: Tier Management System (2-3 days)
1. **Tier Calculation Engine**
   - Implement tier qualification logic
   - Create automatic tier progression
   - Add tier retention and demotion rules
   - Implement tier benefit calculations

2. **Tier Analytics**
   - Create tier progression tracking
   - Implement tier performance metrics
   - Add tier distribution analytics
   - Create tier forecasting

### Phase 4: Rewards and Redemption (2-3 days)
1. **Redemption Engine**
   - Implement reward catalog management
   - Create redemption validation logic
   - Add redemption fulfillment tracking
   - Implement redemption analytics

2. **Reward Recommendations**
   - Create personalized reward suggestions
   - Implement reward availability checking
   - Add reward popularity tracking
   - Create reward performance analytics

### Phase 5: Integration and Analytics (2-3 days)
1. **Service Integration**
   - Integrate with customer service for account data
   - Add transaction service integration
   - Create promotional campaigns integration
   - Implement external partner integrations

2. **Analytics and Reporting**
   - Create loyalty program performance metrics
   - Implement customer engagement analytics
   - Add ROI and revenue impact tracking
   - Create predictive analytics for churn prevention

## Estimated Development Time: 12-17 days
