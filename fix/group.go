package fix

import (
	"fmt"
	"strings"
)

type Group struct {
	noTag    string
	template Items

	items []Items
}

func (g *Group) NoTag() string {
	return g.noTag
}

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

func NewGroup(noTags string, tags ...Item) *Group {
	return &Group{
		noTag:    noTags,
		template: tags,
	}
}

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

func (g *Group) AddEntry(v Items) *Group {
	g.items = append(g.items, v)

	return g
}

func (g *Group) Entry(id int) Item {
	return g.items[id]
}

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

func (g *Group) Entries() []Items {
	return g.items
}
