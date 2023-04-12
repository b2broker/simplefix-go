package messages

type MarketDataRequest interface {
	New() MarketDataRequestBuilder
	Build() MarketDataRequestBuilder
}

// MarketDataRequestBuilder is an interface providing functionality to a builder of auto-generated MarketDataRequest messages.
type MarketDataRequestBuilder interface {
	MarketDataRequest
	PipelineBuilder
}
