package fix44

import (
	"github.com/b2broker/simplefix-go/fix"
)

type AltMDSourceGrp struct {
	*fix.Group
}

func NewAltMDSourceGrp() *AltMDSourceGrp {
	return &AltMDSourceGrp{
		fix.NewGroup(FieldNoAltMDSource,
			fix.NewKeyValue(FieldAltMDSourceID, &fix.String{}),
		),
	}
}

func (group *AltMDSourceGrp) AddEntry(entry *AltMDSourceEntry) *AltMDSourceGrp {
	group.Group.AddEntry(entry.Items())

	return group
}

type AltMDSourceEntry struct {
	*fix.Component
}

func makeAltMDSourceEntry() *AltMDSourceEntry {
	return &AltMDSourceEntry{fix.NewComponent(
		fix.NewKeyValue(FieldAltMDSourceID, &fix.String{}),
	)}
}

func NewAltMDSourceEntry() *AltMDSourceEntry {
	return makeAltMDSourceEntry()
}

func (altMDSourceEntry *AltMDSourceEntry) AltMDSourceID() string {
	kv := altMDSourceEntry.Get(0)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (altMDSourceEntry *AltMDSourceEntry) SetAltMDSourceID(altMDSourceID string) *AltMDSourceEntry {
	kv := altMDSourceEntry.Get(0).(*fix.KeyValue)
	_ = kv.Load().Set(altMDSourceID)
	return altMDSourceEntry
}
