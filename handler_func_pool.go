package simplefixgo

import (
	"errors"
	"sync"
)

// ErrHandleNotFound is returned when a required handler is not found.
var ErrHandleNotFound = errors.New("handler not found")

// HandlerPool is used for managing the pool of message handlers.
type HandlerPool struct {
	mu       sync.RWMutex
	handlers map[string][]interface{}
	counter  *int64
}

// NewHandlerPool creates a new HandlerPool instance.
func NewHandlerPool() *HandlerPool {
	return &HandlerPool{
		handlers: make(map[string][]interface{}),
		counter:  new(int64),
	}
}

func (p *HandlerPool) free(msgType string) {
	if len(p.handlers[msgType]) != 0 {
		return
	}

	delete(p.handlers, msgType)
}

// Remove is used to remove a handler with a specified identifier.
func (p *HandlerPool) Remove(msgType string, _ int64) error {
	if _, ok := p.handlers[msgType]; !ok {
		return ErrHandleNotFound
	}

	p.free(msgType)

	return nil
}

func (p *HandlerPool) handlersByMsgType(msgType string) (result []interface{}) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	handlers, ok := p.handlers[msgType]
	if !ok {
		return
	}

	result = make([]interface{}, 0, len(handlers))
	result = append(result, handlers...)

	return result
}

func (p *HandlerPool) add(msgType string, handle interface{}) int64 {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, ok := p.handlers[msgType]; !ok {
		p.handlers[msgType] = make([]interface{}, 0)
	}

	p.handlers[msgType] = append(p.handlers[msgType], handle)

	return int64(len(p.handlers)) - 1
}

// IncomingHandlerPool is used to manage the pool of incoming messages stored in the form of byte arrays.
type IncomingHandlerPool struct {
	*HandlerPool
}

// NewIncomingHandlerPool creates a new HandlerPool instance.
func NewIncomingHandlerPool() IncomingHandlerPool {
	return IncomingHandlerPool{NewHandlerPool()}
}

// Range is used for traversal through handlers. The traversal stops if any handler returns false.
func (p IncomingHandlerPool) Range(msgType string, f func(IncomingHandlerFunc) bool) {
	for _, handle := range p.handlersByMsgType(msgType) {
		if !f(handle.(IncomingHandlerFunc)) {
			break
		}
	}
}

// Add is used to add a new message handler for the specified message type.
// The function returns the ID of a message for which a handler was added.
func (p *IncomingHandlerPool) Add(msgType string, handle IncomingHandlerFunc) int64 {
	return p.add(msgType, handle)
}

// OutgoingHandlerPool is used to manage the pool of outgoing messages stored as structures.
type OutgoingHandlerPool struct {
	*HandlerPool
}

// NewOutgoingHandlerPool creates a new OutgoingHandlerPool instance.
func NewOutgoingHandlerPool() OutgoingHandlerPool {
	return OutgoingHandlerPool{NewHandlerPool()}
}

// Range is used for traversal through handlers.
// The traversal stops if any handler returns false.
func (p OutgoingHandlerPool) Range(msgType string, f func(OutgoingHandlerFunc) bool) (res bool) {
	for _, handle := range p.handlersByMsgType(msgType) {
		if !f(handle.(OutgoingHandlerFunc)) {
			return false
		}
	}

	return true
}

// Add is used to add a new message handler for the specified message type.
// The function returns the ID of a message for which a handler was added.
func (p *HandlerPool) Add(msgType string, handle OutgoingHandlerFunc) int64 {
	return p.add(msgType, handle)
}
