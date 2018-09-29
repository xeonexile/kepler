package kepler

import (
	"context"
	"log"
	"time"
)

// BatchFunc produces accumulated output
type BatchFunc func(acc Message, in Message) Message

// NewBatch creates new due Time window batcher
func NewBatch(action BatchFunc, seed interface{}, dueTime time.Duration) Pipe {
	return &batch{action: action, seed: seed, dueTime: dueTime, router: NewRouter(false)}
}

type batch struct {
	action  BatchFunc
	seed    interface{}
	dueTime time.Duration
	router  Router
}

// Out outgoing channel
func (p *batch) Out(ctx context.Context, buf chan Message) <-chan Message {
	return buf
}

// In incoming channel
func (p *batch) In(ctx context.Context, input <-chan Message) {
	go func() {
		timer := time.NewTimer(0)

		var dueTimeCh <-chan time.Time
		var acc Message
		for {
			select {
			case msg := <-input:
				if msg != nil {
					if acc == nil {
						acc = NewMessage("seed", p.seed)
					}
					acc = p.action(acc, msg)
					if dueTimeCh == nil {
						timer.Reset(p.dueTime)
						dueTimeCh = timer.C

					}
				}
			case <-dueTimeCh:
				{
					if acc != nil {
						p.router.Send(acc)
						dueTimeCh = nil
						acc = nil
					}
				}
			case <-ctx.Done():
				log.Println("In Done")
				p.router.Close()
				return
			}
		}
	}()
}

// LinkTo add new conditional link
func (p *batch) LinkTo(name string, sink Sink, cond RouteCondition) (closer func()) {
	route := p.router.AddRoute(name, cond)

	sink.In(route.Ctx(), p.Out(nil, route.Buff()))

	return func() { route.Close() }
}
