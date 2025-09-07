package evaluation

import "fmt"

// Service is the main application service for rule evaluation.
type Service struct {
	strategies map[string]EvaluationStrategy
}

// NewService creates a new EvaluationService.
func NewService(strategies map[string]EvaluationStrategy) *Service {
	return &Service{strategies: strategies}
}

// GetStrategyForCategory returns the appropriate evaluation strategy for a given rule category.
func (s *Service) GetStrategyForCategory(category string) (EvaluationStrategy, error) {
	strategy, ok := s.strategies[category]
	if !ok {
		// A default strategy could be returned here if needed.
		return nil, fmt.Errorf("no evaluation strategy found for category: %s", category)
	}
	return strategy, nil
}
