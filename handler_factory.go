package simplefixgo

import (
	"context"
)

// AcceptorHandlerFactory handler factory for Acceptor
type AcceptorHandlerFactory struct {
	bufferSize int
	msgTypeTag string
}

// NewAcceptorHandlerFactory returns new AcceptorHandlerFactory
func NewAcceptorHandlerFactory(msgTypeTag string, bufferSize int) *AcceptorHandlerFactory {
	return &AcceptorHandlerFactory{bufferSize: bufferSize, msgTypeTag: msgTypeTag}
}

// NewAcceptorHandlerFactory makes AcceptorHandler
func (h *AcceptorHandlerFactory) MakeHandler(ctx context.Context) AcceptorHandler {
	return NewAcceptorHandler(ctx, h.msgTypeTag, h.bufferSize)
}
