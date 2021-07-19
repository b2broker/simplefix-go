package messages

type LogonBuilder interface {
	// Flow Message
	Parse(data []byte) (LogonBuilder, error)
	New() LogonBuilder

	// Generated Logon Message
	EncryptMethod() string
	SetFieldEncryptMethod(string) LogonBuilder
	HeartBtInt() int
	SetFieldHeartBtInt(int) LogonBuilder

	Password() string
	SetFieldPassword(string) LogonBuilder
	Username() string
	SetFieldUsername(string) LogonBuilder

	// HeaderBuilder
	HeaderBuilder() HeaderBuilder

	// SendingMessage
	MsgType() string
	ToBytes() ([]byte, error)
}
