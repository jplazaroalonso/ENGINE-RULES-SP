package campaign

import (
	"context"
	"time"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/campaigns-management-service/internal/domain/shared"
)

// Repository defines the contract for campaign persistence
type Repository interface {
	Save(ctx context.Context, campaign *Campaign) error
	FindByID(ctx context.Context, id CampaignID) (*Campaign, error)
	FindByName(ctx context.Context, name string) (*Campaign, error)
	List(ctx context.Context, criteria ListCriteria) ([]*Campaign, error)
	Count(ctx context.Context, criteria ListCriteria) (int64, error)
	Delete(ctx context.Context, id CampaignID) error
	ExistsByName(ctx context.Context, name string) (bool, error)
	FindByStatus(ctx context.Context, status CampaignStatus) ([]*Campaign, error)
	FindByType(ctx context.Context, campaignType CampaignType) ([]*Campaign, error)
	FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*Campaign, error)
	FindByCreatedBy(ctx context.Context, createdBy shared.UserID) ([]*Campaign, error)
}

// ListCriteria defines filtering and pagination options for listing campaigns
type ListCriteria struct {
	Status       *CampaignStatus
	CampaignType *CampaignType
	CreatedBy    *shared.UserID
	StartDate    *time.Time
	EndDate      *time.Time
	SearchQuery  string
	PageSize     int
	PageOffset   int
	SortBy       string
	SortOrder    string
}

// CampaignEventRepository defines the contract for campaign event persistence
type CampaignEventRepository interface {
	Save(ctx context.Context, event CampaignEvent) error
	FindByCampaignID(ctx context.Context, campaignID CampaignID, limit int) ([]CampaignEvent, error)
	FindByCampaignIDAndType(ctx context.Context, campaignID CampaignID, eventType CampaignEventType, limit int) ([]CampaignEvent, error)
	CountByCampaignID(ctx context.Context, campaignID CampaignID) (int64, error)
	CountByCampaignIDAndType(ctx context.Context, campaignID CampaignID, eventType CampaignEventType) (int64, error)
	FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]CampaignEvent, error)
	AggregateMetrics(ctx context.Context, campaignID CampaignID, startDate, endDate time.Time) (*CampaignMetrics, error)
}

// CampaignMetricsRepository defines the contract for campaign metrics persistence
type CampaignMetricsRepository interface {
	Save(ctx context.Context, campaignID CampaignID, metrics CampaignMetrics) error
	FindByCampaignID(ctx context.Context, campaignID CampaignID) (*CampaignMetrics, error)
	UpdateMetrics(ctx context.Context, campaignID CampaignID, metrics CampaignMetrics) error
	ResetMetrics(ctx context.Context, campaignID CampaignID) error
	FindTopPerformers(ctx context.Context, limit int, metric string) ([]*Campaign, error)
	FindUnderperformers(ctx context.Context, limit int, metric string) ([]*Campaign, error)
}

// CampaignTargetingService defines the contract for campaign targeting operations
type CampaignTargetingService interface {
	EvaluateTargeting(ctx context.Context, campaignID CampaignID, customerID shared.CustomerID) (bool, error)
	GetTargetAudience(ctx context.Context, campaignID CampaignID) ([]shared.CustomerID, error)
	UpdateTargetingRules(ctx context.Context, campaignID CampaignID, rules []shared.RuleID) error
	ValidateTargetingRules(ctx context.Context, rules []shared.RuleID) error
}

// CampaignPerformanceService defines the contract for campaign performance operations
type CampaignPerformanceService interface {
	CalculateMetrics(ctx context.Context, campaignID CampaignID) (*CampaignMetrics, error)
	TrackEvent(ctx context.Context, campaignID CampaignID, event CampaignEvent) error
	GetPerformanceReport(ctx context.Context, campaignID CampaignID, period TimePeriod) (*PerformanceReport, error)
	GetPerformanceTrends(ctx context.Context, campaignID CampaignID, period TimePeriod) ([]PerformanceDataPoint, error)
	ComparePerformance(ctx context.Context, campaignIDs []CampaignID, period TimePeriod) (*PerformanceComparison, error)
}

// TimePeriod represents a time period for performance analysis
type TimePeriod struct {
	StartDate   time.Time
	EndDate     time.Time
	Granularity string // "hour", "day", "week", "month"
}

// PerformanceReport represents a comprehensive performance report
type PerformanceReport struct {
	CampaignID      CampaignID             `json:"campaignId"`
	Period          TimePeriod             `json:"period"`
	Metrics         CampaignMetrics        `json:"metrics"`
	Trends          []PerformanceDataPoint `json:"trends"`
	Recommendations []string               `json:"recommendations"`
	Benchmarks      BenchmarkData          `json:"benchmarks"`
	GeneratedAt     time.Time              `json:"generatedAt"`
}

// PerformanceDataPoint represents a single data point in performance trends
type PerformanceDataPoint struct {
	Timestamp   time.Time       `json:"timestamp"`
	Metrics     CampaignMetrics `json:"metrics"`
	Impressions int64           `json:"impressions"`
	Clicks      int64           `json:"clicks"`
	Conversions int64           `json:"conversions"`
	Revenue     shared.Money    `json:"revenue"`
	Cost        shared.Money    `json:"cost"`
}

// PerformanceComparison represents a comparison between multiple campaigns
type PerformanceComparison struct {
	Campaigns      []CampaignComparison `json:"campaigns"`
	Period         TimePeriod           `json:"period"`
	BestPerformer  CampaignID           `json:"bestPerformer"`
	WorstPerformer CampaignID           `json:"worstPerformer"`
	AverageMetrics CampaignMetrics      `json:"averageMetrics"`
	GeneratedAt    time.Time            `json:"generatedAt"`
}

// CampaignComparison represents a single campaign in a performance comparison
type CampaignComparison struct {
	CampaignID CampaignID      `json:"campaignId"`
	Name       string          `json:"name"`
	Metrics    CampaignMetrics `json:"metrics"`
	Rank       int             `json:"rank"`
	Score      float64         `json:"score"`
}

// BenchmarkData represents industry or historical benchmark data
type BenchmarkData struct {
	IndustryCTR            float64         `json:"industryCTR"`
	IndustryConversionRate float64         `json:"industryConversionRate"`
	IndustryCPC            shared.Money    `json:"industryCPC"`
	IndustryROI            float64         `json:"industryROI"`
	HistoricalAverage      CampaignMetrics `json:"historicalAverage"`
	PercentileRank         int             `json:"percentileRank"`
}

// CampaignSchedulingService defines the contract for campaign scheduling operations
type CampaignSchedulingService interface {
	ScheduleCampaign(ctx context.Context, campaignID CampaignID, schedule CampaignSchedule) error
	CancelScheduledCampaign(ctx context.Context, campaignID CampaignID) error
	GetScheduledCampaigns(ctx context.Context, startDate, endDate time.Time) ([]*Campaign, error)
	ProcessScheduledCampaigns(ctx context.Context) error
	UpdateSchedule(ctx context.Context, campaignID CampaignID, schedule CampaignSchedule) error
}

// CampaignSchedule represents campaign scheduling information
type CampaignSchedule struct {
	CampaignID CampaignID `json:"campaignId"`
	StartDate  time.Time  `json:"startDate"`
	EndDate    *time.Time `json:"endDate,omitempty"`
	Recurrence Recurrence `json:"recurrence"`
	Timezone   string     `json:"timezone"`
	IsActive   bool       `json:"isActive"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
}

// Recurrence represents campaign recurrence patterns
type Recurrence struct {
	Type           string     `json:"type"` // "once", "daily", "weekly", "monthly"
	Interval       int        `json:"interval"`
	DaysOfWeek     []int      `json:"daysOfWeek,omitempty"`  // 0=Sunday, 1=Monday, etc.
	DaysOfMonth    []int      `json:"daysOfMonth,omitempty"` // 1-31
	EndDate        *time.Time `json:"endDate,omitempty"`
	MaxOccurrences *int       `json:"maxOccurrences,omitempty"`
}

// CampaignNotificationService defines the contract for campaign notifications
type CampaignNotificationService interface {
	SendCampaignStartedNotification(ctx context.Context, campaignID CampaignID) error
	SendCampaignEndedNotification(ctx context.Context, campaignID CampaignID) error
	SendPerformanceAlert(ctx context.Context, campaignID CampaignID, alert PerformanceAlert) error
	SendBudgetAlert(ctx context.Context, campaignID CampaignID, alert BudgetAlert) error
	SendCustomNotification(ctx context.Context, campaignID CampaignID, notification CustomNotification) error
}

// PerformanceAlert represents a performance-related alert
type PerformanceAlert struct {
	Type         string  `json:"type"` // "low_ctr", "low_conversion", "high_cpc", "negative_roi"
	Threshold    float64 `json:"threshold"`
	CurrentValue float64 `json:"currentValue"`
	Message      string  `json:"message"`
	Severity     string  `json:"severity"` // "low", "medium", "high", "critical"
}

// BudgetAlert represents a budget-related alert
type BudgetAlert struct {
	Type         string  `json:"type"` // "approaching_limit", "exceeded", "depleted"
	Threshold    float64 `json:"threshold"`
	CurrentValue float64 `json:"currentValue"`
	Message      string  `json:"message"`
	Severity     string  `json:"severity"` // "low", "medium", "high", "critical"
}

// CustomNotification represents a custom notification
type CustomNotification struct {
	Title      string                 `json:"title"`
	Message    string                 `json:"message"`
	Type       string                 `json:"type"` // "info", "warning", "error", "success"
	Recipients []string               `json:"recipients"`
	Data       map[string]interface{} `json:"data,omitempty"`
	Priority   string                 `json:"priority"` // "low", "normal", "high", "urgent"
}
