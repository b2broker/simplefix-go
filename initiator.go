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

	ctx    context.Context
	cancel context.CancelFunc
}

// NewInitiator creates new Initiator
func NewInitiator(conn net.Conn, handler InitiatorHandler, bufSize int) *Initiator {
	c := &Initiator{handler: handler}
	c.ctx, c.cancel = context.WithCancel(context.Background())

	c.conn = NewConn(c.ctx, conn, bufSize)

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
	}()

	handlerErr := make(chan error)
	go func() {
		handlerErr <- c.handler.Run()
	}()

	defer c.handler.StopWithError(fmt.Errorf("initiator closed"))
	defer c.conn.Close()

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
				continue
			}
			return nil

		case <-c.ctx.Done():
			return nil

		case msg, ok := <-c.conn.Reader():
			if !ok {
				continue
			}
			c.handler.ServeIncoming(msg)
		}
	}
}
