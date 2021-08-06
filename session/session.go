package session

import (
	"context"
	"errors"
	"fmt"
	simplefixgo "github.com/b2broker/simplefix-go"
	"github.com/b2broker/simplefix-go/fix"
	"github.com/b2broker/simplefix-go/session/messages"
	"github.com/b2broker/simplefix-go/utils"
	"math"
	"strconv"
	"sync/atomic"
	"time"
)

type LogonState int64

var (
	ErrNotLogon       = errors.New("logon message doesnt received")
	ErrTimeoutExpired = errors.New("logon timeout expired")
)

// todo constructor for acceptor and initiator
type LogonSettings struct {
	TargetCompID  string
	SenderCompID  string
	HeartBtInt    int
	EncryptMethod string
	Password      string
	Username      string
	LogonTimeout  time.Duration // todo
	HeartBtLimits *IntLimits
}

const (
	// WaitingLogon connection just started, waiting for Session message or preparing to send it
	WaitingLogon LogonState = iota

	// SuccessfulLogged participants are authenticated, ready to work
	SuccessfulLogged

	// WaitingLogonAnswer waiting for answer to Session
	WaitingLogonAnswer

	// WaitingLogoutAnswer waiting for answer to Logout
	WaitingLogoutAnswer
)

type logonHandler func(request *LogonSettings) (err error)

// todo
type IntLimits struct {
	Min int
	Max int
}

type Handler interface {
	HandleIncoming(msgType string, handle simplefixgo.HandlerFunc) (id int64)
	HandleOutgoing(msgType string, handle simplefixgo.HandlerFunc) (id int64)
	RemoveIncomingHandler(msgType string, id int64) (err error)
	RemoveOutgoingHandler(msgType string, id int64) (err error)
	SendRaw(msgType string, message []byte) error
	Send(message simplefixgo.SendingMessage) error
}

// Session is a service for working with default pipelines of FIX API
// logon, logout, heartbeats, rejects and message sequences
type Session struct {
	*Opts
	side  Side
	state LogonState

	// services
	router Handler

	msgStorageAllHandler    int64
	msgStorageResendHandler int64

	counter      *int64
	eventHandler *utils.EventHandlerPool

	// params
	LogonHandler  logonHandler
	LogonSettings *LogonSettings

	// soon
	// maxMessageSize  int64  // validation
	// encryptedMethod string // validation

	ctx          context.Context
	cancel       context.CancelFunc
	errorHandler func(error)
	timeLocation *time.Location
}

// NewInitiatorSession returns session for an Initiator
func NewInitiatorSession(ctx context.Context, router Handler, params *Opts,
	settings *LogonSettings) (s *Session, err error) {
	s, err = newSession(ctx, params, router, settings)
	if err != nil {
		return
	}

	s.side = sideInitiator
	s.state = WaitingLogonAnswer

	return
}

// NewAcceptorSession returns session for an Acceptor
func NewAcceptorSession(ctx context.Context, params *Opts, router Handler,
	settings *LogonSettings, onLogon logonHandler) (s *Session, err error) {
	s, err = newSession(ctx, params, router, settings)
	if err != nil {
		return
	}

	s.side = sideAcceptor
	s.state = WaitingLogon
	s.LogonHandler = onLogon

	return
}

func newSession(ctx context.Context, params *Opts, router Handler, settings *LogonSettings) (session *Session, err error) {
	session = &Session{
		Opts:         params,
		router:       router,
		counter:      new(int64),
		eventHandler: utils.NewEventHandlerPool(),

		LogonSettings: settings,
	}

	if params.Location != "" {
		session.timeLocation, err = time.LoadLocation(params.Location)
		if err != nil {
			return nil, err
		}
	} else {
		session.timeLocation = time.UTC
	}

	session.ctx, session.cancel = context.WithCancel(ctx)

	return session, nil
}

func (s *Session) changeState(state LogonState) {
	s.state = state

	switch s.state {
	case SuccessfulLogged:
		s.eventHandler.Trigger(utils.EventLogon)
	case WaitingLogoutAnswer:
		s.eventHandler.Trigger(utils.EventRequest)
	}
}

func (s *Session) checkLogonParams(incoming messages.LogonBuilder) (ok bool, tag, reasonCode int) {
	if _, ok := s.AllowedEncryptedMethods[incoming.EncryptMethod()]; !ok {
		return false, s.Tags.EncryptedMethod, s.SessionErrorCodes.IncorrectValue
	}

	if s.LogonSettings.HeartBtLimits == nil {
		return true, 0, 0
	}

	if incoming.HeartBtInt() > s.LogonSettings.HeartBtLimits.Min ||
		incoming.HeartBtInt() < s.LogonSettings.HeartBtLimits.Max {
		return false, s.Tags.HeartBtInt, s.SessionErrorCodes.IncorrectValue
	}

	return true, 0, 0
}

func (s *Session) SetMessageStorage(storage MessageStorage) {
	if s.msgStorageAllHandler > 0 || s.msgStorageResendHandler > 0 {
		_ = s.router.RemoveOutgoingHandler(simplefixgo.AllMsgTypes, s.msgStorageAllHandler)
		_ = s.router.RemoveIncomingHandler(s.ResendRequestBuilder.MsgType(), s.msgStorageResendHandler)
	}

	s.msgStorageAllHandler = s.router.HandleOutgoing(simplefixgo.AllMsgTypes, func(msg []byte) {
		value, _ := fix.ValueByTag(msg, strconv.Itoa(s.Tags.MsgSeqNum))
		id, _ := strconv.Atoi(string(value))

		_ = storage.Save(msg, id)
	})
	s.msgStorageResendHandler = s.router.HandleIncoming(s.ResendRequestBuilder.MsgType(), func(msg []byte) {
		resendMsg, err := s.ResendRequestBuilder.Parse(msg)
		if err != nil {
			s.RejectMessage(msg)
			return
		}

		resendMessages, err := storage.Messages(resendMsg.BeginSeqNo(), resendMsg.EndSeqNo())
		if err != nil {
			return
		}

		for _, message := range resendMessages {
			msgType, _ := fix.ValueByTag(message, strconv.Itoa(s.Tags.MsgType))
			_ = s.router.SendRaw(string(msgType), message)
		}
	})
}

func (s *Session) Logout() error {
	s.changeState(WaitingLogoutAnswer)

	s.sendWithErrorCheck(s.LogoutBuilder.New())

	return nil
}

func (s *Session) OnChangeState(event utils.Event, handle utils.EventHandlerFunc) {
	s.eventHandler.Handle(event, handle)
}

func (s *Session) StartWaiting() {
	s.changeState(WaitingLogon)
}

func (s *Session) LogonRequest() error {
	s.changeState(WaitingLogonAnswer)

	msg := s.LogonBuilder.New().
		SetFieldEncryptMethod(s.LogonSettings.EncryptMethod).
		SetFieldHeartBtInt(s.LogonSettings.HeartBtInt).
		SetFieldPassword(s.LogonSettings.Password).
		SetFieldUsername(s.LogonSettings.Username)

	s.sendWithErrorCheck(msg)
	return nil
}

func (s *Session) handlerError(err error) {
	if s.errorHandler != nil && err != nil {
		s.errorHandler(err)
	}
}

// OnError calls when something wrong, but connection is still working
// you can use it if you want to handler errors in standard process
func (s *Session) OnError(handler func(error)) {
	s.errorHandler = handler
}

func (s *Session) Run() (err error) {
	s.state = WaitingLogon
	if s.side == sideInitiator {
		err = s.LogonRequest()
		if err != nil {
			return fmt.Errorf("sendWithErrorCheck logon request: %w", err)
		}

		err = s.start()
		if err != nil {
			return fmt.Errorf("start heartbeat handler: %w", err)
		}
	}

	s.router.HandleIncoming(s.LogonBuilder.MsgType(), func(msg []byte) {
		incomingLogon, err := s.LogonBuilder.Parse(msg)
		if err != nil {
			s.RejectMessage(msg)
			return
		}

		switch s.state {
		case WaitingLogon:
			if ok, tag, reasonCode := s.checkLogonParams(incomingLogon); !ok {
				s.MakeReject(reasonCode, tag, incomingLogon.HeaderBuilder().MsgSeqNum())
			}

			s.LogonSettings = &LogonSettings{
				HeartBtInt:    incomingLogon.HeartBtInt(),
				EncryptMethod: incomingLogon.EncryptMethod(),
				Password:      incomingLogon.Password(),
				Username:      incomingLogon.Username(),
				TargetCompID:  incomingLogon.HeaderBuilder().TargetCompID(),
				SenderCompID:  incomingLogon.HeaderBuilder().SenderCompID(),
			}

			err := s.LogonHandler(s.LogonSettings)
			if err != nil {
				s.MakeReject(99, 0, incomingLogon.HeaderBuilder().MsgSeqNum())
			}

			err = s.start()
			if err != nil {
				s.MakeReject(99, s.Tags.HeartBtInt, incomingLogon.HeaderBuilder().MsgSeqNum())
				return
			}

			answer := s.LogonBuilder.New()

			s.state = SuccessfulLogged

			s.sendWithErrorCheck(answer)
			return

		case WaitingLogonAnswer:
			s.changeState(SuccessfulLogged)

		case SuccessfulLogged:
			s.MakeReject(99, 0, incomingLogon.HeaderBuilder().MsgSeqNum())
		}
	})
	s.router.HandleIncoming(s.LogoutBuilder.MsgType(), func(msg []byte) {
		_, err := s.LogoutBuilder.Parse(msg)
		if err != nil {
			s.RejectMessage(msg)
			return
		}

		switch s.state {
		case WaitingLogoutAnswer:
			s.changeState(WaitingLogon)

		case SuccessfulLogged:
			s.changeState(WaitingLogoutAnswer)

			s.sendWithErrorCheck(s.LogoutBuilder.New())

		default:
			s.RejectMessage(msg)
		}

		if s.side == sideInitiator {
			s.changeState(WaitingLogonAnswer)
		} else {
			s.changeState(WaitingLogon)
		}
	})
	s.router.HandleIncoming(s.HeartbeatBuilder.MsgType(), func(msg []byte) {
		_, err := s.HeartbeatBuilder.Parse(msg)
		if err != nil {
			s.RejectMessage(msg)
			return
		}

		if !s.IsLogged() {
			s.RejectMessage(msg)
			return
		}
	})
	s.router.HandleIncoming(s.TestRequestBuilder.MsgType(), func(msg []byte) {
		incomingTestRequest, err := s.TestRequestBuilder.Parse(msg)
		if err != nil {
			s.RejectMessage(msg)
			return
		}

		if !s.IsLogged() {
			s.RejectMessage(msg)
			return
		}

		s.sendWithErrorCheck(s.HeartbeatBuilder.New().
			SetFieldTestReqID(incomingTestRequest.TestReqID()))

	})

	return nil
}

func (s *Session) start() error {
	tolerance := int(math.Max(float64(s.LogonSettings.HeartBtInt/20), 1))
	incomingMsgTimer, err := utils.NewTimer(time.Second * time.Duration(s.LogonSettings.HeartBtInt+tolerance))
	if err != nil {
		return err
	}

	outgoingMsgTimer, err := utils.NewTimer(time.Second * time.Duration(s.LogonSettings.HeartBtInt))
	if err != nil {
		return err
	}

	s.router.HandleIncoming(simplefixgo.AllMsgTypes, func(msg []byte) {
		incomingMsgTimer.Refresh()
	})
	s.router.HandleOutgoing(simplefixgo.AllMsgTypes, func(msg []byte) {
		outgoingMsgTimer.Refresh()
	})

	go func() {
		defer incomingMsgTimer.Close()
		testReqCounter := 0
		for {
			select {
			case <-s.ctx.Done():
				return
			default:

			}

			incomingMsgTimer.TakeTimeout()
			testRequest := s.TestRequestBuilder.New()

			testReqCounter++
			expectedTestReq := strconv.Itoa(testReqCounter)
			testRequest.SetFieldTestReqID(expectedTestReq)

			s.sendWithErrorCheck(testRequest)
		}
	}()

	go func() {
		defer outgoingMsgTimer.Close()
		for {
			select {
			case <-s.ctx.Done():
				return
			default:

			}

			outgoingMsgTimer.TakeTimeout()

			heartbeat := s.HeartbeatBuilder.New()

			s.sendWithErrorCheck(heartbeat)
		}
	}()

	return nil
}

func (s *Session) RejectMessage(msg []byte) {
	reject := s.MakeReject(99, 0, 0)

	seqNumB, err := fix.ValueByTag(msg, strconv.Itoa(s.Tags.MsgSeqNum))
	if err != nil {
		reject.SetFieldRefTagID(s.Tags.MsgSeqNum)
		s.sendWithErrorCheck(reject)
		return
	}

	seqNum, err := strconv.Atoi(string(seqNumB))
	if err != nil {
		reject.SetFieldSessionRejectReason(strconv.Itoa(5)) // Value is incorrect (out of range) for this tag
		reject.SetFieldRefTagID(s.Tags.MsgSeqNum)
		s.sendWithErrorCheck(reject)
		return
	}

	reject.SetFieldRefSeqNum(seqNum)

	s.sendWithErrorCheck(reject)
}

func (s *Session) currentTime() time.Time {
	return time.Now().In(s.timeLocation)
}

// Send sends message with preparing header tags:
// - sequence number with counter
// - targetCompIDm senderCompID
// - sending time with current time zone
// if you want to send message with custom fields please use Send method at Handler
func (s *Session) Send(msg messages.Message) error {
	return s.send(msg)
}

func (s *Session) send(msg messages.Message) error {
	msg.HeaderBuilder().
		SetFieldMsgSeqNum(int(atomic.AddInt64(s.counter, 1))).
		SetFieldTargetCompID(s.LogonSettings.TargetCompID).
		SetFieldSenderCompID(s.LogonSettings.SenderCompID).
		SetFieldSendingTime(s.currentTime().Format(fix.TimeLayout))

	return s.router.Send(msg)
}

func (s *Session) sendWithErrorCheck(msg messages.Message) {
	s.handlerError(s.send(msg))
}

func (s *Session) IsLogged() bool {
	return s.state == SuccessfulLogged
}

func (s *Session) MakeReject(reasonCode, tag, seqNum int) messages.RejectBuilder {
	msg := s.RejectBuilder.New().
		SetFieldRefSeqNum(seqNum).
		SetFieldSessionRejectReason(strconv.Itoa(reasonCode))

	if tag != 0 {
		msg.SetFieldRefTagID(tag)
	}

	return msg
}
