package main

import (
	"context"
	"log"
	"time"

	"github.com/lastexile/kepler"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	produce := kepler.NewSpring(func(ctx context.Context, c chan<- kepler.Message) {

		for i := 0; i < 100; i = i + 3 {
			select {
			case <-ctx.Done():
				log.Println("Done")
				close(c)
				return
			case c <- kepler.NewMessage("produce", []int{i, i + 1, i + 2}):
			}
			time.Sleep(1 * time.Second)
		}
	})

	flatten := kepler.NewRoll(func(m kepler.Message) (res []kepler.Message) {
		for _, v := range m.Value().([]int) {
			res = append(res, kepler.NewMessage(m.Topic(), v))
		}
		return
	})

	broadcast := kepler.NewBroadcastPipe(func(m kepler.Message) kepler.Message {
		return m
	})

	batch1 := kepler.NewBatch(func(a kepler.Message, m kepler.Message) kepler.Message {
		return kepler.NewMessage("batch1", append(a.Value().([]int), m.Value().(int)))
	}, []int{}, 5*time.Second)

	batch2 := kepler.NewBatch(func(a kepler.Message, m kepler.Message) kepler.Message {
		return kepler.NewMessage("batch2", append(a.Value().([]int), m.Value().(int)))
	}, []int{}, 10*time.Second)

	end := kepler.NewSink(func(m kepler.Message) {
		time.Sleep(4 * time.Second)
		log.Println("log1: " + m.String())
	})

	produce.LinkTo(".", flatten, kepler.Allways)
	flatten.LinkTo(".", broadcast, kepler.Allways)
	broadcast.LinkTo(".", batch1, kepler.Allways)
	broadcast.LinkTo(".", batch2, kepler.Allways)
	batch1.LinkTo(".", end, kepler.Allways)
	batch2.LinkTo(".", end, kepler.Allways)

	kepler.Await()
}
