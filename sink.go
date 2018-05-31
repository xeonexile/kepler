package kepler

type Sink interface {
	Name() string
	In(input <-chan Message)
}

type sinkImpl struct {
	name   string
	action SinkFunction
}

type SinkFunction func(msg Message)

func (s *sinkImpl) In(input <-chan Message) {
	go func() {
		for msg := range input {
			s.action(msg)
		}
	}()
}

func (s *sinkImpl) Name() string {
	return s.name
}

func NewSink(name string, action SinkFunction) Sink {
	return &sinkImpl{name: name, action: action}
}
