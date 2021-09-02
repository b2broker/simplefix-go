package fix44

import (
	"github.com/b2broker/simplefix-go/fix"
	"github.com/b2broker/simplefix-go/session/messages"
)

const MsgTypeMarketDataSnapshotFullRefresh = "W"

type MarketDataSnapshotFullRefresh struct {
	*fix.Message
}

func makeMarketDataSnapshotFullRefresh() *MarketDataSnapshotFullRefresh {
	msg := &MarketDataSnapshotFullRefresh{
		Message: fix.NewMessage(FieldBeginString, FieldBodyLength, FieldCheckSum, FieldMsgType, beginString, MsgTypeMarketDataSnapshotFullRefresh).
			SetBody(
				fix.NewKeyValue(FieldMDReqID, &fix.String{}),
				makeInstrument().Component,
				NewUnderlyingsGrp().Group,
				NewLegsGrp().Group,
				fix.NewKeyValue(FieldFinancialStatus, &fix.String{}),
				fix.NewKeyValue(FieldCorporateAction, &fix.String{}),
				fix.NewKeyValue(FieldNetChgPrevDay, &fix.Float{}),
				NewMDEntriesGrp().Group,
				fix.NewKeyValue(FieldApplQueueDepth, &fix.Int{}),
				fix.NewKeyValue(FieldApplQueueResolution, &fix.String{}),
			),
	}

	msg.SetHeader(makeHeader().AsComponent())
	msg.SetTrailer(makeTrailer().AsComponent())

	return msg
}

func NewMarketDataSnapshotFullRefresh(instrument *Instrument, noMDEntries *MDEntriesGrp) *MarketDataSnapshotFullRefresh {
	msg := makeMarketDataSnapshotFullRefresh().
		SetInstrument(instrument).
		SetMDEntriesGrp(noMDEntries)

	return msg
}

func ParseMarketDataSnapshotFullRefresh(data []byte) (*MarketDataSnapshotFullRefresh, error) {
	msg := fix.NewMessage(FieldBeginString, FieldBodyLength, FieldCheckSum, FieldMsgType, beginString, MsgTypeMarketDataSnapshotFullRefresh).
		SetBody(makeMarketDataSnapshotFullRefresh().Body()...).
		SetHeader(makeHeader().AsComponent()).
		SetTrailer(makeTrailer().AsComponent())

	if err := msg.Unmarshal(data); err != nil {
		return nil, err
	}

	return &MarketDataSnapshotFullRefresh{
		Message: msg,
	}, nil
}

func (marketDataSnapshotFullRefresh *MarketDataSnapshotFullRefresh) Header() *Header {
	header := marketDataSnapshotFullRefresh.Message.Header()

	return &Header{header}
}

func (marketDataSnapshotFullRefresh *MarketDataSnapshotFullRefresh) HeaderBuilder() messages.HeaderBuilder {
	return marketDataSnapshotFullRefresh.Header()
}

func (marketDataSnapshotFullRefresh *MarketDataSnapshotFullRefresh) Trailer() *Trailer {
	trailer := marketDataSnapshotFullRefresh.Message.Trailer()

	return &Trailer{trailer}
}

func (marketDataSnapshotFullRefresh *MarketDataSnapshotFullRefresh) MDReqID() string {
	kv := marketDataSnapshotFullRefresh.Get(0)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (marketDataSnapshotFullRefresh *MarketDataSnapshotFullRefresh) SetMDReqID(mDReqID string) *MarketDataSnapshotFullRefresh {
	kv := marketDataSnapshotFullRefresh.Get(0).(*fix.KeyValue)
	_ = kv.Load().Set(mDReqID)
	return marketDataSnapshotFullRefresh
}

func (marketDataSnapshotFullRefresh *MarketDataSnapshotFullRefresh) Instrument() *Instrument {
	component := marketDataSnapshotFullRefresh.Get(1).(*fix.Component)

	return &Instrument{component}
}

func (marketDataSnapshotFullRefresh *MarketDataSnapshotFullRefresh) SetInstrument(instrument *Instrument) *MarketDataSnapshotFullRefresh {
	marketDataSnapshotFullRefresh.Set(1, instrument.Component)

	return marketDataSnapshotFullRefresh
}

func (marketDataSnapshotFullRefresh *MarketDataSnapshotFullRefresh) UnderlyingsGrp() *UnderlyingsGrp {
	group := marketDataSnapshotFullRefresh.Get(2).(*fix.Group)

	return &UnderlyingsGrp{group}
}

func (marketDataSnapshotFullRefresh *MarketDataSnapshotFullRefresh) SetUnderlyingsGrp(noUnderlyings *UnderlyingsGrp) *MarketDataSnapshotFullRefresh {
	marketDataSnapshotFullRefresh.Set(2, noUnderlyings.Group)

	return marketDataSnapshotFullRefresh
}

func (marketDataSnapshotFullRefresh *MarketDataSnapshotFullRefresh) LegsGrp() *LegsGrp {
	group := marketDataSnapshotFullRefresh.Get(3).(*fix.Group)

	return &LegsGrp{group}
}

func (marketDataSnapshotFullRefresh *MarketDataSnapshotFullRefresh) SetLegsGrp(noLegs *LegsGrp) *MarketDataSnapshotFullRefresh {
	marketDataSnapshotFullRefresh.Set(3, noLegs.Group)

	return marketDataSnapshotFullRefresh
}

func (marketDataSnapshotFullRefresh *MarketDataSnapshotFullRefresh) FinancialStatus() string {
	kv := marketDataSnapshotFullRefresh.Get(4)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (marketDataSnapshotFullRefresh *MarketDataSnapshotFullRefresh) SetFinancialStatus(financialStatus string) *MarketDataSnapshotFullRefresh {
	kv := marketDataSnapshotFullRefresh.Get(4).(*fix.KeyValue)
	_ = kv.Load().Set(financialStatus)
	return marketDataSnapshotFullRefresh
}

func (marketDataSnapshotFullRefresh *MarketDataSnapshotFullRefresh) CorporateAction() string {
	kv := marketDataSnapshotFullRefresh.Get(5)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (marketDataSnapshotFullRefresh *MarketDataSnapshotFullRefresh) SetCorporateAction(corporateAction string) *MarketDataSnapshotFullRefresh {
	kv := marketDataSnapshotFullRefresh.Get(5).(*fix.KeyValue)
	_ = kv.Load().Set(corporateAction)
	return marketDataSnapshotFullRefresh
}

func (marketDataSnapshotFullRefresh *MarketDataSnapshotFullRefresh) NetChgPrevDay() float64 {
	kv := marketDataSnapshotFullRefresh.Get(6)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(float64)
}

func (marketDataSnapshotFullRefresh *MarketDataSnapshotFullRefresh) SetNetChgPrevDay(netChgPrevDay float64) *MarketDataSnapshotFullRefresh {
	kv := marketDataSnapshotFullRefresh.Get(6).(*fix.KeyValue)
	_ = kv.Load().Set(netChgPrevDay)
	return marketDataSnapshotFullRefresh
}

func (marketDataSnapshotFullRefresh *MarketDataSnapshotFullRefresh) MDEntriesGrp() *MDEntriesGrp {
	group := marketDataSnapshotFullRefresh.Get(7).(*fix.Group)

	return &MDEntriesGrp{group}
}

func (marketDataSnapshotFullRefresh *MarketDataSnapshotFullRefresh) SetMDEntriesGrp(noMDEntries *MDEntriesGrp) *MarketDataSnapshotFullRefresh {
	marketDataSnapshotFullRefresh.Set(7, noMDEntries.Group)

	return marketDataSnapshotFullRefresh
}

func (marketDataSnapshotFullRefresh *MarketDataSnapshotFullRefresh) ApplQueueDepth() int {
	kv := marketDataSnapshotFullRefresh.Get(8)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(int)
}

func (marketDataSnapshotFullRefresh *MarketDataSnapshotFullRefresh) SetApplQueueDepth(applQueueDepth int) *MarketDataSnapshotFullRefresh {
	kv := marketDataSnapshotFullRefresh.Get(8).(*fix.KeyValue)
	_ = kv.Load().Set(applQueueDepth)
	return marketDataSnapshotFullRefresh
}

func (marketDataSnapshotFullRefresh *MarketDataSnapshotFullRefresh) ApplQueueResolution() string {
	kv := marketDataSnapshotFullRefresh.Get(9)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (marketDataSnapshotFullRefresh *MarketDataSnapshotFullRefresh) SetApplQueueResolution(applQueueResolution string) *MarketDataSnapshotFullRefresh {
	kv := marketDataSnapshotFullRefresh.Get(9).(*fix.KeyValue)
	_ = kv.Load().Set(applQueueResolution)
	return marketDataSnapshotFullRefresh
}
