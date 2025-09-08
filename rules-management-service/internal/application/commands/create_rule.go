package commands

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/rules-management-service/internal/domain/rule"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/rules-management-service/internal/domain/shared"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/rules-management-service/internal/infrastructure/telemetry"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// CreateRuleCommand represents the command to create a new rule
type CreateRuleCommand struct {
	Name        string   `json:"name" validate:"required,min=3,max=100"`
	Description string   `json:"description" validate:"max=500"`
	DSLContent  string   `json:"dsl_content" validate:"required"`
	Priority    string   `json:"priority" validate:"required,oneof=LOW MEDIUM HIGH CRITICAL"`
	CreatedBy   string   `json:"created_by" validate:"required"`
	Category    string   `json:"category"`
	Tags        []string `json:"tags"`
}

// CreateRuleResult represents the result of creating a rule
type CreateRuleResult struct {
	RuleID  string `json:"rule_id"`
	Name    string `json:"name"`
	Status  string `json:"status"`
	Version int    `json:"version"`
}

// CreateRuleHandler handles rule creation commands
type CreateRuleHandler struct {
	ruleRepo           rule.Repository
	validator          shared.Validator
	eventBus           shared.EventBus
	replicationEnabled bool
	validationService  rule.ValidationService
}

func NewCreateRuleHandler(
	ruleRepo rule.Repository,
	validator shared.Validator,
	eventBus shared.EventBus,
	replicationEnabled bool,
	validationService rule.ValidationService,
) *CreateRuleHandler {
	return &CreateRuleHandler{
		ruleRepo:           ruleRepo,
		validator:          validator,
		eventBus:           eventBus,
		replicationEnabled: replicationEnabled,
		validationService:  validationService,
	}
}

// Handle processes the create rule command
func (h *CreateRuleHandler) Handle(ctx context.Context, cmd CreateRuleCommand) (*CreateRuleResult, error) {
	tr := otel.Tracer("application")
	ctx, span := tr.Start(ctx, "CreateRuleHandler.Handle")
	defer span.End()

	telemetry.RulesCreated.Inc()
	startTime := time.Now()
	defer func() {
		telemetry.RuleCreationDuration.Observe(time.Since(startTime).Seconds())
	}()

	span.SetAttributes(
		attribute.String("rule.name", cmd.Name),
		attribute.String("rule.category", cmd.Category),
	)

	// Validate command input
	if err := h.validator.Validate(cmd); err != nil {
		return nil, shared.NewValidationError("invalid create rule command", err)
	}

	// Validate DSL content
	isValid, validationErrors := h.validationService.Validate(cmd.DSLContent)
	if !isValid {
		// In a real app, you might want a more structured error here.
		return nil, shared.NewValidationError("DSL validation failed", errors.New(strings.Join(validationErrors, "; ")))
	}

	exists, err := h.ruleRepo.ExistsByName(ctx, cmd.Name)
	if err != nil {
		return nil, shared.NewInfrastructureError("failed to check rule existence", err)
	}
	if exists {
		return nil, shared.NewBusinessError("rule name already exists", nil)
	}

	priority := rule.Priority(cmd.Priority)

	newRule, err := rule.NewRule(cmd.Name, cmd.Description, cmd.DSLContent, cmd.CreatedBy, priority, cmd.Category, cmd.Tags)
	if err != nil {
		return nil, err // Domain error
	}

	if err := h.ruleRepo.Save(ctx, newRule); err != nil {
		return nil, shared.NewInfrastructureError("failed to save rule", err)
	}

	// Publish event if replication is enabled
	if h.replicationEnabled {
		event := rule.RuleCreatedEvent{
			RuleID:    newRule.ID().String(),
			Name:      newRule.Name(),
			CreatedBy: newRule.CreatedBy(),
			CreatedAt: newRule.CreatedAt(),
		}
		if err := h.eventBus.Publish(event); err != nil {
			// In a real app, you might want to handle this more gracefully,
			// e.g., using an outbox pattern to ensure the event is eventually sent.
			// For now, we'll just log it.
			log.Printf("Warning: failed to publish RuleCreatedEvent: %v", err)
		}
	}

	return &CreateRuleResult{
		RuleID:  newRule.ID().String(),
		Name:    newRule.Name(),
		Status:  string(newRule.Status()),
		Version: newRule.Version(),
	}, nil
}
