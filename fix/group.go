package fix

import (
	"bytes"
	"fmt"
	"strings"
)

// Group is a structure used to implement FIX group types.
type Group struct {
	noTag    string
	template Items

	items []Items
}

// NoTag returns a tag value indicating the number of elements in a group.
func (g *Group) NoTag() string {
	return g.noTag
}

// String returns a string representation of a Group.
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

// NewGroup is used to create a new group based on:
// - the tag value specifying the number of elements, and
// - the list of tags
func NewGroup(noTags string, tags ...Item) *Group {
	return &Group{
		noTag:    noTags,
		template: tags,
	}
}

// ToBytes returns a byte representation of a Group.
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

// always have size tag
func (g *Group) IsEmpty() bool {
	return false
}

func (g *Group) WriteBytes(writer *bytes.Buffer) bool {

	if len(g.items) == 0 {
		return false
	}
	intVal := NewInt(len(g.items))
	NewKeyValue(g.noTag, intVal).WriteBytes(writer)

	_ = writer.WriteByte(DelimiterChar)

	addDelimeter := false
	for i, item := range g.items {
		if !item.IsEmpty() && addDelimeter {
			_ = writer.WriteByte(DelimiterChar)
			addDelimeter = false
		}
		if item.WriteBytes(writer) {
			if i <= len(g.items)-1 {
				addDelimeter = true
			}
		}
	}
	return true
}

// AddEntry adds a new entry with the same list of tags as in a specified group.
// To generate all tags for a newly created entry, use the AsTemplate method.
func (g *Group) AddEntry(v Items) *Group {
	g.items = append(g.items, v)

	return g
}

// Entry returns a group entry by its sequence number.
func (g *Group) Entry(id int) Item {
	return g.items[id]
}

// AsTemplate returns a list of group tags as an Items object.
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

// Entries returns all entries belonging to a group as an array of Items objects.
func (g *Group) Entries() []Items {
	return g.items
}
