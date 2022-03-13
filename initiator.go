package simplefixgo

import (
	"context"
	"fmt"
	"net"
)

// InitiatorHandler is an interface implementing basic methods required for handling the Initiator object.
type InitiatorHandler interface {
	ServeIncoming(msg []byte)
	Outgoing() <-chan []byte
	Run() error
	StopWithError(err error)
	Send(message SendingMessage) error
}

// Initiator provides the client-side service functionality.
type Initiator struct {
	conn    *Conn
	handler InitiatorHandler

	ctx    context.Context
	cancel context.CancelFunc
}

// NewInitiator creates a new Initiator instance.
func NewInitiator(conn net.Conn, handler InitiatorHandler, bufSize int) *Initiator {
	c := &Initiator{handler: handler}
	c.ctx, c.cancel = context.WithCancel(context.Background())

	c.conn = NewConn(c.ctx, conn, bufSize)

	return c
}

// Close is used to cancel the specified Initiator context.
func (c *Initiator) Close() {
	c.cancel()
}

// Send is used to send a FIX message.
func (c *Initiator) Send(message SendingMessage) error {
	return c.handler.Send(message)
}

// Serve is used to initiate the procedure of delivering messages.
func (c *Initiator) Serve() error {
	connErr := make(chan error)
	go func() {
		connErr <- c.conn.serve()
		close(connErr)
	}()

	handlerErr := make(chan error, 1)
	go func() {
		handlerErr <- c.handler.Run()
		close(handlerErr)
	}()

	defer func() {
		c.cancel()
		c.conn.Close()
	}()

	go func() {
		for {
			select {
			case <-c.ctx.Done():
				return

			case msg, ok := <-c.handler.Outgoing():
				if !ok {
					return
				}
				c.conn.Write(msg)
			}
		}
	}()

	for {
		select {
		case err := <-handlerErr:
			return fmt.Errorf("handler error: %w", err)

		case err := <-connErr:
			if err != nil {
				c.handler.StopWithError(ErrConnClosed)
				return fmt.Errorf("%w: %s", ErrConnClosed, err)
			}

		case <-c.ctx.Done():
			return nil

		case msg, ok := <-c.conn.Reader():
			if !ok {
				return fmt.Errorf("the connection reader channel was closed")
			}
			c.handler.ServeIncoming(msg)
		}
	}
}
