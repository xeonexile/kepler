package kepler

import (
	"context"
	"log"
)

// Sink is a base endpoint block that consumes incoming message
type Sink interface {
	In(ctx context.Context, input <-chan Message)
}

type sinkImpl struct {
	action SinkFunction
}

// SinkFunction defines incoming message processing action
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
				//TODO: propagate to attached
				return
			}
		}
	}()
}

// NewSink creates new instance of sink
func NewSink(action SinkFunction) Sink {
	return &sinkImpl{action: action}
}
