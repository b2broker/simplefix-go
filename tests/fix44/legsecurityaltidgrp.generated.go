package fix44

import (
	"github.com/b2broker/simplefix-go/fix"
)

type LegSecurityAltIDGrp struct {
	*fix.Group
}

func NewLegSecurityAltIDGrp() *LegSecurityAltIDGrp {
	return &LegSecurityAltIDGrp{
		fix.NewGroup(FieldNoLegSecurityAltID,
			fix.NewKeyValue(FieldLegSecurityAltID, &fix.String{}),
			fix.NewKeyValue(FieldLegSecurityAltIDSource, &fix.String{}),
		),
	}
}

func (group *LegSecurityAltIDGrp) AddEntry(entry *LegSecurityAltIDEntry) *LegSecurityAltIDGrp {
	group.Group.AddEntry(entry.Items())

	return group
}

type LegSecurityAltIDEntry struct {
	*fix.Component
}

func makeLegSecurityAltIDEntry() *LegSecurityAltIDEntry {
	return &LegSecurityAltIDEntry{fix.NewComponent(
		fix.NewKeyValue(FieldLegSecurityAltID, &fix.String{}),
		fix.NewKeyValue(FieldLegSecurityAltIDSource, &fix.String{}),
	)}
}

func NewLegSecurityAltIDEntry() *LegSecurityAltIDEntry {
	return makeLegSecurityAltIDEntry()
}

func (legSecurityAltIDEntry *LegSecurityAltIDEntry) LegSecurityAltID() string {
	kv := legSecurityAltIDEntry.Get(0)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (legSecurityAltIDEntry *LegSecurityAltIDEntry) SetLegSecurityAltID(legSecurityAltID string) *LegSecurityAltIDEntry {
	kv := legSecurityAltIDEntry.Get(0).(*fix.KeyValue)
	_ = kv.Load().Set(legSecurityAltID)
	return legSecurityAltIDEntry
}

func (legSecurityAltIDEntry *LegSecurityAltIDEntry) LegSecurityAltIDSource() string {
	kv := legSecurityAltIDEntry.Get(1)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (legSecurityAltIDEntry *LegSecurityAltIDEntry) SetLegSecurityAltIDSource(legSecurityAltIDSource string) *LegSecurityAltIDEntry {
	kv := legSecurityAltIDEntry.Get(1).(*fix.KeyValue)
	_ = kv.Load().Set(legSecurityAltIDSource)
	return legSecurityAltIDEntry
}
