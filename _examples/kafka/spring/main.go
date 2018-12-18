package main

import (
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/lastexile/kepler"
	kkafka "github.com/lastexile/kepler/kafka"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	log.Println("starting...")

	cfg := &kafka.ConfigMap{
		"metadata.broker.list":            "localhost:9092",
		"group.id":                        "odd",
		"go.events.channel.enable":        true,
		"go.application.rebalance.enable": true,
		"default.topic.config":            kafka.ConfigMap{"auto.offset.reset": "earliest"},
	}

	s, err := kkafka.NewSpring("test", 0, 0, cfg, func(m *kafka.Message) (kepler.Message, error) {
		return kepler.NewMessage("foo", string(m.Value)), nil
	})

	if err != nil {
		log.Fatalf("Unable to create kafkaspring: %v\n", err)
	}

	logSink := kepler.NewSink(func(m kepler.Message) {

		log.Println(m.String())
	})

	s.LinkTo(".", logSink, kepler.Allways)

	kepler.Await()
}
