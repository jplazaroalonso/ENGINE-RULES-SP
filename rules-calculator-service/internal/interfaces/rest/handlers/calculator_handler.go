package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/rules-calculator-service/internal/application"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/rules-calculator-service/internal/interfaces/rest/dto"
)

// CalculatorHandler handles HTTP requests for the calculator service.
type CalculatorHandler struct {
	handler *application.CalculateRulesHandler
}

// NewCalculatorHandler creates a new CalculatorHandler.
func NewCalculatorHandler(handler *application.CalculateRulesHandler) *CalculatorHandler {
	return &CalculatorHandler{
		handler: handler,
	}
}

// Calculate handles a request to calculate rules.
func (h *CalculatorHandler) Calculate(c *gin.Context) {
	var req dto.CalculationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	cmd := application.CalculateRulesCommand{
		RuleIDs: req.RuleIDs,
		Context: req.Context,
	}

	result, err := h.handler.Handle(c.Request.Context(), cmd)
	if err != nil {
		// This is a simplification. In a real app, we'd have more nuanced error handling.
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.CalculationResponse{
		CalculationID: result.CalculationID,
		Value:         result.Value,
		Breakdown:     result.Breakdown,
	})
}
