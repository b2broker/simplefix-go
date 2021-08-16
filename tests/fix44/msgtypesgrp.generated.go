package fix44

import (
	"github.com/b2broker/simplefix-go/fix"
)

type MsgTypesGrp struct {
	*fix.Group
}

func NewMsgTypesGrp() *MsgTypesGrp {
	return &MsgTypesGrp{
		fix.NewGroup(FieldNoMsgTypes,
			fix.NewKeyValue(FieldRefMsgType, &fix.String{}),
			fix.NewKeyValue(FieldMsgDirection, &fix.String{}),
		),
	}
}

func (group *MsgTypesGrp) AddEntry(entry *MsgTypesEntry) *MsgTypesGrp {
	group.Group.AddEntry(entry.Items())

	return group
}

func (group *MsgTypesGrp) Entries() []*MsgTypesEntry {
	items := make([]*MsgTypesEntry, len(group.Group.Entries()))

	for i, item := range group.Group.Entries() {
		items[i] = &MsgTypesEntry{fix.NewComponent(item...)}
	}

	return items
}

type MsgTypesEntry struct {
	*fix.Component
}

func makeMsgTypesEntry() *MsgTypesEntry {
	return &MsgTypesEntry{fix.NewComponent(
		fix.NewKeyValue(FieldRefMsgType, &fix.String{}),
		fix.NewKeyValue(FieldMsgDirection, &fix.String{}),
	)}
}

func NewMsgTypesEntry() *MsgTypesEntry {
	return makeMsgTypesEntry()
}

func (msgTypesEntry *MsgTypesEntry) RefMsgType() string {
	kv := msgTypesEntry.Get(0)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (msgTypesEntry *MsgTypesEntry) SetRefMsgType(refMsgType string) *MsgTypesEntry {
	kv := msgTypesEntry.Get(0).(*fix.KeyValue)
	_ = kv.Load().Set(refMsgType)
	return msgTypesEntry
}

func (msgTypesEntry *MsgTypesEntry) MsgDirection() string {
	kv := msgTypesEntry.Get(1)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (msgTypesEntry *MsgTypesEntry) SetMsgDirection(msgDirection string) *MsgTypesEntry {
	kv := msgTypesEntry.Get(1).(*fix.KeyValue)
	_ = kv.Load().Set(msgDirection)
	return msgTypesEntry
}
