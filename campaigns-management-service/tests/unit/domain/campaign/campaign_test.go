package campaign

import (
	"testing"
	"time"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/campaigns-management-service/internal/domain/shared"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCampaign(t *testing.T) {
	tests := []struct {
		name           string
		campaignName   string
		description    string
		campaignType   CampaignType
		targetingRules []shared.RuleID
		startDate      time.Time
		endDate        *time.Time
		budget         *shared.Money
		createdBy      shared.UserID
		settings       CampaignSettings
		expectError    bool
		errorMsg       string
	}{
		{
			name:           "valid campaign creation",
			campaignName:   "Summer Sale 2024",
			description:    "Summer promotion campaign",
			campaignType:   CampaignTypePromotion,
			targetingRules: []shared.RuleID{shared.NewRuleID()},
			startDate:      time.Now().Add(24 * time.Hour),
			endDate:        nil,
			budget:         &shared.Money{Amount: 1000.0, Currency: "EUR"},
			createdBy:      shared.NewUserID(),
			settings:       createValidSettings(),
			expectError:    false,
		},
		{
			name:           "empty campaign name",
			campaignName:   "",
			description:    "Test campaign",
			campaignType:   CampaignTypePromotion,
			targetingRules: []shared.RuleID{shared.NewRuleID()},
			startDate:      time.Now().Add(24 * time.Hour),
			endDate:        nil,
			budget:         nil,
			createdBy:      shared.NewUserID(),
			settings:       createValidSettings(),
			expectError:    true,
			errorMsg:       "campaign name cannot be empty",
		},
		{
			name:           "invalid campaign type",
			campaignName:   "Test Campaign",
			description:    "Test campaign",
			campaignType:   CampaignType("INVALID"),
			targetingRules: []shared.RuleID{shared.NewRuleID()},
			startDate:      time.Now().Add(24 * time.Hour),
			endDate:        nil,
			budget:         nil,
			createdBy:      shared.NewUserID(),
			settings:       createValidSettings(),
			expectError:    true,
			errorMsg:       "invalid campaign type",
		},
		{
			name:           "no targeting rules",
			campaignName:   "Test Campaign",
			description:    "Test campaign",
			campaignType:   CampaignTypePromotion,
			targetingRules: []shared.RuleID{},
			startDate:      time.Now().Add(24 * time.Hour),
			endDate:        nil,
			budget:         nil,
			createdBy:      shared.NewUserID(),
			settings:       createValidSettings(),
			expectError:    true,
			errorMsg:       "at least one targeting rule must be specified",
		},
		{
			name:           "end date before start date",
			campaignName:   "Test Campaign",
			description:    "Test campaign",
			campaignType:   CampaignTypePromotion,
			targetingRules: []shared.RuleID{shared.NewRuleID()},
			startDate:      time.Now().Add(24 * time.Hour),
			endDate:        &[]time.Time{time.Now()}[0],
			budget:         nil,
			createdBy:      shared.NewUserID(),
			settings:       createValidSettings(),
			expectError:    true,
			errorMsg:       "end date cannot be before start date",
		},
		{
			name:           "negative budget amount",
			campaignName:   "Test Campaign",
			description:    "Test campaign",
			campaignType:   CampaignTypePromotion,
			targetingRules: []shared.RuleID{shared.NewRuleID()},
			startDate:      time.Now().Add(24 * time.Hour),
			endDate:        nil,
			budget:         &shared.Money{Amount: -100.0, Currency: "EUR"},
			createdBy:      shared.NewUserID(),
			settings:       createValidSettings(),
			expectError:    true,
			errorMsg:       "budget amount must be positive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			campaign, err := NewCampaign(
				tt.campaignName,
				tt.description,
				tt.campaignType,
				tt.targetingRules,
				tt.startDate,
				tt.endDate,
				tt.budget,
				tt.createdBy,
				tt.settings,
			)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Nil(t, campaign)
			} else {
				require.NoError(t, err)
				require.NotNil(t, campaign)
				assert.Equal(t, tt.campaignName, campaign.Name())
				assert.Equal(t, tt.description, campaign.Description())
				assert.Equal(t, tt.campaignType, campaign.CampaignType())
				assert.Equal(t, CampaignStatusDraft, campaign.Status())
				assert.Equal(t, tt.startDate, campaign.StartDate())
				assert.Equal(t, tt.endDate, campaign.EndDate())
				assert.Equal(t, tt.budget, campaign.Budget())
				assert.Equal(t, tt.createdBy, campaign.CreatedBy())
				assert.Equal(t, 1, campaign.Version())
				assert.NotEmpty(t, campaign.ID().String())
				assert.False(t, campaign.CreatedAt().IsZero())
				assert.False(t, campaign.UpdatedAt().IsZero())
			}
		})
	}
}

func TestCampaign_Activate(t *testing.T) {
	tests := []struct {
		name        string
		status      CampaignStatus
		startDate   time.Time
		expectError bool
		errorMsg    string
	}{
		{
			name:        "activate draft campaign with future start date",
			status:      CampaignStatusDraft,
			startDate:   time.Now().Add(24 * time.Hour),
			expectError: true,
			errorMsg:    "cannot activate campaign with future start date",
		},
		{
			name:        "activate draft campaign with past start date",
			status:      CampaignStatusDraft,
			startDate:   time.Now().Add(-24 * time.Hour),
			expectError: false,
		},
		{
			name:        "activate paused campaign",
			status:      CampaignStatusPaused,
			startDate:   time.Now().Add(-24 * time.Hour),
			expectError: false,
		},
		{
			name:        "activate active campaign",
			status:      CampaignStatusActive,
			startDate:   time.Now().Add(-24 * time.Hour),
			expectError: true,
			errorMsg:    "cannot activate campaign in current status",
		},
		{
			name:        "activate completed campaign",
			status:      CampaignStatusCompleted,
			startDate:   time.Now().Add(-24 * time.Hour),
			expectError: true,
			errorMsg:    "cannot activate campaign in current status",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			campaign := createTestCampaign(t)
			campaign.status = tt.status
			campaign.startDate = tt.startDate

			err := campaign.Activate()

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
				assert.Equal(t, CampaignStatusActive, campaign.Status())
				assert.True(t, campaign.UpdatedAt().After(campaign.CreatedAt()))
			}
		})
	}
}

func TestCampaign_Pause(t *testing.T) {
	tests := []struct {
		name        string
		status      CampaignStatus
		expectError bool
		errorMsg    string
	}{
		{
			name:        "pause active campaign",
			status:      CampaignStatusActive,
			expectError: false,
		},
		{
			name:        "pause draft campaign",
			status:      CampaignStatusDraft,
			expectError: true,
			errorMsg:    "cannot pause campaign in current status",
		},
		{
			name:        "pause paused campaign",
			status:      CampaignStatusPaused,
			expectError: true,
			errorMsg:    "cannot pause campaign in current status",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			campaign := createTestCampaign(t)
			campaign.status = tt.status

			err := campaign.Pause()

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
				assert.Equal(t, CampaignStatusPaused, campaign.Status())
				assert.True(t, campaign.UpdatedAt().After(campaign.CreatedAt()))
			}
		})
	}
}

func TestCampaign_Complete(t *testing.T) {
	tests := []struct {
		name        string
		status      CampaignStatus
		expectError bool
		errorMsg    string
	}{
		{
			name:        "complete active campaign",
			status:      CampaignStatusActive,
			expectError: false,
		},
		{
			name:        "complete paused campaign",
			status:      CampaignStatusPaused,
			expectError: false,
		},
		{
			name:        "complete draft campaign",
			status:      CampaignStatusDraft,
			expectError: false,
		},
		{
			name:        "complete already completed campaign",
			status:      CampaignStatusCompleted,
			expectError: true,
			errorMsg:    "campaign already completed or cancelled",
		},
		{
			name:        "complete cancelled campaign",
			status:      CampaignStatusCancelled,
			expectError: true,
			errorMsg:    "campaign already completed or cancelled",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			campaign := createTestCampaign(t)
			campaign.status = tt.status

			err := campaign.Complete()

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
				assert.Equal(t, CampaignStatusCompleted, campaign.Status())
				assert.True(t, campaign.UpdatedAt().After(campaign.CreatedAt()))
			}
		})
	}
}

func TestCampaign_Cancel(t *testing.T) {
	tests := []struct {
		name        string
		status      CampaignStatus
		reason      string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "cancel active campaign with reason",
			status:      CampaignStatusActive,
			reason:      "Budget exceeded",
			expectError: false,
		},
		{
			name:        "cancel draft campaign with reason",
			status:      CampaignStatusDraft,
			reason:      "Strategy change",
			expectError: false,
		},
		{
			name:        "cancel campaign without reason",
			status:      CampaignStatusActive,
			reason:      "",
			expectError: true,
			errorMsg:    "cancellation reason cannot be empty",
		},
		{
			name:        "cancel already completed campaign",
			status:      CampaignStatusCompleted,
			reason:      "Test reason",
			expectError: true,
			errorMsg:    "campaign already completed or cancelled",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			campaign := createTestCampaign(t)
			campaign.status = tt.status

			err := campaign.Cancel(tt.reason)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
				assert.Equal(t, CampaignStatusCancelled, campaign.Status())
				assert.True(t, campaign.UpdatedAt().After(campaign.CreatedAt()))
			}
		})
	}
}

func TestCampaign_UpdateTargetingRules(t *testing.T) {
	campaign := createTestCampaign(t)
	originalRules := campaign.TargetingRules()

	newRules := []shared.RuleID{
		shared.NewRuleID(),
		shared.NewRuleID(),
	}

	err := campaign.UpdateTargetingRules(newRules)
	require.NoError(t, err)

	assert.Equal(t, newRules, campaign.TargetingRules())
	assert.True(t, campaign.UpdatedAt().After(campaign.CreatedAt()))
	assert.Greater(t, campaign.Version(), 1)

	// Test with empty rules
	err = campaign.UpdateTargetingRules([]shared.RuleID{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "at least one targeting rule must be specified")
}

func TestCampaign_UpdateBudget(t *testing.T) {
	campaign := createTestCampaign(t)

	newBudget := &shared.Money{Amount: 2000.0, Currency: "USD"}
	err := campaign.UpdateBudget(newBudget)
	require.NoError(t, err)

	assert.Equal(t, newBudget, campaign.Budget())
	assert.True(t, campaign.UpdatedAt().After(campaign.CreatedAt()))
	assert.Greater(t, campaign.Version(), 1)

	// Test with negative budget
	negativeBudget := &shared.Money{Amount: -100.0, Currency: "EUR"}
	err = campaign.UpdateBudget(negativeBudget)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "budget amount must be positive")
}

func TestCampaign_UpdateSettings(t *testing.T) {
	campaign := createTestCampaign(t)

	newSettings := createValidSettings()
	newSettings.TargetAudience = []string{"new-segment"}

	err := campaign.UpdateSettings(newSettings)
	require.NoError(t, err)

	assert.Equal(t, newSettings, campaign.Settings())
	assert.True(t, campaign.UpdatedAt().After(campaign.CreatedAt()))
	assert.Greater(t, campaign.Version(), 1)
}

func TestCampaign_TrackEvent(t *testing.T) {
	campaign := createTestCampaign(t)
	customerID := shared.NewCustomerID()

	eventData := map[string]interface{}{
		"source": "email",
		"device": "mobile",
	}

	err := campaign.TrackEvent(CampaignEventTypeImpression, &customerID, eventData)
	require.NoError(t, err)

	// Verify metrics were updated
	metrics := campaign.Metrics()
	assert.Equal(t, int64(1), metrics.Impressions)
	assert.True(t, metrics.LastUpdated.After(campaign.CreatedAt()))
}

func TestCampaign_GetEvents(t *testing.T) {
	campaign := createTestCampaign(t)

	// Initially no events
	events := campaign.GetEvents()
	assert.Empty(t, events)

	// Activate campaign to generate event
	err := campaign.Activate()
	require.NoError(t, err)

	// Get events
	events = campaign.GetEvents()
	assert.Len(t, events, 1)
	assert.Equal(t, "CampaignActivated", events[0].EventType())

	// Events should be cleared after retrieval
	events = campaign.GetEvents()
	assert.Empty(t, events)
}

func TestCampaign_IsActive(t *testing.T) {
	campaign := createTestCampaign(t)

	// Draft campaign is not active
	assert.False(t, campaign.IsActive())

	// Activate campaign
	err := campaign.Activate()
	require.NoError(t, err)
	assert.True(t, campaign.IsActive())

	// Pause campaign
	err = campaign.Pause()
	require.NoError(t, err)
	assert.False(t, campaign.IsActive())
}

func TestCampaign_HasExceededBudget(t *testing.T) {
	campaign := createTestCampaign(t)
	campaign.budget = &shared.Money{Amount: 1000.0, Currency: "EUR"}

	// Initially no cost
	assert.False(t, campaign.HasExceededBudget())

	// Add cost that exceeds budget
	campaign.metrics.Cost = shared.Money{Amount: 1100.0, Currency: "EUR"}
	assert.True(t, campaign.HasExceededBudget())
}

func TestCampaign_IsApproachingBudgetLimit(t *testing.T) {
	campaign := createTestCampaign(t)
	campaign.budget = &shared.Money{Amount: 1000.0, Currency: "EUR"}

	// Initially no cost
	assert.False(t, campaign.IsApproachingBudgetLimit())

	// Add cost that approaches budget (80% threshold)
	campaign.metrics.Cost = shared.Money{Amount: 800.0, Currency: "EUR"}
	assert.True(t, campaign.IsApproachingBudgetLimit())
}

func TestCampaignType_String(t *testing.T) {
	tests := []struct {
		campaignType CampaignType
		expected     string
	}{
		{CampaignTypePromotion, "PROMOTION"},
		{CampaignTypeLoyalty, "LOYALTY"},
		{CampaignTypeCoupon, "COUPON"},
		{CampaignTypeSegmentation, "SEGMENTATION"},
		{CampaignTypeRetargeting, "RETARGETING"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.campaignType.String())
		})
	}
}

func TestCampaignStatus_String(t *testing.T) {
	tests := []struct {
		status   CampaignStatus
		expected string
	}{
		{CampaignStatusDraft, "DRAFT"},
		{CampaignStatusActive, "ACTIVE"},
		{CampaignStatusPaused, "PAUSED"},
		{CampaignStatusCompleted, "COMPLETED"},
		{CampaignStatusCancelled, "CANCELLED"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.status.String())
		})
	}
}

func TestParseCampaignType(t *testing.T) {
	tests := []struct {
		input    string
		expected CampaignType
		hasError bool
	}{
		{"PROMOTION", CampaignTypePromotion, false},
		{"LOYALTY", CampaignTypeLoyalty, false},
		{"COUPON", CampaignTypeCoupon, false},
		{"SEGMENTATION", CampaignTypeSegmentation, false},
		{"RETARGETING", CampaignTypeRetargeting, false},
		{"INVALID", "", true},
		{"", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := ParseCampaignType(tt.input)
			if tt.hasError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestParseCampaignStatus(t *testing.T) {
	tests := []struct {
		input    string
		expected CampaignStatus
		hasError bool
	}{
		{"DRAFT", CampaignStatusDraft, false},
		{"ACTIVE", CampaignStatusActive, false},
		{"PAUSED", CampaignStatusPaused, false},
		{"COMPLETED", CampaignStatusCompleted, false},
		{"CANCELLED", CampaignStatusCancelled, false},
		{"INVALID", "", true},
		{"", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := ParseCampaignStatus(tt.input)
			if tt.hasError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

// Helper functions

func createTestCampaign(t *testing.T) *Campaign {
	campaign, err := NewCampaign(
		"Test Campaign",
		"Test campaign description",
		CampaignTypePromotion,
		[]shared.RuleID{shared.NewRuleID()},
		time.Now().Add(24*time.Hour),
		nil,
		&shared.Money{Amount: 1000.0, Currency: "EUR"},
		shared.NewUserID(),
		createValidSettings(),
	)
	require.NoError(t, err)
	return campaign
}

func createValidSettings() CampaignSettings {
	return CampaignSettings{
		TargetAudience: []string{"test-audience"},
		Channels:       []Channel{ChannelEmail, ChannelWeb},
		Frequency:      FrequencyDaily,
		MaxImpressions: 1000,
		Personalization: PersonalizationConfig{
			Enabled: false,
			Rules:   []shared.RuleID{},
		},
	}
}
