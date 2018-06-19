package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/lastexile/kepler"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	log.Println("starting...")

	s := kepler.NewSpring("range", func(ctx context.Context, c chan<- kepler.Message) {

		for i := 0; i < 10; i++ {
			select {
			case <-ctx.Done():
				log.Println("Done")
				close(c)
				return
			case c <- kepler.NewValueMessage("range", i):
			}
		}

		time.Sleep(10 * time.Second)
		for i := 10; i < 20; i++ {
			select {
			case <-ctx.Done():
				log.Println("Done")
				close(c)
				return
			case c <- kepler.NewValueMessage("range2", i):
			}
		}
	})

	mux := kepler.NewBroadcastPipe("mux", func(m kepler.Message) kepler.Message {
		log.Println("mux: " + m.String())
		return m
	})

	// t1 := kepler.NewSink("t", func(m kepler.Message) {
	// 	log.Println("t1_1>: " + m.String())
	// 	time.Sleep(4 * time.Second)
	// 	log.Println("t1_1<: " + m.String())
	// })

	// t2 := kepler.NewSink("t", func(m kepler.Message) {
	// 	log.Println("t1_2>: " + m.String())
	// 	time.Sleep(8 * time.Second)
	// 	log.Println("t1_2<: " + m.String())
	// })

	var even [3]kepler.Sink
	for i, _ := range even {
		name := fmt.Sprintf("even_%v", i)
		even[i] = kepler.NewSink("even", func(m kepler.Message) {
			log.Println(name + "> " + m.String())
			time.Sleep(3 * time.Second)
			log.Println(name + "< " + m.String())
		})

		mux.LinkTo(even[i], func(m kepler.Message) bool { return m.Value().(int)%2 == 0 })
	}

	var odd [8]kepler.Sink
	for i, _ := range odd {
		name := fmt.Sprintf("odd_%v", i)
		odd[i] = kepler.NewSink("odd", func(m kepler.Message) {
			log.Println(name + "> " + m.String())
			time.Sleep(10 * time.Second)
			log.Println(name + "< " + m.String())
		})

		mux.LinkTo(odd[i], func(m kepler.Message) bool { return m.Value().(int)%2 != 0 })
	}

	closeMux := s.LinkTo(mux, kepler.Allways)

	time.Sleep(10 * time.Second)
	closeMux()
	//p.LinkTo(t2, func(m kepler.Message) bool { return m.Value().(int) > 5 })
	// mux.LinkTo(t3, func(m kepler.Message) bool { return m.Value().(int)%3 == 0 })
	// mux.LinkTo(t4, func(m kepler.Message) bool { return m.Value().(int)%3 == 0 })

	// mux.LinkTo(t2, func(m kepler.Message) bool { return m.Value().(int)%3 != 0 })
	// mux.LinkTo(t1, func(m kepler.Message) bool { return m.Value().(int)%3 != 0 })
	// mux.LinkTo(t2, kepler.Allways)
	// mux.LinkTo(t1, kepler.Allways)

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter text: ")
	text, _ := reader.ReadString('\n')
	fmt.Println(text)
}
