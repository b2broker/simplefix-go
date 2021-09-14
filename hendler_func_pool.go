package simplefixgo

import (
	"errors"
	"sync"
	"sync/atomic"
)

// ErrHandleNotFound will be returned when looking handler not found
var ErrHandleNotFound = errors.New("handler not found")

// HandlerPool is a structure for work with pool of message handlers
type HandlerPool struct {
	mu       sync.RWMutex
	handlers map[string]map[int64]interface{}
	counter  *int64
}

// NewHandlerPool creates new HandlerPool
func NewHandlerPool() *HandlerPool {
	return &HandlerPool{
		handlers: make(map[string]map[int64]interface{}),
		counter:  new(int64),
	}
}

func (p *HandlerPool) inc() int64 {
	return atomic.AddInt64(p.counter, 1)
}

func (p *HandlerPool) free(msgType string) {
	if len(p.handlers[msgType]) != 0 {
		return
	}

	delete(p.handlers, msgType)
}

// Remove removes handler by its id
func (p *HandlerPool) Remove(msgType string, id int64) error {
	if _, ok := p.handlers[msgType]; !ok {
		return ErrHandleNotFound
	}

	if _, ok := p.handlers[msgType][id]; !ok {
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
	for _, handler := range handlers {
		result = append(result, handler)
	}

	return result
}

func (p *HandlerPool) add(msgType string, handle interface{}) int64 {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, ok := p.handlers[msgType]; !ok {
		p.handlers[msgType] = make(map[int64]interface{})
	}

	id := p.inc()
	p.handlers[msgType][id] = handle

	return id
}

// IncomingHandlerPool is a pool for incoming messages with bytes
type IncomingHandlerPool struct {
	*HandlerPool
}

// NewHandlerPool creates new HandlerPool
func NewIncomingHandlerPool() IncomingHandlerPool {
	return IncomingHandlerPool{NewHandlerPool()}
}

// Range is a handlers traversal
// it will be stop if one of handlers returns false
func (p IncomingHandlerPool) Range(msgType string, f func(IncomingHandlerFunc) bool) {
	for _, handle := range p.handlersByMsgType(msgType) {
		if !f(handle.(IncomingHandlerFunc)) {
			break
		}
	}
}

// Add adds new message handler for message type
// returns id of added message
func (p *IncomingHandlerPool) Add(msgType string, handle IncomingHandlerFunc) int64 {
	return p.add(msgType, handle)
}

// IncomingHandlerPool is a pool for outgoing messages with structs
type OutgoingHandlerPool struct {
	*HandlerPool
}

// NewHandlerPool creates new HandlerPool
func NewOutgoingHandlerPool() OutgoingHandlerPool {
	return OutgoingHandlerPool{NewHandlerPool()}
}

// Range is a handlers traversal
// it will be stop if one of handlers returns false
func (p OutgoingHandlerPool) Range(msgType string, f func(OutgoingHandlerFunc) bool) {
	for _, handle := range p.handlersByMsgType(msgType) {
		if !f(handle.(OutgoingHandlerFunc)) {
			break
		}
	}
}

// Add adds new message handler for message type
// returns id of added message
func (p *HandlerPool) Add(msgType string, handle OutgoingHandlerFunc) int64 {
	return p.add(msgType, handle)
}
