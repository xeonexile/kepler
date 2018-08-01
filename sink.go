package kepler

import (
	"context"
	"log"
)

// Sink is a base endpoint block that consumes incomming message
type Sink interface {
	Name() string
	In(ctx context.Context, input <-chan Message)
}

type sinkImpl struct {
	name   string
	action SinkFunction
}

// SinkFunction defines incomming message processing action
type SinkFunction func(msg Message)

// MarshallerFunc used to serialize Message to bytes
type MarshallerFunc func(m Message) ([]byte, error)

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

// NewSink creates new instance of sink
func NewSink(name string, action SinkFunction) Sink {
	return &sinkImpl{name: name, action: action}
}
