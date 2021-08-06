package fix

import (
	"fmt"
	"strings"
)

// Group is a structure for FIX-Group types
type Group struct {
	noTag    string
	template Items

	items []Items
}

// NoTag returns tag with number of elements
func (g *Group) NoTag() string {
	return g.noTag
}

// String converts Group to string
func (g *Group) String() string {
	var items []string
	for _, item := range g.items {
		itemStr := item.String()
		if itemStr != "" {
			items = append(items, itemStr)
		}
	}
	return fmt.Sprintf("[%s]", strings.Join(items, ", "))
}

// NewGroup create new group from tag with number of elements and list of tags
func NewGroup(noTags string, tags ...Item) *Group {
	return &Group{
		noTag:    noTags,
		template: tags,
	}
}

// ToBytes convert Group to bytes
func (g *Group) ToBytes() []byte {
	var msg [][]byte

	if len(g.items) == 0 {
		return nil
	}

	msg = append(msg, NewKeyValue(
		g.noTag,
		NewInt(len(g.items)),
	).ToBytes())

	for _, item := range g.items {
		itemB := item.ToBytes()
		if itemB != nil {
			msg = append(msg, item.ToBytes())
		}
	}
	return joinBody(msg...)
}

// AddEntry add entry with same list of tags
// you can receive all tags by AsTemplate method
func (g *Group) AddEntry(v Items) *Group {
	g.items = append(g.items, v)

	return g
}

// Entry returns entry of group by sequence number
func (g *Group) Entry(id int) Item {
	return g.items[id]
}

// AsTemplate returns list of group tags as Items
func (g *Group) AsTemplate() Items {
	tmp := make([]Item, len(g.template))

	for i, item := range g.template {
		switch value := item.(type) {
		case *KeyValue:
			tmp[i] = value.AsTemplate()

		case *Group:
			tmp[i] = NewGroup(value.NoTag(), value.AsTemplate()...)

		case *Component:
			tmp[i] = NewComponent(value.AsTemplate()...)
		}
	}

	return tmp
}

// Entries returns all entries of group as list of Items
func (g *Group) Entries() []Items {
	return g.items
}
