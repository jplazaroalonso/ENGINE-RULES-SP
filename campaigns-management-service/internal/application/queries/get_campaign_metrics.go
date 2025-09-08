package queries

import (
	"context"
	"time"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/campaigns-management-service/internal/domain/campaign"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/campaigns-management-service/internal/domain/shared"
)

// GetCampaignMetricsQuery represents the query to get campaign metrics
type GetCampaignMetricsQuery struct {
	CampaignID  string     `json:"campaignId" validate:"required"`
	StartDate   *time.Time `json:"startDate,omitempty"`
	EndDate     *time.Time `json:"endDate,omitempty"`
	Granularity string     `json:"granularity,omitempty" validate:"omitempty,oneof=hour day week month"`
}

// GetCampaignMetricsResult represents the result of getting campaign metrics
type GetCampaignMetricsResult struct {
	CampaignID      string                         `json:"campaignId"`
	CampaignName    string                         `json:"campaignName"`
	Period          TimePeriodResponse             `json:"period"`
	Metrics         CampaignMetricsResponse        `json:"metrics"`
	Trends          []PerformanceDataPointResponse `json:"trends"`
	Recommendations []string                       `json:"recommendations"`
	Benchmarks      BenchmarkDataResponse          `json:"benchmarks"`
	GeneratedAt     string                         `json:"generatedAt"`
}

// TimePeriodResponse represents a time period in the response
type TimePeriodResponse struct {
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
	Granularity string `json:"granularity"`
}

// PerformanceDataPointResponse represents a performance data point in the response
type PerformanceDataPointResponse struct {
	Timestamp      string         `json:"timestamp"`
	Impressions    int64          `json:"impressions"`
	Clicks         int64          `json:"clicks"`
	Conversions    int64          `json:"conversions"`
	Revenue        BudgetResponse `json:"revenue"`
	Cost           BudgetResponse `json:"cost"`
	CTR            float64        `json:"ctr"`
	ConversionRate float64        `json:"conversionRate"`
	ROI            float64        `json:"roi"`
}

// BenchmarkDataResponse represents benchmark data in the response
type BenchmarkDataResponse struct {
	IndustryCTR            float64                 `json:"industryCTR"`
	IndustryConversionRate float64                 `json:"industryConversionRate"`
	IndustryCPC            BudgetResponse          `json:"industryCPC"`
	IndustryROI            float64                 `json:"industryROI"`
	HistoricalAverage      CampaignMetricsResponse `json:"historicalAverage"`
	PercentileRank         int                     `json:"percentileRank"`
}

// GetCampaignMetricsHandler handles get campaign metrics queries
type GetCampaignMetricsHandler struct {
	campaignService *campaign.CampaignService
	validator       shared.Validator
}

// NewGetCampaignMetricsHandler creates a new GetCampaignMetricsHandler
func NewGetCampaignMetricsHandler(
	campaignService *campaign.CampaignService,
	validator shared.Validator,
) *GetCampaignMetricsHandler {
	return &GetCampaignMetricsHandler{
		campaignService: campaignService,
		validator:       validator,
	}
}

// Handle processes the get campaign metrics query
func (h *GetCampaignMetricsHandler) Handle(ctx context.Context, query GetCampaignMetricsQuery) (*GetCampaignMetricsResult, error) {
	// Validate query input
	if err := h.validator.Validate(query); err != nil {
		return nil, shared.NewValidationError("invalid get campaign metrics query", err)
	}

	// Parse campaign ID
	campaignID, err := campaign.NewCampaignIDFromString(query.CampaignID)
	if err != nil {
		return nil, shared.NewValidationError("invalid campaign ID", err)
	}

	// Get campaign to ensure it exists and get name
	campaign, err := h.campaignService.GetCampaign(ctx, campaignID)
	if err != nil {
		return nil, err
	}

	// Set default time period if not provided
	startDate := query.StartDate
	endDate := query.EndDate
	if startDate == nil {
		// Default to last 30 days
		now := time.Now()
		thirtyDaysAgo := now.AddDate(0, 0, -30)
		startDate = &thirtyDaysAgo
	}
	if endDate == nil {
		now := time.Now()
		endDate = &now
	}

	// Set default granularity
	granularity := query.Granularity
	if granularity == "" {
		granularity = "day"
	}

	// Create time period
	period := campaign.TimePeriod{
		StartDate:   *startDate,
		EndDate:     *endDate,
		Granularity: granularity,
	}

	// Get performance report using domain service
	report, err := h.campaignService.GetCampaignPerformanceReport(ctx, campaignID, period)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	result := h.mapPerformanceReportToResponse(campaign, report)

	return result, nil
}

// mapPerformanceReportToResponse maps a performance report to response format
func (h *GetCampaignMetricsHandler) mapPerformanceReportToResponse(
	campaign *campaign.Campaign,
	report *campaign.PerformanceReport,
) *GetCampaignMetricsResult {
	// Map time period
	timePeriod := TimePeriodResponse{
		StartDate:   report.Period.StartDate.Format("2006-01-02T15:04:05Z"),
		EndDate:     report.Period.EndDate.Format("2006-01-02T15:04:05Z"),
		Granularity: report.Period.Granularity,
	}

	// Map metrics
	metrics := h.mapMetricsToResponse(report.Metrics)

	// Map trends
	trends := make([]PerformanceDataPointResponse, len(report.Trends))
	for i, trend := range report.Trends {
		trends[i] = PerformanceDataPointResponse{
			Timestamp:   trend.Timestamp.Format("2006-01-02T15:04:05Z"),
			Impressions: trend.Impressions,
			Clicks:      trend.Clicks,
			Conversions: trend.Conversions,
			Revenue: BudgetResponse{
				Amount:   trend.Revenue.Amount,
				Currency: trend.Revenue.Currency,
			},
			Cost: BudgetResponse{
				Amount:   trend.Cost.Amount,
				Currency: trend.Cost.Currency,
			},
			CTR:            trend.Metrics.CTR,
			ConversionRate: trend.Metrics.ConversionRate,
			ROI:            trend.Metrics.ROI,
		}
	}

	// Map benchmarks
	benchmarks := BenchmarkDataResponse{
		IndustryCTR:            report.Benchmarks.IndustryCTR,
		IndustryConversionRate: report.Benchmarks.IndustryConversionRate,
		IndustryCPC: BudgetResponse{
			Amount:   report.Benchmarks.IndustryCPC.Amount,
			Currency: report.Benchmarks.IndustryCPC.Currency,
		},
		IndustryROI:       report.Benchmarks.IndustryROI,
		HistoricalAverage: h.mapMetricsToResponse(report.Benchmarks.HistoricalAverage),
		PercentileRank:    report.Benchmarks.PercentileRank,
	}

	return &GetCampaignMetricsResult{
		CampaignID:      campaign.ID().String(),
		CampaignName:    campaign.Name(),
		Period:          timePeriod,
		Metrics:         metrics,
		Trends:          trends,
		Recommendations: report.Recommendations,
		Benchmarks:      benchmarks,
		GeneratedAt:     report.GeneratedAt.Format("2006-01-02T15:04:05Z"),
	}
}

// mapMetricsToResponse maps campaign metrics to response format
func (h *GetCampaignMetricsHandler) mapMetricsToResponse(metrics campaign.CampaignMetrics) CampaignMetricsResponse {
	return CampaignMetricsResponse{
		Impressions: metrics.Impressions,
		Clicks:      metrics.Clicks,
		Conversions: metrics.Conversions,
		Revenue: BudgetResponse{
			Amount:   metrics.Revenue.Amount,
			Currency: metrics.Revenue.Currency,
		},
		Cost: BudgetResponse{
			Amount:   metrics.Cost.Amount,
			Currency: metrics.Cost.Currency,
		},
		CTR:            metrics.CTR,
		ConversionRate: metrics.ConversionRate,
		CostPerClick: BudgetResponse{
			Amount:   metrics.CostPerClick.Amount,
			Currency: metrics.CostPerClick.Currency,
		},
		CostPerConversion: BudgetResponse{
			Amount:   metrics.CostPerConversion.Amount,
			Currency: metrics.CostPerConversion.Currency,
		},
		ROAS:        metrics.ROAS,
		ROI:         metrics.ROI,
		LastUpdated: metrics.LastUpdated.Format("2006-01-02T15:04:05Z"),
	}
}
