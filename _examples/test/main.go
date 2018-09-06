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

	s := kepler.NewSpring(func(ctx context.Context, c chan<- kepler.Message) {

		for i := 0; i < 10; i++ {
			select {
			case <-ctx.Done():
				log.Println("Done")
				close(c)
				return
			case c <- kepler.NewMessage("range", []int{i, i + 1, i + 2}):
			}
		}

		time.Sleep(10 * time.Second)
		for i := 10; i < 20; i++ {
			select {
			case <-ctx.Done():
				log.Println("Done")
				close(c)
				return
			case c <- kepler.NewMessage("range2", []int{i, i + 1, i + 2, i + 3}):
			}
		}
	})

	roll := kepler.NewRoll(func(m kepler.Message) (res []kepler.Message) {
		for _, v := range m.Value().([]int) {
			res = append(res, kepler.NewMessage(m.Topic(), v))
		}
		return
	})

	s.LinkTo(".", roll, kepler.Allways)
	// mux := kepler.NewBroadcastPipe(func(m kepler.Message) kepler.Message {
	// 	log.Println("mux: " + m.String())
	// 	return m
	// })

	t1 := kepler.NewSink(func(m kepler.Message) {
		time.Sleep(4 * time.Second)
		log.Println("t1: " + m.String())
	})

	roll.LinkTo(".", t1, kepler.Allways)

	// t2 := kepler.NewSink("t", func(m kepler.Message) {
	// 	log.Println("t1_2>: " + m.String())
	// 	time.Sleep(8 * time.Second)
	// 	log.Println("t1_2<: " + m.String())
	// })

	//closeMux := s.LinkTo(".", mux, kepler.Allways)

	//time.Sleep(10 * time.Second)
	//closeMux()
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
