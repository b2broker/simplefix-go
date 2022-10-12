package simplefixgo

import (
	"context"
	"errors"
	"fmt"
	"github.com/b2broker/simplefix-go/fix"
	"github.com/b2broker/simplefix-go/session/messages"
	"github.com/b2broker/simplefix-go/utils"
	"sync"
)

const AllMsgTypes = "ALL"

// SendingMessage provides a basic method for sending messages.
type SendingMessage interface {
	HeaderBuilder() messages.HeaderBuilder
	MsgType() string
	ToBytes() ([]byte, error)
}

// DefaultHandler is a standard handler for the Acceptor and Initiator objects.
type DefaultHandler struct {
	mu sync.Mutex

	out      chan []byte
	incoming chan []byte

	incomingHandlers IncomingHandlerPool
	outgoingHandlers OutgoingHandlerPool

	eventHandlers *utils.EventHandlerPool

	msgTypeTag string

	ctx    context.Context
	cancel context.CancelFunc
	errors chan error
}

// NewAcceptorHandler creates a handler for an Acceptor object.
func NewAcceptorHandler(ctx context.Context, msgTypeTag string, bufferSize int) *DefaultHandler {
	sh := &DefaultHandler{
		msgTypeTag:    msgTypeTag,
		eventHandlers: utils.NewEventHandlerPool(),

		out:      make(chan []byte, bufferSize),
		incoming: make(chan []byte, bufferSize),
		errors:   make(chan error),

		incomingHandlers: NewIncomingHandlerPool(),
		outgoingHandlers: NewOutgoingHandlerPool(),
	}

	sh.ctx, sh.cancel = context.WithCancel(ctx)

	return sh
}

// NewInitiatorHandler creates a handler for the Initiator object.
func NewInitiatorHandler(ctx context.Context, msgTypeTag string, bufferSize int) *DefaultHandler {
	sh := &DefaultHandler{
		msgTypeTag:    msgTypeTag,
		eventHandlers: utils.NewEventHandlerPool(),

		out:      make(chan []byte, bufferSize),
		incoming: make(chan []byte, bufferSize),
		errors:   make(chan error),

		incomingHandlers: NewIncomingHandlerPool(),
		outgoingHandlers: NewOutgoingHandlerPool(),
	}

	sh.ctx, sh.cancel = context.WithCancel(ctx)

	return sh
}

func (h *DefaultHandler) sendRaw(data []byte) error {
	select {
	case h.out <- data:
	case <-h.ctx.Done():
		return fmt.Errorf("the handler is stopped")
	}
	return nil
}

func (h *DefaultHandler) send(msg SendingMessage) error {
	ok := h.outgoingHandlers.Range(AllMsgTypes, func(handle OutgoingHandlerFunc) bool {
		return handle(msg)
	})
	if !ok {
		return errors.New("the handler for all message types has refused the message and returned false")
	}

	ok = h.outgoingHandlers.Range(msg.MsgType(), func(handle OutgoingHandlerFunc) bool {
		return handle(msg)
	})
	if !ok {
		return errors.New("the handler for the current type has refused the message and returned false")
	}

	data, err := msg.ToBytes()
	if err != nil {
		return err
	}

	return h.sendRaw(data)
}

// SendRaw sends a message in the byte array format
// without involving any additional handlers.
func (h *DefaultHandler) SendRaw(data []byte) error {
	return h.sendRaw(data)
}

// Send is a function that sends a previously prepared message.
func (h *DefaultHandler) Send(message SendingMessage) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	return h.send(message)
}

// SendBatch is a function that sends previously prepared messages.
func (h *DefaultHandler) SendBatch(messages []SendingMessage) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	for _, message := range messages {
		err := h.send(message)
		if err != nil {
			return err
		}
	}

	return nil
}

// RemoveIncomingHandler removes an existing handler for incoming messages.
func (h *DefaultHandler) RemoveIncomingHandler(msgType string, id int64) (err error) {
	return h.incomingHandlers.Remove(msgType, id)
}

// RemoveOutgoingHandler removes an existing handler for outgoing messages.
func (h *DefaultHandler) RemoveOutgoingHandler(msgType string, id int64) (err error) {
	return h.outgoingHandlers.Remove(msgType, id)
}

// HandleIncoming subscribes a handler function to incoming messages with a specific msgType.
// To subscribe to all messages, specify the AllMsgTypes constant for the msgType field
// (such messages will have a higher priority than the ones assigned to specific handlers).
func (h *DefaultHandler) HandleIncoming(msgType string, handle IncomingHandlerFunc) (id int64) {
	return h.incomingHandlers.Add(msgType, handle)
}

// HandleOutgoing subscribes a handler function to outgoing messages with a specific msgType
// (this may be required for modifying messages before sending).
// To subscribe to all messages, specify the AllMsgTypes constant for the msgType field
// (such messages will have a higher priority than the ones assigned to specific handlers).
func (h *DefaultHandler) HandleOutgoing(msgType string, handle OutgoingHandlerFunc) (id int64) {
	return h.outgoingHandlers.Add(msgType, handle)
}

// ServeIncoming is an internal method for handling incoming messages.
func (h *DefaultHandler) ServeIncoming(msg []byte) {
	h.incoming <- msg
}

func (h *DefaultHandler) serve(msg []byte) (err error) {
	msgTypeB, err := fix.ValueByTag(msg, h.msgTypeTag)
	if err != nil {
		return fmt.Errorf("msg type: %w", err)
	}
	msgType := string(msgTypeB)

	h.incomingHandlers.Range(AllMsgTypes, func(handle IncomingHandlerFunc) bool {
		return handle(msg)
	})

	h.incomingHandlers.Range(msgType, func(handle IncomingHandlerFunc) bool {
		return handle(msg)
	})

	return nil
}

// Run is a function that is used for listening and processing messages.
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

func (h *DefaultHandler) Context() context.Context {
	return h.ctx
}

// Outgoing is a service method that provides an outgoing channel
// to the server or client connection manager.
func (h *DefaultHandler) Outgoing() <-chan []byte {
	return h.out
}

// Stop is a function that enables graceful termination of a session.
func (h *DefaultHandler) Stop() {
	h.cancel()
}

// StopWithError is a function that enables graceful termination of a session with throwing an error.
func (h *DefaultHandler) StopWithError(err error) {
	h.errors <- err
}

// OnDisconnect handles disconnection events.
func (h *DefaultHandler) OnDisconnect(handlerFunc utils.EventHandlerFunc) {
	h.eventHandlers.Handle(utils.EventDisconnect, handlerFunc)
}

// OnConnect handles connection events.
func (h *DefaultHandler) OnConnect(handlerFunc utils.EventHandlerFunc) {
	h.eventHandlers.Handle(utils.EventConnect, handlerFunc)
}

// OnStopped handles session termination events.
func (h *DefaultHandler) OnStopped(handlerFunc utils.EventHandlerFunc) {
	h.eventHandlers.Handle(utils.EventStopped, handlerFunc)
}
