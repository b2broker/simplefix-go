package fix44

import (
	"github.com/b2broker/simplefix-go/fix"
	"github.com/b2broker/simplefix-go/session/messages"
)

const MsgTypeMarketDataRequest = "V"

type MarketDataRequest struct {
	*fix.Message
}

func makeMarketDataRequest() *MarketDataRequest {
	msg := &MarketDataRequest{
		Message: fix.NewMessage(FieldBeginString, FieldBodyLength, FieldCheckSum, FieldMsgType, beginString, MsgTypeMarketDataRequest).
			SetBody(
				fix.NewKeyValue(FieldMDReqID, &fix.String{}),
				fix.NewKeyValue(FieldSubscriptionRequestType, &fix.String{}),
				fix.NewKeyValue(FieldMarketDepth, &fix.Int{}),
				fix.NewKeyValue(FieldMDUpdateType, &fix.String{}),
				fix.NewKeyValue(FieldAggregatedBook, &fix.Bool{}),
				fix.NewKeyValue(FieldOpenCloseSettlFlag, &fix.String{}),
				fix.NewKeyValue(FieldScope, &fix.String{}),
				fix.NewKeyValue(FieldMDImplicitDelete, &fix.String{}),
				NewMDEntryTypesGrp().Group,
				NewRelatedSymGrp().Group,
			),
	}

	msg.SetHeader(makeHeader().AsComponent())
	msg.SetTrailer(makeTrailer().AsComponent())

	return msg
}

func NewMarketDataRequest(header *Header, trailer *Trailer, mDReqID string, subscriptionRequestType string, marketDepth int, noMDEntryTypes *MDEntryTypesGrp, noRelatedSym *RelatedSymGrp) *MarketDataRequest {
	msg := makeMarketDataRequest().
		SetMDReqID(mDReqID).
		SetSubscriptionRequestType(subscriptionRequestType).
		SetMarketDepth(marketDepth).
		SetMDEntryTypesGrp(noMDEntryTypes).
		SetRelatedSymGrp(noRelatedSym)
	msg.SetHeader(header.AsComponent())
	msg.SetTrailer(trailer.AsComponent())
	return msg
}

func ParseMarketDataRequest(data []byte) (*MarketDataRequest, error) {
	msg := fix.NewMessage(FieldBeginString, FieldBodyLength, FieldCheckSum, FieldMsgType, FieldBeginString, beginString).
		SetBody(makeMarketDataRequest().Body()...).
		SetHeader(makeHeader().AsComponent()).
		SetTrailer(makeTrailer().AsComponent())

	if err := msg.Unmarshal(data); err != nil {
		return nil, err
	}

	return &MarketDataRequest{
		Message: msg,
	}, nil
}

func (marketDataRequest *MarketDataRequest) Header() *Header {
	header := marketDataRequest.Message.Header()

	return &Header{header}
}

func (marketDataRequest *MarketDataRequest) HeaderBuilder() messages.HeaderBuilder {
	return marketDataRequest.Header()
}

func (marketDataRequest *MarketDataRequest) Trailer() *Trailer {
	trailer := marketDataRequest.Message.Trailer()

	return &Trailer{trailer}
}

func (marketDataRequest *MarketDataRequest) MDReqID() string {
	kv := marketDataRequest.Get(0)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (marketDataRequest *MarketDataRequest) SetMDReqID(mDReqID string) *MarketDataRequest {
	kv := marketDataRequest.Get(0).(*fix.KeyValue)
	_ = kv.Load().Set(mDReqID)
	return marketDataRequest
}

func (marketDataRequest *MarketDataRequest) SubscriptionRequestType() string {
	kv := marketDataRequest.Get(1)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (marketDataRequest *MarketDataRequest) SetSubscriptionRequestType(subscriptionRequestType string) *MarketDataRequest {
	kv := marketDataRequest.Get(1).(*fix.KeyValue)
	_ = kv.Load().Set(subscriptionRequestType)
	return marketDataRequest
}

func (marketDataRequest *MarketDataRequest) MarketDepth() int {
	kv := marketDataRequest.Get(2)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(int)
}

func (marketDataRequest *MarketDataRequest) SetMarketDepth(marketDepth int) *MarketDataRequest {
	kv := marketDataRequest.Get(2).(*fix.KeyValue)
	_ = kv.Load().Set(marketDepth)
	return marketDataRequest
}

func (marketDataRequest *MarketDataRequest) MDUpdateType() string {
	kv := marketDataRequest.Get(3)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (marketDataRequest *MarketDataRequest) SetMDUpdateType(mDUpdateType string) *MarketDataRequest {
	kv := marketDataRequest.Get(3).(*fix.KeyValue)
	_ = kv.Load().Set(mDUpdateType)
	return marketDataRequest
}

func (marketDataRequest *MarketDataRequest) AggregatedBook() bool {
	kv := marketDataRequest.Get(4)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(bool)
}

func (marketDataRequest *MarketDataRequest) SetAggregatedBook(aggregatedBook bool) *MarketDataRequest {
	kv := marketDataRequest.Get(4).(*fix.KeyValue)
	_ = kv.Load().Set(aggregatedBook)
	return marketDataRequest
}

func (marketDataRequest *MarketDataRequest) OpenCloseSettlFlag() string {
	kv := marketDataRequest.Get(5)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (marketDataRequest *MarketDataRequest) SetOpenCloseSettlFlag(openCloseSettlFlag string) *MarketDataRequest {
	kv := marketDataRequest.Get(5).(*fix.KeyValue)
	_ = kv.Load().Set(openCloseSettlFlag)
	return marketDataRequest
}

func (marketDataRequest *MarketDataRequest) Scope() string {
	kv := marketDataRequest.Get(6)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (marketDataRequest *MarketDataRequest) SetScope(scope string) *MarketDataRequest {
	kv := marketDataRequest.Get(6).(*fix.KeyValue)
	_ = kv.Load().Set(scope)
	return marketDataRequest
}

func (marketDataRequest *MarketDataRequest) MDImplicitDelete() string {
	kv := marketDataRequest.Get(7)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (marketDataRequest *MarketDataRequest) SetMDImplicitDelete(mDImplicitDelete string) *MarketDataRequest {
	kv := marketDataRequest.Get(7).(*fix.KeyValue)
	_ = kv.Load().Set(mDImplicitDelete)
	return marketDataRequest
}

func (marketDataRequest *MarketDataRequest) MDEntryTypesGrp() *MDEntryTypesGrp {
	group := marketDataRequest.Get(8).(*fix.Group)

	return &MDEntryTypesGrp{group}
}

func (marketDataRequest *MarketDataRequest) SetMDEntryTypesGrp(noMDEntryTypes *MDEntryTypesGrp) *MarketDataRequest {
	marketDataRequest.Set(8, noMDEntryTypes.Group)

	return marketDataRequest
}

func (marketDataRequest *MarketDataRequest) RelatedSymGrp() *RelatedSymGrp {
	group := marketDataRequest.Get(9).(*fix.Group)

	return &RelatedSymGrp{group}
}

func (marketDataRequest *MarketDataRequest) SetRelatedSymGrp(noRelatedSym *RelatedSymGrp) *MarketDataRequest {
	marketDataRequest.Set(9, noRelatedSym.Group)

	return marketDataRequest
}
