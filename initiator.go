package simplefixgo

import (
	"context"
	"fmt"
	"net"
)

// InitiatorHandler is a basic methods of handling the initiator
type InitiatorHandler interface {
	ServeIncoming(msg []byte)
	Outgoing() <-chan []byte
	Run() error
	StopWithError(err error)
	Send(message SendingMessage) error
}

// Initiator is a client side service
type Initiator struct {
	conn    *Conn
	handler InitiatorHandler

	writer chan []byte

	ctx    context.Context
	cancel context.CancelFunc
}

// NewInitiator creates new Initiator
func NewInitiator(conn net.Conn, handler InitiatorHandler, bufSize int) *Initiator {
	c := &Initiator{handler: handler}
	c.ctx, c.cancel = context.WithCancel(context.Background())

	writer := make(chan []byte, bufSize)
	c.writer = writer

	c.conn = NewConn(c.ctx, conn, writer, bufSize)

	return c
}

// Close cancels context
func (c *Initiator) Close() {
	c.cancel()
}

// Send message
func (c *Initiator) Send(message SendingMessage) error {
	return c.handler.Send(message)
}

// Serve starts process of serving messages
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
				close(c.writer)
				return

			case msg, ok := <-c.handler.Outgoing():
				if !ok {
					return
				}
				c.writer <- msg
			}
		}
	}()

	for {
		select {
		case err := <-handlerErr:
			if err != nil {
				return fmt.Errorf("handler error: %w", err)
			}
			return fmt.Errorf("handler has been stopped")

		case err := <-connErr:
			if err != nil {
				c.handler.StopWithError(ErrConnClosed)
				return fmt.Errorf("%w: %s", ErrConnClosed, err)
			}

		case <-c.ctx.Done():
			return nil

		case msg, ok := <-c.conn.Reader():
			if !ok {
				return fmt.Errorf("conn reader chan was closed")
			}
			c.handler.ServeIncoming(msg)
		}
	}
}
