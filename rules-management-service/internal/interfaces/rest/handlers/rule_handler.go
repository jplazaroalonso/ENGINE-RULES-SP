package handlers

import (
	"net/http"
	"strconv"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/rules-management-service/internal/application/commands"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/rules-management-service/internal/application/queries"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/rules-management-service/internal/domain/shared"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/rules-management-service/internal/interfaces/rest/dto"

	"github.com/gin-gonic/gin"
)

// RuleHandler handles HTTP requests for rule operations
type RuleHandler struct {
	createRuleHandler   *commands.CreateRuleHandler
	getRuleHandler      *queries.GetRuleHandler
	listRulesHandler    *queries.ListRulesHandler
	validateRuleHandler *commands.ValidateRuleHandler
}

func NewRuleHandler(
	createRuleHandler *commands.CreateRuleHandler,
	getRuleHandler *queries.GetRuleHandler,
	listRulesHandler *queries.ListRulesHandler,
	validateRuleHandler *commands.ValidateRuleHandler,
) *RuleHandler {
	return &RuleHandler{
		createRuleHandler:   createRuleHandler,
		getRuleHandler:      getRuleHandler,
		listRulesHandler:    listRulesHandler,
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

// ListRules handles GET /api/v1/rules
func (h *RuleHandler) ListRules(c *gin.Context) {
	// Parse query parameters
	query := queries.ListRulesQuery{
		Page:      parseIntParam(c, "page", 1),
		Limit:     parseIntParam(c, "limit", 20),
		SortBy:    c.Query("sort_by"),
		SortOrder: c.Query("sort_order"),
		Status:    c.Query("status"),
		Category:  c.Query("category"),
		Search:    c.Query("search"),
	}

	result, err := h.listRulesHandler.Handle(c.Request.Context(), query)
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

// parseIntParam parses an integer parameter from query string with a default value
func parseIntParam(c *gin.Context, param string, defaultValue int) int {
	value := c.Query(param)
	if value == "" {
		return defaultValue
	}
	
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	
	return parsed
}
