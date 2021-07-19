package fix44

import (
	"github.com/b2broker/simplefix-go/fix"
)

type MDEntryTypesGrp struct {
	*fix.Group
}

func NewMDEntryTypesGrp() *MDEntryTypesGrp {
	return &MDEntryTypesGrp{
		fix.NewGroup(FieldNoMDEntryTypes,
			fix.NewKeyValue(FieldMDEntryType, &fix.String{}),
		),
	}
}

func (group *MDEntryTypesGrp) AddEntry(entry *MDEntryTypesEntry) *MDEntryTypesGrp {
	group.Group.AddEntry(entry.Items())

	return group
}

type MDEntryTypesEntry struct {
	*fix.Component
}

func makeMDEntryTypesEntry() *MDEntryTypesEntry {
	return &MDEntryTypesEntry{fix.NewComponent(
		fix.NewKeyValue(FieldMDEntryType, &fix.String{}),
	)}
}

func NewMDEntryTypesEntry() *MDEntryTypesEntry {
	return makeMDEntryTypesEntry()
}

func (mDEntryTypesEntry *MDEntryTypesEntry) MDEntryType() string {
	kv := mDEntryTypesEntry.Get(0)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (mDEntryTypesEntry *MDEntryTypesEntry) SetMDEntryType(mDEntryType string) *MDEntryTypesEntry {
	kv := mDEntryTypesEntry.Get(0).(*fix.KeyValue)
	_ = kv.Load().Set(mDEntryType)
	return mDEntryTypesEntry
}
