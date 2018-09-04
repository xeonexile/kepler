package rmq

import (
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

// default connection retry interval
const connectionRetryInterval = 10 * time.Second

// ConnectionFactoryFunc returns opened connection
type ConnectionFactoryFunc func() (*amqp.Connection, error)

// QueueOptions defines basic creation parameters
type QueueOptions struct {
	Name             string
	ExchangeName     string
	ConsumerName     string
	Durable          bool
	ExchangeDurable  bool
	DeleteWhenUnused bool
	NoWait           bool
}

// Connection returns based on Default dialer connectionFactory
func Connection(dialString string) ConnectionFactoryFunc {
	return func() (*amqp.Connection, error) {
		return amqp.Dial(dialString)
	}
}

func initExchangeChannelQueue(connFactory ConnectionFactoryFunc, queue QueueOptions) (conn *amqp.Connection, ch *amqp.Channel, err error) {
	conn, err = connFactory()
	if err != nil {
		log.Errorf("Failed to open the connection: %v \n", err)
		return
	}

	log.Infoln("Connected")

	ch, err = conn.Channel()
	if err != nil {
		log.Errorf("Channel not opened: %v", err)
		return
	}
	log.Infoln("Channel opened")
	err = ch.ExchangeDeclare(
		queue.ExchangeName,     // name
		"topic",                // type
		queue.ExchangeDurable,  // durable
		queue.DeleteWhenUnused, // auto-deleted
		false,        // internal
		queue.NoWait, // noWait
		nil,          // arguments
	)

	if err != nil {
		log.Errorf("Failed to declare Exchange: %v\n", err)
		return
	}

	_, err = ch.QueueDeclare(
		queue.Name,             // name
		queue.Durable,          // durable
		queue.DeleteWhenUnused, // delete when unused
		false,        // exclusive
		queue.NoWait, // no-wait
		nil,          // arguments
	)

	if err != nil {
		log.Errorf("QueueDeclare failed Channel: %v", err)
		return
	}

	return
}

func resetConnection(conn *amqp.Connection, ch *amqp.Channel) (*amqp.Connection, *amqp.Channel) {
	if conn != nil {
		conn.Close()
	}

	if ch != nil {
		ch.Close()
	}

	return nil, nil
}
