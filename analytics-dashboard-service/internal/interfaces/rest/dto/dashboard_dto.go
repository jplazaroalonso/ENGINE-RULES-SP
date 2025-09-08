package dto

import (
	"time"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/analytics-dashboard-service/internal/domain/analytics"
)

// CreateDashboardRequest represents the request to create a dashboard
type CreateDashboardRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=255"`
	Description string `json:"description" validate:"max=1000"`
	OwnerID     string `json:"ownerId" validate:"required,uuid"`
}

// DashboardResponse represents the response for a dashboard
type DashboardResponse struct {
	ID              string                     `json:"id"`
	Name            string                     `json:"name"`
	Description     string                     `json:"description"`
	Layout          analytics.DashboardLayout  `json:"layout"`
	Widgets         []analytics.Widget         `json:"widgets"`
	Filters         analytics.DashboardFilters `json:"filters"`
	RefreshInterval int                        `json:"refreshInterval"`
	IsPublic        bool                       `json:"isPublic"`
	OwnerID         string                     `json:"ownerId"`
	CreatedAt       time.Time                  `json:"createdAt"`
	UpdatedAt       time.Time                  `json:"updatedAt"`
	Version         int                        `json:"version"`
}

// CreateReportRequest represents the request to create a report
type CreateReportRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=255"`
	Description string `json:"description" validate:"max=1000"`
	Type        string `json:"type" validate:"required,oneof=PERFORMANCE COMPLIANCE BUSINESS CUSTOM"`
	OwnerID     string `json:"ownerId" validate:"required,uuid"`
}

// ReportResponse represents the response for a report
type ReportResponse struct {
	ID            string                     `json:"id"`
	Name          string                     `json:"name"`
	Description   string                     `json:"description"`
	Type          string                     `json:"type"`
	Template      analytics.ReportTemplate   `json:"template"`
	Parameters    analytics.ReportParameters `json:"parameters"`
	Schedule      *analytics.ReportSchedule  `json:"schedule,omitempty"`
	OutputFormat  string                     `json:"outputFormat"`
	Recipients    []string                   `json:"recipients"`
	Status        string                     `json:"status"`
	LastGenerated *time.Time                 `json:"lastGenerated,omitempty"`
	NextRun       *time.Time                 `json:"nextRun,omitempty"`
	OwnerID       string                     `json:"ownerId"`
	CreatedAt     time.Time                  `json:"createdAt"`
	UpdatedAt     time.Time                  `json:"updatedAt"`
	Version       int                        `json:"version"`
}

// CreateMetricRequest represents the request to create a metric
type CreateMetricRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=255"`
	Description string `json:"description" validate:"max=1000"`
	Type        string `json:"type" validate:"required,oneof=COUNTER GAUGE HISTOGRAM SUMMARY"`
	Category    string `json:"category" validate:"required,oneof=PERFORMANCE BUSINESS SYSTEM USER"`
	Unit        string `json:"unit" validate:"max=50"`
}

// MetricResponse represents the response for a metric
type MetricResponse struct {
	ID           string                      `json:"id"`
	Name         string                      `json:"name"`
	Description  string                      `json:"description"`
	Type         string                      `json:"type"`
	Category     string                      `json:"category"`
	Unit         string                      `json:"unit"`
	Aggregation  string                      `json:"aggregation"`
	DataSource   analytics.DataSource        `json:"dataSource"`
	Dimensions   []analytics.Dimension       `json:"dimensions"`
	Filters      analytics.MetricFilters     `json:"filters"`
	Calculation  analytics.MetricCalculation `json:"calculation"`
	IsCalculated bool                        `json:"isCalculated"`
	CreatedAt    time.Time                   `json:"createdAt"`
	UpdatedAt    time.Time                   `json:"updatedAt"`
	Version      int                         `json:"version"`
}

// MetricDataResponse represents the response for metric data
type MetricDataResponse struct {
	ID         string                 `json:"id"`
	MetricID   string                 `json:"metricId"`
	Timestamp  time.Time              `json:"timestamp"`
	Value      float64                `json:"value"`
	Dimensions map[string]interface{} `json:"dimensions"`
	Labels     map[string]string      `json:"labels"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}
