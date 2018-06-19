package kepler

import (
	"context"
	"log"
)

type Sink interface {
	Name() string
	In(ctx context.Context, input <-chan Message)
}

type sinkImpl struct {
	name   string
	action SinkFunction
}

type SinkFunction func(msg Message)

func (s *sinkImpl) In(ctx context.Context, input <-chan Message) {
	go func() {
		for {
			select {
			case msg := <-input:
				if msg != nil {
					s.action(msg)
				}
			case <-ctx.Done():
				log.Println("Sink Done")
				//TODO: propogate to attached
				return
			}
		}
	}()
}

func (s *sinkImpl) Name() string {
	return s.name
}

func NewSink(name string, action SinkFunction) Sink {
	return &sinkImpl{name: name, action: action}
}
