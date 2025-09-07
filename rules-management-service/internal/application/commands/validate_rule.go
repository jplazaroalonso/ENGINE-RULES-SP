package commands

import (
	"context"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/rules-management-service/internal/domain/rule"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/rules-management-service/internal/domain/shared"
)

// ValidateRuleCommand represents the command to validate a rule's DSL.
type ValidateRuleCommand struct {
	DSLContent string `json:"dsl_content" validate:"required"`
}

// ValidateRuleResult represents the result of a rule validation.
type ValidateRuleResult struct {
	IsValid bool     `json:"is_valid"`
	Errors  []string `json:"errors,omitempty"`
}

// ValidateRuleHandler handles rule validation commands.
type ValidateRuleHandler struct {
	validator         shared.Validator
	validationService rule.ValidationService
}

// NewValidateRuleHandler creates a new ValidateRuleHandler.
func NewValidateRuleHandler(validator shared.Validator, validationService rule.ValidationService) *ValidateRuleHandler {
	return &ValidateRuleHandler{
		validator:         validator,
		validationService: validationService,
	}
}

// Handle processes the validate rule command.
func (h *ValidateRuleHandler) Handle(ctx context.Context, cmd ValidateRuleCommand) (*ValidateRuleResult, error) {
	if err := h.validator.Validate(cmd); err != nil {
		return nil, shared.NewValidationError("invalid validate rule command", err)
	}

	isValid, errors := h.validationService.Validate(cmd.DSLContent)

	return &ValidateRuleResult{
		IsValid: isValid,
		Errors:  errors,
	}, nil
}
