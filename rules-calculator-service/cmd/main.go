package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/rules-calculator-service/internal/application"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/rules-calculator-service/internal/infrastructure/adapters"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/rules-calculator-service/internal/interfaces/rest/handlers"
)

func main() {
	// Infrastructure
	// In a real application, the URL for the evaluation service would come from config.
	ruleEvaluator := adapters.NewHTTPEvaluationAdapter("http://localhost:8081")

	// Application
	calculateHandler := application.NewCalculateRulesHandler(ruleEvaluator)

	// Interfaces
	httpHandler := handlers.NewCalculatorHandler(calculateHandler)

	router := gin.Default()
	v1 := router.Group("/v1")
	{
		v1.POST("/calculate", httpHandler.Calculate)
	}

	log.Println("Starting Rules Calculator Service on port 8082")
	if err := router.Run(":8082"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
