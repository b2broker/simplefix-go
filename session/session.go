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

	"github.com/b2broker/simplefix-go/fix/encoding"

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
	ErrMissingLogonSettings    = errors.New("logon settings are missing")       // done
	ErrMissingSessionOts       = errors.New("session options are missing")      // done
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

	// WaitingTestReqAnswer waiting for a response to a Test request before disconnect
	WaitingTestReqAnswer

	// Disconnect need to disconnect session
	Disconnect
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
	SendBuffered(message simplefixgo.SendingMessage) error
	SendBatch(messages []simplefixgo.SendingMessage) error
	Context() context.Context
	Stop()
}

// Session is a service that is used to maintain the default FIX API pipelines,
// which include the logon, logout and heartbeat messages, as well as rejects and message sequences.
type Session struct {
	*Opts
	side Side

	state atomic.Int64

	// Services:
	Router       Handler
	unmarshaller Unmarshaller

	messageStorage MessageStorage
	counter        CounterStorage
	eventHandler   *utils.EventHandlerPool

	// Parameters:
	LogonHandler  logonHandler
	LogonSettings *LogonSettings
	logonRequest  func(*Session) error

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
func NewInitiatorSession(handler Handler, opts *Opts, settings *LogonSettings, cs CounterStorage, ms MessageStorage) (s *Session, err error) {
	s, err = newSession(opts, handler, settings, cs, ms)
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
	s.changeState(WaitingLogonAnswer, true)

	return
}

// NewAcceptorSession returns a session for an Acceptor object.
func NewAcceptorSession(params *Opts, handler Handler, settings *LogonSettings, onLogon logonHandler, cs CounterStorage, ms MessageStorage) (s *Session, err error) {
	s, err = newSession(params, handler, settings, cs, ms)
	if err != nil {
		return
	}

	if len(params.AllowedEncryptedMethods) == 0 {
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
	s.changeState(WaitingLogon, true)
	s.LogonHandler = onLogon

	return
}

func newSession(opts *Opts, handler Handler, settings *LogonSettings, cs CounterStorage, ms MessageStorage) (session *Session, err error) {
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
		Opts:           opts,
		Router:         handler,
		messageStorage: ms,
		counter:        cs,
		eventHandler:   utils.NewEventHandlerPool(),
		unmarshaller:   encoding.NewDefaultUnmarshaller(true),

		LogonSettings: settings,
	}

	session.setStorageCallbacks()

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

func (s *Session) changeState(state LogonState, isEventTriggerRequired bool) {
	s.state.Store(int64(state))

	if !isEventTriggerRequired {
		return
	}

	switch state {
	case SuccessfulLogged:
		s.eventHandler.Trigger(utils.EventLogon)
	case WaitingLogoutAnswer:
		s.eventHandler.Trigger(utils.EventRequest)
	case ReceivedLogoutAnswer:
		s.eventHandler.Trigger(utils.EventLogout)
	case Disconnect:
		s.eventHandler.Trigger(utils.EventDisconnect)
	}
}

func (s *Session) checkLogonParams(incoming messages.LogonBuilder) (ok bool, tag, reasonCode int) {
	if _, ok := s.AllowedEncryptedMethods[incoming.EncryptMethod()]; !ok {
		return false, s.Tags.EncryptedMethod, s.SessionErrorCodes.IncorrectValue
	}

	if s.LogonSettings.HeartBtLimits == nil {
		return true, 0, 0
	}

	if incoming.HeartBtInt() < s.LogonSettings.HeartBtLimits.Min ||
		incoming.HeartBtInt() > s.LogonSettings.HeartBtLimits.Max {
		return false, s.Tags.HeartBtInt, s.SessionErrorCodes.IncorrectValue
	}

	return true, 0, 0
}

func (s *Session) setStorageCallbacks() {
	s.Router.HandleOutgoing(simplefixgo.AllMsgTypes, func(msg simplefixgo.SendingMessage) bool {
		err := s.messageStorage.Save(fix.StorageID{
			Sender: s.LogonSettings.SenderCompID,
			Target: s.LogonSettings.TargetCompID,
			Side:   fix.Outgoing,
		}, msg, msg.HeaderBuilder().MsgSeqNum())
		return err == nil
	})

	s.Router.HandleIncoming(simplefixgo.AllMsgTypes, func(msg []byte) bool {
		if state := s.State(); state != WaitingLogonAnswer && state != WaitingLogon {
			seqNum, err := fix.ValueByTag(msg, strconv.Itoa(s.Tags.MsgSeqNum))
			if err != nil {
				return true
			}
			seqNumInt, err := strconv.Atoi(string(seqNum))
			if err != nil {
				return true
			}

			msgType, err := fix.ValueByTag(msg, strconv.Itoa(s.Tags.MsgType))
			if err != nil {
				return true
			}
			msgTypeStr := string(msgType)
			if s.MessageBuilders.SequenceResetBuilder == nil || msgTypeStr != s.MessageBuilders.SequenceResetBuilder.MsgType() {
				err = s.counter.SetSeqNum(fix.StorageID{
					Sender: s.LogonSettings.SenderCompID,
					Target: s.LogonSettings.TargetCompID,
					Side:   fix.Incoming,
				}, seqNumInt)
			}
			return err == nil
		}

		return true
	})

	s.Router.HandleIncoming(s.MessageBuilders.ResendRequestBuilder.MsgType(), func(data []byte) bool {
		resendMsg := s.MessageBuilders.ResendRequestBuilder.New()
		err := s.unmarshaller.Unmarshal(resendMsg, data)
		if err != nil {
			s.RejectMessage(data)
			return true
		}

		resendMessages, err := s.messageStorage.Messages(fix.StorageID{
			Sender: s.LogonSettings.SenderCompID,
			Target: s.LogonSettings.TargetCompID,
			Side:   fix.Outgoing,
		}, resendMsg.BeginSeqNo(), resendMsg.EndSeqNo())
		if err != nil {
			return true
		}

		_ = s.Router.SendBatch(resendMessages)

		return true
	})
}

func (s *Session) SetLogonRequest(logonRequest func(*Session) error) {
	s.logonRequest = logonRequest
}

func (s *Session) Logout() error {
	s.changeState(WaitingLogoutAnswer, true)

	s.sendWithErrorCheck(s.MessageBuilders.LogoutBuilder.Build())

	return nil
}

func (s *Session) OnChangeState(event utils.Event, handle utils.EventHandlerFunc) {
	s.eventHandler.Handle(event, handle)
}

func (s *Session) StartWaiting() {
	s.changeState(WaitingLogon, true)
}

func (s *Session) LogonRequest() error {
	s.changeState(WaitingLogonAnswer, true)
	if s.logonRequest != nil {
		return s.logonRequest(s)
	}

	msg := s.MessageBuilders.LogonBuilder.Build().
		SetFieldEncryptMethod(s.LogonSettings.EncryptMethod).
		SetFieldHeartBtInt(s.LogonSettings.HeartBtInt).
		SetFieldPassword(s.LogonSettings.Password).
		SetFieldUsername(s.LogonSettings.Username)

	s.sendWithErrorCheck(msg)
	return nil
}

func (s *Session) HandlerError(err error) {
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
	s.changeState(WaitingLogon, true)
	s.OnChangeState(utils.EventDisconnect, func() bool {
		s.cancel()
		s.Router.Stop()
		return true
	})
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

	s.Router.HandleIncoming(s.MessageBuilders.LogonBuilder.MsgType(), func(data []byte) bool {

		incomingLogon := s.MessageBuilders.LogonBuilder.New()
		err := s.unmarshaller.Unmarshal(incomingLogon, data)
		if err != nil {
			s.RejectMessage(data)
			return true
		}

		switch s.State() {
		case WaitingLogon:
			s.LogonSettings = &LogonSettings{
				HeartBtInt:      incomingLogon.HeartBtInt(),
				EncryptMethod:   incomingLogon.EncryptMethod(),
				Password:        incomingLogon.Password(),
				Username:        incomingLogon.Username(),
				ResetSeqNumFlag: incomingLogon.ResetSeqNumFlag(),
				TargetCompID:    incomingLogon.HeaderBuilder().TargetCompID(),
				SenderCompID:    incomingLogon.HeaderBuilder().SenderCompID(),
				LogonTimeout:    s.LogonSettings.LogonTimeout,
				CloseTimeout:    s.LogonSettings.CloseTimeout,
				HeartBtLimits:   s.LogonSettings.HeartBtLimits,
			}

			if s.side == sideAcceptor {
				s.LogonSettings.TargetCompID, s.LogonSettings.SenderCompID = s.LogonSettings.SenderCompID, s.LogonSettings.TargetCompID
			}

			if ok, tag, reasonCode := s.checkLogonParams(incomingLogon); !ok {
				s.sendWithErrorCheck(s.MakeReject(reasonCode, tag, incomingLogon.HeaderBuilder().MsgSeqNum()))
				return true
			}

			err := s.LogonHandler(s.LogonSettings)
			if err != nil {
				s.sendWithErrorCheck(s.MakeReject(s.SessionErrorCodes.Other, 0, incomingLogon.HeaderBuilder().MsgSeqNum()))
				return true
			}

			err = s.start()
			if err != nil {
				s.sendWithErrorCheck(s.MakeReject(s.SessionErrorCodes.IncorrectValue, s.Tags.HeartBtInt, incomingLogon.HeaderBuilder().MsgSeqNum()))
				return true
			}

			answer := s.MessageBuilders.LogonBuilder.Build()
			answer.SetFieldEncryptMethod(s.LogonSettings.EncryptMethod).SetFieldHeartBtInt(s.LogonSettings.HeartBtInt)

			s.changeState(SuccessfulLogged, true)

			s.sendWithErrorCheck(answer)

			s.processIncSeq(incomingLogon)

		case WaitingLogonAnswer:
			s.changeState(SuccessfulLogged, true)

			s.processIncSeq(incomingLogon)
		case SuccessfulLogged:
			s.sendWithErrorCheck(s.MakeReject(s.SessionErrorCodes.Other, 0, incomingLogon.HeaderBuilder().MsgSeqNum()))
		}

		return true
	})
	s.Router.HandleIncoming(s.MessageBuilders.LogoutBuilder.MsgType(), func(data []byte) bool {
		err := s.unmarshaller.Unmarshal(s.MessageBuilders.LogoutBuilder.New(), data)
		if err != nil {
			s.RejectMessage(data)
			return true
		}

		switch s.State() {
		case WaitingLogoutAnswer:
			s.changeState(ReceivedLogoutAnswer, true)
			s.changeState(WaitingLogon, true)

		case SuccessfulLogged:
			s.changeState(WaitingLogoutAnswer, true)

			s.sendWithErrorCheck(s.MessageBuilders.LogoutBuilder.Build())

		default:
			s.RejectMessage(data)
		}

		if s.side == sideInitiator {
			s.changeState(WaitingLogonAnswer, true)
		} else {
			s.changeState(WaitingLogon, true)
		}

		return true
	})
	s.Router.HandleIncoming(s.MessageBuilders.HeartbeatBuilder.MsgType(), func(data []byte) bool {
		heartbeat := s.MessageBuilders.HeartbeatBuilder.New()
		err := s.unmarshaller.Unmarshal(heartbeat, data)
		if err != nil {
			s.RejectMessage(data)
			return true
		}

		if !s.IsLogged() {
			s.RejectMessage(data)
			return true
		}

		if s.State() == WaitingTestReqAnswer {
			// reset SuccessfulLogged statue without event trigger
			s.changeState(SuccessfulLogged, false)
		}

		return true
	})
	s.Router.HandleIncoming(s.MessageBuilders.TestRequestBuilder.MsgType(), func(data []byte) bool {
		testRequest := s.MessageBuilders.TestRequestBuilder.New()
		err := s.unmarshaller.Unmarshal(testRequest, data)
		if err != nil {
			s.RejectMessage(data)
			return true
		}

		if !s.IsLogged() {
			s.RejectMessage(data)
			return true
		}

		s.sendWithErrorCheck(s.MessageBuilders.HeartbeatBuilder.Build().
			SetFieldTestReqID(testRequest.TestReqID()))

		return true
	})

	return nil
}

func (s *Session) processIncSeq(incomingLogon messages.LogonBuilder) {
	incSeqNum := incomingLogon.HeaderBuilder().MsgSeqNum()
	currSeqNum, err := s.counter.GetCurrSeqNum(fix.StorageID{
		Sender: s.LogonSettings.SenderCompID,
		Target: s.LogonSettings.TargetCompID,
		Side:   fix.Incoming,
	})
	if err != nil {
		return
	}

	if currSeqNum+1 < incSeqNum {
		resendMsg := s.MessageBuilders.ResendRequestBuilder.New()
		resendMsg.SetFieldBeginSeqNo(currSeqNum + 1)
		resendMsg.SetFieldEndSeqNo(incSeqNum - 1)
		s.sendWithErrorCheck(resendMsg)
	}

	_ = s.counter.SetSeqNum(fix.StorageID{
		Sender: s.LogonSettings.SenderCompID,
		Target: s.LogonSettings.TargetCompID,
		Side:   fix.Incoming,
	}, incSeqNum)
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

	s.Router.HandleIncoming(simplefixgo.AllMsgTypes, func(msg []byte) bool {
		incomingMsgTimer.Refresh()
		if s.State() == WaitingTestReqAnswer {
			s.changeState(SuccessfulLogged, false)
		}

		return true
	})
	s.Router.HandleOutgoing(simplefixgo.AllMsgTypes, func(msg simplefixgo.SendingMessage) bool {
		outgoingMsgTimer.Refresh()

		return true
	})

	go func() {
		defer incomingMsgTimer.Close()
		testReqCounter := 0
		for {
			incomingMsgTimer.TakeTimeout()
			select {
			case <-s.ctx.Done():
				return
			default:
			}

			if s.State() == WaitingTestReqAnswer {
				s.changeState(Disconnect, true)
				return
			}

			testRequest := s.MessageBuilders.TestRequestBuilder.Build()

			testReqCounter++
			expectedTestReq := strconv.Itoa(testReqCounter)
			testRequest.SetFieldTestReqID(expectedTestReq)

			s.changeState(WaitingTestReqAnswer, true)
			s.sendWithErrorCheck(testRequest)
		}
	}()

	go func() {
		defer outgoingMsgTimer.Close()
		for {
			outgoingMsgTimer.TakeTimeout()
			select {
			case <-s.ctx.Done():
				return
			default:
			}

			heartbeat := s.MessageBuilders.HeartbeatBuilder.Build()

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

func (s *Session) CurrentTime() time.Time {
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

	nextSeqNum, err := s.counter.GetNextSeqNum(fix.StorageID{
		Sender: s.LogonSettings.SenderCompID,
		Target: s.LogonSettings.TargetCompID,
		Side:   fix.Outgoing,
	})
	if err != nil {
		return err
	}
	msg.HeaderBuilder().
		SetFieldMsgSeqNum(nextSeqNum).
		SetFieldTargetCompID(s.LogonSettings.TargetCompID).
		SetFieldSenderCompID(s.LogonSettings.SenderCompID).
		SetFieldSendingTime(s.CurrentTime().Format(fix.TimeLayout))

	return s.Router.Send(msg)
}

func (s *Session) SendBuffered(msg messages.Message) error {
	return s.sendBuffered(msg)
}

func (s *Session) sendBuffered(msg messages.Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	nextSeqNum, err := s.counter.GetNextSeqNum(fix.StorageID{
		Sender: s.LogonSettings.SenderCompID,
		Target: s.LogonSettings.TargetCompID,
		Side:   fix.Outgoing,
	})
	if err != nil {
		return err
	}
	msg.HeaderBuilder().
		SetFieldMsgSeqNum(nextSeqNum).
		SetFieldTargetCompID(s.LogonSettings.TargetCompID).
		SetFieldSenderCompID(s.LogonSettings.SenderCompID).
		SetFieldSendingTime(s.CurrentTime().Format(fix.TimeLayout))

	return s.Router.SendBuffered(msg)
}

func (s *Session) sendWithErrorCheck(msg messages.Message) {
	s.HandlerError(s.send(msg))
}

func (s *Session) IsLogged() bool {
	return s.State() == SuccessfulLogged
}

func (s *Session) Context() context.Context {
	return s.ctx
}

func (s *Session) MakeReject(reasonCode, tag, seqNum int) messages.RejectBuilder {
	msg := s.MessageBuilders.RejectBuilder.Build().
		SetFieldRefSeqNum(seqNum).
		SetFieldSessionRejectReason(strconv.Itoa(reasonCode))

	if tag != 0 {
		msg.SetFieldRefTagID(tag)
	}

	return msg
}

// SetUnmarshaller replaces current unmarshaller buy custom one
// It could be called only before starting Session
func (s *Session) SetUnmarshaller(unmarshaller Unmarshaller) {
	s.unmarshaller = unmarshaller
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

func (s *Session) State() LogonState {
	return LogonState(s.state.Load())
}
