package messages

// LogonBuilder is an interface providing functionality to a builder of auto-generated Logon messages.
type LogonBuilder interface {
	// This message is required as a part of standard pipelines.
	Parse(data []byte) (LogonBuilder, error)
	New() LogonBuilder

	// An auto-generated Logon message.
	EncryptMethod() string
	SetFieldEncryptMethod(string) LogonBuilder
	HeartBtInt() int
	SetFieldHeartBtInt(int) LogonBuilder

	Password() string
	SetFieldPassword(string) LogonBuilder
	Username() string
	SetFieldUsername(string) LogonBuilder

	// The builder of message headers.
	HeaderBuilder() HeaderBuilder

	// The code initiating sending of a message.
	MsgType() string
	ToBytes() ([]byte, error)
}
