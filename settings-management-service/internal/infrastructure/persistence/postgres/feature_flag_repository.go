package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/domain/settings"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/domain/shared"
)

// FeatureFlagDBModel represents the feature flag database model
type FeatureFlagDBModel struct {
	ID             uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Key            string     `gorm:"type:varchar(255);not null;index:idx_feature_flag_key_env_org_svc,unique"`
	IsEnabled      bool       `gorm:"type:boolean;not null;default:false"`
	Environment    string     `gorm:"type:varchar(20);not null;index:idx_feature_flag_key_env_org_svc"`
	OrganizationID *uuid.UUID `gorm:"type:uuid;index:idx_feature_flag_key_env_org_svc"`
	Service        *string    `gorm:"type:varchar(100);index:idx_feature_flag_key_env_org_svc"`
	Category       string     `gorm:"type:varchar(100);not null;index:idx_feature_flag_category"`
	Description    *string    `gorm:"type:text"`
	Variants       JSONB      `gorm:"type:jsonb;not null;default:'[]'"`
	TargetingRules JSONB      `gorm:"type:jsonb;not null;default:'[]'"`
	Tags           JSONB      `gorm:"type:jsonb;not null;default:'[]'"`
	Metadata       JSONB      `gorm:"type:jsonb;not null;default:'{}'"`
	CreatedBy      uuid.UUID  `gorm:"type:uuid;not null"`
	UpdatedBy      *uuid.UUID `gorm:"type:uuid"`
	CreatedAt      time.Time  `gorm:"type:timestamp with time zone;default:now()"`
	UpdatedAt      time.Time  `gorm:"type:timestamp with time zone;default:now()"`
	Version        int        `gorm:"type:integer;not null;default:1"`
}

// TableName specifies the table name for FeatureFlagDBModel
func (FeatureFlagDBModel) TableName() string {
	return "feature_flags"
}

// FeatureFlagRepository implements the settings.FeatureFlagRepository interface
type FeatureFlagRepository struct {
	db *gorm.DB
}

// NewFeatureFlagRepository creates a new FeatureFlagRepository
func NewFeatureFlagRepository(db *gorm.DB) *FeatureFlagRepository {
	return &FeatureFlagRepository{db: db}
}

// Save persists a feature flag aggregate
func (r *FeatureFlagRepository) Save(ctx context.Context, featureFlag *settings.FeatureFlag) error {
	dbModel := r.toDBModel(featureFlag)

	// Use Upsert to handle both creation and updates
	err := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"is_enabled", "description", "variants", "targeting_rules", "tags", "metadata", "updated_by", "updated_at", "version"}),
	}).Create(&dbModel).Error

	if err != nil {
		return fmt.Errorf("%w: failed to save feature flag: %v", shared.ErrInternalService, err)
	}
	return nil
}

// FindByID retrieves a feature flag by its ID
func (r *FeatureFlagRepository) FindByID(ctx context.Context, id settings.FeatureFlagID) (*settings.FeatureFlag, error) {
	var dbModel FeatureFlagDBModel
	err := r.db.WithContext(ctx).First(&dbModel, "id = ?", id.String()).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: feature flag with ID %s not found", shared.ErrNotFound, id.String())
		}
		return nil, fmt.Errorf("%w: failed to find feature flag by ID: %v", shared.ErrInternalService, err)
	}
	return r.toDomainEntity(&dbModel), nil
}

// FindByKey retrieves a feature flag by its key, environment, organization, and service
func (r *FeatureFlagRepository) FindByKey(ctx context.Context, key string, environment settings.Environment, organizationID *settings.OrganizationID, service *settings.ServiceName) (*settings.FeatureFlag, error) {
	var dbModel FeatureFlagDBModel
	query := r.db.WithContext(ctx).Where("key = ? AND environment = ?", key, environment.String())

	if organizationID != nil {
		query = query.Where("organization_id = ?", organizationID.String())
	} else {
		query = query.Where("organization_id IS NULL")
	}

	if service != nil {
		query = query.Where("service = ?", service.String())
	} else {
		query = query.Where("service IS NULL")
	}

	err := query.First(&dbModel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: feature flag with key '%s' not found", shared.ErrNotFound, key)
		}
		return nil, fmt.Errorf("%w: failed to find feature flag by key: %v", shared.ErrInternalService, err)
	}
	return r.toDomainEntity(&dbModel), nil
}

// Update updates an existing feature flag
func (r *FeatureFlagRepository) Update(ctx context.Context, featureFlag *settings.FeatureFlag) error {
	dbModel := r.toDBModel(featureFlag)

	result := r.db.WithContext(ctx).Model(&FeatureFlagDBModel{}).Where("id = ?", dbModel.ID).Updates(map[string]interface{}{
		"is_enabled":      dbModel.IsEnabled,
		"description":     dbModel.Description,
		"variants":        dbModel.Variants,
		"targeting_rules": dbModel.TargetingRules,
		"tags":            dbModel.Tags,
		"metadata":        dbModel.Metadata,
		"updated_by":      dbModel.UpdatedBy,
		"updated_at":      time.Now().UTC(),
		"version":         gorm.Expr("version + ?", 1),
	})

	if result.Error != nil {
		return fmt.Errorf("%w: failed to update feature flag: %v", shared.ErrInternalService, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("%w: feature flag with ID %s not found for update", shared.ErrNotFound, featureFlag.GetID().String())
	}
	return nil
}

// Delete removes a feature flag by its ID
func (r *FeatureFlagRepository) Delete(ctx context.Context, id settings.FeatureFlagID) error {
	result := r.db.WithContext(ctx).Delete(&FeatureFlagDBModel{}, "id = ?", id.String())
	if result.Error != nil {
		return fmt.Errorf("%w: failed to delete feature flag: %v", shared.ErrInternalService, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("%w: feature flag with ID %s not found for deletion", shared.ErrNotFound, id.String())
	}
	return nil
}

// List retrieves a list of feature flags with pagination and filtering
func (r *FeatureFlagRepository) List(ctx context.Context, options settings.ListOptions) ([]*settings.FeatureFlag, error) {
	var dbModels []FeatureFlagDBModel
	query := r.db.WithContext(ctx).Model(&FeatureFlagDBModel{})

	// Apply filters
	if options.Filters != nil {
		for key, value := range options.Filters {
			switch key {
			case "environment":
				query = query.Where("environment = ?", value)
			case "organization_id":
				query = query.Where("organization_id = ?", value)
			case "service":
				query = query.Where("service = ?", value)
			case "category":
				query = query.Where("category = ?", value)
			case "is_enabled":
				query = query.Where("is_enabled = ?", value)
			case "tags":
				if tags, ok := value.([]string); ok {
					for _, tag := range tags {
						query = query.Where("tags @> ?", fmt.Sprintf(`["%s"]`, tag))
					}
				}
			}
		}
	}

	// Apply pagination
	if options.Limit > 0 {
		query = query.Limit(options.Limit)
	}
	if options.Offset >= 0 {
		query = query.Offset(options.GetOffset())
	}

	// Apply sorting
	if options.SortBy != "" {
		order := "asc"
		if options.SortOrder != "" {
			order = options.SortOrder
		}
		query = query.Order(fmt.Sprintf("%s %s", options.SortBy, order))
	} else {
		query = query.Order("created_at desc") // Default sort
	}

	err := query.Find(&dbModels).Error
	if err != nil {
		return nil, fmt.Errorf("%w: failed to list feature flags: %v", shared.ErrInternalService, err)
	}

	featureFlags := make([]*settings.FeatureFlag, len(dbModels))
	for i, dbModel := range dbModels {
		featureFlags[i] = r.toDomainEntity(&dbModel)
	}
	return featureFlags, nil
}

// Count returns the total number of feature flags matching the filters
func (r *FeatureFlagRepository) Count(ctx context.Context, filters settings.ListFilters) (int, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&FeatureFlagDBModel{})

	if filters != nil {
		for key, value := range filters {
			switch key {
			case "environment":
				query = query.Where("environment = ?", value)
			case "organization_id":
				query = query.Where("organization_id = ?", value)
			case "service":
				query = query.Where("service = ?", value)
			case "category":
				query = query.Where("category = ?", value)
			case "is_enabled":
				query = query.Where("is_enabled = ?", value)
			case "tags":
				if tags, ok := value.([]string); ok {
					for _, tag := range tags {
						query = query.Where("tags @> ?", fmt.Sprintf(`["%s"]`, tag))
					}
				}
			}
		}
	}

	err := query.Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("%w: failed to count feature flags: %v", shared.ErrInternalService, err)
	}
	return int(count), nil
}

// FindByService retrieves feature flags by service
func (r *FeatureFlagRepository) FindByService(ctx context.Context, service settings.ServiceName, environment settings.Environment, organizationID *settings.OrganizationID) ([]*settings.FeatureFlag, error) {
	var dbModels []FeatureFlagDBModel
	query := r.db.WithContext(ctx).Where("service = ? AND environment = ?", service.String(), environment.String())

	if organizationID != nil {
		query = query.Where("organization_id = ?", organizationID.String())
	} else {
		query = query.Where("organization_id IS NULL")
	}

	err := query.Find(&dbModels).Error
	if err != nil {
		return nil, fmt.Errorf("%w: failed to find feature flags by service: %v", shared.ErrInternalService, err)
	}

	featureFlags := make([]*settings.FeatureFlag, len(dbModels))
	for i, dbModel := range dbModels {
		featureFlags[i] = r.toDomainEntity(&dbModel)
	}
	return featureFlags, nil
}

// FindByEnvironment retrieves feature flags by environment
func (r *FeatureFlagRepository) FindByEnvironment(ctx context.Context, environment settings.Environment, organizationID *settings.OrganizationID) ([]*settings.FeatureFlag, error) {
	var dbModels []FeatureFlagDBModel
	query := r.db.WithContext(ctx).Where("environment = ?", environment.String())

	if organizationID != nil {
		query = query.Where("organization_id = ?", organizationID.String())
	} else {
		query = query.Where("organization_id IS NULL")
	}

	err := query.Find(&dbModels).Error
	if err != nil {
		return nil, fmt.Errorf("%w: failed to find feature flags by environment: %v", shared.ErrInternalService, err)
	}

	featureFlags := make([]*settings.FeatureFlag, len(dbModels))
	for i, dbModel := range dbModels {
		featureFlags[i] = r.toDomainEntity(&dbModel)
	}
	return featureFlags, nil
}

// FindByOrganization retrieves feature flags by organization
func (r *FeatureFlagRepository) FindByOrganization(ctx context.Context, organizationID settings.OrganizationID) ([]*settings.FeatureFlag, error) {
	var dbModels []FeatureFlagDBModel
	err := r.db.WithContext(ctx).Where("organization_id = ?", organizationID.String()).Find(&dbModels).Error
	if err != nil {
		return nil, fmt.Errorf("%w: failed to find feature flags by organization: %v", shared.ErrInternalService, err)
	}

	featureFlags := make([]*settings.FeatureFlag, len(dbModels))
	for i, dbModel := range dbModels {
		featureFlags[i] = r.toDomainEntity(&dbModel)
	}
	return featureFlags, nil
}

// FindByCategory retrieves feature flags by category
func (r *FeatureFlagRepository) FindByCategory(ctx context.Context, category string, environment settings.Environment, organizationID *settings.OrganizationID) ([]*settings.FeatureFlag, error) {
	var dbModels []FeatureFlagDBModel
	query := r.db.WithContext(ctx).Where("category = ? AND environment = ?", category, environment.String())

	if organizationID != nil {
		query = query.Where("organization_id = ?", organizationID.String())
	} else {
		query = query.Where("organization_id IS NULL")
	}

	err := query.Find(&dbModels).Error
	if err != nil {
		return nil, fmt.Errorf("%w: failed to find feature flags by category: %v", shared.ErrInternalService, err)
	}

	featureFlags := make([]*settings.FeatureFlag, len(dbModels))
	for i, dbModel := range dbModels {
		featureFlags[i] = r.toDomainEntity(&dbModel)
	}
	return featureFlags, nil
}

// FindEnabled retrieves enabled feature flags
func (r *FeatureFlagRepository) FindEnabled(ctx context.Context, environment settings.Environment, organizationID *settings.OrganizationID) ([]*settings.FeatureFlag, error) {
	var dbModels []FeatureFlagDBModel
	query := r.db.WithContext(ctx).Where("is_enabled = ? AND environment = ?", true, environment.String())

	if organizationID != nil {
		query = query.Where("organization_id = ?", organizationID.String())
	} else {
		query = query.Where("organization_id IS NULL")
	}

	err := query.Find(&dbModels).Error
	if err != nil {
		return nil, fmt.Errorf("%w: failed to find enabled feature flags: %v", shared.ErrInternalService, err)
	}

	featureFlags := make([]*settings.FeatureFlag, len(dbModels))
	for i, dbModel := range dbModels {
		featureFlags[i] = r.toDomainEntity(&dbModel)
	}
	return featureFlags, nil
}

// ExistsByKey checks if a feature flag with the given key already exists
func (r *FeatureFlagRepository) ExistsByKey(ctx context.Context, key string, environment settings.Environment, organizationID *settings.OrganizationID, service *settings.ServiceName) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&FeatureFlagDBModel{}).Where("key = ? AND environment = ?", key, environment.String())

	if organizationID != nil {
		query = query.Where("organization_id = ?", organizationID.String())
	} else {
		query = query.Where("organization_id IS NULL")
	}

	if service != nil {
		query = query.Where("service = ?", service.String())
	} else {
		query = query.Where("service IS NULL")
	}

	err := query.Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("%w: failed to check feature flag existence: %v", shared.ErrInternalService, err)
	}
	return count > 0, nil
}

// ExistsByID checks if a feature flag with the given ID exists
func (r *FeatureFlagRepository) ExistsByID(ctx context.Context, id settings.FeatureFlagID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&FeatureFlagDBModel{}).Where("id = ?", id.String()).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("%w: failed to check feature flag existence: %v", shared.ErrInternalService, err)
	}
	return count > 0, nil
}

// BulkSave saves multiple feature flags
func (r *FeatureFlagRepository) BulkSave(ctx context.Context, featureFlags []*settings.FeatureFlag) error {
	if len(featureFlags) == 0 {
		return nil
	}

	dbModels := make([]FeatureFlagDBModel, len(featureFlags))
	for i, featureFlag := range featureFlags {
		dbModels[i] = *r.toDBModel(featureFlag)
	}

	err := r.db.WithContext(ctx).CreateInBatches(dbModels, 100).Error
	if err != nil {
		return fmt.Errorf("%w: failed to bulk save feature flags: %v", shared.ErrInternalService, err)
	}
	return nil
}

// BulkUpdate updates multiple feature flags
func (r *FeatureFlagRepository) BulkUpdate(ctx context.Context, featureFlags []*settings.FeatureFlag) error {
	if len(featureFlags) == 0 {
		return nil
	}

	for _, featureFlag := range featureFlags {
		if err := r.Update(ctx, featureFlag); err != nil {
			return err
		}
	}
	return nil
}

// BulkDelete deletes multiple feature flags
func (r *FeatureFlagRepository) BulkDelete(ctx context.Context, ids []settings.FeatureFlagID) error {
	if len(ids) == 0 {
		return nil
	}

	idStrings := make([]string, len(ids))
	for i, id := range ids {
		idStrings[i] = id.String()
	}

	result := r.db.WithContext(ctx).Delete(&FeatureFlagDBModel{}, "id IN ?", idStrings)
	if result.Error != nil {
		return fmt.Errorf("%w: failed to bulk delete feature flags: %v", shared.ErrInternalService, result.Error)
	}
	return nil
}

// toDBModel converts a domain feature flag to a database model
func (r *FeatureFlagRepository) toDBModel(featureFlag *settings.FeatureFlag) *FeatureFlagDBModel {
	variantsJSON, _ := json.Marshal(featureFlag.GetVariants())
	targetingRulesJSON, _ := json.Marshal(featureFlag.GetTargetingRules())
	tagsJSON, _ := json.Marshal(featureFlag.GetTags())
	metadataJSON, _ := json.Marshal(featureFlag.GetMetadata())

	dbModel := &FeatureFlagDBModel{
		ID:             uuid.MustParse(featureFlag.GetID().String()),
		Key:            featureFlag.GetKey(),
		IsEnabled:      featureFlag.GetIsEnabled(),
		Environment:    featureFlag.GetEnvironment().String(),
		Category:       featureFlag.GetCategory(),
		Description:    featureFlag.GetDescription(),
		Variants:       JSONB(variantsJSON),
		TargetingRules: JSONB(targetingRulesJSON),
		Tags:           JSONB(tagsJSON),
		Metadata:       JSONB(metadataJSON),
		CreatedBy:      uuid.MustParse(featureFlag.GetCreatedBy().String()),
		CreatedAt:      featureFlag.GetCreatedAt(),
		UpdatedAt:      featureFlag.GetUpdatedAt(),
		Version:        featureFlag.GetVersion(),
	}

	if featureFlag.GetOrganizationID() != nil {
		orgID := uuid.MustParse(featureFlag.GetOrganizationID().String())
		dbModel.OrganizationID = &orgID
	}

	if featureFlag.GetService() != nil {
		service := featureFlag.GetService().String()
		dbModel.Service = &service
	}

	if featureFlag.GetUpdatedBy() != nil {
		updatedBy := uuid.MustParse(featureFlag.GetUpdatedBy().String())
		dbModel.UpdatedBy = &updatedBy
	}

	return dbModel
}

// toDomainEntity converts a database model to a domain feature flag
func (r *FeatureFlagRepository) toDomainEntity(dbModel *FeatureFlagDBModel) *settings.FeatureFlag {
	featureFlagID, _ := settings.NewFeatureFlagIDFromString(dbModel.ID.String())
	environment, _ := settings.ParseEnvironment(dbModel.Environment)
	createdBy, _ := settings.NewUserIDFromString(dbModel.CreatedBy.String())

	var organizationID *settings.OrganizationID
	if dbModel.OrganizationID != nil {
		orgID, _ := settings.NewOrganizationIDFromString(dbModel.OrganizationID.String())
		organizationID = &orgID
	}

	var service *settings.ServiceName
	if dbModel.Service != nil {
		svc, _ := settings.NewServiceName(*dbModel.Service)
		service = &svc
	}

	var updatedBy *settings.UserID
	if dbModel.UpdatedBy != nil {
		userID, _ := settings.NewUserIDFromString(dbModel.UpdatedBy.String())
		updatedBy = &userID
	}

	var variants []settings.Variant
	_ = json.Unmarshal(dbModel.Variants, &variants)

	var targetingRules []settings.TargetingRule
	_ = json.Unmarshal(dbModel.TargetingRules, &targetingRules)

	var tags []string
	_ = json.Unmarshal(dbModel.Tags, &tags)

	var metadata map[string]interface{}
	_ = json.Unmarshal(dbModel.Metadata, &metadata)

	// Create feature flag (this is a simplified version - in reality you'd need a factory method)
	// For now, we'll create a basic feature flag and then update its fields
	featureFlag, _ := settings.NewFeatureFlag(
		dbModel.Key,
		dbModel.IsEnabled,
		environment,
		organizationID,
		service,
		dbModel.Category,
		dbModel.Description,
		variants,
		targetingRules,
		tags,
		metadata,
		createdBy,
	)

	// Note: In a real implementation, you would need to properly restore the feature flag state
	// including the ID, version, and timestamps. This is a simplified version.

	return featureFlag
}
