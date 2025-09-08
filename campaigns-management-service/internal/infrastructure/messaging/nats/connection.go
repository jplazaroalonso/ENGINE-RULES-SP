package nats

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

// NATSConfig holds NATS connection configuration
type NATSConfig struct {
	URL      string
	Username string
	Password string
	Token    string
}

// ConnectNATS establishes a connection to NATS server
func ConnectNATS(config NATSConfig) (*nats.Conn, error) {
	opts := []nats.Option{
		nats.Name("campaigns-management-service"),tid
		nats.Timeout(10 * time.Second),
		nats.ReconnectWait(2 * time.Second),
		nats.MaxReconnects(-1), // Unlimited reconnects
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			log.Printf("NATS disconnected: %v", err)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Printf("NATS reconnected to %v", nc.ConnectedUrl())
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			log.Printf("NATS connection closed")
		}),
	}

	// Add authentication if provided
	if config.Username != "" && config.Password != "" {
		opts = append(opts, nats.UserInfo(config.Username, config.Password))
	}
	if config.Token != "" {
		opts = append(opts, nats.Token(config.Token))
	}

	// Connect to NATS
	conn, err := nats.Connect(config.URL, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	log.Printf("Connected to NATS at %s", conn.ConnectedUrl())
	return conn, nil
}
