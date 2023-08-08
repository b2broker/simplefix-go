package messages

type SecurityList interface {
	New() SecurityListBuilder
	Build() SecurityListBuilder
}

// SecurityListBuilder is an interface providing functionality to a builder of auto-generated SecurityList messages.
type SecurityListBuilder interface {
	SecurityList
	PipelineBuilder
}
