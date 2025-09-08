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

// OrganizationSettingDBModel represents the organization setting database model
type OrganizationSettingDBModel struct {
	ID             uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OrganizationID uuid.UUID  `gorm:"type:uuid;not null;index:idx_org_setting_org_cat_key,unique"`
	Category       string     `gorm:"type:varchar(100);not null;index:idx_org_setting_org_cat_key"`
	Key            string     `gorm:"type:varchar(255);not null;index:idx_org_setting_org_cat_key"`
	Value          JSONB      `gorm:"type:jsonb;not null"`
	ParentID       *uuid.UUID `gorm:"type:uuid"`
	Description    *string    `gorm:"type:text"`
	Tags           JSONB      `gorm:"type:jsonb;not null;default:'[]'"`
	Metadata       JSONB      `gorm:"type:jsonb;not null;default:'{}'"`
	CreatedBy      uuid.UUID  `gorm:"type:uuid;not null"`
	UpdatedBy      *uuid.UUID `gorm:"type:uuid"`
	CreatedAt      time.Time  `gorm:"type:timestamp with time zone;default:now()"`
	UpdatedAt      time.Time  `gorm:"type:timestamp with time zone;default:now()"`
	Version        int        `gorm:"type:integer;not null;default:1"`
}

// TableName specifies the table name for OrganizationSettingDBModel
func (OrganizationSettingDBModel) TableName() string {
	return "organization_settings"
}

// OrganizationSettingRepository implements the settings.OrganizationSettingRepository interface
type OrganizationSettingRepository struct {
	db *gorm.DB
}

// NewOrganizationSettingRepository creates a new OrganizationSettingRepository
func NewOrganizationSettingRepository(db *gorm.DB) *OrganizationSettingRepository {
	return &OrganizationSettingRepository{db: db}
}

// Save persists an organization setting aggregate
func (r *OrganizationSettingRepository) Save(ctx context.Context, setting *settings.OrganizationSetting) error {
	dbModel := r.toDBModel(setting)

	// Use Upsert to handle both creation and updates
	err := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"value", "description", "tags", "metadata", "updated_by", "updated_at", "version"}),
	}).Create(&dbModel).Error

	if err != nil {
		return fmt.Errorf("%w: failed to save organization setting: %v", shared.ErrInternalService, err)
	}
	return nil
}

// FindByID retrieves an organization setting by its ID
func (r *OrganizationSettingRepository) FindByID(ctx context.Context, id settings.OrganizationSettingID) (*settings.OrganizationSetting, error) {
	var dbModel OrganizationSettingDBModel
	err := r.db.WithContext(ctx).First(&dbModel, "id = ?", id.String()).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: organization setting with ID %s not found", shared.ErrNotFound, id.String())
		}
		return nil, fmt.Errorf("%w: failed to find organization setting by ID: %v", shared.ErrInternalService, err)
	}
	return r.toDomainEntity(&dbModel), nil
}

// FindByKey retrieves an organization setting by its organization ID, category, and key
func (r *OrganizationSettingRepository) FindByKey(ctx context.Context, organizationID settings.OrganizationID, category string, key string) (*settings.OrganizationSetting, error) {
	var dbModel OrganizationSettingDBModel
	err := r.db.WithContext(ctx).Where("organization_id = ? AND category = ? AND key = ?", organizationID.String(), category, key).First(&dbModel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: organization setting with key '%s' not found for organization '%s' in category '%s'", shared.ErrNotFound, key, organizationID.String(), category)
		}
		return nil, fmt.Errorf("%w: failed to find organization setting by key: %v", shared.ErrInternalService, err)
	}
	return r.toDomainEntity(&dbModel), nil
}

// Update updates an existing organization setting
func (r *OrganizationSettingRepository) Update(ctx context.Context, setting *settings.OrganizationSetting) error {
	dbModel := r.toDBModel(setting)

	result := r.db.WithContext(ctx).Model(&OrganizationSettingDBModel{}).Where("id = ?", dbModel.ID).Updates(map[string]interface{}{
		"value":       dbModel.Value,
		"description": dbModel.Description,
		"tags":        dbModel.Tags,
		"metadata":    dbModel.Metadata,
		"updated_by":  dbModel.UpdatedBy,
		"updated_at":  time.Now().UTC(),
		"version":     gorm.Expr("version + ?", 1),
	})

	if result.Error != nil {
		return fmt.Errorf("%w: failed to update organization setting: %v", shared.ErrInternalService, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("%w: organization setting with ID %s not found for update", shared.ErrNotFound, setting.GetID().String())
	}
	return nil
}

// Delete removes an organization setting by its ID
func (r *OrganizationSettingRepository) Delete(ctx context.Context, id settings.OrganizationSettingID) error {
	result := r.db.WithContext(ctx).Delete(&OrganizationSettingDBModel{}, "id = ?", id.String())
	if result.Error != nil {
		return fmt.Errorf("%w: failed to delete organization setting: %v", shared.ErrInternalService, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("%w: organization setting with ID %s not found for deletion", shared.ErrNotFound, id.String())
	}
	return nil
}

// List retrieves a list of organization settings with pagination and filtering
func (r *OrganizationSettingRepository) List(ctx context.Context, options settings.ListOptions) ([]*settings.OrganizationSetting, error) {
	var dbModels []OrganizationSettingDBModel
	query := r.db.WithContext(ctx).Model(&OrganizationSettingDBModel{})

	// Apply filters
	if options.Filters != nil {
		for key, value := range options.Filters {
			switch key {
			case "organization_id":
				query = query.Where("organization_id = ?", value)
			case "category":
				query = query.Where("category = ?", value)
			case "parent_id":
				query = query.Where("parent_id = ?", value)
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
		return nil, fmt.Errorf("%w: failed to list organization settings: %v", shared.ErrInternalService, err)
	}

	organizationSettings := make([]*settings.OrganizationSetting, len(dbModels))
	for i, dbModel := range dbModels {
		organizationSettings[i] = r.toDomainEntity(&dbModel)
	}
	return organizationSettings, nil
}

// Count returns the total number of organization settings matching the filters
func (r *OrganizationSettingRepository) Count(ctx context.Context, filters settings.ListFilters) (int, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&OrganizationSettingDBModel{})

	if filters != nil {
		for key, value := range filters {
			switch key {
			case "organization_id":
				query = query.Where("organization_id = ?", value)
			case "category":
				query = query.Where("category = ?", value)
			case "parent_id":
				query = query.Where("parent_id = ?", value)
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
		return 0, fmt.Errorf("%w: failed to count organization settings: %v", shared.ErrInternalService, err)
	}
	return int(count), nil
}

// FindByOrganization retrieves organization settings by organization ID
func (r *OrganizationSettingRepository) FindByOrganization(ctx context.Context, organizationID settings.OrganizationID) ([]*settings.OrganizationSetting, error) {
	var dbModels []OrganizationSettingDBModel
	err := r.db.WithContext(ctx).Where("organization_id = ?", organizationID.String()).Find(&dbModels).Error
	if err != nil {
		return nil, fmt.Errorf("%w: failed to find organization settings by organization: %v", shared.ErrInternalService, err)
	}

	organizationSettings := make([]*settings.OrganizationSetting, len(dbModels))
	for i, dbModel := range dbModels {
		organizationSettings[i] = r.toDomainEntity(&dbModel)
	}
	return organizationSettings, nil
}

// FindByOrganizationAndCategory retrieves organization settings by organization ID and category
func (r *OrganizationSettingRepository) FindByOrganizationAndCategory(ctx context.Context, organizationID settings.OrganizationID, category string) ([]*settings.OrganizationSetting, error) {
	var dbModels []OrganizationSettingDBModel
	err := r.db.WithContext(ctx).Where("organization_id = ? AND category = ?", organizationID.String(), category).Find(&dbModels).Error
	if err != nil {
		return nil, fmt.Errorf("%w: failed to find organization settings by organization and category: %v", shared.ErrInternalService, err)
	}

	organizationSettings := make([]*settings.OrganizationSetting, len(dbModels))
	for i, dbModel := range dbModels {
		organizationSettings[i] = r.toDomainEntity(&dbModel)
	}
	return organizationSettings, nil
}

// FindByCategory retrieves organization settings by category
func (r *OrganizationSettingRepository) FindByCategory(ctx context.Context, category string) ([]*settings.OrganizationSetting, error) {
	var dbModels []OrganizationSettingDBModel
	err := r.db.WithContext(ctx).Where("category = ?", category).Find(&dbModels).Error
	if err != nil {
		return nil, fmt.Errorf("%w: failed to find organization settings by category: %v", shared.ErrInternalService, err)
	}

	organizationSettings := make([]*settings.OrganizationSetting, len(dbModels))
	for i, dbModel := range dbModels {
		organizationSettings[i] = r.toDomainEntity(&dbModel)
	}
	return organizationSettings, nil
}

// FindByParentOrganization retrieves organization settings by parent organization ID
func (r *OrganizationSettingRepository) FindByParentOrganization(ctx context.Context, parentID settings.OrganizationID) ([]*settings.OrganizationSetting, error) {
	var dbModels []OrganizationSettingDBModel
	err := r.db.WithContext(ctx).Where("parent_id = ?", parentID.String()).Find(&dbModels).Error
	if err != nil {
		return nil, fmt.Errorf("%w: failed to find organization settings by parent organization: %v", shared.ErrInternalService, err)
	}

	organizationSettings := make([]*settings.OrganizationSetting, len(dbModels))
	for i, dbModel := range dbModels {
		organizationSettings[i] = r.toDomainEntity(&dbModel)
	}
	return organizationSettings, nil
}

// ExistsByKey checks if an organization setting with the given key already exists
func (r *OrganizationSettingRepository) ExistsByKey(ctx context.Context, organizationID settings.OrganizationID, category string, key string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&OrganizationSettingDBModel{}).Where("organization_id = ? AND category = ? AND key = ?", organizationID.String(), category, key).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("%w: failed to check organization setting existence: %v", shared.ErrInternalService, err)
	}
	return count > 0, nil
}

// ExistsByID checks if an organization setting with the given ID exists
func (r *OrganizationSettingRepository) ExistsByID(ctx context.Context, id settings.OrganizationSettingID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&OrganizationSettingDBModel{}).Where("id = ?", id.String()).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("%w: failed to check organization setting existence: %v", shared.ErrInternalService, err)
	}
	return count > 0, nil
}

// BulkSave saves multiple organization settings
func (r *OrganizationSettingRepository) BulkSave(ctx context.Context, settings []*settings.OrganizationSetting) error {
	if len(settings) == 0 {
		return nil
	}

	dbModels := make([]OrganizationSettingDBModel, len(settings))
	for i, setting := range settings {
		dbModels[i] = *r.toDBModel(setting)
	}

	err := r.db.WithContext(ctx).CreateInBatches(dbModels, 100).Error
	if err != nil {
		return fmt.Errorf("%w: failed to bulk save organization settings: %v", shared.ErrInternalService, err)
	}
	return nil
}

// BulkUpdate updates multiple organization settings
func (r *OrganizationSettingRepository) BulkUpdate(ctx context.Context, settings []*settings.OrganizationSetting) error {
	if len(settings) == 0 {
		return nil
	}

	for _, setting := range settings {
		if err := r.Update(ctx, setting); err != nil {
			return err
		}
	}
	return nil
}

// BulkDelete deletes multiple organization settings
func (r *OrganizationSettingRepository) BulkDelete(ctx context.Context, ids []settings.OrganizationSettingID) error {
	if len(ids) == 0 {
		return nil
	}

	idStrings := make([]string, len(ids))
	for i, id := range ids {
		idStrings[i] = id.String()
	}

	result := r.db.WithContext(ctx).Delete(&OrganizationSettingDBModel{}, "id IN ?", idStrings)
	if result.Error != nil {
		return fmt.Errorf("%w: failed to bulk delete organization settings: %v", shared.ErrInternalService, result.Error)
	}
	return nil
}

// toDBModel converts a domain organization setting to a database model
func (r *OrganizationSettingRepository) toDBModel(setting *settings.OrganizationSetting) *OrganizationSettingDBModel {
	valueJSON, _ := json.Marshal(setting.GetValue())
	tagsJSON, _ := json.Marshal(setting.GetTags())
	metadataJSON, _ := json.Marshal(setting.GetMetadata())

	dbModel := &OrganizationSettingDBModel{
		ID:             uuid.MustParse(setting.GetID().String()),
		OrganizationID: uuid.MustParse(setting.GetOrganizationID().String()),
		Category:       setting.GetCategory(),
		Key:            setting.GetKey(),
		Value:          JSONB(valueJSON),
		Description:    setting.GetDescription(),
		Tags:           JSONB(tagsJSON),
		Metadata:       JSONB(metadataJSON),
		CreatedBy:      uuid.MustParse(setting.GetCreatedBy().String()),
		CreatedAt:      setting.GetCreatedAt(),
		UpdatedAt:      setting.GetUpdatedAt(),
		Version:        setting.GetVersion(),
	}

	if setting.GetParentID() != nil {
		parentID := uuid.MustParse(setting.GetParentID().String())
		dbModel.ParentID = &parentID
	}

	if setting.GetUpdatedBy() != nil {
		updatedBy := uuid.MustParse(setting.GetUpdatedBy().String())
		dbModel.UpdatedBy = &updatedBy
	}

	return dbModel
}

// toDomainEntity converts a database model to a domain organization setting
func (r *OrganizationSettingRepository) toDomainEntity(dbModel *OrganizationSettingDBModel) *settings.OrganizationSetting {
	settingID, _ := settings.NewOrganizationSettingIDFromString(dbModel.ID.String())
	organizationID, _ := settings.NewOrganizationIDFromString(dbModel.OrganizationID.String())
	createdBy, _ := settings.NewUserIDFromString(dbModel.CreatedBy.String())

	var parentID *settings.OrganizationSettingID
	if dbModel.ParentID != nil {
		parent, _ := settings.NewOrganizationSettingIDFromString(dbModel.ParentID.String())
		parentID = &parent
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

	// Create organization setting (this is a simplified version - in reality you'd need a factory method)
	// For now, we'll create a basic organization setting and then update its fields
	organizationSetting, _ := settings.NewOrganizationSetting(
		organizationID,
		dbModel.Category,
		dbModel.Key,
		value,
		parentID,
		dbModel.Description,
		tags,
		metadata,
		createdBy,
	)

	// Note: In a real implementation, you would need to properly restore the organization setting state
	// including the ID, version, and timestamps. This is a simplified version.

	return organizationSetting
}
