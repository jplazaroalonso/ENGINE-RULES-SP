package queries

import (
	"context"
	"fmt"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/domain/settings"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/domain/shared"
)

// GetOrganizationSettingQuery represents a query to get an organization setting by ID
type GetOrganizationSettingQuery struct {
	ID string `json:"id" validate:"required,uuid"`
}

// GetOrganizationSettingHandler handles the GetOrganizationSettingQuery
type GetOrganizationSettingHandler struct {
	organizationSettingRepo settings.OrganizationSettingRepository
	validator               shared.StructValidator
}

// NewGetOrganizationSettingHandler creates a new GetOrganizationSettingHandler
func NewGetOrganizationSettingHandler(
	organizationSettingRepo settings.OrganizationSettingRepository,
	validator shared.StructValidator,
) *GetOrganizationSettingHandler {
	return &GetOrganizationSettingHandler{
		organizationSettingRepo: organizationSettingRepo,
		validator:               validator,
	}
}

// Handle executes the GetOrganizationSettingQuery
func (h *GetOrganizationSettingHandler) Handle(ctx context.Context, query GetOrganizationSettingQuery) (*settings.OrganizationSetting, error) {
	// Validate the query
	if err := h.validator.Struct(query); err != nil {
		return nil, fmt.Errorf("%w: %s", shared.ErrValidation, err.Error())
	}

	// Parse organization setting ID
	organizationSettingID, err := settings.NewOrganizationSettingIDFromString(query.ID)
	if err != nil {
		return nil, err
	}

	// Find the organization setting
	organizationSetting, err := h.organizationSettingRepo.FindByID(ctx, organizationSettingID)
	if err != nil {
		return nil, fmt.Errorf("%w: organization setting not found: %v", shared.ErrNotFound, err)
	}

	return organizationSetting, nil
}
