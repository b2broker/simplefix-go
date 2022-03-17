package messages

type Logon interface {
	New() LogonBuilder
	EncryptMethod() string
	SetFieldEncryptMethod(string) LogonBuilder
	HeartBtInt() int
	SetFieldHeartBtInt(int) LogonBuilder

	Password() string
	SetFieldPassword(string) LogonBuilder
	Username() string
	SetFieldUsername(string) LogonBuilder
}

// LogonBuilder is an interface providing functionality to a builder of auto-generated Logon messages.
type LogonBuilder interface {
	Logon
	PipelineBuilder
}
