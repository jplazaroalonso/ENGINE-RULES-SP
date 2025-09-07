package telemetry

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// EvaluationsTotal is a counter for the total number of evaluations processed.
	EvaluationsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "rules_evaluation_evaluations_total",
		Help: "The total number of evaluations processed",
	}, []string{"category", "success"})
	// EvaluationDuration is a histogram of the duration of evaluations.
	EvaluationDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "rules_evaluation_evaluation_duration_seconds",
		Help: "The duration of evaluations",
	}, []string{"category"})
)
