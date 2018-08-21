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
	p, err = kafka.NewProducer(config)
	if err != nil {
		log.Errorf("Unable to create kafka Producer: %v\n", err)
		return
	}

	delivery := make(chan kafka.Event)

	f := func(m kepler.Message) {
		writePipe(m, topic, p, formatter, delivery)
	}
	sink = kepler.NewSink(topic, f)

	return
}

// NewPipe creates transit kafka write pipe
func NewPipe(topic string, config *kafka.ConfigMap, formatter MarshallerFunc) (pipe kepler.Pipe, err error) {

	var p *kafka.Producer
	p, err = kafka.NewProducer(ensureDeliveryTimeouts(config))
	if err != nil {
		log.Errorf("Unable to create kafka Producer: %v\n", err)
		return
	}

	delivery := make(chan kafka.Event)

	f := func(m kepler.Message) kepler.Message {
		return writePipe(m, topic, p, formatter, delivery)
	}
	pipe = kepler.NewPipe(topic, f)

	return
}

func writePipe(m kepler.Message, topic string, p *kafka.Producer, formatter MarshallerFunc, delivery chan kafka.Event) kepler.Message {
	value, err := formatter(m)
	if err != nil {
		log.Error("Unable to marshall message")
		return nil
	}
	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          value,
	}
	for {
		if err := p.Produce(msg, delivery); err != nil {
			log.Errorf("Cannot enqueue message: %v\n", err)
			return nil
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
				return m
			}
		}
	}
}
