package commands_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/application/commands"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/domain/settings"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/domain/shared"
)

// MockConfigurationRepository is a mock implementation of the ConfigurationRepository interface
type MockConfigurationRepository struct {
	mock.Mock
}

func (m *MockConfigurationRepository) Save(config *settings.Configuration) error {
	args := m.Called(config)
	return args.Error(0)
}

func (m *MockConfigurationRepository) FindByID(id uuid.UUID) (*settings.Configuration, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*settings.Configuration), args.Error(1)
}

func (m *MockConfigurationRepository) FindByKey(key string) (*settings.Configuration, error) {
	args := m.Called(key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*settings.Configuration), args.Error(1)
}

func (m *MockConfigurationRepository) List(filters settings.ListFilters, options settings.ListOptions) ([]*settings.Configuration, error) {
	args := m.Called(filters, options)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*settings.Configuration), args.Error(1)
}

func (m *MockConfigurationRepository) Count(filters settings.ListFilters) (int64, error) {
	args := m.Called(filters)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockConfigurationRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockConfigurationRepository) ExistsByKey(key string) (bool, error) {
	args := m.Called(key)
	return args.Bool(0), args.Error(1)
}

func (m *MockConfigurationRepository) FindByCategory(category string) ([]*settings.Configuration, error) {
	args := m.Called(category)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*settings.Configuration), args.Error(1)
}

func (m *MockConfigurationRepository) FindByEnvironment(environment string) ([]*settings.Configuration, error) {
	args := m.Called(environment)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*settings.Configuration), args.Error(1)
}

// MockEventBus is a mock implementation of the EventBus interface
type MockEventBus struct {
	mock.Mock
}

func (m *MockEventBus) Publish(event shared.DomainEvent) error {
	args := m.Called(event)
	return args.Error(0)
}

func (m *MockEventBus) Subscribe(eventType string, handler func(shared.DomainEvent) error) error {
	args := m.Called(eventType, handler)
	return args.Error(0)
}

func TestCreateConfigurationCommandHandler_Handle(t *testing.T) {
	tests := []struct {
		name        string
		command     commands.CreateConfigurationCommand
		setupMocks  func(*MockConfigurationRepository, *MockEventBus)
		expectError bool
		errorMsg    string
	}{
		{
			name: "successful configuration creation",
			command: commands.CreateConfigurationCommand{
				Key:         "database.host",
				Value:       map[string]interface{}{"host": "localhost", "port": 5432},
				Category:    "database",
				Environment: "development",
				Description: "Database connection settings",
				IsEncrypted: false,
			},
			setupMocks: func(repo *MockConfigurationRepository, eventBus *MockEventBus) {
				repo.On("ExistsByKey", "database.host").Return(false, nil)
				repo.On("Save", mock.AnythingOfType("*settings.Configuration")).Return(nil)
				eventBus.On("Publish", mock.AnythingOfType("*settings.ConfigurationCreatedEvent")).Return(nil)
			},
			expectError: false,
		},
		{
			name: "configuration with key already exists",
			command: commands.CreateConfigurationCommand{
				Key:         "database.host",
				Value:       map[string]interface{}{"host": "localhost", "port": 5432},
				Category:    "database",
				Environment: "development",
				Description: "Database connection settings",
				IsEncrypted: false,
			},
			setupMocks: func(repo *MockConfigurationRepository, eventBus *MockEventBus) {
				repo.On("ExistsByKey", "database.host").Return(true, nil)
			},
			expectError: true,
			errorMsg:    "configuration with key 'database.host' already exists",
		},
		{
			name: "repository error when checking key existence",
			command: commands.CreateConfigurationCommand{
				Key:         "database.host",
				Value:       map[string]interface{}{"host": "localhost", "port": 5432},
				Category:    "database",
				Environment: "development",
				Description: "Database connection settings",
				IsEncrypted: false,
			},
			setupMocks: func(repo *MockConfigurationRepository, eventBus *MockEventBus) {
				repo.On("ExistsByKey", "database.host").Return(false, assert.AnError)
			},
			expectError: true,
			errorMsg:    "failed to check if configuration exists",
		},
		{
			name: "repository error when saving configuration",
			command: commands.CreateConfigurationCommand{
				Key:         "database.host",
				Value:       map[string]interface{}{"host": "localhost", "port": 5432},
				Category:    "database",
				Environment: "development",
				Description: "Database connection settings",
				IsEncrypted: false,
			},
			setupMocks: func(repo *MockConfigurationRepository, eventBus *MockEventBus) {
				repo.On("ExistsByKey", "database.host").Return(false, nil)
				repo.On("Save", mock.AnythingOfType("*settings.Configuration")).Return(assert.AnError)
			},
			expectError: true,
			errorMsg:    "failed to save configuration",
		},
		{
			name: "event bus error when publishing event",
			command: commands.CreateConfigurationCommand{
				Key:         "database.host",
				Value:       map[string]interface{}{"host": "localhost", "port": 5432},
				Category:    "database",
				Environment: "development",
				Description: "Database connection settings",
				IsEncrypted: false,
			},
			setupMocks: func(repo *MockConfigurationRepository, eventBus *MockEventBus) {
				repo.On("ExistsByKey", "database.host").Return(false, nil)
				repo.On("Save", mock.AnythingOfType("*settings.Configuration")).Return(nil)
				eventBus.On("Publish", mock.AnythingOfType("*settings.ConfigurationCreatedEvent")).Return(assert.AnError)
			},
			expectError: true,
			errorMsg:    "failed to publish configuration created event",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mocks
			mockRepo := new(MockConfigurationRepository)
			mockEventBus := new(MockEventBus)

			// Setup mocks
			tt.setupMocks(mockRepo, mockEventBus)

			// Create handler
			handler := commands.NewCreateConfigurationCommandHandler(mockRepo, mockEventBus)

			// Execute command
			result, err := handler.Handle(tt.command)

			// Assertions
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.command.Key, result.Key)
				assert.Equal(t, tt.command.Value, result.Value)
				assert.Equal(t, tt.command.Category, result.Category)
				assert.Equal(t, tt.command.Environment, result.Environment)
				assert.Equal(t, tt.command.Description, result.Description)
				assert.Equal(t, tt.command.IsEncrypted, result.IsEncrypted)
			}

			// Verify all expectations were met
			mockRepo.AssertExpectations(t)
			mockEventBus.AssertExpectations(t)
		})
	}
}

func TestCreateConfigurationCommandHandler_Handle_InvalidCommand(t *testing.T) {
	tests := []struct {
		name     string
		command  commands.CreateConfigurationCommand
		errorMsg string
	}{
		{
			name: "empty key",
			command: commands.CreateConfigurationCommand{
				Key:         "",
				Value:       map[string]interface{}{"host": "localhost", "port": 5432},
				Category:    "database",
				Environment: "development",
				Description: "Database connection settings",
				IsEncrypted: false,
			},
			errorMsg: "key cannot be empty",
		},
		{
			name: "empty category",
			command: commands.CreateConfigurationCommand{
				Key:         "database.host",
				Value:       map[string]interface{}{"host": "localhost", "port": 5432},
				Category:    "",
				Environment: "development",
				Description: "Database connection settings",
				IsEncrypted: false,
			},
			errorMsg: "category cannot be empty",
		},
		{
			name: "empty environment",
			command: commands.CreateConfigurationCommand{
				Key:         "database.host",
				Value:       map[string]interface{}{"host": "localhost", "port": 5432},
				Category:    "database",
				Environment: "",
				Description: "Database connection settings",
				IsEncrypted: false,
			},
			errorMsg: "environment cannot be empty",
		},
		{
			name: "nil value",
			command: commands.CreateConfigurationCommand{
				Key:         "database.host",
				Value:       nil,
				Category:    "database",
				Environment: "development",
				Description: "Database connection settings",
				IsEncrypted: false,
			},
			errorMsg: "value cannot be nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mocks
			mockRepo := new(MockConfigurationRepository)
			mockEventBus := new(MockEventBus)

			// Create handler
			handler := commands.NewCreateConfigurationCommandHandler(mockRepo, mockEventBus)

			// Execute command
			result, err := handler.Handle(tt.command)

			// Assertions
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.errorMsg)
			assert.Nil(t, result)

			// Verify no repository or event bus calls were made
			mockRepo.AssertExpectations(t)
			mockEventBus.AssertExpectations(t)
		})
	}
}
