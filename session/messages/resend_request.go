package messages

type ResendRequest interface {
	New() ResendRequestBuilder
	BeginSeqNo() int
	SetFieldBeginSeqNo(int) ResendRequestBuilder
	EndSeqNo() int
	SetFieldEndSeqNo(int) ResendRequestBuilder
}

// ResendRequestBuilder is an interface providing functionality to a builder of auto-generated ResendRequest messages.
type ResendRequestBuilder interface {
	ResendRequest
	PipelineBuilder
}
