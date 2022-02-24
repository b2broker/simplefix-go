package messages

// SequenceResetBuilder is an interface providing functionality to a builder of auto-generated SequenceReset messages.
type SequenceResetBuilder interface {
	// This message is required as a part of standard pipelines.
	Parse(data []byte) (SequenceResetBuilder, error)
	New() SequenceResetBuilder

	// An auto-generated ResendRequest message.
	NewSeqNo() int
	SetFieldNewSeqNo(newSeqNo int) SequenceResetBuilder

	// The builder of message headers.
	HeaderBuilder() HeaderBuilder

	// The code initiating sending of a message.
	MsgType() string
	ToBytes() ([]byte, error)
}
