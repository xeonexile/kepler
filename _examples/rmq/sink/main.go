package main

import (
	"bufio"
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/lastexile/kepler"
	"github.com/lastexile/kepler/rmq"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	log.Println("starting...")

	url := "amqp://user:pass@localhost:5672"

	q := rmq.QueueOptions{"test_ex", "test", "test", true, true, true, true}
	s, err := rmq.NewSink(rmq.Connection(url), q, func(m kepler.Message) ([]byte, error) {
		return []byte(m.Value().(string)), nil
	})
	if err != nil {
		log.Fatalf("Unable to create rmqspring: %v\n", err)
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

	reader := bufio.NewReader(os.Stdin)
	log.Print("Enter text: ")
	text, _ := reader.ReadString('\n')
	log.Println(text)
}
