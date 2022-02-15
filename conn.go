package simplefixgo

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"net"
)

// ErrConnClosed connection error
var ErrConnClosed = fmt.Errorf("reader closed")

const (
	endOfMsgTag = "10="
)

// net.Conn wrapper for working with split messages
type Conn struct {
	reader chan []byte
	writer chan []byte
	conn   net.Conn

	ctx    context.Context
	cancel context.CancelFunc
}

// NewConn creates new Conn
func NewConn(ctx context.Context, conn net.Conn, msgBuffSize int) *Conn {
	c := &Conn{
		reader: make(chan []byte, msgBuffSize),
		writer: make(chan []byte, msgBuffSize),

		conn: conn,
	}

	c.ctx, c.cancel = context.WithCancel(ctx)

	return c
}

// Close cancels Conn context to stop work
func (c *Conn) Close() {
	c.cancel()
}

func (c *Conn) serve() error {
	defer c.conn.Close()

	errCh := make(chan error, 2)
	go c.runWriter(errCh)
	go c.runReader(errCh)

	select {
	case err := <-errCh:
		return err
	case <-c.ctx.Done():
		return nil
	}
}

func (c *Conn) runReader(errCh chan error) {
	r := bufio.NewReader(c.conn)
	defer c.cancel()

	var msg []byte
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
		}

		buff, err := r.ReadBytes(byte(1))
		if err != nil {
			errCh <- err
			return
		}

		msg = append(msg, buff...)
		if len(buff) >= 3 && bytes.Equal(buff[0:3], []byte(endOfMsgTag)) {
			c.reader <- msg
			msg = []byte{}
		}
	}
}

func (c *Conn) runWriter(errCh chan error) {
	defer c.cancel()

	for {
		select {
		case msg := <-c.writer:
			_, err := c.conn.Write(msg)
			if err != nil {
				errCh <- err
				return
			}

		case <-c.ctx.Done():
			errCh <- ErrConnClosed
			return
		}
	}
}

// Reader returns sole chan incoming with messages
func (c *Conn) Reader() <-chan []byte {
	return c.reader
}

// Write sends messages to outgoing socket
func (c *Conn) Write(msg []byte) {
	c.writer <- msg
}
