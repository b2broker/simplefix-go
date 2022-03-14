package utils

import "sync"

type Event int

const (
	// EventDisconnect occurs when the connection is down.
	EventDisconnect Event = iota

	// EventDisconnect occurs when the connection is up.
	EventConnect

	// EventStopped occurs when the handler is stopped.
	EventStopped

	// EventLogon occurs upon receiving the Logon message.
	EventLogon

	// EventRequest occurs upon sending the Logon message,
	// after which the Session awaits an answer from the counterparty.
	EventRequest

	// EventLogout occurs upon receiving the Logout message.
	EventLogout
)

// EventHandlerFunc is called when an event occurs.
type EventHandlerFunc func() bool

// EventHandlerPool is a service required for saving and calling the event handlers.
type EventHandlerPool struct {
	mu   sync.RWMutex
	pool map[Event][]EventHandlerFunc
}

// NewEventHandlerPool creates a new EventHandlerPool instance.
func NewEventHandlerPool() *EventHandlerPool {
	return &EventHandlerPool{pool: make(map[Event][]EventHandlerFunc)}
}

// Handle adds a new handler for an event.
func (evp *EventHandlerPool) Handle(e Event, handle EventHandlerFunc) {
	evp.mu.Lock()
	defer evp.mu.Unlock()

	if _, ok := evp.pool[e]; !ok {
		evp.pool[e] = []EventHandlerFunc{}
	}

	evp.pool[e] = append(evp.pool[e], handle)
}

// Trigger calls all handlers associated with an occurring event.
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

func (evp *EventHandlerPool) Clean() {
	evp.mu.Lock()
	defer evp.mu.Unlock()

	evp.pool = nil
}
