package messages

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
