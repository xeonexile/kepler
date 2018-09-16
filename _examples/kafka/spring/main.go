package main

import (
	"bufio"
	"log"
	"os"

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

	s, err := kkafka.NewSpring("test", 0, 0, cfg, func(d []byte) (kepler.Message, error) {
		return kepler.NewMessage("foo", string(d)), nil
	})

	if err != nil {
		log.Fatalf("Unable to create kafkaspring: %v\n", err)
	}

	logSink := kepler.NewSink("odd", func(m kepler.Message) {

		log.Println(m.String())
	})

	s.LinkTo(logSink, kepler.Allways)

	reader := bufio.NewReader(os.Stdin)
	log.Print("Enter text: ")
	text, _ := reader.ReadString('\n')
	log.Println(text)
}