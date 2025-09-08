package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/analytics-dashboard-service/internal/domain/analytics"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/analytics-dashboard-service/internal/domain/shared"
	"gorm.io/gorm"
)

// DashboardModel represents the database model for Dashboard
type DashboardModel struct {
	ID              string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name            string `gorm:"not null"`
	Description     string
	Layout          json.RawMessage `gorm:"type:jsonb;not null;default:'{}'"`
	Widgets         json.RawMessage `gorm:"type:jsonb;not null;default:'[]'"`
	Filters         json.RawMessage `gorm:"type:jsonb;not null;default:'{}'"`
	RefreshInterval int             `gorm:"not null;default:300"`
	IsPublic        bool            `gorm:"not null;default:false"`
	OwnerID         string          `gorm:"type:uuid;not null"`
	CreatedAt       time.Time       `gorm:"autoCreateTime"`
	UpdatedAt       time.Time       `gorm:"autoUpdateTime"`
	Version         int             `gorm:"not null;default:1"`
}

// TableName returns the table name for DashboardModel
func (DashboardModel) TableName() string {
	return "dashboards"
}

// DashboardRepository implements the DashboardRepository interface
type DashboardRepository struct {
	db *gorm.DB
}

// NewDashboardRepository creates a new dashboard repository
func NewDashboardRepository(db *gorm.DB) *DashboardRepository {
	return &DashboardRepository{db: db}
}

// Save saves a dashboard to the database
func (r *DashboardRepository) Save(ctx context.Context, dashboard *analytics.Dashboard) error {
	model, err := r.domainToModel(dashboard)
	if err != nil {
		return fmt.Errorf("failed to convert domain to model: %w", err)
	}

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("failed to save dashboard: %w", err)
	}

	// Update the domain object with the generated ID
	dashboard.ID = shared.DashboardID(model.ID)
	return nil
}

// FindByID finds a dashboard by ID
func (r *DashboardRepository) FindByID(ctx context.Context, id shared.DashboardID) (*analytics.Dashboard, error) {
	var model DashboardModel
	if err := r.db.WithContext(ctx).Where("id = ?", id.String()).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, shared.ErrNotFound
		}
		return nil, fmt.Errorf("failed to find dashboard: %w", err)
	}

	return r.modelToDomain(&model)
}

// FindByOwnerID finds dashboards by owner ID
func (r *DashboardRepository) FindByOwnerID(ctx context.Context, ownerID shared.UserID) ([]*analytics.Dashboard, error) {
	var models []DashboardModel
	if err := r.db.WithContext(ctx).Where("owner_id = ?", ownerID.String()).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to find dashboards by owner: %w", err)
	}

	dashboards := make([]*analytics.Dashboard, len(models))
	for i, model := range models {
		dashboard, err := r.modelToDomain(&model)
		if err != nil {
			return nil, fmt.Errorf("failed to convert model to domain: %w", err)
		}
		dashboards[i] = dashboard
	}

	return dashboards, nil
}

// FindPublic finds public dashboards
func (r *DashboardRepository) FindPublic(ctx context.Context) ([]*analytics.Dashboard, error) {
	var models []DashboardModel
	if err := r.db.WithContext(ctx).Where("is_public = ?", true).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to find public dashboards: %w", err)
	}

	dashboards := make([]*analytics.Dashboard, len(models))
	for i, model := range models {
		dashboard, err := r.modelToDomain(&model)
		if err != nil {
			return nil, fmt.Errorf("failed to convert model to domain: %w", err)
		}
		dashboards[i] = dashboard
	}

	return dashboards, nil
}

// Delete deletes a dashboard by ID
func (r *DashboardRepository) Delete(ctx context.Context, id shared.DashboardID) error {
	if err := r.db.WithContext(ctx).Where("id = ?", id.String()).Delete(&DashboardModel{}).Error; err != nil {
		return fmt.Errorf("failed to delete dashboard: %w", err)
	}
	return nil
}

// Update updates a dashboard
func (r *DashboardRepository) Update(ctx context.Context, dashboard *analytics.Dashboard) error {
	model, err := r.domainToModel(dashboard)
	if err != nil {
		return fmt.Errorf("failed to convert domain to model: %w", err)
	}

	if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
		return fmt.Errorf("failed to update dashboard: %w", err)
	}

	return nil
}

// domainToModel converts a domain Dashboard to a database model
func (r *DashboardRepository) domainToModel(dashboard *analytics.Dashboard) (*DashboardModel, error) {
	layoutBytes, err := json.Marshal(dashboard.Layout)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal layout: %w", err)
	}

	widgetsBytes, err := json.Marshal(dashboard.Widgets)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal widgets: %w", err)
	}

	filtersBytes, err := json.Marshal(dashboard.Filters)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal filters: %w", err)
	}

	return &DashboardModel{
		ID:              dashboard.ID.String(),
		Name:            dashboard.Name,
		Description:     dashboard.Description,
		Layout:          layoutBytes,
		Widgets:         widgetsBytes,
		Filters:         filtersBytes,
		RefreshInterval: dashboard.RefreshInterval,
		IsPublic:        dashboard.IsPublic,
		OwnerID:         dashboard.OwnerID.String(),
		Version:         dashboard.Version,
	}, nil
}

// modelToDomain converts a database model to a domain Dashboard
func (r *DashboardRepository) modelToDomain(model *DashboardModel) (*analytics.Dashboard, error) {
	var layout analytics.DashboardLayout
	if err := json.Unmarshal(model.Layout, &layout); err != nil {
		return nil, fmt.Errorf("failed to unmarshal layout: %w", err)
	}

	var widgets []analytics.Widget
	if err := json.Unmarshal(model.Widgets, &widgets); err != nil {
		return nil, fmt.Errorf("failed to unmarshal widgets: %w", err)
	}

	var filters analytics.DashboardFilters
	if err := json.Unmarshal(model.Filters, &filters); err != nil {
		return nil, fmt.Errorf("failed to unmarshal filters: %w", err)
	}

	return &analytics.Dashboard{
		ID:              shared.DashboardID(model.ID),
		Name:            model.Name,
		Description:     model.Description,
		Layout:          layout,
		Widgets:         widgets,
		Filters:         filters,
		RefreshInterval: model.RefreshInterval,
		IsPublic:        model.IsPublic,
		OwnerID:         shared.UserID(model.OwnerID),
		CreatedAt:       model.CreatedAt,
		UpdatedAt:       model.UpdatedAt,
		Version:         model.Version,
		Events:          []*shared.DomainEvent{},
	}, nil
}
