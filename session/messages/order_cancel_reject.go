package messages

type OrderCancelReject interface {
	New() OrderCancelRejectBuilder
	Build() OrderCancelRejectBuilder
}

// OrderCancelRejectBuilder is an interface providing functionality to a builder of auto-generated OrderCancelReject messages.
type OrderCancelRejectBuilder interface {
	OrderCancelReject
	PipelineBuilder
}
