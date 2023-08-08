package messages

type SecurityListRequest interface {
	New() SecurityListRequestBuilder
	Build() SecurityListRequestBuilder
}

// SecurityListRequestBuilder is an interface providing functionality to a builder of auto-generated SecurityListRequest messages.
type SecurityListRequestBuilder interface {
	SecurityListRequest
	PipelineBuilder
}
