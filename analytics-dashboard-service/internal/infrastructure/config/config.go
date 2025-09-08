package config

import (
	"os"
	"strconv"
	"time"
)

// Config represents the application configuration
type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	NATS     NATSConfig     `json:"nats"`
	Cache    CacheConfig    `json:"cache"`
	External ExternalConfig `json:"external"`
}

// ServerConfig represents server configuration
type ServerConfig struct {
	Port         string        `json:"port"`
	ReadTimeout  time.Duration `json:"readTimeout"`
	WriteTimeout time.Duration `json:"writeTimeout"`
	IdleTimeout  time.Duration `json:"idleTimeout"`
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	DSN             string        `json:"dsn"`
	MaxOpenConns    int           `json:"maxOpenConns"`
	MaxIdleConns    int           `json:"maxIdleConns"`
	ConnMaxLifetime time.Duration `json:"connMaxLifetime"`
}

// NATSConfig represents NATS configuration
type NATSConfig struct {
	URL       string `json:"url"`
	ClusterID string `json:"clusterId"`
	ClientID  string `json:"clientId"`
}

// CacheConfig represents cache configuration
type CacheConfig struct {
	RedisURL string        `json:"redisUrl"`
	TTL      time.Duration `json:"ttl"`
}

// ExternalConfig represents external service configuration
type ExternalConfig struct {
	RulesServiceURL     string `json:"rulesServiceUrl"`
	CustomerServiceURL  string `json:"customerServiceUrl"`
	CampaignServiceURL  string `json:"campaignServiceUrl"`
	PromotionServiceURL string `json:"promotionServiceUrl"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         getEnv("PORT", "8080"),
			ReadTimeout:  getDurationEnv("READ_TIMEOUT", 30*time.Second),
			WriteTimeout: getDurationEnv("WRITE_TIMEOUT", 30*time.Second),
			IdleTimeout:  getDurationEnv("IDLE_TIMEOUT", 120*time.Second),
		},
		Database: DatabaseConfig{
			DSN:             getEnv("DATABASE_DSN", "postgres://user:password@localhost:5432/analytics_db?sslmode=disable"),
			MaxOpenConns:    getIntEnv("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getIntEnv("DB_MAX_IDLE_CONNS", 25),
			ConnMaxLifetime: getDurationEnv("DB_CONN_MAX_LIFETIME", 5*time.Minute),
		},
		NATS: NATSConfig{
			URL:       getEnv("NATS_URL", "nats://localhost:4222"),
			ClusterID: getEnv("NATS_CLUSTER_ID", "analytics-cluster"),
			ClientID:  getEnv("NATS_CLIENT_ID", "analytics-service"),
		},
		Cache: CacheConfig{
			RedisURL: getEnv("REDIS_URL", "redis://localhost:6379"),
			TTL:      getDurationEnv("CACHE_TTL", 5*time.Minute),
		},
		External: ExternalConfig{
			RulesServiceURL:     getEnv("RULES_SERVICE_URL", "http://localhost:8081"),
			CustomerServiceURL:  getEnv("CUSTOMER_SERVICE_URL", "http://localhost:8082"),
			CampaignServiceURL:  getEnv("CAMPAIGN_SERVICE_URL", "http://localhost:8083"),
			PromotionServiceURL: getEnv("PROMOTION_SERVICE_URL", "http://localhost:8084"),
		},
	}
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getIntEnv gets an integer environment variable with a default value
func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getDurationEnv gets a duration environment variable with a default value
func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
