package external

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/campaigns-management-service/internal/domain/campaign"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/campaigns-management-service/internal/domain/shared"
)

// RulesClient implements the CampaignTargetingService interface
type RulesClient struct {
	baseURL    string
	httpClient *http.Client
	apiKey     string
}

// NewRulesClient creates a new rules client
func NewRulesClient(baseURL, apiKey string) *RulesClient {
	return &RulesClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		apiKey: apiKey,
	}
}

// EvaluateTargeting evaluates customer targeting using rules engine
func (c *RulesClient) EvaluateTargeting(ctx context.Context, campaignID campaign.CampaignID, customerID shared.CustomerID) (bool, error) {
	// Get campaign to retrieve targeting rules
	// Note: In a real implementation, you would need access to the campaign repository
	// For now, we'll assume the targeting rules are passed separately

	// Create evaluation request
	request := EvaluationRequest{
		CustomerID: customerID.String(),
		Context:    make(map[string]interface{}),
	}

	// Make HTTP request to rules engine
	response, err := c.makeRequest(ctx, "POST", "/api/v1/evaluate", request)
	if err != nil {
		return false, shared.NewInfrastructureError("failed to evaluate targeting rules", err)
	}

	// Parse response
	var evalResponse EvaluationResponse
	if err := json.Unmarshal(response, &evalResponse); err != nil {
		return false, shared.NewInfrastructureError("failed to parse evaluation response", err)
	}

	return evalResponse.Result, nil
}

// GetTargetAudience gets the target audience for a campaign
func (c *RulesClient) GetTargetAudience(ctx context.Context, campaignID campaign.CampaignID) ([]shared.CustomerID, error) {
	// Make HTTP request to get target audience
	response, err := c.makeRequest(ctx, "GET", fmt.Sprintf("/api/v1/campaigns/%s/audience", campaignID.String()), nil)
	if err != nil {
		return nil, shared.NewInfrastructureError("failed to get target audience", err)
	}

	// Parse response
	var audienceResponse AudienceResponse
	if err := json.Unmarshal(response, &audienceResponse); err != nil {
		return nil, shared.NewInfrastructureError("failed to parse audience response", err)
	}

	// Convert to domain types
	customerIDs := make([]shared.CustomerID, len(audienceResponse.CustomerIDs))
	for i, id := range audienceResponse.CustomerIDs {
		customerID, err := shared.NewCustomerIDFromString(id)
		if err != nil {
			return nil, shared.NewValidationError("invalid customer ID", err)
		}
		customerIDs[i] = customerID
	}

	return customerIDs, nil
}

// UpdateTargetingRules updates targeting rules for a campaign
func (c *RulesClient) UpdateTargetingRules(ctx context.Context, campaignID campaign.CampaignID, rules []shared.RuleID) error {
	// Convert rule IDs to strings
	ruleStrings := make([]string, len(rules))
	for i, rule := range rules {
		ruleStrings[i] = rule.String()
	}

	request := UpdateTargetingRequest{
		RuleIDs: ruleStrings,
	}

	// Make HTTP request to update targeting rules
	_, err := c.makeRequest(ctx, "PUT", fmt.Sprintf("/api/v1/campaigns/%s/targeting", campaignID.String()), request)
	if err != nil {
		return shared.NewInfrastructureError("failed to update targeting rules", err)
	}

	return nil
}

// ValidateTargetingRules validates targeting rules
func (c *RulesClient) ValidateTargetingRules(ctx context.Context, rules []shared.RuleID) error {
	// Convert rule IDs to strings
	ruleStrings := make([]string, len(rules))
	for i, rule := range rules {
		ruleStrings[i] = rule.String()
	}

	request := ValidationRequest{
		RuleIDs: ruleStrings,
	}

	// Make HTTP request to validate rules
	response, err := c.makeRequest(ctx, "POST", "/api/v1/rules/validate", request)
	if err != nil {
		return shared.NewInfrastructureError("failed to validate targeting rules", err)
	}

	// Parse response
	var validationResponse ValidationResponse
	if err := json.Unmarshal(response, &validationResponse); err != nil {
		return shared.NewInfrastructureError("failed to parse validation response", err)
	}

	// Check if validation passed
	if !validationResponse.Valid {
		return shared.NewValidationError("targeting rules validation failed", fmt.Errorf("validation errors: %v", validationResponse.Errors))
	}

	return nil
}

// makeRequest makes an HTTP request to the rules engine
func (c *RulesClient) makeRequest(ctx context.Context, method, path string, body interface{}) ([]byte, error) {
	url := c.baseURL + path

	var reqBody []byte
	var err error
	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	if reqBody != nil {
		req.Body = http.NoBody // We'll set the body properly
		req.GetBody = func() (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader(reqBody)), nil
		}
		req.ContentLength = int64(len(reqBody))
	}

	// Make request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// Request/Response types for rules engine API

type EvaluationRequest struct {
	CustomerID string                 `json:"customerId"`
	Context    map[string]interface{} `json:"context"`
}

type EvaluationResponse struct {
	Result bool                   `json:"result"`
	Reason string                 `json:"reason,omitempty"`
	Data   map[string]interface{} `json:"data,omitempty"`
}

type AudienceResponse struct {
	CustomerIDs []string `json:"customerIds"`
	Total       int      `json:"total"`
}

type UpdateTargetingRequest struct {
	RuleIDs []string `json:"ruleIds"`
}

type ValidationRequest struct {
	RuleIDs []string `json:"ruleIds"`
}

type ValidationResponse struct {
	Valid    bool     `json:"valid"`
	Errors   []string `json:"errors,omitempty"`
	Warnings []string `json:"warnings,omitempty"`
}

// MockRulesClient is a mock implementation for testing
type MockRulesClient struct {
	evaluationResults map[string]bool
	audienceResults   map[string][]string
	validationResults map[string]bool
}

// NewMockRulesClient creates a new mock rules client
func NewMockRulesClient() *MockRulesClient {
	return &MockRulesClient{
		evaluationResults: make(map[string]bool),
		audienceResults:   make(map[string][]string),
		validationResults: make(map[string]bool),
	}
}

// SetEvaluationResult sets the evaluation result for a customer
func (m *MockRulesClient) SetEvaluationResult(customerID string, result bool) {
	m.evaluationResults[customerID] = result
}

// SetAudienceResult sets the audience result for a campaign
func (m *MockRulesClient) SetAudienceResult(campaignID string, customerIDs []string) {
	m.audienceResults[campaignID] = customerIDs
}

// SetValidationResult sets the validation result for rules
func (m *MockRulesClient) SetValidationResult(ruleIDs []string, result bool) {
	key := fmt.Sprintf("%v", ruleIDs)
	m.validationResults[key] = result
}

// EvaluateTargeting evaluates customer targeting (mock implementation)
func (m *MockRulesClient) EvaluateTargeting(ctx context.Context, campaignID campaign.CampaignID, customerID shared.CustomerID) (bool, error) {
	if result, exists := m.evaluationResults[customerID.String()]; exists {
		return result, nil
	}
	// Default to true for testing
	return true, nil
}

// GetTargetAudience gets the target audience (mock implementation)
func (m *MockRulesClient) GetTargetAudience(ctx context.Context, campaignID campaign.CampaignID) ([]shared.CustomerID, error) {
	if customerIDs, exists := m.audienceResults[campaignID.String()]; exists {
		result := make([]shared.CustomerID, len(customerIDs))
		for i, id := range customerIDs {
			customerID, err := shared.NewCustomerIDFromString(id)
			if err != nil {
				return nil, shared.NewValidationError("invalid customer ID", err)
			}
			result[i] = customerID
		}
		return result, nil
	}
	// Default to empty list for testing
	return []shared.CustomerID{}, nil
}

// UpdateTargetingRules updates targeting rules (mock implementation)
func (m *MockRulesClient) UpdateTargetingRules(ctx context.Context, campaignID campaign.CampaignID, rules []shared.RuleID) error {
	// Mock implementation - always succeeds
	return nil
}

// ValidateTargetingRules validates targeting rules (mock implementation)
func (m *MockRulesClient) ValidateTargetingRules(ctx context.Context, rules []shared.RuleID) error {
	ruleStrings := make([]string, len(rules))
	for i, rule := range rules {
		ruleStrings[i] = rule.String()
	}
	key := fmt.Sprintf("%v", ruleStrings)

	if result, exists := m.validationResults[key]; exists {
		if !result {
			return shared.NewValidationError("mock validation failed", nil)
		}
	}

	// Default to valid for testing
	return nil
}
