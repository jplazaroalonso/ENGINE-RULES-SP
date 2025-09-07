package config

import "time"

// Config holds the application configuration
type Config struct {
	NATS NATSConfig `yaml:"nats"`
	App  AppConfig  `yaml:"app"`
}

// NATSConfig holds NATS-related configuration
type NATSConfig struct {
	URL         string        `yaml:"url"`
	ClusterName string        `yaml:"clusterName"`
	ClientID    string        `yaml:"clientId"`
	Timeout     time.Duration `yaml:"timeout"`
}

// AppConfig holds application-level configuration
type AppConfig struct {
	ReplicationEnabled bool `yaml:"replicationEnabled"`
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		NATS: NATSConfig{
			URL:         "nats://localhost:4222",
			ClusterName: "rules-engine-cluster",
			ClientID:    "rules-management-service",
			Timeout:     5 * time.Second,
		},
		App: AppConfig{
			ReplicationEnabled: true, // Default to enabled
		},
	}
}
