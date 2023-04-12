package messages

type NewOrderSingle interface {
	New() NewOrderSingleBuilder
	Build() NewOrderSingleBuilder
}

// NewOrderSingleBuilder is an interface providing functionality to a builder of auto-generated NewOrderSingle messages.
type NewOrderSingleBuilder interface {
	NewOrderSingle
	PipelineBuilder
}
