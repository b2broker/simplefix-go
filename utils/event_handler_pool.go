package utils

import "sync"

type Event int

const (
	EventDisconnect Event = iota
	EventConnect
	EventStopped
	EventRun

	EventLogon
	EventRequest
	EventLogout
)

type EventHandlerFunc func() bool

type EventHandlerPool struct {
	mu   sync.RWMutex
	pool map[Event][]EventHandlerFunc
}

func NewEventHandlerPool() *EventHandlerPool {
	return &EventHandlerPool{pool: make(map[Event][]EventHandlerFunc)}
}

func (evp *EventHandlerPool) Handle(e Event, handle EventHandlerFunc) {
	evp.mu.Lock()
	defer evp.mu.Unlock()

	if _, ok := evp.pool[e]; !ok {
		evp.pool[e] = []EventHandlerFunc{}
	}

	evp.pool[e] = append(evp.pool[e], handle)
}

func (evp *EventHandlerPool) Trigger(e Event) {
	evp.mu.RLock()
	defer evp.mu.RUnlock()

	handlers, ok := evp.pool[e]
	if !ok {
		return
	}

	for _, handle := range handlers {
		if !handle() {
			return
		}
	}
}
