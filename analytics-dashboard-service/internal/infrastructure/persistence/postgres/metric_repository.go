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

// MetricModel represents the database model for Metric
type MetricModel struct {
	ID           string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name         string `gorm:"not null;unique"`
	Description  string
	Type         string `gorm:"not null"`
	Category     string `gorm:"not null"`
	Unit         string
	Aggregation  string          `gorm:"not null"`
	DataSource   json.RawMessage `gorm:"type:jsonb;not null"`
	Dimensions   json.RawMessage `gorm:"type:jsonb;not null;default:'[]'"`
	Filters      json.RawMessage `gorm:"type:jsonb;not null;default:'{}'"`
	Calculation  json.RawMessage `gorm:"type:jsonb"`
	IsCalculated bool            `gorm:"not null;default:false"`
	CreatedAt    time.Time       `gorm:"autoCreateTime"`
	UpdatedAt    time.Time       `gorm:"autoUpdateTime"`
	Version      int             `gorm:"not null;default:1"`
}

// TableName returns the table name for MetricModel
func (MetricModel) TableName() string {
	return "metrics"
}

// MetricRepository implements the MetricRepository interface
type MetricRepository struct {
	db *gorm.DB
}

// NewMetricRepository creates a new metric repository
func NewMetricRepository(db *gorm.DB) *MetricRepository {
	return &MetricRepository{db: db}
}

// Save saves a metric to the database
func (r *MetricRepository) Save(ctx context.Context, metric *analytics.Metric) error {
	model, err := r.domainToModel(metric)
	if err != nil {
		return fmt.Errorf("failed to convert domain to model: %w", err)
	}

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("failed to save metric: %w", err)
	}

	// Update the domain object with the generated ID
	metric.ID = shared.MetricID(model.ID)
	return nil
}

// FindByID finds a metric by ID
func (r *MetricRepository) FindByID(ctx context.Context, id shared.MetricID) (*analytics.Metric, error) {
	var model MetricModel
	if err := r.db.WithContext(ctx).Where("id = ?", id.String()).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, shared.ErrNotFound
		}
		return nil, fmt.Errorf("failed to find metric: %w", err)
	}

	return r.modelToDomain(&model)
}

// FindByCategory finds metrics by category
func (r *MetricRepository) FindByCategory(ctx context.Context, category analytics.MetricCategory) ([]*analytics.Metric, error) {
	var models []MetricModel
	if err := r.db.WithContext(ctx).Where("category = ?", string(category)).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to find metrics by category: %w", err)
	}

	metrics := make([]*analytics.Metric, len(models))
	for i, model := range models {
		metric, err := r.modelToDomain(&model)
		if err != nil {
			return nil, fmt.Errorf("failed to convert model to domain: %w", err)
		}
		metrics[i] = metric
	}

	return metrics, nil
}

// FindByType finds metrics by type
func (r *MetricRepository) FindByType(ctx context.Context, metricType analytics.MetricType) ([]*analytics.Metric, error) {
	var models []MetricModel
	if err := r.db.WithContext(ctx).Where("type = ?", string(metricType)).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to find metrics by type: %w", err)
	}

	metrics := make([]*analytics.Metric, len(models))
	for i, model := range models {
		metric, err := r.modelToDomain(&model)
		if err != nil {
			return nil, fmt.Errorf("failed to convert model to domain: %w", err)
		}
		metrics[i] = metric
	}

	return metrics, nil
}

// FindAll finds all metrics
func (r *MetricRepository) FindAll(ctx context.Context) ([]*analytics.Metric, error) {
	var models []MetricModel
	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to find all metrics: %w", err)
	}

	metrics := make([]*analytics.Metric, len(models))
	for i, model := range models {
		metric, err := r.modelToDomain(&model)
		if err != nil {
			return nil, fmt.Errorf("failed to convert model to domain: %w", err)
		}
		metrics[i] = metric
	}

	return metrics, nil
}

// Delete deletes a metric by ID
func (r *MetricRepository) Delete(ctx context.Context, id shared.MetricID) error {
	if err := r.db.WithContext(ctx).Where("id = ?", id.String()).Delete(&MetricModel{}).Error; err != nil {
		return fmt.Errorf("failed to delete metric: %w", err)
	}
	return nil
}

// Update updates a metric
func (r *MetricRepository) Update(ctx context.Context, metric *analytics.Metric) error {
	model, err := r.domainToModel(metric)
	if err != nil {
		return fmt.Errorf("failed to convert domain to model: %w", err)
	}

	if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
		return fmt.Errorf("failed to update metric: %w", err)
	}

	return nil
}

// domainToModel converts a domain Metric to a database model
func (r *MetricRepository) domainToModel(metric *analytics.Metric) (*MetricModel, error) {
	dataSourceBytes, err := json.Marshal(metric.DataSource)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data source: %w", err)
	}

	dimensionsBytes, err := json.Marshal(metric.Dimensions)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal dimensions: %w", err)
	}

	filtersBytes, err := json.Marshal(metric.Filters)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal filters: %w", err)
	}

	var calculationBytes json.RawMessage
	if metric.Calculation.Formula != "" {
		calculationBytes, err = json.Marshal(metric.Calculation)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal calculation: %w", err)
		}
	}

	return &MetricModel{
		ID:           metric.ID.String(),
		Name:         metric.Name,
		Description:  metric.Description,
		Type:         string(metric.Type),
		Category:     string(metric.Category),
		Unit:         metric.Unit,
		Aggregation:  string(metric.Aggregation),
		DataSource:   dataSourceBytes,
		Dimensions:   dimensionsBytes,
		Filters:      filtersBytes,
		Calculation:  calculationBytes,
		IsCalculated: metric.IsCalculated,
		Version:      metric.Version,
	}, nil
}

// modelToDomain converts a database model to a domain Metric
func (r *MetricRepository) modelToDomain(model *MetricModel) (*analytics.Metric, error) {
	var dataSource analytics.DataSource
	if err := json.Unmarshal(model.DataSource, &dataSource); err != nil {
		return nil, fmt.Errorf("failed to unmarshal data source: %w", err)
	}

	var dimensions []analytics.Dimension
	if err := json.Unmarshal(model.Dimensions, &dimensions); err != nil {
		return nil, fmt.Errorf("failed to unmarshal dimensions: %w", err)
	}

	var filters analytics.MetricFilters
	if err := json.Unmarshal(model.Filters, &filters); err != nil {
		return nil, fmt.Errorf("failed to unmarshal filters: %w", err)
	}

	var calculation analytics.MetricCalculation
	if len(model.Calculation) > 0 {
		if err := json.Unmarshal(model.Calculation, &calculation); err != nil {
			return nil, fmt.Errorf("failed to unmarshal calculation: %w", err)
		}
	}

	return &analytics.Metric{
		ID:           shared.MetricID(model.ID),
		Name:         model.Name,
		Description:  model.Description,
		Type:         analytics.MetricType(model.Type),
		Category:     analytics.MetricCategory(model.Category),
		Unit:         model.Unit,
		Aggregation:  analytics.AggregationType(model.Aggregation),
		DataSource:   dataSource,
		Dimensions:   dimensions,
		Filters:      filters,
		Calculation:  calculation,
		IsCalculated: model.IsCalculated,
		CreatedAt:    model.CreatedAt,
		UpdatedAt:    model.UpdatedAt,
		Version:      model.Version,
		Events:       []*shared.DomainEvent{},
	}, nil
}
