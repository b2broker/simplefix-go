package messages

// HeartbeatBuilder is an interface for generated Heartbeat message builder
type HeartbeatBuilder interface {
	// Flow Message
	Parse(data []byte) (HeartbeatBuilder, error)
	New() HeartbeatBuilder

	// Generated Heartbeat Message
	TestReqID() string
	SetFieldTestReqID(string) HeartbeatBuilder

	// HeaderBuilder
	HeaderBuilder() HeaderBuilder

	// SendingMessage
	MsgType() string
	ToBytes() ([]byte, error)
}
