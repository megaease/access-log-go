package eventhub

import (
	"encoding/json"
	"fmt"
	"megaease/access-log-go/accesslog/api"
	"megaease/access-log-go/accesslog/utils"
	"time"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

const (
	topicAccessLog = ""
)

type (
	eventHubKafka struct {
		config        *Config
		kafkaProducer sarama.SyncProducer

		done        chan struct{}
		accessLogCh chan *api.AccessLog
	}
)

// New creates a new EventHub.
func newKafka(config *Config) (EventHub, error) {
	conf := config.Kafka
	tlsConfig, err := utils.LoadTLSConfig("", conf.Certfile, conf.Keyfile)
	if err != nil {
		return nil, fmt.Errorf("load tls config failed: %v", err)
	}

	logrus.Infof("hub center address: %v", conf.Addresses)
	logrus.Infof("cert files: %s, %s", conf.Certfile, conf.Keyfile)

	// Kafka producer configuration
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.RequiredAcks = sarama.WaitForAll
	kafkaConfig.Producer.Retry.Max = 5
	kafkaConfig.Producer.Return.Successes = true
	kafkaConfig.Net.TLS.Enable = true
	kafkaConfig.Net.TLS.Config = tlsConfig
	kafkaConfig.Net.SASL.Enable = true
	kafkaConfig.Net.SASL.User = conf.Username
	kafkaConfig.Net.SASL.Password = conf.Password
	kafkaConfig.Net.SASL.Mechanism = sarama.SASLTypePlaintext

	producer, err := sarama.NewSyncProducer(conf.Addresses, kafkaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %v", err)
	}

	h := &eventHubKafka{
		config:        config,
		kafkaProducer: producer,
		done:          make(chan struct{}),
		accessLogCh:   make(chan *api.AccessLog, 10000),
	}

	go h.run()
	return h, nil
}

func (h *eventHubKafka) run() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	allLogs := []*api.AccessLog{}
	for {
		select {
		case <-h.done:
			h.sendEvent(allLogs)
			return
		case accessLog := <-h.accessLogCh:
			allLogs = append(allLogs, accessLog)
		case <-ticker.C:
			h.sendEvent(allLogs)
			allLogs = []*api.AccessLog{}
			ticker.Reset(5 * time.Second)
		}
	}
}

func (h *eventHubKafka) sendEvent(accessLogs []*api.AccessLog) {
	if len(accessLogs) == 0 {
		return
	}
	buff, err := json.Marshal(accessLogs)
	if err != nil {
		logrus.Errorf("marshal %s %#v to json failed: %v", topicAccessLog, accessLogs, err)
		return
	}

	kafkaMsg := &sarama.ProducerMessage{
		Topic: topicAccessLog,
		Value: sarama.ByteEncoder(buff),
	}

	_, _, err = h.kafkaProducer.SendMessage(kafkaMsg)
	if err != nil {
	}
}

func (h *eventHubKafka) Send(accessLog *api.AccessLog) error {
	select {
	case h.accessLogCh <- accessLog:
		return nil
	default:
		return fmt.Errorf("access log channel is full")
	}
}

func (h *eventHubKafka) Close() {
	close(h.done)
}