package messages

import "github.com/b2broker/simplefix-go/fix"

// ComponentConverter is an interface providing functionality to a builder of trailer messages.
type ComponentConverter interface {
	AsComponent() *fix.Component
}

// HeaderBuilder is an interface providing functionality to a builder of header messages.
type HeaderBuilder interface {
	New() HeaderBuilder

	SenderCompID() string
	SetFieldSenderCompID(senderCompID string) HeaderBuilder
	TargetCompID() string
	SetFieldTargetCompID(targetCompID string) HeaderBuilder
	MsgSeqNum() int
	SetFieldMsgSeqNum(msgSeqNum int) HeaderBuilder
	SendingTime() string
	SetFieldSendingTime(string) HeaderBuilder

	AsComponent() *fix.Component
}
