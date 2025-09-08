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
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/analytics-dashboard-service/internal/application/commands"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/analytics-dashboard-service/internal/application/queries"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/analytics-dashboard-service/internal/domain/shared"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/analytics-dashboard-service/internal/infrastructure/config"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/analytics-dashboard-service/internal/infrastructure/messaging/nats"
	persistence "github.com/juanpablolazaro/ENGINE-RULES-SP/analytics-dashboard-service/internal/infrastructure/persistence/postgres"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/analytics-dashboard-service/internal/infrastructure/persistence/postgres/migrations"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/analytics-dashboard-service/internal/interfaces/rest/handlers"
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

	// Database setup
	db, err := gorm.Open(postgres.Open(cfg.Database.DSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	if err := migrations.ApplyMigrations(db); err != nil {
		log.Fatalf("failed to apply migrations: %v", err)
	}

	// Infrastructure
	dashboardRepo := persistence.NewDashboardRepository(db)
	reportRepo := persistence.NewReportRepository(db)
	metricRepo := persistence.NewMetricRepository(db)

	var eventBus shared.EventBus
	if cfg.NATS.URL != "" {
		publisher, err := nats.NewEventPublisher(cfg.NATS)
		if err != nil {
			log.Printf("Warning: failed to create event publisher: %v", err)
			eventBus = &nats.NoOpEventPublisher{}
		} else {
			eventBus = publisher
		}
	} else {
		log.Println("NATS URL not configured, event publishing disabled")
		eventBus = &nats.NoOpEventPublisher{}
	}

	// Application

	// Command handlers
	createDashboardHandler := commands.NewCreateDashboardHandler(dashboardRepo, eventBus)
	createReportHandler := commands.NewCreateReportHandler(reportRepo, eventBus)
	createMetricHandler := commands.NewCreateMetricHandler(metricRepo, eventBus)

	// Query handlers
	getDashboardHandler := queries.NewGetDashboardHandler(dashboardRepo)
	listDashboardsHandler := queries.NewListDashboardsHandler(dashboardRepo)
	getMetricsHandler := queries.NewGetMetricsHandler(metricRepo)

	// Interfaces
	dashboardHandler := handlers.NewDashboardHandler(createDashboardHandler, getDashboardHandler, listDashboardsHandler)
	reportHandler := handlers.NewReportHandler(createReportHandler)
	metricHandler := handlers.NewMetricHandler(createMetricHandler, getMetricsHandler)

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "analytics-dashboard-service",
			"version": "1.0.0",
		})
	})

	// Prometheus metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// API routes
	v1 := router.Group("/api/v1")
	{
		// Dashboard routes
		dashboards := v1.Group("/dashboards")
		{
			dashboards.GET("", dashboardHandler.ListDashboards)
			dashboards.POST("", dashboardHandler.CreateDashboard)
			dashboards.GET("/:id", dashboardHandler.GetDashboard)
			dashboards.PUT("/:id", dashboardHandler.UpdateDashboard)
			dashboards.DELETE("/:id", dashboardHandler.DeleteDashboard)
		}

		// Report routes
		reports := v1.Group("/reports")
		{
			reports.GET("", reportHandler.ListReports)
			reports.POST("", reportHandler.CreateReport)
			reports.GET("/:id", reportHandler.GetReport)
			reports.PUT("/:id", reportHandler.UpdateReport)
			reports.DELETE("/:id", reportHandler.DeleteReport)
			reports.POST("/:id/generate", reportHandler.GenerateReport)
		}

		// Metric routes
		metrics := v1.Group("/metrics")
		{
			metrics.GET("", metricHandler.ListMetrics)
			metrics.POST("", metricHandler.CreateMetric)
			metrics.GET("/:id", metricHandler.GetMetric)
			metrics.PUT("/:id", metricHandler.UpdateMetric)
			metrics.DELETE("/:id", metricHandler.DeleteMetric)
			metrics.GET("/:id/data", metricHandler.GetMetricData)
		}

		// Analytics routes
		analytics := v1.Group("/analytics")
		{
			analytics.GET("/real-time", dashboardHandler.GetRealTimeAnalytics)
			analytics.GET("/performance", dashboardHandler.GetPerformanceMetrics)
			analytics.GET("/business", dashboardHandler.GetBusinessMetrics)
			analytics.GET("/compliance", dashboardHandler.GetComplianceMetrics)
		}
	}

	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Graceful shutdown
	go func() {
		log.Printf("Starting Analytics Dashboard Service on port %s", cfg.Server.Port)
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
