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
	})

	p := kepler.NewPipe("mux", func(m kepler.Message) kepler.Message {
		return m
	})

	t1 := kepler.NewSink("t", func(m kepler.Message) {
		log.Println("t1: " + m.String())
		time.Sleep(4 * time.Second)
	})

	t2 := kepler.NewSink("t", func(m kepler.Message) {
		log.Println("t2: " + m.String())
		time.Sleep(8 * time.Second)
	})

	s.LinkTo(p, kepler.Allways)
	//p.LinkTo(t2, func(m kepler.Message) bool { return m.Value().(int) > 5 })
	p.LinkTo(t2, kepler.Allways)
	p.LinkTo(t1, kepler.Allways)

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter text: ")
	text, _ := reader.ReadString('\n')
	fmt.Println(text)
}
