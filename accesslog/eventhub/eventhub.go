package eventhub

import "megaease/access-log-go/accesslog/api"

type EventHubType string

const (
	EventHubTypeKafka EventHubType = "kafka"
)

type (
	Config struct {
		Type  EventHubType
		Kafka Kafka
	}

	Kafka struct {
		Addresses []string
		Certfile  string
		Keyfile   string
		Username  string
		Password  string
		Topic     string
	}

	EventHub interface {
		Send(*api.AccessLog) error
		Close()
	}
)

func New(config *Config) (EventHub, error) {
	switch config.Type {
	case EventHubTypeKafka:
		return newKafka(config)
	}

	return nil, nil
}
