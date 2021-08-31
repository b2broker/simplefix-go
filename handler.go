package simplefixgo

import (
	"context"
	"errors"
	"fmt"
	"github.com/b2broker/simplefix-go/fix"
	"github.com/b2broker/simplefix-go/utils"
)

const AllMsgTypes = "ALL"

// SendingMessage basic method for sending message
type SendingMessage interface {
	MsgType() string
	ToBytes() ([]byte, error)
}

// DefaultHandler is a standard handler for
type DefaultHandler struct {
	out      chan []byte
	incoming chan []byte

	incomingHandlers *HandlerPool
	outgoingHandlers *HandlerPool

	eventHandlers *utils.EventHandlerPool

	msgTypeTag string

	ctx    context.Context
	cancel context.CancelFunc
	errors chan error
}

// NewAcceptorHandler creates handler for Acceptor
func NewAcceptorHandler(ctx context.Context, msgTypeTag string, bufferSize int) *DefaultHandler {
	sh := &DefaultHandler{
		msgTypeTag:    msgTypeTag,
		eventHandlers: utils.NewEventHandlerPool(),

		out:      make(chan []byte, bufferSize),
		incoming: make(chan []byte, bufferSize),
		errors:   make(chan error),

		incomingHandlers: NewHandlerPool(),
		outgoingHandlers: NewHandlerPool(),
	}

	sh.ctx, sh.cancel = context.WithCancel(ctx)

	return sh
}

// NewInitiatorHandler creates handler for Initiator
func NewInitiatorHandler(ctx context.Context, msgTypeTag string, bufferSize int) *DefaultHandler {
	sh := &DefaultHandler{
		msgTypeTag:    msgTypeTag,
		eventHandlers: utils.NewEventHandlerPool(),

		out:      make(chan []byte, bufferSize),
		incoming: make(chan []byte, bufferSize),
		errors:   make(chan error),

		incomingHandlers: NewHandlerPool(),
		outgoingHandlers: NewHandlerPool(),
	}

	sh.ctx, sh.cancel = context.WithCancel(ctx)

	return sh
}

func (h *DefaultHandler) send(msgType string, data []byte) error {
	h.outgoingHandlers.Range(AllMsgTypes, func(handle HandlerFunc) bool {
		return handle(data)
	})

	h.outgoingHandlers.Range(msgType, func(handle HandlerFunc) bool {
		return handle(data)
	})

	h.out <- data

	return nil
}

// SendRaw sends raw message
func (h *DefaultHandler) SendRaw(msgType string, message []byte) error {
	return h.send(msgType, message)
}

// Send sends prepared message
func (h *DefaultHandler) Send(message SendingMessage) error {
	data, err := message.ToBytes()
	if err != nil {
		return fmt.Errorf("converting to bytes: %s", err)
	}

	return h.send(message.MsgType(), data)
}

// RemoveOutgoingHandler removes existing incoming handler
func (h *DefaultHandler) RemoveIncomingHandler(msgType string, id int64) (err error) {
	return h.incomingHandlers.Remove(msgType, id)
}

// RemoveOutgoingHandler removes existing outgoing handler
func (h *DefaultHandler) RemoveOutgoingHandler(msgType string, id int64) (err error) {
	return h.outgoingHandlers.Remove(msgType, id)
}

// HandleIncoming subscribes handler function to incoming messages with specific msgType
// For subscription to all messages use AllMsgTypes constant for field msgType
// in this case your messages will have high priority
func (h *DefaultHandler) HandleIncoming(msgType string, handle HandlerFunc) (id int64) {
	return h.incomingHandlers.Add(msgType, handle)
}

// HandleOutgoing subscribes handler function to outgoing messages with specific msgType
// for modification before sending
// For subscription to all messages use AllMsgTypes constant for field msgType
// in this case your messages will have high priority
func (h *DefaultHandler) HandleOutgoing(msgType string, handle HandlerFunc) (id int64) {
	return h.outgoingHandlers.Add(msgType, handle)
}

// ServeIncoming is inner method for handle incoming messages
func (h *DefaultHandler) ServeIncoming(msg []byte) {
	h.incoming <- msg
}

func (h *DefaultHandler) serve(msg []byte) (err error) {
	msgTypeB, err := fix.ValueByTag(msg, h.msgTypeTag)
	if err != nil {
		return fmt.Errorf("msg type: %w", err)
	}
	msgType := string(msgTypeB)

	h.incomingHandlers.Range(AllMsgTypes, func(handle HandlerFunc) bool {
		return handle(msg)
	})

	h.incomingHandlers.Range(msgType, func(handle HandlerFunc) bool {
		return handle(msg)
	})

	return nil
}

// Run starts listen and serve messages
func (h *DefaultHandler) Run() (err error) {
	h.eventHandlers.Trigger(utils.EventConnect)

	for {
		select {
		case msg, ok := <-h.incoming:
			if !ok {
				return ErrConnClosed
			}

			err = h.serve(msg)
			if err != nil {
				return err
			}

		case <-h.ctx.Done():
			h.eventHandlers.Trigger(utils.EventStopped)

			return

		case err := <-h.errors:
			if errors.Is(err, ErrConnClosed) {
				h.eventHandlers.Trigger(utils.EventDisconnect)
			}

			return err
		}
	}
}

// Outgoing is service method for returning outgoing chan to server or client connection manager
func (h *DefaultHandler) Outgoing() <-chan []byte {
	return h.out
}

// Stop is a graceful stop
func (h *DefaultHandler) Stop() {
	h.cancel()
}

// Stop is a graceful stop with error
func (h *DefaultHandler) StopWithError(err error) {
	h.errors <- err
}

// OnDisconnect handles disconnect event
func (h *DefaultHandler) OnDisconnect(handlerFunc utils.EventHandlerFunc) {
	h.eventHandlers.Handle(utils.EventDisconnect, handlerFunc)
}

// OnDisconnect handles disconnect event
func (h *DefaultHandler) OnConnect(handlerFunc utils.EventHandlerFunc) {
	h.eventHandlers.Handle(utils.EventConnect, handlerFunc)
}

// OnStopped handles Close event
func (h *DefaultHandler) OnStopped(handlerFunc utils.EventHandlerFunc) {
	h.eventHandlers.Handle(utils.EventStopped, handlerFunc)
}
