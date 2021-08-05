package tests

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	simplefixgo "github.com/b2broker/simplefix-go"
	"github.com/b2broker/simplefix-go/fix"
	"github.com/b2broker/simplefix-go/session"
	"github.com/b2broker/simplefix-go/session/storages/memory"
	fixgen "github.com/b2broker/simplefix-go/tests/fix44"
	"github.com/b2broker/simplefix-go/utils"
	"net"
	"sync"
	"testing"
	"time"
)

func TestHeartbeat(t *testing.T) {
	const (
		countOfHeartbeats = 4
		heartBtInt        = 1
		port              = 9991
	)

	// close acceptor after work
	acceptor := RunAcceptor(port, t, memory.NewStorage(100, 100))
	defer acceptor.Close()
	go func() {
		err := acceptor.ListenAndServe()
		if err != nil && !errors.Is(err, simplefixgo.ErrServerClosed) {
			panic(err)
		}
	}()

	initiatorSession, initiatorHandler := RunNewInitiator(port, t, session.LogonSettings{
		TargetCompID:  "Server",
		SenderCompID:  "Client",
		HeartBtInt:    heartBtInt,
		EncryptMethod: fixgen.EnumEncryptMethodNoneother,
	})

	waitHeartbeats := utils.TimedWaitGroup{}
	waitHeartbeats.Add(countOfHeartbeats)
	heartbeats := 4

	initiatorHandler.HandleIncoming(fixgen.MsgTypeHeartbeat, func(msg []byte) {
		if heartbeats <= 0 {
			return
		}
		heartbeats--
		waitHeartbeats.Done()
	})

	initiatorSession.OnChangeState(utils.EventLogon, func() bool {
		t.Log("client connected to server")
		return true
	})

	err := waitHeartbeats.WaitWithTimeout(time.Second * countOfHeartbeats * heartBtInt * 2)
	if err != nil {
		t.Fatalf("wait heartbeats: %s", err)
	}
}

func TestTestRequest(t *testing.T) {
	const (
		heartBtInt = 5
		testReqID  = "aloha"
		port       = 9992
	)

	// close acceptor after work
	acceptor := RunAcceptor(port, t, memory.NewStorage(100, 100))
	defer acceptor.Close()
	go func() {
		err := acceptor.ListenAndServe()
		if err != nil && !errors.Is(err, simplefixgo.ErrServerClosed) {
			panic(err)
		}
	}()

	initiatorSession, initiatorHandler := RunNewInitiator(port, t, session.LogonSettings{
		TargetCompID:  "Server",
		SenderCompID:  "Client",
		HeartBtInt:    heartBtInt,
		EncryptMethod: fixgen.EnumEncryptMethodNoneother,
	})

	waitHeartbeats := utils.TimedWaitGroup{}
	waitHeartbeats.Add(1)

	initiatorHandler.HandleIncoming(fixgen.MsgTypeHeartbeat, func(msg []byte) {
		heartbeatMsg, err := fixgen.ParseHeartbeat(msg)
		if err != nil {
			t.Fatalf("parse heartbeat: %s", err)
		}

		if heartbeatMsg.TestReqID() == testReqID {
			waitHeartbeats.Done()
		}
	})

	initiatorSession.OnChangeState(utils.EventLogon, func() bool {
		t.Log("client connected to server")
		t.Log("send test request")

		testRequestMsg := fixgen.TestRequest{}.New()
		testRequestMsg.SetFieldTestReqID(testReqID)

		err := initiatorSession.Send(testRequestMsg)
		if err != nil {
			panic(err)
		}

		return true
	})

	err := waitHeartbeats.WaitWithTimeout(time.Second * heartBtInt * 2)
	if err != nil {
		t.Fatalf("wait heartbeats: %s", err)
	}
}

func TestResendSequence(t *testing.T) {
	const (
		waitingResend       = time.Second * 6
		beforeResendRequest = time.Second * 3
		port                = 9993
		resendBegin         = 1
		resendEnd           = 3
	)

	var countOfResending = resendEnd - resendBegin + 1 // including

	// close acceptor after work
	acceptor := RunAcceptor(port, t, memory.NewStorage(100, 100))
	defer acceptor.Close()
	go func() {
		err := acceptor.ListenAndServe()
		if err != nil && !errors.Is(err, simplefixgo.ErrServerClosed) {
			panic(err)
		}
	}()

	initiatorSession, initiatorHandler := RunNewInitiator(port, t, session.LogonSettings{
		TargetCompID:  "Server",
		SenderCompID:  "Client",
		HeartBtInt:    1,
		EncryptMethod: fixgen.EnumEncryptMethodNoneother,
	})

	waitRepeats := utils.TimedWaitGroup{}
	waitRepeats.Add(countOfResending)
	messages := new(sync.Map)

	initiatorHandler.HandleIncoming(simplefixgo.AllMsgTypes, func(msg []byte) {
		msgSeqNumB, err := fix.ValueByTag(msg, fixgen.FieldMsgSeqNum)
		if err != nil {
			t.Fatalf("message sequence num parsing: %s", err)
		}

		msgSeqNum := string(msgSeqNumB)

		old, ok := messages.Load(msgSeqNum)
		if ok {
			if !bytes.Equal(old.([]byte), msg) {
				t.Log("> incoming", string(msg))
				t.Log("> saved", string(old.([]byte)))
				t.Fatalf("> different messages with same sequence number")
			} else {
				defer waitRepeats.Done()
			}
		} else {
			messages.Store(msgSeqNum, msg)
		}
	})

	initiatorSession.OnChangeState(utils.EventLogon, func() bool {
		t.Log("client connected to server")
		return true
	})

	time.Sleep(beforeResendRequest)
	err := initiatorSession.Send(fixgen.ResendRequest{}.New().SetFieldBeginSeqNo(resendBegin).SetFieldEndSeqNo(resendEnd))
	if err != nil {
		panic(err)
	}

	defer acceptor.Close()
	err = waitRepeats.WaitWithTimeout(waitingResend)
	if err != nil {
		t.Fatalf("wait heartbeats: %s", err)
	}
}

func TestCloseInitiatorConn(t *testing.T) {
	const (
		port = 9993
	)

	// close acceptor after work
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		t.Fatalf("listen error: %s", err)
	}

	waitClientClosed := make(chan struct{})
	handlerFactory := simplefixgo.NewAcceptorHandlerFactory(fixgen.FieldMsgType, 10)
	server := simplefixgo.NewAcceptor(listener, handlerFactory, func(handler simplefixgo.AcceptorHandler) {
		s, err := session.NewAcceptorSession(
			context.Background(),
			PseudoGeneratedOpts,
			handler,
			session.LogonSettings{HeartBtInt: 30, LogonTimeout: time.Second * 30},
			func(request session.LogonSettings) (err error) { return nil },
		)
		if err != nil {
			panic(err)
		}

		err = s.Run()
		if err != nil {
			t.Fatalf("run s: %s", err)
		}

		handler.OnDisconnect(func() bool {
			t.Log("client disconnected")
			waitClientClosed <- struct{}{}
			return true
		})
	})

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()

	conn, err := net.Dial("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		t.Fatalf("could not dial: %s", err)
	}

	handler := simplefixgo.NewInitiatorHandler(context.Background(), fixgen.FieldMsgType, 10)
	client := simplefixgo.NewInitiator(conn, handler, 10)

	s, err := session.NewInitiatorSession(
		context.Background(),
		handler,
		PseudoGeneratedOpts,
		session.LogonSettings{
			TargetCompID:  "Server",
			SenderCompID:  "Client",
			HeartBtInt:    1,
			EncryptMethod: fixgen.EnumEncryptMethodNoneother,
		},
	)
	if err != nil {
		panic(err)
	}

	go func() {
		err := client.Serve()
		if err != nil {
			panic(fmt.Errorf("serve client: %s", err))
		}
	}()

	err = s.Run()
	if err != nil {
		t.Fatalf("run session: %s", err)
	}

	client.Close()

	select {
	case <-waitClientClosed:
	case <-time.After(time.Second * 3):
		t.Fatalf("too long time waiting close")
	}
}

func TestCloseAcceptorConn(t *testing.T) {
	const (
		port = 9994
	)

	// close acceptor after work
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		t.Fatalf("listen error: %s", err)
	}

	waitServerDisconnect := make(chan struct{})
	handlerFactory := simplefixgo.NewAcceptorHandlerFactory(fixgen.FieldMsgType, 10)
	server := simplefixgo.NewAcceptor(listener, handlerFactory, func(handler simplefixgo.AcceptorHandler) {
		s, err := session.NewAcceptorSession(
			context.Background(),
			PseudoGeneratedOpts,
			handler,
			session.LogonSettings{HeartBtInt: 30, LogonTimeout: time.Second * 30},
			func(request session.LogonSettings) (err error) { return nil },
		)
		if err != nil {
			panic(err)
		}

		err = s.Run()
		if err != nil {
			t.Fatalf("run s: %s", err)
		}

		handler.OnConnect(func() bool {
			t.Log("server: client connected")
			return true
		})
	})

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()

	conn, err := net.Dial("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		t.Fatalf("could not dial: %s", err)
	}

	initiatorHandler := simplefixgo.NewInitiatorHandler(context.Background(), fixgen.FieldMsgType, 10)
	client := simplefixgo.NewInitiator(conn, initiatorHandler, 10)

	s, err := session.NewInitiatorSession(
		context.Background(),
		initiatorHandler,
		PseudoGeneratedOpts,
		session.LogonSettings{
			TargetCompID:  "Server",
			SenderCompID:  "Client",
			HeartBtInt:    1,
			EncryptMethod: fixgen.EnumEncryptMethodNoneother,
		},
	)
	if err != nil {
		panic(err)
	}

	initiatorHandler.OnConnect(func() bool {
		t.Log("client: connected to server")
		server.Close()

		return true
	})

	initiatorHandler.OnDisconnect(func() bool {
		t.Log("server disconnected")
		waitServerDisconnect <- struct{}{}
		return true
	})

	go func() {
		err := client.Serve()
		if !errors.Is(err, simplefixgo.ErrConnClosed) {
			panic(fmt.Errorf("serve client: %s", err))
		}
		t.Log("server closed")
	}()

	err = s.Run()
	if err != nil {
		t.Fatalf("run session: %s", err)
	}

	select {
	case <-waitServerDisconnect:
	case <-time.After(time.Second * 3):
		t.Fatalf("too long time waiting close")
	}
}
