package kepler

import (
	"fmt"
	"reflect"
	"time"
)

// Message is a base communication structure
type Message interface {
	Topic() string
	Value() interface{}

	SetValue(value interface{}) Message
	String() string
}

type valueMessage struct {
	topic     string
	createdAt time.Time
	value     interface{}
}

func (m *valueMessage) Value() interface{} {
	return m.value
}

func (m *valueMessage) SetValue(value interface{}) Message {

	if reflect.TypeOf(m.value) != reflect.TypeOf(value) {
		panic("types must be the same")
	}

	return NewMessage(m.topic, value)
}

func (m *valueMessage) Topic() string {
	return m.topic
}

func (m *valueMessage) String() string {
	return fmt.Sprintf("[%s]@%s: %v", m.topic, m.createdAt.String(), m.value)
}

// NewMessage creates message with specified name and value
func NewMessage(topic string, value interface{}) Message {
	return &valueMessage{topic: topic, createdAt: time.Now().UTC(), value: value}
}
