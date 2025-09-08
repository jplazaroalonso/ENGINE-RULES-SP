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
	"github.com/go-playground/validator/v10"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	// "go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/rules-management-service/internal/application/commands"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/rules-management-service/internal/application/queries"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/rules-management-service/internal/infrastructure/config"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/rules-management-service/internal/infrastructure/dsl"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/rules-management-service/internal/infrastructure/messaging/nats"
	persistence "github.com/juanpablolazaro/ENGINE-RULES-SP/rules-management-service/internal/infrastructure/persistence/postgres"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/rules-management-service/internal/infrastructure/persistence/postgres/migrations"

	// "github.com/juanpablolazaro/ENGINE-RULES-SP/rules-management-service/internal/infrastructure/telemetry"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/rules-management-service/internal/domain/shared"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/rules-management-service/internal/infrastructure/validation"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/rules-management-service/internal/interfaces/rest/handlers"
)

// AppValidator wraps the go-playground/validator.
type AppValidator struct {
	validate *validator.Validate
}

func (v *AppValidator) Validate(s interface{}) error {
	return v.validate.Struct(s)
}

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

	// Database setup
	db, err := gorm.Open(postgres.Open(cfg.Database.DSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	if err := migrations.ApplyMigrations(db); err != nil {
		log.Fatalf("failed to apply migrations: %v", err)
	}

	// Infrastructure
	ruleRepo := persistence.NewRuleRepository(db)
	var eventPublisher shared.EventBus
	if cfg.NATS.URL != "" {
		publisher, err := nats.NewEventPublisher(cfg.NATS)
		if err != nil {
			log.Printf("Warning: failed to create event publisher: %v", err)
			eventPublisher = &nats.NoOpEventPublisher{}
		} else {
			eventPublisher = publisher
		}
	} else {
		log.Println("NATS URL not configured, event publishing disabled")
		eventPublisher = &nats.NoOpEventPublisher{}
	}

	// Application
	validator := validation.NewStructValidator()
	validationService := dsl.NewSimpleValidator()
	createRuleHandler := commands.NewCreateRuleHandler(ruleRepo, validator, eventPublisher, true, validationService) // Assuming replication is enabled
	getRuleHandler := queries.NewGetRuleHandler(ruleRepo)
	listRulesHandler := queries.NewListRulesHandler(ruleRepo)
	validateRuleHandler := commands.NewValidateRuleHandler(validator, validationService)

	// Interfaces
	ruleHandler := handlers.NewRuleHandler(createRuleHandler, getRuleHandler, listRulesHandler, validateRuleHandler)

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	// router.Use(otelgin.Middleware(cfg.Telemetry.ServiceName))

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "service": "rules-management-service", "version": "2.0"})
	})
	
	// Test endpoint to verify code updates
	router.GET("/test", func(c *gin.Context) {
		log.Println("TEST ENDPOINT CALLED - CODE IS UPDATED")
		c.JSON(200, gin.H{"message": "test endpoint working", "timestamp": "2025-09-07-21:30"})
	})

	// Prometheus metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	v1 := router.Group("/v1")
	{
		log.Println("Registering v1.GET /rules route")
		v1.GET("/rules", ruleHandler.ListRules)
		v1.POST("/rules", ruleHandler.CreateRule)
		v1.GET("/rules/:id", ruleHandler.GetRule)
	}

	// API Gateway routes
	apiV1 := router.Group("/api/v1")
	{
		log.Println("Registering apiV1.GET /rules route")
		apiV1.GET("/rules", ruleHandler.ListRules)
		apiV1.POST("/rules", ruleHandler.CreateRule)
		apiV1.GET("/rules/:id", ruleHandler.GetRule)
	}

	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		log.Printf("Starting Rules Management Service on port %s", cfg.Server.Port)
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
