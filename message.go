package kepler

import (
	"fmt"
	"reflect"
	"time"
)

type Message interface {
	Topic() string
	Value() interface{}

	SetValue(value interface{}) Message
	String() string
}

type ValueMessage struct {
	topic     string
	createdAt time.Time
	value     interface{}
}

func (m *ValueMessage) Value() interface{} {
	return m.value
}

func (m *ValueMessage) SetValue(value interface{}) Message {

	if reflect.TypeOf(m.value) != reflect.TypeOf(value) {
		panic("types must be the same")
	}

	return NewValueMessage(m.topic, value)
}

func (m *ValueMessage) Topic() string {
	return m.topic
}

func (m *ValueMessage) String() string {
	return fmt.Sprintf("[%s]@%s: %v", m.topic, m.createdAt.String(), m.value)
}

func NewValueMessage(topic string, value interface{}) Message {
	return &ValueMessage{topic: topic, createdAt: time.Now().UTC(), value: value}
}
