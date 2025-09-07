package config

import (
	"os"
)

// Config holds the application configuration.
type Config struct {
	Server    ServerConfig
	Telemetry TelemetryConfig
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

// DefaultConfig returns the default configuration.
func DefaultConfig() *Config {
	// Get environment variables with defaults
	serverPort := getEnv("SERVER_PORT", "8082")
	telemetryServiceName := getEnv("TELEMETRY_SERVICE_NAME", "rules-calculator-service")
	telemetryExporter := getEnv("TELEMETRY_EXPORTER", "stdout")

	return &Config{
		Server: ServerConfig{
			Port: serverPort,
		},
		Telemetry: TelemetryConfig{
			ServiceName: telemetryServiceName,
			Exporter:    telemetryExporter,
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
