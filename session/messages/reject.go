package messages

// RejectBuilder is an interface providing functionality to a builder of auto-generated Reject messages.
type RejectBuilder interface {
	// This message is required as a part of standard pipelines.
	Parse(data []byte) (RejectBuilder, error)
	New() RejectBuilder

	// An auto-generated Reject message.
	RefTagID() int
	SetFieldRefTagID(int) RejectBuilder
	RefSeqNum() int
	SetFieldRefSeqNum(int) RejectBuilder
	SessionRejectReason() string
	SetFieldSessionRejectReason(string) RejectBuilder

	// The builder of message headers.
	HeaderBuilder() HeaderBuilder

	// The code initiating sending of a message.
	MsgType() string
	ToBytes() ([]byte, error)
}
