package kepler

import (
	"testing"
)

func TestRouteClosesForTwoSprings(t *testing.T) {

	b1 := NewBroadcastPipe(func(in Message) Message {
		return in
	}).(*pipeImpl)

	b2 := NewBroadcastPipe(func(in Message) Message {
		return in
	}).(*pipeImpl)

	sink := NewSink(func(in Message) {

	})

	closer1 := b1.LinkTo(".", sink, Allways)
	closer2 := b2.LinkTo(".", sink, Allways)

	closer1()
	closer2()
}
