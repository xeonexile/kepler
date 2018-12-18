package kafka

import (
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/lastexile/kepler"
)

// MarshallerFunc used to serialize Message to bytes
type MarshallerFunc func(m kepler.Message) ([]byte, []byte, error)

// UnmarshalFunction used to translate kafka Message to kepler
type UnmarshalFunction func(in *kafka.Message) (kepler.Message, error)

// default connection retry interval
const connectionRetryInterval = 10 * time.Second

func ensureDeliveryTimeouts(config *kafka.ConfigMap) *kafka.ConfigMap {
	config.SetKey("api.version.request", "false")
	config.SetKey("message.timeout.ms", 100)
	config.SetKey("session.timeout.ms", 100)
	config.SetKey("request.timeout.ms", 100)
	config.SetKey("message.send.max.retries", 3)

	return config
}
