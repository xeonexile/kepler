package kepler

import (
	"log"
	"reflect"
)

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

func NewRoute(name string, rc func(m Message) bool) *route {
	return &route{name, make(chan Message), rc}
}

// RouteCondition filtrates outgoing messages
type RouteCondition func(m Message) bool

type Router struct {
	routes map[string][]Route
	cases  map[string][]reflect.SelectCase
}

func NewRouter() *Router {
	return &Router{make(map[string][]Route), make(map[string][]reflect.SelectCase)}
}

// Add route to map by its name
func (r *Router) Add(route Route) {
	var (
		b  []Route
		c  []reflect.SelectCase
		ok bool
	)
	if b, ok = r.routes[route.Name()]; !ok {
		b = make([]Route, 0)
		c = make([]reflect.SelectCase, 0)
	}
	c = r.cases[route.Name()]
	r.routes[route.Name()] = append(b, route)

	sc := reflect.SelectCase{
		Dir:  reflect.SelectSend,
		Chan: reflect.ValueOf(route.Buff())}
	r.cases[route.Name()] = append(c, sc)
}

func (r *Router) ByName(name string) []Route {
	if b, ok := r.routes[name]; ok {
		return b
	}
	return empty
}

// ByCond returns first random routes collection, that can accept Message m
func (r *Router) ByCond(m Message) []Route {
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

// Cases returns first random select cases collection, that can accept Message m
func (r *Router) Cases(m Message) []reflect.SelectCase {
	routes := r.ByCond(m)
	if (len(routes)) == 0 {
		return emptyCases
	}

	cases := r.cases[routes[0].Name()]
	for i := 0; i < len(cases); i++ {
		cases[i].Send = reflect.ValueOf(m)
	}
	return cases
}

// Allways true condition
func Allways(_ Message) bool {
	return true
}

var (
	empty      = make([]Route, 0)
	emptyCases = make([]reflect.SelectCase, 0)
)

func (r *Router) Send(m Message, routes []Route) int {
	switch len(routes) {
	case 1:
		return send1(m, routes[0].Buff())
	case 2:
		return send2(m, routes[0].Buff(), routes[1].Buff())
	case 3:
		return send3(m, routes[0].Buff(), routes[1].Buff(), routes[2].Buff())
	case 4:
		return send4(m, routes[0].Buff(), routes[1].Buff(), routes[2].Buff(), routes[3].Buff())
	case 5:
		return send5(m, routes[0].Buff(), routes[1].Buff(), routes[2].Buff(), routes[3].Buff(), routes[4].Buff())
	default:
		return -2
	}
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
