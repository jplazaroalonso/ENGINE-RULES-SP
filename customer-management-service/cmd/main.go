package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/customer-management-service/internal/infrastructure/external"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/customer-management-service/internal/infrastructure/messaging/nats"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/customer-management-service/internal/infrastructure/persistence/postgres"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/customer-management-service/internal/infrastructure/validation"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/customer-management-service/internal/interfaces/rest/handlers"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/customer-management-service/internal/interfaces/rest/middleware"
)

func main() {
	// Load configuration from environment variables
	config := loadConfig()

	// Initialize database connection
	db, err := postgres.ConnectToPostgreSQL(config.DatabaseDSN)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run database migrations
	if err := postgres.AutoMigrate(db); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}

	// Initialize NATS connection
	natsConn, err := nats.ConnectToNATS(config.NATSURL)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer natsConn.Close()

	// Initialize JetStream
	if err := nats.InitializeJetStream(natsConn); err != nil {
		log.Fatalf("Failed to initialize JetStream: %v", err)
	}

	// Initialize event bus
	eventBus, err := nats.NewEventBus(natsConn)
	if err != nil {
		log.Fatalf("Failed to create event bus: %v", err)
	}

	// Initialize repositories
	customerRepo := postgres.NewCustomerRepository(db)

	// Initialize external clients
	rulesClient := external.NewRulesClient(config.RulesEngineURL, config.RulesEngineAPIKey)

	// Initialize validators
	validator := validation.NewValidator()

	// Initialize handlers
	customerHandler := handlers.NewCustomerHandler(customerRepo, eventBus, validator, rulesClient)

	// Initialize Gin router
	router := setupRouter(customerHandler)

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", config.Port),
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting customer management service on port %s", config.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Give outstanding requests 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}

// Config represents the application configuration
type Config struct {
	Port              string
	DatabaseDSN       string
	NATSURL           string
	RulesEngineURL    string
	RulesEngineAPIKey string
}

// loadConfig loads configuration from environment variables
func loadConfig() *Config {
	return &Config{
		Port:              getEnv("PORT", "8080"),
		DatabaseDSN:       getEnv("DATABASE_DSN", "postgres://postgres:password@localhost:5432/customer_management?sslmode=disable"),
		NATSURL:           getEnv("NATS_URL", "nats://localhost:4222"),
		RulesEngineURL:    getEnv("RULES_ENGINE_URL", "http://localhost:8081"),
		RulesEngineAPIKey: getEnv("RULES_ENGINE_API_KEY", ""),
	}
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// setupRouter sets up the Gin router with middleware and routes
func setupRouter(customerHandler *handlers.CustomerHandler) *gin.Engine {
	// Set Gin mode
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Add middleware
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())
	router.Use(middleware.CORS())
	router.Use(middleware.RequestID())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().UTC(),
			"service":   "customer-management-service",
		})
	})

	// API routes
	api := router.Group("/api/v1")
	{
		// Customer routes
		customers := api.Group("/customers")
		{
			customers.GET("", customerHandler.ListCustomers)
			customers.POST("", customerHandler.CreateCustomer)
			customers.GET("/:id", customerHandler.GetCustomer)
			customers.PUT("/:id", customerHandler.UpdateCustomer)
			customers.DELETE("/:id", customerHandler.DeleteCustomer)

			// Customer analytics routes
			customers.GET("/:id/analytics", customerHandler.GetCustomerAnalytics)
			customers.GET("/:id/insights", customerHandler.GetCustomerInsights)
			customers.POST("/:id/track", customerHandler.TrackCustomerEvent)
			customers.GET("/:id/segments", customerHandler.GetCustomerSegments)

			// Customer privacy & GDPR routes
			customers.GET("/:id/data", customerHandler.ExportCustomerData)
			customers.DELETE("/:id/data", customerHandler.DeleteCustomerData)
			customers.PUT("/:id/consent", customerHandler.UpdatePrivacyConsent)
			customers.GET("/:id/consent", customerHandler.GetPrivacyConsent)
			customers.POST("/:id/anonymize", customerHandler.AnonymizeCustomerData)
		}

		// Customer segment routes
		segments := api.Group("/customers/segments")
		{
			segments.GET("", customerHandler.ListSegments)
			segments.POST("", customerHandler.CreateSegment)
			segments.GET("/:id", customerHandler.GetSegment)
			segments.PUT("/:id", customerHandler.UpdateSegment)
			segments.DELETE("/:id", customerHandler.DeleteSegment)
			segments.POST("/:id/calculate", customerHandler.CalculateSegment)
			segments.GET("/:id/customers", customerHandler.GetSegmentCustomers)
		}

		// Bulk operations routes
		bulk := api.Group("/customers/bulk")
		{
			bulk.POST("/update", customerHandler.BulkUpdateCustomers)
			bulk.POST("/delete", customerHandler.BulkDeleteCustomers)
			bulk.POST("/segments", customerHandler.BulkAssignSegments)
		}

		// Import/Export routes
		importExport := api.Group("/customers")
		{
			importExport.GET("/export", customerHandler.ExportCustomers)
			importExport.POST("/import", customerHandler.ImportCustomers)
		}

		segmentsImportExport := api.Group("/customers/segments")
		{
			segmentsImportExport.GET("/export", customerHandler.ExportSegments)
			segmentsImportExport.POST("/import", customerHandler.ImportSegments)
		}
	}

	return router
}
