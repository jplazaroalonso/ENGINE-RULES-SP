package nats

import (
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

// ConnectToNATS connects to NATS server
func ConnectToNATS(url string) (*nats.Conn, error) {
	opts := []nats.Option{
		nats.Name("Customer Management Service"),
		nats.Timeout(10 * time.Second),
		nats.ReconnectWait(2 * time.Second),
		nats.MaxReconnects(-1),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			fmt.Printf("Disconnected from NATS: %v\n", err)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			fmt.Printf("Reconnected to NATS at %s\n", nc.ConnectedUrl())
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			fmt.Printf("NATS connection closed\n")
		}),
	}

	conn, err := nats.Connect(url, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	return conn, nil
}

// InitializeJetStream initializes JetStream with required streams
func InitializeJetStream(conn *nats.Conn) error {
	js, err := conn.JetStream()
	if err != nil {
		return fmt.Errorf("failed to create JetStream context: %w", err)
	}

	// Create streams for different event types
	streams := []struct {
		name     string
		subjects []string
	}{
		{
			name:     "CUSTOMER_EVENTS",
			subjects: []string{"events.Customer.*"},
		},
		{
			name:     "SEGMENT_EVENTS",
			subjects: []string{"events.CustomerSegment.*"},
		},
		{
			name:     "ANALYTICS_EVENTS",
			subjects: []string{"events.CustomerAnalytics.*"},
		},
	}

	for _, stream := range streams {
		streamConfig := &nats.StreamConfig{
			Name:      stream.name,
			Subjects:  stream.subjects,
			Retention: nats.WorkQueuePolicy,
			MaxAge:    24 * time.Hour,
			Storage:   nats.FileStorage,
			Replicas:  1,
		}

		_, err := js.AddStream(streamConfig)
		if err != nil {
			// Stream might already exist, check if it's a different error
			if err != nats.ErrStreamNameAlreadyInUse {
				return fmt.Errorf("failed to create stream %s: %w", stream.name, err)
			}
		}
	}

	return nil
}
