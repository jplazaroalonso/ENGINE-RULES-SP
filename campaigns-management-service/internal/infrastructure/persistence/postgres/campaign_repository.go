package postgres

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/campaigns-management-service/internal/domain/campaign"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/campaigns-management-service/internal/domain/shared"
	"gorm.io/gorm"
)

// CampaignDBModel represents the database model for campaigns
type CampaignDBModel struct {
	ID             uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name           string         `gorm:"type:varchar(255);not null;unique"`
	Description    string         `gorm:"type:text"`
	Status         string         `gorm:"type:varchar(50);not null;default:'DRAFT'"`
	CampaignType   string         `gorm:"type:varchar(50);not null"`
	TargetingRules JSONB          `gorm:"type:jsonb;not null;default:'[]'"`
	StartDate      time.Time      `gorm:"type:timestamp with time zone;not null"`
	EndDate        *time.Time     `gorm:"type:timestamp with time zone"`
	BudgetAmount   *float64       `gorm:"type:decimal(15,2)"`
	BudgetCurrency *string        `gorm:"type:varchar(3)"`
	CreatedBy      string         `gorm:"type:varchar(255);not null"`
	CreatedAt      time.Time      `gorm:"type:timestamp with time zone;default:now()"`
	UpdatedAt      time.Time      `gorm:"type:timestamp with time zone;default:now()"`
	Settings       JSONB          `gorm:"type:jsonb;not null;default:'{}'"`
	Version        int            `gorm:"type:integer;not null;default:1"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

// CampaignMetricsDBModel represents the database model for campaign metrics
type CampaignMetricsDBModel struct {
	ID                uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CampaignID        uuid.UUID `gorm:"type:uuid;not null;index"`
	Impressions       int64     `gorm:"type:bigint;not null;default:0"`
	Clicks            int64     `gorm:"type:bigint;not null;default:0"`
	Conversions       int64     `gorm:"type:bigint;not null;default:0"`
	RevenueAmount     float64   `gorm:"type:decimal(15,2);not null;default:0"`
	RevenueCurrency   string    `gorm:"type:varchar(3);not null;default:'EUR'"`
	CostAmount        float64   `gorm:"type:decimal(15,2);not null;default:0"`
	CostCurrency      string    `gorm:"type:varchar(3);not null;default:'EUR'"`
	CTR               float64   `gorm:"type:decimal(10,4);not null;default:0"`
	ConversionRate    float64   `gorm:"type:decimal(10,4);not null;default:0"`
	CostPerClick      float64   `gorm:"type:decimal(15,2);not null;default:0"`
	CostPerConversion float64   `gorm:"type:decimal(15,2);not null;default:0"`
	ROAS              float64   `gorm:"type:decimal(10,4);not null;default:0"`
	ROI               float64   `gorm:"type:decimal(10,4);not null;default:0"`
	LastUpdated       time.Time `gorm:"type:timestamp with time zone;default:now()"`
}

// CampaignEventDBModel represents the database model for campaign events
type CampaignEventDBModel struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CampaignID      uuid.UUID `gorm:"type:uuid;not null;index"`
	EventType       string    `gorm:"type:varchar(50);not null"`
	CustomerID      *string   `gorm:"type:varchar(255)"`
	EventData       JSONB     `gorm:"type:jsonb;not null;default:'{}'"`
	OccurredAt      time.Time `gorm:"type:timestamp with time zone;default:now()"`
	RevenueAmount   *float64  `gorm:"type:decimal(15,2)"`
	RevenueCurrency *string   `gorm:"type:varchar(3)"`
	CostAmount      *float64  `gorm:"type:decimal(15,2)"`
	CostCurrency    *string   `gorm:"type:varchar(3)"`
}

// JSONB is a custom type for PostgreSQL JSONB fields
type JSONB map[string]interface{}

// Scan implements the sql.Scanner interface for JSONB
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = make(map[string]interface{})
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot scan %T into JSONB", value)
	}

	return json.Unmarshal(bytes, j)
}

// Value implements the driver.Valuer interface for JSONB
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// CampaignRepository implements the campaign repository interface using PostgreSQL
type CampaignRepository struct {
	db *gorm.DB
}

// NewCampaignRepository creates a new campaign repository
func NewCampaignRepository(db *gorm.DB) *CampaignRepository {
	return &CampaignRepository{db: db}
}

// Save saves a campaign to the database
func (r *CampaignRepository) Save(ctx context.Context, campaign *campaign.Campaign) error {
	campaignDB := r.toDBModel(campaign)

	// Use GORM's Save method which will insert or update based on primary key
	if err := r.db.WithContext(ctx).Save(campaignDB).Error; err != nil {
		return shared.NewInfrastructureError("failed to save campaign", err)
	}

	// Update the campaign's version
	campaign = r.toDomainEntity(campaignDB)

	return nil
}

// FindByID finds a campaign by ID
func (r *CampaignRepository) FindByID(ctx context.Context, id campaign.CampaignID) (*campaign.Campaign, error) {
	var campaignDB CampaignDBModel

	if err := r.db.WithContext(ctx).Where("id = ?", id.String()).First(&campaignDB).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, shared.NewInfrastructureError("failed to find campaign by ID", err)
	}

	return r.toDomainEntity(&campaignDB), nil
}

// FindByName finds a campaign by name
func (r *CampaignRepository) FindByName(ctx context.Context, name string) (*campaign.Campaign, error) {
	var campaignDB CampaignDBModel

	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&campaignDB).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, shared.NewInfrastructureError("failed to find campaign by name", err)
	}

	return r.toDomainEntity(&campaignDB), nil
}

// List lists campaigns with filtering and pagination
func (r *CampaignRepository) List(ctx context.Context, criteria campaign.ListCriteria) ([]*campaign.Campaign, error) {
	var campaignsDB []CampaignDBModel
	query := r.db.WithContext(ctx).Model(&CampaignDBModel{})

	// Apply filters
	if criteria.Status != nil {
		query = query.Where("status = ?", criteria.Status.String())
	}

	if criteria.CampaignType != nil {
		query = query.Where("campaign_type = ?", criteria.CampaignType.String())
	}

	if criteria.CreatedBy != nil {
		query = query.Where("created_by = ?", criteria.CreatedBy.String())
	}

	if criteria.StartDate != nil {
		query = query.Where("start_date >= ?", *criteria.StartDate)
	}

	if criteria.EndDate != nil {
		query = query.Where("end_date <= ?", *criteria.EndDate)
	}

	if criteria.SearchQuery != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?",
			"%"+criteria.SearchQuery+"%", "%"+criteria.SearchQuery+"%")
	}

	// Apply sorting
	sortOrder := "ASC"
	if criteria.SortOrder == "desc" {
		sortOrder = "DESC"
	}
	query = query.Order(fmt.Sprintf("%s %s", criteria.SortBy, sortOrder))

	// Apply pagination
	query = query.Offset(criteria.PageOffset).Limit(criteria.PageSize)

	if err := query.Find(&campaignsDB).Error; err != nil {
		return nil, shared.NewInfrastructureError("failed to list campaigns", err)
	}

	// Convert to domain entities
	campaigns := make([]*campaign.Campaign, len(campaignsDB))
	for i, campaignDB := range campaignsDB {
		campaigns[i] = r.toDomainEntity(&campaignDB)
	}

	return campaigns, nil
}

// Count counts campaigns with filtering
func (r *CampaignRepository) Count(ctx context.Context, criteria campaign.ListCriteria) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&CampaignDBModel{})

	// Apply filters (same as List method)
	if criteria.Status != nil {
		query = query.Where("status = ?", criteria.Status.String())
	}

	if criteria.CampaignType != nil {
		query = query.Where("campaign_type = ?", criteria.CampaignType.String())
	}

	if criteria.CreatedBy != nil {
		query = query.Where("created_by = ?", criteria.CreatedBy.String())
	}

	if criteria.StartDate != nil {
		query = query.Where("start_date >= ?", *criteria.StartDate)
	}

	if criteria.EndDate != nil {
		query = query.Where("end_date <= ?", *criteria.EndDate)
	}

	if criteria.SearchQuery != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?",
			"%"+criteria.SearchQuery+"%", "%"+criteria.SearchQuery+"%")
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, shared.NewInfrastructureError("failed to count campaigns", err)
	}

	return count, nil
}

// Delete soft deletes a campaign
func (r *CampaignRepository) Delete(ctx context.Context, id campaign.CampaignID) error {
	if err := r.db.WithContext(ctx).Where("id = ?", id.String()).Delete(&CampaignDBModel{}).Error; err != nil {
		return shared.NewInfrastructureError("failed to delete campaign", err)
	}

	return nil
}

// ExistsByName checks if a campaign with the given name exists
func (r *CampaignRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	var count int64

	if err := r.db.WithContext(ctx).Model(&CampaignDBModel{}).Where("name = ?", name).Count(&count).Error; err != nil {
		return false, shared.NewInfrastructureError("failed to check campaign name existence", err)
	}

	return count > 0, nil
}

// FindByStatus finds campaigns by status
func (r *CampaignRepository) FindByStatus(ctx context.Context, status campaign.CampaignStatus) ([]*campaign.Campaign, error) {
	var campaignsDB []CampaignDBModel

	if err := r.db.WithContext(ctx).Where("status = ?", status.String()).Find(&campaignsDB).Error; err != nil {
		return nil, shared.NewInfrastructureError("failed to find campaigns by status", err)
	}

	campaigns := make([]*campaign.Campaign, len(campaignsDB))
	for i, campaignDB := range campaignsDB {
		campaigns[i] = r.toDomainEntity(&campaignDB)
	}

	return campaigns, nil
}

// FindByType finds campaigns by type
func (r *CampaignRepository) FindByType(ctx context.Context, campaignType campaign.CampaignType) ([]*campaign.Campaign, error) {
	var campaignsDB []CampaignDBModel

	if err := r.db.WithContext(ctx).Where("campaign_type = ?", campaignType.String()).Find(&campaignsDB).Error; err != nil {
		return nil, shared.NewInfrastructureError("failed to find campaigns by type", err)
	}

	campaigns := make([]*campaign.Campaign, len(campaignsDB))
	for i, campaignDB := range campaignsDB {
		campaigns[i] = r.toDomainEntity(&campaignDB)
	}

	return campaigns, nil
}

// FindByDateRange finds campaigns within a date range
func (r *CampaignRepository) FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*campaign.Campaign, error) {
	var campaignsDB []CampaignDBModel

	query := r.db.WithContext(ctx).Where("start_date <= ? AND (end_date IS NULL OR end_date >= ?)", endDate, startDate)

	if err := query.Find(&campaignsDB).Error; err != nil {
		return nil, shared.NewInfrastructureError("failed to find campaigns by date range", err)
	}

	campaigns := make([]*campaign.Campaign, len(campaignsDB))
	for i, campaignDB := range campaignsDB {
		campaigns[i] = r.toDomainEntity(&campaignDB)
	}

	return campaigns, nil
}

// FindByCreatedBy finds campaigns by creator
func (r *CampaignRepository) FindByCreatedBy(ctx context.Context, createdBy shared.UserID) ([]*campaign.Campaign, error) {
	var campaignsDB []CampaignDBModel

	if err := r.db.WithContext(ctx).Where("created_by = ?", createdBy.String()).Find(&campaignsDB).Error; err != nil {
		return nil, shared.NewInfrastructureError("failed to find campaigns by creator", err)
	}

	campaigns := make([]*campaign.Campaign, len(campaignsDB))
	for i, campaignDB := range campaignsDB {
		campaigns[i] = r.toDomainEntity(&campaignDB)
	}

	return campaigns, nil
}

// toDBModel converts a domain campaign to database model
func (r *CampaignRepository) toDBModel(campaign *campaign.Campaign) *CampaignDBModel {
	// Parse campaign ID
	campaignID, _ := uuid.Parse(campaign.ID().String())

	// Convert targeting rules to JSONB
	targetingRules := make([]string, len(campaign.TargetingRules()))
	for i, rule := range campaign.TargetingRules() {
		targetingRules[i] = rule.String()
	}
	targetingRulesJSONB := JSONB{"rules": targetingRules}

	// Convert settings to JSONB
	settingsJSONB := r.settingsToJSONB(campaign.Settings())

	// Convert budget
	var budgetAmount *float64
	var budgetCurrency *string
	if campaign.Budget() != nil {
		budgetAmount = &campaign.Budget().Amount
		budgetCurrency = &campaign.Budget().Currency
	}

	return &CampaignDBModel{
		ID:             campaignID,
		Name:           campaign.Name(),
		Description:    campaign.Description(),
		Status:         campaign.Status().String(),
		CampaignType:   campaign.CampaignType().String(),
		TargetingRules: targetingRulesJSONB,
		StartDate:      campaign.StartDate(),
		EndDate:        campaign.EndDate(),
		BudgetAmount:   budgetAmount,
		BudgetCurrency: budgetCurrency,
		CreatedBy:      campaign.CreatedBy().String(),
		CreatedAt:      campaign.CreatedAt(),
		UpdatedAt:      campaign.UpdatedAt(),
		Settings:       settingsJSONB,
		Version:        campaign.Version(),
	}
}

// toDomainEntity converts a database model to domain campaign
func (r *CampaignRepository) toDomainEntity(campaignDB *CampaignDBModel) *campaign.Campaign {
	// Parse campaign type
	campaignType, _ := campaign.ParseCampaignType(campaignDB.CampaignType)

	// Convert targeting rules from JSONB
	targetingRules := r.jsonBToTargetingRules(campaignDB.TargetingRules)

	// Convert budget
	var budget *shared.Money
	if campaignDB.BudgetAmount != nil && campaignDB.BudgetCurrency != nil {
		budget = &shared.Money{
			Amount:   *campaignDB.BudgetAmount,
			Currency: *campaignDB.BudgetCurrency,
		}
	}

	// Convert settings from JSONB
	settings := r.jsonBToSettings(campaignDB.Settings)

	// Parse created by user ID
	createdBy, _ := shared.NewUserIDFromString(campaignDB.CreatedBy)

	// Create campaign (this is a simplified version - in reality you'd need a factory method)
	// For now, we'll create a basic campaign and then update its fields
	campaign, _ := campaign.NewCampaign(
		campaignDB.Name,
		campaignDB.Description,
		campaignType,
		targetingRules,
		campaignDB.StartDate,
		campaignDB.EndDate,
		budget,
		createdBy,
		settings,
	)

	// Note: In a real implementation, you would need to properly restore the campaign state
	// including the ID, status, version, and timestamps. This is a simplified version.

	return campaign
}

// settingsToJSONB converts campaign settings to JSONB
func (r *CampaignRepository) settingsToJSONB(settings campaign.CampaignSettings) JSONB {
	// Convert channels
	channels := make([]string, len(settings.Channels))
	for i, channel := range settings.Channels {
		channels[i] = channel.String()
	}

	// Convert frequency
	frequency := settings.Frequency.String()

	// Convert A/B test settings
	var abTestSettings map[string]interface{}
	if settings.ABTestSettings != nil {
		variants := make([]map[string]interface{}, len(settings.ABTestSettings.Variants))
		for i, variant := range settings.ABTestSettings.Variants {
			variants[i] = map[string]interface{}{
				"id":          variant.ID,
				"name":        variant.Name,
				"description": variant.Description,
				"settings":    variant.Settings,
				"weight":      variant.Weight,
			}
		}

		abTestSettings = map[string]interface{}{
			"enabled":       settings.ABTestSettings.Enabled,
			"variants":      variants,
			"trafficSplit":  settings.ABTestSettings.TrafficSplit,
			"successMetric": settings.ABTestSettings.SuccessMetric,
			"duration":      settings.ABTestSettings.Duration,
		}
	}

	// Convert scheduling rules
	schedulingRules := make([]map[string]interface{}, len(settings.SchedulingRules))
	for i, rule := range settings.SchedulingRules {
		conditions := make([]map[string]interface{}, len(rule.Conditions))
		for j, condition := range rule.Conditions {
			conditions[j] = map[string]interface{}{
				"type":     condition.Type,
				"operator": condition.Operator,
				"value":    condition.Value,
				"metadata": condition.Metadata,
			}
		}

		actions := make([]map[string]interface{}, len(rule.Actions))
		for j, action := range rule.Actions {
			actions[j] = map[string]interface{}{
				"type":       action.Type,
				"parameters": action.Parameters,
			}
		}

		schedulingRules[i] = map[string]interface{}{
			"id":          rule.ID,
			"name":        rule.Name,
			"description": rule.Description,
			"conditions":  conditions,
			"actions":     actions,
			"isActive":    rule.IsActive,
		}
	}

	// Convert personalization config
	personalization := map[string]interface{}{
		"enabled":     settings.Personalization.Enabled,
		"rules":       settings.Personalization.Rules,
		"fallback":    settings.Personalization.Fallback,
		"maxVariants": settings.Personalization.MaxVariants,
	}

	// Convert budget limit
	var budgetLimit map[string]interface{}
	if settings.BudgetLimit != nil {
		budgetLimit = map[string]interface{}{
			"amount":   settings.BudgetLimit.Amount,
			"currency": settings.BudgetLimit.Currency,
		}
	}

	return JSONB{
		"targetAudience":  settings.TargetAudience,
		"channels":        channels,
		"frequency":       frequency,
		"maxImpressions":  settings.MaxImpressions,
		"budgetLimit":     budgetLimit,
		"abTestSettings":  abTestSettings,
		"schedulingRules": schedulingRules,
		"personalization": personalization,
	}
}

// jsonBToTargetingRules converts JSONB to targeting rules
func (r *CampaignRepository) jsonBToTargetingRules(jsonb JSONB) []shared.RuleID {
	rules, ok := jsonb["rules"].([]interface{})
	if !ok {
		return []shared.RuleID{}
	}

	targetingRules := make([]shared.RuleID, len(rules))
	for i, rule := range rules {
		if ruleStr, ok := rule.(string); ok {
			if ruleID, err := shared.NewRuleIDFromString(ruleStr); err == nil {
				targetingRules[i] = ruleID
			}
		}
	}

	return targetingRules
}

// jsonBToSettings converts JSONB to campaign settings
func (r *CampaignRepository) jsonBToSettings(jsonb JSONB) campaign.CampaignSettings {
	// This is a simplified version - in reality you'd need to properly parse all fields
	// and handle type conversions and validation

	targetAudience := []string{}
	if ta, ok := jsonb["targetAudience"].([]interface{}); ok {
		for _, item := range ta {
			if str, ok := item.(string); ok {
				targetAudience = append(targetAudience, str)
			}
		}
	}

	channels := []campaign.Channel{}
	if ch, ok := jsonb["channels"].([]interface{}); ok {
		for _, item := range ch {
			if str, ok := item.(string); ok {
				if channel, err := campaign.ParseChannel(str); err == nil {
					channels = append(channels, channel)
				}
			}
		}
	}

	frequency := campaign.FrequencyOnce
	if f, ok := jsonb["frequency"].(string); ok {
		if parsed, err := campaign.ParseFrequency(f); err == nil {
			frequency = parsed
		}
	}

	// Create basic settings - in reality you'd parse all fields
	settings, _ := campaign.NewCampaignSettings(
		targetAudience,
		channels,
		frequency,
		nil,                              // maxImpressions
		nil,                              // budgetLimit
		nil,                              // abTestSettings
		[]campaign.SchedulingRule{},      // schedulingRules
		campaign.PersonalizationConfig{}, // personalization
	)

	return settings
}
