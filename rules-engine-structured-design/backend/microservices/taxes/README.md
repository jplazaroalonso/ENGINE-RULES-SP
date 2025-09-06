# Taxes Service

## Overview
The Taxes Service manages tax calculation, jurisdiction rules, compliance reporting, and regulatory requirements. It handles complex tax scenarios including multi-jurisdictional taxation, tax exemptions, and real-time tax rate updates.

## Domain Model

### Core Entities

#### Tax Rule
```go
type TaxRule struct {
    ID                    string                 `json:"id" gorm:"primaryKey"`
    Name                  string                 `json:"name" gorm:"not null"`
    Description           string                 `json:"description"`
    TaxType               TaxType               `json:"tax_type" gorm:"not null"`
    
    // Jurisdiction
    JurisdictionID        string                 `json:"jurisdiction_id" gorm:"index"`
    Jurisdiction          TaxJurisdiction        `json:"jurisdiction" gorm:"foreignKey:JurisdictionID"`
    
    // Tax Configuration
    TaxRate               float64                `json:"tax_rate" gorm:"not null"`
    IsPercentage          bool                   `json:"is_percentage" gorm:"default:true"`
    MinimumAmount         *float64               `json:"minimum_amount,omitempty"`
    MaximumAmount         *float64               `json:"maximum_amount,omitempty"`
    
    // Applicability
    ApplicableCategories  []string               `json:"applicable_categories" gorm:"serializer:json"`
    ExcludedCategories    []string               `json:"excluded_categories" gorm:"serializer:json"`
    ApplicableProducts    []ProductTaxCriteria   `json:"applicable_products" gorm:"serializer:json"`
    
    // Tax Exemptions
    Exemptions            []TaxExemption         `json:"exemptions" gorm:"foreignKey:TaxRuleID"`
    
    // Effective Period
    EffectiveDate         time.Time              `json:"effective_date"`
    ExpiryDate            *time.Time             `json:"expiry_date,omitempty"`
    
    // Status
    Status                TaxRuleStatus          `json:"status" gorm:"default:'ACTIVE'"`
    Priority              int                    `json:"priority" gorm:"default:100"`
    
    // Metadata
    CreatedAt             time.Time              `json:"created_at"`
    UpdatedAt             time.Time              `json:"updated_at"`
    CreatedBy             string                 `json:"created_by"`
    Version               int                    `json:"version" gorm:"default:1"`
}

type TaxType string
const (
    TaxTypeSales       TaxType = "SALES_TAX"
    TaxTypeVAT         TaxType = "VAT"
    TaxTypeGST         TaxType = "GST"
    TaxTypeService     TaxType = "SERVICE_TAX"
    TaxTypeLuxury      TaxType = "LUXURY_TAX"
    TaxTypeExcise      TaxType = "EXCISE_TAX"
    TaxTypeImport      TaxType = "IMPORT_TAX"
    TaxTypeEnvironment TaxType = "ENVIRONMENT_TAX"
)

type TaxRuleStatus string
const (
    StatusActive   TaxRuleStatus = "ACTIVE"
    StatusInactive TaxRuleStatus = "INACTIVE"
    StatusPending  TaxRuleStatus = "PENDING"
    StatusArchived TaxRuleStatus = "ARCHIVED"
    StatusDraft    TaxRuleStatus = "DRAFT"
)

type ProductTaxCriteria struct {
    CriteriaType string   `json:"criteria_type"` // "category", "brand", "sku", "price_range", "weight"
    Values       []string `json:"values"`
    Operator     string   `json:"operator"`      // "in", "not_in", "equals", "greater_than", "less_than"
}
```

#### Tax Jurisdiction
```go
type TaxJurisdiction struct {
    ID                    string                 `json:"id" gorm:"primaryKey"`
    Name                  string                 `json:"name" gorm:"not null"`
    Code                  string                 `json:"code" gorm:"uniqueIndex;not null"`
    Type                  JurisdictionType      `json:"type" gorm:"not null"`
    
    // Hierarchy
    ParentID              *string                `json:"parent_id,omitempty"`
    Parent                *TaxJurisdiction       `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
    Children              []TaxJurisdiction      `json:"children,omitempty" gorm:"foreignKey:ParentID"`
    Level                 int                    `json:"level" gorm:"default:0"`
    Path                  string                 `json:"path"` // Hierarchical path like "US/CA/SF"
    
    // Geographic Boundaries
    Boundaries            *GeographicBoundary    `json:"boundaries,omitempty" gorm:"embedded"`
    
    // Tax Configuration
    DefaultTaxRate        float64                `json:"default_tax_rate" gorm:"default:0"`
    CompoundTax           bool                   `json:"compound_tax" gorm:"default:false"`
    TaxOnTax              bool                   `json:"tax_on_tax" gorm:"default:false"`
    
    // Registration Requirements
    RegistrationThreshold *float64               `json:"registration_threshold,omitempty"`
    RequiresRegistration  bool                   `json:"requires_registration" gorm:"default:false"`
    
    // Status
    IsActive              bool                   `json:"is_active" gorm:"default:true"`
    
    // Metadata
    CreatedAt             time.Time              `json:"created_at"`
    UpdatedAt             time.Time              `json:"updated_at"`
}

type JurisdictionType string
const (
    JurisdictionCountry JurisdictionType = "COUNTRY"
    JurisdictionState   JurisdictionType = "STATE"
    JurisdictionCounty  JurisdictionType = "COUNTY"
    JurisdictionCity    JurisdictionType = "CITY"
    JurisdictionZone    JurisdictionType = "ZONE"
    JurisdictionSpecial JurisdictionType = "SPECIAL"
)

type GeographicBoundary struct {
    Country     string  `json:"country"`
    State       string  `json:"state,omitempty"`
    County      string  `json:"county,omitempty"`
    City        string  `json:"city,omitempty"`
    PostalCodes []string `json:"postal_codes,omitempty" gorm:"serializer:json"`
    Coordinates []Coordinate `json:"coordinates,omitempty" gorm:"serializer:json"`
}

type Coordinate struct {
    Latitude  float64 `json:"latitude"`
    Longitude float64 `json:"longitude"`
}
```

#### Tax Calculation
```go
type TaxCalculation struct {
    ID                    string                 `json:"id" gorm:"primaryKey"`
    TransactionID         string                 `json:"transaction_id" gorm:"index"`
    CustomerID            string                 `json:"customer_id" gorm:"index"`
    
    // Calculation Input
    TaxableAmount         float64                `json:"taxable_amount"`
    ShippingAmount        float64                `json:"shipping_amount"`
    Currency              string                 `json:"currency" gorm:"default:'USD'"`
    
    // Location Context
    BillingAddress        Address                `json:"billing_address" gorm:"embedded;embeddedPrefix:billing_"`
    ShippingAddress       Address                `json:"shipping_address" gorm:"embedded;embeddedPrefix:shipping_"`
    OriginAddress         Address                `json:"origin_address" gorm:"embedded;embeddedPrefix:origin_"`
    
    // Products
    LineItems             []TaxLineItem          `json:"line_items" gorm:"serializer:json"`
    
    // Applied Tax Rules
    AppliedRules          []AppliedTaxRule       `json:"applied_rules" gorm:"serializer:json"`
    
    // Calculation Results
    TotalTaxAmount        float64                `json:"total_tax_amount"`
    TotalAmount           float64                `json:"total_amount"`
    EffectiveTaxRate      float64                `json:"effective_tax_rate"`
    
    // Exemptions Applied
    ExemptionsApplied     []AppliedExemption     `json:"exemptions_applied" gorm:"serializer:json"`
    ExemptAmount          float64                `json:"exempt_amount"`
    
    // Compliance
    ComplianceInfo        ComplianceInfo         `json:"compliance_info" gorm:"embedded"`
    
    // Status
    Status                CalculationStatus      `json:"status" gorm:"default:'COMPLETED'"`
    
    // Metadata
    CalculatedAt          time.Time              `json:"calculated_at"`
    CreatedAt             time.Time              `json:"created_at"`
    Version               string                 `json:"version" gorm:"default:'1.0'"`
}

type Address struct {
    Street1    string `json:"street1"`
    Street2    string `json:"street2,omitempty"`
    City       string `json:"city"`
    State      string `json:"state"`
    PostalCode string `json:"postal_code"`
    Country    string `json:"country"`
}

type TaxLineItem struct {
    ProductID       string  `json:"product_id"`
    SKU             string  `json:"sku"`
    Description     string  `json:"description"`
    Category        string  `json:"category"`
    UnitPrice       float64 `json:"unit_price"`
    Quantity        int     `json:"quantity"`
    TotalAmount     float64 `json:"total_amount"`
    TaxableAmount   float64 `json:"taxable_amount"`
    TaxAmount       float64 `json:"tax_amount"`
    ExemptAmount    float64 `json:"exempt_amount"`
    TaxRate         float64 `json:"tax_rate"`
}

type AppliedTaxRule struct {
    RuleID          string  `json:"rule_id"`
    RuleName        string  `json:"rule_name"`
    TaxType         string  `json:"tax_type"`
    JurisdictionID  string  `json:"jurisdiction_id"`
    JurisdictionName string `json:"jurisdiction_name"`
    TaxRate         float64 `json:"tax_rate"`
    TaxableAmount   float64 `json:"taxable_amount"`
    TaxAmount       float64 `json:"tax_amount"`
    Priority        int     `json:"priority"`
}

type AppliedExemption struct {
    ExemptionID     string  `json:"exemption_id"`
    ExemptionType   string  `json:"exemption_type"`
    ExemptionReason string  `json:"exemption_reason"`
    ExemptAmount    float64 `json:"exempt_amount"`
    CertificateNumber string `json:"certificate_number,omitempty"`
}

type ComplianceInfo struct {
    TaxIDNumber       string `json:"tax_id_number,omitempty"`
    VATNumber         string `json:"vat_number,omitempty"`
    ComplianceLevel   string `json:"compliance_level"`
    AuditTrailID      string `json:"audit_trail_id"`
    ReportingRequired bool   `json:"reporting_required"`
    DocumentationRequired bool `json:"documentation_required"`
}

type CalculationStatus string
const (
    CalculationStatusPending   CalculationStatus = "PENDING"
    CalculationStatusCompleted CalculationStatus = "COMPLETED"
    CalculationStatusFailed    CalculationStatus = "FAILED"
    CalculationStatusReviewed  CalculationStatus = "REVIEWED"
)
```

## REST API Endpoints

### Tax Calculation

#### Calculate Tax
```
POST /api/v1/tax/calculate
Content-Type: application/json

{
  "transaction_id": "txn-123",
  "customer_id": "customer-456",
  "billing_address": {
    "street1": "123 Main St",
    "city": "San Francisco",
    "state": "CA",
    "postal_code": "94105",
    "country": "US"
  },
  "shipping_address": {
    "street1": "456 Oak Ave",
    "city": "Los Angeles",
    "state": "CA",
    "postal_code": "90210",
    "country": "US"
  },
  "line_items": [
    {
      "product_id": "product-789",
      "sku": "LAPTOP-001",
      "description": "Gaming Laptop",
      "category": "electronics",
      "unit_price": 1200.00,
      "quantity": 1,
      "total_amount": 1200.00
    }
  ],
  "shipping_amount": 25.00,
  "exemptions": [
    {
      "exemption_type": "RESALE",
      "certificate_number": "RESALE-123456"
    }
  ]
}

Response: 200 OK
{
  "calculation_id": "calc-789",
  "total_tax_amount": 108.75,
  "total_amount": 1333.75,
  "effective_tax_rate": 0.087,
  "applied_rules": [
    {
      "rule_id": "rule-ca-state",
      "rule_name": "California State Sales Tax",
      "tax_type": "SALES_TAX",
      "jurisdiction_name": "California",
      "tax_rate": 0.075,
      "taxable_amount": 1200.00,
      "tax_amount": 90.00
    },
    {
      "rule_id": "rule-sf-local",
      "rule_name": "San Francisco Local Tax",
      "tax_type": "SALES_TAX", 
      "jurisdiction_name": "San Francisco",
      "tax_rate": 0.0125,
      "taxable_amount": 1200.00,
      "tax_amount": 15.00
    }
  ],
  "line_items": [
    {
      "product_id": "product-789",
      "taxable_amount": 1200.00,
      "tax_amount": 105.00,
      "tax_rate": 0.0875
    }
  ],
  "compliance_info": {
    "compliance_level": "STANDARD",
    "reporting_required": true,
    "audit_trail_id": "audit-456"
  }
}
```

#### Get Tax Jurisdictions
```
GET /api/v1/tax/jurisdictions?address=123+Main+St,San+Francisco,CA,94105

Response: 200 OK
{
  "jurisdictions": [
    {
      "id": "jurisdiction-us",
      "name": "United States",
      "code": "US",
      "type": "COUNTRY",
      "level": 0,
      "default_tax_rate": 0.0
    },
    {
      "id": "jurisdiction-ca",
      "name": "California",
      "code": "CA",
      "type": "STATE",
      "level": 1,
      "parent_id": "jurisdiction-us",
      "default_tax_rate": 0.075
    },
    {
      "id": "jurisdiction-sf",
      "name": "San Francisco",
      "code": "SF",
      "type": "CITY",
      "level": 2,
      "parent_id": "jurisdiction-ca",
      "default_tax_rate": 0.0125
    }
  ],
  "total_jurisdictions": 3,
  "combined_tax_rate": 0.0875
}
```

### Tax Rules Management

#### Create Tax Rule
```
POST /api/v1/tax/rules
Content-Type: application/json

{
  "name": "Electronics Luxury Tax",
  "description": "Additional tax on high-value electronics",
  "tax_type": "LUXURY_TAX",
  "jurisdiction_id": "jurisdiction-ca",
  "tax_rate": 0.05,
  "minimum_amount": 1000.00,
  "applicable_categories": ["electronics", "computers"],
  "effective_date": "2024-01-01T00:00:00Z",
  "priority": 200
}

Response: 201 Created
{
  "id": "rule-luxury-electronics",
  "name": "Electronics Luxury Tax",
  "status": "ACTIVE",
  "created_at": "2024-01-15T10:00:00Z"
}
```

### Compliance and Reporting

#### Generate Tax Report
```
GET /api/v1/tax/reports?start_date=2024-01-01&end_date=2024-01-31&jurisdiction=CA

Response: 200 OK
{
  "report_id": "report-202401-ca",
  "period": {
    "start_date": "2024-01-01T00:00:00Z",
    "end_date": "2024-01-31T23:59:59Z"
  },
  "jurisdiction": {
    "id": "jurisdiction-ca",
    "name": "California"
  },
  "summary": {
    "total_transactions": 15420,
    "total_taxable_amount": 2456780.50,
    "total_tax_collected": 215216.94,
    "total_exempt_amount": 125430.25,
    "effective_tax_rate": 0.0876
  },
  "tax_breakdown": [
    {
      "tax_type": "SALES_TAX",
      "tax_rate": 0.075,
      "taxable_amount": 2456780.50,
      "tax_amount": 184258.54
    },
    {
      "tax_type": "LOCAL_TAX",
      "tax_rate": 0.0125,
      "taxable_amount": 2456780.50,
      "tax_amount": 30709.76
    }
  ],
  "exemptions_summary": [
    {
      "exemption_type": "RESALE",
      "exempt_amount": 98765.43,
      "transaction_count": 156
    }
  ]
}
```

## Implementation Tasks

### Phase 1: Core Tax Engine (3-4 days)
1. **Domain Model and Database**
   - Implement tax rule, jurisdiction, and calculation entities
   - Create hierarchical jurisdiction structure
   - Add repository implementations with GORM
   - Implement tax rate and rule management

2. **Tax Calculation Engine**
   - Create tax calculation algorithms
   - Implement multi-jurisdictional tax logic
   - Add tax rate aggregation and compounding
   - Create exemption processing logic

### Phase 2: Jurisdiction Management (2-3 days)
1. **Jurisdiction Hierarchy**
   - Implement jurisdiction tree structure
   - Create address-to-jurisdiction mapping
   - Add geographic boundary checking
   - Implement jurisdiction inheritance rules

2. **Address Resolution**
   - Create address validation and normalization
   - Implement geocoding for address lookup
   - Add postal code to jurisdiction mapping
   - Create address-based tax determination

### Phase 3: Tax Rules Engine (2-3 days)
1. **Rule Management**
   - Implement tax rule CRUD operations
   - Create rule priority and conflict resolution
   - Add rule effective date management
   - Implement rule versioning and history

2. **Product Categorization**
   - Create product tax category management
   - Implement category-based tax rules
   - Add product-specific tax overrides
   - Create tax category inheritance

### Phase 4: Compliance and Reporting (2-3 days)
1. **Compliance Framework**
   - Implement tax compliance checking
   - Create audit trail generation
   - Add compliance documentation
   - Implement regulatory reporting

2. **Tax Reporting**
   - Create tax calculation reports
   - Implement jurisdiction-specific reports
   - Add compliance summary reports
   - Create export functionality for tax filings

### Phase 5: Integration and Optimization (2-3 days)
1. **External Integration**
   - Integrate with tax rate providers
   - Add address validation services
   - Create compliance API integrations
   - Implement real-time tax rate updates

2. **Performance Optimization**
   - Implement tax calculation caching
   - Add jurisdiction lookup optimization
   - Create tax rate caching strategies
   - Implement bulk calculation processing

## Estimated Development Time: 11-16 days
