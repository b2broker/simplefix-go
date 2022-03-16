package messages

import "github.com/b2broker/simplefix-go/fix"

// LogoutBuilder is an interface providing functionality to a builder of auto-generated Logout messages.
type LogoutBuilder interface {
	New() LogoutBuilder
	Items() fix.Items
	HeaderBuilder() HeaderBuilder
	MsgType() string
	ToBytes() ([]byte, error)
}
