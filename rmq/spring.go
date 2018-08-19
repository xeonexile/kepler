package rmq

import (
	"context"
	"time"

	"github.com/lastexile/kepler"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

// NewSpring - return generic amqp consumer spring
func NewSpring(connFactory ConnectionFactoryFunc, queue QueueOptions, formatter kepler.UnmarshalFunction) (kepler.Spring, error) {

	return kepler.NewSpring(queue.Name, func(ctx context.Context, out chan<- kepler.Message) {
		var conn *amqp.Connection
		var ch *amqp.Channel
		var err error
		var msgs <-chan amqp.Delivery
		var d amqp.Delivery
		var msg kepler.Message
		for {
			if conn == nil || ch == nil {
				log.Warn("Reconnecting")
				for {
					conn, ch, err = initExchangeChannelQueue(connFactory, queue)

					if err != nil {
						log.Errorf("Failed to init channel: %v\n", err)
						conn, ch = resetConnection(conn, ch)

						time.Sleep(connectionRetryInterval)
						continue
					}

					msgs, err = ch.Consume(
						queue.Name,         // queue
						queue.ConsumerName, // consumer
						false,              // auto-ack
						false,              // exclusive
						false,              // no-local
						queue.NoWait,       // no-wait
						nil,                // args
					)

					if err != nil {
						log.Errorf("Consume failed Channel: %v", err)
						conn, ch = resetConnection(conn, ch)
						continue
					}
					log.Info("Connected")
					break
				}
			}
			d = <-msgs

			if d.Body == nil {
				log.Errorf("Consume failed Channel: %v", err)
				conn, ch = resetConnection(conn, ch)
				continue
			}

			log.Infof("Received message: %s", d.Body)
			if msg, err = formatter(d.Body); err != nil {
				log.Errorf("Formatter error: %s", err)
			} else {
				out <- msg
				d.Ack(true)
			}
			log.Infof("mssss: %v", msg)
		}
	}), nil
}
