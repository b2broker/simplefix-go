package messages

import "github.com/b2broker/simplefix-go/fix"

// TestRequestBuilder is an interface providing functionality to a builder of auto-generated TestRequest messages.
type TestRequestBuilder interface {
	New() TestRequestBuilder
	Items() fix.Items
	TestReqID() string
	SetFieldTestReqID(string) TestRequestBuilder
	HeaderBuilder() HeaderBuilder
	MsgType() string
	ToBytes() ([]byte, error)
}
