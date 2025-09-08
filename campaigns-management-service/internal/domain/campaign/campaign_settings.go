package campaign

import (
	"fmt"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/campaigns-management-service/internal/domain/shared"
)

// CampaignSettings represents campaign configuration and settings
type CampaignSettings struct {
	TargetAudience  []string              `json:"targetAudience"`
	Channels        []Channel             `json:"channels"`
	Frequency       Frequency             `json:"frequency"`
	MaxImpressions  *int                  `json:"maxImpressions,omitempty"`
	BudgetLimit     *shared.Money         `json:"budgetLimit,omitempty"`
	ABTestSettings  *ABTestSettings       `json:"abTestSettings,omitempty"`
	SchedulingRules []SchedulingRule      `json:"schedulingRules"`
	Personalization PersonalizationConfig `json:"personalization"`
}

// Channel represents the communication channel
type Channel int

const (
	ChannelEmail Channel = iota
	ChannelSMS
	ChannelPush
	ChannelWeb
	ChannelSocial
	ChannelDisplay
)

func (c Channel) String() string {
	switch c {
	case ChannelEmail:
		return "EMAIL"
	case ChannelSMS:
		return "SMS"
	case ChannelPush:
		return "PUSH"
	case ChannelWeb:
		return "WEB"
	case ChannelSocial:
		return "SOCIAL"
	case ChannelDisplay:
		return "DISPLAY"
	default:
		return "UNKNOWN"
	}
}

func ParseChannel(channel string) (Channel, error) {
	switch channel {
	case "EMAIL":
		return ChannelEmail, nil
	case "SMS":
		return ChannelSMS, nil
	case "PUSH":
		return ChannelPush, nil
	case "WEB":
		return ChannelWeb, nil
	case "SOCIAL":
		return ChannelSocial, nil
	case "DISPLAY":
		return ChannelDisplay, nil
	default:
		return ChannelEmail, shared.NewValidationError(fmt.Sprintf("invalid channel: %s", channel), nil)
	}
}

// Frequency represents the campaign frequency
type Frequency int

const (
	FrequencyOnce Frequency = iota
	FrequencyDaily
	FrequencyWeekly
	FrequencyMonthly
)

func (f Frequency) String() string {
	switch f {
	case FrequencyOnce:
		return "ONCE"
	case FrequencyDaily:
		return "DAILY"
	case FrequencyWeekly:
		return "WEEKLY"
	case FrequencyMonthly:
		return "MONTHLY"
	default:
		return "UNKNOWN"
	}
}

func ParseFrequency(frequency string) (Frequency, error) {
	switch frequency {
	case "ONCE":
		return FrequencyOnce, nil
	case "DAILY":
		return FrequencyDaily, nil
	case "WEEKLY":
		return FrequencyWeekly, nil
	case "MONTHLY":
		return FrequencyMonthly, nil
	default:
		return FrequencyOnce, shared.NewValidationError(fmt.Sprintf("invalid frequency: %s", frequency), nil)
	}
}

// ABTestSettings represents A/B testing configuration
type ABTestSettings struct {
	Enabled       bool      `json:"enabled"`
	Variants      []Variant `json:"variants"`
	TrafficSplit  float64   `json:"trafficSplit"` // Percentage of traffic for variant A
	SuccessMetric string    `json:"successMetric"`
	Duration      int       `json:"duration"` // Duration in days
}

// Variant represents an A/B test variant
type Variant struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Settings    map[string]interface{} `json:"settings"`
	Weight      float64                `json:"weight"` // Traffic weight (0.0 to 1.0)
}

// SchedulingRule represents campaign scheduling rules
type SchedulingRule struct {
	ID          string                `json:"id"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Conditions  []SchedulingCondition `json:"conditions"`
	Actions     []SchedulingAction    `json:"actions"`
	IsActive    bool                  `json:"isActive"`
}

// SchedulingCondition represents a condition for scheduling
type SchedulingCondition struct {
	Type     string                 `json:"type"`     // "time", "date", "event", "metric"
	Operator string                 `json:"operator"` // "equals", "greater_than", "less_than", "contains"
	Value    interface{}            `json:"value"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// SchedulingAction represents an action to take when conditions are met
type SchedulingAction struct {
	Type       string                 `json:"type"` // "activate", "pause", "stop", "update_settings"
	Parameters map[string]interface{} `json:"parameters"`
}

// PersonalizationConfig represents personalization settings
type PersonalizationConfig struct {
	Enabled     bool     `json:"enabled"`
	Rules       []string `json:"rules"`       // Rule IDs for personalization
	Fallback    string   `json:"fallback"`    // Fallback content when personalization fails
	MaxVariants int      `json:"maxVariants"` // Maximum number of variants
}

// NewCampaignSettings creates new campaign settings with validation
func NewCampaignSettings(
	targetAudience []string,
	channels []Channel,
	frequency Frequency,
	maxImpressions *int,
	budgetLimit *shared.Money,
	abTestSettings *ABTestSettings,
	schedulingRules []SchedulingRule,
	personalization PersonalizationConfig,
) (CampaignSettings, error) {
	settings := CampaignSettings{
		TargetAudience:  targetAudience,
		Channels:        channels,
		Frequency:       frequency,
		MaxImpressions:  maxImpressions,
		BudgetLimit:     budgetLimit,
		ABTestSettings:  abTestSettings,
		SchedulingRules: schedulingRules,
		Personalization: personalization,
	}

	if err := settings.Validate(); err != nil {
		return CampaignSettings{}, err
	}

	return settings, nil
}

// Validate validates campaign settings
func (cs CampaignSettings) Validate() error {
	// Validate channels
	if len(cs.Channels) == 0 {
		return shared.NewValidationError("at least one channel must be specified", nil)
	}

	// Validate max impressions
	if cs.MaxImpressions != nil && *cs.MaxImpressions <= 0 {
		return shared.NewValidationError("max impressions must be positive", nil)
	}

	// Validate budget limit
	if cs.BudgetLimit != nil && !cs.BudgetLimit.IsPositive() {
		return shared.NewValidationError("budget limit must be positive", nil)
	}

	// Validate A/B test settings
	if cs.ABTestSettings != nil {
		if err := cs.ABTestSettings.Validate(); err != nil {
			return err
		}
	}

	// Validate scheduling rules
	for i, rule := range cs.SchedulingRules {
		if err := rule.Validate(); err != nil {
			return fmt.Errorf("invalid scheduling rule at index %d: %w", i, err)
		}
	}

	// Validate personalization config
	if err := cs.Personalization.Validate(); err != nil {
		return err
	}

	return nil
}

// Validate validates A/B test settings
func (ab ABTestSettings) Validate() error {
	if !ab.Enabled {
		return nil
	}

	// Validate variants
	if len(ab.Variants) < 2 {
		return shared.NewValidationError("A/B test must have at least 2 variants", nil)
	}

	if len(ab.Variants) > 10 {
		return shared.NewValidationError("A/B test cannot have more than 10 variants", nil)
	}

	// Validate traffic split
	if ab.TrafficSplit < 0.0 || ab.TrafficSplit > 1.0 {
		return shared.NewValidationError("traffic split must be between 0.0 and 1.0", nil)
	}

	// Validate success metric
	if ab.SuccessMetric == "" {
		return shared.NewValidationError("success metric is required for A/B test", nil)
	}

	// Validate duration
	if ab.Duration <= 0 {
		return shared.NewValidationError("A/B test duration must be positive", nil)
	}

	// Validate variants
	totalWeight := 0.0
	for i, variant := range ab.Variants {
		if err := variant.Validate(); err != nil {
			return fmt.Errorf("invalid variant at index %d: %w", i, err)
		}
		totalWeight += variant.Weight
	}

	// Check if weights sum to approximately 1.0 (allowing for floating point precision)
	if totalWeight < 0.99 || totalWeight > 1.01 {
		return shared.NewValidationError("variant weights must sum to 1.0", nil)
	}

	return nil
}

// Validate validates a variant
func (v Variant) Validate() error {
	if v.ID == "" {
		return shared.NewValidationError("variant ID is required", nil)
	}

	if v.Name == "" {
		return shared.NewValidationError("variant name is required", nil)
	}

	if v.Weight < 0.0 || v.Weight > 1.0 {
		return shared.NewValidationError("variant weight must be between 0.0 and 1.0", nil)
	}

	return nil
}

// Validate validates a scheduling rule
func (sr SchedulingRule) Validate() error {
	if sr.ID == "" {
		return shared.NewValidationError("scheduling rule ID is required", nil)
	}

	if sr.Name == "" {
		return shared.NewValidationError("scheduling rule name is required", nil)
	}

	// Validate conditions
	if len(sr.Conditions) == 0 {
		return shared.NewValidationError("scheduling rule must have at least one condition", nil)
	}

	for i, condition := range sr.Conditions {
		if err := condition.Validate(); err != nil {
			return fmt.Errorf("invalid condition at index %d: %w", i, err)
		}
	}

	// Validate actions
	if len(sr.Actions) == 0 {
		return shared.NewValidationError("scheduling rule must have at least one action", nil)
	}

	for i, action := range sr.Actions {
		if err := action.Validate(); err != nil {
			return fmt.Errorf("invalid action at index %d: %w", i, err)
		}
	}

	return nil
}

// Validate validates a scheduling condition
func (sc SchedulingCondition) Validate() error {
	if sc.Type == "" {
		return shared.NewValidationError("condition type is required", nil)
	}

	validTypes := []string{"time", "date", "event", "metric"}
	if !contains(validTypes, sc.Type) {
		return shared.NewValidationError(fmt.Sprintf("invalid condition type: %s", sc.Type), nil)
	}

	if sc.Operator == "" {
		return shared.NewValidationError("condition operator is required", nil)
	}

	validOperators := []string{"equals", "greater_than", "less_than", "contains", "not_equals"}
	if !contains(validOperators, sc.Operator) {
		return shared.NewValidationError(fmt.Sprintf("invalid condition operator: %s", sc.Operator), nil)
	}

	return nil
}

// Validate validates a scheduling action
func (sa SchedulingAction) Validate() error {
	if sa.Type == "" {
		return shared.NewValidationError("action type is required", nil)
	}

	validTypes := []string{"activate", "pause", "stop", "update_settings", "send_notification"}
	if !contains(validTypes, sa.Type) {
		return shared.NewValidationError(fmt.Sprintf("invalid action type: %s", sa.Type), nil)
	}

	return nil
}

// Validate validates personalization config
func (pc PersonalizationConfig) Validate() error {
	if !pc.Enabled {
		return nil
	}

	if pc.MaxVariants <= 0 {
		return shared.NewValidationError("max variants must be positive", nil)
	}

	if pc.MaxVariants > 100 {
		return shared.NewValidationError("max variants cannot exceed 100", nil)
	}

	return nil
}

// Helper function to check if slice contains value
func contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}
