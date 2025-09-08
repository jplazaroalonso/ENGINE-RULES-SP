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

// MetricDataModel represents the database model for MetricData
type MetricDataModel struct {
	ID         string          `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	MetricID   string          `gorm:"type:uuid;not null"`
	Timestamp  time.Time       `gorm:"not null"`
	Value      float64         `gorm:"type:decimal(20,6);not null"`
	Dimensions json.RawMessage `gorm:"type:jsonb;not null;default:'{}'"`
	Labels     json.RawMessage `gorm:"type:jsonb;not null;default:'{}'"`
}

// TableName returns the table name for MetricDataModel
func (MetricDataModel) TableName() string {
	return "metric_data"
}

// MetricDataRepository implements the MetricDataRepository interface
type MetricDataRepository struct {
	db *gorm.DB
}

// NewMetricDataRepository creates a new metric data repository
func NewMetricDataRepository(db *gorm.DB) *MetricDataRepository {
	return &MetricDataRepository{db: db}
}

// Save saves metric data to the database
func (r *MetricDataRepository) Save(ctx context.Context, data *analytics.MetricData) error {
	model, err := r.domainToModel(data)
	if err != nil {
		return fmt.Errorf("failed to convert domain to model: %w", err)
	}

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("failed to save metric data: %w", err)
	}

	// Update the domain object with the generated ID
	data.ID = model.ID
	return nil
}

// SaveBatch saves multiple metric data points in a batch
func (r *MetricDataRepository) SaveBatch(ctx context.Context, data []*analytics.MetricData) error {
	if len(data) == 0 {
		return nil
	}

	models := make([]MetricDataModel, len(data))
	for i, d := range data {
		model, err := r.domainToModel(d)
		if err != nil {
			return fmt.Errorf("failed to convert domain to model: %w", err)
		}
		models[i] = *model
	}

	if err := r.db.WithContext(ctx).CreateInBatches(models, 1000).Error; err != nil {
		return fmt.Errorf("failed to save metric data batch: %w", err)
	}

	// Update the domain objects with the generated IDs
	for i, model := range models {
		data[i].ID = model.ID
	}

	return nil
}

// FindByMetricID finds metric data by metric ID within a time range
func (r *MetricDataRepository) FindByMetricID(ctx context.Context, metricID shared.MetricID, timeRange *shared.TimeRange) ([]*analytics.MetricData, error) {
	var models []MetricDataModel
	query := r.db.WithContext(ctx).Where("metric_id = ?", metricID.String())

	if timeRange != nil {
		query = query.Where("timestamp >= ? AND timestamp <= ?", timeRange.Start, timeRange.End)
	}

	if err := query.Order("timestamp ASC").Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to find metric data: %w", err)
	}

	data := make([]*analytics.MetricData, len(models))
	for i, model := range models {
		metricData, err := r.modelToDomain(&model)
		if err != nil {
			return nil, fmt.Errorf("failed to convert model to domain: %w", err)
		}
		data[i] = metricData
	}

	return data, nil
}

// FindByMetricIDAndDimensions finds metric data by metric ID and dimensions within a time range
func (r *MetricDataRepository) FindByMetricIDAndDimensions(ctx context.Context, metricID shared.MetricID, dimensions map[string]interface{}, timeRange *shared.TimeRange) ([]*analytics.MetricData, error) {
	var models []MetricDataModel
	query := r.db.WithContext(ctx).Where("metric_id = ?", metricID.String())

	if timeRange != nil {
		query = query.Where("timestamp >= ? AND timestamp <= ?", timeRange.Start, timeRange.End)
	}

	// Add dimension filters
	for key, value := range dimensions {
		query = query.Where("dimensions->>? = ?", key, value)
	}

	if err := query.Order("timestamp ASC").Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to find metric data by dimensions: %w", err)
	}

	data := make([]*analytics.MetricData, len(models))
	for i, model := range models {
		metricData, err := r.modelToDomain(&model)
		if err != nil {
			return nil, fmt.Errorf("failed to convert model to domain: %w", err)
		}
		data[i] = metricData
	}

	return data, nil
}

// AggregateByMetricID aggregates metric data by metric ID
func (r *MetricDataRepository) AggregateByMetricID(ctx context.Context, metricID shared.MetricID, aggregation analytics.AggregationType, timeRange *shared.TimeRange) (*analytics.MetricData, error) {
	var result struct {
		Value     float64   `json:"value"`
		Timestamp time.Time `json:"timestamp"`
	}

	query := r.db.WithContext(ctx).Model(&MetricDataModel{}).Where("metric_id = ?", metricID.String())

	if timeRange != nil {
		query = query.Where("timestamp >= ? AND timestamp <= ?", timeRange.Start, timeRange.End)
	}

	switch aggregation {
	case analytics.AggregationTypeSum:
		if err := query.Select("SUM(value) as value, MAX(timestamp) as timestamp").Scan(&result).Error; err != nil {
			return nil, fmt.Errorf("failed to aggregate sum: %w", err)
		}
	case analytics.AggregationTypeAvg:
		if err := query.Select("AVG(value) as value, MAX(timestamp) as timestamp").Scan(&result).Error; err != nil {
			return nil, fmt.Errorf("failed to aggregate avg: %w", err)
		}
	case analytics.AggregationTypeMin:
		if err := query.Select("MIN(value) as value, MAX(timestamp) as timestamp").Scan(&result).Error; err != nil {
			return nil, fmt.Errorf("failed to aggregate min: %w", err)
		}
	case analytics.AggregationTypeMax:
		if err := query.Select("MAX(value) as value, MAX(timestamp) as timestamp").Scan(&result).Error; err != nil {
			return nil, fmt.Errorf("failed to aggregate max: %w", err)
		}
	case analytics.AggregationTypeCount:
		if err := query.Select("COUNT(*) as value, MAX(timestamp) as timestamp").Scan(&result).Error; err != nil {
			return nil, fmt.Errorf("failed to aggregate count: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported aggregation type: %s", aggregation)
	}

	return &analytics.MetricData{
		ID:         fmt.Sprintf("aggregated-%s", metricID.String()),
		MetricID:   metricID,
		Timestamp:  result.Timestamp,
		Value:      result.Value,
		Dimensions: make(map[string]interface{}),
		Labels:     make(map[string]string),
	}, nil
}

// DeleteOldData deletes metric data older than the specified time
func (r *MetricDataRepository) DeleteOldData(ctx context.Context, before time.Time) error {
	if err := r.db.WithContext(ctx).Where("timestamp < ?", before).Delete(&MetricDataModel{}).Error; err != nil {
		return fmt.Errorf("failed to delete old metric data: %w", err)
	}
	return nil
}

// domainToModel converts a domain MetricData to a database model
func (r *MetricDataRepository) domainToModel(data *analytics.MetricData) (*MetricDataModel, error) {
	dimensionsBytes, err := json.Marshal(data.Dimensions)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal dimensions: %w", err)
	}

	labelsBytes, err := json.Marshal(data.Labels)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal labels: %w", err)
	}

	return &MetricDataModel{
		ID:         data.ID,
		MetricID:   data.MetricID.String(),
		Timestamp:  data.Timestamp,
		Value:      data.Value,
		Dimensions: dimensionsBytes,
		Labels:     labelsBytes,
	}, nil
}

// modelToDomain converts a database model to a domain MetricData
func (r *MetricDataRepository) modelToDomain(model *MetricDataModel) (*analytics.MetricData, error) {
	var dimensions map[string]interface{}
	if err := json.Unmarshal(model.Dimensions, &dimensions); err != nil {
		return nil, fmt.Errorf("failed to unmarshal dimensions: %w", err)
	}

	var labels map[string]string
	if err := json.Unmarshal(model.Labels, &labels); err != nil {
		return nil, fmt.Errorf("failed to unmarshal labels: %w", err)
	}

	return &analytics.MetricData{
		ID:         model.ID,
		MetricID:   shared.MetricID(model.MetricID),
		Timestamp:  model.Timestamp,
		Value:      model.Value,
		Dimensions: dimensions,
		Labels:     labels,
	}, nil
}
