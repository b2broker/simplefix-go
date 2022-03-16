package messages

import "github.com/b2broker/simplefix-go/fix"

// LogonBuilder is an interface providing functionality to a builder of auto-generated Logon messages.
type LogonBuilder interface {
	New() LogonBuilder
	Items() fix.Items

	EncryptMethod() string
	SetFieldEncryptMethod(string) LogonBuilder
	HeartBtInt() int
	SetFieldHeartBtInt(int) LogonBuilder

	Password() string
	SetFieldPassword(string) LogonBuilder
	Username() string
	SetFieldUsername(string) LogonBuilder

	HeaderBuilder() HeaderBuilder

	MsgType() string
	ToBytes() ([]byte, error)
}
