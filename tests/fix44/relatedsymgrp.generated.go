package fix44

import (
	"github.com/b2broker/simplefix-go/fix"
)

type RelatedSymGrp struct {
	*fix.Group
}

func NewRelatedSymGrp() *RelatedSymGrp {
	return &RelatedSymGrp{
		fix.NewGroup(FieldNoRelatedSym,
			makeInstrument().Component,
			NewUnderlyingsGrp().Group,
			NewLegsGrp().Group,
			NewTradingSessionsGrp().Group,
			fix.NewKeyValue(FieldApplQueueAction, &fix.String{}),
			fix.NewKeyValue(FieldApplQueueMax, &fix.Int{}),
		),
	}
}

func (group *RelatedSymGrp) AddEntry(entry *RelatedSymEntry) *RelatedSymGrp {
	group.Group.AddEntry(entry.Items())

	return group
}

type RelatedSymEntry struct {
	*fix.Component
}

func makeRelatedSymEntry() *RelatedSymEntry {
	return &RelatedSymEntry{fix.NewComponent(
		makeInstrument().Component,
		NewUnderlyingsGrp().Group,
		NewLegsGrp().Group,
		NewTradingSessionsGrp().Group,
		fix.NewKeyValue(FieldApplQueueAction, &fix.String{}),
		fix.NewKeyValue(FieldApplQueueMax, &fix.Int{}),
	)}
}

func NewRelatedSymEntry() *RelatedSymEntry {
	return makeRelatedSymEntry()
}

func (relatedSymEntry *RelatedSymEntry) Instrument() *Instrument {
	component := relatedSymEntry.Get(0).(*fix.Component)

	return &Instrument{component}
}

func (relatedSymEntry *RelatedSymEntry) SetInstrument(instrument *Instrument) *RelatedSymEntry {
	relatedSymEntry.Set(0, instrument.Component)

	return relatedSymEntry
}

func (relatedSymEntry *RelatedSymEntry) UnderlyingsGrp() *UnderlyingsGrp {
	group := relatedSymEntry.Get(1).(*fix.Group)

	return &UnderlyingsGrp{group}
}

func (relatedSymEntry *RelatedSymEntry) SetUnderlyingsGrp(noUnderlyings *UnderlyingsGrp) *RelatedSymEntry {
	relatedSymEntry.Set(1, noUnderlyings.Group)

	return relatedSymEntry
}

func (relatedSymEntry *RelatedSymEntry) LegsGrp() *LegsGrp {
	group := relatedSymEntry.Get(2).(*fix.Group)

	return &LegsGrp{group}
}

func (relatedSymEntry *RelatedSymEntry) SetLegsGrp(noLegs *LegsGrp) *RelatedSymEntry {
	relatedSymEntry.Set(2, noLegs.Group)

	return relatedSymEntry
}

func (relatedSymEntry *RelatedSymEntry) TradingSessionsGrp() *TradingSessionsGrp {
	group := relatedSymEntry.Get(3).(*fix.Group)

	return &TradingSessionsGrp{group}
}

func (relatedSymEntry *RelatedSymEntry) SetTradingSessionsGrp(noTradingSessions *TradingSessionsGrp) *RelatedSymEntry {
	relatedSymEntry.Set(3, noTradingSessions.Group)

	return relatedSymEntry
}

func (relatedSymEntry *RelatedSymEntry) ApplQueueAction() string {
	kv := relatedSymEntry.Get(4)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (relatedSymEntry *RelatedSymEntry) SetApplQueueAction(applQueueAction string) *RelatedSymEntry {
	kv := relatedSymEntry.Get(4).(*fix.KeyValue)
	_ = kv.Load().Set(applQueueAction)
	return relatedSymEntry
}

func (relatedSymEntry *RelatedSymEntry) ApplQueueMax() int {
	kv := relatedSymEntry.Get(5)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(int)
}

func (relatedSymEntry *RelatedSymEntry) SetApplQueueMax(applQueueMax int) *RelatedSymEntry {
	kv := relatedSymEntry.Get(5).(*fix.KeyValue)
	_ = kv.Load().Set(applQueueMax)
	return relatedSymEntry
}
