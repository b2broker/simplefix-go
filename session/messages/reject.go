package messages

type RejectBuilder interface {
	// Flow Message
	Parse(data []byte) (RejectBuilder, error)
	New() RejectBuilder

	// Generated Reject Message
	RefTagID() int
	SetFieldRefTagID(int) RejectBuilder
	RefSeqNum() int
	SetFieldRefSeqNum(int) RejectBuilder
	SessionRejectReason() string
	SetFieldSessionRejectReason(string) RejectBuilder

	// HeaderBuilder
	HeaderBuilder() HeaderBuilder

	// SendingMessage
	MsgType() string
	ToBytes() ([]byte, error)
}
