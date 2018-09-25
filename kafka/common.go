package kafka

import (
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/lastexile/kepler"
	"github.com/streadway/amqp"
)

// MarshallerFunc used to serialize Message to bytes
type MarshallerFunc func(m kepler.Message) ([]byte, []byte, error)

// default connection retry interval
const connectionRetryInterval = 10 * time.Second

// ConnectionFactoryFunc returns opened connection
type ConnectionFactoryFunc func() (*amqp.Connection, error)

func ensureDeliveryTimeouts(config *kafka.ConfigMap) *kafka.ConfigMap {
	config.SetKey("api.version.request", "false")
	config.SetKey("message.timeout.ms", 100)
	config.SetKey("session.timeout.ms", 100)
	config.SetKey("request.timeout.ms", 100)
	config.SetKey("message.send.max.retries", 3)

	return config
}
