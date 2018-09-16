package kepler

import (
	"context"
	"log"
)

// Route encapsulates channel-based communication mechanism
type Route interface {
	Name() string
	Buff() chan Message
	Ctx() context.Context
	Cond() func(m Message) bool

	// Close current route, and remove it from parent router
	Close()
}

type route struct {
	name   string
	buff   chan Message
	ctx    context.Context
	cond   func(m Message) bool
	router Router
	cancel context.CancelFunc
}

func (r *route) Name() string {
	return r.name
}

func (r *route) Buff() chan Message {
	return r.buff
}

func (r *route) Ctx() context.Context {
	return r.ctx
}

func (r *route) Cond() func(m Message) bool {
	return r.cond
}

func (r *route) Close() {
	log.Println("Route Closing")
	r.cancel()

	r.router.RemoveRoute(r)
	//close(r.buff)
	log.Println("Route Closed")
}

// NewRoute creates new instance of named conditional Route
func newRoute(name string, rc func(m Message) bool, router Router) Route {
	ctx, cancel := context.WithCancel(context.Background())
	return &route{name, make(chan Message), ctx, rc, router, cancel}
}
