package simplefixgo

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net"
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
}

// NewConn is called to create a new connection.
func NewConn(ctx context.Context, conn net.Conn, msgBuffSize int) *Conn {
	c := &Conn{
		reader: make(chan []byte, msgBuffSize),
		writer: make(chan []byte, msgBuffSize),

		conn: conn,
	}

	c.ctx, c.cancel = context.WithCancel(ctx)

	return c
}

// Close is called to cancel a connection context and close a connection.
func (c *Conn) Close() {
	c.conn.Close()
	c.cancel()
}

func (c *Conn) serve() error {
	defer close(c.writer)
	defer close(c.reader)

	eg := errgroup.Group{}

	eg.Go(c.runWriter)
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

func (c *Conn) runWriter() error {
	defer c.cancel()

	for {
		select {
		case msg := <-c.writer:
			_, err := c.conn.Write(msg)
			if err != nil {
				return fmt.Errorf("write error: %w", err)
			}

		case <-c.ctx.Done():
			return nil
		}
	}
}

// Reader returns a separate channel for handing incoming messages.
func (c *Conn) Reader() <-chan []byte {
	return c.reader
}

// Write is called to send messages to an outgoing socket.
func (c *Conn) Write(msg []byte) {
	select {
	case <-c.ctx.Done():
		return
	default:
	}
	c.writer <- msg
}
