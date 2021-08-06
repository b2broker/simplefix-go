package messages

// LogoutBuilder is an interface for generated Logout message builder
type LogoutBuilder interface {

	// Flow Message
	Parse(data []byte) (LogoutBuilder, error)
	New() LogoutBuilder

	// HeaderBuilder
	HeaderBuilder() HeaderBuilder

	// SendingMessage
	MsgType() string
	ToBytes() ([]byte, error)
}
