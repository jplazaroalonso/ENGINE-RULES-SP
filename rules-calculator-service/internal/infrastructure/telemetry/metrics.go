package telemetry

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// CalculationsTotal is a counter for the total number of calculations processed.
	CalculationsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "rules_calculator_calculations_total",
		Help: "The total number of calculations processed",
	})
	// CalculationDuration is a histogram of the duration of calculations.
	CalculationDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "rules_calculator_calculation_duration_seconds",
		Help: "The duration of calculations",
	})
	// RuleEvaluationDuration is a histogram of the duration of individual rule evaluations.
	RuleEvaluationDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "rules_calculator_rule_evaluation_duration_seconds",
		Help: "The duration of individual rule evaluations.",
	}, []string{"rule_id"})
)
