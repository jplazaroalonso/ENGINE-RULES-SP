package analytics

import (
	"time"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/analytics-dashboard-service/internal/domain/shared"
)

// Metric represents a metric aggregate
type Metric struct {
	ID           shared.MetricID       `json:"id"`
	Name         string                `json:"name"`
	Description  string                `json:"description"`
	Type         MetricType            `json:"type"`
	Category     MetricCategory        `json:"category"`
	Unit         string                `json:"unit"`
	Aggregation  AggregationType       `json:"aggregation"`
	DataSource   DataSource            `json:"dataSource"`
	Dimensions   []Dimension           `json:"dimensions"`
	Filters      MetricFilters         `json:"filters"`
	Calculation  MetricCalculation     `json:"calculation"`
	IsCalculated bool                  `json:"isCalculated"`
	CreatedAt    time.Time             `json:"createdAt"`
	UpdatedAt    time.Time             `json:"updatedAt"`
	Version      int                   `json:"version"`
	Events       []*shared.DomainEvent `json:"-"`
}

// MetricType represents the type of metric
type MetricType string

const (
	MetricTypeCounter   MetricType = "COUNTER"
	MetricTypeGauge     MetricType = "GAUGE"
	MetricTypeHistogram MetricType = "HISTOGRAM"
	MetricTypeSummary   MetricType = "SUMMARY"
)

// MetricCategory represents the category of metric
type MetricCategory string

const (
	MetricCategoryPerformance MetricCategory = "PERFORMANCE"
	MetricCategoryBusiness    MetricCategory = "BUSINESS"
	MetricCategorySystem      MetricCategory = "SYSTEM"
	MetricCategoryUser        MetricCategory = "USER"
)

// AggregationType represents the type of aggregation
type AggregationType string

const (
	AggregationTypeSum      AggregationType = "SUM"
	AggregationTypeAvg      AggregationType = "AVG"
	AggregationTypeMin      AggregationType = "MIN"
	AggregationTypeMax      AggregationType = "MAX"
	AggregationTypeCount    AggregationType = "COUNT"
	AggregationTypeDistinct AggregationType = "DISTINCT"
)

// Dimension represents a dimension for metric analysis
type Dimension struct {
	Name        string `json:"name"`
	Type        string `json:"type"` // "string", "number", "date", "boolean"
	Description string `json:"description"`
	Required    bool   `json:"required"`
}

// MetricFilters represents filters for metric data
type MetricFilters map[string]interface{}

// MetricCalculation represents calculation configuration for derived metrics
type MetricCalculation struct {
	Formula       string                 `json:"formula"`
	Variables     map[string]interface{} `json:"variables"`
	Conditions    []CalculationCondition `json:"conditions"`
	IsValid       bool                   `json:"isValid"`
	LastValidated *time.Time             `json:"lastValidated,omitempty"`
}

// CalculationCondition represents a condition in metric calculation
type CalculationCondition struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"` // "=", "!=", ">", "<", ">=", "<=", "IN", "NOT_IN"
	Value    interface{} `json:"value"`
}

// MetricData represents a data point for a metric
type MetricData struct {
	ID         string                 `json:"id"`
	MetricID   shared.MetricID        `json:"metricId"`
	Timestamp  time.Time              `json:"timestamp"`
	Value      float64                `json:"value"`
	Dimensions map[string]interface{} `json:"dimensions"`
	Labels     map[string]string      `json:"labels"`
}

// NewMetric creates a new metric
func NewMetric(name, description string, metricType MetricType, category MetricCategory) *Metric {
	now := time.Now()
	return &Metric{
		ID:           shared.NewMetricID(),
		Name:         name,
		Description:  description,
		Type:         metricType,
		Category:     category,
		Unit:         "",
		Aggregation:  AggregationTypeSum,
		DataSource:   DataSource{},
		Dimensions:   []Dimension{},
		Filters:      make(MetricFilters),
		Calculation:  MetricCalculation{},
		IsCalculated: false,
		CreatedAt:    now,
		UpdatedAt:    now,
		Version:      1,
		Events:       []*shared.DomainEvent{},
	}
}

// SetUnit sets the unit for the metric
func (m *Metric) SetUnit(unit string) {
	m.Unit = unit
	m.UpdatedAt = time.Now()
	m.Version++

	event := shared.NewDomainEvent("UnitChanged", m.ID.String(), m.Version, map[string]interface{}{
		"unit": unit,
	})
	m.Events = append(m.Events, event)
}

// SetAggregation sets the aggregation type for the metric
func (m *Metric) SetAggregation(aggregation AggregationType) error {
	if !m.isValidAggregation(aggregation) {
		return shared.NewDomainError("INVALID_AGGREGATION", "Invalid aggregation type", "")
	}

	m.Aggregation = aggregation
	m.UpdatedAt = time.Now()
	m.Version++

	event := shared.NewDomainEvent("AggregationChanged", m.ID.String(), m.Version, map[string]interface{}{
		"aggregation": aggregation,
	})
	m.Events = append(m.Events, event)

	return nil
}

// SetDataSource sets the data source for the metric
func (m *Metric) SetDataSource(dataSource DataSource) {
	m.DataSource = dataSource
	m.UpdatedAt = time.Now()
	m.Version++

	event := shared.NewDomainEvent("DataSourceChanged", m.ID.String(), m.Version, map[string]interface{}{
		"dataSource": dataSource,
	})
	m.Events = append(m.Events, event)
}

// AddDimension adds a dimension to the metric
func (m *Metric) AddDimension(dimension Dimension) error {
	if dimension.Name == "" {
		return shared.NewDomainError("INVALID_DIMENSION", "Dimension name cannot be empty", "")
	}

	// Check for duplicate dimension names
	for _, existingDim := range m.Dimensions {
		if existingDim.Name == dimension.Name {
			return shared.NewDomainError("DUPLICATE_DIMENSION", "Dimension name already exists", "")
		}
	}

	m.Dimensions = append(m.Dimensions, dimension)
	m.UpdatedAt = time.Now()
	m.Version++

	event := shared.NewDomainEvent("DimensionAdded", m.ID.String(), m.Version, map[string]interface{}{
		"dimension": dimension,
	})
	m.Events = append(m.Events, event)

	return nil
}

// RemoveDimension removes a dimension from the metric
func (m *Metric) RemoveDimension(dimensionName string) error {
	for i, dimension := range m.Dimensions {
		if dimension.Name == dimensionName {
			m.Dimensions = append(m.Dimensions[:i], m.Dimensions[i+1:]...)
			m.UpdatedAt = time.Now()
			m.Version++

			event := shared.NewDomainEvent("DimensionRemoved", m.ID.String(), m.Version, map[string]interface{}{
				"dimensionName": dimensionName,
			})
			m.Events = append(m.Events, event)

			return nil
		}
	}
	return shared.ErrNotFound
}

// SetCalculation sets the calculation configuration for derived metrics
func (m *Metric) SetCalculation(calculation MetricCalculation) error {
	if calculation.Formula == "" {
		return shared.NewDomainError("INVALID_CALCULATION", "Calculation formula cannot be empty", "")
	}

	// Validate the calculation formula
	if err := m.validateCalculation(calculation); err != nil {
		return err
	}

	m.Calculation = calculation
	m.Calculation.IsValid = true
	now := time.Now()
	m.Calculation.LastValidated = &now
	m.IsCalculated = true
	m.UpdatedAt = now
	m.Version++

	event := shared.NewDomainEvent("CalculationUpdated", m.ID.String(), m.Version, map[string]interface{}{
		"calculation": calculation,
	})
	m.Events = append(m.Events, event)

	return nil
}

// AddFilter adds a filter to the metric
func (m *Metric) AddFilter(key string, value interface{}) {
	m.Filters[key] = value
	m.UpdatedAt = time.Now()
	m.Version++

	event := shared.NewDomainEvent("FilterAdded", m.ID.String(), m.Version, map[string]interface{}{
		"key":   key,
		"value": value,
	})
	m.Events = append(m.Events, event)
}

// RemoveFilter removes a filter from the metric
func (m *Metric) RemoveFilter(key string) {
	if _, exists := m.Filters[key]; exists {
		delete(m.Filters, key)
		m.UpdatedAt = time.Now()
		m.Version++

		event := shared.NewDomainEvent("FilterRemoved", m.ID.String(), m.Version, map[string]interface{}{
			"key": key,
		})
		m.Events = append(m.Events, event)
	}
}

// validateCalculation validates the calculation formula
func (m *Metric) validateCalculation(calculation MetricCalculation) error {
	// Basic validation - in a real implementation, this would parse and validate the formula
	if len(calculation.Formula) < 3 {
		return shared.NewDomainError("INVALID_FORMULA", "Formula is too short", "")
	}

	// Check for required variables
	for _, condition := range calculation.Conditions {
		if condition.Field == "" || condition.Operator == "" {
			return shared.NewDomainError("INVALID_CONDITION", "Condition must have field and operator", "")
		}
	}

	return nil
}

// isValidAggregation checks if the aggregation type is valid
func (m *Metric) isValidAggregation(aggregation AggregationType) bool {
	validAggregations := []AggregationType{
		AggregationTypeSum, AggregationTypeAvg, AggregationTypeMin,
		AggregationTypeMax, AggregationTypeCount, AggregationTypeDistinct,
	}

	for _, validAgg := range validAggregations {
		if aggregation == validAgg {
			return true
		}
	}
	return false
}

// ClearEvents clears the domain events
func (m *Metric) ClearEvents() {
	m.Events = []*shared.DomainEvent{}
}
