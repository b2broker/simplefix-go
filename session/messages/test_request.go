package messages

type TestRequestBuilder interface {
	// Flow Message
	Parse(data []byte) (TestRequestBuilder, error)
	New() TestRequestBuilder

	// Generated TestRequest Message
	TestReqID() string
	SetFieldTestReqID(string) TestRequestBuilder

	// HeaderBuilder
	HeaderBuilder() HeaderBuilder

	// SendingMessage
	MsgType() string
	ToBytes() ([]byte, error)
}
