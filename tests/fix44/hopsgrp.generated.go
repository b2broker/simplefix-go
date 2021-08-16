package fix44

import (
	"github.com/b2broker/simplefix-go/fix"
)

type HopsGrp struct {
	*fix.Group
}

func NewHopsGrp() *HopsGrp {
	return &HopsGrp{
		fix.NewGroup(FieldNoHops,
			fix.NewKeyValue(FieldHopCompID, &fix.String{}),
			fix.NewKeyValue(FieldHopSendingTime, &fix.String{}),
			fix.NewKeyValue(FieldHopRefID, &fix.Int{}),
		),
	}
}

func (group *HopsGrp) AddEntry(entry *HopsEntry) *HopsGrp {
	group.Group.AddEntry(entry.Items())

	return group
}

func (group *HopsGrp) Entries() []*HopsEntry {
	items := make([]*HopsEntry, len(group.Group.Entries()))

	for i, item := range group.Group.Entries() {
		items[i] = &HopsEntry{fix.NewComponent(item...)}
	}

	return items
}

type HopsEntry struct {
	*fix.Component
}

func makeHopsEntry() *HopsEntry {
	return &HopsEntry{fix.NewComponent(
		fix.NewKeyValue(FieldHopCompID, &fix.String{}),
		fix.NewKeyValue(FieldHopSendingTime, &fix.String{}),
		fix.NewKeyValue(FieldHopRefID, &fix.Int{}),
	)}
}

func NewHopsEntry() *HopsEntry {
	return makeHopsEntry()
}

func (hopsEntry *HopsEntry) HopCompID() string {
	kv := hopsEntry.Get(0)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (hopsEntry *HopsEntry) SetHopCompID(hopCompID string) *HopsEntry {
	kv := hopsEntry.Get(0).(*fix.KeyValue)
	_ = kv.Load().Set(hopCompID)
	return hopsEntry
}

func (hopsEntry *HopsEntry) HopSendingTime() string {
	kv := hopsEntry.Get(1)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (hopsEntry *HopsEntry) SetHopSendingTime(hopSendingTime string) *HopsEntry {
	kv := hopsEntry.Get(1).(*fix.KeyValue)
	_ = kv.Load().Set(hopSendingTime)
	return hopsEntry
}

func (hopsEntry *HopsEntry) HopRefID() int {
	kv := hopsEntry.Get(2)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(int)
}

func (hopsEntry *HopsEntry) SetHopRefID(hopRefID int) *HopsEntry {
	kv := hopsEntry.Get(2).(*fix.KeyValue)
	_ = kv.Load().Set(hopRefID)
	return hopsEntry
}
