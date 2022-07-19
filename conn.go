package simplefixgo

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"net"
	"time"

	"golang.org/x/sync/errgroup"
)

// ErrConnClosed handles connection errors.
var ErrConnClosed = fmt.Errorf("the reader is closed")

const (
	endOfMsgTag = "10="
)

// Conn is a net.Conn wrapper that is used for handling split messages.
type Conn struct {
	reader chan []byte
	writer chan []byte
	conn   net.Conn

	ctx    context.Context
	cancel context.CancelFunc

	writeDeadline time.Duration
}

// NewConn is called to create a new connection.
func NewConn(ctx context.Context, conn net.Conn, msgBuffSize int, writeDeadline time.Duration) *Conn {
	c := &Conn{
		reader:        make(chan []byte, msgBuffSize),
		writer:        make(chan []byte, msgBuffSize),
		writeDeadline: writeDeadline,
		conn:          conn,
	}

	c.ctx, c.cancel = context.WithCancel(ctx)

	return c
}

// Close is called to cancel a connection context and close a connection.
func (c *Conn) Close() {
	_ = c.conn.Close()
	c.cancel()
}

func (c *Conn) serve() error {
	defer close(c.writer)
	defer close(c.reader)

	eg := errgroup.Group{}

	eg.Go(c.runReader)

	return eg.Wait()
}

func (c *Conn) runReader() error {
	defer c.cancel()
	r := bufio.NewReader(c.conn)

	var msg []byte
	for {
		select {
		case <-c.ctx.Done():
			return nil
		default:
		}

		buff, err := r.ReadBytes(byte(1))
		if err != nil {
			return fmt.Errorf("read error: %w", err)
		}

		msg = append(msg, buff...)
		if len(buff) >= 3 && bytes.Equal(buff[0:3], []byte(endOfMsgTag)) {
			c.reader <- msg
			msg = []byte{}
		}
	}
}

// Reader returns a separate channel for handing incoming messages.
func (c *Conn) Reader() <-chan []byte {
	return c.reader
}

// Write is called to send messages to an outgoing socket.
func (c *Conn) Write(msg []byte) error {
	select {
	case <-c.ctx.Done():
		return ErrConnClosed
	default:
	}

	if err := c.conn.SetWriteDeadline(time.Now().Add(c.writeDeadline)); err != nil {
		c.cancel()
		return fmt.Errorf("set write deadline error: %w", err)
	}
	_, err := c.conn.Write(msg)
	if err != nil {
		c.cancel()
		return fmt.Errorf("write error: %w", err)
	}

	return nil
}
