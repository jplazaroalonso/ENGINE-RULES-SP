package campaign

import (
	"time"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/campaigns-management-service/internal/domain/shared"
)

// mustNewMoney creates a Money value and panics if there's an error
// This is used for known valid values where we don't expect errors
func mustNewMoney(amount float64, currency string) shared.Money {
	money, err := shared.NewMoney(amount, currency)
	if err != nil {
		panic(err)
	}
	return money
}

// CampaignMetrics represents campaign performance metrics
type CampaignMetrics struct {
	Impressions       int64        `json:"impressions"`
	Clicks            int64        `json:"clicks"`
	Conversions       int64        `json:"conversions"`
	Revenue           shared.Money `json:"revenue"`
	Cost              shared.Money `json:"cost"`
	CTR               float64      `json:"ctr"`               // Click-through rate
	ConversionRate    float64      `json:"conversionRate"`    // Conversion rate
	CostPerClick      shared.Money `json:"costPerClick"`      // CPC
	CostPerConversion shared.Money `json:"costPerConversion"` // CPA
	ROAS              float64      `json:"roas"`              // Return on ad spend
	ROI               float64      `json:"roi"`               // Return on investment
	LastUpdated       time.Time    `json:"lastUpdated"`
}

// CampaignEventType represents the type of campaign event
type CampaignEventType int

const (
	CampaignEventTypeImpression CampaignEventType = iota
	CampaignEventTypeClick
	CampaignEventTypeConversion
	CampaignEventTypeBounce
	CampaignEventTypeUnsubscribe
)

func (et CampaignEventType) String() string {
	switch et {
	case CampaignEventTypeImpression:
		return "IMPRESSION"
	case CampaignEventTypeClick:
		return "CLICK"
	case CampaignEventTypeConversion:
		return "CONVERSION"
	case CampaignEventTypeBounce:
		return "BOUNCE"
	case CampaignEventTypeUnsubscribe:
		return "UNSUBSCRIBE"
	default:
		return "UNKNOWN"
	}
}

func ParseCampaignEventType(eventType string) (CampaignEventType, error) {
	switch eventType {
	case "IMPRESSION":
		return CampaignEventTypeImpression, nil
	case "CLICK":
		return CampaignEventTypeClick, nil
	case "CONVERSION":
		return CampaignEventTypeConversion, nil
	case "BOUNCE":
		return CampaignEventTypeBounce, nil
	case "UNSUBSCRIBE":
		return CampaignEventTypeUnsubscribe, nil
	default:
		return CampaignEventTypeImpression, shared.NewValidationError("invalid campaign event type", nil)
	}
}

// CampaignEvent represents a campaign event
type CampaignEvent struct {
	ID         string                 `json:"id"`
	CampaignID CampaignID             `json:"campaignId"`
	EventType  CampaignEventType      `json:"eventType"`
	CustomerID *shared.CustomerID     `json:"customerId,omitempty"`
	EventData  map[string]interface{} `json:"eventData"`
	OccurredAt time.Time              `json:"occurredAt"`
	Revenue    *shared.Money          `json:"revenue,omitempty"`
	Cost       *shared.Money          `json:"cost,omitempty"`
}

// NewCampaignMetrics creates new campaign metrics
func NewCampaignMetrics() CampaignMetrics {
	return CampaignMetrics{
		Impressions:       0,
		Clicks:            0,
		Conversions:       0,
		Revenue:           mustNewMoney(0, "EUR"),
		Cost:              mustNewMoney(0, "EUR"),
		CTR:               0.0,
		ConversionRate:    0.0,
		CostPerClick:      mustNewMoney(0, "EUR"),
		CostPerConversion: mustNewMoney(0, "EUR"),
		ROAS:              0.0,
		ROI:               0.0,
		LastUpdated:       time.Now(),
	}
}

// NewCampaignEvent creates a new campaign event
func NewCampaignEvent(
	campaignID CampaignID,
	eventType CampaignEventType,
	customerID *shared.CustomerID,
	eventData map[string]interface{},
) CampaignEvent {
	return CampaignEvent{
		ID:         NewCampaignID().String(), // Generate unique ID
		CampaignID: campaignID,
		EventType:  eventType,
		CustomerID: customerID,
		EventData:  eventData,
		OccurredAt: time.Now(),
	}
}

// UpdateFromEvent updates metrics based on a campaign event
func (cm *CampaignMetrics) UpdateFromEvent(event CampaignEvent) {
	switch event.EventType {
	case CampaignEventTypeImpression:
		cm.Impressions++
	case CampaignEventTypeClick:
		cm.Clicks++
	case CampaignEventTypeConversion:
		cm.Conversions++
		if event.Revenue != nil {
			cm.Revenue, _ = cm.Revenue.Add(*event.Revenue)
		}
	case CampaignEventTypeBounce:
		// Bounce events don't directly affect metrics but can be used for analysis
	case CampaignEventTypeUnsubscribe:
		// Unsubscribe events don't directly affect metrics but can be used for analysis
	}

	// Update cost if provided
	if event.Cost != nil {
		cm.Cost, _ = cm.Cost.Add(*event.Cost)
	}

	// Recalculate derived metrics
	cm.calculateDerivedMetrics()
	cm.LastUpdated = time.Now()
}

// calculateDerivedMetrics calculates CTR, conversion rate, CPC, CPA, ROAS, and ROI
func (cm *CampaignMetrics) calculateDerivedMetrics() {
	// Calculate CTR (Click-through rate)
	if cm.Impressions > 0 {
		cm.CTR = float64(cm.Clicks) / float64(cm.Impressions) * 100
	} else {
		cm.CTR = 0.0
	}

	// Calculate Conversion Rate
	if cm.Clicks > 0 {
		cm.ConversionRate = float64(cm.Conversions) / float64(cm.Clicks) * 100
	} else {
		cm.ConversionRate = 0.0
	}

	// Calculate CPC (Cost Per Click)
	if cm.Clicks > 0 && !cm.Cost.IsZero() {
		cm.CostPerClick, _ = cm.Cost.Divide(float64(cm.Clicks))
	} else {
		cm.CostPerClick = mustNewMoney(0, cm.Cost.Currency)
	}

	// Calculate CPA (Cost Per Acquisition/Conversion)
	if cm.Conversions > 0 && !cm.Cost.IsZero() {
		cm.CostPerConversion, _ = cm.Cost.Divide(float64(cm.Conversions))
	} else {
		cm.CostPerConversion = mustNewMoney(0, cm.Cost.Currency)
	}

	// Calculate ROAS (Return on Ad Spend)
	if !cm.Cost.IsZero() {
		cm.ROAS = cm.Revenue.Amount / cm.Cost.Amount
	} else {
		cm.ROAS = 0.0
	}

	// Calculate ROI (Return on Investment)
	if !cm.Cost.IsZero() {
		profit := cm.Revenue.Amount - cm.Cost.Amount
		cm.ROI = (profit / cm.Cost.Amount) * 100
	} else {
		cm.ROI = 0.0
	}
}

// AddMetrics adds metrics from another campaign metrics instance
func (cm *CampaignMetrics) AddMetrics(other CampaignMetrics) error {
	// Validate currency compatibility
	if cm.Revenue.Currency != other.Revenue.Currency {
		return shared.NewValidationError("cannot add metrics with different revenue currencies", nil)
	}

	if cm.Cost.Currency != other.Cost.Currency {
		return shared.NewValidationError("cannot add metrics with different cost currencies", nil)
	}

	// Add basic metrics
	cm.Impressions += other.Impressions
	cm.Clicks += other.Clicks
	cm.Conversions += other.Conversions

	// Add monetary values
	cm.Revenue, _ = cm.Revenue.Add(other.Revenue)
	cm.Cost, _ = cm.Cost.Add(other.Cost)

	// Recalculate derived metrics
	cm.calculateDerivedMetrics()
	cm.LastUpdated = time.Now()

	return nil
}

// Reset resets all metrics to zero
func (cm *CampaignMetrics) Reset() {
	cm.Impressions = 0
	cm.Clicks = 0
	cm.Conversions = 0
	cm.Revenue = mustNewMoney(0, cm.Revenue.Currency)
	cm.Cost = mustNewMoney(0, cm.Cost.Currency)
	cm.CTR = 0.0
	cm.ConversionRate = 0.0
	cm.CostPerClick = mustNewMoney(0, cm.Cost.Currency)
	cm.CostPerConversion = mustNewMoney(0, cm.Cost.Currency)
	cm.ROAS = 0.0
	cm.ROI = 0.0
	cm.LastUpdated = time.Now()
}

// GetPerformanceScore calculates an overall performance score (0-100)
func (cm CampaignMetrics) GetPerformanceScore() float64 {
	score := 0.0

	// CTR score (0-30 points)
	if cm.CTR >= 5.0 {
		score += 30.0
	} else if cm.CTR >= 2.0 {
		score += 20.0
	} else if cm.CTR >= 1.0 {
		score += 10.0
	}

	// Conversion rate score (0-30 points)
	if cm.ConversionRate >= 10.0 {
		score += 30.0
	} else if cm.ConversionRate >= 5.0 {
		score += 20.0
	} else if cm.ConversionRate >= 2.0 {
		score += 10.0
	}

	// ROI score (0-40 points)
	if cm.ROI >= 200.0 {
		score += 40.0
	} else if cm.ROI >= 100.0 {
		score += 30.0
	} else if cm.ROI >= 50.0 {
		score += 20.0
	} else if cm.ROI >= 0.0 {
		score += 10.0
	}

	return score
}

// GetPerformanceGrade returns a letter grade based on performance score
func (cm CampaignMetrics) GetPerformanceGrade() string {
	score := cm.GetPerformanceScore()

	switch {
	case score >= 90:
		return "A+"
	case score >= 80:
		return "A"
	case score >= 70:
		return "B+"
	case score >= 60:
		return "B"
	case score >= 50:
		return "C+"
	case score >= 40:
		return "C"
	case score >= 30:
		return "D"
	default:
		return "F"
	}
}

// IsPerformingWell returns true if the campaign is performing well
func (cm CampaignMetrics) IsPerformingWell() bool {
	return cm.GetPerformanceScore() >= 60.0
}

// NeedsAttention returns true if the campaign needs attention
func (cm CampaignMetrics) NeedsAttention() bool {
	score := cm.GetPerformanceScore()
	return score < 40.0
}

// GetRecommendations returns performance improvement recommendations
func (cm CampaignMetrics) GetRecommendations() []string {
	var recommendations []string

	// CTR recommendations
	if cm.CTR < 1.0 {
		recommendations = append(recommendations, "Improve ad creative and targeting to increase click-through rate")
	}

	// Conversion rate recommendations
	if cm.ConversionRate < 2.0 {
		recommendations = append(recommendations, "Optimize landing page and user experience to improve conversion rate")
	}

	// ROI recommendations
	if cm.ROI < 50.0 {
		recommendations = append(recommendations, "Review campaign costs and optimize bidding strategy")
	}

	// Volume recommendations
	if cm.Impressions < 1000 {
		recommendations = append(recommendations, "Increase campaign reach and budget to generate more impressions")
	}

	// Cost recommendations
	if !cm.CostPerClick.IsZero() && cm.CostPerClick.Amount > 2.0 {
		recommendations = append(recommendations, "Optimize targeting and ad quality to reduce cost per click")
	}

	return recommendations
}
