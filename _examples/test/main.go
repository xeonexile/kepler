package main

import (
	"context"
	"log"
	"time"

	"github.com/lastexile/kepler"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	log.Println("starting...")

	s := kepler.NewSpring(func(ctx context.Context, c chan<- kepler.Message) {

		for i := 0; i < 100; i = i + 3 {
			select {
			case <-ctx.Done():
				log.Println("Done")
				close(c)
				return
			case c <- kepler.NewMessage("range", []int{i, i + 1, i + 2}):
			}
			time.Sleep(1 * time.Second)
		}
	})

	roll := kepler.NewRoll(func(m kepler.Message) (res []kepler.Message) {
		for _, v := range m.Value().([]int) {
			res = append(res, kepler.NewMessage(m.Topic(), v))
		}
		return
	})

	batch := kepler.NewBatch(func(a kepler.Message, m kepler.Message) kepler.Message {
		return kepler.NewMessage(a.Topic(), append(a.Value().([]int), m.Value().(int)))
	}, []int{}, 5*time.Second)

	s.LinkTo(".", roll, kepler.Allways)
	roll.LinkTo(".", batch, kepler.Allways)
	// mux := kepler.NewBroadcastPipe(func(m kepler.Message) kepler.Message {
	// 	log.Println("mux: " + m.String())
	// 	return m
	// })

	log1 := kepler.NewSink(func(m kepler.Message) {
		time.Sleep(4 * time.Second)
		log.Println("log1: " + m.String())
	})

	batch.LinkTo(".", log1, kepler.Allways)

	kepler.Await()
}
