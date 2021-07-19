package simplefixgo

import (
	"context"
)

type AcceptorHandlerFactory struct {
	bufferSize int
	msgTypeTag string
}

func NewAcceptorHandlerFactory(msgTypeTag string, bufferSize int) *AcceptorHandlerFactory {
	return &AcceptorHandlerFactory{bufferSize: bufferSize, msgTypeTag: msgTypeTag}
}

func (h *AcceptorHandlerFactory) MakeHandler(ctx context.Context) AcceptorHandler {
	return NewAcceptorHandler(ctx, h.msgTypeTag, h.bufferSize)
}
