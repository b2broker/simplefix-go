package simplefixgo

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/b2broker/simplefix-go/utils"
)

// Sender is an interface implemented by any structure that can issue a SendingMessage.
type Sender interface {
	Send(message SendingMessage) error
}

// OutgoingHandlerFunc is used for handling outgoing messages.
type OutgoingHandlerFunc func(msg SendingMessage) bool

// IncomingHandlerFunc is used for handling incoming messages.
type IncomingHandlerFunc func(data []byte) bool

// AcceptorHandler is a set of methods providing basic functionality to the Acceptor.
type AcceptorHandler interface {
	ServeIncoming(msg []byte)
	Outgoing() <-chan []byte
	Run() error
	StopWithError(err error)
	Send(message SendingMessage) error
	SendRaw(data []byte) error
	RemoveIncomingHandler(msgType string, id int64) (err error)
	RemoveOutgoingHandler(msgType string, id int64) (err error)
	HandleIncoming(msgType string, handle IncomingHandlerFunc) (id int64)
	HandleOutgoing(msgType string, handle OutgoingHandlerFunc) (id int64)
	OnDisconnect(handlerFunc utils.EventHandlerFunc)
	OnConnect(handlerFunc utils.EventHandlerFunc)
	OnStopped(handlerFunc utils.EventHandlerFunc)
	Context() context.Context
}

// HandlerFactory creates handlers for the Acceptor.
type HandlerFactory interface {
	MakeHandler(ctx context.Context) AcceptorHandler
}

// Acceptor is a server-side service used for handling client connections.
type Acceptor struct {
	listener        net.Listener
	factory         HandlerFactory
	size            int
	handleNewClient func(handler AcceptorHandler)

	ctx    context.Context
	cancel context.CancelFunc
}

// NewAcceptor is used to create a new Acceptor instance.
func NewAcceptor(listener net.Listener, factory HandlerFactory, handleNewClient func(handler AcceptorHandler)) *Acceptor {
	s := &Acceptor{
		factory:         factory,
		listener:        listener,
		handleNewClient: handleNewClient,
	}

	s.ctx, s.cancel = context.WithCancel(context.Background())

	return s
}

// Close is called to cancel the Acceptor context and close a connection.
func (s *Acceptor) Close() {
	s.cancel()
}

// ListenAndServe is used for listening and maintaining connections.
// It verifies and accepts connection requests from new clients.
func (s *Acceptor) ListenAndServe() error {
	listenErr := make(chan error, 1)
	defer s.Close()
	defer s.listener.Close()

	go func() {
		for {
			conn, err := s.listener.Accept()
			if err != nil {
				listenErr <- err
				return
			}

			go s.serve(s.ctx, conn)
		}
	}()

	for {
		select {
		case err := <-listenErr:
			return fmt.Errorf("could not accept connection: %w", err)

		case <-s.ctx.Done():
			return nil
		}
	}
}

// serve is used for listening and maintaining connections.
// It handles client connections opened for ClientConn instances.
func (s *Acceptor) serve(parentCtx context.Context, netConn net.Conn) {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	conn := NewConn(parentCtx, netConn, s.size)

	handler := s.factory.MakeHandler(ctx)

	connErr := make(chan error, 1)
	go func() {
		connErr <- conn.serve()
	}()

	handlerErr := make(chan error, 1)

	defer conn.Close()

	if s.handleNewClient != nil {
		s.handleNewClient(handler)
	}

	go func() {
		handlerErr <- handler.Run()
	}()

	go func() {
		for {
			select {
			case <-s.ctx.Done():
				return

			case msg, ok := <-handler.Outgoing():
				if !ok {
					return
				}
				conn.Write(msg)
			}
		}
	}()

	for {
		select {
		case <-handlerErr:
			return

		case err := <-connErr:
			if errors.Is(err, io.EOF) {
				handler.StopWithError(ErrConnClosed)
				return
			}
			handler.StopWithError(err)
			return

		case <-s.ctx.Done():
			return

		case msg, ok := <-conn.Reader():
			if !ok {
				return
			}
			handler.ServeIncoming(msg)
		}
	}
}
