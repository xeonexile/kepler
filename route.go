package kepler

import (
	"log"
	"reflect"
)

type Router interface {
	AddRoute(name string, rc func(m Message) bool) chan Message
	Send(Message) int
}

type Route interface {
	Name() string
	Buff() chan Message
	Cond() func(m Message) bool
}

type route struct {
	name string
	buff chan Message
	cond func(m Message) bool
}

func (r *route) Name() string {
	return r.name
}

func (r *route) Buff() chan Message {
	return r.buff
}

func (r *route) Cond() func(m Message) bool {
	return r.cond
}

// NewRoute creates new instance of named conditional Route
func NewRoute(name string, rc func(m Message) bool) Route {
	return &route{name, make(chan Message), rc}
}

// RouteCondition filtrates outgoing messages
type RouteCondition func(m Message) bool

type router struct {
	routes map[string][]Route
	ring   int
}

// NewRouter creates new Router instance
func NewRouter() Router {
	return &router{make(map[string][]Route), 0}
}

// Add route to map by its name
func (r *router) AddRoute(name string, rc func(m Message) bool) chan Message {
	route := NewRoute(name, rc)
	var (
		b  []Route
		ok bool
	)
	if b, ok = r.routes[route.Name()]; !ok {
		b = make([]Route, 0)
	}
	r.routes[route.Name()] = append(b, route)

	return route.Buff()
}

// Allways true condition
func Allways(_ Message) bool {
	return true
}

var (
	empty = make([]Route, 0)
)

// Send message to first free Cond aplicable route
func (r *router) Send(m Message) int {
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
	for _, v := range r.routes {
		if len(v) > 0 && v[0].Cond()(m) {
			return v
		}
	}
	return empty
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
