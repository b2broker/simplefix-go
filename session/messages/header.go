package messages

import "github.com/b2broker/simplefix-go/fix"

// ComponentConverter is an interface for Trailer message builder
type ComponentConverter interface {
	AsComponent() *fix.Component
}

// HeaderBuilder is an interface for Header message builder
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
