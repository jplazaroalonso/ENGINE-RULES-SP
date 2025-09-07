package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"rules-management-service/internal/application/commands"
	"rules-management-service/internal/application/queries"
	"rules-management-service/internal/infrastructure/config"
	"rules-management-service/internal/infrastructure/dsl"
	"rules-management-service/internal/infrastructure/messaging/nats"
	db "rules-management-service/internal/infrastructure/persistence/postgres"
	"rules-management-service/internal/interfaces/rest/handlers"
)

// AppValidator wraps the go-playground/validator.
type AppValidator struct {
	validate *validator.Validate
}

func (v *AppValidator) Validate(s interface{}) error {
	return v.validate.Struct(s)
}

func main() {
	// Load configuration
	cfg := config.DefaultConfig()

	// In a real app, you would connect to PostgreSQL.
	// DSN might come from environment variables.
	dsn := "host=localhost user=user password=password dbname=rules_db port=5432 sslmode=disable"
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("failed to connect to postgresql, falling back to sqlite: %v", err)
		// Fallback to SQLite for local development if PostgreSQL is not available.
	}
	// In a production environment, you would run migrations using a tool like goose or golang-migrate.
	// For example: `goose -dir ./migrations up`
	// For simplicity in this example, we continue to use AutoMigrate.
	err = gormDB.AutoMigrate(&db.RuleDBModel{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	// Initialize validator
	validate := validator.New()
	appValidator := &AppValidator{validate: validate}

	// Initialize services
	validationService := dsl.NewSimpleValidator()

	// Initialize repositories
	ruleRepo := db.NewRuleRepository(gormDB)

	// Initialize NATS Publisher (Event Bus)
	eventPublisher, err := nats.NewEventPublisher(cfg.NATS)
	if err != nil {
		log.Fatalf("failed to create NATS publisher: %v", err)
	}
	defer eventPublisher.Close()

	// Initialize application handlers
	createRuleHandler := commands.NewCreateRuleHandler(ruleRepo, appValidator, eventPublisher, cfg.App.ReplicationEnabled, validationService)
	getRuleHandler := queries.NewGetRuleHandler(ruleRepo)
	validateRuleHandler := commands.NewValidateRuleHandler(appValidator, validationService)

	// Initialize NATS Subscriber
	commandSubscriber, err := nats.NewCommandSubscriber(cfg.NATS, createRuleHandler)
	if err != nil {
		log.Fatalf("failed to create NATS subscriber: %v", err)
	}
	commandSubscriber.Start()
	defer commandSubscriber.Close()

	// Initialize HTTP handlers
	ruleHandler := handlers.NewRuleHandler(createRuleHandler, getRuleHandler, validateRuleHandler)

	// Setup Gin router
	router := gin.Default()
	v1 := router.Group("/v1")
	{
		rules := v1.Group("/rules")
		{
			rules.POST("", ruleHandler.CreateRule)
			rules.POST("/validate", ruleHandler.ValidateRule)
			rules.GET("/:id", ruleHandler.GetRule)
		}
	}

	// Start server in a goroutine
	go func() {
		if err := router.Run(":8080"); err != nil {
			log.Fatalf("failed to run server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
}
