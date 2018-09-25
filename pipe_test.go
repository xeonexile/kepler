package kepler

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouteClosesForTwoSprings(t *testing.T) {

	res := 0
	var wg sync.WaitGroup
	wg.Add(1)
	s := NewSpring(func(ctx context.Context, in chan<- Message) {
		in <- NewMessage("a", 1)
	})

	b1 := NewBroadcastPipe(func(in Message) Message {
		return in.SetValue(in.Value().(int) + 1)
	}).(*pipeImpl)

	b2 := NewPipe(func(in Message) Message {
		return in.SetValue(in.Value().(int) + 1)
	}).(*pipeImpl)

	sink := NewSink(func(in Message) {
		res = in.Value().(int)
		wg.Done()
	})

	s.LinkTo(".", b1, Allways)
	closer1 := b1.LinkTo(".", b2, Allways)
	closer2 := b2.LinkTo(".", sink, Allways)

	wg.Wait()
	assert.Equal(t, 3, res)
	closer1()
	closer2()
}

func TestRoutes(t *testing.T) {

	res := 0
	var wg sync.WaitGroup
	wg.Add(1)
	s1 := NewSpring(func(ctx context.Context, in chan<- Message) {
		in <- NewMessage("a", 1)
	})

	s := NewPipe(func(in Message) Message {
		return in.SetValue(in.Value().(int) + 1)
	}).(*pipeImpl)

	p1 := NewPipe(func(in Message) Message {
		return in.SetValue(in.Value().(int) + 1)
	}).(*pipeImpl)
	p2 := NewPipe(func(in Message) Message {
		return in.SetValue(in.Value().(int) + 1)
	}).(*pipeImpl)
	p3 := NewPipe(func(in Message) Message {
		return in.SetValue(in.Value().(int) + 1)
	}).(*pipeImpl)
	p4 := NewPipe(func(in Message) Message {
		return in.SetValue(in.Value().(int) + 1)
	}).(*pipeImpl)
	p5 := NewPipe(func(in Message) Message {
		return in.SetValue(in.Value().(int) + 1)
	}).(*pipeImpl)
	p6 := NewPipe(func(in Message) Message {
		return in.SetValue(in.Value().(int) + 1)
	}).(*pipeImpl)

	sink := NewSink(func(in Message) {
		res = in.Value().(int)
		wg.Done()
	})

	s1.LinkTo(".", s, Allways)

	s.LinkTo(".", p1, Allways)
	p1.LinkTo(".", sink, Allways)

	s.LinkTo(".", p2, Allways)
	p2.LinkTo(".", sink, Allways)

	s.LinkTo(".", p3, Allways)
	p3.LinkTo(".", sink, Allways)

	s.LinkTo(".", p4, Allways)
	p4.LinkTo(".", sink, Allways)

	s.LinkTo(".", p5, Allways)
	p5.LinkTo(".", sink, Allways)

	s.LinkTo(".", p6, Allways)
	p6.LinkTo(".", sink, Allways)

	wg.Wait()
	assert.Equal(t, 3, res)
}

func TestMessage(t *testing.T) {
	a := NewMessage("a", 1)

	assert.Equal(t, "a", a.Topic())
	assert.Equal(t, 1, a.Value().(int))

	b := a.SetValue(2)
	assert.Equal(t, "a", a.Topic())
	assert.Equal(t, "a", b.Topic())
	assert.Equal(t, 1, a.Value().(int))
	assert.Equal(t, 2, b.Value().(int))
}
