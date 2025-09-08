package commands

import (
	"context"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/campaigns-management-service/internal/domain/campaign"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/campaigns-management-service/internal/domain/shared"
)

// PauseCampaignCommand represents the command to pause a campaign
type PauseCampaignCommand struct {
	CampaignID string `json:"campaignId" validate:"required"`
	PausedBy   string `json:"pausedBy" validate:"required"`
}

// PauseCampaignResult represents the result of pausing a campaign
type PauseCampaignResult struct {
	CampaignID string `json:"campaignId"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	PausedBy   string `json:"pausedBy"`
	PausedAt   string `json:"pausedAt"`
}

// PauseCampaignHandler handles campaign pause commands
type PauseCampaignHandler struct {
	campaignService *campaign.CampaignService
	validator       shared.Validator
}

// NewPauseCampaignHandler creates a new PauseCampaignHandler
func NewPauseCampaignHandler(
	campaignService *campaign.CampaignService,
	validator shared.Validator,
) *PauseCampaignHandler {
	return &PauseCampaignHandler{
		campaignService: campaignService,
		validator:       validator,
	}
}

// Handle processes the pause campaign command
func (h *PauseCampaignHandler) Handle(ctx context.Context, cmd PauseCampaignCommand) (*PauseCampaignResult, error) {
	// Validate command input
	if err := h.validator.Validate(cmd); err != nil {
		return nil, shared.NewValidationError("invalid pause campaign command", err)
	}

	// Parse campaign ID
	campaignID, err := campaign.NewCampaignIDFromString(cmd.CampaignID)
	if err != nil {
		return nil, shared.NewValidationError("invalid campaign ID", err)
	}

	// Parse user ID
	userID, err := shared.NewUserIDFromString(cmd.PausedBy)
	if err != nil {
		return nil, shared.NewValidationError("invalid user ID", err)
	}

	// Pause campaign using domain service
	if err := h.campaignService.PauseCampaign(ctx, campaignID, userID); err != nil {
		return nil, err
	}

	// Get updated campaign to return result
	updatedCampaign, err := h.campaignService.GetCampaign(ctx, campaignID)
	if err != nil {
		return nil, err
	}

	return &PauseCampaignResult{
		CampaignID: updatedCampaign.ID().String(),
		Name:       updatedCampaign.Name(),
		Status:     updatedCampaign.Status().String(),
		PausedBy:   cmd.PausedBy,
		PausedAt:   updatedCampaign.UpdatedAt().Format("2006-01-02T15:04:05Z"),
	}, nil
}
