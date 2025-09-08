package postgres

import (
	"context"
	"database/sql/driver"
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

// JSONB is a custom type for handling JSONB columns in PostgreSQL
type JSONB json.RawMessage

// Scan implements the sql.Scanner interface for JSONB
func (j *JSONB) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}
	result := json.RawMessage{}
	err := json.Unmarshal(bytes, &result)
	*j = JSONB(result)
	return err
}

// Value implements the driver.Valuer interface for JSONB
func (j JSONB) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.RawMessage(j).MarshalJSON()
}

// ConfigurationDBModel represents the configuration database model
type ConfigurationDBModel struct {
	ID             uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Key            string     `gorm:"type:varchar(255);not null;index:idx_config_key_env_org_svc,unique"`
	Value          JSONB      `gorm:"type:jsonb;not null"`
	Environment    string     `gorm:"type:varchar(20);not null;index:idx_config_key_env_org_svc"`
	OrganizationID *uuid.UUID `gorm:"type:uuid;index:idx_config_key_env_org_svc"`
	Service        *string    `gorm:"type:varchar(100);index:idx_config_key_env_org_svc"`
	Category       string     `gorm:"type:varchar(100);not null;index:idx_config_category"`
	Description    *string    `gorm:"type:text"`
	Tags           JSONB      `gorm:"type:jsonb;not null;default:'[]'"`
	Metadata       JSONB      `gorm:"type:jsonb;not null;default:'{}'"`
	CreatedBy      uuid.UUID  `gorm:"type:uuid;not null"`
	UpdatedBy      *uuid.UUID `gorm:"type:uuid"`
	CreatedAt      time.Time  `gorm:"type:timestamp with time zone;default:now()"`
	UpdatedAt      time.Time  `gorm:"type:timestamp with time zone;default:now()"`
	Version        int        `gorm:"type:integer;not null;default:1"`
}

// TableName specifies the table name for ConfigurationDBModel
func (ConfigurationDBModel) TableName() string {
	return "configurations"
}

// ConfigurationRepository implements the settings.ConfigurationRepository interface
type ConfigurationRepository struct {
	db *gorm.DB
}

// NewConfigurationRepository creates a new ConfigurationRepository
func NewConfigurationRepository(db *gorm.DB) *ConfigurationRepository {
	return &ConfigurationRepository{db: db}
}

// Save persists a configuration aggregate
func (r *ConfigurationRepository) Save(ctx context.Context, config *settings.Configuration) error {
	dbModel := r.toDBModel(config)

	// Use Upsert to handle both creation and updates
	err := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"value", "description", "tags", "metadata", "updated_by", "updated_at", "version"}),
	}).Create(&dbModel).Error

	if err != nil {
		return fmt.Errorf("%w: failed to save configuration: %v", shared.ErrInternalService, err)
	}
	return nil
}

// FindByID retrieves a configuration by its ID
func (r *ConfigurationRepository) FindByID(ctx context.Context, id settings.ConfigurationID) (*settings.Configuration, error) {
	var dbModel ConfigurationDBModel
	err := r.db.WithContext(ctx).First(&dbModel, "id = ?", id.String()).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: configuration with ID %s not found", shared.ErrNotFound, id.String())
		}
		return nil, fmt.Errorf("%w: failed to find configuration by ID: %v", shared.ErrInternalService, err)
	}
	return r.toDomainEntity(&dbModel), nil
}

// FindByKey retrieves a configuration by its key, environment, organization, and service
func (r *ConfigurationRepository) FindByKey(ctx context.Context, key string, environment settings.Environment, organizationID *settings.OrganizationID, service *settings.ServiceName) (*settings.Configuration, error) {
	var dbModel ConfigurationDBModel
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
			return nil, fmt.Errorf("%w: configuration with key '%s' not found", shared.ErrNotFound, key)
		}
		return nil, fmt.Errorf("%w: failed to find configuration by key: %v", shared.ErrInternalService, err)
	}
	return r.toDomainEntity(&dbModel), nil
}

// Update updates an existing configuration
func (r *ConfigurationRepository) Update(ctx context.Context, config *settings.Configuration) error {
	dbModel := r.toDBModel(config)

	result := r.db.WithContext(ctx).Model(&ConfigurationDBModel{}).Where("id = ?", dbModel.ID).Updates(map[string]interface{}{
		"value":       dbModel.Value,
		"description": dbModel.Description,
		"tags":        dbModel.Tags,
		"metadata":    dbModel.Metadata,
		"updated_by":  dbModel.UpdatedBy,
		"updated_at":  time.Now().UTC(),
		"version":     gorm.Expr("version + ?", 1),
	})

	if result.Error != nil {
		return fmt.Errorf("%w: failed to update configuration: %v", shared.ErrInternalService, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("%w: configuration with ID %s not found for update", shared.ErrNotFound, config.GetID().String())
	}
	return nil
}

// Delete removes a configuration by its ID
func (r *ConfigurationRepository) Delete(ctx context.Context, id settings.ConfigurationID) error {
	result := r.db.WithContext(ctx).Delete(&ConfigurationDBModel{}, "id = ?", id.String())
	if result.Error != nil {
		return fmt.Errorf("%w: failed to delete configuration: %v", shared.ErrInternalService, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("%w: configuration with ID %s not found for deletion", shared.ErrNotFound, id.String())
	}
	return nil
}

// List retrieves a list of configurations with pagination and filtering
func (r *ConfigurationRepository) List(ctx context.Context, options settings.ListOptions) ([]*settings.Configuration, error) {
	var dbModels []ConfigurationDBModel
	query := r.db.WithContext(ctx).Model(&ConfigurationDBModel{})

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
		return nil, fmt.Errorf("%w: failed to list configurations: %v", shared.ErrInternalService, err)
	}

	configurations := make([]*settings.Configuration, len(dbModels))
	for i, dbModel := range dbModels {
		configurations[i] = r.toDomainEntity(&dbModel)
	}
	return configurations, nil
}

// Count returns the total number of configurations matching the filters
func (r *ConfigurationRepository) Count(ctx context.Context, filters settings.ListFilters) (int, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&ConfigurationDBModel{})

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
		return 0, fmt.Errorf("%w: failed to count configurations: %v", shared.ErrInternalService, err)
	}
	return int(count), nil
}

// FindByService retrieves configurations by service
func (r *ConfigurationRepository) FindByService(ctx context.Context, service settings.ServiceName, environment settings.Environment, organizationID *settings.OrganizationID) ([]*settings.Configuration, error) {
	var dbModels []ConfigurationDBModel
	query := r.db.WithContext(ctx).Where("service = ? AND environment = ?", service.String(), environment.String())

	if organizationID != nil {
		query = query.Where("organization_id = ?", organizationID.String())
	} else {
		query = query.Where("organization_id IS NULL")
	}

	err := query.Find(&dbModels).Error
	if err != nil {
		return nil, fmt.Errorf("%w: failed to find configurations by service: %v", shared.ErrInternalService, err)
	}

	configurations := make([]*settings.Configuration, len(dbModels))
	for i, dbModel := range dbModels {
		configurations[i] = r.toDomainEntity(&dbModel)
	}
	return configurations, nil
}

// FindByEnvironment retrieves configurations by environment
func (r *ConfigurationRepository) FindByEnvironment(ctx context.Context, environment settings.Environment, organizationID *settings.OrganizationID) ([]*settings.Configuration, error) {
	var dbModels []ConfigurationDBModel
	query := r.db.WithContext(ctx).Where("environment = ?", environment.String())

	if organizationID != nil {
		query = query.Where("organization_id = ?", organizationID.String())
	} else {
		query = query.Where("organization_id IS NULL")
	}

	err := query.Find(&dbModels).Error
	if err != nil {
		return nil, fmt.Errorf("%w: failed to find configurations by environment: %v", shared.ErrInternalService, err)
	}

	configurations := make([]*settings.Configuration, len(dbModels))
	for i, dbModel := range dbModels {
		configurations[i] = r.toDomainEntity(&dbModel)
	}
	return configurations, nil
}

// FindByOrganization retrieves configurations by organization
func (r *ConfigurationRepository) FindByOrganization(ctx context.Context, organizationID settings.OrganizationID) ([]*settings.Configuration, error) {
	var dbModels []ConfigurationDBModel
	err := r.db.WithContext(ctx).Where("organization_id = ?", organizationID.String()).Find(&dbModels).Error
	if err != nil {
		return nil, fmt.Errorf("%w: failed to find configurations by organization: %v", shared.ErrInternalService, err)
	}

	configurations := make([]*settings.Configuration, len(dbModels))
	for i, dbModel := range dbModels {
		configurations[i] = r.toDomainEntity(&dbModel)
	}
	return configurations, nil
}

// FindByCategory retrieves configurations by category
func (r *ConfigurationRepository) FindByCategory(ctx context.Context, category string, environment settings.Environment, organizationID *settings.OrganizationID) ([]*settings.Configuration, error) {
	var dbModels []ConfigurationDBModel
	query := r.db.WithContext(ctx).Where("category = ? AND environment = ?", category, environment.String())

	if organizationID != nil {
		query = query.Where("organization_id = ?", organizationID.String())
	} else {
		query = query.Where("organization_id IS NULL")
	}

	err := query.Find(&dbModels).Error
	if err != nil {
		return nil, fmt.Errorf("%w: failed to find configurations by category: %v", shared.ErrInternalService, err)
	}

	configurations := make([]*settings.Configuration, len(dbModels))
	for i, dbModel := range dbModels {
		configurations[i] = r.toDomainEntity(&dbModel)
	}
	return configurations, nil
}

// ExistsByKey checks if a configuration with the given key already exists
func (r *ConfigurationRepository) ExistsByKey(ctx context.Context, key string, environment settings.Environment, organizationID *settings.OrganizationID, service *settings.ServiceName) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&ConfigurationDBModel{}).Where("key = ? AND environment = ?", key, environment.String())

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
		return false, fmt.Errorf("%w: failed to check configuration existence: %v", shared.ErrInternalService, err)
	}
	return count > 0, nil
}

// ExistsByID checks if a configuration with the given ID exists
func (r *ConfigurationRepository) ExistsByID(ctx context.Context, id settings.ConfigurationID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&ConfigurationDBModel{}).Where("id = ?", id.String()).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("%w: failed to check configuration existence: %v", shared.ErrInternalService, err)
	}
	return count > 0, nil
}

// BulkSave saves multiple configurations
func (r *ConfigurationRepository) BulkSave(ctx context.Context, configurations []*settings.Configuration) error {
	if len(configurations) == 0 {
		return nil
	}

	dbModels := make([]ConfigurationDBModel, len(configurations))
	for i, config := range configurations {
		dbModels[i] = *r.toDBModel(config)
	}

	err := r.db.WithContext(ctx).CreateInBatches(dbModels, 100).Error
	if err != nil {
		return fmt.Errorf("%w: failed to bulk save configurations: %v", shared.ErrInternalService, err)
	}
	return nil
}

// BulkUpdate updates multiple configurations
func (r *ConfigurationRepository) BulkUpdate(ctx context.Context, configurations []*settings.Configuration) error {
	if len(configurations) == 0 {
		return nil
	}

	for _, config := range configurations {
		if err := r.Update(ctx, config); err != nil {
			return err
		}
	}
	return nil
}

// BulkDelete deletes multiple configurations
func (r *ConfigurationRepository) BulkDelete(ctx context.Context, ids []settings.ConfigurationID) error {
	if len(ids) == 0 {
		return nil
	}

	idStrings := make([]string, len(ids))
	for i, id := range ids {
		idStrings[i] = id.String()
	}

	result := r.db.WithContext(ctx).Delete(&ConfigurationDBModel{}, "id IN ?", idStrings)
	if result.Error != nil {
		return fmt.Errorf("%w: failed to bulk delete configurations: %v", shared.ErrInternalService, result.Error)
	}
	return nil
}

// toDBModel converts a domain configuration to a database model
func (r *ConfigurationRepository) toDBModel(config *settings.Configuration) *ConfigurationDBModel {
	valueJSON, _ := json.Marshal(config.GetValue())
	tagsJSON, _ := json.Marshal(config.GetTags())
	metadataJSON, _ := json.Marshal(config.GetMetadata())

	dbModel := &ConfigurationDBModel{
		ID:          uuid.MustParse(config.GetID().String()),
		Key:         config.GetKey(),
		Value:       JSONB(valueJSON),
		Environment: config.GetEnvironment().String(),
		Category:    config.GetCategory(),
		Description: config.GetDescription(),
		Tags:        JSONB(tagsJSON),
		Metadata:    JSONB(metadataJSON),
		CreatedBy:   uuid.MustParse(config.GetCreatedBy().String()),
		CreatedAt:   config.GetCreatedAt(),
		UpdatedAt:   config.GetUpdatedAt(),
		Version:     config.GetVersion(),
	}

	if config.GetOrganizationID() != nil {
		orgID := uuid.MustParse(config.GetOrganizationID().String())
		dbModel.OrganizationID = &orgID
	}

	if config.GetService() != nil {
		service := config.GetService().String()
		dbModel.Service = &service
	}

	if config.GetUpdatedBy() != nil {
		updatedBy := uuid.MustParse(config.GetUpdatedBy().String())
		dbModel.UpdatedBy = &updatedBy
	}

	return dbModel
}

// toDomainEntity converts a database model to a domain configuration
func (r *ConfigurationRepository) toDomainEntity(dbModel *ConfigurationDBModel) *settings.Configuration {
	configID, _ := settings.NewConfigurationIDFromString(dbModel.ID.String())
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

	var value interface{}
	_ = json.Unmarshal(dbModel.Value, &value)

	var tags []string
	_ = json.Unmarshal(dbModel.Tags, &tags)

	var metadata map[string]interface{}
	_ = json.Unmarshal(dbModel.Metadata, &metadata)

	// Create configuration (this is a simplified version - in reality you'd need a factory method)
	// For now, we'll create a basic configuration and then update its fields
	configuration, _ := settings.NewConfiguration(
		dbModel.Key,
		value,
		environment,
		organizationID,
		service,
		dbModel.Category,
		dbModel.Description,
		tags,
		metadata,
		createdBy,
	)

	// Note: In a real implementation, you would need to properly restore the configuration state
	// including the ID, version, and timestamps. This is a simplified version.

	return configuration
}
