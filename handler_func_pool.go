package simplefixgo

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrHandleNotFound = errors.New("handler not found")

type HandlerPool struct {
	mu       sync.RWMutex
	handlers map[string]map[int64]HandlerFunc
	counter  *int64
}

func NewHandlerPool() *HandlerPool {
	return &HandlerPool{
		handlers: make(map[string]map[int64]HandlerFunc),
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

func (p *HandlerPool) Range(msgType string, f func(HandlerFunc) bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if handlers, ok := p.handlers[msgType]; ok {
		for _, handle := range handlers {
			if !f(handle) {
				break
			}
		}
	}
}

func (p *HandlerPool) Add(msgType string, handle HandlerFunc) int64 {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, ok := p.handlers[msgType]; !ok {
		p.handlers[msgType] = make(map[int64]HandlerFunc)
	}

	id := p.inc()
	p.handlers[msgType][id] = handle

	return id
}
