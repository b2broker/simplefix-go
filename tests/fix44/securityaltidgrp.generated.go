package fix44

import (
	"github.com/b2broker/simplefix-go/fix"
)

type SecurityAltIDGrp struct {
	*fix.Group
}

func NewSecurityAltIDGrp() *SecurityAltIDGrp {
	return &SecurityAltIDGrp{
		fix.NewGroup(FieldNoSecurityAltID,
			fix.NewKeyValue(FieldSecurityAltID, &fix.String{}),
			fix.NewKeyValue(FieldSecurityAltIDSource, &fix.String{}),
		),
	}
}

func (group *SecurityAltIDGrp) AddEntry(entry *SecurityAltIDEntry) *SecurityAltIDGrp {
	group.Group.AddEntry(entry.Items())

	return group
}

type SecurityAltIDEntry struct {
	*fix.Component
}

func makeSecurityAltIDEntry() *SecurityAltIDEntry {
	return &SecurityAltIDEntry{fix.NewComponent(
		fix.NewKeyValue(FieldSecurityAltID, &fix.String{}),
		fix.NewKeyValue(FieldSecurityAltIDSource, &fix.String{}),
	)}
}

func NewSecurityAltIDEntry() *SecurityAltIDEntry {
	return makeSecurityAltIDEntry()
}

func (securityAltIDEntry *SecurityAltIDEntry) SecurityAltID() string {
	kv := securityAltIDEntry.Get(0)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (securityAltIDEntry *SecurityAltIDEntry) SetSecurityAltID(securityAltID string) *SecurityAltIDEntry {
	kv := securityAltIDEntry.Get(0).(*fix.KeyValue)
	_ = kv.Load().Set(securityAltID)
	return securityAltIDEntry
}

func (securityAltIDEntry *SecurityAltIDEntry) SecurityAltIDSource() string {
	kv := securityAltIDEntry.Get(1)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (securityAltIDEntry *SecurityAltIDEntry) SetSecurityAltIDSource(securityAltIDSource string) *SecurityAltIDEntry {
	kv := securityAltIDEntry.Get(1).(*fix.KeyValue)
	_ = kv.Load().Set(securityAltIDSource)
	return securityAltIDEntry
}
