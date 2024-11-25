package accesslog

type (
	Config struct {
		Kafka Kafka
	}

	Kafka struct {
		KafkaAddresses []string
		KafkaCertfile  string
		KafkaKeyfile   string
		KafkaUsername  string
		KafkaPassword  string
	}
)
