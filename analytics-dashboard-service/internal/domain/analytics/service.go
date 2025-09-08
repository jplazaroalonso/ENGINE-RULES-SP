package analytics

import (
	"context"
	"time"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/analytics-dashboard-service/internal/domain/shared"
)

// DashboardRepository interface for dashboard persistence
type DashboardRepository interface {
	Save(ctx context.Context, dashboard *Dashboard) error
	FindByID(ctx context.Context, id shared.DashboardID) (*Dashboard, error)
	FindByOwnerID(ctx context.Context, ownerID shared.UserID) ([]*Dashboard, error)
	FindPublic(ctx context.Context) ([]*Dashboard, error)
	Delete(ctx context.Context, id shared.DashboardID) error
	Update(ctx context.Context, dashboard *Dashboard) error
}

// ReportRepository interface for report persistence
type ReportRepository interface {
	Save(ctx context.Context, report *Report) error
	FindByID(ctx context.Context, id shared.ReportID) (*Report, error)
	FindByOwnerID(ctx context.Context, ownerID shared.UserID) ([]*Report, error)
	FindByStatus(ctx context.Context, status ReportStatus) ([]*Report, error)
	FindScheduled(ctx context.Context, before time.Time) ([]*Report, error)
	Delete(ctx context.Context, id shared.ReportID) error
	Update(ctx context.Context, report *Report) error
}

// MetricRepository interface for metric persistence
type MetricRepository interface {
	Save(ctx context.Context, metric *Metric) error
	FindByID(ctx context.Context, id shared.MetricID) (*Metric, error)
	FindByCategory(ctx context.Context, category MetricCategory) ([]*Metric, error)
	FindByType(ctx context.Context, metricType MetricType) ([]*Metric, error)
	FindAll(ctx context.Context) ([]*Metric, error)
	Delete(ctx context.Context, id shared.MetricID) error
	Update(ctx context.Context, metric *Metric) error
}

// MetricDataRepository interface for metric data persistence
type MetricDataRepository interface {
	Save(ctx context.Context, data *MetricData) error
	SaveBatch(ctx context.Context, data []*MetricData) error
	FindByMetricID(ctx context.Context, metricID shared.MetricID, timeRange *shared.TimeRange) ([]*MetricData, error)
	FindByMetricIDAndDimensions(ctx context.Context, metricID shared.MetricID, dimensions map[string]interface{}, timeRange *shared.TimeRange) ([]*MetricData, error)
	AggregateByMetricID(ctx context.Context, metricID shared.MetricID, aggregation AggregationType, timeRange *shared.TimeRange) (*MetricData, error)
	DeleteOldData(ctx context.Context, before time.Time) error
}

// DataAggregator interface for aggregating data from external sources
type DataAggregator interface {
	AggregateRulesData(ctx context.Context, timeRange *shared.TimeRange) (map[string]interface{}, error)
	AggregateCustomerData(ctx context.Context, timeRange *shared.TimeRange) (map[string]interface{}, error)
	AggregateCampaignData(ctx context.Context, timeRange *shared.TimeRange) (map[string]interface{}, error)
	AggregatePromotionData(ctx context.Context, timeRange *shared.TimeRange) (map[string]interface{}, error)
}

// ReportGenerator interface for generating reports
type ReportGenerator interface {
	GenerateReport(ctx context.Context, report *Report, data map[string]interface{}) ([]byte, error)
	GeneratePDF(ctx context.Context, report *Report, data map[string]interface{}) ([]byte, error)
	GenerateExcel(ctx context.Context, report *Report, data map[string]interface{}) ([]byte, error)
	GenerateCSV(ctx context.Context, report *Report, data map[string]interface{}) ([]byte, error)
	GenerateJSON(ctx context.Context, report *Report, data map[string]interface{}) ([]byte, error)
	GenerateHTML(ctx context.Context, report *Report, data map[string]interface{}) ([]byte, error)
}

// NotificationService interface for sending notifications
type NotificationService interface {
	SendReportNotification(ctx context.Context, report *Report, recipients []string, attachment []byte) error
	SendErrorNotification(ctx context.Context, report *Report, error error) error
}

// CacheService interface for caching analytics data
type CacheService interface {
	Get(ctx context.Context, key string) (interface{}, error)
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Clear(ctx context.Context) error
}

// AnalyticsService interface for the main analytics service
type AnalyticsService interface {
	// Dashboard operations
	CreateDashboard(ctx context.Context, name, description string, ownerID shared.UserID) (*Dashboard, error)
	GetDashboard(ctx context.Context, id shared.DashboardID) (*Dashboard, error)
	UpdateDashboard(ctx context.Context, dashboard *Dashboard) error
	DeleteDashboard(ctx context.Context, id shared.DashboardID) error
	ListDashboards(ctx context.Context, ownerID shared.UserID) ([]*Dashboard, error)

	// Widget operations
	AddWidget(ctx context.Context, dashboardID shared.DashboardID, widget Widget) error
	UpdateWidget(ctx context.Context, dashboardID shared.DashboardID, widget Widget) error
	RemoveWidget(ctx context.Context, dashboardID shared.DashboardID, widgetID shared.WidgetID) error
	RefreshWidget(ctx context.Context, dashboardID shared.DashboardID, widgetID shared.WidgetID) (map[string]interface{}, error)

	// Report operations
	CreateReport(ctx context.Context, name, description string, reportType ReportType, ownerID shared.UserID) (*Report, error)
	GetReport(ctx context.Context, id shared.ReportID) (*Report, error)
	UpdateReport(ctx context.Context, report *Report) error
	DeleteReport(ctx context.Context, id shared.ReportID) error
	GenerateReport(ctx context.Context, id shared.ReportID) ([]byte, error)
	ListReports(ctx context.Context, ownerID shared.UserID) ([]*Report, error)

	// Metric operations
	CreateMetric(ctx context.Context, name, description string, metricType MetricType, category MetricCategory) (*Metric, error)
	GetMetric(ctx context.Context, id shared.MetricID) (*Metric, error)
	UpdateMetric(ctx context.Context, metric *Metric) error
	DeleteMetric(ctx context.Context, id shared.MetricID) error
	ListMetrics(ctx context.Context) ([]*Metric, error)
	GetMetricData(ctx context.Context, metricID shared.MetricID, timeRange *shared.TimeRange) ([]*MetricData, error)

	// Real-time analytics
	GetRealTimeAnalytics(ctx context.Context) (map[string]interface{}, error)
	GetPerformanceMetrics(ctx context.Context) (map[string]interface{}, error)
	GetBusinessMetrics(ctx context.Context) (map[string]interface{}, error)
	GetComplianceMetrics(ctx context.Context) (map[string]interface{}, error)
}
