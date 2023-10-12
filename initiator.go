package simplefixgo

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

// InitiatorHandler is an interface implementing basic methods required for handling the Initiator object.
type InitiatorHandler interface {
	ServeIncoming(msg []byte)
	Outgoing() <-chan []byte
	Run() error
	StopWithError(err error)
	Send(message SendingMessage) error
	Context() context.Context
	Stop()
}

// Initiator provides the client-side service functionality.
type Initiator struct {
	conn    *Conn
	handler InitiatorHandler

	ctx    context.Context
	cancel context.CancelFunc
}

// NewInitiator creates a new Initiator instance.
func NewInitiator(conn net.Conn, handler InitiatorHandler, bufSize int, writeDeadline time.Duration) *Initiator {
	c := &Initiator{handler: handler}
	c.ctx, c.cancel = context.WithCancel(context.Background())

	c.conn = NewConn(c.ctx, conn, bufSize, writeDeadline)

	return c
}

// Close is used to cancel the specified Initiator context.
func (c *Initiator) Close() {
	c.conn.Close()
	c.cancel()
}

// Send is used to send a FIX message.
func (c *Initiator) Send(message SendingMessage) error {
	return c.handler.Send(message)
}

// Serve is used to initiate the procedure of delivering messages.
func (c *Initiator) Serve() error {
	eg := errgroup.Group{}
	defer c.Close()

	stopHandler := sync.Once{}

	eg.Go(func() error {
		defer c.Close()

		err := c.conn.serve()
		if err != nil {
			err = fmt.Errorf("%s: %w", err, ErrConnClosed)
			defer stopHandler.Do(func() {
				c.handler.StopWithError(err)
			})
		}

		return err
	})

	eg.Go(func() error {
		defer c.Close()

		return c.handler.Run()
	})

	eg.Go(func() error {
		defer c.Close()

		for {
			select {
			case <-c.ctx.Done():
				return nil

			case msg, ok := <-c.handler.Outgoing():
				if !ok {
					return fmt.Errorf("outgoing chan is closed")
				}

				err := c.conn.Write(msg)
				if err != nil {
					c.handler.Stop()

					return ErrConnClosed
				}
			}
		}
	})

	eg.Go(func() error {
		defer c.Close()
		select {
		case <-c.handler.Context().Done():
			stopHandler.Do(func() {})
		case <-c.ctx.Done():
			stopHandler.Do(func() {
				c.handler.StopWithError(nil)
			})
		}

		return nil
	})

	eg.Go(func() error {
		defer c.Close()

		for {
			select {
			case <-c.ctx.Done():
				return nil

			case msg, ok := <-c.conn.Reader():
				if !ok {
					continue
				}
				c.handler.ServeIncoming(msg)
			}
		}
	})

	err := eg.Wait()
	if err != nil {
		stopHandler.Do(func() {
			c.handler.StopWithError(err)
		})
		return fmt.Errorf("stop handler: %w", err)
	}

	return nil
}
