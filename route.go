package kepler

import (
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

type RoutesMap struct {
	routes map[string][]Route
	cases  map[string][]reflect.SelectCase
}

func NewRoutesMap() *RoutesMap {
	return &RoutesMap{make(map[string][]Route), make(map[string][]reflect.SelectCase)}
}

// Add route to map by its name
func (r *RoutesMap) Add(route Route) {
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

func (r *RoutesMap) ByName(name string) []Route {
	if b, ok := r.routes[name]; ok {
		return b
	}
	return empty
}

// ByCond returns first random routes collection, that can accept Message m
func (r *RoutesMap) ByCond(m Message) []Route {
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
func (r *RoutesMap) Cases(m Message) []reflect.SelectCase {
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
