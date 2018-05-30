package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/lastexile/kepler"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	log.Println("starting...")

	s := kepler.NewSpring("range", func(c chan<- kepler.Message) {
		for i := 0; i < 10; i++ {
			c <- kepler.NewValueMessage("range", i)
		}

		time.Sleep(20 * time.Second)
		for i := 10; i < 20; i++ {
			c <- kepler.NewValueMessage("range2", i)
		}
	})

	mux := kepler.NewPipe("mux", func(m kepler.Message) kepler.Message {
		return m
	})

	t1 := kepler.NewSink("t", func(m kepler.Message) {
		log.Println("t1_1>: " + m.String())
		time.Sleep(4 * time.Second)
		log.Println("t1_1<: " + m.String())
	})

	t2 := kepler.NewSink("t", func(m kepler.Message) {
		log.Println("t1_2>: " + m.String())
		time.Sleep(8 * time.Second)
		log.Println("t1_2<: " + m.String())
	})

	t3 := kepler.NewSink("t3", func(m kepler.Message) {
		log.Println("t3_1>: " + m.String())
		time.Sleep(10 * time.Second)
		log.Println("t3_1<: " + m.String())
	})

	t4 := kepler.NewSink("t3", func(m kepler.Message) {
		log.Println("t3_2>: " + m.String())
		time.Sleep(10 * time.Second)
		log.Println("t3_2<: " + m.String())
	})

	s.LinkTo(mux, kepler.Allways)
	//p.LinkTo(t2, func(m kepler.Message) bool { return m.Value().(int) > 5 })
	mux.LinkTo(t3, func(m kepler.Message) bool { return m.Value().(int)%3 == 0 })
	mux.LinkTo(t4, func(m kepler.Message) bool { return m.Value().(int)%3 == 0 })

	mux.LinkTo(t2, func(m kepler.Message) bool { return m.Value().(int)%3 != 0 })
	mux.LinkTo(t1, func(m kepler.Message) bool { return m.Value().(int)%3 != 0 })
	// mux.LinkTo(t2, kepler.Allways)
	// mux.LinkTo(t1, kepler.Allways)

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter text: ")
	text, _ := reader.ReadString('\n')
	fmt.Println(text)
}
