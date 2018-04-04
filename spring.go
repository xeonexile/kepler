package kepler

import (
	"sync"
)

// Spring of outgong messages
type Spring interface {
	Out(Route) <-chan Message
	LinkTo(Sink, RouteCondition)
}

func (s *springImpl) Out(route Route) <-chan Message {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		s.action(route.Buff())
		wg.Done()
	}()

	go func() {
		wg.Wait()
		close(route.Buff())
	}()

	return route.Buff()
}

func (s *springImpl) LinkTo(sink Sink, cond RouteCondition) {
	route := s.addRoute(sink.Name(), cond)
	sink.In(s.Out(route))
}

// SpringFunction out generator function
type SpringFunction func(out chan<- Message)

// NewSpring creates new Spring
func NewSpring(name string, action SpringFunction) Spring {
	return &springImpl{name: name, action: action, routes: NewRoutesMap()}
}

type springImpl struct {
	name   string
	action SpringFunction
	routes *RoutesMap
}

func (s *springImpl) addRoute(name string, rc RouteCondition) (res *route) {
	res = NewRoute(name, rc)
	s.routes.Add(res)
	return res
}
