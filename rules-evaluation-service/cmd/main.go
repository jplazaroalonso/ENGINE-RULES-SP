package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"rules-evaluation-service/internal/application"
	"rules-evaluation-service/internal/domain/evaluation"
	"rules-evaluation-service/internal/infrastructure/strategies"
	"rules-evaluation-service/internal/interfaces/rest/handlers"
)

func main() {
	// Initialize strategies
	promotionsStrategy := strategies.NewPromotionsStrategy()
	taxesStrategy := strategies.NewTaxesStrategy()

	strategyMap := map[string]evaluation.EvaluationStrategy{
		"PROMOTIONS": promotionsStrategy,
		"TAXES":      taxesStrategy,
		// Register other strategies here
	}

	// Initialize domain service
	evaluationService := evaluation.NewService(strategyMap)

	// Initialize application handler
	evaluateRuleHandler := application.NewEvaluateRuleHandler(evaluationService)

	// Initialize HTTP handler
	evaluationHandler := handlers.NewEvaluationHandler(evaluateRuleHandler)

	// Setup Gin router
	router := gin.Default()
	v1 := router.Group("/v1")
	{
		v1.POST("/evaluate", evaluationHandler.EvaluateRule)
	}

	// Start server on a different port than the management service
	if err := router.Run(":8081"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
