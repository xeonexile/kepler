package kepler

import (
	"context"
)

// Spring of outgong messages
type Spring interface {
	Out(ctx context.Context, out chan Message) <-chan Message
	LinkTo(Sink, RouteCondition) (closer func())
}

func (s *springImpl) Out(ctx context.Context, o chan Message) <-chan Message {

	go func() {
		s.action(ctx, o)
	}()

	return o
}

func (s *springImpl) LinkTo(sink Sink, cond RouteCondition) (closer func()) {
	route := s.router.AddRoute(sink.Name(), cond)

	//pass linked conext to Sink
	inCtx, inClose := context.WithCancel(route.Ctx())
	sink.In(inCtx, s.Out(route.Ctx(), route.Buff()))

	//send signal to Sink and close original route
	return func() { inClose(); route.Close() }
}

// SpringFunction out generator function
type SpringFunction func(ctx context.Context, out chan<- Message)

// NewSpring creates new Spring
func NewSpring(name string, action SpringFunction) Spring {
	return &springImpl{name: name, action: action, router: NewRouter(false)}
}

type springImpl struct {
	name   string
	action SpringFunction
	router Router
}
