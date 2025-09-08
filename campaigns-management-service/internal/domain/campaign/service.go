package campaign

import (
	"context"
	"fmt"
	"time"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/campaigns-management-service/internal/domain/shared"
)

// CampaignService provides domain services for campaign operations
type CampaignService struct {
	repo                Repository
	eventRepo           CampaignEventRepository
	metricsRepo         CampaignMetricsRepository
	targetingService    CampaignTargetingService
	performanceService  CampaignPerformanceService
	schedulingService   CampaignSchedulingService
	notificationService CampaignNotificationService
	eventBus            shared.EventBus
}

// NewCampaignService creates a new campaign service
func NewCampaignService(
	repo Repository,
	eventRepo CampaignEventRepository,
	metricsRepo CampaignMetricsRepository,
	targetingService CampaignTargetingService,
	performanceService CampaignPerformanceService,
	schedulingService CampaignSchedulingService,
	notificationService CampaignNotificationService,
	eventBus shared.EventBus,
) *CampaignService {
	return &CampaignService{
		repo:                repo,
		eventRepo:           eventRepo,
		metricsRepo:         metricsRepo,
		targetingService:    targetingService,
		performanceService:  performanceService,
		schedulingService:   schedulingService,
		notificationService: notificationService,
		eventBus:            eventBus,
	}
}

// CreateCampaign creates a new campaign with validation and event publishing
func (s *CampaignService) CreateCampaign(
	ctx context.Context,
	name, description string,
	campaignType CampaignType,
	targetingRules []shared.RuleID,
	startDate time.Time,
	endDate *time.Time,
	budget *shared.Money,
	createdBy shared.UserID,
	settings CampaignSettings,
) (*Campaign, error) {
	// Check if campaign name already exists
	exists, err := s.repo.ExistsByName(ctx, name)
	if err != nil {
		return nil, shared.NewInfrastructureError("failed to check campaign name existence", err)
	}
	if exists {
		return nil, shared.NewBusinessError("campaign name already exists", nil)
	}

	// Validate targeting rules
	if err := s.targetingService.ValidateTargetingRules(ctx, targetingRules); err != nil {
		return nil, shared.NewValidationError("invalid targeting rules", err)
	}

	// Create campaign
	campaign, err := NewCampaign(
		name, description, campaignType, targetingRules,
		startDate, endDate, budget, createdBy, settings,
	)
	if err != nil {
		return nil, err
	}

	// Save campaign
	if err := s.repo.Save(ctx, campaign); err != nil {
		return nil, shared.NewInfrastructureError("failed to save campaign", err)
	}

	// Initialize metrics
	metrics := NewCampaignMetrics()
	if err := s.metricsRepo.Save(ctx, campaign.ID(), metrics); err != nil {
		return nil, shared.NewInfrastructureError("failed to initialize campaign metrics", err)
	}

	// Publish domain events
	events := campaign.GetEvents()
	for _, event := range events {
		if err := s.eventBus.Publish(ctx, event); err != nil {
			// Log error but don't fail the operation
			// In production, consider implementing outbox pattern
			fmt.Printf("Failed to publish event: %v\n", err)
		}
	}

	return campaign, nil
}

// ActivateCampaign activates a campaign with proper validation
func (s *CampaignService) ActivateCampaign(ctx context.Context, campaignID CampaignID, activatedBy shared.UserID) error {
	// Get campaign
	campaign, err := s.repo.FindByID(ctx, campaignID)
	if err != nil {
		return shared.NewInfrastructureError("failed to find campaign", err)
	}
	if campaign == nil {
		return shared.NewNotFoundError("campaign not found", nil)
	}

	// Activate campaign
	if err := campaign.Activate(); err != nil {
		return err
	}

	// Save updated campaign
	if err := s.repo.Save(ctx, campaign); err != nil {
		return shared.NewInfrastructureError("failed to save campaign", err)
	}

	// Publish domain events
	events := campaign.GetEvents()
	for _, event := range events {
		if err := s.eventBus.Publish(ctx, event); err != nil {
			fmt.Printf("Failed to publish event: %v\n", err)
		}
	}

	// Send notification
	if err := s.notificationService.SendCampaignStartedNotification(ctx, campaignID); err != nil {
		fmt.Printf("Failed to send campaign started notification: %v\n", err)
	}

	return nil
}

// PauseCampaign pauses a campaign
func (s *CampaignService) PauseCampaign(ctx context.Context, campaignID CampaignID, pausedBy shared.UserID) error {
	campaign, err := s.repo.FindByID(ctx, campaignID)
	if err != nil {
		return shared.NewInfrastructureError("failed to find campaign", err)
	}
	if campaign == nil {
		return shared.NewNotFoundError("campaign not found", nil)
	}

	if err := campaign.Pause(); err != nil {
		return err
	}

	if err := s.repo.Save(ctx, campaign); err != nil {
		return shared.NewInfrastructureError("failed to save campaign", err)
	}

	// Publish domain events
	events := campaign.GetEvents()
	for _, event := range events {
		if err := s.eventBus.Publish(ctx, event); err != nil {
			fmt.Printf("Failed to publish event: %v\n", err)
		}
	}

	return nil
}

// ResumeCampaign resumes a paused campaign
func (s *CampaignService) ResumeCampaign(ctx context.Context, campaignID CampaignID, resumedBy shared.UserID) error {
	campaign, err := s.repo.FindByID(ctx, campaignID)
	if err != nil {
		return shared.NewInfrastructureError("failed to find campaign", err)
	}
	if campaign == nil {
		return shared.NewNotFoundError("campaign not found", nil)
	}

	if err := campaign.Resume(); err != nil {
		return err
	}

	if err := s.repo.Save(ctx, campaign); err != nil {
		return shared.NewInfrastructureError("failed to save campaign", err)
	}

	// Publish domain events
	events := campaign.GetEvents()
	for _, event := range events {
		if err := s.eventBus.Publish(ctx, event); err != nil {
			fmt.Printf("Failed to publish event: %v\n", err)
		}
	}

	return nil
}

// CompleteCampaign completes a campaign
func (s *CampaignService) CompleteCampaign(ctx context.Context, campaignID CampaignID, completedBy shared.UserID) error {
	campaign, err := s.repo.FindByID(ctx, campaignID)
	if err != nil {
		return shared.NewInfrastructureError("failed to find campaign", err)
	}
	if campaign == nil {
		return shared.NewNotFoundError("campaign not found", nil)
	}

	if err := campaign.Complete(); err != nil {
		return err
	}

	if err := s.repo.Save(ctx, campaign); err != nil {
		return shared.NewInfrastructureError("failed to save campaign", err)
	}

	// Publish domain events
	events := campaign.GetEvents()
	for _, event := range events {
		if err := s.eventBus.Publish(ctx, event); err != nil {
			fmt.Printf("Failed to publish event: %v\n", err)
		}
	}

	// Send notification
	if err := s.notificationService.SendCampaignEndedNotification(ctx, campaignID); err != nil {
		fmt.Printf("Failed to send campaign ended notification: %v\n", err)
	}

	return nil
}

// CancelCampaign cancels a campaign
func (s *CampaignService) CancelCampaign(ctx context.Context, campaignID CampaignID, reason string, cancelledBy shared.UserID) error {
	campaign, err := s.repo.FindByID(ctx, campaignID)
	if err != nil {
		return shared.NewInfrastructureError("failed to find campaign", err)
	}
	if campaign == nil {
		return shared.NewNotFoundError("campaign not found", nil)
	}

	if err := campaign.Cancel(reason); err != nil {
		return err
	}

	if err := s.repo.Save(ctx, campaign); err != nil {
		return shared.NewInfrastructureError("failed to save campaign", err)
	}

	// Publish domain events
	events := campaign.GetEvents()
	for _, event := range events {
		if err := s.eventBus.Publish(ctx, event); err != nil {
			fmt.Printf("Failed to publish event: %v\n", err)
		}
	}

	return nil
}

// TrackCampaignEvent tracks a campaign event and updates metrics
func (s *CampaignService) TrackCampaignEvent(
	ctx context.Context,
	campaignID CampaignID,
	eventType CampaignEventType,
	customerID *shared.CustomerID,
	eventData map[string]interface{},
	revenue *shared.Money,
	cost *shared.Money,
) error {
	// Get campaign to ensure it exists
	campaign, err := s.repo.FindByID(ctx, campaignID)
	if err != nil {
		return shared.NewInfrastructureError("failed to find campaign", err)
	}
	if campaign == nil {
		return shared.NewNotFoundError("campaign not found", nil)
	}

	// Create campaign event
	event := NewCampaignEvent(campaignID, eventType, customerID, eventData)
	event.Revenue = revenue
	event.Cost = cost

	// Save event
	if err := s.eventRepo.Save(ctx, event); err != nil {
		return shared.NewInfrastructureError("failed to save campaign event", err)
	}

	// Update campaign metrics
	if err := campaign.TrackEvent(eventType, customerID, eventData); err != nil {
		return err
	}

	// Save updated campaign
	if err := s.repo.Save(ctx, campaign); err != nil {
		return shared.NewInfrastructureError("failed to save campaign", err)
	}

	// Update metrics repository
	metrics := campaign.Metrics()
	if err := s.metricsRepo.UpdateMetrics(ctx, campaignID, metrics); err != nil {
		return shared.NewInfrastructureError("failed to update campaign metrics", err)
	}

	// Publish domain events
	events := campaign.GetEvents()
	for _, event := range events {
		if err := s.eventBus.Publish(ctx, event); err != nil {
			fmt.Printf("Failed to publish event: %v\n", err)
		}
	}

	// Check for performance alerts
	if err := s.checkPerformanceAlerts(ctx, campaign); err != nil {
		fmt.Printf("Failed to check performance alerts: %v\n", err)
	}

	// Check for budget alerts
	if err := s.checkBudgetAlerts(ctx, campaign); err != nil {
		fmt.Printf("Failed to check budget alerts: %v\n", err)
	}

	return nil
}

// GetCampaignPerformanceReport generates a comprehensive performance report
func (s *CampaignService) GetCampaignPerformanceReport(
	ctx context.Context,
	campaignID CampaignID,
	period TimePeriod,
) (*PerformanceReport, error) {
	// Get campaign
	campaign, err := s.repo.FindByID(ctx, campaignID)
	if err != nil {
		return nil, shared.NewInfrastructureError("failed to find campaign", err)
	}
	if campaign == nil {
		return nil, shared.NewNotFoundError("campaign not found", nil)
	}

	// Get performance report from performance service
	report, err := s.performanceService.GetPerformanceReport(ctx, campaignID, period)
	if err != nil {
		return nil, shared.NewInfrastructureError("failed to generate performance report", err)
	}

	return report, nil
}

// checkPerformanceAlerts checks for performance-related alerts
func (s *CampaignService) checkPerformanceAlerts(ctx context.Context, campaign *Campaign) error {
	metrics := campaign.Metrics()

	// Check CTR alert
	if metrics.CTR < 1.0 && metrics.Impressions > 1000 {
		alert := PerformanceAlert{
			Type:         "low_ctr",
			Threshold:    1.0,
			CurrentValue: metrics.CTR,
			Message:      fmt.Sprintf("Campaign %s has low CTR: %.2f%%", campaign.Name(), metrics.CTR),
			Severity:     "medium",
		}
		if err := s.notificationService.SendPerformanceAlert(ctx, campaign.ID(), alert); err != nil {
			return err
		}
	}

	// Check conversion rate alert
	if metrics.ConversionRate < 2.0 && metrics.Clicks > 100 {
		alert := PerformanceAlert{
			Type:         "low_conversion",
			Threshold:    2.0,
			CurrentValue: metrics.ConversionRate,
			Message:      fmt.Sprintf("Campaign %s has low conversion rate: %.2f%%", campaign.Name(), metrics.ConversionRate),
			Severity:     "high",
		}
		if err := s.notificationService.SendPerformanceAlert(ctx, campaign.ID(), alert); err != nil {
			return err
		}
	}

	// Check ROI alert
	if metrics.ROI < 0 {
		alert := PerformanceAlert{
			Type:         "negative_roi",
			Threshold:    0.0,
			CurrentValue: metrics.ROI,
			Message:      fmt.Sprintf("Campaign %s has negative ROI: %.2f%%", campaign.Name(), metrics.ROI),
			Severity:     "critical",
		}
		if err := s.notificationService.SendPerformanceAlert(ctx, campaign.ID(), alert); err != nil {
			return err
		}
	}

	return nil
}

// checkBudgetAlerts checks for budget-related alerts
func (s *CampaignService) checkBudgetAlerts(ctx context.Context, campaign *Campaign) error {
	if campaign.Budget() == nil {
		return nil // No budget set, no alerts needed
	}

	metrics := campaign.Metrics()
	budget := campaign.Budget()

	// Check if approaching budget limit (90%)
	if campaign.IsApproachingBudgetLimit() {
		alert := BudgetAlert{
			Type:         "approaching_limit",
			Threshold:    budget.Amount * 0.9,
			CurrentValue: metrics.Cost.Amount,
			Message:      fmt.Sprintf("Campaign %s is approaching budget limit: %.2f of %.2f", campaign.Name(), metrics.Cost.Amount, budget.Amount),
			Severity:     "medium",
		}
		if err := s.notificationService.SendBudgetAlert(ctx, campaign.ID(), alert); err != nil {
			return err
		}
	}

	// Check if budget exceeded
	if campaign.HasExceededBudget() {
		alert := BudgetAlert{
			Type:         "exceeded",
			Threshold:    budget.Amount,
			CurrentValue: metrics.Cost.Amount,
			Message:      fmt.Sprintf("Campaign %s has exceeded budget: %.2f of %.2f", campaign.Name(), metrics.Cost.Amount, budget.Amount),
			Severity:     "critical",
		}
		if err := s.notificationService.SendBudgetAlert(ctx, campaign.ID(), alert); err != nil {
			return err
		}
	}

	return nil
}

// GetCampaign retrieves a campaign by ID
func (s *CampaignService) GetCampaign(ctx context.Context, campaignID CampaignID) (*Campaign, error) {
	campaign, err := s.repo.FindByID(ctx, campaignID)
	if err != nil {
		return nil, shared.NewInfrastructureError("failed to find campaign", err)
	}
	if campaign == nil {
		return nil, shared.NewNotFoundError("campaign not found", nil)
	}
	return campaign, nil
}

// ListCampaigns retrieves campaigns based on criteria
func (s *CampaignService) ListCampaigns(ctx context.Context, criteria ListCriteria) ([]*Campaign, error) {
	campaigns, err := s.repo.List(ctx, criteria)
	if err != nil {
		return nil, shared.NewInfrastructureError("failed to list campaigns", err)
	}
	return campaigns, nil
}

// CountCampaigns counts campaigns based on criteria
func (s *CampaignService) CountCampaigns(ctx context.Context, criteria ListCriteria) (int64, error) {
	count, err := s.repo.Count(ctx, criteria)
	if err != nil {
		return 0, shared.NewInfrastructureError("failed to count campaigns", err)
	}
	return count, nil
}

// SaveCampaign saves a campaign
func (s *CampaignService) SaveCampaign(ctx context.Context, campaign *Campaign) error {
	if err := s.repo.Save(ctx, campaign); err != nil {
		return shared.NewInfrastructureError("failed to save campaign", err)
	}
	return nil
}

// DeleteCampaign deletes a campaign
func (s *CampaignService) DeleteCampaign(ctx context.Context, campaignID CampaignID) error {
	// Get campaign to ensure it exists
	campaign, err := s.repo.FindByID(ctx, campaignID)
	if err != nil {
		return shared.NewInfrastructureError("failed to find campaign", err)
	}
	if campaign == nil {
		return shared.NewNotFoundError("campaign not found", nil)
	}

	// Delete campaign
	if err := s.repo.Delete(ctx, campaignID); err != nil {
		return shared.NewInfrastructureError("failed to delete campaign", err)
	}

	// Publish domain events
	events := campaign.GetEvents()
	for _, event := range events {
		if err := s.eventBus.Publish(ctx, event); err != nil {
			fmt.Printf("Failed to publish event: %v\n", err)
		}
	}

	return nil
}

// ExistsByName checks if a campaign with the given name exists
func (s *CampaignService) ExistsByName(ctx context.Context, name string) (bool, error) {
	exists, err := s.repo.ExistsByName(ctx, name)
	if err != nil {
		return false, shared.NewInfrastructureError("failed to check campaign name existence", err)
	}
	return exists, nil
}
