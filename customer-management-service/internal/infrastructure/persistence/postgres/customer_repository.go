package postgres

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/customer-management-service/internal/domain/customer"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/customer-management-service/internal/domain/shared"
	"gorm.io/gorm"
)

// CustomerRepository implements the customer repository interface using PostgreSQL
type CustomerRepository struct {
	db *gorm.DB
}

// NewCustomerRepository creates a new customer repository
func NewCustomerRepository(db *gorm.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

// CustomerDBModel represents the customer database model
type CustomerDBModel struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Email        string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	Name         string    `gorm:"type:varchar(255);not null"`
	Age          *int      `gorm:"type:integer;check:age >= 0 AND age <= 150"`
	Gender       *string   `gorm:"type:varchar(20);check:gender IN ('MALE', 'FEMALE', 'OTHER', 'UNKNOWN')"`
	Location     JSONB     `gorm:"type:jsonb"`
	Preferences  JSONB     `gorm:"type:jsonb;not null;default:'{}'"`
	Segments     JSONB     `gorm:"type:jsonb;not null;default:'[]'"`
	Tags         JSONB     `gorm:"type:jsonb;not null;default:'[]'"`
	Status       string    `gorm:"type:varchar(20);not null;default:'ACTIVE';check:status IN ('ACTIVE', 'INACTIVE', 'SUSPENDED', 'DELETED')"`
	CreatedAt    time.Time `gorm:"type:timestamp with time zone;default:now()"`
	UpdatedAt    time.Time `gorm:"type:timestamp with time zone;default:now()"`
	LastActivity time.Time `gorm:"type:timestamp with time zone;default:now()"`
	Metadata     JSONB     `gorm:"type:jsonb;not null;default:'{}'"`
	Version      int       `gorm:"type:integer;not null;default:1"`
}

// JSONB represents a JSONB type for PostgreSQL
type JSONB map[string]interface{}

// Scan implements the sql.Scanner interface
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = make(map[string]interface{})
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot scan %T into JSONB", value)
	}

	return json.Unmarshal(bytes, j)
}

// Value implements the driver.Valuer interface
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}

	return json.Marshal(j)
}

// Save saves a customer to the database
func (r *CustomerRepository) Save(ctx context.Context, customer *customer.Customer) error {
	customerDB := r.toDBModel(customer)

	if err := r.db.WithContext(ctx).Create(customerDB).Error; err != nil {
		return shared.NewInfrastructureError("failed to save customer", err)
	}

	return nil
}

// FindByID finds a customer by ID
func (r *CustomerRepository) FindByID(ctx context.Context, id customer.CustomerID) (*customer.Customer, error) {
	var customerDB CustomerDBModel

	if err := r.db.WithContext(ctx).Where("id = ?", id.String()).First(&customerDB).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, shared.NewNotFoundError("customer not found", err)
		}
		return nil, shared.NewInfrastructureError("failed to find customer", err)
	}

	return r.toDomainEntity(&customerDB)
}

// FindByEmail finds a customer by email
func (r *CustomerRepository) FindByEmail(ctx context.Context, email shared.EmailAddress) (*customer.Customer, error) {
	var customerDB CustomerDBModel

	if err := r.db.WithContext(ctx).Where("email = ?", email.String()).First(&customerDB).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, shared.NewNotFoundError("customer not found", err)
		}
		return nil, shared.NewInfrastructureError("failed to find customer", err)
	}

	return r.toDomainEntity(&customerDB)
}

// Update updates a customer in the database
func (r *CustomerRepository) Update(ctx context.Context, customer *customer.Customer) error {
	customerDB := r.toDBModel(customer)

	if err := r.db.WithContext(ctx).Save(customerDB).Error; err != nil {
		return shared.NewInfrastructureError("failed to update customer", err)
	}

	return nil
}

// Delete deletes a customer from the database
func (r *CustomerRepository) Delete(ctx context.Context, id customer.CustomerID) error {
	if err := r.db.WithContext(ctx).Where("id = ?", id.String()).Delete(&CustomerDBModel{}).Error; err != nil {
		return shared.NewInfrastructureError("failed to delete customer", err)
	}

	return nil
}

// List lists customers based on criteria
func (r *CustomerRepository) List(ctx context.Context, criteria customer.ListCriteria) ([]*customer.Customer, error) {
	var customerDBs []CustomerDBModel

	query := r.db.WithContext(ctx).Model(&CustomerDBModel{})

	// Apply filters
	if criteria.Status != nil {
		query = query.Where("status = ?", criteria.Status.String())
	}

	if len(criteria.SegmentIDs) > 0 {
		segmentIDs := make([]string, len(criteria.SegmentIDs))
		for i, segmentID := range criteria.SegmentIDs {
			segmentIDs[i] = segmentID.String()
		}
		query = query.Where("segments ?| ?", segmentIDs)
	}

	if len(criteria.Tags) > 0 {
		query = query.Where("tags ?| ?", criteria.Tags)
	}

	if criteria.DateRange != nil {
		query = query.Where("created_at BETWEEN ? AND ?", criteria.DateRange.StartDate, criteria.DateRange.EndDate)
	}

	// Apply sorting
	if criteria.SortBy != "" {
		order := criteria.SortBy
		if criteria.SortOrder == "desc" {
			order += " DESC"
		} else {
			order += " ASC"
		}
		query = query.Order(order)
	}

	// Apply pagination
	if criteria.Limit > 0 {
		query = query.Limit(criteria.Limit)
	}
	if criteria.Page > 0 {
		offset := (criteria.Page - 1) * criteria.Limit
		query = query.Offset(offset)
	}

	if err := query.Find(&customerDBs).Error; err != nil {
		return nil, shared.NewInfrastructureError("failed to list customers", err)
	}

	// Convert to domain entities
	customers := make([]*customer.Customer, len(customerDBs))
	for i, customerDB := range customerDBs {
		customerEntity, err := r.toDomainEntity(&customerDB)
		if err != nil {
			return nil, err
		}
		customers[i] = customerEntity
	}

	return customers, nil
}

// Count counts customers based on criteria
func (r *CustomerRepository) Count(ctx context.Context, criteria customer.ListCriteria) (int64, error) {
	var count int64

	query := r.db.WithContext(ctx).Model(&CustomerDBModel{})

	// Apply filters
	if criteria.Status != nil {
		query = query.Where("status = ?", criteria.Status.String())
	}

	if len(criteria.SegmentIDs) > 0 {
		segmentIDs := make([]string, len(criteria.SegmentIDs))
		for i, segmentID := range criteria.SegmentIDs {
			segmentIDs[i] = segmentID.String()
		}
		query = query.Where("segments ?| ?", segmentIDs)
	}

	if len(criteria.Tags) > 0 {
		query = query.Where("tags ?| ?", criteria.Tags)
	}

	if criteria.DateRange != nil {
		query = query.Where("created_at BETWEEN ? AND ?", criteria.DateRange.StartDate, criteria.DateRange.EndDate)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, shared.NewInfrastructureError("failed to count customers", err)
	}

	return count, nil
}

// Search searches customers based on query
func (r *CustomerRepository) Search(ctx context.Context, query string, criteria customer.ListCriteria) ([]*customer.Customer, error) {
	var customerDBs []CustomerDBModel

	dbQuery := r.db.WithContext(ctx).Model(&CustomerDBModel{})

	// Apply search query
	if query != "" {
		dbQuery = dbQuery.Where("name ILIKE ? OR email ILIKE ?", "%"+query+"%", "%"+query+"%")
	}

	// Apply other filters (same as List method)
	if criteria.Status != nil {
		dbQuery = dbQuery.Where("status = ?", criteria.Status.String())
	}

	if len(criteria.SegmentIDs) > 0 {
		segmentIDs := make([]string, len(criteria.SegmentIDs))
		for i, segmentID := range criteria.SegmentIDs {
			segmentIDs[i] = segmentID.String()
		}
		dbQuery = dbQuery.Where("segments ?| ?", segmentIDs)
	}

	if len(criteria.Tags) > 0 {
		dbQuery = dbQuery.Where("tags ?| ?", criteria.Tags)
	}

	if criteria.DateRange != nil {
		dbQuery = dbQuery.Where("created_at BETWEEN ? AND ?", criteria.DateRange.StartDate, criteria.DateRange.EndDate)
	}

	// Apply sorting
	if criteria.SortBy != "" {
		order := criteria.SortBy
		if criteria.SortOrder == "desc" {
			order += " DESC"
		} else {
			order += " ASC"
		}
		dbQuery = dbQuery.Order(order)
	}

	// Apply pagination
	if criteria.Limit > 0 {
		dbQuery = dbQuery.Limit(criteria.Limit)
	}
	if criteria.Page > 0 {
		offset := (criteria.Page - 1) * criteria.Limit
		dbQuery = dbQuery.Offset(offset)
	}

	if err := dbQuery.Find(&customerDBs).Error; err != nil {
		return nil, shared.NewInfrastructureError("failed to search customers", err)
	}

	// Convert to domain entities
	customers := make([]*customer.Customer, len(customerDBs))
	for i, customerDB := range customerDBs {
		customerEntity, err := r.toDomainEntity(&customerDB)
		if err != nil {
			return nil, err
		}
		customers[i] = customerEntity
	}

	return customers, nil
}

// FindByStatus finds customers by status
func (r *CustomerRepository) FindByStatus(ctx context.Context, status customer.CustomerStatus) ([]*customer.Customer, error) {
	var customerDBs []CustomerDBModel

	if err := r.db.WithContext(ctx).Where("status = ?", status.String()).Find(&customerDBs).Error; err != nil {
		return nil, shared.NewInfrastructureError("failed to find customers by status", err)
	}

	// Convert to domain entities
	customers := make([]*customer.Customer, len(customerDBs))
	for i, customerDB := range customerDBs {
		customerEntity, err := r.toDomainEntity(&customerDB)
		if err != nil {
			return nil, err
		}
		customers[i] = customerEntity
	}

	return customers, nil
}

// FindBySegment finds customers by segment
func (r *CustomerRepository) FindBySegment(ctx context.Context, segmentID customer.SegmentID) ([]*customer.Customer, error) {
	var customerDBs []CustomerDBModel

	if err := r.db.WithContext(ctx).Where("segments ? ?", segmentID.String()).Find(&customerDBs).Error; err != nil {
		return nil, shared.NewInfrastructureError("failed to find customers by segment", err)
	}

	// Convert to domain entities
	customers := make([]*customer.Customer, len(customerDBs))
	for i, customerDB := range customerDBs {
		customerEntity, err := r.toDomainEntity(&customerDB)
		if err != nil {
			return nil, err
		}
		customers[i] = customerEntity
	}

	return customers, nil
}

// FindByTags finds customers by tags
func (r *CustomerRepository) FindByTags(ctx context.Context, tags []string) ([]*customer.Customer, error) {
	var customerDBs []CustomerDBModel

	if err := r.db.WithContext(ctx).Where("tags ?| ?", tags).Find(&customerDBs).Error; err != nil {
		return nil, shared.NewInfrastructureError("failed to find customers by tags", err)
	}

	// Convert to domain entities
	customers := make([]*customer.Customer, len(customerDBs))
	for i, customerDB := range customerDBs {
		customerEntity, err := r.toDomainEntity(&customerDB)
		if err != nil {
			return nil, err
		}
		customers[i] = customerEntity
	}

	return customers, nil
}

// FindByDateRange finds customers by date range
func (r *CustomerRepository) FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*customer.Customer, error) {
	var customerDBs []CustomerDBModel

	if err := r.db.WithContext(ctx).Where("created_at BETWEEN ? AND ?", startDate, endDate).Find(&customerDBs).Error; err != nil {
		return nil, shared.NewInfrastructureError("failed to find customers by date range", err)
	}

	// Convert to domain entities
	customers := make([]*customer.Customer, len(customerDBs))
	for i, customerDB := range customerDBs {
		customerEntity, err := r.toDomainEntity(&customerDB)
		if err != nil {
			return nil, err
		}
		customers[i] = customerEntity
	}

	return customers, nil
}

// ExistsByEmail checks if a customer exists by email
func (r *CustomerRepository) ExistsByEmail(ctx context.Context, email shared.EmailAddress) (bool, error) {
	var count int64

	if err := r.db.WithContext(ctx).Model(&CustomerDBModel{}).Where("email = ?", email.String()).Count(&count).Error; err != nil {
		return false, shared.NewInfrastructureError("failed to check email existence", err)
	}

	return count > 0, nil
}

// ExistsByID checks if a customer exists by ID
func (r *CustomerRepository) ExistsByID(ctx context.Context, id customer.CustomerID) (bool, error) {
	var count int64

	if err := r.db.WithContext(ctx).Model(&CustomerDBModel{}).Where("id = ?", id.String()).Count(&count).Error; err != nil {
		return false, shared.NewInfrastructureError("failed to check customer existence", err)
	}

	return count > 0, nil
}

// BulkUpdate performs bulk updates on customers
func (r *CustomerRepository) BulkUpdate(ctx context.Context, updates []customer.BulkUpdate) error {
	// This is a simplified implementation
	// In a real scenario, you might want to use batch operations or transactions
	for _, update := range updates {
		if err := r.db.WithContext(ctx).Model(&CustomerDBModel{}).Where("id = ?", update.CustomerID.String()).Updates(update.Updates).Error; err != nil {
			return shared.NewInfrastructureError("failed to bulk update customer", err)
		}
	}

	return nil
}

// BulkDelete performs bulk deletion of customers
func (r *CustomerRepository) BulkDelete(ctx context.Context, ids []customer.CustomerID) error {
	idStrings := make([]string, len(ids))
	for i, id := range ids {
		idStrings[i] = id.String()
	}

	if err := r.db.WithContext(ctx).Where("id IN ?", idStrings).Delete(&CustomerDBModel{}).Error; err != nil {
		return shared.NewInfrastructureError("failed to bulk delete customers", err)
	}

	return nil
}

// toDBModel converts a domain customer to database model
func (r *CustomerRepository) toDBModel(customer *customer.Customer) *CustomerDBModel {
	customerID, _ := uuid.Parse(customer.GetID().String())

	// Convert location to JSONB
	var location JSONB
	if customer.GetLocation() != nil {
		location = JSONB{
			"country":    customer.GetLocation().Country,
			"city":       customer.GetLocation().City,
			"region":     customer.GetLocation().Region,
			"postalCode": customer.GetLocation().PostalCode,
			"timezone":   customer.GetLocation().Timezone,
			"latitude":   customer.GetLocation().Latitude,
			"longitude":  customer.GetLocation().Longitude,
		}
	}

	// Convert preferences to JSONB
	preferences := JSONB{
		"language":              customer.GetPreferences().Language,
		"currency":              customer.GetPreferences().Currency,
		"timezone":              customer.GetPreferences().Timezone,
		"notificationSettings":  customer.GetPreferences().NotificationSettings,
		"privacySettings":       customer.GetPreferences().PrivacySettings,
		"marketingConsent":      customer.GetPreferences().MarketingConsent,
		"dataProcessingConsent": customer.GetPreferences().DataProcessingConsent,
		"customPreferences":     customer.GetPreferences().CustomPreferences,
	}

	// Convert segments to JSONB
	segments := make([]string, len(customer.GetSegments()))
	for i, segment := range customer.GetSegments() {
		segments[i] = segment.String()
	}

	// Convert metadata to JSONB
	metadata := JSONB{
		"source":             customer.GetMetadata().Source,
		"acquisitionDate":    customer.GetMetadata().AcquisitionDate,
		"lifetimeValue":      customer.GetMetadata().LifetimeValue,
		"purchaseHistory":    customer.GetMetadata().PurchaseHistory,
		"interactionHistory": customer.GetMetadata().InteractionHistory,
		"deviceInfo":         customer.GetMetadata().DeviceInfo,
		"referralSource":     customer.GetMetadata().ReferralSource,
		"lastLogin":          customer.GetMetadata().LastLogin,
		"loginCount":         customer.GetMetadata().LoginCount,
	}

	return &CustomerDBModel{
		ID:           customerID,
		Email:        customer.GetEmail().String(),
		Name:         customer.GetName(),
		Age:          customer.GetAge(),
		Gender:       (*string)(customer.GetGender()),
		Location:     location,
		Preferences:  preferences,
		Segments:     JSONB{"segments": segments},
		Tags:         JSONB{"tags": customer.GetTags()},
		Status:       customer.GetStatus().String(),
		CreatedAt:    customer.GetCreatedAt(),
		UpdatedAt:    customer.GetUpdatedAt(),
		LastActivity: customer.GetLastActivity(),
		Metadata:     metadata,
		Version:      customer.GetVersion(),
	}
}

// toDomainEntity converts a database model to domain customer
func (r *CustomerRepository) toDomainEntity(customerDB *CustomerDBModel) (*customer.Customer, error) {
	// Parse customer ID
	customerID, _ := shared.NewCustomerIDFromString(customerDB.ID.String())

	// Parse email
	email, _ := shared.NewEmailAddress(customerDB.Email)

	// Parse gender
	var gender *customer.Gender
	if customerDB.Gender != nil {
		parsedGender, _ := customer.ParseGender(*customerDB.Gender)
		gender = &parsedGender
	}

	// Parse location
	var location *customer.CustomerLocation
	if customerDB.Location != nil && len(customerDB.Location) > 0 {
		locationData := customerDB.Location
		country, _ := locationData["country"].(string)
		city, _ := locationData["city"].(string)
		region, _ := locationData["region"].(string)
		timezone, _ := locationData["timezone"].(string)

		var postalCode *string
		if pc, ok := locationData["postalCode"].(string); ok {
			postalCode = &pc
		}

		var latitude *float64
		if lat, ok := locationData["latitude"].(float64); ok {
			latitude = &lat
		}

		var longitude *float64
		if lng, ok := locationData["longitude"].(float64); ok {
			longitude = &lng
		}

		location, _ = customer.NewCustomerLocation(country, city, region, timezone, postalCode, latitude, longitude)
	}

	// Parse preferences
	preferencesData := customerDB.Preferences
	language, _ := preferencesData["language"].(string)
	currency, _ := preferencesData["currency"].(string)
	timezone, _ := preferencesData["timezone"].(string)

	// Parse notification settings
	notificationSettingsData, _ := preferencesData["notificationSettings"].(map[string]interface{})
	notificationSettings := customer.NotificationSettings{
		EmailNotifications: getBoolValue(notificationSettingsData, "emailNotifications"),
		SMSNotifications:   getBoolValue(notificationSettingsData, "smsNotifications"),
		PushNotifications:  getBoolValue(notificationSettingsData, "pushNotifications"),
		MarketingEmails:    getBoolValue(notificationSettingsData, "marketingEmails"),
		SystemAlerts:       getBoolValue(notificationSettingsData, "systemAlerts"),
	}

	// Parse privacy settings
	privacySettingsData, _ := preferencesData["privacySettings"].(map[string]interface{})
	privacySettings := customer.PrivacySettings{
		DataSharing:       getBoolValue(privacySettingsData, "dataSharing"),
		AnalyticsTracking: getBoolValue(privacySettingsData, "analyticsTracking"),
		Personalization:   getBoolValue(privacySettingsData, "personalization"),
		ThirdPartySharing: getBoolValue(privacySettingsData, "thirdPartySharing"),
	}

	marketingConsent := getBoolValue(preferencesData, "marketingConsent")
	dataProcessingConsent := getBoolValue(preferencesData, "dataProcessingConsent")

	customPreferences, _ := preferencesData["customPreferences"].(map[string]interface{})

	preferences, _ := customer.NewCustomerPreferences(
		language,
		currency,
		timezone,
		notificationSettings,
		privacySettings,
		marketingConsent,
		dataProcessingConsent,
		customPreferences,
	)

	// Parse segments
	var segments []customer.SegmentID
	if segmentsData, ok := customerDB.Segments["segments"].([]interface{}); ok {
		for _, segmentStr := range segmentsData {
			if segmentIDStr, ok := segmentStr.(string); ok {
				if segmentID, err := shared.NewSegmentIDFromString(segmentIDStr); err == nil {
					segments = append(segments, segmentID)
				}
			}
		}
	}

	// Parse tags
	var tags []string
	if tagsData, ok := customerDB.Tags["tags"].([]interface{}); ok {
		for _, tag := range tagsData {
			if tagStr, ok := tag.(string); ok {
				tags = append(tags, tagStr)
			}
		}
	}

	// Parse metadata
	metadataData := customerDB.Metadata
	source, _ := metadataData["source"].(string)

	var acquisitionDate time.Time
	if ad, ok := metadataData["acquisitionDate"].(string); ok {
		acquisitionDate, _ = time.Parse(time.RFC3339, ad)
	}

	// Parse lifetime value
	var lifetimeValue shared.Money
	if lvData, ok := metadataData["lifetimeValue"].(map[string]interface{}); ok {
		amount, _ := lvData["amount"].(float64)
		currency, _ := lvData["currency"].(string)
		lifetimeValue, _ = shared.NewMoney(amount, currency)
	}

	var referralSource *string
	if rs, ok := metadataData["referralSource"].(string); ok {
		referralSource = &rs
	}

	var lastLogin *time.Time
	if ll, ok := metadataData["lastLogin"].(string); ok {
		if parsed, err := time.Parse(time.RFC3339, ll); err == nil {
			lastLogin = &parsed
		}
	}

	loginCount := 0
	if lc, ok := metadataData["loginCount"].(float64); ok {
		loginCount = int(lc)
	}

	metadata := customer.NewCustomerMetadata(
		source,
		acquisitionDate,
		lifetimeValue,
		[]customer.PurchaseRecord{},    // Simplified for now
		[]customer.InteractionRecord{}, // Simplified for now
		[]customer.DeviceInfo{},        // Simplified for now
		referralSource,
		lastLogin,
		loginCount,
	)

	// Create customer (this is a simplified version - in reality you'd need a factory method)
	// For now, we'll create a basic customer and then update its fields
	customerEntity, _ := customer.NewCustomer(
		email,
		customerDB.Name,
		customerDB.Age,
		gender,
		location,
		preferences,
		tags,
		metadata,
	)

	// Note: In a real implementation, you would need to properly restore the customer state
	// including the ID, status, version, and timestamps. This is a simplified version.

	return customerEntity, nil
}

// getBoolValue gets a boolean value from a map, returning false if not found
func getBoolValue(data map[string]interface{}, key string) bool {
	if value, ok := data[key].(bool); ok {
		return value
	}
	return false
}
