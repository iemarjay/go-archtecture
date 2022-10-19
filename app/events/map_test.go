package events

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestEvent_Emit(t *testing.T) {
	event := NewEvent()
	listener := &testListener{}
	event.Listen(name, listener)
	event.Emit(name, testEventInstance)

	val1 := listener.Payload.(*testEvent)
	assert.Equal(t, val1, testEventInstance)
}

var name = Name("test")

type testListener struct {
	Payload interface{}
}

func (t *testListener) Handle(event interface{}) {
	t.Payload = event
}

type testEvent struct {
	Data bool
}

var testEventInstance = &testEvent{
	Data: true,
}
