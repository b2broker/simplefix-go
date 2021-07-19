package session

import (
	"github.com/b2broker/simplefix-go/session/messages"
)

type MessageParser struct {
	SenderCompID  string
	TargetCompID  string
	MsgSeqCounter *int

	Reject messages.RejectBuilder
	Logon  messages.LogonBuilder

	Header  messages.HeaderBuilder
	Trailer messages.ComponentConverter
}

type MessageMaker interface {
	SetHeader(messages.HeaderBuilder)

	MakeReject(reasonCode, tag, seqNum int) messages.RejectBuilder
	RejectMessage(msg []byte) messages.RejectBuilder
}
