package external

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/customer-management-service/internal/domain/shared"
)

// RulesClient implements the rules engine client interface
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

// ValidateRule validates a rule with the rules engine
func (rc *RulesClient) ValidateRule(ctx context.Context, ruleID shared.RuleID) error {
	url := fmt.Sprintf("%s/api/v1/rules/%s/validate", rc.baseURL, ruleID.String())

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+rc.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := rc.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to validate rule: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("rule validation failed with status: %d", resp.StatusCode)
	}

	return nil
}

// GetRule retrieves a rule from the rules engine
func (rc *RulesClient) GetRule(ctx context.Context, ruleID shared.RuleID) (*RuleResponse, error) {
	url := fmt.Sprintf("%s/api/v1/rules/%s", rc.baseURL, ruleID.String())

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+rc.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := rc.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get rule: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, shared.NewNotFoundError("rule not found", nil)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get rule with status: %d", resp.StatusCode)
	}

	var ruleResponse RuleResponse
	if err := json.NewDecoder(resp.Body).Decode(&ruleResponse); err != nil {
		return nil, fmt.Errorf("failed to decode rule response: %w", err)
	}

	return &ruleResponse, nil
}

// EvaluateRule evaluates a rule with given context
func (rc *RulesClient) EvaluateRule(ctx context.Context, ruleID shared.RuleID, context map[string]interface{}) (*RuleEvaluationResponse, error) {
	url := fmt.Sprintf("%s/api/v1/rules/%s/evaluate", rc.baseURL, ruleID.String())

	requestBody := RuleEvaluationRequest{
		Context: context,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+rc.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := rc.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate rule: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("rule evaluation failed with status: %d", resp.StatusCode)
	}

	var evaluationResponse RuleEvaluationResponse
	if err := json.NewDecoder(resp.Body).Decode(&evaluationResponse); err != nil {
		return nil, fmt.Errorf("failed to decode evaluation response: %w", err)
	}

	return &evaluationResponse, nil
}

// ListRules lists rules from the rules engine
func (rc *RulesClient) ListRules(ctx context.Context, criteria RuleListCriteria) (*RuleListResponse, error) {
	url := fmt.Sprintf("%s/api/v1/rules", rc.baseURL)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add query parameters
	q := req.URL.Query()
	if criteria.Page > 0 {
		q.Add("page", fmt.Sprintf("%d", criteria.Page))
	}
	if criteria.Limit > 0 {
		q.Add("limit", fmt.Sprintf("%d", criteria.Limit))
	}
	if criteria.Status != "" {
		q.Add("status", criteria.Status)
	}
	if criteria.Category != "" {
		q.Add("category", criteria.Category)
	}
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Authorization", "Bearer "+rc.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := rc.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to list rules: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list rules with status: %d", resp.StatusCode)
	}

	var listResponse RuleListResponse
	if err := json.NewDecoder(resp.Body).Decode(&listResponse); err != nil {
		return nil, fmt.Errorf("failed to decode list response: %w", err)
	}

	return &listResponse, nil
}

// RuleResponse represents a rule response from the rules engine
type RuleResponse struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Category    string                 `json:"category"`
	Status      string                 `json:"status"`
	Definition  map[string]interface{} `json:"definition"`
	CreatedAt   time.Time              `json:"createdAt"`
	UpdatedAt   time.Time              `json:"updatedAt"`
	Version     int                    `json:"version"`
}

// RuleEvaluationRequest represents a rule evaluation request
type RuleEvaluationRequest struct {
	Context map[string]interface{} `json:"context"`
}

// RuleEvaluationResponse represents a rule evaluation response
type RuleEvaluationResponse struct {
	RuleID      string                 `json:"ruleId"`
	Result      bool                   `json:"result"`
	Score       float64                `json:"score"`
	Details     map[string]interface{} `json:"details"`
	EvaluatedAt time.Time              `json:"evaluatedAt"`
}

// RuleListCriteria represents criteria for listing rules
type RuleListCriteria struct {
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
	Status   string `json:"status"`
	Category string `json:"category"`
}

// RuleListResponse represents a rule list response
type RuleListResponse struct {
	Rules      []RuleResponse `json:"rules"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	TotalPages int            `json:"totalPages"`
}
