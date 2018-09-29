package main

import (
	"log"

	"github.com/lastexile/kepler"
	"github.com/lastexile/kepler/rmq"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	log.Println("starting...")

	url := "amqp://user:pass@localhost:5672"

	q := rmq.QueueOptions{"test_ex", "test", "test", true, true, true, true}
	s, err := rmq.NewSpring(rmq.Connection(url), q, func(d []byte) (kepler.Message, error) {
		return kepler.NewMessage("foo", string(d)), nil
	})
	if err != nil {
		log.Fatalf("Unable to create rmqspring: %v\n", err)
	}

	logSink := kepler.NewSink(func(m kepler.Message) {

		log.Println(m.String())
	})

	s.LinkTo(".", logSink, kepler.Allways)

	kepler.Await()
}
