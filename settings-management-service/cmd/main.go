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
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/application/commands"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/application/queries"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/infrastructure/cache/redis"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/infrastructure/messaging/nats"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/infrastructure/persistence/postgres"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/infrastructure/validation"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/interfaces/rest/handlers"
)

func main() {
	log.Println("Starting Settings Management Service...")

	// Initialize infrastructure
	db, err := initializeDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer func() {
		if err := postgres.Close(db); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	natsConn, err := initializeNATS()
	if err != nil {
		log.Fatalf("Failed to initialize NATS: %v", err)
	}
	defer func() {
		if err := nats.Close(natsConn); err != nil {
			log.Printf("Error closing NATS connection: %v", err)
		}
	}()

	redisClient, err := initializeRedis()
	if err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}
	defer func() {
		if err := redis.Close(redisClient); err != nil {
			log.Printf("Error closing Redis connection: %v", err)
		}
	}()

	// Initialize event bus
	eventBus, err := nats.NewEventBus(natsConn)
	if err != nil {
		log.Fatalf("Failed to initialize event bus: %v", err)
	}
	defer func() {
		if err := eventBus.Close(); err != nil {
			log.Printf("Error closing event bus: %v", err)
		}
	}()

	// Initialize repositories
	configRepo := postgres.NewConfigurationRepository(db)
	featureFlagRepo := postgres.NewFeatureFlagRepository(db)
	userPreferenceRepo := postgres.NewUserPreferenceRepository(db)
	organizationSettingRepo := postgres.NewOrganizationSettingRepository(db)
	cacheRepo := redis.NewCacheRepository(redisClient, "settings")

	// Initialize validator
	validator := validation.NewStructValidator()

	// Initialize command handlers
	createConfigHandler := commands.NewCreateConfigurationHandler(configRepo, eventBus, validator)
	updateConfigHandler := commands.NewUpdateConfigurationHandler(configRepo, eventBus, validator)
	deleteConfigHandler := commands.NewDeleteConfigurationHandler(configRepo, eventBus, validator)
	createFeatureFlagHandler := commands.NewCreateFeatureFlagHandler(featureFlagRepo, eventBus, validator)
	updateFeatureFlagHandler := commands.NewUpdateFeatureFlagHandler(featureFlagRepo, eventBus, validator)
	deleteFeatureFlagHandler := commands.NewDeleteFeatureFlagHandler(featureFlagRepo, eventBus, validator)
	createUserPreferenceHandler := commands.NewCreateUserPreferenceHandler(userPreferenceRepo, eventBus, validator)
	updateUserPreferenceHandler := commands.NewUpdateUserPreferenceHandler(userPreferenceRepo, eventBus, validator)
	deleteUserPreferenceHandler := commands.NewDeleteUserPreferenceHandler(userPreferenceRepo, eventBus, validator)
	createOrganizationSettingHandler := commands.NewCreateOrganizationSettingHandler(organizationSettingRepo, eventBus, validator)
	updateOrganizationSettingHandler := commands.NewUpdateOrganizationSettingHandler(organizationSettingRepo, eventBus, validator)
	deleteOrganizationSettingHandler := commands.NewDeleteOrganizationSettingHandler(organizationSettingRepo, eventBus, validator)

	// Initialize query handlers
	getConfigHandler := queries.NewGetConfigurationHandler(configRepo, validator)
	listConfigsHandler := queries.NewListConfigurationsHandler(configRepo, validator)
	getFeatureFlagHandler := queries.NewGetFeatureFlagHandler(featureFlagRepo, validator)
	listFeatureFlagsHandler := queries.NewListFeatureFlagsHandler(featureFlagRepo, validator)
	getUserPreferenceHandler := queries.NewGetUserPreferenceHandler(userPreferenceRepo, validator)
	listUserPreferencesHandler := queries.NewListUserPreferencesHandler(userPreferenceRepo, validator)
	getOrganizationSettingHandler := queries.NewGetOrganizationSettingHandler(organizationSettingRepo, validator)
	listOrganizationSettingsHandler := queries.NewListOrganizationSettingsHandler(organizationSettingRepo, validator)

	// Initialize REST handlers
	configHandler := handlers.NewConfigurationHandler(
		createConfigHandler,
		updateConfigHandler,
		deleteConfigHandler,
		getConfigHandler,
		listConfigsHandler,
	)
	featureFlagHandler := handlers.NewFeatureFlagHandler(
		createFeatureFlagHandler,
		updateFeatureFlagHandler,
		deleteFeatureFlagHandler,
		getFeatureFlagHandler,
		listFeatureFlagsHandler,
	)
	userPreferenceHandler := handlers.NewUserPreferenceHandler(
		createUserPreferenceHandler,
		updateUserPreferenceHandler,
		deleteUserPreferenceHandler,
		getUserPreferenceHandler,
		listUserPreferencesHandler,
	)
	organizationSettingHandler := handlers.NewOrganizationSettingHandler(
		createOrganizationSettingHandler,
		updateOrganizationSettingHandler,
		deleteOrganizationSettingHandler,
		getOrganizationSettingHandler,
		listOrganizationSettingsHandler,
	)

	// Initialize HTTP server
	router := setupRouter(configHandler, featureFlagHandler, userPreferenceHandler, organizationSettingHandler)
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Println("Starting HTTP server on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Give outstanding requests a deadline for completion
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

// initializeDatabase initializes the PostgreSQL database connection
func initializeDatabase() (*gorm.DB, error) {
	config := postgres.NewConnectionConfig()
	db, err := postgres.Connect(config)
	if err != nil {
		return nil, err
	}

	// Run migrations
	if err := postgres.AutoMigrate(db); err != nil {
		return nil, err
	}

	// Create indexes
	if err := postgres.CreateIndexes(db); err != nil {
		return nil, err
	}

	return db, nil
}

// initializeNATS initializes the NATS connection
func initializeNATS() (*nats.Conn, error) {
	config := nats.NewConnectionConfig()
	conn, err := nats.Connect(config)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// initializeRedis initializes the Redis connection
func initializeRedis() (*redis.Client, error) {
	config := redis.NewConnectionConfig()
	client, err := redis.Connect(config)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// setupRouter configures the HTTP router
func setupRouter(
	configHandler *handlers.ConfigurationHandler,
	featureFlagHandler *handlers.FeatureFlagHandler,
	userPreferenceHandler *handlers.UserPreferenceHandler,
	organizationSettingHandler *handlers.OrganizationSettingHandler,
) *gin.Engine {
	// Set Gin mode
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().UTC(),
			"service":   "settings-management-service",
		})
	})

	// API routes
	api := router.Group("/api/v1")
	{
		// Configuration routes
		configurations := api.Group("/configurations")
		{
			configurations.POST("", configHandler.CreateConfiguration)
			configurations.GET("", configHandler.ListConfigurations)
			configurations.GET("/:id", configHandler.GetConfiguration)
			configurations.PUT("/:id", configHandler.UpdateConfiguration)
			configurations.DELETE("/:id", configHandler.DeleteConfiguration)
		}

		// Feature flag routes
		featureFlags := api.Group("/feature-flags")
		{
			featureFlags.POST("", featureFlagHandler.CreateFeatureFlag)
			featureFlags.GET("", featureFlagHandler.ListFeatureFlags)
			featureFlags.GET("/:id", featureFlagHandler.GetFeatureFlag)
			featureFlags.PUT("/:id", featureFlagHandler.UpdateFeatureFlag)
			featureFlags.DELETE("/:id", featureFlagHandler.DeleteFeatureFlag)
		}

		// User preference routes
		userPreferences := api.Group("/user-preferences")
		{
			userPreferences.POST("", userPreferenceHandler.CreateUserPreference)
			userPreferences.GET("", userPreferenceHandler.ListUserPreferences)
			userPreferences.GET("/:id", userPreferenceHandler.GetUserPreference)
			userPreferences.PUT("/:id", userPreferenceHandler.UpdateUserPreference)
			userPreferences.DELETE("/:id", userPreferenceHandler.DeleteUserPreference)
		}

		// Organization setting routes
		organizationSettings := api.Group("/organization-settings")
		{
			organizationSettings.POST("", organizationSettingHandler.CreateOrganizationSetting)
			organizationSettings.GET("", organizationSettingHandler.ListOrganizationSettings)
			organizationSettings.GET("/:id", organizationSettingHandler.GetOrganizationSetting)
			organizationSettings.PUT("/:id", organizationSettingHandler.UpdateOrganizationSetting)
			organizationSettings.DELETE("/:id", organizationSettingHandler.DeleteOrganizationSetting)
		}
	}

	return router
}

// corsMiddleware adds CORS headers
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
