package messages

type TestRequest interface {
	New() TestRequestBuilder
	SetFieldTestReqID(string) TestRequestBuilder
	TestReqID() string
}

// TestRequestBuilder is an interface providing functionality to a builder of auto-generated TestRequest messages.
type TestRequestBuilder interface {
	TestRequest
	PipelineBuilder
}
