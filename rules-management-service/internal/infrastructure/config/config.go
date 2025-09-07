package config

import (
	"os"
)

// Config holds the application configuration.
type Config struct {
	Server    ServerConfig
	Telemetry TelemetryConfig
	Database  DatabaseConfig
	NATS      NATSConfig
}

// ServerConfig holds the server configuration.
type ServerConfig struct {
	Port string
}

// TelemetryConfig holds the telemetry configuration.
type TelemetryConfig struct {
	ServiceName string
	Exporter    string // e.g., "stdout", "jaeger", "otlp"
}

// DatabaseConfig holds the database configuration.
type DatabaseConfig struct {
	DSN string
}

// NATSConfig holds the NATS configuration.
type NATSConfig struct {
	URL string
}

// DefaultConfig returns the default configuration.
func DefaultConfig() *Config {
	// Get environment variables with defaults
	serverPort := getEnv("SERVER_PORT", "8080")
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbName := getEnv("DB_NAME", "rules_dev")
	dbUser := getEnv("DB_USER", "user")
	dbPassword := getEnv("DB_PASSWORD", "password")
	dbSSLMode := getEnv("DB_SSL_MODE", "disable")
	natsURL := getEnv("NATS_URL", "nats://localhost:4222")
	telemetryServiceName := getEnv("TELEMETRY_SERVICE_NAME", "rules-management-service")
	telemetryExporter := getEnv("TELEMETRY_EXPORTER", "stdout")

	// Build DSN
	dsn := "host=" + dbHost + " user=" + dbUser + " password=" + dbPassword + " dbname=" + dbName + " port=" + dbPort + " sslmode=" + dbSSLMode + " TimeZone=UTC"

	return &Config{
		Server: ServerConfig{
			Port: serverPort,
		},
		Telemetry: TelemetryConfig{
			ServiceName: telemetryServiceName,
			Exporter:    telemetryExporter,
		},
		Database: DatabaseConfig{
			DSN: dsn,
		},
		NATS: NATSConfig{
			URL: natsURL,
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
