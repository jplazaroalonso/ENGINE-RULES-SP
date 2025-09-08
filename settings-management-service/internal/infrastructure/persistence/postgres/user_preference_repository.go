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

// UserPreferenceDBModel represents the user preference database model
type UserPreferenceDBModel struct {
	ID             uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID         uuid.UUID  `gorm:"type:uuid;not null;index:idx_user_pref_user_org_cat_key,unique"`
	Category       string     `gorm:"type:varchar(100);not null;index:idx_user_pref_user_org_cat_key"`
	Key            string     `gorm:"type:varchar(255);not null;index:idx_user_pref_user_org_cat_key"`
	Value          JSONB      `gorm:"type:jsonb;not null"`
	OrganizationID *uuid.UUID `gorm:"type:uuid;index:idx_user_pref_user_org_cat_key"`
	Description    *string    `gorm:"type:text"`
	Tags           JSONB      `gorm:"type:jsonb;not null;default:'[]'"`
	Metadata       JSONB      `gorm:"type:jsonb;not null;default:'{}'"`
	CreatedBy      uuid.UUID  `gorm:"type:uuid;not null"`
	UpdatedBy      *uuid.UUID `gorm:"type:uuid"`
	CreatedAt      time.Time  `gorm:"type:timestamp with time zone;default:now()"`
	UpdatedAt      time.Time  `gorm:"type:timestamp with time zone;default:now()"`
	Version        int        `gorm:"type:integer;not null;default:1"`
}

// TableName specifies the table name for UserPreferenceDBModel
func (UserPreferenceDBModel) TableName() string {
	return "user_preferences"
}

// UserPreferenceRepository implements the settings.UserPreferenceRepository interface
type UserPreferenceRepository struct {
	db *gorm.DB
}

// NewUserPreferenceRepository creates a new UserPreferenceRepository
func NewUserPreferenceRepository(db *gorm.DB) *UserPreferenceRepository {
	return &UserPreferenceRepository{db: db}
}

// Save persists a user preference aggregate
func (r *UserPreferenceRepository) Save(ctx context.Context, preference *settings.UserPreference) error {
	dbModel := r.toDBModel(preference)

	// Use Upsert to handle both creation and updates
	err := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"value", "description", "tags", "metadata", "updated_by", "updated_at", "version"}),
	}).Create(&dbModel).Error

	if err != nil {
		return fmt.Errorf("%w: failed to save user preference: %v", shared.ErrInternalService, err)
	}
	return nil
}

// FindByID retrieves a user preference by its ID
func (r *UserPreferenceRepository) FindByID(ctx context.Context, id settings.UserPreferenceID) (*settings.UserPreference, error) {
	var dbModel UserPreferenceDBModel
	err := r.db.WithContext(ctx).First(&dbModel, "id = ?", id.String()).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: user preference with ID %s not found", shared.ErrNotFound, id.String())
		}
		return nil, fmt.Errorf("%w: failed to find user preference by ID: %v", shared.ErrInternalService, err)
	}
	return r.toDomainEntity(&dbModel), nil
}

// FindByKey retrieves a user preference by its user ID, category, key, and organization
func (r *UserPreferenceRepository) FindByKey(ctx context.Context, userID settings.UserID, category string, key string, organizationID *settings.OrganizationID) (*settings.UserPreference, error) {
	var dbModel UserPreferenceDBModel
	query := r.db.WithContext(ctx).Where("user_id = ? AND category = ? AND key = ?", userID.String(), category, key)

	if organizationID != nil {
		query = query.Where("organization_id = ?", organizationID.String())
	} else {
		query = query.Where("organization_id IS NULL")
	}

	err := query.First(&dbModel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: user preference with key '%s' not found for user '%s' in category '%s'", shared.ErrNotFound, key, userID.String(), category)
		}
		return nil, fmt.Errorf("%w: failed to find user preference by key: %v", shared.ErrInternalService, err)
	}
	return r.toDomainEntity(&dbModel), nil
}

// Update updates an existing user preference
func (r *UserPreferenceRepository) Update(ctx context.Context, preference *settings.UserPreference) error {
	dbModel := r.toDBModel(preference)

	result := r.db.WithContext(ctx).Model(&UserPreferenceDBModel{}).Where("id = ?", dbModel.ID).Updates(map[string]interface{}{
		"value":       dbModel.Value,
		"description": dbModel.Description,
		"tags":        dbModel.Tags,
		"metadata":    dbModel.Metadata,
		"updated_by":  dbModel.UpdatedBy,
		"updated_at":  time.Now().UTC(),
		"version":     gorm.Expr("version + ?", 1),
	})

	if result.Error != nil {
		return fmt.Errorf("%w: failed to update user preference: %v", shared.ErrInternalService, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("%w: user preference with ID %s not found for update", shared.ErrNotFound, preference.GetID().String())
	}
	return nil
}

// Delete removes a user preference by its ID
func (r *UserPreferenceRepository) Delete(ctx context.Context, id settings.UserPreferenceID) error {
	result := r.db.WithContext(ctx).Delete(&UserPreferenceDBModel{}, "id = ?", id.String())
	if result.Error != nil {
		return fmt.Errorf("%w: failed to delete user preference: %v", shared.ErrInternalService, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("%w: user preference with ID %s not found for deletion", shared.ErrNotFound, id.String())
	}
	return nil
}

// List retrieves a list of user preferences with pagination and filtering
func (r *UserPreferenceRepository) List(ctx context.Context, options settings.ListOptions) ([]*settings.UserPreference, error) {
	var dbModels []UserPreferenceDBModel
	query := r.db.WithContext(ctx).Model(&UserPreferenceDBModel{})

	// Apply filters
	if options.Filters != nil {
		for key, value := range options.Filters {
			switch key {
			case "user_id":
				query = query.Where("user_id = ?", value)
			case "organization_id":
				query = query.Where("organization_id = ?", value)
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
		return nil, fmt.Errorf("%w: failed to list user preferences: %v", shared.ErrInternalService, err)
	}

	userPreferences := make([]*settings.UserPreference, len(dbModels))
	for i, dbModel := range dbModels {
		userPreferences[i] = r.toDomainEntity(&dbModel)
	}
	return userPreferences, nil
}

// Count returns the total number of user preferences matching the filters
func (r *UserPreferenceRepository) Count(ctx context.Context, filters settings.ListFilters) (int, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&UserPreferenceDBModel{})

	if filters != nil {
		for key, value := range filters {
			switch key {
			case "user_id":
				query = query.Where("user_id = ?", value)
			case "organization_id":
				query = query.Where("organization_id = ?", value)
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
		return 0, fmt.Errorf("%w: failed to count user preferences: %v", shared.ErrInternalService, err)
	}
	return int(count), nil
}

// FindByUser retrieves user preferences by user ID
func (r *UserPreferenceRepository) FindByUser(ctx context.Context, userID settings.UserID, organizationID *settings.OrganizationID) ([]*settings.UserPreference, error) {
	var dbModels []UserPreferenceDBModel
	query := r.db.WithContext(ctx).Where("user_id = ?", userID.String())

	if organizationID != nil {
		query = query.Where("organization_id = ?", organizationID.String())
	} else {
		query = query.Where("organization_id IS NULL")
	}

	err := query.Find(&dbModels).Error
	if err != nil {
		return nil, fmt.Errorf("%w: failed to find user preferences by user: %v", shared.ErrInternalService, err)
	}

	userPreferences := make([]*settings.UserPreference, len(dbModels))
	for i, dbModel := range dbModels {
		userPreferences[i] = r.toDomainEntity(&dbModel)
	}
	return userPreferences, nil
}

// FindByUserAndCategory retrieves user preferences by user ID and category
func (r *UserPreferenceRepository) FindByUserAndCategory(ctx context.Context, userID settings.UserID, category string, organizationID *settings.OrganizationID) ([]*settings.UserPreference, error) {
	var dbModels []UserPreferenceDBModel
	query := r.db.WithContext(ctx).Where("user_id = ? AND category = ?", userID.String(), category)

	if organizationID != nil {
		query = query.Where("organization_id = ?", organizationID.String())
	} else {
		query = query.Where("organization_id IS NULL")
	}

	err := query.Find(&dbModels).Error
	if err != nil {
		return nil, fmt.Errorf("%w: failed to find user preferences by user and category: %v", shared.ErrInternalService, err)
	}

	userPreferences := make([]*settings.UserPreference, len(dbModels))
	for i, dbModel := range dbModels {
		userPreferences[i] = r.toDomainEntity(&dbModel)
	}
	return userPreferences, nil
}

// FindByOrganization retrieves user preferences by organization
func (r *UserPreferenceRepository) FindByOrganization(ctx context.Context, organizationID settings.OrganizationID) ([]*settings.UserPreference, error) {
	var dbModels []UserPreferenceDBModel
	err := r.db.WithContext(ctx).Where("organization_id = ?", organizationID.String()).Find(&dbModels).Error
	if err != nil {
		return nil, fmt.Errorf("%w: failed to find user preferences by organization: %v", shared.ErrInternalService, err)
	}

	userPreferences := make([]*settings.UserPreference, len(dbModels))
	for i, dbModel := range dbModels {
		userPreferences[i] = r.toDomainEntity(&dbModel)
	}
	return userPreferences, nil
}

// FindByCategory retrieves user preferences by category
func (r *UserPreferenceRepository) FindByCategory(ctx context.Context, category string, organizationID *settings.OrganizationID) ([]*settings.UserPreference, error) {
	var dbModels []UserPreferenceDBModel
	query := r.db.WithContext(ctx).Where("category = ?", category)

	if organizationID != nil {
		query = query.Where("organization_id = ?", organizationID.String())
	} else {
		query = query.Where("organization_id IS NULL")
	}

	err := query.Find(&dbModels).Error
	if err != nil {
		return nil, fmt.Errorf("%w: failed to find user preferences by category: %v", shared.ErrInternalService, err)
	}

	userPreferences := make([]*settings.UserPreference, len(dbModels))
	for i, dbModel := range dbModels {
		userPreferences[i] = r.toDomainEntity(&dbModel)
	}
	return userPreferences, nil
}

// ExistsByKey checks if a user preference with the given key already exists
func (r *UserPreferenceRepository) ExistsByKey(ctx context.Context, userID settings.UserID, category string, key string, organizationID *settings.OrganizationID) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&UserPreferenceDBModel{}).Where("user_id = ? AND category = ? AND key = ?", userID.String(), category, key)

	if organizationID != nil {
		query = query.Where("organization_id = ?", organizationID.String())
	} else {
		query = query.Where("organization_id IS NULL")
	}

	err := query.Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("%w: failed to check user preference existence: %v", shared.ErrInternalService, err)
	}
	return count > 0, nil
}

// ExistsByID checks if a user preference with the given ID exists
func (r *UserPreferenceRepository) ExistsByID(ctx context.Context, id settings.UserPreferenceID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&UserPreferenceDBModel{}).Where("id = ?", id.String()).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("%w: failed to check user preference existence: %v", shared.ErrInternalService, err)
	}
	return count > 0, nil
}

// BulkSave saves multiple user preferences
func (r *UserPreferenceRepository) BulkSave(ctx context.Context, preferences []*settings.UserPreference) error {
	if len(preferences) == 0 {
		return nil
	}

	dbModels := make([]UserPreferenceDBModel, len(preferences))
	for i, preference := range preferences {
		dbModels[i] = *r.toDBModel(preference)
	}

	err := r.db.WithContext(ctx).CreateInBatches(dbModels, 100).Error
	if err != nil {
		return fmt.Errorf("%w: failed to bulk save user preferences: %v", shared.ErrInternalService, err)
	}
	return nil
}

// BulkUpdate updates multiple user preferences
func (r *UserPreferenceRepository) BulkUpdate(ctx context.Context, preferences []*settings.UserPreference) error {
	if len(preferences) == 0 {
		return nil
	}

	for _, preference := range preferences {
		if err := r.Update(ctx, preference); err != nil {
			return err
		}
	}
	return nil
}

// BulkDelete deletes multiple user preferences
func (r *UserPreferenceRepository) BulkDelete(ctx context.Context, ids []settings.UserPreferenceID) error {
	if len(ids) == 0 {
		return nil
	}

	idStrings := make([]string, len(ids))
	for i, id := range ids {
		idStrings[i] = id.String()
	}

	result := r.db.WithContext(ctx).Delete(&UserPreferenceDBModel{}, "id IN ?", idStrings)
	if result.Error != nil {
		return fmt.Errorf("%w: failed to bulk delete user preferences: %v", shared.ErrInternalService, result.Error)
	}
	return nil
}

// toDBModel converts a domain user preference to a database model
func (r *UserPreferenceRepository) toDBModel(preference *settings.UserPreference) *UserPreferenceDBModel {
	valueJSON, _ := json.Marshal(preference.GetValue())
	tagsJSON, _ := json.Marshal(preference.GetTags())
	metadataJSON, _ := json.Marshal(preference.GetMetadata())

	dbModel := &UserPreferenceDBModel{
		ID:          uuid.MustParse(preference.GetID().String()),
		UserID:      uuid.MustParse(preference.GetUserID().String()),
		Category:    preference.GetCategory(),
		Key:         preference.GetKey(),
		Value:       JSONB(valueJSON),
		Description: preference.GetDescription(),
		Tags:        JSONB(tagsJSON),
		Metadata:    JSONB(metadataJSON),
		CreatedBy:   uuid.MustParse(preference.GetCreatedBy().String()),
		CreatedAt:   preference.GetCreatedAt(),
		UpdatedAt:   preference.GetUpdatedAt(),
		Version:     preference.GetVersion(),
	}

	if preference.GetOrganizationID() != nil {
		orgID := uuid.MustParse(preference.GetOrganizationID().String())
		dbModel.OrganizationID = &orgID
	}

	if preference.GetUpdatedBy() != nil {
		updatedBy := uuid.MustParse(preference.GetUpdatedBy().String())
		dbModel.UpdatedBy = &updatedBy
	}

	return dbModel
}

// toDomainEntity converts a database model to a domain user preference
func (r *UserPreferenceRepository) toDomainEntity(dbModel *UserPreferenceDBModel) *settings.UserPreference {
	preferenceID, _ := settings.NewUserPreferenceIDFromString(dbModel.ID.String())
	userID, _ := settings.NewUserIDFromString(dbModel.UserID.String())
	createdBy, _ := settings.NewUserIDFromString(dbModel.CreatedBy.String())

	var organizationID *settings.OrganizationID
	if dbModel.OrganizationID != nil {
		orgID, _ := settings.NewOrganizationIDFromString(dbModel.OrganizationID.String())
		organizationID = &orgID
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

	// Create user preference (this is a simplified version - in reality you'd need a factory method)
	// For now, we'll create a basic user preference and then update its fields
	userPreference, _ := settings.NewUserPreference(
		userID,
		dbModel.Category,
		dbModel.Key,
		value,
		organizationID,
		dbModel.Description,
		tags,
		metadata,
		createdBy,
	)

	// Note: In a real implementation, you would need to properly restore the user preference state
	// including the ID, version, and timestamps. This is a simplified version.

	return userPreference
}
