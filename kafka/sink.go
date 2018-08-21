package kafka

import (
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/lastexile/kepler"
	log "github.com/sirupsen/logrus"
)

// NewSink creates new Kafka outgoing sink
func NewSink(topic string, config *kafka.ConfigMap, formatter MarshallerFunc) (sink kepler.Sink, err error) {

	var p *kafka.Producer

	config.SetKey("api.version.request", "false")
	config.SetKey("message.timeout.ms", 100)
	config.SetKey("session.timeout.ms", 100)
	config.SetKey("request.timeout.ms", 100)
	config.SetKey("message.send.max.retries", 3)
	p, err = kafka.NewProducer(config)
	if err != nil {
		log.Errorf("Unable to create kafka Producer: %v\n", err)
		return
	}

	// TODO: consider closing of producer
	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Errorf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					log.Debugf("Delivered message to %v", ev.TopicPartition)
				}
			}
		}
	}()

	delivery := make(chan kafka.Event)
	sink = kepler.NewSink(topic, func(m kepler.Message) {
		value, err := formatter(m)
		if err != nil {
			log.Error("Unable to marshall message")
			return
		}
		msg := &kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          value,
		}
		for {
			if err := p.Produce(msg, delivery); err != nil {
				log.Errorf("Cannot enqueue message: %v\n", err)
				return
			}

			//TODO track status synchronously
			e := <-delivery

			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Errorf("Delivery failed: %v\n", ev.TopicPartition)
					time.Sleep(connectionRetryInterval)
					continue

				} else {
					log.Debugf("Delivered message to %v", ev.TopicPartition)
					return
				}
			}
		}

	})

	return
}
