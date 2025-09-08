package settings

import (
	"time"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/domain/shared"
)

// RolloutStrategy represents the rollout strategy for feature flags
type RolloutStrategy string

const (
	RolloutStrategyAll        RolloutStrategy = "ALL"
	RolloutStrategyPercentage RolloutStrategy = "PERCENTAGE"
	RolloutStrategyUserList   RolloutStrategy = "USER_LIST"
	RolloutStrategyRules      RolloutStrategy = "RULES"
)

// String returns the string representation of the rollout strategy
func (r RolloutStrategy) String() string {
	return string(r)
}

// ParseRolloutStrategy parses a string to RolloutStrategy
func ParseRolloutStrategy(strategy string) (RolloutStrategy, error) {
	switch strategy {
	case "ALL":
		return RolloutStrategyAll, nil
	case "PERCENTAGE":
		return RolloutStrategyPercentage, nil
	case "USER_LIST":
		return RolloutStrategyUserList, nil
	case "RULES":
		return RolloutStrategyRules, nil
	default:
		return "", shared.NewValidationError("invalid rollout strategy", nil)
	}
}

// TargetAudience represents the target audience for feature flags
type TargetAudience struct {
	UserIDs       []string               `json:"userIds,omitempty"`
	UserGroups    []string               `json:"userGroups,omitempty"`
	Organizations []string               `json:"organizations,omitempty"`
	Regions       []string               `json:"regions,omitempty"`
	Custom        map[string]interface{} `json:"custom,omitempty"`
}

// NewTargetAudience creates a new target audience
func NewTargetAudience(userIDs, userGroups, organizations, regions []string, custom map[string]interface{}) TargetAudience {
	return TargetAudience{
		UserIDs:       userIDs,
		UserGroups:    userGroups,
		Organizations: organizations,
		Regions:       regions,
		Custom:        custom,
	}
}

// Variant represents a feature flag variant
type Variant struct {
	Key         string      `json:"key"`
	Name        string      `json:"name"`
	Value       interface{} `json:"value"`
	Weight      int         `json:"weight"`
	Description string      `json:"description"`
}

// NewVariant creates a new variant
func NewVariant(key, name string, value interface{}, weight int, description string) (Variant, error) {
	if key == "" {
		return Variant{}, shared.NewValidationError("variant key is required", nil)
	}

	if name == "" {
		return Variant{}, shared.NewValidationError("variant name is required", nil)
	}

	if weight < 0 || weight > 100 {
		return Variant{}, shared.NewValidationError("variant weight must be between 0 and 100", nil)
	}

	return Variant{
		Key:         key,
		Name:        name,
		Value:       value,
		Weight:      weight,
		Description: description,
	}, nil
}

// TargetingRule represents a targeting rule for feature flags
type TargetingRule struct {
	ID         string               `json:"id"`
	Name       string               `json:"name"`
	Conditions []TargetingCondition `json:"conditions"`
	Variant    string               `json:"variant"`
	IsEnabled  bool                 `json:"isEnabled"`
}

// NewTargetingRule creates a new targeting rule
func NewTargetingRule(id, name string, conditions []TargetingCondition, variant string, isEnabled bool) (TargetingRule, error) {
	if id == "" {
		return TargetingRule{}, shared.NewValidationError("targeting rule ID is required", nil)
	}

	if name == "" {
		return TargetingRule{}, shared.NewValidationError("targeting rule name is required", nil)
	}

	if len(conditions) == 0 {
		return TargetingRule{}, shared.NewValidationError("targeting rule must have at least one condition", nil)
	}

	if variant == "" {
		return TargetingRule{}, shared.NewValidationError("targeting rule variant is required", nil)
	}

	return TargetingRule{
		ID:         id,
		Name:       name,
		Conditions: conditions,
		Variant:    variant,
		IsEnabled:  isEnabled,
	}, nil
}

// TargetingCondition represents a targeting condition
type TargetingCondition struct {
	Attribute string      `json:"attribute"`
	Operator  string      `json:"operator"`
	Value     interface{} `json:"value"`
}

// NewTargetingCondition creates a new targeting condition
func NewTargetingCondition(attribute, operator string, value interface{}) (TargetingCondition, error) {
	if attribute == "" {
		return TargetingCondition{}, shared.NewValidationError("targeting condition attribute is required", nil)
	}

	if operator == "" {
		return TargetingCondition{}, shared.NewValidationError("targeting condition operator is required", nil)
	}

	validOperators := []string{"equals", "not_equals", "contains", "not_contains", "greater_than", "less_than", "in", "not_in"}
	valid := false
	for _, op := range validOperators {
		if operator == op {
			valid = true
			break
		}
	}

	if !valid {
		return TargetingCondition{}, shared.NewValidationError("invalid targeting condition operator", nil)
	}

	return TargetingCondition{
		Attribute: attribute,
		Operator:  operator,
		Value:     value,
	}, nil
}

// FeatureFlag represents a feature flag aggregate
type FeatureFlag struct {
	id              FeatureFlagID
	name            string
	description     string
	key             string
	isEnabled       bool
	rolloutStrategy RolloutStrategy
	targetAudience  TargetAudience
	variants        []Variant
	rules           []TargetingRule
	environment     Environment
	organizationID  *OrganizationID
	service         *ServiceName
	createdBy       UserID
	createdAt       time.Time
	updatedAt       time.Time
	version         int
	events          []shared.DomainEvent
}

// NewFeatureFlag creates a new feature flag
func NewFeatureFlag(
	name string,
	description string,
	key string,
	isEnabled bool,
	rolloutStrategy RolloutStrategy,
	targetAudience TargetAudience,
	variants []Variant,
	rules []TargetingRule,
	environment Environment,
	organizationID *OrganizationID,
	service *ServiceName,
	createdBy UserID,
) (*FeatureFlag, error) {
	if name == "" {
		return nil, shared.NewValidationError("feature flag name is required", nil)
	}

	if description == "" {
		return nil, shared.NewValidationError("feature flag description is required", nil)
	}

	if key == "" {
		return nil, shared.NewValidationError("feature flag key is required", nil)
	}

	if createdBy.IsEmpty() {
		return nil, shared.NewValidationError("created by user ID is required", nil)
	}

	// Validate key format (should be unique and follow naming convention)
	if !isValidFeatureFlagKey(key) {
		return nil, shared.NewValidationError("invalid feature flag key format", nil)
	}

	// Validate variants
	if len(variants) == 0 {
		return nil, shared.NewValidationError("feature flag must have at least one variant", nil)
	}

	// Validate variant weights sum to 100 for percentage rollout
	if rolloutStrategy == RolloutStrategyPercentage {
		totalWeight := 0
		for _, variant := range variants {
			totalWeight += variant.Weight
		}
		if totalWeight != 100 {
			return nil, shared.NewValidationError("variant weights must sum to 100 for percentage rollout", nil)
		}
	}

	now := time.Now()

	featureFlag := &FeatureFlag{
		id:              shared.NewFeatureFlagID(),
		name:            name,
		description:     description,
		key:             key,
		isEnabled:       isEnabled,
		rolloutStrategy: rolloutStrategy,
		targetAudience:  targetAudience,
		variants:        variants,
		rules:           rules,
		environment:     environment,
		organizationID:  organizationID,
		service:         service,
		createdBy:       createdBy,
		createdAt:       now,
		updatedAt:       now,
		version:         1,
		events:          []shared.DomainEvent{},
	}

	// Add feature flag created event
	featureFlag.addEvent(NewFeatureFlagCreatedEvent(featureFlag))

	return featureFlag, nil
}

// GetID returns the feature flag ID
func (f *FeatureFlag) GetID() FeatureFlagID {
	return f.id
}

// GetName returns the feature flag name
func (f *FeatureFlag) GetName() string {
	return f.name
}

// GetDescription returns the feature flag description
func (f *FeatureFlag) GetDescription() string {
	return f.description
}

// GetKey returns the feature flag key
func (f *FeatureFlag) GetKey() string {
	return f.key
}

// IsEnabled returns whether the feature flag is enabled
func (f *FeatureFlag) IsEnabled() bool {
	return f.isEnabled
}

// GetRolloutStrategy returns the rollout strategy
func (f *FeatureFlag) GetRolloutStrategy() RolloutStrategy {
	return f.rolloutStrategy
}

// GetTargetAudience returns the target audience
func (f *FeatureFlag) GetTargetAudience() TargetAudience {
	return f.targetAudience
}

// GetVariants returns the variants
func (f *FeatureFlag) GetVariants() []Variant {
	return f.variants
}

// GetRules returns the targeting rules
func (f *FeatureFlag) GetRules() []TargetingRule {
	return f.rules
}

// GetEnvironment returns the environment
func (f *FeatureFlag) GetEnvironment() Environment {
	return f.environment
}

// GetOrganizationID returns the organization ID
func (f *FeatureFlag) GetOrganizationID() *OrganizationID {
	return f.organizationID
}

// GetService returns the service name
func (f *FeatureFlag) GetService() *ServiceName {
	return f.service
}

// GetCreatedBy returns the user who created the feature flag
func (f *FeatureFlag) GetCreatedBy() UserID {
	return f.createdBy
}

// GetCreatedAt returns the creation timestamp
func (f *FeatureFlag) GetCreatedAt() time.Time {
	return f.createdAt
}

// GetUpdatedAt returns the last update timestamp
func (f *FeatureFlag) GetUpdatedAt() time.Time {
	return f.updatedAt
}

// GetVersion returns the feature flag version
func (f *FeatureFlag) GetVersion() int {
	return f.version
}

// GetEvents returns the domain events
func (f *FeatureFlag) GetEvents() []shared.DomainEvent {
	return f.events
}

// ClearEvents clears the domain events
func (f *FeatureFlag) ClearEvents() {
	f.events = []shared.DomainEvent{}
}

// Enable enables the feature flag
func (f *FeatureFlag) Enable(updatedBy UserID) error {
	if updatedBy.IsEmpty() {
		return shared.NewValidationError("updated by user ID is required", nil)
	}

	if f.isEnabled {
		return shared.NewValidationError("feature flag is already enabled", nil)
	}

	f.isEnabled = true
	f.updatedAt = time.Now()
	f.version++

	f.addEvent(NewFeatureFlagEnabledEvent(f))

	return nil
}

// Disable disables the feature flag
func (f *FeatureFlag) Disable(updatedBy UserID) error {
	if updatedBy.IsEmpty() {
		return shared.NewValidationError("updated by user ID is required", nil)
	}

	if !f.isEnabled {
		return shared.NewValidationError("feature flag is already disabled", nil)
	}

	f.isEnabled = false
	f.updatedAt = time.Now()
	f.version++

	f.addEvent(NewFeatureFlagDisabledEvent(f))

	return nil
}

// UpdateRolloutStrategy updates the rollout strategy
func (f *FeatureFlag) UpdateRolloutStrategy(strategy RolloutStrategy, updatedBy UserID) error {
	if updatedBy.IsEmpty() {
		return shared.NewValidationError("updated by user ID is required", nil)
	}

	// Validate variant weights for percentage rollout
	if strategy == RolloutStrategyPercentage {
		totalWeight := 0
		for _, variant := range f.variants {
			totalWeight += variant.Weight
		}
		if totalWeight != 100 {
			return shared.NewValidationError("variant weights must sum to 100 for percentage rollout", nil)
		}
	}

	f.rolloutStrategy = strategy
	f.updatedAt = time.Now()
	f.version++

	f.addEvent(NewFeatureFlagUpdatedEvent(f))

	return nil
}

// UpdateTargetAudience updates the target audience
func (f *FeatureFlag) UpdateTargetAudience(audience TargetAudience, updatedBy UserID) error {
	if updatedBy.IsEmpty() {
		return shared.NewValidationError("updated by user ID is required", nil)
	}

	f.targetAudience = audience
	f.updatedAt = time.Now()
	f.version++

	f.addEvent(NewFeatureFlagUpdatedEvent(f))

	return nil
}

// AddVariant adds a new variant
func (f *FeatureFlag) AddVariant(variant Variant, updatedBy UserID) error {
	if updatedBy.IsEmpty() {
		return shared.NewValidationError("updated by user ID is required", nil)
	}

	// Check if variant key already exists
	for _, existingVariant := range f.variants {
		if existingVariant.Key == variant.Key {
			return shared.NewValidationError("variant key already exists", nil)
		}
	}

	f.variants = append(f.variants, variant)
	f.updatedAt = time.Now()
	f.version++

	f.addEvent(NewFeatureFlagUpdatedEvent(f))

	return nil
}

// RemoveVariant removes a variant
func (f *FeatureFlag) RemoveVariant(variantKey string, updatedBy UserID) error {
	if updatedBy.IsEmpty() {
		return shared.NewValidationError("updated by user ID is required", nil)
	}

	if len(f.variants) <= 1 {
		return shared.NewValidationError("cannot remove the last variant", nil)
	}

	for i, variant := range f.variants {
		if variant.Key == variantKey {
			f.variants = append(f.variants[:i], f.variants[i+1:]...)
			f.updatedAt = time.Now()
			f.version++

			f.addEvent(NewFeatureFlagUpdatedEvent(f))
			return nil
		}
	}

	return shared.NewValidationError("variant not found", nil)
}

// AddTargetingRule adds a new targeting rule
func (f *FeatureFlag) AddTargetingRule(rule TargetingRule, updatedBy UserID) error {
	if updatedBy.IsEmpty() {
		return shared.NewValidationError("updated by user ID is required", nil)
	}

	// Check if rule ID already exists
	for _, existingRule := range f.rules {
		if existingRule.ID == rule.ID {
			return shared.NewValidationError("targeting rule ID already exists", nil)
		}
	}

	f.rules = append(f.rules, rule)
	f.updatedAt = time.Now()
	f.version++

	f.addEvent(NewFeatureFlagUpdatedEvent(f))

	return nil
}

// RemoveTargetingRule removes a targeting rule
func (f *FeatureFlag) RemoveTargetingRule(ruleID string, updatedBy UserID) error {
	if updatedBy.IsEmpty() {
		return shared.NewValidationError("updated by user ID is required", nil)
	}

	for i, rule := range f.rules {
		if rule.ID == ruleID {
			f.rules = append(f.rules[:i], f.rules[i+1:]...)
			f.updatedAt = time.Now()
			f.version++

			f.addEvent(NewFeatureFlagUpdatedEvent(f))
			return nil
		}
	}

	return shared.NewValidationError("targeting rule not found", nil)
}

// Delete deletes the feature flag
func (f *FeatureFlag) Delete(deletedBy UserID) error {
	if deletedBy.IsEmpty() {
		return shared.NewValidationError("deleted by user ID is required", nil)
	}

	f.updatedAt = time.Now()
	f.version++

	f.addEvent(NewFeatureFlagDeletedEvent(f))

	return nil
}

// addEvent adds a domain event to the feature flag
func (f *FeatureFlag) addEvent(event shared.DomainEvent) {
	f.events = append(f.events, event)
}

// isValidFeatureFlagKey validates the feature flag key format
func isValidFeatureFlagKey(key string) bool {
	// Feature flag keys should be lowercase, alphanumeric with hyphens and underscores
	// and should not start or end with special characters
	if len(key) == 0 || len(key) > 100 {
		return false
	}

	// Check first and last characters
	if key[0] == '-' || key[0] == '_' || key[len(key)-1] == '-' || key[len(key)-1] == '_' {
		return false
	}

	// Check all characters are valid
	for _, char := range key {
		if !((char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '-' || char == '_') {
			return false
		}
	}

	return true
}
