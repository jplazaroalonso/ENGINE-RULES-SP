package commands

import (
	"context"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/customer-management-service/internal/domain/customer"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/customer-management-service/internal/domain/shared"
)

// UpdateCustomerCommand represents a command to update a customer
type UpdateCustomerCommand struct {
	CustomerID  string                            `json:"customerId" validate:"required"`
	Name        *string                           `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Age         *int                              `json:"age,omitempty" validate:"omitempty,min=0,max=150"`
	Gender      *string                           `json:"gender,omitempty" validate:"omitempty,oneof=MALE FEMALE OTHER UNKNOWN"`
	Location    *UpdateCustomerLocationRequest    `json:"location,omitempty"`
	Preferences *UpdateCustomerPreferencesRequest `json:"preferences,omitempty"`
	Tags        *[]string                         `json:"tags,omitempty"`
	UpdatedBy   string                            `json:"updatedBy" validate:"required"`
}

// UpdateCustomerLocationRequest represents location data for customer update
type UpdateCustomerLocationRequest struct {
	Country    *string  `json:"country,omitempty"`
	City       *string  `json:"city,omitempty"`
	Region     *string  `json:"region,omitempty"`
	PostalCode *string  `json:"postalCode,omitempty"`
	Timezone   *string  `json:"timezone,omitempty"`
	Latitude   *float64 `json:"latitude,omitempty"`
	Longitude  *float64 `json:"longitude,omitempty"`
}

// UpdateCustomerPreferencesRequest represents preferences data for customer update
type UpdateCustomerPreferencesRequest struct {
	Language              *string                            `json:"language,omitempty"`
	Currency              *string                            `json:"currency,omitempty"`
	Timezone              *string                            `json:"timezone,omitempty"`
	NotificationSettings  *UpdateNotificationSettingsRequest `json:"notificationSettings,omitempty"`
	PrivacySettings       *UpdatePrivacySettingsRequest      `json:"privacySettings,omitempty"`
	MarketingConsent      *bool                              `json:"marketingConsent,omitempty"`
	DataProcessingConsent *bool                              `json:"dataProcessingConsent,omitempty"`
	CustomPreferences     *map[string]interface{}            `json:"customPreferences,omitempty"`
}

// UpdateNotificationSettingsRequest represents notification settings for customer update
type UpdateNotificationSettingsRequest struct {
	EmailNotifications *bool `json:"emailNotifications,omitempty"`
	SMSNotifications   *bool `json:"smsNotifications,omitempty"`
	PushNotifications  *bool `json:"pushNotifications,omitempty"`
	MarketingEmails    *bool `json:"marketingEmails,omitempty"`
	SystemAlerts       *bool `json:"systemAlerts,omitempty"`
}

// UpdatePrivacySettingsRequest represents privacy settings for customer update
type UpdatePrivacySettingsRequest struct {
	DataSharing       *bool `json:"dataSharing,omitempty"`
	AnalyticsTracking *bool `json:"analyticsTracking,omitempty"`
	Personalization   *bool `json:"personalization,omitempty"`
	ThirdPartySharing *bool `json:"thirdPartySharing,omitempty"`
}

// UpdateCustomerResult represents the result of updating a customer
type UpdateCustomerResult struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	Version   int    `json:"version"`
	UpdatedAt string `json:"updatedAt"`
}

// UpdateCustomerHandler handles the update customer command
type UpdateCustomerHandler struct {
	customerRepo    customer.CustomerRepository
	eventBus        shared.EventBus
	validator       shared.StructValidator
	segmentationSvc customer.CustomerSegmentationService
}

// NewUpdateCustomerHandler creates a new update customer handler
func NewUpdateCustomerHandler(
	customerRepo customer.CustomerRepository,
	eventBus shared.EventBus,
	validator shared.StructValidator,
	segmentationSvc customer.CustomerSegmentationService,
) *UpdateCustomerHandler {
	return &UpdateCustomerHandler{
		customerRepo:    customerRepo,
		eventBus:        eventBus,
		validator:       validator,
		segmentationSvc: segmentationSvc,
	}
}

// Handle handles the update customer command
func (h *UpdateCustomerHandler) Handle(ctx context.Context, cmd UpdateCustomerCommand) (*UpdateCustomerResult, error) {
	// Validate the command
	if err := h.validator.Validate(cmd); err != nil {
		return nil, shared.NewValidationError("invalid update customer command", err)
	}

	// Parse customer ID
	customerID, err := shared.NewCustomerIDFromString(cmd.CustomerID)
	if err != nil {
		return nil, err
	}

	// Find the customer
	customerEntity, err := h.customerRepo.FindByID(ctx, customerID)
	if err != nil {
		return nil, shared.NewNotFoundError("customer not found", err)
	}

	// Update name if provided
	if cmd.Name != nil {
		if err := customerEntity.UpdateName(*cmd.Name); err != nil {
			return nil, err
		}
	}

	// Update age if provided
	if cmd.Age != nil {
		if err := customerEntity.UpdateAge(cmd.Age); err != nil {
			return nil, err
		}
	}

	// Update gender if provided
	if cmd.Gender != nil {
		gender, err := customer.ParseGender(*cmd.Gender)
		if err != nil {
			return nil, err
		}
		if err := customerEntity.UpdateGender(&gender); err != nil {
			return nil, err
		}
	}

	// Update location if provided
	if cmd.Location != nil {
		location, err := h.updateLocation(customerEntity.GetLocation(), cmd.Location)
		if err != nil {
			return nil, err
		}
		if err := customerEntity.UpdateLocation(location); err != nil {
			return nil, err
		}
	}

	// Update preferences if provided
	if cmd.Preferences != nil {
		preferences, err := h.updatePreferences(customerEntity.GetPreferences(), cmd.Preferences)
		if err != nil {
			return nil, err
		}
		if err := customerEntity.UpdatePreferences(preferences); err != nil {
			return nil, err
		}
	}

	// Update tags if provided
	if cmd.Tags != nil {
		// Remove all existing tags and add new ones
		existingTags := customerEntity.GetTags()
		for _, tag := range existingTags {
			if err := customerEntity.RemoveTag(tag); err != nil {
				// Ignore error if tag doesn't exist
			}
		}
		for _, tag := range *cmd.Tags {
			if err := customerEntity.AddTag(tag); err != nil {
				return nil, err
			}
		}
	}

	// Save the updated customer
	if err := h.customerRepo.Update(ctx, customerEntity); err != nil {
		return nil, shared.NewInfrastructureError("failed to update customer", err)
	}

	// Publish domain events
	for _, event := range customerEntity.GetEvents() {
		if err := h.eventBus.Publish(ctx, event); err != nil {
			// Log error but don't fail the operation
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

	return &UpdateCustomerResult{
		ID:        customerEntity.GetID().String(),
		Name:      customerEntity.GetName(),
		Status:    customerEntity.GetStatus().String(),
		Version:   customerEntity.GetVersion(),
		UpdatedAt: customerEntity.GetUpdatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// updateLocation updates the customer location based on the update request
func (h *UpdateCustomerHandler) updateLocation(currentLocation *customer.CustomerLocation, update *UpdateCustomerLocationRequest) (*customer.CustomerLocation, error) {
	// If no current location, create a new one
	if currentLocation == nil {
		country := ""
		if update.Country != nil {
			country = *update.Country
		}
		city := ""
		if update.City != nil {
			city = *update.City
		}
		timezone := ""
		if update.Timezone != nil {
			timezone = *update.Timezone
		}

		return customer.NewCustomerLocation(
			country,
			city,
			getStringValue(update.Region),
			timezone,
			update.PostalCode,
			update.Latitude,
			update.Longitude,
		)
	}

	// Update existing location
	country := currentLocation.Country
	if update.Country != nil {
		country = *update.Country
	}

	city := currentLocation.City
	if update.City != nil {
		city = *update.City
	}

	region := currentLocation.Region
	if update.Region != nil {
		region = *update.Region
	}

	timezone := currentLocation.Timezone
	if update.Timezone != nil {
		timezone = *update.Timezone
	}

	postalCode := currentLocation.PostalCode
	if update.PostalCode != nil {
		postalCode = update.PostalCode
	}

	latitude := currentLocation.Latitude
	if update.Latitude != nil {
		latitude = update.Latitude
	}

	longitude := currentLocation.Longitude
	if update.Longitude != nil {
		longitude = update.Longitude
	}

	return customer.NewCustomerLocation(
		country,
		city,
		region,
		timezone,
		postalCode,
		latitude,
		longitude,
	)
}

// updatePreferences updates the customer preferences based on the update request
func (h *UpdateCustomerHandler) updatePreferences(currentPreferences customer.CustomerPreferences, update *UpdateCustomerPreferencesRequest) (customer.CustomerPreferences, error) {
	language := currentPreferences.Language
	if update.Language != nil {
		language = *update.Language
	}

	currency := currentPreferences.Currency
	if update.Currency != nil {
		currency = *update.Currency
	}

	timezone := currentPreferences.Timezone
	if update.Timezone != nil {
		timezone = *update.Timezone
	}

	// Update notification settings
	notificationSettings := currentPreferences.NotificationSettings
	if update.NotificationSettings != nil {
		if update.NotificationSettings.EmailNotifications != nil {
			notificationSettings.EmailNotifications = *update.NotificationSettings.EmailNotifications
		}
		if update.NotificationSettings.SMSNotifications != nil {
			notificationSettings.SMSNotifications = *update.NotificationSettings.SMSNotifications
		}
		if update.NotificationSettings.PushNotifications != nil {
			notificationSettings.PushNotifications = *update.NotificationSettings.PushNotifications
		}
		if update.NotificationSettings.MarketingEmails != nil {
			notificationSettings.MarketingEmails = *update.NotificationSettings.MarketingEmails
		}
		if update.NotificationSettings.SystemAlerts != nil {
			notificationSettings.SystemAlerts = *update.NotificationSettings.SystemAlerts
		}
	}

	// Update privacy settings
	privacySettings := currentPreferences.PrivacySettings
	if update.PrivacySettings != nil {
		if update.PrivacySettings.DataSharing != nil {
			privacySettings.DataSharing = *update.PrivacySettings.DataSharing
		}
		if update.PrivacySettings.AnalyticsTracking != nil {
			privacySettings.AnalyticsTracking = *update.PrivacySettings.AnalyticsTracking
		}
		if update.PrivacySettings.Personalization != nil {
			privacySettings.Personalization = *update.PrivacySettings.Personalization
		}
		if update.PrivacySettings.ThirdPartySharing != nil {
			privacySettings.ThirdPartySharing = *update.PrivacySettings.ThirdPartySharing
		}
	}

	marketingConsent := currentPreferences.MarketingConsent
	if update.MarketingConsent != nil {
		marketingConsent = *update.MarketingConsent
	}

	dataProcessingConsent := currentPreferences.DataProcessingConsent
	if update.DataProcessingConsent != nil {
		dataProcessingConsent = *update.DataProcessingConsent
	}

	customPreferences := currentPreferences.CustomPreferences
	if update.CustomPreferences != nil {
		customPreferences = *update.CustomPreferences
	}

	return customer.NewCustomerPreferences(
		language,
		currency,
		timezone,
		notificationSettings,
		privacySettings,
		marketingConsent,
		dataProcessingConsent,
		customPreferences,
	)
}

// getStringValue returns the string value or empty string if nil
func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
