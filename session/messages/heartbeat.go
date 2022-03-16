package messages

import "github.com/b2broker/simplefix-go/fix"

// HeartbeatBuilder is an interface providing functionality to a builder of auto-generated Heartbeat messages.
type HeartbeatBuilder interface {
	New() HeartbeatBuilder
	Items() fix.Items
	TestReqID() string
	SetFieldTestReqID(string) HeartbeatBuilder
	HeaderBuilder() HeaderBuilder
	MsgType() string
	ToBytes() ([]byte, error)
}
