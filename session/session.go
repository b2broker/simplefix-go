package session

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	simplefixgo "github.com/b2broker/simplefix-go"
	"github.com/b2broker/simplefix-go/fix"
	"github.com/b2broker/simplefix-go/session/messages"
	"github.com/b2broker/simplefix-go/utils"
)

type LogonState int64

var (
	ErrMissingHandler          = errors.New("a handler is missing")
	ErrMissingRequiredTag      = errors.New("a required tag is missing in the tags list")
	ErrMissingEncryptedMethods = errors.New("a list of supported encryption methods is missing")
	ErrMissingErrorCodes       = errors.New("a list of error codes is missing")
	ErrMissingMessageBuilder   = errors.New("a required message builder is missing")
	ErrInvalidHeartBtLimits    = errors.New("an invalid heartbeat value exceeding the allowed limit")
	ErrInvalidHeartBtInt       = errors.New("an invalid integer value assigned to the heartbeat field")
	ErrInvalidLogonTimeout     = errors.New("the logon request timeout is too small")
	ErrMissingEncryptMethod    = errors.New("the encryption method is missing") // done
	ErrMissingLogonSettings    = errors.New("logon settings are missing") // done
	ErrMissingSessionOts       = errors.New("session options are missing")   // done
)

const (
	// WaitingLogon the connection has just started, waiting for a Session message or preparing to send it.
	WaitingLogon LogonState = iota

	// SuccessfulLogged session participants are authenticated, ready to work.
	SuccessfulLogged

	// WaitingLogonAnswer waiting for a response to a Logon message before starting the session.
	WaitingLogonAnswer

	// WaitingLogoutAnswer waiting for a response to a Logout message before terminating the session.
	WaitingLogoutAnswer

	// ReceivedLogoutAnswer a response to a Logout message was received.
	ReceivedLogoutAnswer
)

const (
	MinLogonTimeout = time.Millisecond
)

type logonHandler func(request *LogonSettings) (err error)

// todo
type IntLimits struct {
	Min int
	Max int
}

type Handler interface {
	HandleIncoming(msgType string, handle simplefixgo.IncomingHandlerFunc) (id int64)
	HandleOutgoing(msgType string, handle simplefixgo.OutgoingHandlerFunc) (id int64)
	RemoveIncomingHandler(msgType string, id int64) (err error)
	RemoveOutgoingHandler(msgType string, id int64) (err error)
	SendRaw(data []byte) error
	Send(message simplefixgo.SendingMessage) error
	Context() context.Context
}

// Session is a service that is used to maintain the default FIX API pipelines,
// which include the logon, logout and heartbeat messages, as well as rejects and message sequences.
type Session struct {
	*Opts
	side Side

	state   LogonState
	stateMu sync.RWMutex

	// Services:
	router Handler

	msgStorageAllHandler    int64
	msgStorageResendHandler int64

	counter      *int64
	eventHandler *utils.EventHandlerPool

	// Parameters:
	LogonHandler  logonHandler
	LogonSettings *LogonSettings

	// soon
	// maxMessageSize  int64  // validation
	// encryptedMethod string // validation

	ctx          context.Context
	cancel       context.CancelFunc
	errorHandler func(error)
	timeLocation *time.Location
	mu           sync.Mutex
}

// NewInitiatorSession returns a session for an Initiator object.
func NewInitiatorSession(handler Handler, opts *Opts, settings *LogonSettings) (s *Session, err error) {
	s, err = newSession(opts, handler, settings)
	if err != nil {
		return
	}

	if settings.HeartBtInt == 0 {
		return nil, ErrInvalidHeartBtInt
	}

	if settings.EncryptMethod == "" {
		return nil, ErrMissingEncryptMethod
	}

	s.side = sideInitiator
	s.changeState(WaitingLogonAnswer)

	return
}

// NewAcceptorSession returns a session for an Acceptor object.
func NewAcceptorSession(params *Opts, handler Handler, settings *LogonSettings, onLogon logonHandler) (s *Session, err error) {
	s, err = newSession(params, handler, settings)
	if err != nil {
		return
	}

	if params.AllowedEncryptedMethods == nil || len(params.AllowedEncryptedMethods) == 0 {
		return nil, ErrMissingEncryptedMethods
	}

	if settings.HeartBtLimits == nil || settings.HeartBtLimits.Min > settings.HeartBtLimits.Max ||
		settings.HeartBtLimits.Max == 0 || settings.HeartBtLimits.Min == 0 {
		return nil, ErrInvalidHeartBtLimits
	}

	if settings.LogonTimeout < MinLogonTimeout {
		return nil, ErrInvalidLogonTimeout
	}

	s.side = sideAcceptor
	s.changeState(WaitingLogon)
	s.LogonHandler = onLogon

	return
}

func newSession(opts *Opts, handler Handler, settings *LogonSettings) (session *Session, err error) {
	if handler == nil {
		return nil, ErrMissingHandler
	}

	if settings == nil {
		return nil, ErrMissingLogonSettings
	}

	err = opts.validate()
	if err != nil {
		return nil, err
	}

	session = &Session{
		Opts:         opts,
		router:       handler,
		counter:      new(int64),
		eventHandler: utils.NewEventHandlerPool(),

		LogonSettings: settings,
	}

	if opts.Location != "" {
		session.timeLocation, err = time.LoadLocation(opts.Location)
		if err != nil {
			return nil, err
		}
	} else {
		session.timeLocation = time.UTC
	}

	session.ctx, session.cancel = context.WithCancel(handler.Context())

	return session, nil
}

func (s *Session) changeState(state LogonState) {
	s.stateMu.Lock()
	s.state = state
	s.stateMu.Unlock()

	switch state {
	case SuccessfulLogged:
		s.eventHandler.Trigger(utils.EventLogon)
	case WaitingLogoutAnswer:
		s.eventHandler.Trigger(utils.EventRequest)
	case ReceivedLogoutAnswer:
		s.eventHandler.Trigger(utils.EventLogout)
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
		_ = s.router.RemoveIncomingHandler(s.MessageBuilders.ResendRequestBuilder.MsgType(), s.msgStorageResendHandler)
	}

	s.msgStorageAllHandler = s.router.HandleOutgoing(simplefixgo.AllMsgTypes, func(msg simplefixgo.SendingMessage) bool {
		_ = storage.Save(msg, msg.HeaderBuilder().MsgSeqNum())

		return true
	})
	s.msgStorageResendHandler = s.router.HandleIncoming(s.MessageBuilders.ResendRequestBuilder.MsgType(), func(data []byte) bool {
		resendMsg, err := s.MessageBuilders.ResendRequestBuilder.Parse(data)
		if err != nil {
			s.RejectMessage(data)
			return true
		}

		resendMessages, err := storage.Messages(resendMsg.BeginSeqNo(), resendMsg.EndSeqNo())
		if err != nil {
			return true
		}

		for _, message := range resendMessages {
			_ = s.router.Send(message)
		}

		return true
	})
}

func (s *Session) Logout() error {
	s.changeState(WaitingLogoutAnswer)

	s.sendWithErrorCheck(s.MessageBuilders.LogoutBuilder.New())

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

	msg := s.MessageBuilders.LogonBuilder.New().
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

// OnError is called when something goes wrong but the connection is still working.
// You can use it for handling errors that might occur as part of standard session procedures.
func (s *Session) OnError(handler func(error)) {
	s.errorHandler = handler
}

func (s *Session) Run() (err error) {
	s.changeState(WaitingLogon)
	if s.side == sideInitiator {
		err = s.LogonRequest()
		if err != nil {
			return fmt.Errorf("sendWithErrorCheck logon request: %w", err)
		}

		s.OnChangeState(utils.EventLogon, func() bool {
			_ = s.start()

			return true
		})
	}

	s.router.HandleIncoming(s.MessageBuilders.LogonBuilder.MsgType(), func(msg []byte) bool {
		incomingLogon, err := s.MessageBuilders.LogonBuilder.Parse(msg)
		if err != nil {
			s.RejectMessage(msg)
			return true
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
				s.MakeReject(s.SessionErrorCodes.Other, 0, incomingLogon.HeaderBuilder().MsgSeqNum())
				return true
			}

			err = s.start()
			if err != nil {
				s.MakeReject(s.SessionErrorCodes.IncorrectValue, s.Tags.HeartBtInt, incomingLogon.HeaderBuilder().MsgSeqNum())
				return true
			}

			answer := s.MessageBuilders.LogonBuilder.New()

			s.changeState(SuccessfulLogged)

			s.sendWithErrorCheck(answer)
			return true

		case WaitingLogonAnswer:
			s.changeState(SuccessfulLogged)

		case SuccessfulLogged:
			s.MakeReject(s.SessionErrorCodes.Other, 0, incomingLogon.HeaderBuilder().MsgSeqNum())
		}

		return true
	})
	s.router.HandleIncoming(s.MessageBuilders.LogoutBuilder.MsgType(), func(msg []byte) bool {
		_, err := s.MessageBuilders.LogoutBuilder.Parse(msg)
		if err != nil {
			s.RejectMessage(msg)
			return true
		}

		switch s.state {
		case WaitingLogoutAnswer:
			s.changeState(ReceivedLogoutAnswer)
			s.changeState(WaitingLogon)

		case SuccessfulLogged:
			s.changeState(WaitingLogoutAnswer)

			s.sendWithErrorCheck(s.MessageBuilders.LogoutBuilder.New())

		default:
			s.RejectMessage(msg)
		}

		if s.side == sideInitiator {
			s.changeState(WaitingLogonAnswer)
		} else {
			s.changeState(WaitingLogon)
		}

		return true
	})
	s.router.HandleIncoming(s.MessageBuilders.HeartbeatBuilder.MsgType(), func(msg []byte) bool {
		_, err := s.MessageBuilders.HeartbeatBuilder.Parse(msg)
		if err != nil {
			s.RejectMessage(msg)
			return true
		}

		if !s.IsLogged() {
			s.RejectMessage(msg)
			return true
		}

		return true
	})
	s.router.HandleIncoming(s.MessageBuilders.TestRequestBuilder.MsgType(), func(msg []byte) bool {
		incomingTestRequest, err := s.MessageBuilders.TestRequestBuilder.Parse(msg)
		if err != nil {
			s.RejectMessage(msg)
			return true
		}

		if !s.IsLogged() {
			s.RejectMessage(msg)
			return true
		}

		s.sendWithErrorCheck(s.MessageBuilders.HeartbeatBuilder.New().
			SetFieldTestReqID(incomingTestRequest.TestReqID()))

		return true
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

	s.router.HandleIncoming(simplefixgo.AllMsgTypes, func(msg []byte) bool {
		incomingMsgTimer.Refresh()

		return true
	})
	s.router.HandleOutgoing(simplefixgo.AllMsgTypes, func(msg simplefixgo.SendingMessage) bool {
		outgoingMsgTimer.Refresh()

		return true
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
			testRequest := s.MessageBuilders.TestRequestBuilder.New()

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

			heartbeat := s.MessageBuilders.HeartbeatBuilder.New()

			s.sendWithErrorCheck(heartbeat)
		}
	}()

	return nil
}

func (s *Session) RejectMessage(msg []byte) {
	reject := s.MakeReject(s.SessionErrorCodes.Other, 0, 0)

	seqNumB, err := fix.ValueByTag(msg, strconv.Itoa(s.Tags.MsgSeqNum))
	if err != nil {
		reject.SetFieldRefTagID(s.Tags.MsgSeqNum)
		s.sendWithErrorCheck(reject)
		return
	}

	seqNum, err := strconv.Atoi(string(seqNumB))
	if err != nil {
		reject.SetFieldSessionRejectReason(strconv.Itoa(5)) // An incorrect (out of range) value for this tag.
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

// Send is used to send a message after preparing its header tags:
// - the sequence number with a counter
// - the targetCompID and senderCompID fields
// - the sending time, with the current time zone indicated
// To send a message with custom fields, call the Send method for a Handler instead.
func (s *Session) Send(msg messages.Message) error {
	return s.send(msg)
}

func (s *Session) send(msg messages.Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()

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
	s.stateMu.RLock()
	defer s.stateMu.RUnlock()

	return s.state == SuccessfulLogged
}

func (s *Session) Context() context.Context {
	return s.ctx
}

func (s *Session) MakeReject(reasonCode, tag, seqNum int) messages.RejectBuilder {
	msg := s.MessageBuilders.RejectBuilder.New().
		SetFieldRefSeqNum(seqNum).
		SetFieldSessionRejectReason(strconv.Itoa(reasonCode))

	if tag != 0 {
		msg.SetFieldRefTagID(tag)
	}

	return msg
}

func (s *Session) Stop() (err error) {
	defer func() {
		s.eventHandler.Clean()
	}()

	err = s.Logout()
	if err != nil {
		return fmt.Errorf("sendWithErrorCheck logout request: %w", err)
	}

	delayTimer := time.AfterFunc(s.LogonSettings.CloseTimeout, func() {
		s.cancel()
	})

	s.OnChangeState(utils.EventLogout, func() bool {
		delayTimer.Stop()
		s.cancel()

		return true
	})

	return nil
}
