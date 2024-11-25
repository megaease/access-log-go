package accesslog

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

const (
	topicAccessLog = ""
)

type (
	eventHub struct {
		config        *Config
		kafkaProducer sarama.SyncProducer

		done        chan struct{}
		accessLogCh chan *AccessLog
	}
)

// New creates a new EventHub.
func newEventHub(config *Config) (*eventHub, error) {
	conf := config.Kafka
	tlsConfig, err := loadTLSConfig("", conf.KafkaCertfile, conf.KafkaKeyfile)
	if err != nil {
		return nil, fmt.Errorf("load tls config failed: %v", err)
	}

	logrus.Infof("hub center address: %v", conf.KafkaAddresses)
	logrus.Infof("cert files: %s, %s", conf.KafkaCertfile, conf.KafkaKeyfile)

	// Kafka producer configuration
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.RequiredAcks = sarama.WaitForAll
	kafkaConfig.Producer.Retry.Max = 5
	kafkaConfig.Producer.Return.Successes = true
	kafkaConfig.Net.TLS.Enable = true
	kafkaConfig.Net.TLS.Config = tlsConfig
	kafkaConfig.Net.SASL.Enable = true
	kafkaConfig.Net.SASL.User = conf.KafkaUsername
	kafkaConfig.Net.SASL.Password = conf.KafkaPassword
	kafkaConfig.Net.SASL.Mechanism = sarama.SASLTypePlaintext

	producer, err := sarama.NewSyncProducer(conf.KafkaAddresses, kafkaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %v", err)
	}

	h := &eventHub{
		config:        config,
		kafkaProducer: producer,
		done:          make(chan struct{}),
		accessLogCh:   make(chan *AccessLog, 10000),
	}

	go h.run()
	return h, nil
}

func (h *eventHub) run() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	allLogs := []*AccessLog{}
	for {
		select {
		case <-h.done:
			h.sendEvent(allLogs)
			return
		case accessLog := <-h.accessLogCh:
			allLogs = append(allLogs, accessLog)
		case <-ticker.C:
			h.sendEvent(allLogs)
			allLogs = []*AccessLog{}
			ticker.Reset(5 * time.Second)
		}
	}
}

func (h *eventHub) sendEvent(accessLogs []*AccessLog) {
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

func (h *eventHub) Send(accessLog *AccessLog) {
	select {
	case h.accessLogCh <- accessLog:
	default:
		logrus.Error("event hub is full")
	}
}

func (h *eventHub) Close() {
	close(h.done)
}
