package eventhub

import (
	"fmt"
	"megaease/access-log-go/accesslog/api"
)

// EventHubType is the type of event hub.
type EventHubType string

const (
	// EventHubTypeKafka is the Kafka event hub.
	EventHubTypeKafka EventHubType = "kafka"
	// EventHubTypeMock is the mock event hub.
	EventHubTypeMock EventHubType = "mock"
)

type (
	// Config is the configuration of event hub.
	Config struct {
		Type  EventHubType
		Kafka Kafka
	}

	// Kafka is the configuration of Kafka.
	Kafka struct {
		Addresses []string
		Certfile  string
		Keyfile   string
		Username  string
		Password  string

		Topic string
	}

	// EventHub is the interface of event hub.
	EventHub interface {
		Send(*api.AccessLog) error
		Close()
	}
)

// New creates a new EventHub.
func New(config *Config) (EventHub, error) {
	switch config.Type {
	case EventHubTypeKafka:
		return newKafka(config)
	case EventHubTypeMock:
		return newEventHubMock()
	}
	return nil, fmt.Errorf("unsupported event hub type: %s", config.Type)
}
