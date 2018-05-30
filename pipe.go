package kepler

type Pipe interface {
	Spring
	Sink
}

type PipeImpl struct {
	name   string
	action PipeFunction
	router Router
}

type PipeFunction func(in Message) Message

func (p *PipeImpl) Out(buf chan Message) <-chan Message {
	return buf
}

func (p *PipeImpl) In(input <-chan Message) {
	go func() {
		for msg := range input {
			m := p.action(msg)
			if m != nil {
				p.router.Send(m)
			}
		}
	}()
}

func (p *PipeImpl) Name() string {
	return p.name
}

func (p *PipeImpl) LinkTo(sink Sink, cond RouteCondition) {
	sink.In(p.Out(p.router.AddRoute(sink.Name(), cond)))
}

func NewPipe(name string, action PipeFunction) Pipe {
	return &PipeImpl{name: name, action: action, router: NewRouter()}
}
