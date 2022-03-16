package messages

import "github.com/b2broker/simplefix-go/fix"

// RejectBuilder is an interface providing functionality to a builder of auto-generated Reject messages.
type RejectBuilder interface {
	New() RejectBuilder
	Items() fix.Items
	RefTagID() int
	SetFieldRefTagID(int) RejectBuilder
	RefSeqNum() int
	SetFieldRefSeqNum(int) RejectBuilder
	SessionRejectReason() string
	SetFieldSessionRejectReason(string) RejectBuilder
	HeaderBuilder() HeaderBuilder
	MsgType() string
	ToBytes() ([]byte, error)
}
