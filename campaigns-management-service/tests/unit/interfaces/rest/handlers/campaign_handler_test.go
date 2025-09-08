package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/campaigns-management-service/internal/application/commands"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/campaigns-management-service/internal/application/queries"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/campaigns-management-service/internal/domain/campaign"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/campaigns-management-service/internal/domain/shared"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/campaigns-management-service/internal/interfaces/rest/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockCreateCampaignHandler is a mock implementation of the create campaign handler
type MockCreateCampaignHandler struct {
	mock.Mock
}

func (m *MockCreateCampaignHandler) Handle(ctx context.Context, cmd commands.CreateCampaignCommand) (*commands.CreateCampaignResult, error) {
	args := m.Called(ctx, cmd)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*commands.CreateCampaignResult), args.Error(1)
}

// MockGetCampaignHandler is a mock implementation of the get campaign handler
type MockGetCampaignHandler struct {
	mock.Mock
}

func (m *MockGetCampaignHandler) Handle(ctx context.Context, query queries.GetCampaignQuery) (*queries.CampaignDTO, error) {
	args := m.Called(ctx, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*queries.CampaignDTO), args.Error(1)
}

// MockListCampaignsHandler is a mock implementation of the list campaigns handler
type MockListCampaignsHandler struct {
	mock.Mock
}

func (m *MockListCampaignsHandler) Handle(ctx context.Context, query queries.ListCampaignsQuery) (*queries.ListCampaignsResult, error) {
	args := m.Called(ctx, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*queries.ListCampaignsResult), args.Error(1)
}

// MockUpdateCampaignHandler is a mock implementation of the update campaign handler
type MockUpdateCampaignHandler struct {
	mock.Mock
}

func (m *MockUpdateCampaignHandler) Handle(ctx context.Context, cmd commands.UpdateCampaignCommand) (*commands.UpdateCampaignResult, error) {
	args := m.Called(ctx, cmd)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*commands.UpdateCampaignResult), args.Error(1)
}

// MockActivateCampaignHandler is a mock implementation of the activate campaign handler
type MockActivateCampaignHandler struct {
	mock.Mock
}

func (m *MockActivateCampaignHandler) Handle(ctx context.Context, cmd commands.ActivateCampaignCommand) (*commands.ActivateCampaignResult, error) {
	args := m.Called(ctx, cmd)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*commands.ActivateCampaignResult), args.Error(1)
}

// MockPauseCampaignHandler is a mock implementation of the pause campaign handler
type MockPauseCampaignHandler struct {
	mock.Mock
}

func (m *MockPauseCampaignHandler) Handle(ctx context.Context, cmd commands.PauseCampaignCommand) (*commands.PauseCampaignResult, error) {
	args := m.Called(ctx, cmd)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*commands.PauseCampaignResult), args.Error(1)
}

// MockDeleteCampaignHandler is a mock implementation of the delete campaign handler
type MockDeleteCampaignHandler struct {
	mock.Mock
}

func (m *MockDeleteCampaignHandler) Handle(ctx context.Context, cmd commands.DeleteCampaignCommand) (*commands.DeleteCampaignResult, error) {
	args := m.Called(ctx, cmd)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*commands.DeleteCampaignResult), args.Error(1)
}

// MockGetCampaignMetricsHandler is a mock implementation of the get campaign metrics handler
type MockGetCampaignMetricsHandler struct {
	mock.Mock
}

func (m *MockGetCampaignMetricsHandler) Handle(ctx context.Context, query queries.GetCampaignMetricsQuery) (*queries.CampaignMetricsDTO, error) {
	args := m.Called(ctx, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*queries.CampaignMetricsDTO), args.Error(1)
}

func TestCampaignHandler_CreateCampaign(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		setupMocks     func(*MockCreateCampaignHandler)
		expectedStatus int
		expectedError  string
	}{
		{
			name: "successful campaign creation",
			requestBody: dto.CreateCampaignRequest{
				Name:           "Test Campaign",
				Description:    "Test campaign description",
				Type:           "PROMOTION",
				TargetingRules: []string{"rule-123", "rule-456"},
				StartDate:      time.Now().Add(24 * time.Hour),
				EndDate:        nil,
				Budget: &dto.MoneyDTO{
					Amount:   1000.0,
					Currency: "EUR",
				},
				CreatedBy: "john.doe",
				Settings: &dto.CampaignSettingsRequest{
					TargetAudience: []string{"test-audience"},
					Channels:       []campaign.Channel{campaign.ChannelEmail, campaign.ChannelWeb},
					Frequency:      campaign.FrequencyDaily,
					MaxImpressions: 1000,
				},
			},
			setupMocks: func(mockHandler *MockCreateCampaignHandler) {
				mockHandler.On("Handle", mock.Anything, mock.AnythingOfType("commands.CreateCampaignCommand")).Return(
					&commands.CreateCampaignResult{
						ID:      "campaign-123",
						Name:    "Test Campaign",
						Status:  "DRAFT",
						Version: 1,
					}, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "invalid request body",
			requestBody: map[string]interface{}{
				"name": "", // Invalid: empty name
			},
			setupMocks: func(mockHandler *MockCreateCampaignHandler) {
				// No mock setup needed as validation should fail before handler is called
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid request payload",
		},
		{
			name: "handler error",
			requestBody: dto.CreateCampaignRequest{
				Name:           "Test Campaign",
				Description:    "Test campaign description",
				Type:           "PROMOTION",
				TargetingRules: []string{"rule-123"},
				StartDate:      time.Now().Add(24 * time.Hour),
				CreatedBy:      "john.doe",
				Settings: &dto.CampaignSettingsRequest{
					Channels:  []campaign.Channel{campaign.ChannelEmail},
					Frequency: campaign.FrequencyDaily,
				},
			},
			setupMocks: func(mockHandler *MockCreateCampaignHandler) {
				mockHandler.On("Handle", mock.Anything, mock.AnythingOfType("commands.CreateCampaignCommand")).Return(
					nil, shared.NewBusinessError("campaign name already exists", "Test Campaign"))
			},
			expectedStatus: http.StatusConflict,
			expectedError:  "campaign name already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			gin.SetMode(gin.TestMode)
			router := gin.New()

			mockCreateHandler := new(MockCreateCampaignHandler)
			mockGetHandler := new(MockGetCampaignHandler)
			mockListHandler := new(MockListCampaignsHandler)
			mockUpdateHandler := new(MockUpdateCampaignHandler)
			mockActivateHandler := new(MockActivateCampaignHandler)
			mockPauseHandler := new(MockPauseCampaignHandler)
			mockDeleteHandler := new(MockDeleteCampaignHandler)
			mockMetricsHandler := new(MockGetCampaignMetricsHandler)

			tt.setupMocks(mockCreateHandler)

			handler := NewCampaignHandler(
				mockCreateHandler,
				mockUpdateHandler,
				mockActivateHandler,
				mockPauseHandler,
				mockDeleteHandler,
				mockGetHandler,
				mockListHandler,
				mockMetricsHandler,
			)

			router.POST("/campaigns", handler.CreateCampaign)

			// Create request body
			requestBody, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req, err := http.NewRequest("POST", "/campaigns", bytes.NewBuffer(requestBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			// Act
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				var response map[string]interface{}
				err = json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Contains(t, response["error"], tt.expectedError)
			} else {
				var response dto.CreateCampaignResponse
				err = json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.NotEmpty(t, response.ID)
				assert.Equal(t, "Test Campaign", response.Name)
				assert.Equal(t, "DRAFT", response.Status)
			}

			mockCreateHandler.AssertExpectations(t)
		})
	}
}

func TestCampaignHandler_GetCampaign(t *testing.T) {
	tests := []struct {
		name           string
		campaignID     string
		setupMocks     func(*MockGetCampaignHandler)
		expectedStatus int
		expectedError  string
	}{
		{
			name:       "successful campaign retrieval",
			campaignID: "campaign-123",
			setupMocks: func(mockHandler *MockGetCampaignHandler) {
				mockHandler.On("Handle", mock.Anything, mock.AnythingOfType("queries.GetCampaignQuery")).Return(
					&queries.CampaignDTO{
						ID:          "campaign-123",
						Name:        "Test Campaign",
						Description: "Test campaign description",
						Status:      "ACTIVE",
						Type:        "PROMOTION",
						CreatedBy:   "john.doe",
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
						Version:     1,
					}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:       "campaign not found",
			campaignID: "non-existent",
			setupMocks: func(mockHandler *MockGetCampaignHandler) {
				mockHandler.On("Handle", mock.Anything, mock.AnythingOfType("queries.GetCampaignQuery")).Return(
					nil, shared.NewNotFoundError("campaign not found", nil))
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "campaign not found",
		},
		{
			name:       "invalid campaign ID",
			campaignID: "invalid-uuid",
			setupMocks: func(mockHandler *MockGetCampaignHandler) {
				// No mock setup needed as validation should fail before handler is called
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid campaign ID format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			gin.SetMode(gin.TestMode)
			router := gin.New()

			mockCreateHandler := new(MockCreateCampaignHandler)
			mockGetHandler := new(MockGetCampaignHandler)
			mockListHandler := new(MockListCampaignsHandler)
			mockUpdateHandler := new(MockUpdateCampaignHandler)
			mockActivateHandler := new(MockActivateCampaignHandler)
			mockPauseHandler := new(MockPauseCampaignHandler)
			mockDeleteHandler := new(MockDeleteCampaignHandler)
			mockMetricsHandler := new(MockGetCampaignMetricsHandler)

			tt.setupMocks(mockGetHandler)

			handler := NewCampaignHandler(
				mockCreateHandler,
				mockUpdateHandler,
				mockActivateHandler,
				mockPauseHandler,
				mockDeleteHandler,
				mockGetHandler,
				mockListHandler,
				mockMetricsHandler,
			)

			router.GET("/campaigns/:id", handler.GetCampaign)

			// Act
			req, err := http.NewRequest("GET", "/campaigns/"+tt.campaignID, nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				var response map[string]interface{}
				err = json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Contains(t, response["error"], tt.expectedError)
			} else {
				var response dto.CampaignResponse
				err = json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Equal(t, "campaign-123", response.ID)
				assert.Equal(t, "Test Campaign", response.Name)
			}

			mockGetHandler.AssertExpectations(t)
		})
	}
}

func TestCampaignHandler_ListCampaigns(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		setupMocks     func(*MockListCampaignsHandler)
		expectedStatus int
		expectedCount  int
	}{
		{
			name:        "successful campaign listing",
			queryParams: "?page=1&limit=10",
			setupMocks: func(mockHandler *MockListCampaignsHandler) {
				mockHandler.On("Handle", mock.Anything, mock.AnythingOfType("queries.ListCampaignsQuery")).Return(
					&queries.ListCampaignsResult{
						Campaigns: []*queries.CampaignDTO{
							{
								ID:     "campaign-1",
								Name:   "Campaign 1",
								Status: "ACTIVE",
							},
							{
								ID:     "campaign-2",
								Name:   "Campaign 2",
								Status: "DRAFT",
							},
						},
						Total:      2,
						Page:       1,
						Limit:      10,
						TotalPages: 1,
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name:        "empty campaign list",
			queryParams: "?page=1&limit=10",
			setupMocks: func(mockHandler *MockListCampaignsHandler) {
				mockHandler.On("Handle", mock.Anything, mock.AnythingOfType("queries.ListCampaignsQuery")).Return(
					&queries.ListCampaignsResult{
						Campaigns:  []*queries.CampaignDTO{},
						Total:      0,
						Page:       1,
						Limit:      10,
						TotalPages: 0,
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCount:  0,
		},
		{
			name:        "invalid page parameter",
			queryParams: "?page=invalid&limit=10",
			setupMocks: func(mockHandler *MockListCampaignsHandler) {
				// No mock setup needed as validation should fail before handler is called
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			gin.SetMode(gin.TestMode)
			router := gin.New()

			mockCreateHandler := new(MockCreateCampaignHandler)
			mockGetHandler := new(MockGetCampaignHandler)
			mockListHandler := new(MockListCampaignsHandler)
			mockUpdateHandler := new(MockUpdateCampaignHandler)
			mockActivateHandler := new(MockActivateCampaignHandler)
			mockPauseHandler := new(MockPauseCampaignHandler)
			mockDeleteHandler := new(MockDeleteCampaignHandler)
			mockMetricsHandler := new(MockGetCampaignMetricsHandler)

			tt.setupMocks(mockListHandler)

			handler := NewCampaignHandler(
				mockCreateHandler,
				mockUpdateHandler,
				mockActivateHandler,
				mockPauseHandler,
				mockDeleteHandler,
				mockGetHandler,
				mockListHandler,
				mockMetricsHandler,
			)

			router.GET("/campaigns", handler.ListCampaigns)

			// Act
			req, err := http.NewRequest("GET", "/campaigns"+tt.queryParams, nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var response dto.ListCampaignsResponse
				err = json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Len(t, response.Campaigns, tt.expectedCount)
				assert.Equal(t, tt.expectedCount, response.Total)
			}

			mockListHandler.AssertExpectations(t)
		})
	}
}

func TestCampaignHandler_ActivateCampaign(t *testing.T) {
	tests := []struct {
		name           string
		campaignID     string
		requestBody    interface{}
		setupMocks     func(*MockActivateCampaignHandler)
		expectedStatus int
		expectedError  string
	}{
		{
			name:       "successful campaign activation",
			campaignID: "campaign-123",
			requestBody: dto.ActivateCampaignRequest{
				ActivatedBy: "john.doe",
			},
			setupMocks: func(mockHandler *MockActivateCampaignHandler) {
				mockHandler.On("Handle", mock.Anything, mock.AnythingOfType("commands.ActivateCampaignCommand")).Return(
					&commands.ActivateCampaignResult{
						ID:        "campaign-123",
						Status:    "ACTIVE",
						UpdatedAt: time.Now(),
					}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:       "campaign not found",
			campaignID: "non-existent",
			requestBody: dto.ActivateCampaignRequest{
				ActivatedBy: "john.doe",
			},
			setupMocks: func(mockHandler *MockActivateCampaignHandler) {
				mockHandler.On("Handle", mock.Anything, mock.AnythingOfType("commands.ActivateCampaignCommand")).Return(
					nil, shared.NewNotFoundError("campaign not found", nil))
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "campaign not found",
		},
		{
			name:       "invalid request body",
			campaignID: "campaign-123",
			requestBody: map[string]interface{}{
				"activatedBy": "", // Invalid: empty activatedBy
			},
			setupMocks: func(mockHandler *MockActivateCampaignHandler) {
				// No mock setup needed as validation should fail before handler is called
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid request payload",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			gin.SetMode(gin.TestMode)
			router := gin.New()

			mockCreateHandler := new(MockCreateCampaignHandler)
			mockGetHandler := new(MockGetCampaignHandler)
			mockListHandler := new(MockListCampaignsHandler)
			mockUpdateHandler := new(MockUpdateCampaignHandler)
			mockActivateHandler := new(MockActivateCampaignHandler)
			mockPauseHandler := new(MockPauseCampaignHandler)
			mockDeleteHandler := new(MockDeleteCampaignHandler)
			mockMetricsHandler := new(MockGetCampaignMetricsHandler)

			tt.setupMocks(mockActivateHandler)

			handler := NewCampaignHandler(
				mockCreateHandler,
				mockUpdateHandler,
				mockActivateHandler,
				mockPauseHandler,
				mockDeleteHandler,
				mockGetHandler,
				mockListHandler,
				mockMetricsHandler,
			)

			router.POST("/campaigns/:id/activate", handler.ActivateCampaign)

			// Create request body
			requestBody, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req, err := http.NewRequest("POST", "/campaigns/"+tt.campaignID+"/activate", bytes.NewBuffer(requestBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			// Act
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				var response map[string]interface{}
				err = json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Contains(t, response["error"], tt.expectedError)
			} else {
				var response dto.ActivateCampaignResponse
				err = json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Equal(t, "campaign-123", response.ID)
				assert.Equal(t, "ACTIVE", response.Status)
			}

			mockActivateHandler.AssertExpectations(t)
		})
	}
}

func TestCampaignHandler_GetCampaignMetrics(t *testing.T) {
	tests := []struct {
		name           string
		campaignID     string
		setupMocks     func(*MockGetCampaignMetricsHandler)
		expectedStatus int
		expectedError  string
	}{
		{
			name:       "successful metrics retrieval",
			campaignID: "campaign-123",
			setupMocks: func(mockHandler *MockGetCampaignMetricsHandler) {
				mockHandler.On("Handle", mock.Anything, mock.AnythingOfType("queries.GetCampaignMetricsQuery")).Return(
					&queries.CampaignMetricsDTO{
						Impressions:    1000,
						Clicks:         50,
						Conversions:    5,
						Revenue:        shared.Money{Amount: 500.0, Currency: "EUR"},
						Cost:           shared.Money{Amount: 100.0, Currency: "EUR"},
						CTR:            5.0,
						ConversionRate: 10.0,
						LastUpdated:    time.Now(),
					}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:       "campaign not found",
			campaignID: "non-existent",
			setupMocks: func(mockHandler *MockGetCampaignMetricsHandler) {
				mockHandler.On("Handle", mock.Anything, mock.AnythingOfType("queries.GetCampaignMetricsQuery")).Return(
					nil, shared.NewNotFoundError("campaign not found", nil))
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "campaign not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			gin.SetMode(gin.TestMode)
			router := gin.New()

			mockCreateHandler := new(MockCreateCampaignHandler)
			mockGetHandler := new(MockGetCampaignHandler)
			mockListHandler := new(MockListCampaignsHandler)
			mockUpdateHandler := new(MockUpdateCampaignHandler)
			mockActivateHandler := new(MockActivateCampaignHandler)
			mockPauseHandler := new(MockPauseCampaignHandler)
			mockDeleteHandler := new(MockDeleteCampaignHandler)
			mockMetricsHandler := new(MockGetCampaignMetricsHandler)

			tt.setupMocks(mockMetricsHandler)

			handler := NewCampaignHandler(
				mockCreateHandler,
				mockUpdateHandler,
				mockActivateHandler,
				mockPauseHandler,
				mockDeleteHandler,
				mockGetHandler,
				mockListHandler,
				mockMetricsHandler,
			)

			router.GET("/campaigns/:id/metrics", handler.GetCampaignMetrics)

			// Act
			req, err := http.NewRequest("GET", "/campaigns/"+tt.campaignID+"/metrics", nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				var response map[string]interface{}
				err = json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Contains(t, response["error"], tt.expectedError)
			} else {
				var response dto.CampaignMetricsResponse
				err = json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Equal(t, int64(1000), response.Impressions)
				assert.Equal(t, int64(50), response.Clicks)
				assert.Equal(t, int64(5), response.Conversions)
			}

			mockMetricsHandler.AssertExpectations(t)
		})
	}
}
