package kepler

import (
	"context"
	"log"
	"reflect"
	"sync"
)

const (
	// DefaultRouteGroup name
	DefaultRouteGroup string = "."
)

// Router responsible for message's routing among connected routes
//can work in two modes as FanOut and Broadcast
type Router interface {
	AddRoute(name string, rc func(m Message) bool) Route
	RemoveRoute(Route)
	Send(Message) int
	Close()
}

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

// RouteCondition filtrates outgoing messages
type RouteCondition func(m Message) bool

type router struct {
	mx          sync.Mutex
	routes      map[string][]Route
	ring        int
	broadcaster bool
}

// NewRouter creates new Router instance
func NewRouter(broadcaster bool) Router {
	return &router{sync.Mutex{}, make(map[string][]Route), 0, broadcaster}
}

// Add route to map by its name
func (r *router) AddRoute(name string, rc func(m Message) bool) Route {
	r.mx.Lock()
	defer r.mx.Unlock()

	route := newRoute(name, rc, r)
	var (
		b  []Route
		ok bool
	)
	if b, ok = r.routes[route.Name()]; !ok {
		b = make([]Route, 0)
	}
	r.routes[route.Name()] = append(b, route)

	return route
}

func (r *router) Close() {
	r.mx.Lock()
	defer r.mx.Unlock()

	for _, rg := range r.routes {
		for _, ri := range rg {
			ri.Close()
			//close channel since this method is called from Input
			close(ri.Buff())
		}
	}
}

func (r *router) RemoveRoute(route Route) {
	r.mx.Lock()
	defer r.mx.Unlock()

	if rg, ok := r.routes[route.Name()]; ok {
		var nrg []Route
		for _, v := range rg {
			if v == route {
				continue
			} else {
				nrg = append(nrg, v)
			}
		}
		r.routes[route.Name()] = nrg
	}
}

// Allways true condition
func Allways(_ Message) bool {
	return true
}

var (
	empty = make([]Route, 0)
)

func (r *router) broadcast(m Message) int {
	routes := r.byCond(m)

	for _, rt := range routes {
		rt.Buff() <- m
	}
	return 1
}

// Send message to first free Cond applicable route
func (r *router) Send(m Message) int {
	r.mx.Lock()
	defer r.mx.Unlock()

	if r.broadcaster {
		return r.broadcast(m)
	}

	c := -1
	routes := r.byCond(m)
	switch len(routes) {
	case 1:
		c = send1(m, routes[0].Buff())
	case 2:
		c = send2(m, routes[0].Buff(), routes[1].Buff())
	case 3:
		c = send3(m, routes[0].Buff(), routes[1].Buff(), routes[2].Buff())
	case 4:
		c = send4(m, routes[0].Buff(), routes[1].Buff(), routes[2].Buff(), routes[3].Buff())
	case 5:
		c = send5(m, routes[0].Buff(), routes[1].Buff(), routes[2].Buff(), routes[3].Buff(), routes[4].Buff())
	case 6:
		c = send6(m, routes[0].Buff(), routes[1].Buff(), routes[2].Buff(), routes[3].Buff(), routes[4].Buff(), routes[5].Buff())
	case 7:
		c = send7(m, routes[0].Buff(), routes[1].Buff(), routes[2].Buff(), routes[3].Buff(), routes[4].Buff(), routes[5].Buff(), routes[6].Buff())
	default:
		c = sendDefault(m, routes)
	}

	if c == -1 {
		if r.ring++; r.ring >= len(routes) {
			r.ring = 0
		}

		routes[r.ring].Buff() <- m
	}
	return r.ring
}

// ByName returns list of routes by name
func (r *router) byName(name string) []Route {
	if b, ok := r.routes[name]; ok {
		return b
	}
	return empty
}

// ByCond returns first random routes collection, that can accept Message m
func (r *router) byCond(m Message) []Route {
	if m == nil {
		return empty
	}
	res := empty
	for _, v := range r.routes {
		if len(v) > 0 && v[0].Cond()(m) {
			res = v
			if v[0].Name() == DefaultRouteGroup {
				continue
			}
			return v
		}
	}
	return res
}

func sendDefault(m Message, in []Route) int {
	cases := make([]reflect.SelectCase, len(in))

	for i, c := range in {
		cases[i] = reflect.SelectCase{
			Dir:  reflect.SelectSend,
			Chan: reflect.ValueOf(c.Buff()),
			Send: reflect.ValueOf(m)}
	}

	idx, _, _ := reflect.Select(cases)
	return idx
}

func send1(m Message, in ...chan Message) int {
	select {
	case in[0] <- m:
		return 0
	default:
		return -1
	}

}

func send2(m Message, in ...chan Message) int {
	select {
	case in[0] <- m:
		log.Print("2 done 0")
		return 0
	case in[1] <- m:
		log.Print("2 done 1")
		return 1
	default:
		//log.Print("2 wait")
		return -1
	}
}

func send3(m Message, in ...chan Message) int {
	select {
	case in[0] <- m:
		return 0
	case in[1] <- m:
		return 1
	case in[2] <- m:
		return 2
	default:
		return -1
	}
}
func send4(m Message, in ...chan Message) int {
	select {
	case in[0] <- m:
		return 0
	case in[1] <- m:
		return 1
	case in[2] <- m:
		return 2
	case in[3] <- m:
		return 3
	default:
		return -1
	}
}
func send5(m Message, in ...chan Message) int {
	select {
	case in[0] <- m:
		return 0
	case in[1] <- m:
		return 1
	case in[2] <- m:
		return 2
	case in[3] <- m:
		return 3
	case in[4] <- m:
		return 4
	default:
		return -1
	}
}

func send6(m Message, in ...chan Message) int {
	select {
	case in[0] <- m:
		return 0
	case in[1] <- m:
		return 1
	case in[2] <- m:
		return 2
	case in[3] <- m:
		return 3
	case in[4] <- m:
		return 4
	case in[5] <- m:
		return 5
	default:
		return -1
	}
}

func send7(m Message, in ...chan Message) int {
	select {
	case in[0] <- m:
		return 0
	case in[1] <- m:
		return 1
	case in[2] <- m:
		return 2
	case in[3] <- m:
		return 3
	case in[4] <- m:
		return 4
	case in[5] <- m:
		return 5
	case in[6] <- m:
		return 6
	default:
		return -1
	}
}
