package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/lastexile/kepler"
	log "github.com/sirupsen/logrus"
)

// NewKafkaSink creates new Kafka outgoing sink
func NewSink(topic string, config *kafka.ConfigMap, formatter MarshallerFunc) (sink kepler.Sink, err error) {

	var p *kafka.Producer
	p, err = kafka.NewProducer(config)
	if err != nil {
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
					log.Info("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	sink = kepler.NewSink(topic, func(m kepler.Message) {

		value, err := formatter(m)
		if err != nil {
			log.Error("Unable to marshall message")
			return
		}
		p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          value,
		}, nil)

		log.Info("message sent to kafka")
	})

	return
}
