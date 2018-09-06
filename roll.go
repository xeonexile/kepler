package kepler

import (
	"context"
	"log"
)

type rollImpl struct {
	action    RollFunction
	router    Router
	broadcast bool
}

// RollFunction defines transform func that flattens input message
type RollFunction func(in Message) []Message

// Out outgoing channel
func (p *rollImpl) Out(ctx context.Context, buf chan Message) <-chan Message {
	return buf
}

// In incoming channel
func (p *rollImpl) In(ctx context.Context, input <-chan Message) {
	go func() {
		for {
			select {
			case msg := <-input:
				if msg != nil {
					for _, m := range p.action(msg) {
						if m != nil {
							p.router.Send(m)
						}
					}
				}
			case <-ctx.Done():
				log.Println("In Done")
				//TODO: propagate to attached
				p.router.Close()
				return
			}
		}
	}()
}

// LinkTo add new conditional link
func (p *rollImpl) LinkTo(name string, sink Sink, cond RouteCondition) (closer func()) {
	route := p.router.AddRoute(name, cond)

	sink.In(route.Ctx(), p.Out(nil, route.Buff()))

	return func() { route.Close() }
}

// NewRoll creates new instance of roll pipe that flattens incoming message
func NewRoll(action RollFunction) Pipe {
	return &rollImpl{action: action, router: NewRouter(false)}
}
