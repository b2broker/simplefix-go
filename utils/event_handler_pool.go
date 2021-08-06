package utils

import "sync"

type Event int

const (
	// EventDisconnect calls when connection down
	EventDisconnect Event = iota

	// EventDisconnect calls when connection up
	EventConnect

	// EventStopped calls when handler stopped
	EventStopped

	// EventLogon calls when logon message received
	EventLogon

	// EventRequest calls when logon sent and Session waiting for answer from other side
	EventRequest
)

// EventHandlerFunc is a function for calling when event happened
type EventHandlerFunc func() bool

// EventHandlerPool is a service for saving and calling handlers
type EventHandlerPool struct {
	mu   sync.RWMutex
	pool map[Event][]EventHandlerFunc
}

// NewEventHandlerPool creates new EventHandlerPool
func NewEventHandlerPool() *EventHandlerPool {
	return &EventHandlerPool{pool: make(map[Event][]EventHandlerFunc)}
}

// Handle adds new handler for event
func (evp *EventHandlerPool) Handle(e Event, handle EventHandlerFunc) {
	evp.mu.Lock()
	defer evp.mu.Unlock()

	if _, ok := evp.pool[e]; !ok {
		evp.pool[e] = []EventHandlerFunc{}
	}

	evp.pool[e] = append(evp.pool[e], handle)
}

// Trigger calls all handlers for received event
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
