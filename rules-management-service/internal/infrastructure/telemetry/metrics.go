package telemetry

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// RulesCreated is a counter for the total number of rules created.
	RulesCreated = promauto.NewCounter(prometheus.CounterOpts{
		Name: "rules_management_rules_created_total",
		Help: "The total number of rules created",
	})
	// RuleCreationDuration is a histogram of the duration of rule creation.
	RuleCreationDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "rules_management_rule_creation_duration_seconds",
		Help: "The duration of rule creation",
	})
	// DBQueryDuration is a histogram of the duration of database queries.
	DBQueryDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "rules_management_db_query_duration_seconds",
		Help: "The duration of database queries.",
	}, []string{"operation"})
)
