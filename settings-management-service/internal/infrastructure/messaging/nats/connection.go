package nats

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/nats-io/nats.go"
)

// ConnectionConfig holds the NATS connection configuration
type ConnectionConfig struct {
	URL              string
	MaxReconnects    int
	ReconnectWait    time.Duration
	Timeout          time.Duration
	PingInterval     time.Duration
	MaxPingsOut      int
	ReconnectBufSize int
}

// NewConnectionConfig creates a new connection config from environment variables
func NewConnectionConfig() *ConnectionConfig {
	return &ConnectionConfig{
		URL:              getEnv("NATS_URL", "nats://localhost:4222"),
		MaxReconnects:    getEnvInt("NATS_MAX_RECONNECTS", -1),
		ReconnectWait:    getEnvDuration("NATS_RECONNECT_WAIT", 2*time.Second),
		Timeout:          getEnvDuration("NATS_TIMEOUT", 5*time.Second),
		PingInterval:     getEnvDuration("NATS_PING_INTERVAL", 2*time.Minute),
		MaxPingsOut:      getEnvInt("NATS_MAX_PINGS_OUT", 2),
		ReconnectBufSize: getEnvInt("NATS_RECONNECT_BUF_SIZE", 8*1024*1024), // 8MB
	}
}

// Connect establishes a connection to NATS
func Connect(config *ConnectionConfig) (*nats.Conn, error) {
	// Set up connection options
	opts := []nats.Option{
		nats.Name("Settings Management Service"),
		nats.MaxReconnects(config.MaxReconnects),
		nats.ReconnectWait(config.ReconnectWait),
		nats.Timeout(config.Timeout),
		nats.PingInterval(config.PingInterval),
		nats.MaxPingsOutstanding(config.MaxPingsOut),
		nats.ReconnectBufSize(config.ReconnectBufSize),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			log.Printf("NATS disconnected: %v", err)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Printf("NATS reconnected to %v", nc.ConnectedUrl())
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			log.Printf("NATS connection closed")
		}),
		nats.ErrorHandler(func(nc *nats.Conn, sub *nats.Subscription, err error) {
			log.Printf("NATS error: %v", err)
		}),
	}

	// Connect to NATS
	conn, err := nats.Connect(config.URL, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	// Wait for connection to be established
	if !conn.IsConnected() {
		return nil, fmt.Errorf("NATS connection not established")
	}

	log.Printf("Successfully connected to NATS at %s", config.URL)
	return conn, nil
}

// HealthCheck performs a health check on the NATS connection
func HealthCheck(conn *nats.Conn) error {
	if conn == nil {
		return fmt.Errorf("NATS connection is nil")
	}

	if !conn.IsConnected() {
		return fmt.Errorf("NATS connection is not connected")
	}

	// Try to publish a test message to verify the connection
	testSubject := "health.check"
	testData := []byte("ping")

	// Use a timeout for the health check
	timeout := time.NewTimer(5 * time.Second)
	defer timeout.Stop()

	done := make(chan error, 1)
	go func() {
		err := conn.Publish(testSubject, testData)
		done <- err
	}()

	select {
	case err := <-done:
		if err != nil {
			return fmt.Errorf("NATS health check failed: %w", err)
		}
		return nil
	case <-timeout.C:
		return fmt.Errorf("NATS health check timed out")
	}
}

// Close closes the NATS connection
func Close(conn *nats.Conn) error {
	if conn != nil {
		conn.Close()
		log.Println("NATS connection closed")
	}
	return nil
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt gets an environment variable as an integer with a default value
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvDuration gets an environment variable as a duration with a default value
func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
