package commands

import (
	"context"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/campaigns-management-service/internal/domain/campaign"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/campaigns-management-service/internal/domain/shared"
)

// ActivateCampaignCommand represents the command to activate a campaign
type ActivateCampaignCommand struct {
	CampaignID  string `json:"campaignId" validate:"required"`
	ActivatedBy string `json:"activatedBy" validate:"required"`
}

// ActivateCampaignResult represents the result of activating a campaign
type ActivateCampaignResult struct {
	CampaignID  string `json:"campaignId"`
	Name        string `json:"name"`
	Status      string `json:"status"`
	ActivatedBy string `json:"activatedBy"`
	ActivatedAt string `json:"activatedAt"`
}

// ActivateCampaignHandler handles campaign activation commands
type ActivateCampaignHandler struct {
	campaignService *campaign.CampaignService
	validator       shared.Validator
}

// NewActivateCampaignHandler creates a new ActivateCampaignHandler
func NewActivateCampaignHandler(
	campaignService *campaign.CampaignService,
	validator shared.Validator,
) *ActivateCampaignHandler {
	return &ActivateCampaignHandler{
		campaignService: campaignService,
		validator:       validator,
	}
}

// Handle processes the activate campaign command
func (h *ActivateCampaignHandler) Handle(ctx context.Context, cmd ActivateCampaignCommand) (*ActivateCampaignResult, error) {
	// Validate command input
	if err := h.validator.Validate(cmd); err != nil {
		return nil, shared.NewValidationError("invalid activate campaign command", err)
	}

	// Parse campaign ID
	campaignID, err := campaign.NewCampaignIDFromString(cmd.CampaignID)
	if err != nil {
		return nil, shared.NewValidationError("invalid campaign ID", err)
	}

	// Parse user ID
	userID, err := shared.NewUserIDFromString(cmd.ActivatedBy)
	if err != nil {
		return nil, shared.NewValidationError("invalid user ID", err)
	}

	// Activate campaign using domain service
	if err := h.campaignService.ActivateCampaign(ctx, campaignID, userID); err != nil {
		return nil, err
	}

	// Get updated campaign to return result
	updatedCampaign, err := h.campaignService.GetCampaign(ctx, campaignID)
	if err != nil {
		return nil, err
	}

	return &ActivateCampaignResult{
		CampaignID:  updatedCampaign.ID().String(),
		Name:        updatedCampaign.Name(),
		Status:      updatedCampaign.Status().String(),
		ActivatedBy: cmd.ActivatedBy,
		ActivatedAt: updatedCampaign.UpdatedAt().Format("2006-01-02T15:04:05Z"),
	}, nil
}
