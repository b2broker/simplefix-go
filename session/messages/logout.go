package messages

// LogoutBuilder is an interface providing functionality to a builder of auto-generated Logout messages.
type LogoutBuilder interface {

	// This message is required as a part of standard pipelines.
	Parse(data []byte) (LogoutBuilder, error)
	New() LogoutBuilder

	// The builder of message headers.
	HeaderBuilder() HeaderBuilder

	// The code initiating sending of a message.
	MsgType() string
	ToBytes() ([]byte, error)
}
