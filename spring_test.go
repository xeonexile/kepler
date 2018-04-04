package kepler

import (
	"fmt"
	"testing"
)

func TestSpringFanout(t *testing.T) {
	s := NewSpring("range", func(c chan<- Message) {
		for i := 0; i < 10; i++ {
			c <- NewValueMessage("range", i)
		}
	})

	t1 := NewSink("t1", func(m Message) {
		fmt.Println("t1: " + m.String())
	})

	t2 := NewSink("t2", func(m Message) {
		fmt.Println("t2: " + m.String())
	})

	s.LinkTo(t2, Allways)
	s.LinkTo(t1, Allways)
}