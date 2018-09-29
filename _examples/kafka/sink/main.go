package main

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/lastexile/kepler"
	kkafka "github.com/lastexile/kepler/kafka"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	log.Println("starting...")

	cfg := &kafka.ConfigMap{
		"metadata.broker.list": "localhost:9092",
		"group.id":             "odd",
		"default.topic.config": kafka.ConfigMap{"auto.offset.reset": "earliest"},
	}

	s, err := kkafka.NewSink("test", cfg, func(m kepler.Message) ([]byte, []byte, error) {
		return []byte(m.Value().(string)), []byte(strconv.Itoa(1)), nil
	})

	if err != nil {
		log.Fatalf("Unable to create kafkasink: %v\n", err)
	}

	spring := kepler.NewSpring(func(ctx context.Context, ch chan<- kepler.Message) {

		i := 1
		for {
			ch <- kepler.NewMessage("odd", strconv.Itoa(i))
			i++
			time.Sleep(1 * time.Second)
			log.Println(i)
		}
	})

	spring.LinkTo(".", s, kepler.Allways)

	kepler.Await()
}
