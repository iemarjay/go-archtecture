package events

type Listener interface {
	Handle(interface{})
}

type Name string

type Event struct {
	listeners map[Name][]Listener
}

func NewEvent() *Event {
	return &Event{listeners: map[Name][]Listener{}}
}

func (e *Event) Listen(name Name, listener Listener) {
	l := e.listeners[name]
	e.listeners[name] = append(l, listener)
}

func (e *Event) Emit(name Name, payload interface{}) {
	for _, listener := range e.listeners[name] {
		listener.Handle(payload)
	}
}
