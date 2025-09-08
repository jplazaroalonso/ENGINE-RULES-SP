package customer

import (
	"time"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/customer-management-service/internal/domain/shared"
)

// CustomerID represents a customer identifier
type CustomerID = shared.CustomerID

// Money represents a monetary value
type Money = shared.Money

// CustomerMetrics represents customer metrics
type CustomerMetrics struct {
	CustomerID            CustomerID `json:"customerId"`
	TotalPurchases        int        `json:"totalPurchases"`
	TotalSpent            Money      `json:"totalSpent"`
	AverageOrderValue     Money      `json:"averageOrderValue"`
	LastPurchaseDate      *time.Time `json:"lastPurchaseDate,omitempty"`
	DaysSinceLastPurchase *int       `json:"daysSinceLastPurchase,omitempty"`
	PurchaseFrequency     float64    `json:"purchaseFrequency"` // purchases per month
	CustomerLifetime      int        `json:"customerLifetime"`  // days since first purchase
	ChurnRisk             float64    `json:"churnRisk"`         // 0-1 scale
	EngagementScore       float64    `json:"engagementScore"`   // 0-100 scale
	LastCalculated        time.Time  `json:"lastCalculated"`
}

// NewCustomerMetrics creates new customer metrics
func NewCustomerMetrics(customerID CustomerID, totalPurchases int, totalSpent Money, averageOrderValue Money, lastPurchaseDate *time.Time, daysSinceLastPurchase *int, purchaseFrequency float64, customerLifetime int, churnRisk float64, engagementScore float64) CustomerMetrics {
	return CustomerMetrics{
		CustomerID:            customerID,
		TotalPurchases:        totalPurchases,
		TotalSpent:            totalSpent,
		AverageOrderValue:     averageOrderValue,
		LastPurchaseDate:      lastPurchaseDate,
		DaysSinceLastPurchase: daysSinceLastPurchase,
		PurchaseFrequency:     purchaseFrequency,
		CustomerLifetime:      customerLifetime,
		ChurnRisk:             churnRisk,
		EngagementScore:       engagementScore,
		LastCalculated:        time.Now(),
	}
}

// CustomerInsights represents customer insights
type CustomerInsights struct {
	CustomerID       CustomerID        `json:"customerId"`
	BehaviorPatterns []BehaviorPattern `json:"behaviorPatterns"`
	Recommendations  []Recommendation  `json:"recommendations"`
	RiskFactors      []RiskFactor      `json:"riskFactors"`
	Opportunities    []Opportunity     `json:"opportunities"`
	LastCalculated   time.Time         `json:"lastCalculated"`
}

// BehaviorPattern represents a behavior pattern
type BehaviorPattern struct {
	Type        string                 `json:"type"`
	Description string                 `json:"description"`
	Confidence  float64                `json:"confidence"` // 0-1 scale
	Data        map[string]interface{} `json:"data"`
}

// Recommendation represents a recommendation
type Recommendation struct {
	Type        string                 `json:"type"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Priority    string                 `json:"priority"` // HIGH, MEDIUM, LOW
	Data        map[string]interface{} `json:"data"`
}

// RiskFactor represents a risk factor
type RiskFactor struct {
	Type        string                 `json:"type"`
	Description string                 `json:"description"`
	Severity    string                 `json:"severity"` // HIGH, MEDIUM, LOW
	Data        map[string]interface{} `json:"data"`
}

// Opportunity represents an opportunity
type Opportunity struct {
	Type        string                 `json:"type"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Potential   string                 `json:"potential"` // HIGH, MEDIUM, LOW
	Data        map[string]interface{} `json:"data"`
}

// NewCustomerInsights creates new customer insights
func NewCustomerInsights(customerID CustomerID, behaviorPatterns []BehaviorPattern, recommendations []Recommendation, riskFactors []RiskFactor, opportunities []Opportunity) CustomerInsights {
	return CustomerInsights{
		CustomerID:       customerID,
		BehaviorPatterns: behaviorPatterns,
		Recommendations:  recommendations,
		RiskFactors:      riskFactors,
		Opportunities:    opportunities,
		LastCalculated:   time.Now(),
	}
}

// CustomerEvent represents a customer event
type CustomerEvent struct {
	ID         string                 `json:"id"`
	CustomerID CustomerID             `json:"customerId"`
	EventType  string                 `json:"eventType"`
	EventData  map[string]interface{} `json:"eventData"`
	OccurredAt time.Time              `json:"occurredAt"`
	SessionID  *string                `json:"sessionId,omitempty"`
	DeviceInfo map[string]interface{} `json:"deviceInfo,omitempty"`
}

// NewCustomerEvent creates a new customer event
func NewCustomerEvent(customerID CustomerID, eventType string, eventData map[string]interface{}, sessionID *string, deviceInfo map[string]interface{}) CustomerEvent {
	return CustomerEvent{
		ID:         shared.NewEventID().String(),
		CustomerID: customerID,
		EventType:  eventType,
		EventData:  eventData,
		OccurredAt: time.Now(),
		SessionID:  sessionID,
		DeviceInfo: deviceInfo,
	}
}

// ReportCriteria represents criteria for generating reports
type ReportCriteria struct {
	CustomerIDs []CustomerID           `json:"customerIds,omitempty"`
	SegmentIDs  []SegmentID            `json:"segmentIds,omitempty"`
	DateRange   *DateRange             `json:"dateRange,omitempty"`
	Metrics     []string               `json:"metrics,omitempty"`
	GroupBy     []string               `json:"groupBy,omitempty"`
	Filters     map[string]interface{} `json:"filters,omitempty"`
}

// DateRange represents a date range
type DateRange struct {
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
}

// NewDateRange creates a new date range
func NewDateRange(startDate, endDate time.Time) (DateRange, error) {
	if startDate.After(endDate) {
		return DateRange{}, shared.NewValidationError("start date cannot be after end date", nil)
	}

	return DateRange{
		StartDate: startDate,
		EndDate:   endDate,
	}, nil
}

// CustomerReport represents a customer report
type CustomerReport struct {
	ID          string                 `json:"id"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Criteria    ReportCriteria         `json:"criteria"`
	Data        map[string]interface{} `json:"data"`
	GeneratedAt time.Time              `json:"generatedAt"`
	GeneratedBy string                 `json:"generatedBy"`
}

// NewCustomerReport creates a new customer report
func NewCustomerReport(title, description string, criteria ReportCriteria, data map[string]interface{}, generatedBy string) CustomerReport {
	return CustomerReport{
		ID:          shared.NewEventID().String(),
		Title:       title,
		Description: description,
		Criteria:    criteria,
		Data:        data,
		GeneratedAt: time.Now(),
		GeneratedBy: generatedBy,
	}
}

// CustomerDataExport represents a customer data export (GDPR compliance)
type CustomerDataExport struct {
	CustomerID  CustomerID             `json:"customerId"`
	ExportID    string                 `json:"exportId"`
	Data        map[string]interface{} `json:"data"`
	RequestedAt time.Time              `json:"requestedAt"`
	ExportedAt  time.Time              `json:"exportedAt"`
	ExpiresAt   time.Time              `json:"expiresAt"`
	DownloadURL *string                `json:"downloadUrl,omitempty"`
}

// NewCustomerDataExport creates a new customer data export
func NewCustomerDataExport(customerID CustomerID, data map[string]interface{}, downloadURL *string) CustomerDataExport {
	now := time.Now()
	expiresAt := now.Add(30 * 24 * time.Hour) // 30 days

	return CustomerDataExport{
		CustomerID:  customerID,
		ExportID:    shared.NewEventID().String(),
		Data:        data,
		RequestedAt: now,
		ExportedAt:  now,
		ExpiresAt:   expiresAt,
		DownloadURL: downloadURL,
	}
}

// PrivacyConsent represents privacy consent information
type PrivacyConsent struct {
	CustomerID               CustomerID `json:"customerId"`
	MarketingConsent         bool       `json:"marketingConsent"`
	DataProcessingConsent    bool       `json:"dataProcessingConsent"`
	AnalyticsConsent         bool       `json:"analyticsConsent"`
	PersonalizationConsent   bool       `json:"personalizationConsent"`
	ThirdPartySharingConsent bool       `json:"thirdPartySharingConsent"`
	UpdatedAt                time.Time  `json:"updatedAt"`
	UpdatedBy                string     `json:"updatedBy"`
}

// NewPrivacyConsent creates new privacy consent
func NewPrivacyConsent(customerID CustomerID, marketingConsent, dataProcessingConsent, analyticsConsent, personalizationConsent, thirdPartySharingConsent bool, updatedBy string) PrivacyConsent {
	return PrivacyConsent{
		CustomerID:               customerID,
		MarketingConsent:         marketingConsent,
		DataProcessingConsent:    dataProcessingConsent,
		AnalyticsConsent:         analyticsConsent,
		PersonalizationConsent:   personalizationConsent,
		ThirdPartySharingConsent: thirdPartySharingConsent,
		UpdatedAt:                time.Now(),
		UpdatedBy:                updatedBy,
	}
}
