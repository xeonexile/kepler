package kepler

type Pipe interface {
	Spring
	Sink
}

type pipeImpl struct {
	name   string
	action PipeFunction
	router Router
}

// PipeFunction defines transform func that will be performed within block
type PipeFunction func(in Message) Message

// Out outgoing channel
func (p *pipeImpl) Out(buf chan Message) <-chan Message {
	return buf
}

// In incomming channel
func (p *pipeImpl) In(input <-chan Message) {
	go func() {
		for msg := range input {
			m := p.action(msg)
			if m != nil {
				p.router.Send(m)
			}
		}
	}()
}

// Name of this pipe
func (p *pipeImpl) Name() string {
	return p.name
}

// LinkTo add new conditional link
func (p *pipeImpl) LinkTo(sink Sink, cond RouteCondition) {
	sink.In(p.Out(p.router.AddRoute(sink.Name(), cond)))
}

// NewPipe creates new instance of pipe with defined transform action
func NewPipe(name string, action PipeFunction) Pipe {
	return &pipeImpl{name: name, action: action, router: NewRouter()}
}
