package kepler

import (
	"context"
	"log"
)

// Pipe is a composition of Spring - Map - Sink. Act as a simple transform block
type Pipe interface {
	Spring
	Sink
}

type pipeImpl struct {
	action    PipeFunction
	router    Router
	broadcast bool
}

// PipeFunction defines transform func that will be performed within block
type PipeFunction func(in Message) Message

// Out outgoing channel
func (p *pipeImpl) Out(ctx context.Context, buf chan Message) <-chan Message {
	return buf
}

// In incoming channel
func (p *pipeImpl) In(ctx context.Context, input <-chan Message) {
	go func() {
		for {
			select {
			case msg := <-input:
				if msg != nil {
					m := p.action(msg)
					if m != nil {
						p.router.Send(m)
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
func (p *pipeImpl) LinkTo(name string, sink Sink, cond RouteCondition) (closer func()) {
	route := p.router.AddRoute(name, cond)
	sink.In(route.Ctx(), p.Out(nil, route.Buff()))
	return func() { route.Close() }
}

// NewPipe creates new instance of pipe with defined transform action
func NewPipe(action PipeFunction) Pipe {
	return &pipeImpl{action: action, router: NewRouter(false)}
}

// NewBroadcastPipe creates broadcast pipe
func NewBroadcastPipe(action PipeFunction) Pipe {
	return &pipeImpl{action: action, router: NewRouter(true)}
}
