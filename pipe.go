package kepler

import (
	"reflect"
)

type Pipe interface {
	Spring
	Sink
}

type PipeImpl struct {
	name   string
	action PipeFunction
	routes *RoutesMap
}

type PipeFunction func(in Message) Message

func (p *PipeImpl) Out(route Route) <-chan Message {
	return route.Buff()
}

func (p *PipeImpl) In(input <-chan Message) {
	go func() {
		for msg := range input {
			reflect.Select(p.routes.Cases(p.action(msg)))
		}
	}()
}

func (p *PipeImpl) Name() string {
	return p.name
}

func (p *PipeImpl) LinkTo(sink Sink, cond RouteCondition) {
	route := p.addRoute(sink.Name(), cond)
	sink.In(p.Out(route))
}

func NewPipe(name string, action PipeFunction) Pipe {
	return &PipeImpl{name: name, action: action, routes: NewRoutesMap()}
}

func (p *PipeImpl) addRoute(name string, rc RouteCondition) (res *route) {
	res = NewRoute(name, rc)
	p.routes.Add(res)
	return res
}
