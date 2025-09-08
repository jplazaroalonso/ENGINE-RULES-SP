package commands

import (
	"context"
	"testing"
	"time"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/campaigns-management-service/internal/domain/campaign"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/campaigns-management-service/internal/domain/shared"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockRepository is a mock implementation of the campaign repository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Save(ctx context.Context, campaign *campaign.Campaign) error {
	args := m.Called(ctx, campaign)
	return args.Error(0)
}

func (m *MockRepository) FindByID(ctx context.Context, id campaign.CampaignID) (*campaign.Campaign, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*campaign.Campaign), args.Error(1)
}

func (m *MockRepository) FindByName(ctx context.Context, name string) (*campaign.Campaign, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*campaign.Campaign), args.Error(1)
}

func (m *MockRepository) List(ctx context.Context, criteria campaign.ListCriteria) ([]*campaign.Campaign, error) {
	args := m.Called(ctx, criteria)
	return args.Get(0).([]*campaign.Campaign), args.Error(1)
}

func (m *MockRepository) Count(ctx context.Context, criteria campaign.ListCriteria) (int64, error) {
	args := m.Called(ctx, criteria)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRepository) Delete(ctx context.Context, id campaign.CampaignID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	args := m.Called(ctx, name)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepository) FindByStatus(ctx context.Context, status campaign.CampaignStatus) ([]*campaign.Campaign, error) {
	args := m.Called(ctx, status)
	return args.Get(0).([]*campaign.Campaign), args.Error(1)
}

func (m *MockRepository) FindByType(ctx context.Context, campaignType campaign.CampaignType) ([]*campaign.Campaign, error) {
	args := m.Called(ctx, campaignType)
	return args.Get(0).([]*campaign.Campaign), args.Error(1)
}

func (m *MockRepository) FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*campaign.Campaign, error) {
	args := m.Called(ctx, startDate, endDate)
	return args.Get(0).([]*campaign.Campaign), args.Error(1)
}

func (m *MockRepository) FindByCreatedBy(ctx context.Context, createdBy shared.UserID) ([]*campaign.Campaign, error) {
	args := m.Called(ctx, createdBy)
	return args.Get(0).([]*campaign.Campaign), args.Error(1)
}

// MockEventBus is a mock implementation of the event bus
type MockEventBus struct {
	mock.Mock
}

func (m *MockEventBus) Publish(ctx context.Context, event shared.DomainEvent) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

// MockValidator is a mock implementation of the validator
type MockValidator struct {
	mock.Mock
}

func (m *MockValidator) Validate(s interface{}) error {
	args := m.Called(s)
	return args.Error(0)
}

func TestCreateCampaignHandler_Handle(t *testing.T) {
	tests := []struct {
		name        string
		command     CreateCampaignCommand
		setupMocks  func(*MockRepository, *MockEventBus, *MockValidator)
		expectError bool
		errorMsg    string
	}{
		{
			name: "successful campaign creation",
			command: CreateCampaignCommand{
				Name:           "Test Campaign",
				Description:    "Test campaign description",
				Type:           "PROMOTION",
				TargetingRules: []string{"rule-123", "rule-456"},
				StartDate:      time.Now().Add(24 * time.Hour),
				EndDate:        nil,
				Budget:         &shared.Money{Amount: 1000.0, Currency: "EUR"},
				CreatedBy:      "john.doe",
				Settings: CreateCampaignSettings{
					TargetAudience: []string{"test-audience"},
					Channels:       []string{"EMAIL", "WEB"},
					Frequency:      "DAILY",
					MaxImpressions: 1000,
				},
			},
			setupMocks: func(repo *MockRepository, eventBus *MockEventBus, validator *MockValidator) {
				validator.On("Validate", mock.AnythingOfType("CreateCampaignCommand")).Return(nil)
				repo.On("ExistsByName", mock.Anything, "Test Campaign").Return(false, nil)
				repo.On("Save", mock.Anything, mock.AnythingOfType("*campaign.Campaign")).Return(nil)
				eventBus.On("Publish", mock.Anything, mock.AnythingOfType("shared.DomainEvent")).Return(nil)
			},
			expectError: false,
		},
		{
			name: "validation error",
			command: CreateCampaignCommand{
				Name:           "", // Invalid: empty name
				Description:    "Test campaign description",
				Type:           "PROMOTION",
				TargetingRules: []string{"rule-123"},
				StartDate:      time.Now().Add(24 * time.Hour),
				CreatedBy:      "john.doe",
				Settings: CreateCampaignSettings{
					Channels:  []string{"EMAIL"},
					Frequency: "DAILY",
				},
			},
			setupMocks: func(repo *MockRepository, eventBus *MockEventBus, validator *MockValidator) {
				validator.On("Validate", mock.AnythingOfType("CreateCampaignCommand")).Return(shared.NewValidationError("name is required", nil))
			},
			expectError: true,
			errorMsg:    "name is required",
		},
		{
			name: "campaign name already exists",
			command: CreateCampaignCommand{
				Name:           "Existing Campaign",
				Description:    "Test campaign description",
				Type:           "PROMOTION",
				TargetingRules: []string{"rule-123"},
				StartDate:      time.Now().Add(24 * time.Hour),
				CreatedBy:      "john.doe",
				Settings: CreateCampaignSettings{
					Channels:  []string{"EMAIL"},
					Frequency: "DAILY",
				},
			},
			setupMocks: func(repo *MockRepository, eventBus *MockEventBus, validator *MockValidator) {
				validator.On("Validate", mock.AnythingOfType("CreateCampaignCommand")).Return(nil)
				repo.On("ExistsByName", mock.Anything, "Existing Campaign").Return(true, nil)
			},
			expectError: true,
			errorMsg:    "campaign name already exists",
		},
		{
			name: "repository error when checking name existence",
			command: CreateCampaignCommand{
				Name:           "Test Campaign",
				Description:    "Test campaign description",
				Type:           "PROMOTION",
				TargetingRules: []string{"rule-123"},
				StartDate:      time.Now().Add(24 * time.Hour),
				CreatedBy:      "john.doe",
				Settings: CreateCampaignSettings{
					Channels:  []string{"EMAIL"},
					Frequency: "DAILY",
				},
			},
			setupMocks: func(repo *MockRepository, eventBus *MockEventBus, validator *MockValidator) {
				validator.On("Validate", mock.AnythingOfType("CreateCampaignCommand")).Return(nil)
				repo.On("ExistsByName", mock.Anything, "Test Campaign").Return(false, shared.NewInfrastructureError("database error", nil))
			},
			expectError: true,
			errorMsg:    "database error",
		},
		{
			name: "repository error when saving",
			command: CreateCampaignCommand{
				Name:           "Test Campaign",
				Description:    "Test campaign description",
				Type:           "PROMOTION",
				TargetingRules: []string{"rule-123"},
				StartDate:      time.Now().Add(24 * time.Hour),
				CreatedBy:      "john.doe",
				Settings: CreateCampaignSettings{
					Channels:  []string{"EMAIL"},
					Frequency: "DAILY",
				},
			},
			setupMocks: func(repo *MockRepository, eventBus *MockEventBus, validator *MockValidator) {
				validator.On("Validate", mock.AnythingOfType("CreateCampaignCommand")).Return(nil)
				repo.On("ExistsByName", mock.Anything, "Test Campaign").Return(false, nil)
				repo.On("Save", mock.Anything, mock.AnythingOfType("*campaign.Campaign")).Return(shared.NewInfrastructureError("save error", nil))
			},
			expectError: true,
			errorMsg:    "save error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := new(MockRepository)
			mockEventBus := new(MockEventBus)
			mockValidator := new(MockValidator)

			tt.setupMocks(mockRepo, mockEventBus, mockValidator)

			handler := NewCreateCampaignHandler(mockRepo, mockEventBus, mockValidator)
			ctx := context.Background()

			// Act
			result, err := handler.Handle(ctx, tt.command)

			// Assert
			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, tt.command.Name, result.Name)
				assert.Equal(t, "DRAFT", result.Status)
				assert.Equal(t, 1, result.Version)
				assert.NotEmpty(t, result.ID)
			}

			// Verify all expectations were met
			mockRepo.AssertExpectations(t)
			mockEventBus.AssertExpectations(t)
			mockValidator.AssertExpectations(t)
		})
	}
}

func TestCreateCampaignHandler_Handle_InvalidCampaignType(t *testing.T) {
	// Arrange
	mockRepo := new(MockRepository)
	mockEventBus := new(MockEventBus)
	mockValidator := new(MockValidator)

	mockValidator.On("Validate", mock.AnythingOfType("CreateCampaignCommand")).Return(nil)

	handler := NewCreateCampaignHandler(mockRepo, mockEventBus, mockValidator)
	ctx := context.Background()

	command := CreateCampaignCommand{
		Name:           "Test Campaign",
		Description:    "Test campaign description",
		Type:           "INVALID_TYPE", // Invalid campaign type
		TargetingRules: []string{"rule-123"},
		StartDate:      time.Now().Add(24 * time.Hour),
		CreatedBy:      "john.doe",
		Settings: CreateCampaignSettings{
			Channels:  []string{"EMAIL"},
			Frequency: "DAILY",
		},
	}

	// Act
	result, err := handler.Handle(ctx, command)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid campaign type")
	assert.Nil(t, result)

	mockValidator.AssertExpectations(t)
}

func TestCreateCampaignHandler_Handle_InvalidTargetingRules(t *testing.T) {
	// Arrange
	mockRepo := new(MockRepository)
	mockEventBus := new(MockEventBus)
	mockValidator := new(MockValidator)

	mockValidator.On("Validate", mock.AnythingOfType("CreateCampaignCommand")).Return(nil)

	handler := NewCreateCampaignHandler(mockRepo, mockEventBus, mockValidator)
	ctx := context.Background()

	command := CreateCampaignCommand{
		Name:           "Test Campaign",
		Description:    "Test campaign description",
		Type:           "PROMOTION",
		TargetingRules: []string{}, // Invalid: empty targeting rules
		StartDate:      time.Now().Add(24 * time.Hour),
		CreatedBy:      "john.doe",
		Settings: CreateCampaignSettings{
			Channels:  []string{"EMAIL"},
			Frequency: "DAILY",
		},
	}

	// Act
	result, err := handler.Handle(ctx, command)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "at least one targeting rule must be specified")
	assert.Nil(t, result)

	mockValidator.AssertExpectations(t)
}

func TestCreateCampaignHandler_Handle_InvalidBudget(t *testing.T) {
	// Arrange
	mockRepo := new(MockRepository)
	mockEventBus := new(MockEventBus)
	mockValidator := new(MockValidator)

	mockValidator.On("Validate", mock.AnythingOfType("CreateCampaignCommand")).Return(nil)

	handler := NewCreateCampaignHandler(mockRepo, mockEventBus, mockValidator)
	ctx := context.Background()

	command := CreateCampaignCommand{
		Name:           "Test Campaign",
		Description:    "Test campaign description",
		Type:           "PROMOTION",
		TargetingRules: []string{"rule-123"},
		StartDate:      time.Now().Add(24 * time.Hour),
		Budget:         &shared.Money{Amount: -100.0, Currency: "EUR"}, // Invalid: negative amount
		CreatedBy:      "john.doe",
		Settings: CreateCampaignSettings{
			Channels:  []string{"EMAIL"},
			Frequency: "DAILY",
		},
	}

	// Act
	result, err := handler.Handle(ctx, command)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "budget amount must be positive")
	assert.Nil(t, result)

	mockValidator.AssertExpectations(t)
}

func TestCreateCampaignHandler_Handle_EventPublishingError(t *testing.T) {
	// Arrange
	mockRepo := new(MockRepository)
	mockEventBus := new(MockEventBus)
	mockValidator := new(MockValidator)

	mockValidator.On("Validate", mock.AnythingOfType("CreateCampaignCommand")).Return(nil)
	mockRepo.On("ExistsByName", mock.Anything, "Test Campaign").Return(false, nil)
	mockRepo.On("Save", mock.Anything, mock.AnythingOfType("*campaign.Campaign")).Return(nil)
	mockEventBus.On("Publish", mock.Anything, mock.AnythingOfType("shared.DomainEvent")).Return(shared.NewInfrastructureError("event publishing error", nil))

	handler := NewCreateCampaignHandler(mockRepo, mockEventBus, mockValidator)
	ctx := context.Background()

	command := CreateCampaignCommand{
		Name:           "Test Campaign",
		Description:    "Test campaign description",
		Type:           "PROMOTION",
		TargetingRules: []string{"rule-123"},
		StartDate:      time.Now().Add(24 * time.Hour),
		CreatedBy:      "john.doe",
		Settings: CreateCampaignSettings{
			Channels:  []string{"EMAIL"},
			Frequency: "DAILY",
		},
	}

	// Act
	result, err := handler.Handle(ctx, command)

	// Assert
	// Note: Event publishing errors should not fail the command in this implementation
	// In a real scenario, you might want to implement an outbox pattern
	require.NoError(t, err)
	require.NotNil(t, result)

	mockRepo.AssertExpectations(t)
	mockEventBus.AssertExpectations(t)
	mockValidator.AssertExpectations(t)
}
