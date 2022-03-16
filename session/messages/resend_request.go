package messages

import "github.com/b2broker/simplefix-go/fix"

// ResendRequestBuilder is an interface providing functionality to a builder of auto-generated ResendRequest messages.
type ResendRequestBuilder interface {
	New() ResendRequestBuilder
	Items() fix.Items
	BeginSeqNo() int
	SetFieldBeginSeqNo(int) ResendRequestBuilder
	EndSeqNo() int
	SetFieldEndSeqNo(int) ResendRequestBuilder
	HeaderBuilder() HeaderBuilder
	MsgType() string
	ToBytes() ([]byte, error)
}
