package fix44

import (
	"github.com/b2broker/simplefix-go/fix"
)

type TradingSessionsGrp struct {
	*fix.Group
}

func NewTradingSessionsGrp() *TradingSessionsGrp {
	return &TradingSessionsGrp{
		fix.NewGroup(FieldNoTradingSessions,
			fix.NewKeyValue(FieldTradingSessionID, &fix.String{}),
			fix.NewKeyValue(FieldTradingSessionSubID, &fix.String{}),
		),
	}
}

func (group *TradingSessionsGrp) AddEntry(entry *TradingSessionsEntry) *TradingSessionsGrp {
	group.Group.AddEntry(entry.Items())

	return group
}

func (group *TradingSessionsGrp) Entries() []*TradingSessionsEntry {
	items := make([]*TradingSessionsEntry, len(group.Group.Entries()))

	for i, item := range group.Group.Entries() {
		items[i] = &TradingSessionsEntry{fix.NewComponent(item...)}
	}

	return items
}

type TradingSessionsEntry struct {
	*fix.Component
}

func makeTradingSessionsEntry() *TradingSessionsEntry {
	return &TradingSessionsEntry{fix.NewComponent(
		fix.NewKeyValue(FieldTradingSessionID, &fix.String{}),
		fix.NewKeyValue(FieldTradingSessionSubID, &fix.String{}),
	)}
}

func NewTradingSessionsEntry() *TradingSessionsEntry {
	return makeTradingSessionsEntry()
}

func (tradingSessionsEntry *TradingSessionsEntry) TradingSessionID() string {
	kv := tradingSessionsEntry.Get(0)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (tradingSessionsEntry *TradingSessionsEntry) SetTradingSessionID(tradingSessionID string) *TradingSessionsEntry {
	kv := tradingSessionsEntry.Get(0).(*fix.KeyValue)
	_ = kv.Load().Set(tradingSessionID)
	return tradingSessionsEntry
}

func (tradingSessionsEntry *TradingSessionsEntry) TradingSessionSubID() string {
	kv := tradingSessionsEntry.Get(1)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (tradingSessionsEntry *TradingSessionsEntry) SetTradingSessionSubID(tradingSessionSubID string) *TradingSessionsEntry {
	kv := tradingSessionsEntry.Get(1).(*fix.KeyValue)
	_ = kv.Load().Set(tradingSessionSubID)
	return tradingSessionsEntry
}
