package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	// "go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"rules-evaluation-service/internal/application"
	"rules-evaluation-service/internal/domain/evaluation"
	"rules-evaluation-service/internal/infrastructure/config"
	"rules-evaluation-service/internal/infrastructure/strategies"

	// "rules-evaluation-service/internal/infrastructure/telemetry"
	"rules-evaluation-service/internal/interfaces/rest/handlers"
)

func main() {
	cfg := config.DefaultConfig()

	// Setup OpenTelemetry - temporarily disabled
	// tp, err := telemetry.InitTracer(cfg.Telemetry.ServiceName, cfg.Telemetry.Exporter)
	// if err != nil {
	// 	log.Fatalf("failed to initialize tracer: %v", err)
	// }
	// defer func() {
	// 	if err := tp.Shutdown(context.Background()); err != nil {
	// 		log.Printf("Error shutting down tracer provider: %v", err)
	// 	}
	// }()

	// Infrastructure
	promotionsStrategy := &strategies.PromotionsStrategy{}
	taxesStrategy := &strategies.TaxesStrategy{}

	// Domain
	evaluationService := evaluation.NewService(map[string]evaluation.EvaluationStrategy{
		"PROMOTIONS": promotionsStrategy,
		"TAXES":      taxesStrategy,
	})

	// Application
	evaluateRuleHandler := application.NewEvaluateRuleHandler(evaluationService)

	// Interfaces
	evaluationHandler := handlers.NewEvaluationHandler(evaluateRuleHandler)

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	// router.Use(otelgin.Middleware(cfg.Telemetry.ServiceName))

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "service": "rules-evaluation-service"})
	})

	// Prometheus metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	v1 := router.Group("/v1")
	{
		v1.POST("/evaluate", evaluationHandler.EvaluateRule)
	}

	// API Gateway routes
	apiV1 := router.Group("/api/v1")
	{
		apiV1.POST("/evaluate", evaluationHandler.EvaluateRule)
	}

	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		log.Printf("Starting Rules Evaluation Service on port %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
