package nats

import (
	"context"
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/rules-management-service/internal/application/commands"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/rules-management-service/internal/infrastructure/config"
)

// CommandSubscriber listens for commands on NATS subjects and executes them.
type CommandSubscriber struct {
	conn              *nats.Conn
	js                nats.JetStreamContext
	createRuleHandler *commands.CreateRuleHandler
}

// NewCommandSubscriber creates a new NATS command subscriber.
func NewCommandSubscriber(cfg config.NATSConfig, createRuleHandler *commands.CreateRuleHandler) (*CommandSubscriber, error) {
	conn, err := nats.Connect(cfg.URL)
	if err != nil {
		return nil, err
	}

	js, err := conn.JetStream()
	if err != nil {
		return nil, err
	}

	return &CommandSubscriber{
		conn:              conn,
		js:                js,
		createRuleHandler: createRuleHandler,
	}, nil
}

// Start starts the subscriber to listen for commands.
func (s *CommandSubscriber) Start() {
	// Subscribe to the subject for creating rules
	_, err := s.js.Subscribe("commands.rule.create", func(msg *nats.Msg) {
		s.handleCreateRuleCommand(msg)
	}, nats.Durable("rule-command-consumer"), nats.AckWait(s.conn.Opts.Timeout))
	if err != nil {
		log.Printf("Error subscribing to commands.rule.create: %v", err)
	}

	log.Println("Listening for rule creation commands on NATS subject 'commands.rule.create'")
}

func (s *CommandSubscriber) handleCreateRuleCommand(msg *nats.Msg) {
	var cmd commands.CreateRuleCommand
	if err := json.Unmarshal(msg.Data, &cmd); err != nil {
		log.Printf("Error unmarshaling create rule command: %v", err)
		if err := msg.Nak(); err != nil {
			log.Printf("Error sending Nak: %v", err)
		}
		return
	}

	log.Printf("Received create rule command for: %s", cmd.Name)
	if _, err := s.createRuleHandler.Handle(context.Background(), cmd); err != nil {
		log.Printf("Error handling create rule command: %v", err)
		if err := msg.Nak(); err != nil {
			log.Printf("Error sending Nak: %v", err)
		}
		return
	}

	if err := msg.Ack(); err != nil {
		log.Printf("Error sending Ack: %v", err)
	}
	log.Printf("Successfully processed create rule command for: %s", cmd.Name)
}

// Close closes the NATS connection.
func (s *CommandSubscriber) Close() {
	if s.conn != nil {
		s.conn.Close()
	}
}
