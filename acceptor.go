package simplefixgo

import (
	"context"
	"errors"
	"fmt"
	"github.com/b2broker/simplefix-go/utils"
	"io"
	"net"
)

//type DefaultHandler interface {
//	Listen(ctx context.Context, conn *Conn) error
//}

type Sender interface {
	Send(message SendingMessage) error
}

type HandlerFunc func(msg []byte)

type AcceptorHandler interface {
	ServeIncoming(msg []byte)
	Outgoing() <-chan []byte
	Run() error
	StopWithError(err error)
	Send(message SendingMessage) error
	SendRaw(msgType string, message []byte) error
	RemoveIncomingHandler(msgType string, id int64) (err error)
	RemoveOutgoingHandler(msgType string, id int64) (err error)
	HandleIncoming(msgType string, handle HandlerFunc) (id int64)
	HandleOutgoing(msgType string, handle HandlerFunc) (id int64)
	OnDisconnect(handlerFunc utils.EventHandlerFunc)
	OnConnect(handlerFunc utils.EventHandlerFunc)
	OnStopped(handlerFunc utils.EventHandlerFunc)
}

type HandlerFactory interface {
	MakeHandler(ctx context.Context) AcceptorHandler
}

type Acceptor struct {
	listener        net.Listener
	factory         HandlerFactory
	size            int
	handleNewClient func(handler AcceptorHandler)

	ctx    context.Context
	cancel context.CancelFunc
}

var ErrServerClosed = errors.New("server closed")
var ErrClientDisconnect = errors.New("client disconnected")
var ErrHandlerStopped = errors.New("handler stopped")

func NewAcceptor(listener net.Listener, factory HandlerFactory, handleNewClient func(handler AcceptorHandler)) *Acceptor {
	s := &Acceptor{
		factory:         factory,
		listener:        listener,
		handleNewClient: handleNewClient,
	}

	s.ctx, s.cancel = context.WithCancel(context.Background())

	return s
}

func (s *Acceptor) Close() {
	s.cancel()
}

// ListenAndServe run listening and serving for connection
// start accepting connections of new clients
func (s *Acceptor) ListenAndServe() error {
	listenErr := make(chan error)
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
			return fmt.Errorf("could not accept conn: %w", err)

		case <-s.ctx.Done():
			return nil
		}
	}
}

// serve run listening and serving connection
// handle ClientConn's connection
func (s *Acceptor) serve(parentCtx context.Context, netConn net.Conn) {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	conn := NewConn(parentCtx, netConn, s.size)

	handler := s.factory.MakeHandler(ctx)

	connErr := make(chan error)
	go func() {
		connErr <- conn.serve()
	}()

	handlerErr := make(chan error)
	go func() {
		handlerErr <- handler.Run()
	}()

	defer conn.Close()

	if s.handleNewClient != nil {
		s.handleNewClient(handler)
	}

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

		case msg, ok := <-handler.Outgoing():
			if !ok {
				return
			}
			conn.Write(msg)
		}
	}
}
