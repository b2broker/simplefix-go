package simplefixgo

import (
	"context"
	"fmt"
	"net"
)

type OutgoingMessage interface {
	ToBytes() ([]byte, error)
}

type InitiatorHandler interface {
	ServeIncoming(msg []byte)
	Outgoing() <-chan []byte
	Run() error
	StopWithError(err error)
	Send(message SendingMessage) error
}

type Initiator struct {
	conn    *Conn
	handler InitiatorHandler

	ctx    context.Context
	cancel context.CancelFunc
}

func NewInitiator(conn net.Conn, handler InitiatorHandler, bufSize int) *Initiator {
	c := &Initiator{handler: handler}
	c.ctx, c.cancel = context.WithCancel(context.Background())

	c.conn = NewConn(c.ctx, conn, bufSize)

	return c
}

func (c *Initiator) Close() {
	c.cancel()
}

func (c *Initiator) Send(message SendingMessage) error {
	return c.handler.Send(message)
}

func (c *Initiator) Serve() error {
	connErr := make(chan error)
	go func() {
		connErr <- c.conn.serve()
	}()

	handlerErr := make(chan error)
	go func() {
		handlerErr <- c.handler.Run()
	}()

	defer c.conn.Close()

	for {
		select {
		case err := <-handlerErr:
			return fmt.Errorf("handler error: %w", err)

		case err := <-connErr:
			c.handler.StopWithError(ErrConnClosed)
			return fmt.Errorf("%w: %s", ErrConnClosed, err)

		case <-c.ctx.Done():
			return nil

		case msg, ok := <-c.conn.Reader():
			if !ok {
				return fmt.Errorf("conn reader chan was closed")
			}
			c.handler.ServeIncoming(msg)

		case msg, ok := <-c.handler.Outgoing():
			if !ok {
				return fmt.Errorf("handler chan was closed")
			}
			c.conn.Write(msg)
		}
	}
}
