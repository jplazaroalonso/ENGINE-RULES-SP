package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/campaigns-management-service/internal/application/commands"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/campaigns-management-service/internal/application/queries"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/campaigns-management-service/internal/domain/campaign"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/campaigns-management-service/internal/infrastructure/messaging/nats"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/campaigns-management-service/internal/infrastructure/persistence/postgres"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/campaigns-management-service/internal/infrastructure/validation"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/campaigns-management-service/internal/interfaces/rest/handlers"
)

func main() {
	log.Println("Starting Campaigns Management Service...")

	// Load configuration from environment variables
	config := loadConfig()

	// Initialize database
	db, err := postgres.ConnectDatabase(postgres.DatabaseConfig{
		Host:     config.DatabaseHost,
		Port:     config.DatabasePort,
		User:     config.DatabaseUser,
		Password: config.DatabasePassword,
		DBName:   config.DatabaseName,
		SSLMode:  config.DatabaseSSLMode,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run database migrations
	if err := postgres.AutoMigrate(db); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}

	// Initialize NATS connection
	natsConn, err := nats.ConnectNATS(nats.NATSConfig{
		URL:      config.NATSURL,
		Username: config.NATSUsername,
		Password: config.NATSPassword,
		Token:    config.NATSToken,
	})
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer natsConn.Close()

	// Initialize infrastructure components
	campaignRepo := postgres.NewCampaignRepository(db)
	eventBus := nats.NewNATSEventBus(natsConn)
	validator := validation.NewStructValidator()

	// Initialize campaign service with minimal dependencies
	campaignService := campaign.NewCampaignService(
		campaignRepo,
		nil, // eventRepo - will be implemented later
		nil, // metricsRepo - will be implemented later
		nil, // targetingService - will be implemented later
		nil, // performanceService - will be implemented later
		nil, // schedulingService - will be implemented later
		nil, // notificationService - will be implemented later
		eventBus,
	)

	// Initialize application handlers
	createHandler := commands.NewCreateCampaignHandler(campaignService, validator)
	updateHandler := commands.NewUpdateCampaignHandler(campaignService, validator)
	activateHandler := commands.NewActivateCampaignHandler(campaignService, validator)
	pauseHandler := commands.NewPauseCampaignHandler(campaignService, validator)
	deleteHandler := commands.NewDeleteCampaignHandler(campaignService, validator)
	getHandler := queries.NewGetCampaignHandler(campaignService, validator)
	listHandler := queries.NewListCampaignsHandler(campaignService, validator)
	metricsHandler := queries.NewGetCampaignMetricsHandler(campaignService, validator)

	// Initialize REST handlers
	campaignHandler := handlers.NewCampaignHandler(
		createHandler,
		updateHandler,
		activateHandler,
		pauseHandler,
		deleteHandler,
		getHandler,
		listHandler,
		metricsHandler,
	)

	// Initialize Gin router
	router := setupRouter(campaignHandler)

	// Start HTTP server
	server := &http.Server{
		Addr:    ":" + config.Port,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", config.Port)
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
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

// Config holds application configuration
type Config struct {
	Port             string
	DatabaseHost     string
	DatabasePort     int
	DatabaseUser     string
	DatabasePassword string
	DatabaseName     string
	DatabaseSSLMode  string
	NATSURL          string
	NATSUsername     string
	NATSPassword     string
	NATSToken        string
	RulesServiceURL  string
}

// loadConfig loads configuration from environment variables
func loadConfig() Config {
	return Config{
		Port:             getEnv("PORT", "8080"),
		DatabaseHost:     getEnv("DB_HOST", "localhost"),
		DatabasePort:     getEnvAsInt("DB_PORT", 5432),
		DatabaseUser:     getEnv("DB_USER", "postgres"),
		DatabasePassword: getEnv("DB_PASSWORD", "password"),
		DatabaseName:     getEnv("DB_NAME", "campaigns_db"),
		DatabaseSSLMode:  getEnv("DB_SSLMODE", "disable"),
		NATSURL:          getEnv("NATS_URL", "nats://localhost:4222"),
		NATSUsername:     getEnv("NATS_USERNAME", ""),
		NATSPassword:     getEnv("NATS_PASSWORD", ""),
		NATSToken:        getEnv("NATS_TOKEN", ""),
		RulesServiceURL:  getEnv("RULES_SERVICE_URL", "http://localhost:8081"),
	}
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt gets an environment variable as integer with a default value
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// setupRouter configures the Gin router with all routes
func setupRouter(campaignHandler *handlers.CampaignHandler) *gin.Engine {
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "healthy",
			"service":   "campaigns-management-service",
			"timestamp": time.Now().UTC(),
		})
	})

	// API routes
	api := router.Group("/api/v1")
	{
		campaigns := api.Group("/campaigns")
		{
			campaigns.POST("", campaignHandler.CreateCampaign)
			campaigns.GET("", campaignHandler.ListCampaigns)
			campaigns.GET("/:id", campaignHandler.GetCampaign)
			campaigns.PUT("/:id", campaignHandler.UpdateCampaign)
			campaigns.POST("/:id/activate", campaignHandler.ActivateCampaign)
			campaigns.POST("/:id/pause", campaignHandler.PauseCampaign)
			campaigns.DELETE("/:id", campaignHandler.DeleteCampaign)
			campaigns.GET("/:id/metrics", campaignHandler.GetCampaignMetrics)
		}
	}

	return router
}
