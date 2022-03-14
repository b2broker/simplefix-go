package simplefixgo

import (
	"context"
)

// AcceptorHandlerFactory is a handler factory for an Acceptor object.
type AcceptorHandlerFactory struct {
	bufferSize int
	msgTypeTag string
}

// NewAcceptorHandlerFactory returns a new AcceptorHandlerFactory instance.
func NewAcceptorHandlerFactory(msgTypeTag string, bufferSize int) *AcceptorHandlerFactory {
	return &AcceptorHandlerFactory{bufferSize: bufferSize, msgTypeTag: msgTypeTag}
}

// MakeHandler creates a new AcceptorHandler instance.
func (h *AcceptorHandlerFactory) MakeHandler(ctx context.Context) AcceptorHandler {
	return NewAcceptorHandler(ctx, h.msgTypeTag, h.bufferSize)
}
