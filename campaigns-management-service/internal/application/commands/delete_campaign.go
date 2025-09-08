package commands

import (
	"context"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/campaigns-management-service/internal/domain/campaign"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/campaigns-management-service/internal/domain/shared"
)

// DeleteCampaignCommand represents the command to delete a campaign
type DeleteCampaignCommand struct {
	CampaignID string `json:"campaignId" validate:"required"`
	DeletedBy  string `json:"deletedBy" validate:"required"`
	Reason     string `json:"reason,omitempty" validate:"omitempty,max=200"`
}

// DeleteCampaignResult represents the result of deleting a campaign
type DeleteCampaignResult struct {
	CampaignID string `json:"campaignId"`
	Name       string `json:"name"`
	DeletedBy  string `json:"deletedBy"`
	DeletedAt  string `json:"deletedAt"`
	Reason     string `json:"reason,omitempty"`
}

// DeleteCampaignHandler handles campaign deletion commands
type DeleteCampaignHandler struct {
	campaignService *campaign.CampaignService
	validator       shared.Validator
}

// NewDeleteCampaignHandler creates a new DeleteCampaignHandler
func NewDeleteCampaignHandler(
	campaignService *campaign.CampaignService,
	validator shared.Validator,
) *DeleteCampaignHandler {
	return &DeleteCampaignHandler{
		campaignService: campaignService,
		validator:       validator,
	}
}

// Handle processes the delete campaign command
func (h *DeleteCampaignHandler) Handle(ctx context.Context, cmd DeleteCampaignCommand) (*DeleteCampaignResult, error) {
	// Validate command input
	if err := h.validator.Validate(cmd); err != nil {
		return nil, shared.NewValidationError("invalid delete campaign command", err)
	}

	// Parse campaign ID
	campaignID, err := campaign.NewCampaignIDFromString(cmd.CampaignID)
	if err != nil {
		return nil, shared.NewValidationError("invalid campaign ID", err)
	}

	// Get campaign to return name in result
	campaign, err := h.campaignService.GetCampaign(ctx, campaignID)
	if err != nil {
		return nil, err
	}

	// Check if campaign can be deleted
	if campaign.Status() == campaign.CampaignStatusActive {
		return nil, shared.NewBusinessError("cannot delete active campaign", nil)
	}

	// Delete campaign using domain service
	if err := h.campaignService.DeleteCampaign(ctx, campaignID); err != nil {
		return nil, err
	}

	return &DeleteCampaignResult{
		CampaignID: campaign.ID().String(),
		Name:       campaign.Name(),
		DeletedBy:  cmd.DeletedBy,
		DeletedAt:  campaign.UpdatedAt().Format("2006-01-02T15:04:05Z"),
		Reason:     cmd.Reason,
	}, nil
}
