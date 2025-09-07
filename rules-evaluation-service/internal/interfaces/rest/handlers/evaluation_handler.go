package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"rules-evaluation-service/internal/application"
	"rules-evaluation-service/internal/interfaces/rest/dto"
)

// EvaluationHandler handles HTTP requests for rule evaluation.
type EvaluationHandler struct {
	evaluateRuleHandler *application.EvaluateRuleHandler
}

func NewEvaluationHandler(evaluateRuleHandler *application.EvaluateRuleHandler) *EvaluationHandler {
	return &EvaluationHandler{evaluateRuleHandler: evaluateRuleHandler}
}

// EvaluateRule handles POST /v1/evaluate
func (h *EvaluationHandler) EvaluateRule(c *gin.Context) {
	var req dto.EvaluationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid request body",
			Message: err.Error(),
		})
		return
	}

	cmd := application.EvaluateRuleCommand{
		RuleCategory: req.RuleCategory,
		DSLContent:   req.DSLContent,
		Context:      req.Context,
	}

	result, err := h.evaluateRuleHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		// In a real app, you would have more sophisticated error handling.
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.EvaluationResponse{Result: result.Result})
}
