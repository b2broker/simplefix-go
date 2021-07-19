package tests

import (
	"bytes"
	simplefixgo "github.com/b2broker/simplefix-go"
	"github.com/b2broker/simplefix-go/fix"
	"github.com/b2broker/simplefix-go/session"
	"github.com/b2broker/simplefix-go/session/storages/memory"
	fixgen "github.com/b2broker/simplefix-go/tests/fix44"
	"github.com/b2broker/simplefix-go/utils"
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
	RunAcceptor(port, t, memory.NewStorage(100, 100))
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
	RunAcceptor(port, t, memory.NewStorage(100, 100))
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

		initiatorSession.Send(testRequestMsg)

		return true
	})

	err := waitHeartbeats.WaitWithTimeout(time.Second * heartBtInt * 2)
	if err != nil {
		t.Fatalf("wait heartbeats: %s", err)
	}
}

func TestResendSequence(t *testing.T) {
	const (
		waitingResend       = time.Second * 5
		beforeResendRequest = time.Second * 3
		port                = 9993
		resendBegin         = 1
		resendEnd           = 3
	)

	var countOfResending = resendEnd - resendBegin + 1 // including

	// close acceptor after work
	RunAcceptor(port, t, memory.NewStorage(100, 100))
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
	initiatorSession.Send(fixgen.ResendRequest{}.New().SetFieldBeginSeqNo(resendBegin).SetFieldEndSeqNo(resendEnd))

	err := waitRepeats.WaitWithTimeout(waitingResend)
	if err != nil {
		t.Fatalf("wait heartbeats: %s", err)
	}
}
