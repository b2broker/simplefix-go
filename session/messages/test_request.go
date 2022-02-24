package messages

// TestRequestBuilder is an interface providing functionality to a builder of auto-generated TestRequest messages.
type TestRequestBuilder interface {
	// This message is required as a part of standard pipelines.
	Parse(data []byte) (TestRequestBuilder, error)
	New() TestRequestBuilder

	// An auto-generated TestRequest message.
	TestReqID() string
	SetFieldTestReqID(string) TestRequestBuilder

	// The builder of message headers.
	HeaderBuilder() HeaderBuilder

	// The code initiating sending of a message.
	MsgType() string
	ToBytes() ([]byte, error)
}
