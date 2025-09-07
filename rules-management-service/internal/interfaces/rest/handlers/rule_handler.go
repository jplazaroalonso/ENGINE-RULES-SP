package handlers

import (
	"net/http"

	"rules-management-service/internal/application/commands"
	"rules-management-service/internal/application/queries"
	"rules-management-service/internal/domain/shared"
	"rules-management-service/internal/interfaces/rest/dto"

	"github.com/gin-gonic/gin"
)

// RuleHandler handles HTTP requests for rule operations
type RuleHandler struct {
	createRuleHandler   *commands.CreateRuleHandler
	getRuleHandler      *queries.GetRuleHandler
	validateRuleHandler *commands.ValidateRuleHandler
}

func NewRuleHandler(
	createRuleHandler *commands.CreateRuleHandler,
	getRuleHandler *queries.GetRuleHandler,
	validateRuleHandler *commands.ValidateRuleHandler,
) *RuleHandler {
	return &RuleHandler{
		createRuleHandler:   createRuleHandler,
		getRuleHandler:      getRuleHandler,
		validateRuleHandler: validateRuleHandler,
	}
}

// CreateRule handles POST /api/v1/rules
func (h *RuleHandler) CreateRule(c *gin.Context) {
	var req dto.CreateRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid request body",
			Message: err.Error(),
		})
		return
	}

	// In a real app, user would come from auth middleware
	const createdBy = "user@example.com"

	cmd := commands.CreateRuleCommand{
		Name:        req.Name,
		Description: req.Description,
		DSLContent:  req.DSLContent,
		Priority:    req.Priority,
		CreatedBy:   createdBy,
		Category:    req.Category,
		Tags:        req.Tags,
	}

	result, err := h.createRuleHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, result)
}

// ValidateRule handles POST /api/v1/rules/validate
func (h *RuleHandler) ValidateRule(c *gin.Context) {
	var req dto.ValidateRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid request body",
			Message: err.Error(),
		})
		return
	}

	cmd := commands.ValidateRuleCommand{
		DSLContent: req.DSLContent,
	}

	result, err := h.validateRuleHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetRule handles GET /api/v1/rules/:id
func (h *RuleHandler) GetRule(c *gin.Context) {
	ruleID := c.Param("id")

	query := queries.GetRuleQuery{RuleID: ruleID}
	result, err := h.getRuleHandler.Handle(c.Request.Context(), query)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// handleError maps domain errors to appropriate HTTP responses
func handleError(c *gin.Context, err error) {
	switch err.(type) {
	case *shared.ValidationError:
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Validation failed", Message: err.Error()})
	case *shared.BusinessError:
		c.JSON(http.StatusConflict, dto.ErrorResponse{Error: "Business rule violation", Message: err.Error()})
	case *shared.NotFoundError:
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "Resource not found", Message: err.Error()})
	case *shared.InfrastructureError:
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Internal server error"})
	default:
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "An unexpected error occurred"})
	}
}
