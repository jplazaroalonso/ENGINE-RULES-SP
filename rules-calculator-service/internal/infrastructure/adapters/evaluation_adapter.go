package adapters

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// HTTPEvaluationAdapter is an adapter to the rule evaluation service.
type HTTPEvaluationAdapter struct {
	baseURL string
	client  *http.Client
}

// NewHTTPEvaluationAdapter creates a new HTTPEvaluationAdapter.
func NewHTTPEvaluationAdapter(baseURL string) *HTTPEvaluationAdapter {
	return &HTTPEvaluationAdapter{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

type evaluationRequest struct {
	RuleID  string                 `json:"rule_id"`
	Context map[string]interface{} `json:"context"`
}

type evaluationResponse struct {
	Value float64 `json:"value"`
}

// Evaluate evaluates a rule using the rule evaluation service.
func (a *HTTPEvaluationAdapter) Evaluate(ctx context.Context, ruleID string, context map[string]interface{}) (float64, error) {
	// This is a mock implementation.
	// It does not actually call the evaluation service.
	// It returns a dummy value.
	// In a real implementation, this would make an HTTP call to the evaluation service.

	// Seed the random number generator
	// rand.Seed(time.Now().UnixNano())

	// Generate a random float between 10 and 100
	// randomFloat := 10 + rand.Float64()*(100-10)

	// return randomFloat, nil

	reqBody := evaluationRequest{
		RuleID:  ruleID,
		Context: context,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal evaluation request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", a.baseURL+"/v1/evaluate", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return 0, fmt.Errorf("failed to create evaluation request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to call evaluation service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("evaluation service returned non-OK status: %d", resp.StatusCode)
	}

	var evalResp evaluationResponse
	if err := json.NewDecoder(resp.Body).Decode(&evalResp); err != nil {
		return 0, fmt.Errorf("failed to decode evaluation response: %w", err)
	}

	return evalResp.Value, nil
}
