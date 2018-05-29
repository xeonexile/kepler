package kepler

type Sink interface {
	Name() string
	In(input <-chan Message)
}

type SinkImpl struct {
	name   string
	action SinkFunction
}

type SinkFunction func(msg Message)

func (s *SinkImpl) In(input <-chan Message) {
	go func() {
		for msg := range input {
			s.action(msg)
		}
	}()
}

func (s *SinkImpl) Name() string {
	return s.name
}

func NewSink(name string, action SinkFunction) Sink {
	return &SinkImpl{name: name, action: action}
}
