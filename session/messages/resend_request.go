package messages

// ResendRequestBuilder is an interface providing functionality to a builder of auto-generated ResendRequest messages.
type ResendRequestBuilder interface {
	// This message is required as a part of standard pipelines.
	Parse(data []byte) (ResendRequestBuilder, error)
	New() ResendRequestBuilder

	// An auto-generated ResendRequest message.
	BeginSeqNo() int
	SetFieldBeginSeqNo(int) ResendRequestBuilder
	EndSeqNo() int
	SetFieldEndSeqNo(int) ResendRequestBuilder

	// The builder of message headers.
	HeaderBuilder() HeaderBuilder

	// The code initiating sending of a message.
	MsgType() string
	ToBytes() ([]byte, error)
}
