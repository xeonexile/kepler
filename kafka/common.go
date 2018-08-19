package kafka

import (
	"time"

	"github.com/lastexile/kepler"
	"github.com/streadway/amqp"
)

// MarshallerFunc used to serialize Message to bytes
type MarshallerFunc func(m kepler.Message) ([]byte, error)

// default connection retry interval
const connectionRetryInterval = 10 * time.Second

// ConnectionFactoryFunc returns opened connection
type ConnectionFactoryFunc func() (*amqp.Connection, error)
