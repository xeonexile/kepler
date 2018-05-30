package kepler

import (
	"sync"
)

// Spring of outgong messages
type Spring interface {
	Out(out chan Message) <-chan Message
	LinkTo(Sink, RouteCondition)
}

func (s *springImpl) Out(o chan Message) <-chan Message {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		s.action(o)
		wg.Done()
	}()

	go func() {
		wg.Wait()
		close(o)
	}()

	return o
}

func (s *springImpl) LinkTo(sink Sink, cond RouteCondition) {
	sink.In(s.Out(s.routes.AddRoute(sink.Name(), cond)))
}

// SpringFunction out generator function
type SpringFunction func(out chan<- Message)

// NewSpring creates new Spring
func NewSpring(name string, action SpringFunction) Spring {
	return &springImpl{name: name, action: action, routes: NewRouter()}
}

type springImpl struct {
	name   string
	action SpringFunction
	routes Router
}
