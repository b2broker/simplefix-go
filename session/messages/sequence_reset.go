package messages

import "github.com/b2broker/simplefix-go/fix"

// SequenceResetBuilder is an interface providing functionality to a builder of auto-generated SequenceReset messages.
type SequenceResetBuilder interface {
	New() SequenceResetBuilder
	Items() fix.Items
	NewSeqNo() int
	SetFieldNewSeqNo(newSeqNo int) SequenceResetBuilder
	HeaderBuilder() HeaderBuilder
	MsgType() string
	ToBytes() ([]byte, error)
}
