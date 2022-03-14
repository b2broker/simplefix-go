package messages

// HeartbeatBuilder is an interface providing functionality to a builder of auto-generated Heartbeat messages.
type HeartbeatBuilder interface {
	// This message is required as a part of standard pipelines.
	Parse(data []byte) (HeartbeatBuilder, error)
	New() HeartbeatBuilder

	// An auto-generated heartbeat message.
	TestReqID() string
	SetFieldTestReqID(string) HeartbeatBuilder

	// The builder of message headers.
	HeaderBuilder() HeaderBuilder

	// The code initiating sending of a message.
	MsgType() string
	ToBytes() ([]byte, error)
}
