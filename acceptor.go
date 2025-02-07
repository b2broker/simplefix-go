package simplefixgo

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"

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
	Stop()
	StopWithError(err error)
	CloseErrorChan()
	Send(message SendingMessage) error
	SendBatch(messages []SendingMessage) error
	SendBuffered(message SendingMessage) error
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
	writeTimeout    time.Duration

	ctx    context.Context
	cancel context.CancelFunc
}

// NewAcceptor is used to create a new Acceptor instance.
func NewAcceptor(listener net.Listener, factory HandlerFactory, writeTimeout time.Duration, handleNewClient func(handler AcceptorHandler)) *Acceptor {
	s := &Acceptor{
		factory:         factory,
		listener:        listener,
		handleNewClient: handleNewClient,
		writeTimeout:    writeTimeout,
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

	conn := NewConn(parentCtx, netConn, s.size, s.writeTimeout)
	defer conn.Close()

	cancelFun := func() {
		conn.Close()
		cancel()
	}

	handler := s.factory.MakeHandler(ctx)
	defer handler.CloseErrorChan()

	eg := errgroup.Group{}

	eg.Go(func() error {
		defer cancelFun()

		err := conn.serve()
		if err != nil {
			err = fmt.Errorf("%s: %w", err, ErrConnClosed)
			if !strings.Contains(err.Error(), "use of closed network connection") {
				handler.StopWithError(err)
			}
		}

		return err
	})

	if s.handleNewClient != nil {
		s.handleNewClient(handler)
	}

	eg.Go(func() error {
		defer cancelFun()

		return handler.Run()
	})

	eg.Go(func() error {
		defer cancelFun()

		for {
			select {
			case <-ctx.Done():
				return nil

			case msg, ok := <-handler.Outgoing():
				if !ok {
					return nil
				}

				err := conn.Write(msg)
				if err != nil {
					return err
				}
			}
		}
	})

	eg.Go(func() error {
		defer cancelFun()

		for {
			select {
			case <-ctx.Done():
				return nil

			case msg, ok := <-conn.Reader():
				if !ok {
					return nil
				}
				handler.ServeIncoming(msg)
			}
		}
	})

	_ = eg.Wait()
}
