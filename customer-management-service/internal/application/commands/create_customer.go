package commands

import (
	"context"
	"time"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/customer-management-service/internal/domain/customer"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/customer-management-service/internal/domain/shared"
)

// CreateCustomerCommand represents a command to create a new customer
type CreateCustomerCommand struct {
	Email       string                           `json:"email" validate:"required,email"`
	Name        string                           `json:"name" validate:"required,min=1,max=255"`
	Age         *int                             `json:"age,omitempty" validate:"omitempty,min=0,max=150"`
	Gender      *string                          `json:"gender,omitempty" validate:"omitempty,oneof=MALE FEMALE OTHER UNKNOWN"`
	Location    *CreateCustomerLocationRequest   `json:"location,omitempty"`
	Preferences CreateCustomerPreferencesRequest `json:"preferences" validate:"required"`
	Tags        []string                         `json:"tags,omitempty"`
	Metadata    CreateCustomerMetadataRequest    `json:"metadata,omitempty"`
}

// CreateCustomerLocationRequest represents location data for customer creation
type CreateCustomerLocationRequest struct {
	Country    string   `json:"country" validate:"required"`
	City       string   `json:"city" validate:"required"`
	Region     string   `json:"region,omitempty"`
	PostalCode *string  `json:"postalCode,omitempty"`
	Timezone   string   `json:"timezone" validate:"required"`
	Latitude   *float64 `json:"latitude,omitempty"`
	Longitude  *float64 `json:"longitude,omitempty"`
}

// CreateCustomerPreferencesRequest represents preferences data for customer creation
type CreateCustomerPreferencesRequest struct {
	Language              string                            `json:"language" validate:"required"`
	Currency              string                            `json:"currency" validate:"required"`
	Timezone              string                            `json:"timezone" validate:"required"`
	NotificationSettings  CreateNotificationSettingsRequest `json:"notificationSettings" validate:"required"`
	PrivacySettings       CreatePrivacySettingsRequest      `json:"privacySettings" validate:"required"`
	MarketingConsent      bool                              `json:"marketingConsent"`
	DataProcessingConsent bool                              `json:"dataProcessingConsent"`
	CustomPreferences     map[string]interface{}            `json:"customPreferences,omitempty"`
}

// CreateNotificationSettingsRequest represents notification settings for customer creation
type CreateNotificationSettingsRequest struct {
	EmailNotifications bool `json:"emailNotifications"`
	SMSNotifications   bool `json:"smsNotifications"`
	PushNotifications  bool `json:"pushNotifications"`
	MarketingEmails    bool `json:"marketingEmails"`
	SystemAlerts       bool `json:"systemAlerts"`
}

// CreatePrivacySettingsRequest represents privacy settings for customer creation
type CreatePrivacySettingsRequest struct {
	DataSharing       bool `json:"dataSharing"`
	AnalyticsTracking bool `json:"analyticsTracking"`
	Personalization   bool `json:"personalization"`
	ThirdPartySharing bool `json:"thirdPartySharing"`
}

// CreateCustomerMetadataRequest represents metadata for customer creation
type CreateCustomerMetadataRequest struct {
	Source         string  `json:"source" validate:"required"`
	ReferralSource *string `json:"referralSource,omitempty"`
}

// CreateCustomerResult represents the result of creating a customer
type CreateCustomerResult struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	Version   int       `json:"version"`
}

// CreateCustomerHandler handles the create customer command
type CreateCustomerHandler struct {
	customerRepo    customer.CustomerRepository
	eventBus        shared.EventBus
	validator       shared.StructValidator
	segmentationSvc customer.CustomerSegmentationService
}

// NewCreateCustomerHandler creates a new create customer handler
func NewCreateCustomerHandler(
	customerRepo customer.CustomerRepository,
	eventBus shared.EventBus,
	validator shared.StructValidator,
	segmentationSvc customer.CustomerSegmentationService,
) *CreateCustomerHandler {
	return &CreateCustomerHandler{
		customerRepo:    customerRepo,
		eventBus:        eventBus,
		validator:       validator,
		segmentationSvc: segmentationSvc,
	}
}

// Handle handles the create customer command
func (h *CreateCustomerHandler) Handle(ctx context.Context, cmd CreateCustomerCommand) (*CreateCustomerResult, error) {
	// Validate the command
	if err := h.validator.Validate(cmd); err != nil {
		return nil, shared.NewValidationError("invalid create customer command", err)
	}

	// Create email address value object
	email, err := shared.NewEmailAddress(cmd.Email)
	if err != nil {
		return nil, err
	}

	// Check if customer with this email already exists
	exists, err := h.customerRepo.ExistsByEmail(ctx, email)
	if err != nil {
		return nil, shared.NewInfrastructureError("failed to check email existence", err)
	}
	if exists {
		return nil, shared.NewConflictError("customer with this email already exists", "EMAIL_EXISTS")
	}

	// Parse gender if provided
	var gender *customer.Gender
	if cmd.Gender != nil {
		parsedGender, err := customer.ParseGender(*cmd.Gender)
		if err != nil {
			return nil, err
		}
		gender = &parsedGender
	}

	// Create location if provided
	var location *customer.CustomerLocation
	if cmd.Location != nil {
		location, err = customer.NewCustomerLocation(
			cmd.Location.Country,
			cmd.Location.City,
			cmd.Location.Region,
			cmd.Location.Timezone,
			cmd.Location.PostalCode,
			cmd.Location.Latitude,
			cmd.Location.Longitude,
		)
		if err != nil {
			return nil, err
		}
	}

	// Create notification settings
	notificationSettings := customer.NotificationSettings{
		EmailNotifications: cmd.Preferences.NotificationSettings.EmailNotifications,
		SMSNotifications:   cmd.Preferences.NotificationSettings.SMSNotifications,
		PushNotifications:  cmd.Preferences.NotificationSettings.PushNotifications,
		MarketingEmails:    cmd.Preferences.NotificationSettings.MarketingEmails,
		SystemAlerts:       cmd.Preferences.NotificationSettings.SystemAlerts,
	}

	// Create privacy settings
	privacySettings := customer.PrivacySettings{
		DataSharing:       cmd.Preferences.PrivacySettings.DataSharing,
		AnalyticsTracking: cmd.Preferences.PrivacySettings.AnalyticsTracking,
		Personalization:   cmd.Preferences.PrivacySettings.Personalization,
		ThirdPartySharing: cmd.Preferences.PrivacySettings.ThirdPartySharing,
	}

	// Create preferences
	preferences, err := customer.NewCustomerPreferences(
		cmd.Preferences.Language,
		cmd.Preferences.Currency,
		cmd.Preferences.Timezone,
		notificationSettings,
		privacySettings,
		cmd.Preferences.MarketingConsent,
		cmd.Preferences.DataProcessingConsent,
		cmd.Preferences.CustomPreferences,
	)
	if err != nil {
		return nil, err
	}

	// Create metadata
	now := time.Now()
	zeroMoney, _ := shared.NewMoney(0, cmd.Preferences.Currency)
	metadata := customer.NewCustomerMetadata(
		cmd.Metadata.Source,
		now, // acquisition date
		zeroMoney,
		[]customer.PurchaseRecord{},
		[]customer.InteractionRecord{},
		[]customer.DeviceInfo{},
		cmd.Metadata.ReferralSource,
		nil, // last login
		0,   // login count
	)

	// Create the customer
	customerEntity, err := customer.NewCustomer(
		email,
		cmd.Name,
		cmd.Age,
		gender,
		location,
		preferences,
		cmd.Tags,
		metadata,
	)
	if err != nil {
		return nil, err
	}

	// Save the customer
	if err := h.customerRepo.Save(ctx, customerEntity); err != nil {
		return nil, shared.NewInfrastructureError("failed to save customer", err)
	}

	// Publish domain events
	for _, event := range customerEntity.GetEvents() {
		if err := h.eventBus.Publish(ctx, event); err != nil {
			// Log error but don't fail the operation
			// In a real implementation, you might want to implement an outbox pattern
		}
	}

	// Update segment membership (this will be done asynchronously)
	go func() {
		if err := h.segmentationSvc.UpdateSegmentMembership(context.Background(), customerEntity.GetID()); err != nil {
			// Log error but don't fail the operation
		}
	}()

	// Clear events after publishing
	customerEntity.ClearEvents()

	return &CreateCustomerResult{
		ID:        customerEntity.GetID().String(),
		Email:     customerEntity.GetEmail().String(),
		Name:      customerEntity.GetName(),
		Status:    customerEntity.GetStatus().String(),
		CreatedAt: customerEntity.GetCreatedAt(),
		Version:   customerEntity.GetVersion(),
	}, nil
}
