package fix44

import (
	"github.com/b2broker/simplefix-go/fix"
)

type UnderlyingsGrp struct {
	*fix.Group
}

func NewUnderlyingsGrp() *UnderlyingsGrp {
	return &UnderlyingsGrp{
		fix.NewGroup(FieldNoUnderlyings,
			makeUnderlyingInstrument().Component,
		),
	}
}

func (group *UnderlyingsGrp) AddEntry(entry *UnderlyingsEntry) *UnderlyingsGrp {
	group.Group.AddEntry(entry.Items())

	return group
}

type UnderlyingsEntry struct {
	*fix.Component
}

func makeUnderlyingsEntry() *UnderlyingsEntry {
	return &UnderlyingsEntry{fix.NewComponent(
		makeUnderlyingInstrument().Component,
	)}
}

func NewUnderlyingsEntry() *UnderlyingsEntry {
	return makeUnderlyingsEntry()
}

func (underlyingsEntry *UnderlyingsEntry) UnderlyingInstrument() *UnderlyingInstrument {
	component := underlyingsEntry.Get(0).(*fix.Component)

	return &UnderlyingInstrument{component}
}

func (underlyingsEntry *UnderlyingsEntry) SetUnderlyingInstrument(underlyingInstrument *UnderlyingInstrument) *UnderlyingsEntry {
	underlyingsEntry.Set(0, underlyingInstrument.Component)

	return underlyingsEntry
}
